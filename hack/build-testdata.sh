#!/usr/bin/env bash


set -eu
set -o pipefail

cd $(dirname $0)/..

pushd() {
    command pushd "$@" > /dev/null
}

popd() {
    command popd "$@" > /dev/null
}

clean() {
  rm -rf bin
  rm -rf testdata/apis/go
  rm -rf testdata/apis/graphql
}

init() {
  mkdir testdata/apis/go
  mkdir testdata/apis/graphql

  go generate ./tools.go
  go build -o ./bin/protoc-gen-graphql ./cmd/protoc-gen-graphql
}

execProtoc() {
  for protoDir in ./testdata/apis/proto/*; do
    protoc -I testdata/apis/proto -I . --include_source_info --include_imports --descriptor_set_out=${protoDir}/descriptor_set.pb ${protoDir}/*.proto
    protoc -I testdata/apis/proto -I . --plugin=protoc-gen-graphql='./bin/protoc-gen-graphql' --graphql_out=testdata/apis/graphql  ${protoDir}/*.proto
    protoc -I testdata/apis/proto -I . --go_out=testdata ${protoDir}/*.proto
  done
}

initGoMod() {
  for pbgoDir in ./testdata/apis/go/*; do
    pushd $pbgoDir
    go mod init apis/go/$(basename $pbgoDir)
    echo 'replace github.com/izumin5210/remixer => ../../../..' >> go.mod
    go mod tidy
    popd
  done
}

_main() {
  clean
  init
  execProtoc
  initGoMod
}

_main
