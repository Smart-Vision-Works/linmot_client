//go:build tools
// +build tools

package proto

//go:generate protoc --experimental_allow_proto3_optional --proto_path=../../proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true stage_primer.proto
