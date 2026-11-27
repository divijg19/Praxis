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
