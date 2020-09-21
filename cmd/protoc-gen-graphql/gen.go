// Generate `testdata/*.protoset` files for testing
//go:generate protoc -I . -I ../.. --include_source_info --include_imports --descriptor_set_out=./testdata/user.protoset ./testdata/user.proto

package main
