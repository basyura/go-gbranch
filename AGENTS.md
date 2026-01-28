# Repository Guidelines

## 計画

- 修正を始める前に計画をマークダウンファイルで .plans フォルダ配下に生成してください。
- 計画のファイル名は連番とし、1つ目を 001 始まりとして修正にあった適切なファイル名としてください。
- 計画のフォーマットは `001-chop-by-width.md` を参照すること。
- 具体的なファイル編集をする前に、修正案を提示すること。

## Project Structure & Module Organization

- main_gbranch.go contains all application logic for the gbranch CLI.
- go.mod and go.sum define the Go module and dependencies.
- gbranch.exe is a built binary artifact (do not edit by hand).
- There are no separate cmd/, internal/, or test/ directories in this repository.

## Build, Test, and Development Commands

- go build -o gbranch.exe . builds the Windows binary in the repo root.
- go run . runs the CLI from source for quick checks.
- go test ./... runs all Go tests (currently none are present).
- gofmt -w main_gbranch.go formats the source file.

## Coding Style & Naming Conventions

- Follow standard Go conventions (tabs for indentation, CamelCase for types, lowerCamel for variables and functions).
- Keep functions small and focused; prefer clear names over comments.
- Avoid introducing new dependencies unless needed for CLI output or Git parsing.
- Keep line output deterministic to preserve CLI behavior and parsing expectations.

## Testing Guidelines

- No test framework is currently used; add *_test.go files under the repo root if you introduce tests.
- Name tests TestXxx and run with go test ./....
- If you add tests, keep them fast and independent of network access.

## Commit & Pull Request Guidelines

- Commit messages in this repo are short, Japanese summaries (example: 出力フォーマット変更).
- Keep commits focused on a single change; avoid mixing refactors and behavior changes.
- Pull requests should include a summary, rationale, and a sample output snippet for CLI changes.

## Configuration & Usage Notes

- Environment variables control output: $GBRANCH_SYMBOL for the current branch marker and $GBRANCH_FG for color (example values: red, hiblue).
- If you change output formatting, document the exact format and update any related examples.
