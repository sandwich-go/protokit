package protokit

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sandwich-go/boost/xmap"
	"github.com/sandwich-go/boost/xslice"
)

type Marker interface {
	AddVersion(kv ...string)
	AddSource(source ...string)
	Format()
}
type MarkerInfo struct {
	GeneratedBy                  string
	Versions                     map[string]string
	Sources                      []string
	MarkerLeadingWithDoubleSlash string
	MarkerLeadingWithDoubleDash  string
	MarkerLeadingWithHexKey      string
	MarkerForHTML                string
}

func (b *MarkerInfo) AddVersion(kv ...string) {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("AddVersion: Pairs got the odd number of input pairs for name-version: %d", len(kv)))
	}
	if b.Versions == nil {
		b.Versions = make(map[string]string)
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = strings.ToLower(s)
			continue
		}
		b.Versions[key] = s
	}
}
func (b *MarkerInfo) AddSource(source ...string) { b.Sources = append(b.Sources, source...) }
func (b *MarkerInfo) Format() {
	b.GeneratedBy = "Code generated by protokitgo. DO NOT EDIT."
	var formatList []string
	formatList = append(formatList, b.GeneratedBy)
	if len(b.Versions) != 0 {
		formatList = append(formatList, "versions:")
		xmap.WalkStringStringMapDeterministic(b.Versions, func(k string, v string) bool {
			formatList = append(formatList, fmt.Sprintf("    %s : %s", k, v))
			return true
		})
	}
	if len(b.Sources) != 0 {
		sort.Strings(b.Sources)
		formatList = append(formatList, fmt.Sprintf("source: %s", strings.Join(b.Sources, ", ")))
	}
	b.MarkerLeadingWithDoubleSlash = strings.Join(xslice.StringsAddPrefix(formatList, "// "), "\n")
	b.MarkerLeadingWithDoubleDash = strings.Join(xslice.StringsAddPrefix(formatList, "-- "), "\n")
	b.MarkerLeadingWithHexKey = strings.Join(xslice.StringsAddPrefix(formatList, "# "), "\n")
	b.MarkerForHTML = strings.Join(xslice.StringsAddSuffix(xslice.StringsAddPrefix(formatList, "<!-- "), " -->"), "\n")
}
