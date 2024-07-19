package util

import (
	"strings"
)

func GetFilePath(paths ...string) string {
    var builder strings.Builder
    for _, path := range paths {
        builder.WriteString(path)
    }
    return builder.String()
}