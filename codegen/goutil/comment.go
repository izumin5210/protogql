package goutil

import "strings"

func ToComment(desc string) string {
	if desc == "" {
		return ""
	}
	var lines []string
	for _, l := range strings.Split(desc, "\n") {
		lines = append(lines, "// "+l)
	}
	return strings.Join(lines, "\n")
}
