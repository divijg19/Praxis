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
	binPath = filepath.Join(os.TempDir(), "praxis_test_v040")
	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "build failed:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func runPraxis(t *testing.T, args ...string) (string, int) { return runPraxisStdin(t, "", args...) }

func runPraxisStdin(t *testing.T, stdin string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr), exitErr.ExitCode()
		}
		t.Fatal(err)
	}
	return string(out), 0
}

func canonicalInput(t *testing.T, id string) string {
	t.Helper()
	tr, ok := content.TraceFor(id)
	if !ok {
		t.Fatalf("no trace for %s", id)
	}
	var keys []string
	for _, r := range tr.Keys {
		keys = append(keys, string(r))
	}
	b, err := json.Marshal(map[string]any{"keys": keys})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCatalogOutputStable(t *testing.T) {
	out, code := runPraxis(t, "catalog")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	all := content.All()
	if len(lines) != len(all) {
		t.Fatalf("got %d names, want %d", len(lines), len(all))
	}
	for i, c := range all {
		if lines[i] != c.Name {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], c.Name)
		}
	}
}

func TestCatalogFlags(t *testing.T) {
	ids, code := runPraxis(t, "catalog", "--ids")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if len(strings.Split(strings.TrimRight(ids, "\n"), "\n")) != len(content.All()) {
		t.Fatal("catalog --ids count mismatch")
	}
	core, code := runPraxis(t, "catalog", "--core")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if len(strings.Split(strings.TrimRight(core, "\n"), "\n")) != len(content.CoreTutorialIDs()) {
		t.Fatal("catalog --core count mismatch")
	}
}

func TestCatalogJSONHasNoSpoilers(t *testing.T) {
	out, code := runPraxis(t, "catalog", "--json")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	for _, banned := range []string{`"content"`, `"result"`, `"target"`, `"keys"`, `"hints"`, `"start_row"`} {
		if strings.Contains(out, banned) {
			t.Errorf("catalog --json leaked %s", banned)
		}
	}
}

func TestDescribeFields(t *testing.T) {
	out, code := runPraxis(t, "describe", "grid_rush")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	var d struct {
		StartRow int `json:"start_row"`
		Par      int `json:"par"`
	}
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	if d.StartRow < 1 || d.Par < 1 {
		t.Errorf("describe fields = %+v", d)
	}
}

func TestTraceGate(t *testing.T) {
	if _, code := runPraxis(t, "trace", "grid_rush"); code != 1 {
		t.Error("trace must be dev-gated")
	}
	t.Setenv("PRAXIS_DEV", "1")
	out, code := runPraxis(t, "trace", "grid_rush")
	if code != 0 {
		t.Fatalf("exit code %d: %s", code, out)
	}
	var tr struct {
		Keys string `json:"keys"`
	}
	if err := json.Unmarshal([]byte(out), &tr); err != nil || tr.Keys == "" {
		t.Errorf("bad trace %q: %v", out, err)
	}
	if _, code := runPraxis(t, "trace", "does_not_exist"); code != 1 {
		t.Error("unknown id must exit 1")
	}
}

func TestStartAndComplete(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_ = os.MkdirAll(filepath.Join(d, "praxis"), 0o755)
	_ = os.WriteFile(filepath.Join(d, "praxis", "progress.json"), []byte(`{"format":1,"entries":{}}`), 0o644)
	if _, code := runPraxis(t, "start", "grid_rush"); code != 0 {
		t.Fatalf("start exit %d", code)
	}
	out, code := runPraxisStdin(t, canonicalInput(t, "grid_rush"), "complete", "grid_rush")
	if code != 0 {
		t.Fatalf("complete exit %d: %s", code, out)
	}
	var verdict struct {
		Verified    bool `json:"verified"`
		Marks       int  `json:"marks"`
		Completions int  `json:"completions"`
	}
	if err := json.Unmarshal([]byte(out), &verdict); err != nil {
		t.Fatal(err)
	}
	if !verdict.Verified || verdict.Marks < 1 || verdict.Completions != 1 {
		t.Errorf("bad verdict: %+v", verdict)
	}
}

func TestTaintedCompletionRejected(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	_ = os.MkdirAll(filepath.Join(d, "praxis"), 0o755)
	_ = os.WriteFile(filepath.Join(d, "praxis", "progress.json"), []byte(`{"format":1,"entries":{}}`), 0o644)
	input := `{"keys":["j"],"tainted":true}`
	out, code := runPraxisStdin(t, input, "complete", "grid_rush")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, `"verified":false`) {
		t.Errorf("tainted input verified: %s", out)
	}
}

func TestRecordRemoved(t *testing.T) {
	if _, code := runPraxis(t, "record", "grid_rush", "1", "1"); code != 1 {
		t.Error("removed command `record` must exit 1")
	}
}

func TestReset(t *testing.T) {
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	out, code := runPraxis(t, "reset", "--yes")
	if code != 0 || !strings.Contains(out, "reset") {
		t.Errorf("reset output %q code %d", out, code)
	}
}

