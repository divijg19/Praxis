# `Praxis`

Praxis teaches Vim through short, self-contained challenges: Tutorials, Training, and Trials, all inside Neovim. It is deliberate practice, not a course, and is meant to become a place you return to in order to improve.

## Install

Install the binary (make sure `$(go env GOPATH)/bin` is on `PATH`):

```bash
go install github.com/divijg19/Praxis/cmd/praxis@latest
```

Add this repository's `nvim/` directory to Neovim's `runtimepath`, for example:

```bash
ln -sfn /path/to/Praxis/nvim "${XDG_CONFIG_HOME:-$HOME/.config}/nvim"
```

> **Upgrading from v0.3.x:** progress files are not migrated. `praxis reset --yes` starts fresh; a backup may be kept as `progress.json.corrupt` until you reset.

## Quick Start

Open Neovim and run `:praxis`.

Inside a challenge, `[<leader>r]` retries. Training and Trial challenges also
offer `[<leader>h]` for an opt-in hint: the first press previews the cost, the
second reveals the text. Revealed hints add key-equivalents to that attempt's
marks and are cleared on retry.

Run `praxis report` for a private, read-only summary of verified practice. It
writes nothing and produces no score or ranking.

## Documentation

- [CURRICULUM.md](docs/CURRICULUM.md): what Praxis teaches
- [ARCHITECTURE.md](docs/ARCHITECTURE.md): how it is built
- [REFERENCE.md](docs/REFERENCE.md): challenge model, CLI, stats, persistence contract
- [TESTING.md](docs/TESTING.md): how it is verified
- [TRACES_AND_PAR.md](docs/TRACES_AND_PAR.md): canonical solutions and verified par
- [SCORING.md](docs/SCORING.md): marks from verified runs
- [DESIGN_PRINCIPLES.md](docs/DESIGN_PRINCIPLES.md): the principles that guide curriculum decisions

## Verify

```bash
go run ./tools/verify
```

## Contribute

Contributions follow the standard flow: branch, implement, run `go run ./tools/verify`, and open a pull request.
