package proto

//go:generate protoc --go_out=identity --go_opt=paths=source_relative --go-grpc_out=identity --go-grpc_opt=paths=source_relative identity.proto
//go:generate protoc --go_out=storage --go_opt=paths=source_relative --go-grpc_out=storage --go-grpc_opt=paths=source_relative storage.proto
