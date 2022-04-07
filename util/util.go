package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

func FilePathWalkFollowLink(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(GetActuallyDir(root), walkFn)
}

func GetActuallyDir(root string) string {
	dirInfo, err := os.Lstat(root)
	if err != nil {
		return root
	}
	if dirInfo.Mode()&os.ModeSymlink != 0 {
		dirLinkTo, err := os.Readlink(root)
		if err != nil {
			return root
		}
		return dirLinkTo
	}
	return root
}

func StringSetAdd(s []string, v ...string) []string {
	for _, vv := range v {
		if !ContainsString(s, vv) {
			s = append(s, vv)
		}
	}
	return s
}

// ContainsString returns true if a string is present in a iteratee.
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
func PanicIfTrue(condation bool, fmtStr string, args ...interface{}) {
	if !condation {
		return
	}
	panic(fmt.Errorf(fmtStr, args...))
}

// And now lots of helper functions.

// Is c an ASCII lower-case letter?
func IsASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func IsASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && IsASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if IsASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if IsASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && IsASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// CamelCaseSlice is like CamelCase, but the argument is a slice of strings to
// be joined with "_".
func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

// PanicIfErrorAsFisrt err不为nil则wrap并panic，将err作为第一个fmt的参数
func PanicIfErrorAsFisrt(err error, fmtStr string, args ...interface{}) {
	if err == nil {
		return
	}
	var argList []interface{}
	argList = append(argList, err)
	argList = append(argList, args...)
	panic(fmt.Errorf(fmtStr, argList...))
}
func FileGetContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
func FileWalkFuncWithExcludeFilter(files *[]string, excluded func(f string) bool, ext ...string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if excluded != nil && excluded(path) {
			return err
		}
		if len(ext) > 0 {
			if ContainsString(ext, filepath.Ext(path)) {
				*files = append(*files, path)
			}
		} else {
			*files = append(*files, path)
		}
		return err
	}
}
