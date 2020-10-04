package gqls

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/izumin5210/remixer/cmd/protoc-gen-graphql/protoutil"
	"github.com/izumin5210/remixer/options"
)

type Query struct {
	Proto   protoreflect.MethodDescriptor
	Options *options.GraphqlQueryOptions
}

func NewQuery(md protoreflect.MethodDescriptor) (*Query, bool) {
	if !proto.HasExtension(md.Options(), options.E_GraphqlQuery) {
		return nil, false
	}
	ext := proto.GetExtension(md.Options(), options.E_GraphqlQuery)
	opts, ok := ext.(*options.GraphqlQueryOptions)
	if !ok {
		return nil, false
	}
	return &Query{Proto: md, Options: opts}, true
}

func (q *Query) FieldDefinitionAST() (*ast.FieldDefinition, error) {
	def := &ast.FieldDefinition{
		Name:       q.Options.GetName(),
		Directives: rpcDirectivesAST(q.Proto),
	}

	var err error

	def.Arguments, err = q.argsAST()
	if err != nil {
		return nil, err
	}

	outputType, err := q.outputType()
	if err != nil {
		return nil, err
	}
	def.Type = outputType.TypeAST()

	return def, nil
}

func (q *Query) outputType() (Type, error) {
	if name := q.Options.GetOutput(); name != "" {
		var typ Type
		err := protoutil.RangeFields(q.Proto.Output(), func(fd protoreflect.FieldDescriptor) error {
			if string(fd.Name()) == name {
				var err error
				typ, err = TypeFromProtoField(fd)
				if err != nil {
					// TODO: handing
					return err
				}
				return protoutil.BreakRange
			}
			return fmt.Errorf("%s is not found in %s", name, fd.FullName())
		})
		if err != nil {
			return nil, err
		}
		return typ, nil
	} else {
		typ, err := TypeFromProto(q.Proto.Output())
		if err != nil {
			return nil, err
		}
		return typ, nil
	}
}

func (q *Query) argsAST() (ast.ArgumentDefinitionList, error) {
	var args ast.ArgumentDefinitionList
	err := protoutil.RangeFields(q.Proto.Input(), func(fd protoreflect.FieldDescriptor) error {
		typ, err := TypeFromProtoField(fd)
		if err != nil {
			// TODO: handing
			return err
		}
		args = append(args, &ast.ArgumentDefinition{
			Name: strcase.ToLowerCamel(string(fd.Name())),
			Type: typ.TypeAST(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return args, nil
}
