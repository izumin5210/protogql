package types_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/izumin5210/protogql/gqlruntime/types"
)

func TestUint32(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		var buf bytes.Buffer
		types.MarshalUint32(123).MarshalGQL(&buf)
		if got, want := buf.String(), "123"; got != want {
			t.Errorf("MarshalUint32(123) returns %q, want %q", got, want)
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalUint32(in)
			if err != nil {
				t.Errorf("UnmarshalUint32(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("UnmarshalUint32(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), uint32(123))
		testUnmarshal(t, int64(123), uint32(123))
		testUnmarshal(t, json.Number("123"), uint32(123))
		testUnmarshal(t, "123", uint32(123))
	})
}

func TestUint64(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		var buf bytes.Buffer
		types.MarshalUint64(123).MarshalGQL(&buf)
		if got, want := buf.String(), "123"; got != want {
			t.Errorf("MarshalUint64(123) returns %q, want %q", got, want)
		}
	})
	t.Run("unmarshal", func(t *testing.T) {
		testUnmarshal := func(t *testing.T, in interface{}, want interface{}) {
			got, err := types.UnmarshalUint64(in)
			if err != nil {
				t.Errorf("UnmarshalUint64(%v) returns an error: %v", in, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("UnmarshalUint64(%v) returns %v, want %v", in, got, want)
			}
		}

		testUnmarshal(t, int(123), uint64(123))
		testUnmarshal(t, int64(123), uint64(123))
		testUnmarshal(t, json.Number("123"), uint64(123))
		testUnmarshal(t, "123", uint64(123))
	})
}
