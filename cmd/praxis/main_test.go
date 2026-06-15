package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/divijg19/Praxis/internal/content"
	"github.com/divijg19/Praxis/internal/stats"
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
	if !strings.Contains(out, "Completions: 2") {
		t.Errorf("expected Completions: 2, got: %s", out)
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
	// direct record() call bypasses normal attempt tracking;
	// this state (Completions>0, Attempts=0) cannot occur through the Neovim frontend
	if !strings.Contains(out, "Confidence: —") {
		t.Errorf("expected Confidence: — (no attempts), got: %s", out)
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
	if !strings.Contains(out, "Completions: 1") {
		t.Errorf("expected Completions: 1, got: %s", out)
	}
	if !strings.Contains(out, "Success Rate: 0%") {
		t.Errorf("expected success rate 0%% in output, got: %s", out)
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
	// direct record() call bypasses normal attempt tracking;
	// this state (Completions>0, Attempts=0) cannot occur through the Neovim frontend
	if !strings.Contains(out, "Confidence: —") {
		t.Errorf("expected Confidence: — (no attempts), got: %s", out)
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
	if !strings.Contains(out, "Total Attempts: 0") {
		t.Errorf("expected Total Attempts: 0, got: %s", out)
	}
	if !strings.Contains(out, "Mastery:") {
		t.Errorf("expected Mastery header, got: %s", out)
	}
	if !strings.Contains(out, "Unseen: 39") {
		t.Errorf("expected Unseen: 39, got: %s", out)
	}
	if !strings.Contains(out, "Learning: 2") {
		t.Errorf("expected Learning: 2, got: %s", out)
	}
	if !strings.Contains(out, "Highest Tier: Learning") {
		t.Errorf("expected Highest Tier: Learning, got: %s", out)
	}
	if !strings.Contains(out, "Next Challenge:") {
		t.Errorf("expected Next Challenge header, got: %s", out)
	}
	if !strings.Contains(out, "motion_rush") {
		t.Errorf("expected motion_rush as next challenge (still Learning at 2 completions), got: %s", out)
	}
	if strings.Contains(out, "Recommended Review:") {
		t.Errorf("unexpected Recommended Review section (no Practiced+ challenges), got: %s", out)
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

func TestAttemptCommand(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "attempt", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	out, code := runPraxis(t, "stats", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "Attempts: 1") {
		t.Errorf("expected Attempts: 1, got: %s", out)
	}
	if !strings.Contains(out, "Completions: 0") {
		t.Errorf("expected Completions: 0, got: %s", out)
	}
}

func TestAttemptWithRecord(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	runPraxis(t, "attempt", "motion_rush")
	runPraxis(t, "record", "motion_rush", "2", "180")
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
	if !strings.Contains(out, "Confidence: High") {
		t.Errorf("expected Confidence: High, got: %s", out)
	}
}

func TestAttemptUnknown(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "attempt", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestStatsCommandConfidenceLevels(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	m := map[string]stats.Stats{
		"motion_rush": {Attempts: 0, Completions: 0},
		"grid_rush":   {Attempts: 10, Completions: 5},
		"find_hunter": {Attempts: 10, Completions: 6},
		"word_hunter": {Attempts: 10, Completions: 8},
	}
	if err := stats.Save(m); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		id   string
		want string
	}{
		{"motion_rush", "Confidence: —"},
		{"grid_rush", "Confidence: Low"},
		{"find_hunter", "Confidence: Medium"},
		{"word_hunter", "Confidence: High"},
	}
	for _, tc := range cases {
		out, _ := runPraxis(t, "stats", tc.id)
		if !strings.Contains(out, tc.want) {
			t.Fatalf("praxis stats %s: expected %q in output, got:\n%s",
				tc.id, tc.want, out)
		}
	}
}

func TestNextCommand(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)

	out, code := runPraxis(t, "next")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "motion_rush\n" {
		t.Errorf("expected 'motion_rush', got %q", out)
	}
}

func TestNextCommandAfterCompletion(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)

	m := map[string]stats.Stats{
		"motion_rush": {Attempts: 3, Completions: 3},
		"grid_rush":   {Attempts: 3, Completions: 3},
	}
	if err := stats.Save(m); err != nil {
		t.Fatal(err)
	}

	out, code := runPraxis(t, "next")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "find_hunter\n" {
		t.Errorf("expected 'find_hunter', got %q", out)
	}
}

func TestNextCommandComplete(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)

	m := make(map[string]stats.Stats)
	for _, c := range content.All() {
		m[c.ID] = stats.Stats{Attempts: 10, Completions: 10}
	}
	if err := stats.Save(m); err != nil {
		t.Fatal(err)
	}

	out, code := runPraxis(t, "next")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestStageCommand(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)

	out, code := runPraxis(t, "stage", "motion_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "Movement\n" {
		t.Errorf("expected 'Movement', got %q", out)
	}

	out, code = runPraxis(t, "stage", "find_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if out != "Search\n" {
		t.Errorf("expected 'Search', got %q", out)
	}
}

func TestStageCommandUnknown(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)

	_, code := runPraxis(t, "stage", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1 for unknown challenge, got %d", code)
	}
}
