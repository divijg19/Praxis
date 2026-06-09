# Release Checklist

This document is the mandatory release process for every Praxis release.

## 10-Step Release Process

### 1. Format

```bash
gofmt -l .
```

Must produce no output. If it lists files, run `gofmt -w` on them.

### 2. Build

```bash
go build ./...
```

All packages must compile without errors.

### 3. Vet

```bash
go vet ./...
```

All packages must pass static analysis.

### 4. Test

```bash
go test ./...
```

All tests must pass. The final line must report `PASS` or `ok` for every package.

### 5. Replay

```bash
tools/replay/replay.sh
```

Must report `ALL 41/41 REPLAY TESTS PASS`.

### 6. Documentation

If challenge content changed:

```bash
go run scripts/generate_catalog.go > docs/CHALLENGES.md
```

Update `docs/07_RELEASES.md` with the new version row. Verify all documentation references are consistent.

### 7. Stage

```bash
git add -A
git status
```

Review staged files. Confirm no unintended changes are included.

### 8. Commit

```bash
git commit
```

Write a descriptive commit message with:
- Title: version tag and one-line summary
- Body: categorized changes (tests, docs, content, fixes)
- Discipline section: what did NOT change

### 9. Tag

```bash
git tag v0.1.<N>
```

The tag version must match the release plan.

### 10. Push and Release

```bash
git push origin <branch> v0.1.<N>
gh issue create --title "v0.1.<N>: <title>" --body "<scope>"
gh issue create --title "v0.1.<N>: verify <title>" --body "<checklist>"
```

Create two issues: a release issue describing the scope and a verification issue with a checklist.

## Additional Notes

- Archive issues are optional but recommended for milestone releases
- If replay fails, fix the issue before proceeding to stage
- If tests fail, fix the issue before proceeding to stage
- Do not skip steps. Every release follows the same process.
