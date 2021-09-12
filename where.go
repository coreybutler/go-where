package where

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	fs "github.com/coreybutler/go-fsutil"
)

var str string

// Find the first path containing the executable
func Find(executable string, recursive ...bool) (string, error) {
	result := FindAll(executable, recursive...)

	if len(result) == 0 {
		return str, errors.New("not found")
	}

	return result[0], nil
}

// Find all known locations of an executable
func FindAll(executable string, recursive ...bool) []string {
	paths := strings.Split(os.Getenv("PATH"), ";")
	results := make([]string, 0)

	r := false
	if len(recursive) > 0 {
		r = recursive[0]
	}

	for _, dir := range paths {
		// If file exists, add the path
		if fs.Exists(filepath.Join(dir, executable)) {
			if fs.IsExecutable(filepath.Join(dir, executable)) || contains(Executables, filepath.Ext(executable)) {
				results = append(results, filepath.Join(dir, executable))
			}
		} else {
			// Expand any environment variables
			dir = Expand(dir)

			// If the file does not exist
			file := executable

			// Support all file extensions if none is specified
			if len(strings.Split(file, ".")) == 0 {
				file = file + ".*"
			}

			// Support recursive lookups
			if r {
				dir = filepath.Join(dir, "**")
			}

			// Use glob matching to find the executable
			matches, err := filepath.Glob(filepath.Join(dir, file))
			if err == nil {
				for _, file := range matches {
					// Determine whether the file is executable or not
					if fs.IsExecutable(file) {
						results = append(results, file)
					} else {
						if contains(Executables, filepath.Ext(file)) || file == filepath.Join(dir, executable) {
							results = append(results, file)
						}
					}
				}
			}
		}
	}

	return results
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if value == item {
			return true
		}
	}

	return false
}
