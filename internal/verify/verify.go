// Package verify replays a learner's key log in a clean, config-free Neovim
// and returns a verdict. No progress is written here.
package verify

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/divijg19/Praxis/internal/content"
)

//go:embed replay.lua
var scripts embed.FS

type Input struct {
	Keys    []string `json:"keys"`
	Tainted bool     `json:"tainted,omitempty"`
}

type Result struct {
	Verified   bool     `json:"verified"`
	Reason     string   `json:"reason"`
	Keys       int      `json:"keys"`
	Edits      int      `json:"edits"`
	Par        int      `json:"par"`
	FinalLines []string `json:"final_lines,omitempty"`
}

// namedKeys is immutable and shared by every replay; keep it off the hot path.
var namedKeys = map[string]string{
	"<Esc>": "\x1b", "<CR>": "\r", "<BS>": "\x08", "<Tab>": "\t",
	"<Left>": "\x1b[D", "<Right>": "\x1b[C", "<Up>": "\x1b[A", "<Down>": "\x1b[B",
	"<Home>": "\x1b[H", "<End>": "\x1b[F", "<Del>": "\x1b[3~",
}

func Run(id string, in Input) (Result, error) {
	if in.Tainted {
		return Result{Reason: "input was tainted by non-typed keys", Par: content.ParFor(id)}, nil
	}
	if len(in.Keys) == 0 {
		return Result{Reason: "empty key log", Par: content.ParFor(id)}, nil
	}
	if len(in.Keys) > 10000 {
		return Result{Reason: "key log too large", Par: content.ParFor(id)}, nil
	}
	keys := make([]string, len(in.Keys))
	for i, k := range in.Keys {
		if normalized, ok := namedKeys[k]; ok {
			keys[i] = normalized
			continue
		}
		if k == "" || len(k) > 16 {
			return Result{Reason: "invalid key in log", Par: content.ParFor(id)}, nil
		}
		// Angle brackets are never allowed: joining raw keys before
		// nvim_replace_termcodes would let `<`,`C`,`h`,… reassemble into a
		// termcode such as <Char-58> and execute an ex command.
		r := []rune(k)
		if len(r) != 1 || r[0] == ':' || r[0] == '<' || r[0] == '>' || r[0] == 'Q' || r[0] == 'q' ||
			(r[0] < 0x20 && r[0] != '\x1b' && r[0] != '\r' && r[0] != '\x08' && r[0] != '\t') {
			return Result{Reason: "invalid key in log", Par: content.ParFor(id)}, nil
		}
		keys[i] = k
	}
	d, ok := content.DescriptionFor(id)
	if !ok {
		return Result{}, fmt.Errorf("unknown challenge %q", id)
	}
	descJSON, err := json.Marshal(d)
	if err != nil {
		return Result{}, err
	}
	keysJSON, err := json.Marshal(keys)
	if err != nil {
		return Result{}, err
	}
	tmp, err := os.MkdirTemp("", "praxis-verify-")
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	keysPath := filepath.Join(tmp, "keys.json")
	if err := os.WriteFile(keysPath, keysJSON, 0o600); err != nil {
		return Result{}, err
	}

	script, err := scripts.ReadFile("replay.lua")
	if err != nil {
		return Result{}, err
	}
	scriptPath := filepath.Join(tmp, "replay.lua")
	if err := os.WriteFile(scriptPath, script, 0o600); err != nil {
		return Result{}, err
	}
	outPath := filepath.Join(tmp, "verdict.json")

	env := append(os.Environ(),
		"PRAXIS_VERIFY_DESC="+string(descJSON),
		"PRAXIS_VERIFY_KEYS="+keysPath,
		"PRAXIS_VERIFY_OUT="+outPath,
		"XDG_CONFIG_HOME="+filepath.Join(tmp, "cfg"),
		"XDG_DATA_HOME="+filepath.Join(tmp, "data"),
		"XDG_STATE_HOME="+filepath.Join(tmp, "state"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nvim", "--clean", "-u", "NONE", "-i", "NONE", "--noplugin", "-n", "--headless", "-l", scriptPath)
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		return Result{Reason: "clean replay failed"}, fmt.Errorf("clean replay: %w: %s", err, out)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return Result{}, fmt.Errorf("read verdict: %w", err)
	}
	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return Result{}, fmt.Errorf("decode verdict: %w", err)
	}
	return result, nil
}
