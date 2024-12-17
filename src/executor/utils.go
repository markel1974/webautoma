package executor

import (
	"strconv"
	"strings"
)

// parseInt converts a string to an integer. Returns the integer value or zero if the conversion fails.
func parseInt(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

// computeFileId extracts the base name (without extension) from a given file path string.
func computeFileId(filename string) string {
	var pos = strings.LastIndex(filename, "\\")
	if pos < 0 {
		pos = strings.LastIndex(filename, "/")
	}
	if pos < 0 {
		pos = 0
	} else {
		pos++
	}
	name := filename[pos:]
	if pos = strings.LastIndex(name, "."); pos >= 0 {
		name = name[0:pos]
	}
	return name
}
