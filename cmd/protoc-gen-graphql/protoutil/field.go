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

var kindMap = map[protoreflect.Kind]JSONKind{
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
	return kindMap[k]
}
