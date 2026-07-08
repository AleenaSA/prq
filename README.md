# prq

A minimalist CLI to see your GitHub PR queue at a glance.

```
$ prq

 Pending Reviews (2)

 ●  org/repo#42  "Add auth middleware"  @alice  +120/-30  2h ago
 ●  org/repo#55  "Fix rate limiter"     @bob    +8/-3     5h ago

 My PRs (3)

 ✓ ● ◯  org/repo#38  "Migrate database"   +45/-12   3d ago
 ✗ ● ⊘  org/repo#41  "Add caching layer"  +200/-50  1d ago
 ◌ ● ?  org/repo#43  "Update deps"        +15/-10   2h ago
```

## Status symbols

For your own PRs, the three symbols represent (left to right):

| Position | Symbol       | Meaning           |
| -------- | ------------ | ----------------- |
| CI       | `✓`          | Passing           |
| CI       | `✗`          | Failing           |
| CI       | `◌`          | Pending           |
| CI       | `–`          | No checks         |
| Review   | `●` (green)  | Approved          |
| Review   | `●` (yellow) | Pending review    |
| Review   | `●` (red)    | Changes requested |
| Merge    | `◯`          | Clean, mergeable  |
| Merge    | `⊘`          | Has conflicts     |
| Merge    | `?`          | Unknown           |

## Install

```bash
# Homebrew
brew install AleenaSA/tap/prq

# Go
go install github.com/AleenaSA/prq@latest

# Binary (see releases)
curl -sSL https://github.com/AleenaSA/prq/releases/latest/download/prq_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz | tar xz
```

## Auth

`prq` looks for a GitHub token in this order:

1. `gh auth token` (reuses your GitHub CLI session)
2. `GITHUB_TOKEN` environment variable

## Usage

```
prq              # compact colored output
prq --table      # aligned table columns
prq --no-color   # no ANSI escape codes
```

Also respects the `NO_COLOR` environment variable.

## License

MIT
