package util

import (
    "fmt"
	"strings"
)

const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    TmplDir = "web/templates/"
    HtmlExtension = ".html"
)

func PrintError(err error) {
    fmt.Println(fmt.Sprintf("%s%s%s", Red, err.Error(), Reset))
}

func PrintMessage(messages ...string) {
    appended := GetFilePath(messages...)
    fmt.Println(fmt.Sprintf("%s%s%s", Yellow, appended, Reset))
}

func PrintSuccess(messages ...string) {
    appended := GetFilePath(messages...)
    fmt.Println(fmt.Sprintf("%s%s%s", Green, appended, Reset))
}

func GetTmplPath(tmpl string) string {
    return GetFilePath(TmplDir, tmpl, HtmlExtension)
}

func GetFilePath(paths ...string) string {
    var builder strings.Builder
    for _, path := range paths {
        builder.WriteString(path)
    }
    return builder.String()
}