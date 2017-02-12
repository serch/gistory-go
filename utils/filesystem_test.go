package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOrDirExists(t *testing.T) {
	t.Parallel()
	tmpDir, _ := ioutil.TempDir("", "")

	if exists, _ := FileOrDirExists(tmpDir); !exists {
		t.Errorf("Expected directory %s to exist.", tmpDir)
	}

	tmpFile := createTemporaryFile(tmpDir)

	if exists, _ := FileOrDirExists(tmpFile); !exists {
		t.Errorf("Expected file %s to exist.", tmpFile)
	}

	os.RemoveAll(tmpDir)

	if exists, _ := FileOrDirExists(tmpFile); exists {
		t.Errorf("Expected file %s to not exist.", tmpFile)
	}

	if exists, _ := FileOrDirExists(tmpDir); exists {
		t.Errorf("Expected directory %s to not exist.", tmpDir)
	}
}

func createTemporaryFile(dir string) string {
	content := []byte("temporary file's content")
	tmpFile := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpFile, content, 0666); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}
