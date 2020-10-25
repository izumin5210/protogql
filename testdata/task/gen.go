//go:generate sh -c "protoc -I proto -I ../.. --go_out=paths=source_relative:api proto/*.proto"
//go:generate sh -c "protoc -I proto -I ../.. --plugin=protoc-gen-graphql='./bin/protoc-gen-graphql' --graphql_out=paths=source_relative:graph proto/*.proto"

package main
