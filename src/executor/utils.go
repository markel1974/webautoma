package executor

import (
	"strconv"
	"strings"
)

func parseInt(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

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
