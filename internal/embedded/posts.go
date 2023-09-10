package embedded

import (
	"bytes"
	_ "embed"
	"encoding/json"

	"xeiaso.net/v4/internal"
)

//go:generate go run generate.go

var (
	//go:embed posts.json
	postJSON []byte

	Posts []*internal.Post
)

func init() {
	if err := json.NewDecoder(bytes.NewReader(postJSON)).Decode(&Posts); err != nil {
		panic(err)
	}
}
