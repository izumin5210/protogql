package gqlutil

import (
	"errors"
	"fmt"
)

var ErrInvalidDirective = errors.New("invalid argument")

type ProtoDirective struct {
	FullName  string
	Package   string
	Name      string
	GoPackage string
	GoName    string
}

func (d *ProtoDirective) IsValid() bool {
	return d.FullName != "" && d.Package != "" && d.Name != "" && d.GoPackage != "" && d.GoName != ""
}

type ProtoFieldDirective struct {
	Name   string
	Type   string
	GoName string
	GoType string
}

func (d *ProtoFieldDirective) IsValid() bool {
	return d.Name != "" && d.Type != "" && d.GoName != "" && d.GoType != ""
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
				}
			}
			if !out.IsValid() {
				return nil, fmt.Errorf("invalid proto directive: %w", ErrInvalidDirective)
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
				case arg.Name == "goType" && arg.Value.Kind == ast.StringValue:
					out.GoType = arg.Value.Raw
				}
			}
			if !out.IsValid() {
				return nil, fmt.Errorf("invalid protoField directive: %w", ErrInvalidDirective)
			}
			return out, nil
		}
	}
	return nil, nil
}
