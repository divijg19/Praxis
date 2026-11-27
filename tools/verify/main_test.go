package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// A broken describe path must fail certification rather than reporting 0/0 PASS.
func TestReplayFailsWithBrokenBinary(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the failing-binary stub is a POSIX shell script")
	}
	if _, err := exec.LookPath("nvim"); err != nil {
		t.Skip("nvim is required for this integration test")
	}
	stub := filepath.Join(t.TempDir(), "praxis-failing")
	if err := os.WriteFile(stub, []byte("#!/bin/sh\nexit 1\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	root, err := findRoot()
	if err != nil {
		t.Fatal(err)
	}

	env := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "PRAXIS_BIN=") {
			env = append(env, entry)
		}
	}
	env = append(env, "PRAXIS_BIN="+stub)

	ctx, cancel := context.WithTimeout(context.Background(), harnessTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nvim", "--headless", "-u", "NONE", "-l", filepath.Join(root, "tools", "replay", "replay.lua"))
	cmd.Dir = root
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("replay passed with a broken binary:\n%s", out)
	}
	text := string(out)
	if !strings.Contains(text, "Describe failures:") && !strings.Contains(text, "catalog --ids returned no challenges") {
		t.Fatalf("expected a catalog or describe failure in output, got:\n%s", out)
	}
}
