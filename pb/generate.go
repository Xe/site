package pb

import (
	"embed"
)

func init() {}

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --twirp_out=. --twirp_opt=paths=source_relative xesite.proto
//go:generate go run ../cmd/twirp-openapi-gen --verbose --in=xesite.proto --path-prefix=/api --servers=https://xeiaso.net --title=xeiaso.net --out=openapi.json

//go:embed xesite.proto openapi.json external/*.proto
var Proto embed.FS
