package protokit

import (
	"html/template"
	"sort"
	"strings"

	"github.com/sandwich-go/boost/xstrings"
)

func unescaped(str string) template.HTML       { return template.HTML(str) }
func Slash2Underline(s string) string          { return strings.ReplaceAll(s, "/", "_") }
func Slash2Dot(s string) string                { return strings.ReplaceAll(s, "/", ".") }
func Underline2Dot(s string) string            { return strings.ReplaceAll(s, "_", ".") }
func Dot2Underline(s string) string            { return strings.ReplaceAll(s, ".", "_") }
func Dot2Slash(s string) string                { return strings.ReplaceAll(s, ".", "/") }
func ReplaceEmpty(s string, old string) string { return strings.ReplaceAll(s, old, "") }
func SortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func SortedVals(m map[string]string) []string {
	var vals []string
	for _, v := range m {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	return vals
}

func SortStrings(m []string) []string {
	sort.Strings(m)
	return m
}

func ReplaceEmptyAndTile(s string, old string) string {
	return strings.Title(strings.ReplaceAll(s, old, ""))
}

func StringTitle(s string) string { return strings.Title(s) }

func ReplaceString(s string, old string, new string) string { return strings.ReplaceAll(s, old, new) }
func Join(elems []string, sep string) string                { return strings.Join(elems, sep) }

var funcMap = template.FuncMap{
	"unescaped":           unescaped,
	"slash2Underline":     Slash2Underline,
	"slash2Dot":           Slash2Dot,
	"underline2Dot":       Underline2Dot,
	"dot2Underline":       Dot2Underline,
	"dot2Slash":           Dot2Slash,
	"ReplaceEmpty":        ReplaceEmpty,
	"ReplaceEmptyAndTile": ReplaceEmptyAndTile,
	"ReplaceString":       ReplaceString,
	"CamelCase":           xstrings.CamelCase,
	"SnakeCase":           xstrings.SnakeCase,
	"FirstLower":          xstrings.FirstLower,
	"FirstUpper":          xstrings.FirstUpper,
	"SortedKeys":          SortedKeys,
	"SortedVals":          SortedVals,
	"Join":                Join,
	"SortStrings":         SortStrings,
}
