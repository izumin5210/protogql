// Generate `testdata/*.protoset` files for testing
//go:generate sh -c "protoc -I . -I ../.. --include_source_info --include_imports --descriptor_set_out=./testdata/user.protoset ./testdata/user.proto"
//go:generate sh -c "protoc -I . -I ../.. --include_source_info --include_imports --descriptor_set_out=./testdata/starwars.protoset ./testdata/starwars/*.proto"

package main
