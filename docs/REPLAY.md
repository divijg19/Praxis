# Replay Verification

## Purpose

Replay verification runs every challenge through a real Neovim session, ensuring that challenge content, results, and validator behavior are consistent.

This is the **canonical integration verification** for Praxis. It is the single source of truth for end-to-end correctness: it exercises the entire pipeline from Go binary to Lua frontend to Neovim buffer state.

## Running

```bash
tools/replay/replay.sh
```

The script:
1. Builds the Go binary to `/tmp/praxis`
2. Runs Neovim headless with `tools/replay/replay.lua`
3. Prints PASS/FAIL for each challenge
4. Reports summary: `ALL 41/41 REPLAY TESTS PASS`

## Interpreting Results

Each challenge reports `PASS <id>` or `FAIL <id>`.

- A **cursor challenge** passes when navigating to the target character completes successfully
- A **buffer challenge** passes when setting the buffer content to the result matches exactly
- The **UTF-8 challenge** passes when byte-to-character normalization correctly identifies the star at bytecol=9, charcol=6

The summary line shows the final count. All 41 must pass.

## Adding a New Challenge

When adding a new challenge:

1. Add its ID to the appropriate list in `tools/replay/replay.lua`:
   - Cursor challenges go in `cursor_ids`
   - Buffer challenges go in `buffer_ids`
2. Run the replay to verify
3. Commit the updated `replay.lua`

## Limitations

- Buffer challenges are verified by directly setting the buffer to the result, not by simulating keystrokes
- This means replay tests validate that content+result pairs are consistent, not that the teaching technique works
- Full technique validation requires manual testing within Neovim
