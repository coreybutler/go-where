package where

import (
	"os"
)

var Executables = []string{".bin", ".pkg", ".sh", "", ".bash", ".zsh", ".command", ".run"}

func Expand(txt string) string {
	return os.ExpandEnv(txt)
}
