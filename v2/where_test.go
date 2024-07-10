package where

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	file := "node.exe"
	// makeTestExecutable(file)

	out, err := Find(file)
	if err != nil {
		t.Log("File error: " + err.Error())
		t.Fail()
	}

	t.Log(out[0])
}

func TestFindException(t *testing.T) {
	file := "node.exe"
	// makeTestExecutable(file)

	out, err := Find(file)
	if err != nil {
		t.Log("File error: " + err.Error())
		t.Fail()
	}

	output, err := Find(file, Options{
		Except: out,
	})
	if err == nil {
		t.Logf("expected no results, got %d", len(output))
		t.Fail()
	}

	// if out != filepath.Join(root(), file) {
	// 	t.Logf("Expected %v, Received %v", filepath.Join(root(), file), out)
	// 	t.Fail()
	// }

	// clear(file)
}

func root() string {
	return strings.Split(os.Getenv("PATH"), ";")[0]
}

func makeTestExecutable(filename string) {
	file := filepath.Join(root(), filename)
	perm := os.ModePerm

	os.WriteFile(file, []byte("test"), perm)
	log.Print("Created temp file: " + file)
}

func clear(filename string) {
	file := filepath.Join(root(), filename)
	os.RemoveAll(file)
	log.Print("Removed temp file: " + file)
}
