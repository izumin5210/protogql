package gqlutil

import (
	stderrors "errors"
	"strings"

	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
)

var ErrInvalidDirective = stderrors.New("invalid argument")

type ProtoDirective struct {
	FullName  string
	Package   string
	Name      string
	GoPackage string
	GoName    string
	Oneof     *ProtoDirectiveOneof
}

type ProtoDirectiveOneof struct {
	Fields []*ProtoDirectiveOneofField
}

type ProtoDirectiveOneofField struct {
	Name   string
	GoName string
}

func (d *ProtoDirective) IsValid() bool {
	return d.FullName != "" && d.Package != "" && d.Name != "" && d.GoPackage != "" && d.GoName != ""
}

type ProtoFieldDirective struct {
	Name          string
	Type          string
	GoName        string
	GoTypeName    string
	GoTypePackage string
	OneofName     string
	OneofGoName   string
}

func (d *ProtoFieldDirective) IsValid() bool {
	return d.Name != "" && d.Type != "" && d.GoName != "" && d.GoTypeName != ""
}

func (f *ProtoDirectiveOneofField) IsValid() bool {
	return f.Name != "" && f.GoName != ""
}

func (d *ProtoFieldDirective) IsWellKnownType() bool {
	switch d.Type {
	case "google.protobuf.Int32Value", "google.protobuf.Int64Value",
		"google.protobuf.UInt32Value", "google.protobuf.UInt64Value",
		"google.protobuf.FloatValue", "google.protobuf.DoubleValue",
		"google.protobuf.BoolValue",
		"google.protobuf.StringValue",
		"google.protobuf.Timestamp":
		return true
	}
	return false
}

func (d *ProtoFieldDirective) IsGoBuiltinType() bool {
	return strings.ToLower(d.Type) == d.Type
}

func ExtractProtoDirective(directives ast.DirectiveList) (*ProtoDirective, error) {
	for _, d := range directives {
		if d.Name == "proto" {
			out := new(ProtoDirective)
			for _, arg := range d.Arguments {
				switch {
				case arg.Name == "fullName" && arg.Value.Kind == ast.StringValue:
					out.FullName = arg.Value.Raw
				case arg.Name == "package" && arg.Value.Kind == ast.StringValue:
					out.Package = arg.Value.Raw
				case arg.Name == "name" && arg.Value.Kind == ast.StringValue:
					out.Name = arg.Value.Raw
				case arg.Name == "goPackage" && arg.Value.Kind == ast.StringValue:
					out.GoPackage = arg.Value.Raw
				case arg.Name == "goName" && arg.Value.Kind == ast.StringValue:
					out.GoName = arg.Value.Raw
				case arg.Name == "oneof" && arg.Value.Kind == ast.ObjectValue:
					oneof, err := extractProtoDirectiveOneof(arg.Value.Children)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					out.Oneof = oneof
				}
			}
			if !out.IsValid() {
				return nil, ErrInvalidDirective
			}
			return out, nil
		}
	}
	return nil, nil
}

func ExtractProtoFieldDirective(directives ast.DirectiveList) (*ProtoFieldDirective, error) {
	for _, d := range directives {
		if d.Name == "protoField" {
			out := new(ProtoFieldDirective)
			for _, arg := range d.Arguments {
				switch {
				case arg.Name == "name" && arg.Value.Kind == ast.StringValue:
					out.Name = arg.Value.Raw
				case arg.Name == "type" && arg.Value.Kind == ast.StringValue:
					out.Type = arg.Value.Raw
				case arg.Name == "goName" && arg.Value.Kind == ast.StringValue:
					out.GoName = arg.Value.Raw
				case arg.Name == "goTypeName" && arg.Value.Kind == ast.StringValue:
					out.GoTypeName = arg.Value.Raw
				case arg.Name == "goTypePackage" && arg.Value.Kind == ast.StringValue:
					out.GoTypePackage = arg.Value.Raw
				case arg.Name == "oneofName" && arg.Value.Kind == ast.StringValue:
					out.OneofName = arg.Value.Raw
				case arg.Name == "oneofGoName" && arg.Value.Kind == ast.StringValue:
					out.OneofGoName = arg.Value.Raw
				}
			}
			if !out.IsValid() {
				return nil, ErrInvalidDirective
			}
			return out, nil
		}
	}
	return nil, nil
}

func extractProtoDirectiveOneof(children ast.ChildValueList) (*ProtoDirectiveOneof, error) {
	oneof := &ProtoDirectiveOneof{}
	for _, child := range children {
		switch {
		case child.Name == "fields" && child.Value.Kind == ast.ListValue:
			for _, child := range child.Value.Children {
				f := &ProtoDirectiveOneofField{}
				for _, child := range child.Value.Children {
					switch {
					case child.Name == "name" && child.Value.Kind == ast.StringValue:
						f.Name = child.Value.Raw
					case child.Name == "goName" && child.Value.Kind == ast.StringValue:
						f.GoName = child.Value.Raw
					}
				}
				if !f.IsValid() {
					return nil, ErrInvalidDirective
				}
				oneof.Fields = append(oneof.Fields, f)
			}
		}
	}
	return oneof, nil
}
