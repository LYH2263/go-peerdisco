package validate

import "strings"

func ID(s string) bool {
	if s == "" || len(s) > 128 {
		return false
	}
	return !strings.ContainsAny(s, " \t\r\n")
}

func Addr(s string) bool {
	return s != "" && strings.Contains(s, ":")
}
