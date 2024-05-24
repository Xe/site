package protofeed

//go:generate protoc --proto_path=. --proto_path=.. --go_out=./protofeed --go_opt=paths=source_relative ./protofeed.proto
//go:generate protoc --proto_path=. --proto_path=.. --go_out=./mimi/announce --go_opt=paths=source_relative --twirp_out=./mimi/announce --twirp_opt=paths=source_relative ./mimi-announce.proto
//go:generate protoc --proto_path=. --proto_path=.. --go_out=./mi --go_opt=paths=source_relative --twirp_out=./mi --twirp_opt=paths=source_relative ./mi.proto
