package protoutil

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func FormatComments(set protogen.CommentSet) string {
	var lines []string
	for _, l := range strings.Split(strings.TrimSpace(string(set.Leading)), "\n") {
		lines = append(lines, strings.TrimSpace(l))
	}
	return strings.Join(lines, "\n")
}
