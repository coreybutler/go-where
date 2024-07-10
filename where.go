package where

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	fs "github.com/coreybutler/go-fsutil"
)

var str string

// Find the first path containing the executable
func Find(executable string, recursive ...bool) (string, error) {
	executable = filepath.Base(executable)
	result, _ := FindAll(executable, recursive...)

	if len(result) == 0 {
		if runtime.GOOS == "windows" {
			// Terminal
			path, err := exec.Command("where", executable).Output()
			if err == nil {
				return strings.TrimSpace(strings.Split(string(path), "\n")[0]), nil
			}

			// Powershell
			path, err = exec.Command("powershell", "-Command", "Get-Command", executable, "| Select-Object -ExpandProperty Source").Output()
			if err == nil {
				return strings.TrimSpace(strings.Split(string(path), "\n")[0]), nil
			}
		}

		return str, errors.New("not found")
	}

	return result[0], nil
}

func FindExcept(executable string, recursive bool, except ...string) (string, error) {
	executable = filepath.Base(executable)
	result, _ := FindAll(executable, recursive)

	if len(result) == 0 {
		return str, errors.New("not found")
	}

	for _, path := range result {
		if !contains(except, path) {
			return path, nil
		}
	}

	return str, errors.New("not found")
}

func FindAllExcept(executable string, recursive bool, except ...string) ([]string, error) {
	executable = filepath.Base(executable)
	result, _ := FindAll(executable, recursive)

	if len(result) == 0 {
		return []string{}, errors.New("not found")
	}

	results := make([]string, 0)
	for _, path := range result {
		if !contains(except, path) {
			results = append(results, path)
		}
	}

	return results, nil
}

// Find all known locations of an executable
func FindAll(executable string, recursive ...bool) ([]string, error) {
	executable = filepath.Base(executable)
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
				newPath := filepath.Join(dir, executable)
				if !contains(results, newPath) {
					results = append(results, newPath)
				}
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
					if !contains(results, file) {
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
	}

	return results, nil
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
