package main

import (
	"encoding/json"
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

func TestCatalogOutputStable(t *testing.T) {
	out, code := runPraxis(t, "catalog")
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

func TestDescribeCommand(t *testing.T) {
	out, code := runPraxis(t, "describe", "delete_word_hunter")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	var d content.Description
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if d.ID != "delete_word_hunter" {
		t.Errorf("ID = %q, want %q", d.ID, "delete_word_hunter")
	}
	if d.Verify != "buffer" {
		t.Errorf("Verify = %q, want %q", d.Verify, "buffer")
	}
	if d.Evaluation != nil {
		t.Errorf("expected nil Evaluation for buffer challenge, got %+v", d.Evaluation)
	}
}

func TestDescribeComposite(t *testing.T) {
	out, code := runPraxis(t, "describe", "find_diw_combo")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	var d content.Description
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if d.Verify != "composite" {
		t.Errorf("Verify = %q, want %q", d.Verify, "composite")
	}
	if d.Evaluation == nil {
		t.Fatal("expected non-nil Evaluation for composite challenge")
	}
	if d.Evaluation.MaxMoves <= 0 {
		t.Errorf("MaxMoves = %d, want > 0", d.Evaluation.MaxMoves)
	}
	if d.Layer != "Training" {
		t.Errorf("Layer = %q, want %q", d.Layer, "Training")
	}
}

func TestDescribeUnknown(t *testing.T) {
	out, code := runPraxis(t, "describe", "nope")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(out, "unknown challenge") {
		t.Errorf("expected stderr to contain 'unknown challenge', got: %s", out)
	}
}

func TestHelpCommand(t *testing.T) {
	out, code := runPraxis(t, "help")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "describe") {
		t.Errorf("help output missing 'describe', got:\n%s", out)
	}
	if !strings.Contains(out, "catalog") {
		t.Errorf("help output missing 'catalog', got:\n%s", out)
	}
	if !strings.Contains(out, "attempt") {
		t.Errorf("help output missing 'attempt', got:\n%s", out)
	}
	if !strings.Contains(out, "record") {
		t.Errorf("help output missing 'record', got:\n%s", out)
	}
	if !strings.Contains(out, "next") {
		t.Errorf("help output missing 'next', got:\n%s", out)
	}
	if !strings.Contains(out, "stats") {
		t.Errorf("help output missing 'stats', got:\n%s", out)
	}
	if strings.Contains(out, "list") {
		t.Errorf("help output should not mention removed 'list' command")
	}
	if strings.Contains(out, "challenge") {
		t.Errorf("help output should not mention removed 'challenge' command")
	}
	if strings.Contains(out, "verify") {
		t.Errorf("help output should not mention removed 'verify' command")
	}
	if strings.Contains(out, "target") {
		t.Errorf("help output should not mention removed 'target' command")
	}
	if strings.Contains(out, "result") {
		t.Errorf("help output should not mention removed 'result' command")
	}
	if strings.Contains(out, "stage") {
		t.Errorf("help output should not mention removed 'stage' command")
	}
}

func TestBarePraxis(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	out, code := runPraxis(t)
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if !strings.Contains(out, "motion_rush") {
		t.Errorf("bare output missing next challenge ID, got:\n%s", out)
	}
	if !strings.Contains(out, "praxis help") {
		t.Errorf("bare output missing 'praxis help' hint, got:\n%s", out)
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
	if !strings.Contains(out, "Best Time: 180 ms") {
		t.Errorf("expected Best Time: 180 ms, got: %s", out)
	}
	if !strings.Contains(out, "Mastery: Learning") {
		t.Errorf("expected Mastery: Learning, got: %s", out)
	}
	if !strings.Contains(out, "Confidence: —") {
		t.Errorf("expected Confidence: — (no attempts), got: %s", out)
	}
}

func TestRecordUnknownChallenge(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "record", "bogus", "1", "100")
	if code != 1 {
		t.Errorf("expected exit code 1 for unknown challenge, got %d", code)
	}
}

func TestRecordRejectsNonNumeric(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	if _, code := runPraxis(t, "record", "motion_rush", "4", "380"); code != 0 {
		t.Fatalf("seed record exit code %d", code)
	}
	out, code := runPraxis(t, "record", "motion_rush", "abc", "xyz")
	if code != 1 {
		t.Errorf("expected exit code 1 for non-numeric moves/time, got %d", code)
	}
	if strings.Contains(out, "unknown") {
		t.Errorf("non-numeric args should not be treated as unknown challenge: %s", out)
	}
	loaded, _ := stats.Load()
	s := loaded["motion_rush"]
	if s.BestMoves != 4 {
		t.Errorf("stats corrupted: BestMoves = %d, want 4", s.BestMoves)
	}
	if s.BestTimeMs != 380 {
		t.Errorf("stats corrupted: BestTimeMs = %d, want 380", s.BestTimeMs)
	}
}

func TestRecordRejectsNegative(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "record", "motion_rush", "-1", "-1")
	if code != 1 {
		t.Errorf("expected exit code 1 for negative moves/time, got %d", code)
	}
	loaded, _ := stats.Load()
	s := loaded["motion_rush"]
	if s.BestMoves < 0 {
		t.Errorf("stats corrupted with negative BestMoves = %d", s.BestMoves)
	}
	if s.BestTimeMs < 0 {
		t.Errorf("stats corrupted with negative BestTimeMs = %d", s.BestTimeMs)
	}
	if s.Completions != 0 {
		t.Errorf("negative record should not create a completion, got %d", s.Completions)
	}
}

