package protoutil

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type JSONKind int

const (
	_ JSONKind = iota
	JSONInt
	JSONFloat
	JSONString
	JSONBoolean
	JSONBase64String
	JSONEnumString
	JSONObject
)

var JSONKindMap = map[protoreflect.Kind]JSONKind{
	protoreflect.DoubleKind:  JSONFloat,
	protoreflect.FloatKind:   JSONFloat,
	protoreflect.Int64Kind:   JSONString,
	protoreflect.Uint64Kind:  JSONString,
	protoreflect.Int32Kind:   JSONInt,
	protoreflect.Fixed64Kind: JSONString,
	protoreflect.Fixed32Kind: JSONInt,
	protoreflect.BoolKind:    JSONBoolean,
	protoreflect.StringKind:  JSONString,
	// protoreflect.GroupKind:
	protoreflect.MessageKind:  JSONObject,
	protoreflect.BytesKind:    JSONBase64String,
	protoreflect.Uint32Kind:   JSONInt,
	protoreflect.EnumKind:     JSONEnumString,
	protoreflect.Sfixed32Kind: JSONInt,
	protoreflect.Sfixed64Kind: JSONString,
	protoreflect.Sint32Kind:   JSONInt,
	protoreflect.Sint64Kind:   JSONString,
}

func JSONKindFrom(k protoreflect.Kind) JSONKind {
	return JSONKindMap[k]
}

type GoKind int

const (
	_ GoKind = iota
	GoInt32
	GoInt64
	GoUint32
	GoUint64
	GoFloat32
	GoFloat64
	GoBool
	GoString
	GoBytes
)

var GoKindMap = map[protoreflect.Kind]GoKind{
	protoreflect.DoubleKind:  GoFloat64,
	protoreflect.FloatKind:   GoFloat32,
	protoreflect.Int64Kind:   GoInt64,
	protoreflect.Uint64Kind:  GoUint64,
	protoreflect.Int32Kind:   GoInt32,
	protoreflect.Fixed64Kind: GoUint64,
	protoreflect.Fixed32Kind: GoUint32,
	protoreflect.BoolKind:    GoBool,
	protoreflect.StringKind:  GoString,
	// protoreflect.GroupKind:
	// protoreflect.MessageKind:,
	protoreflect.BytesKind:  GoBytes,
	protoreflect.Uint32Kind: GoUint32,
	// protoreflect.EnumKind:,
	protoreflect.Sfixed32Kind: GoInt32,
	protoreflect.Sfixed64Kind: GoInt64,
	protoreflect.Sint32Kind:   GoInt32,
	protoreflect.Sint64Kind:   GoInt64,
}

func GoKindFrom(k protoreflect.Kind) GoKind {
	return GoKindMap[k]
}

func (k GoKind) String() string {
	switch k {
	case GoInt32:
		return "int32"
	case GoInt64:
		return "int64"
	case GoUint32:
		return "uint32"
	case GoUint64:
		return "uint64"
	case GoFloat32:
		return "float32"
	case GoFloat64:
		return "float64"
	case GoBool:
		return "bool"
	case GoString:
		return "string"
	case GoBytes:
		return "[]byte"
	}
	panic("unreachable")
}
