package main

import (
	"os"
	"testing"
)

func TestMonitorScriptExistsAndExecutable(t *testing.T) {
	path := "scripts/monitor.sh"
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("monitor script missing at %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file, found directory at %s", path)
	}
	mode := info.Mode()
	if mode&0o111 == 0 {
		t.Fatalf("monitor script is not executable: mode=%v", mode)
	}
}

func TestBackupScriptExistsAndExecutable(t *testing.T) {
	path := "scripts/backup.sh"
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("backup script missing at %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file, found directory at %s", path)
	}
	mode := info.Mode()
	if mode&0o111 == 0 {
		t.Fatalf("backup script is not executable: mode=%v", mode)
	}
}
