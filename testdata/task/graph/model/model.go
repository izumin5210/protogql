package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(i, 10))
	})
}

func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case json.Number:
		return strconv.ParseUint(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}
