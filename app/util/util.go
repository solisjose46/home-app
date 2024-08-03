package util

import (
    "fmt"
)

const (
    TmplDir = "web/templates/"
    WebStaticDir = "web/static/"
    StaticDir = "/static/"
    HtmlExtension = ".html"
)

func GetTmplPath(tmpl string) string {
    return fmt.Sprintf("%s%s%s", TmplDir, tmpl, HtmlExtension)
}