package types

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
)

var dateTimeFmts = []string{time.RFC3339Nano, time.RFC3339}

func MarshalTimestamp(ts *timestamp.Timestamp) graphql.Marshaler {
	if ts == nil {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(ptypes.TimestampString(ts)))
	})
}

func UnmarshalTimestamp(v interface{}) (*timestamp.Timestamp, error) {
	switch v := v.(type) {
	case string:
		var err error
		for _, f := range dateTimeFmts {
			var t time.Time
			t, err = time.Parse(f, v)
			if err == nil {
				return ptypes.TimestampProto(t)
			}
		}
		return nil, err
	default:
		return nil, fmt.Errorf("%T is not an DateTime", v)
	}
}

func MarshalInt32Value(in *wrappers.Int32Value) graphql.Marshaler {
	if in == nil {
		return graphql.Null
	}
	return graphql.MarshalInt32(in.GetValue())
}

func UnmarshalInt32Value(in interface{}) (*wrappers.Int32Value, error) {
	v, err := graphql.UnmarshalInt32(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.Int32Value{Value: v}, nil
}

func MarshalInt64Value(v *wrappers.Int64Value) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphql.MarshalInt64(v.GetValue())
}

func UnmarshalInt64Value(in interface{}) (*wrappers.Int64Value, error) {
	v, err := graphql.UnmarshalInt64(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.Int64Value{Value: v}, nil
}

func MarshalUInt32Value(in *wrappers.UInt32Value) graphql.Marshaler {
	if in == nil {
		return graphql.Null
	}
	return MarshalUint32(in.GetValue())
}

func UnmarshalUInt32Value(in interface{}) (*wrappers.UInt32Value, error) {
	v, err := UnmarshalUint32(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.UInt32Value{Value: v}, nil
}

func MarshalUInt64Value(v *wrappers.UInt64Value) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return MarshalUint64(v.GetValue())
}

func UnmarshalUInt64Value(in interface{}) (*wrappers.UInt64Value, error) {
	v, err := UnmarshalUint64(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.UInt64Value{Value: v}, nil
}

func MarshalFloatValue(v *wrappers.FloatValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%g", v.GetValue()))
	})
}

func UnmarshalFloatValue(in interface{}) (*wrappers.FloatValue, error) {
	v, err := graphql.UnmarshalFloat(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.FloatValue{Value: float32(v)}, nil
}

func MarshalDoubleValue(v *wrappers.DoubleValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphql.MarshalFloat(v.GetValue())
}

func UnmarshalDoubleValue(in interface{}) (*wrappers.DoubleValue, error) {
	v, err := graphql.UnmarshalFloat(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.DoubleValue{Value: v}, nil
}

func MarshalBoolValue(v *wrappers.BoolValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphql.MarshalBoolean(v.GetValue())
}

func UnmarshalBoolValue(in interface{}) (*wrappers.BoolValue, error) {
	v, err := graphql.UnmarshalBoolean(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.BoolValue{Value: v}, nil
}

func MarshalStringValue(v *wrappers.StringValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return graphql.MarshalString(v.GetValue())
}

func UnmarshalStringValue(in interface{}) (*wrappers.StringValue, error) {
	v, err := graphql.UnmarshalString(in)
	if err != nil {
		return nil, err
	}
	return &wrappers.StringValue{Value: v}, nil
}
