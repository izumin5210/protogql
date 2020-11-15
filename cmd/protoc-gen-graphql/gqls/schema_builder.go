package gqls

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SchemaBuilder struct {
	fd protoreflect.FileDescriptor
	*Schema
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{}
}

func collectTypes(msgs []*protogen.Message, enums []*protogen.Enum) (allMsgs []*protogen.Message, allEnums []*protogen.Enum) {
	allMsgs = append(allMsgs, msgs...)
	allEnums = append(allEnums, enums...)

	for _, m := range msgs {
		submsgs, subenums := collectTypes(m.Messages, m.Enums)
		allMsgs = append(allMsgs, submsgs...)
		allEnums = append(allEnums, subenums...)
	}

	return
}

func (b *SchemaBuilder) Build(f *protogen.File) (*Schema, error) {
	fd := f.Desc
	b.fd = fd
	b.Schema = newSchema()

	var types []Type

	msgs, enums := collectTypes(f.Messages, f.Enums)
	for _, m := range msgs {
		t, err := TypeFromProtoMessage(m)
		if err != nil {
			return nil, err
		}
		types = append(types, t)

		// collect oneofs
		for _, o := range m.Oneofs {
			types = append(types, NewUnionType(o))
		}
	}
	for _, e := range enums {
		types = append(types, NewEnumType(e))
	}

	for _, t := range types {
		err := b.AddType(t)
		if err != nil {
			return nil, err
		}
	}

	for _, s := range f.Services {
		for _, m := range s.Methods {
			if m.Input.Desc.Name() == m.Desc.Name()+"Request" {
				delete(b.Types, string(m.Input.Desc.Name()))
			}
			if m.Output.Desc.Name() == m.Desc.Name()+"Response" {
				delete(b.Types, string(m.Output.Desc.Name()))
			}
		}
	}

	for _, t := range b.Types {
		if ot, ok := t.(*ObjectType); ok {
			it := NewInputObjectType(ot)
			b.Types[it.Name()] = it
		}
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

	return nil
}
