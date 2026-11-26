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

## Quick Start

Open Neovim and run `:Praxis`.

## Documentation

- [CURRICULUM.md](docs/CURRICULUM.md): what Praxis teaches
- [ARCHITECTURE.md](docs/ARCHITECTURE.md): how it is built
- [REFERENCE.md](docs/REFERENCE.md): challenge model, CLI, stats, persistence contract
- [TESTING.md](docs/TESTING.md): how it is verified
- [DESIGN_PRINCIPLES.md](docs/DESIGN_PRINCIPLES.md): the principles that guide curriculum decisions

## Verify

```bash
go run ./tools/verify
```

## Contribute

Contributions follow the standard flow: branch, implement, run `go run ./tools/verify`, and open a pull request.
