package asar // import "layeh.com/asar"

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"strconv"
)

type entryEncoder struct {
	Contents      []io.Reader
	CurrentOffset int64
	Header        bytes.Buffer
	Encoder       *json.Encoder
}

func (enc *entryEncoder) Write(v interface{}) {
	enc.Encoder.Encode(v)
	enc.Header.Truncate(enc.Header.Len() - 1) // cut off trailing new line
}

func (enc *entryEncoder) WriteField(key string, v interface{}) {
	enc.Write(key)
	enc.Header.WriteByte(':')
	enc.Write(v)
}

func (enc *entryEncoder) Encode(e *Entry) error {
	enc.Header.WriteByte('{')
	if e.Flags&FlagDir != 0 {
		enc.Write("files")
		enc.Header.WriteString(":{")
		for i, child := range e.Children {
			if i > 0 {
				enc.Header.WriteByte(',')
			}
			if !validFilename(child.Name) {
				panic(errHeader)
			}
			enc.Write(child.Name)
			enc.Header.WriteByte(':')
			if err := enc.Encode(child); err != nil {
				return err
			}
		}
		enc.Header.WriteByte('}')
	} else {
		enc.Write("size")
		enc.Header.WriteByte(':')
		enc.Write(e.Size)

		if e.Flags&FlagExecutable != 0 {
			enc.Header.WriteByte(',')
			enc.WriteField("executable", true)
		}

		enc.Header.WriteByte(',')
		if e.Flags&FlagUnpacked == 0 {
			enc.WriteField("offset", strconv.FormatInt(enc.CurrentOffset, 10))
			enc.CurrentOffset += e.Size
			enc.Contents = append(enc.Contents, io.NewSectionReader(e.r, e.baseOffset, e.Size))
		} else {
			enc.WriteField("unpacked", true)
		}
	}
	enc.Header.WriteByte('}')
	return nil
}

// EncodeTo writes an ASAR archive containing Entry's descendants. This function
// is usally called on the root entry.
func (e *Entry) EncodeTo(w io.Writer) (n int64, err error) {

	defer func() {
		if r := recover(); r != nil {
			if e := r.(error); e != nil {
				err = e
			} else {
				panic(r)
			}
		}

	}()

	encoder := entryEncoder{}
	{
		var reserve [16]byte
		encoder.Header.Write(reserve[:])
	}
	encoder.Encoder = json.NewEncoder(&encoder.Header)
	if err = encoder.Encode(e); err != nil {
		return
	}

	{
		var padding [3]byte
		encoder.Header.Write(padding[:encoder.Header.Len()%4])
	}

	header := encoder.Header.Bytes()
	binary.LittleEndian.PutUint32(header[:4], 4)
	binary.LittleEndian.PutUint32(header[4:8], 8+uint32(encoder.Header.Len()))
	binary.LittleEndian.PutUint32(header[8:12], 4+uint32(encoder.Header.Len()))
	binary.LittleEndian.PutUint32(header[12:16], uint32(encoder.Header.Len()))

	n, err = encoder.Header.WriteTo(w)
	if err != nil {
		return
	}

	for _, chunk := range encoder.Contents {
		var written int64
		written, err = io.Copy(w, chunk)
		n += written
		if err != nil {
			return
		}
	}

	return
}
