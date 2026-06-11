package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/divijg19/Praxis/internal/content"
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

func TestListOutputStable(t *testing.T) {
	out, code := runPraxis(t, "list")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	expected := make([]string, len(content.All()))
	for i, c := range content.All() {
		expected[i] = c.Name
	}
	if len(lines) != len(expected) {
		t.Fatalf("got %d lines, want %d", len(lines), len(expected))
	}
	for i := range lines {
		if lines[i] != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], expected[i])
		}
	}
}

func TestTargetOutputStable(t *testing.T) {
	out, code := runPraxis(t, "target", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "★\n" {
		t.Errorf("target output: got %q, want \"★\\n\"", out)
	}
}

func TestVerifyOutputStable(t *testing.T) {
	out, code := runPraxis(t, "verify", "delete_word_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "buffer\n" {
		t.Errorf("verify output: got %q, want \"buffer\\n\"", out)
	}
}

func TestRecordStats(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "record", "motion_rush", "4", "380")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if _, err := os.Stat(filepath.Join(d, "praxis", "stats.json")); err != nil {
		t.Fatalf("stats.json not created: %v", err)
	}
	_, code = runPraxis(t, "record", "motion_rush", "2", "180")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	out, code := runPraxis(t, "stats", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "Attempts: 2") {
		t.Errorf("expected Attempts: 2, got: %s", out)
	}
	if !strings.Contains(out, "Best Moves: 2") {
		t.Errorf("expected Best Moves: 2, got: %s", out)
	}
	if !strings.Contains(out, "Best Time: 180ms") {
		t.Errorf("expected Best Time: 180ms, got: %s", out)
	}
	if !strings.Contains(out, "Mastery: Learning") {
		t.Errorf("expected Mastery: Learning, got: %s", out)
	}
}

func TestStatsCommand(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	runPraxis(t, "record", "motion_rush", "4", "380")
	out, code := runPraxis(t, "stats", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "Attempts: 1") {
		t.Errorf("expected Attempts: 1, got: %s", out)
	}
	if !strings.Contains(out, "Completions: 1") {
		t.Errorf("expected Completions: 1, got: %s", out)
	}
	if !strings.Contains(out, "Best Moves: 4") {
		t.Errorf("expected Best Moves: 4, got: %s", out)
	}
	if !strings.Contains(out, "Best Time: 380ms") {
		t.Errorf("expected Best Time: 380ms, got: %s", out)
	}
	if !strings.Contains(out, "Mastery: Learning") {
		t.Errorf("expected Mastery: Learning, got: %s", out)
	}
}

func TestStatsSummary(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	runPraxis(t, "record", "motion_rush", "3", "200")
	runPraxis(t, "record", "grid_rush", "5", "300")
	runPraxis(t, "record", "motion_rush", "2", "150")
	out, code := runPraxis(t, "stats")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "2/41") {
		t.Errorf("expected 2/41 completed, got: %s", out)
	}
	if !strings.Contains(out, "Total Attempts: 3") {
		t.Errorf("expected Total Attempts: 3, got: %s", out)
	}
}

func TestStatsUnknownChallenge(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "stats", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
