package gqls

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
)

type SchemaBuilder struct {
	fd protoreflect.FileDescriptor
	*Schema
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{}
}

func (b *SchemaBuilder) Build(fd protoreflect.FileDescriptor) (*Schema, error) {
	b.fd = fd
	b.Schema = newSchema()

	typeDescriptors, err := protoutil.TypeDFS(fd)
	if err != nil {
		return nil, err
	}
	for _, td := range typeDescriptors {
		t, err := TypeFromProto(td)
		if err != nil {
			return nil, err
		}
		err = b.AddType(t)
		if err != nil {
			return nil, err
		}
	}

	err = protoutil.RangeServices(fd, func(sd protoreflect.ServiceDescriptor) error {
		err := protoutil.RangeMethods(sd, func(md protoreflect.MethodDescriptor) error {
			if q, ok := NewQuery(md); ok {
				b.AddQuery(q)
			}
			if m, ok := NewMutation(md); ok {
				b.AddMutation(m)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return b.Schema, nil
}

func (b *SchemaBuilder) AddType(t Type) error {
	dt, ok := t.(interface {
		Type
		Definable
	})
	if !ok {
		return nil
	}
	if dt.ProtoDescriptor().ParentFile() != b.fd {
		return nil
	}
	// TODO: should handle collisions
	b.Types[dt.Name()] = dt

	if ot, ok := dt.(*ObjectType); ok {
		it := NewInputObjectType(ot)
		b.Types[it.Name()] = it
	}

	return nil
}

func (b *SchemaBuilder) AddQuery(q *Query) {
	b.Queries = append(b.Queries, q)
}

func (b *SchemaBuilder) AddMutation(m *Mutation) {
	b.Mutations = append(b.Mutations, m)
}
