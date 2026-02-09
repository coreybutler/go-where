package where

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	fs "github.com/coreybutler/go-fsutil"
)

type Options struct {
	All            bool          `json:"all"`
	Recursive      bool          `json:"recursive"`
	FollowSymlinks bool          `json:"follow_symlinks"`
	Except         []string      `json:"except"`
	Timeout        time.Duration `json:"timeout"` // Default: 5 seconds
}

var emptybool bool
var Extensions []string
var AltPaths []string
var DisableExtensionChecking bool = false

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
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	executable = filepath.Base(executable)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	result, _ := seek(ctx, executable, cfg)

	if len(result) == 0 {
		return []string{}, errors.New("not found")
	}

	return result, nil
}

func isExecutableExtension(exe string) bool {
	if DisableExtensionChecking {
		return false
	}

	return (contains(Executables, filepath.Ext(exe)) || contains(Extensions, filepath.Ext(exe)))
}

// fileExists checks if a file exists using Lstat (doesn't follow symlinks in the path)
func fileExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// isFileExecutable checks if a file is executable, following symlinks in the file itself
func isFileExecutable(path string) bool {
	// Use fs.IsExecutable which handles the os.Stat call for us
	return fs.IsExecutable(path)
}

func seek(ctx context.Context, exe string, opts Options) ([]string, error) {
	paths := strings.Split(os.Getenv("PATH"), ";")
	results := []string{}
	resultMap := make(map[string]bool) // Track unique results

	for _, path := range AltPaths {
		paths = append(paths, path)
	}

	for _, dir := range paths {
		// Check for timeout
		select {
		case <-ctx.Done():
			return []string{}, errors.New("operation timed out")
		default:
		}

		// Expand any environment variables
		expandedDir := Expand(dir)

		// First, try direct search in the original expanded directory
		fullPath := filepath.Join(expandedDir, exe)
		if fileExists(fullPath) {
			if isFileExecutable(fullPath) {
				if isExecutableExtension(exe) {
					if !resultMap[fullPath] {
						results = append(results, fullPath)
						resultMap[fullPath] = true
					}
				}
			}
		}

		// Always do glob search when the executable name has no extension
		// This ensures we find all variations including those in symlinked directories
		if len(strings.Split(strings.TrimSpace(exe), ".")) <= 1 {
			file := strings.TrimSpace(exe)

			// Support all file extensions if none is specified
			file = strings.TrimSpace(file) + ".*"

			// Support recursive lookups
			globSearchDir := expandedDir
			if opts.Recursive {
				globSearchDir = filepath.Join(expandedDir, "**")
			}

			// Use glob matching to find the executable in original directory
			matches, err := filepath.Glob(filepath.Join(globSearchDir, file))
			if err == nil {
				for _, matchedFile := range matches {
					// Skip if already found
					if resultMap[matchedFile] {
						continue
					}

					// Determine whether the file is executable or not
					if isFileExecutable(matchedFile) {
						results = append(results, matchedFile)
						resultMap[matchedFile] = true
					} else {
						// Check if the matched file has an executable extension
						if isExecutableExtension(matchedFile) || isExecutableExtension(exe) {
							results = append(results, matchedFile)
							resultMap[matchedFile] = true
						}
					}
				}
			}

			// Also manually search in resolved symlink path if different and if FollowSymlinks is enabled
			if opts.FollowSymlinks {
				absPath, _ := filepath.Abs(expandedDir)
				evalPath, _ := filepath.EvalSymlinks(absPath)
				if evalPath != "" && evalPath != absPath {
					// Read directory contents from the resolved path
					entries, err := os.ReadDir(evalPath)
					if err == nil {
						for _, entry := range entries {
							// Skip directories unless recursive
							if entry.IsDir() && !opts.Recursive {
								continue
							}

							// Check if entry matches the pattern
							matched, _ := filepath.Match(file, entry.Name())
							if matched {
								resolvedFile := filepath.Join(evalPath, entry.Name())
								// Convert back to original symlinked path
								relPath, _ := filepath.Rel(evalPath, resolvedFile)
								resultPath := filepath.Join(expandedDir, relPath)

								// Skip if already found
								if resultMap[resultPath] {
									continue
								}

								// Determine whether the file is executable or not
								if isFileExecutable(resolvedFile) {
									results = append(results, resultPath)
									resultMap[resultPath] = true
								} else {
									// Check if the matched file has an executable extension
									if isExecutableExtension(resultPath) || isExecutableExtension(exe) {
										results = append(results, resultPath)
										resultMap[resultPath] = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if len(opts.Except) > 0 {
		// Remove exceptions while preserving order
		filtered := []string{}
		for _, path := range results {
			shouldExclude := false
			for _, pattern := range opts.Except {
				matched, _ := filepath.Match(pattern, path)
				if matched {
					shouldExclude = true
					break
				}
			}
			if !shouldExclude {
				filtered = append(filtered, path)
			}
		}
		results = filtered
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