func TestRecordWrongArgCount(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_, code := runPraxis(t, "record", "motion_rush", "4")
	if code != 1 {
		t.Errorf("expected exit code 1 for missing time arg, got %d", code)
	}
}

func TestDescribeWrongArgCount(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	if _, code := runPraxis(t, "describe"); code != 1 {
		t.Errorf("expected exit code 1 for missing id, got %d", code)
	}
	if _, code := runPraxis(t, "describe", "a", "b"); code != 1 {
		t.Errorf("expected exit code 1 for extra args, got %d", code)
	}
}

func TestResetMissingFile(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	if _, err := os.Stat(filepath.Join(d, "praxis", "stats.json")); err == nil {
		t.Fatalf("precondition failed: stats.json already exists")
	}
	out, code := runPraxis(t, "reset", "--yes")
	if code != 0 {
		t.Fatalf("expected exit code 0 for reset with missing file, got %d: %s", code, out)
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
	if !strings.Contains(out, "Best Time: 380 ms") {
		t.Errorf("expected Best Time: 380 ms, got: %s", out)
	}
	if !strings.Contains(out, "Mastery: Learning") {
		t.Errorf("expected Mastery: Learning, got: %s", out)
	}
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
	total := len(content.All())
	if !strings.Contains(out, fmt.Sprintf("2/%d", total)) {
		t.Errorf("expected 2/%d completed, got: %s", total, out)
	}
	if !strings.Contains(out, "Total Attempts: 0") {
		t.Errorf("expected Total Attempts: 0, got: %s", out)
	}
	if !strings.Contains(out, "Mastery:") {
		t.Errorf("expected Mastery header, got: %s", out)
	}
	if !strings.Contains(out, fmt.Sprintf("Unseen: %d", total-2)) {
		t.Errorf("expected Unseen: %d, got: %s", total-2, out)
	}
	if !strings.Contains(out, "Learning: 2") {
		t.Errorf("expected Learning: 2, got: %s", out)
	}
	if !strings.Contains(out, "Most mastered: Learning") {
		t.Errorf("expected Most mastered: Learning, got: %s", out)
	}
	if !strings.Contains(out, "Next:") {
		t.Errorf("expected Next: header, got: %s", out)
	}
	if !strings.Contains(out, "motion_rush") {
		t.Errorf("expected motion_rush as next challenge (still Learning at 2 completions), got: %s", out)
	}
	if strings.Contains(out, "Recommended Review:") {
		t.Errorf("unexpected Recommended Review section (no Practiced+ challenges), got: %s", out)
	}
}

func TestStatsSummaryIgnoresStaleKeys(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	data := `{"motion_rush":{"attempts":1,"completions":1},"removed_challenge":{"attempts":3,"completions":3}}`
	dir := filepath.Join(d, "praxis")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "stats.json"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	out, code := runPraxis(t, "stats")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	total := len(content.All())
	if !strings.Contains(out, fmt.Sprintf("Unseen: %d", total-1)) {
		t.Errorf("expected Unseen: %d (stale key ignored), got: %s", total-1, out)
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

func TestResetCommand(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	m := map[string]stats.Stats{
		"motion_rush": {Attempts: 5, Completions: 4},
	}
	if err := stats.Save(m); err != nil {
		t.Fatal(err)
	}
	out, code := runPraxis(t, "reset", "--yes")
	if code != 0 {
		t.Fatalf("exit code %d, output: %s", code, out)
	}
	if !strings.Contains(out, "reset") {
		t.Errorf("expected reset message, got: %s", out)
	}
	loaded, _ := stats.Load()
	if len(loaded) != 0 {
		t.Errorf("stats not cleared, got %d entries", len(loaded))
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
