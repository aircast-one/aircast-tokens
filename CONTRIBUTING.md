# Contributing

aircast-tokens is the single source of truth for the Aircast colour palette and radius
scale. `tokens.css` is hand-written — nothing here generates it. Everything else
(`palette.go`, `material.go`, `cmd/aircast-tokens`) derives *from* it for consumers that
aren't CSS: the Go embed, the Material 3 role mapping used by the Android/iOS/Qt heads,
and the CLI that renders those formats.

## Changing a token

1. Edit `tokens.css` directly.
2. Run `go test ./...` — `palette_test.go` checks the OKLCH→sRGB conversion against the
   hex values Chrome actually draws, so a bad edit to the conversion math fails loudly.
3. If the change should reach native heads, spot-check the derived output:
   ```
   go run ./cmd/aircast-tokens -format material-css
   go run ./cmd/aircast-tokens -format kotlin -package one.aircast.map
   go run ./cmd/aircast-tokens -format swift
   go run ./cmd/aircast-tokens -format qt
   ```
4. Commit and push to `master` — this repo has no CI and no branch/PR workflow, consumers
   pin by commit SHA.
5. Bump the pin in each consumer:
   - `aircast-flasher`: `"@aircast/tokens": "github:aircast-one/aircast-tokens#<commit>"`
     in `package.json`.
   - `aircastd`: `go get -tool github.com/aircast-one/aircast-tokens/cmd/aircast-tokens@<commit>`.
   - Native heads (Android/iOS/Qt): re-run the relevant `go run ... @<commit>` command above
     and commit the regenerated file.

## Scope

- `tokens.css` — the palette and radius scale only. Fonts, type scale, and the `dark`
  variant selector stay in each consuming app.
- `material.go` — maps those tokens onto Material 3 roles once, shared by every native
  format. Don't duplicate role logic in a per-format renderer; add it here.
- `palette.go` — OKLCH math and parsing. Changes here are the highest-risk part of this
  repo since every format depends on them — always run `go test ./...` before pushing.

## Adding a new output format

Add a case to `render()` in `cmd/aircast-tokens/main.go` and a corresponding function
alongside `MaterialCSS()` in `material.go`. Keep the CLI a thin dispatcher; put the actual
rendering logic in the package so it stays testable without the CLI.

Requires Go 1.26.3+ (pinned in `go.mod`).
