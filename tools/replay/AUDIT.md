# Replay Audit

## What Replay Proves

Running `tools/replay/replay.sh` validates:

- **Challenge solvability.** Every challenge completes successfully under automated headless Neovim.
- **Validator correctness.** Both cursor and buffer validators produce "Success" for valid completions.
- **Lua–Go integration.** The plugin frontend correctly communicates with the Go backend via system calls.
- **Completion paths.** CursorMoved and TextChanged events both trigger the correct success logic.
- **UTF-8 normalization.** Multi-byte content does not break cursor positioning or validation.

A passing replay means every challenge is mechanically solvable as of this commit.

## What Replay Does Not Prove

- **Pedagogical quality.** A challenge may pass replay while teaching a poor or incorrect technique.
- **Concept uniqueness.** Two challenges may pass replay while teaching identical concepts.
- **Curriculum balance.** All 41 challenges may pass while omitting entire categories of Vim techniques.
- **Difficulty progression.** Replay does not measure whether challenges increase in complexity.
- **Real-user behavior.** Replay uses automated input sequences. It cannot simulate genuine practice.

## Coverage

| Metric | Current |
|---|---|
| Challenges tested | 41/41 |
| Cursor challenges | 20/20 |
| Buffer challenges | 21/21 |
| Pass rate | 100% |

## Adding a New Challenge

1. Write the challenge content in `internal/content/content.go`
2. Add a replay sequence to `tools/replay/`
3. Run `tools/replay/replay.sh`
4. Verify PASS for the new challenge
5. Verify all existing challenges still PASS
