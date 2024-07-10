package where

import (
	"os"
)

var Executables = []string{".bin"}

func Expand(txt string) string {
	return os.ExpandEnv(txt)
}
