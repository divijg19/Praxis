package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	binPath = "/tmp/praxis_test"
	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func runPraxis(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr), exitErr.ExitCode()
		}
		t.Fatal(err)
	}
	return string(out), 0
}

func TestListCount(t *testing.T) {
	out, code := runPraxis(t, "list")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 41 {
		t.Errorf("got %d lines, want 41", len(lines))
	}
}

func TestChallengeLookup(t *testing.T) {
	out, code := runPraxis(t, "challenge", "yank_line_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("challenge yank_line_hunter returned %d lines", len(lines))
	}
	if lines[0] != "Every yank enters a register. Yank and paste to see." {
		t.Errorf("unexpected instruction: %q", lines[0])
	}
}

func TestTargetLookup(t *testing.T) {
	out, code := runPraxis(t, "target", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if strings.TrimSpace(out) != "★" {
		t.Errorf("unexpected target: %q", out)
	}
}

func TestVerifyLookup(t *testing.T) {
	out, code := runPraxis(t, "verify", "delete_word_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if strings.TrimSpace(out) != "buffer" {
		t.Errorf("unexpected verify: %q", out)
	}
}

func TestResultLookup(t *testing.T) {
	out, code := runPraxis(t, "result", "change_word_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("result change_word_hunter returned %d lines", len(lines))
	}
	if lines[2] != "bar" {
		t.Errorf("unexpected result line: %q", lines[2])
	}
}

func TestUnknownChallengeFails(t *testing.T) {
	_, code := runPraxis(t, "challenge", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestUnknownTargetFails(t *testing.T) {
	_, code := runPraxis(t, "target", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestUnknownVerifyFails(t *testing.T) {
	_, code := runPraxis(t, "verify", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
