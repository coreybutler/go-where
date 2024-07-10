package where

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	fs "github.com/coreybutler/go-fsutil"
)

type Options struct {
	All       bool     `json:"all"`
	Recursive bool     `json:"recursive"`
	Except    []string `json:"except"`
}

var emptybool bool

// Find the first path containing the executable
func Find(executable string, opts ...Options) ([]string, error) {
	var cfg Options
	if len(opts) > 0 {
		cfg = opts[0]
	} else {
		cfg = Options{}
	}

	if cfg.All == emptybool {
		cfg.All = false
	}
	if cfg.Recursive == emptybool {
		cfg.Recursive = true
	}
	if cfg.Except == nil {
		cfg.Except = []string{}
	}

	executable = filepath.Base(executable)
	result, _ := seek(executable, cfg)

	if len(result) == 0 {
		return []string{}, errors.New("not found")
	}

	return result, nil
}

func seek(exe string, opts Options) ([]string, error) {
	paths := strings.Split(os.Getenv("PATH"), ";")
	results := make(map[string]bool)

	for _, dir := range paths {
		// If file exists, add the path
		if fs.Exists(filepath.Join(dir, exe)) {
			if fs.IsExecutable(filepath.Join(dir, exe)) || contains(Executables, filepath.Ext(exe)) {
				results[filepath.Join(dir, exe)] = true
			}
		} else {
			// Expand any environment variables
			dir = Expand(dir)

			// If the file does not exist
			file := exe

			// Support all file extensions if none is specified
			if len(strings.Split(file, ".")) == 0 {
				file = file + ".*"
			}

			// Support recursive lookups
			if opts.Recursive {
				dir = filepath.Join(dir, "**")
			}

			// Use glob matching to find the executable
			matches, err := filepath.Glob(filepath.Join(dir, file))
			if err == nil {
				for _, file := range matches {
					// Determine whether the file is executable or not
					if fs.IsExecutable(file) {
						results[file] = true
					} else {
						if contains(Executables, filepath.Ext(file)) || file == filepath.Join(dir, exe) {
							results[file] = true
						}
					}
				}
			}
		}
	}

	if len(opts.Except) > 0 {
		// Remove exceptions
		for path, _ := range results {
			for _, pattern := range opts.Except {
				matched, _ := filepath.Match(pattern, path)
				if matched {
					delete(results, path)
				}
			}
		}
	}

	final := []string{}
	for path := range results {
		final = append(final, path)
	}

	return final, nil
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if value == item {
			return true
		}
		matched, err := filepath.Match(item, value)
		if err == nil && matched {
			return true
		}
	}

	return false
}
