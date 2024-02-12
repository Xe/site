package adminpb

func init() {}

//go:generate protoc --proto_path=. --proto_path=../../pb --go_out=. --go_opt=paths=source_relative --twirp_out=. --twirp_opt=paths=source_relative internal.proto
