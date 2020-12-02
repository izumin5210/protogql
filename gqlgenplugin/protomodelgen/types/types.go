package types

import "fmt"

type ProtoType interface {
	ProtoLikeType
	PbGoTypeName() string
}

type ProtoWrapperType interface {
	ProtoLikeType
	GoWrapperTypeName() string
}

type ProtoLikeType interface {
	Type
	FuncNameFromProto() string
	FuncNameFromRepeatedProto() string
	FuncNameToProto() string
	FuncNameToRepeatedProto() string
}

type Type interface {
	GoTypeName() string
}

type UnionMemberType interface {
	Type() Type
}

func GoWrapperTypeName(t Type) string {
	switch t := t.(type) {
	case ProtoType:
		return t.PbGoTypeName()
	case ProtoWrapperType:
		return t.GoWrapperTypeName()
	default:
		return t.GoTypeName()
	}
}

func UnwrapStatement(t Type, varName string) string {
	switch t := t.(type) {
	case ProtoLikeType:
		return fmt.Sprintf("%s(%s)", t.FuncNameToProto(), varName)
	default:
		return varName
	}
}

func IsProtoType(t Type) bool {
	_, ok := t.(ProtoType)
	return ok
}
