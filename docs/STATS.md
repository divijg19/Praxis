# Progress Tracking

## Overview

Praxis now remembers your practice. Every challenge completion records:

- **Attempts** — total completions for this challenge
- **Best Moves** — fewest moves across all completions
- **Best Time** — fastest completion time in milliseconds
- **Last Played** — date of most recent completion

## Storage

Stats are stored in:

```
~/.local/share/praxis/stats.json
```

Or `$XDG_DATA_HOME/praxis/stats.json` if the environment variable is set.

The file is simple JSON, human-readable and editable:

```json
{
  "motion_rush": {
    "attempts": 12,
    "completions": 10,
    "best_moves": 2,
    "best_time_ms": 180,
    "last_played": "2026-06-09"
  }
}
```

No database. No migrations. If the file is missing or corrupt, Praxis silently starts fresh.

## CLI

### Per-challenge stats

```bash
praxis stats motion_rush
```

Output:

```
Attempts: 12
Completions: 10
Best Moves: 2
Best Time: 180ms
```

### Summary

```bash
praxis stats
```

Output:

```
Challenges Completed: 31/41
Total Attempts: 245
```

### Recording (internal)

Used by the Neovim frontend on challenge completion:

```bash
praxis record motion_rush 4 380
```

Silent on success. Not intended for manual use.

## Neovim Integration

When a challenge is completed, the result screen now shows best stats:

```
Success

Moves: 4
Time: 380ms

Best Moves: 3
Best Time: 250ms

Press r to replay
```

## Design Decisions

- Stats are recorded only on successful completion
- Best values track minima (fewest moves, fastest time)
- First completion always sets the best value
- Subsequent completions only update if better
- No per-attempt history — only the best is preserved

## Future

This is the data foundation for:

- v0.1.2 Practice Sessions (filter challenges by completion status)
- v0.1.3 Mastery Ratings (score based on best performance)
- v0.1.4 Challenge Integrity (identify outliers)
