package metrics

import "strings"

// LabelKey 规范化标签键。
func LabelKey(parts ...string) string {
	return strings.Join(parts, ".")
}

// Sanitize 去掉空白。
func Sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '	' {
			return '_'
		}
		return r
	}, s)
}
