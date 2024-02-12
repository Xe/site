package pb

func init() {}

//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --twirp_out=. --twirp_opt=paths=source_relative xesite.proto
