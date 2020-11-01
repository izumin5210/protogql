package types_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/izumin5210/remixer/gqlruntime/types"
)

func TestTimestampValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *timestamp.Timestamp, want string) {
			var buf bytes.Buffer
			types.MarshalTimestamp(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalTimestamp(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &timestamp.Timestamp{Seconds: 1604238614, Nanos: 877253000}, `"2020-11-01T13:50:14.877253Z"`)
		testMarshal(t, &timestamp.Timestamp{Seconds: 1604238614, Nanos: 0}, `"2020-11-01T13:50:14Z"`)
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, wantSecs int64, wantNanos int32) {
			got, err := types.UnmarshalTimestamp(in)
			if err != nil {
				t.Errorf("UnmarshalTimestamp(%v) returns an error: %v", in, err)
			}
			if got, want := got.GetSeconds(), wantSecs; got != want {
				t.Errorf("UnmarshalTimestamp(%v) returns %v seconds, want %v", in, got, want)
			}
			if got, want := got.GetNanos(), wantNanos; got != want {
				t.Errorf("UnmarshalTimestamp(%v) returns %v nanos, want %v", in, got, want)
			}
		}

		testUnmarshal(t, "2020-11-01T13:50:14.877253Z", 1604238614, 877253000)
		testUnmarshal(t, "2020-11-01T22:50:14.877253+09:00", 1604238614, 877253000)
		testUnmarshal(t, "2020-11-01T13:50:14Z", 1604238614, 0)
	})
}

func TestInt32Value(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.Int32Value, want string) {
			var buf bytes.Buffer
			types.MarshalInt32Value(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalInt32Value(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.Int32Value{Value: 123}, "123")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalInt32Value(in)
			if err != nil {
				t.Errorf("UnmarshalInt32Value(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalInt32Value(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), int32(123))
		testUnmarshal(t, int64(123), int32(123))
		testUnmarshal(t, json.Number("123"), int32(123))
		testUnmarshal(t, "123", int32(123))
	})
}

func TestInt64Value(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.Int64Value, want string) {
			var buf bytes.Buffer
			types.MarshalInt64Value(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalInt64Value(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.Int64Value{Value: 123}, "123")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalInt64Value(in)
			if err != nil {
				t.Errorf("UnmarshalInt64Value(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalInt64Value(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), int64(123))
		testUnmarshal(t, int64(123), int64(123))
		testUnmarshal(t, json.Number("123"), int64(123))
		testUnmarshal(t, "123", int64(123))
	})
}

func TestUInt32Value(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.UInt32Value, want string) {
			var buf bytes.Buffer
			types.MarshalUInt32Value(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalUInt32Value(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.UInt32Value{Value: 123}, "123")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalUInt32Value(in)
			if err != nil {
				t.Errorf("UnmarshalUInt32Value(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalUInt32Value(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), uint32(123))
		testUnmarshal(t, int64(123), uint32(123))
		testUnmarshal(t, json.Number("123"), uint32(123))
		testUnmarshal(t, "123", uint32(123))
	})
}

func TestUInt64Value(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.UInt64Value, want string) {
			var buf bytes.Buffer
			types.MarshalUInt64Value(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalUInt64Value(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.UInt64Value{Value: 123}, "123")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalUInt64Value(in)
			if err != nil {
				t.Errorf("UnmarshalUInt64Value(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalUInt64Value(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), uint64(123))
		testUnmarshal(t, int64(123), uint64(123))
		testUnmarshal(t, json.Number("123"), uint64(123))
		testUnmarshal(t, "123", uint64(123))
	})
}

func TestFloatValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.FloatValue, want string) {
			var buf bytes.Buffer
			types.MarshalFloatValue(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalFloatValue(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.FloatValue{Value: float32(1.23)}, "1.23")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalFloatValue(in)
			if err != nil {
				t.Errorf("UnmarshalFloatValue(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalFloatValue(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), float32(123))
		testUnmarshal(t, int64(123), float32(123))
		testUnmarshal(t, float64(12.3), float32(12.3))
		testUnmarshal(t, json.Number("12.3"), float32(12.3))
		testUnmarshal(t, "12.3", float32(12.3))
	})
}

func TestDoubleValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.DoubleValue, want string) {
			var buf bytes.Buffer
			types.MarshalDoubleValue(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalDoubleValue(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.DoubleValue{Value: float64(1.23)}, "1.23")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalDoubleValue(in)
			if err != nil {
				t.Errorf("UnmarshalDoubleValue(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalDoubleValue(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), float64(123))
		testUnmarshal(t, int64(123), float64(123))
		testUnmarshal(t, float64(12.3), float64(12.3))
		testUnmarshal(t, json.Number("12.3"), float64(12.3))
		testUnmarshal(t, "12.3", float64(12.3))
	})
}

func TestBoolValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.BoolValue, want string) {
			var buf bytes.Buffer
			types.MarshalBoolValue(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalBoolValue(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.BoolValue{Value: true}, "true")
		testMarshal(t, &wrappers.BoolValue{Value: false}, "false")
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalBoolValue(in)
			if err != nil {
				t.Errorf("UnmarshalBoolValue(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalBoolValue(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, true, true)
		testUnmarshal(t, "true", true)
	})
}

func TestStringValue(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testMarshal := func(t *testing.T, in *wrappers.StringValue, want string) {
			var buf bytes.Buffer
			types.MarshalStringValue(in).MarshalGQL(&buf)
			if got, want := buf.String(), want; got != want {
				t.Errorf("MarshalStringValue(%v) returns %q, want %q", in, got, want)
			}
		}

		testMarshal(t, &wrappers.StringValue{Value: "foo"}, `"foo"`)
		testMarshal(t, nil, "null")
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalStringValue(in)
			if err != nil {
				t.Errorf("UnmarshalStringValue(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got.GetValue(), want) {
				t.Errorf("UnmarshalStringValue(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, "foo", "foo")
	})
}
