package where

import (
	"os"
	"slices"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var Executables = func() []string {
	exes := slices.DeleteFunc(strings.Split(strings.ToLower(os.Getenv("PATHEXT")), ";"), func(e string) bool {
		return strings.TrimSpace(e) == ""
	})

	if len(exes) == 0 {
		return []string{".exe", ".cmd", ".com", ".bat"}
	}

	return exes
}()

func Expand(txt string) string {
	value, err := registry.ExpandString(txt)
	if err != nil {
		return txt
	}

	return value
}
