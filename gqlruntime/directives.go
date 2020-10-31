package gqlruntime

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func GrpcDirective(ctx context.Context, obj interface{}, next graphql.Resolver, service string, rpc string) (res interface{}, err error) {
	return next(ctx)
}

func ProtoDirective(ctx context.Context, obj interface{}, next graphql.Resolver, fullName string, packageArg string, name string, goPackage string, goName string) (res interface{}, err error) {
	return next(ctx)
}

func ProtoFieldDirective(ctx context.Context, obj interface{}, next graphql.Resolver, name string, typeArg string, goName string, goTypeName string, goTypePackage *string) (res interface{}, err error) {
	return next(ctx)
}
