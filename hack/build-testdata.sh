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
  rm -rf go
  rm -rf graphql
}

init() {
  mkdir go
  mkdir graphql

  go generate ./tools.go
}

execProtoc() {
  for protoDir in ./proto/*; do
    protoc -I proto -I ../.. --include_source_info --include_imports --descriptor_set_out=${protoDir}/descriptor_set.pb ${protoDir}/*.proto
    protoc -I proto -I ../.. --plugin=protoc-gen-graphql='./bin/protoc-gen-graphql' --graphql_out=graphql  ${protoDir}/*.proto
    protoc -I proto -I ../.. --go_out=.. ${protoDir}/*.proto
  done
}

initGoMod() {
  for pbgoDir in ./go/*; do
    pushd $pbgoDir
    go mod init apis/go/$(basename $pbgoDir)
    echo 'replace github.com/izumin5210/remixer => ../../../..' >> go.mod
    go mod tidy
    popd
  done
}

_main() {
  pushd testdata/apis

  clean
  init
  execProtoc
  initGoMod

  popd
}

_main
