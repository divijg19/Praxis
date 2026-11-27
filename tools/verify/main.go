// Command verify runs the single Praxis certification gate.
//
// It is the only release gate: build, vet, format, lint, unit/integrity tests,
// replay, and the learner journey. External binaries (golangci-lint, nvim) are
// resolved from PATH rather than pinned as module dependencies.
package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const harnessTimeout = 2 * time.Minute

func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above working directory")
		}
		dir = parent
	}
}

type verifier struct {
	root    string
	binDir  string
	dataDir string
}

func (v *verifier) run(env []string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = v.root
	if env != nil {
		cmd.Env = env
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func (v *verifier) runHarness(env []string, name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), harnessTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = v.root
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%s timed out after %s", name, harnessTimeout)
		}
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func (v *verifier) checkFormat() error {
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = v.root
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gofmt -l .: %w", err)
	}
	if files := strings.TrimSpace(out.String()); files != "" {
		return fmt.Errorf("unformatted files:\n%s", files)
	}
	fmt.Println("PASS: clean")
	return nil
}

// harnessEnv points the Lua harnesses at the freshly built binary and an
// isolated data directory. Existing PATH/PRAXIS_* entries are replaced rather
// than duplicated; PATH keeps its original entries after the build directory
// so helper tools remain available.
func harnessEnv(binDir, binPath, dataDir string) []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "PATH=") ||
			strings.HasPrefix(entry, "PRAXIS_BIN=") ||
			strings.HasPrefix(entry, "PRAXIS_VERIFY_DATA=") ||
			strings.HasPrefix(entry, "PRAXIS_DEV=") {
			continue
		}
		env = append(env, entry)
	}
	return append(env,
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"PRAXIS_BIN="+binPath,
		"PRAXIS_VERIFY_DATA="+dataDir,
		"PRAXIS_DEV=1",
	)
}

func certify() error {
	root, err := findRoot()
	if err != nil {
		return err
	}
	for _, tool := range []string{"go", "gofmt", "golangci-lint", "nvim"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("required tool %q not found on PATH: %w", tool, err)
		}
	}

	tmp, err := os.MkdirTemp("", "praxis-verify-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	v := &verifier{
		root:    root,
		binDir:  filepath.Join(tmp, "bin"),
		dataDir: filepath.Join(tmp, "data"),
	}
	for _, dir := range []string{v.binDir, v.dataDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	binPath := filepath.Join(v.binDir, "praxis")
	env := harnessEnv(v.binDir, binPath, v.dataDir)

	steps := []struct {
		name string
		run  func() error
	}{
		{"go build ./...", func() error { return v.run(nil, "go", "build", "./...") }},
		{"go vet ./...", func() error { return v.run(nil, "go", "vet", "./...") }},
		{"gofmt -l .", v.checkFormat},
		{"golangci-lint run", func() error { return v.run(nil, "golangci-lint", "run") }},
		{"go test ./... -count=1", func() error { return v.run(nil, "go", "test", "./...", "-count=1") }},
		{"build praxis binary", func() error { return v.run(nil, "go", "build", "-o", binPath, "./cmd/praxis") }},
		{"tools/replay/replay.lua", func() error {
			return v.runHarness(env, "nvim", "--headless", "-u", "NONE", "-l", filepath.Join(root, "tools", "replay", "replay.lua"))
		}},
		{"tools/journey/journey.lua", func() error {
			return v.runHarness(env, "nvim", "--headless", "-u", "NONE", "-l", filepath.Join(root, "tools", "journey", "journey.lua"))
		}},
	}

	for _, step := range steps {
		fmt.Printf("=== %s ===\n", step.name)
		if err := step.run(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
		fmt.Println()
	}
	return nil
}

func main() {
	if err := certify(); err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("=== PRAXIS CERTIFIED: ALL CHECKS PASS ===")
}
