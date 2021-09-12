package where

import (
	"golang.org/x/sys/windows/registry"
)

var Executables = []string{".exe", ".cmd", ".com", ".bat"}

func Expand(txt string) string {
	value, err := registry.ExpandString(txt)
	if err != nil {
		return txt
	}

	return value
}