func freshProgress(t *testing.T) {
	t.Helper()
	d := t.TempDir()
	t.Setenv("XDG_DATA_HOME", d)
	if err := os.MkdirAll(filepath.Join(d, "praxis"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "praxis", "progress.json"), []byte(`{"format":1,"entries":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDescribeHintsCostOnly(t *testing.T) {
	out, code := runPraxis(t, "describe", "find_diw_combo")
	if code != 0 {
		t.Fatalf("exit code %d", code)
	}
	if strings.Contains(out, content.HintText("find_diw_combo", 1)) {
		t.Fatal("describe leaked hint text")
	}
	if !strings.Contains(out, `"hints":`) || !strings.Contains(out, `"cost":1`) {
		t.Fatalf("describe lacks hint costs: %s", out)
	}
}

func TestHintPreviewAndReveal(t *testing.T) {
	freshProgress(t)
	preview, code := runPraxis(t, "hint", "find_diw_combo", "1", "--preview")
	if code != 0 {
		t.Fatalf("preview exit %d: %s", code, preview)
	}
	if strings.Contains(preview, content.HintText("find_diw_combo", 1)) {
		t.Fatal("preview leaked hint text")
	}
	p, err := stats.Load()
	if err != nil {
		t.Fatal(err)
	}
	if stats.RevealedHintLevel(p, "find_diw_combo") != 0 {
		t.Fatal("preview wrote progress")
	}
	revealed, code := runPraxis(t, "hint", "find_diw_combo", "1")
	if code != 0 || !strings.Contains(revealed, content.HintText("find_diw_combo", 1)) {
		t.Fatalf("reveal failed: %s", revealed)
	}
	p, _ = stats.Load()
	if stats.RevealedHintLevel(p, "find_diw_combo") != 1 {
		t.Fatal("reveal was not persisted")
	}
}

func TestHintCostDerivedAndSpoofProof(t *testing.T) {
	freshProgress(t)
	if _, code := runPraxis(t, "start", "find_diw_combo"); code != 0 {
		t.Fatal("start failed")
	}
	if _, code := runPraxis(t, "hint", "find_diw_combo", "2"); code != 0 {
		t.Fatal("reveal failed")
	}
	payload := canonicalInput(t, "find_diw_combo")
	out, code := runPraxisStdin(t, payload, "complete", "find_diw_combo")
	if code != 0 {
		t.Fatalf("complete exit %d: %s", code, out)
	}
	var verdict struct {
		HintCost int `json:"hint_cost"`
		Marks    int `json:"marks"`
	}
	if err := json.Unmarshal([]byte(out), &verdict); err != nil {
		t.Fatal(err)
	}
	if verdict.HintCost != 3 {
		t.Fatalf("hint cost = %d, want 3", verdict.HintCost)
	}
	p, _ := stats.Load()
	if len(p.Events) != 1 || p.Events[0].HintCost != 3 {
		t.Fatalf("event hint cost wrong: %+v", p.Events)
	}
	// Send a *valid* payload containing spoofed fields. A malformed payload
	// would be silently truncated by json.Decoder and prove nothing.
	var input map[string]any
	if err := json.Unmarshal([]byte(payload), &input); err != nil {
		t.Fatal(err)
	}
	input["hint_cost"] = 0
	input["hint_level"] = 99
	input["marks"] = 3
	spoof, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	out, code = runPraxisStdin(t, string(spoof), "verify", "find_diw_combo")
	if code != 0 {
		t.Fatalf("verify exit %d: %s", code, out)
	}
	var spoofed struct {
		HintCost int `json:"hint_cost"`
		Marks    int `json:"marks"`
	}
	if err := json.Unmarshal([]byte(out), &spoofed); err != nil {
		t.Fatal(err)
	}
	if spoofed.HintCost != 3 || spoofed.Marks != 1 {
		t.Fatalf("client fields not ignored: %s", out)
	}
}

func TestHintNotRetroactiveAndClearedOnStart(t *testing.T) {
	freshProgress(t)
	if _, code := runPraxis(t, "start", "find_diw_combo"); code != 0 {
		t.Fatal("start failed")
	}
	out, code := runPraxisStdin(t, canonicalInput(t, "find_diw_combo"), "complete", "find_diw_combo")
	if code != 0 {
		t.Fatalf("first complete failed: %s", out)
	}
	var first struct {
		Marks       int `json:"marks"`
		Completions int `json:"completions"`
	}
	if err := json.Unmarshal([]byte(out), &first); err != nil {
		t.Fatal(err)
	}
	if _, code := runPraxis(t, "hint", "find_diw_combo", "1"); code != 0 {
		t.Fatal("reveal failed")
	}
	p, _ := stats.Load()
	if len(p.Events) != 1 || p.Events[0].Marks != first.Marks {
		t.Fatal("reveal changed a stored mark")
	}
	if _, code := runPraxis(t, "start", "find_diw_combo"); code != 0 {
		t.Fatal("retry start failed")
	}
	p, _ = stats.Load()
	if stats.RevealedHintLevel(p, "find_diw_combo") != 0 {
		t.Fatal("start did not clear revealed hints")
	}
}
