package asar // import "layeh.com/asar"

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	errHeader = errors.New("asar: invalid file header")
)

type jsonReader struct {
	ASAR       io.ReaderAt
	BaseOffset int64
	D          *json.Decoder
	Token      json.Token
}

func (j *jsonReader) Peek() json.Token {
	if j.Token != nil {
		return j.Token
	}
	tkn, err := j.D.Token()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}
	j.Token = tkn
	return tkn
}

func (j *jsonReader) HasDelimRune(r rune) bool {
	peek := j.Peek()
	ru, ok := peek.(json.Delim)
	if !ok {
		return false
	}
	if rune(ru) != r {
		return false
	}
	return true
}

func (j *jsonReader) Next() json.Token {
	if j.Token != nil {
		t := j.Token
		j.Token = nil
		return t
	}

	tkn, err := j.D.Token()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}
	return tkn
}

func (j *jsonReader) NextDelimRune() rune {
	tkn := j.Next()
	r, ok := tkn.(json.Delim)
	if !ok {
		panic(errHeader)
	}
	return rune(r)
}

func (j *jsonReader) ExpectDelim(r rune) {
	next := j.NextDelimRune()
	if next != r {
		panic(errHeader)
	}
}

func (j *jsonReader) ExpectBool() bool {
	tkn := j.Next()
	b, ok := tkn.(bool)
	if !ok {
		panic(errHeader)
	}
	return b
}

func (j *jsonReader) ExpectString() string {
	next := j.Next()
	str, ok := next.(string)
	if !ok {
		panic(errHeader)
	}
	return str
}

func (j *jsonReader) ExpectStringVal(val string) {
	str := j.ExpectString()
	if str != val {
		panic(errHeader)
	}
}

func (j *jsonReader) ExpectInt64() int64 {
	var number json.Number
	switch j.Peek().(type) {
	case string:
		number = json.Number(j.ExpectString())
	case json.Number:
		number = j.Next().(json.Number)
	default:
		panic(errHeader)
	}
	val, err := number.Int64()
	if err != nil {
		panic(errHeader)
	}
	return val
}

func parseRoot(r *jsonReader) *Entry {
	entry := &Entry{
		Flags: FlagDir,
	}

	r.ExpectDelim('{')
	r.ExpectStringVal("files")
	parseFiles(r, entry)
	r.ExpectDelim('}')
	if r.Next() != nil {
		panic(errHeader)
	}
	return entry
}

func parseFiles(r *jsonReader, parent *Entry) {
	r.ExpectDelim('{')
	for !r.HasDelimRune('}') {
		parseEntry(r, parent)
	}
	r.ExpectDelim('}')
}

func parseEntry(r *jsonReader, parent *Entry) {
	name := r.ExpectString()
	if name == "" {
		panic(errHeader)
	}
	if !validFilename(name) {
		panic(errHeader)
	}

	r.ExpectDelim('{')

	child := &Entry{
		Name:   name,
		Parent: parent,
	}

	for !r.HasDelimRune('}') {
		switch r.ExpectString() {
		case "files":
			child.Flags |= FlagDir
			parseFiles(r, child)
		case "size":
			child.Size = r.ExpectInt64()
		case "offset":
			child.Offset = r.ExpectInt64()
		case "unpacked":
			if r.ExpectBool() {
				child.Flags |= FlagUnpacked
			}
		case "executable":
			if r.ExpectBool() {
				child.Flags |= FlagExecutable
			}
		default:
			panic(errHeader)
		}
	}

	if child.Flags&FlagDir == 0 {
		child.r = r.ASAR
		child.baseOffset = r.BaseOffset
	}

	parent.Children = append(parent.Children, child)
	r.ExpectDelim('}')
}

func decodeHeader(asar io.ReaderAt, header *io.SectionReader, offset int64) (entry *Entry, err error) {
	decoder := json.NewDecoder(header)
	decoder.UseNumber()
	reader := jsonReader{
		ASAR:       asar,
		BaseOffset: offset,
		D:          decoder,
	}
	defer func() {
		if r := recover(); r != nil {
			if e := r.(error); e != nil {
				err = e
			} else {
				panic(r)
			}
		}

	}()
	entry = parseRoot(&reader)
	return
}
