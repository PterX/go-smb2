package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPassthroughFSRenameKeepsOpenHandleUsable(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file.lock"), []byte("data"), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	fs := NewPassthroughFS(root)
	h, err := fs.Open("file.lock", os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer fs.Close(h)

	if err := fs.Rename(h, "file", 1); err != nil {
		t.Fatalf("rename through open handle: %v", err)
	}

	if _, err := fs.GetAttr(h); err != nil {
		t.Fatalf("get attr after rename through open handle: %v", err)
	}

	if err := fs.Unlink(h); err != nil {
		t.Fatalf("unlink renamed file through open handle: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, "file")); !os.IsNotExist(err) {
		t.Fatalf("renamed file still exists after unlink: %v", err)
	}
}
