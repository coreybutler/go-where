package where

import (
	"os"
)

var Executables = []string{".bin", ".pkg", ".sh"}

func Expand(txt string) string {
	return os.ExpandEnv(txt)
}
