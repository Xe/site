package embedded

import (
	"bytes"
	_ "embed"
	"encoding/gob"

	"xeiaso.net/v4/internal"
)

//go:generate go run generate.go

var (
	//go:embed posts.gob
	postGob []byte

	Posts []*internal.Post
)

func init() {
	if err := gob.NewDecoder(bytes.NewReader(postGob)).Decode(&Posts); err != nil {
		panic(err)
	}
}
