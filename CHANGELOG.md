# Changelog

All notable changes to Illygen are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [v0.1.1] — 2026-02-26

### Fixed

- `engine.Run` no longer panics when passed a `nil` Context — treated as empty
- `KnowledgeStore.Add` now returns an error for empty `id` or `domain` instead of silently storing bad data
- Cycle detection threshold changed from magic number `100` to named constant `maxVisits = 50` in internal runtime
- Error message when a node routes to an unknown node now includes actionable hint: `"did you call flow.Add()?"`
- Addressed test and build issues introduced during the sync and example additions; ensure module builds and examples run as tests

### Added

- `Context.Bool(key)` — convenience accessor for bool values
- `Context.Int(key)` — convenience accessor for int values
- `Context.Float(key)` — convenience accessor for float64 values
- `NewNode` panics immediately with a clear message if called with empty `id` or nil `NodeFunc`
- `Flow.Add` panics immediately with a clear message if called with a nil node
- Full unit test suite — 33 tests covering Node, Flow, Engine, Context, KnowledgeStore, and knowledge injection
- `CHANGELOG.md` — this file
- `examples/conversational` — interactive conversational demo showing how to use `KnowledgeStore` and `Engine` (REPL)
- `example_pkg_test.go` — `Example` tests so pkg.go.dev can display runnable examples for `Flow`, `KnowledgeStore`, and `Engine`+`Knowledge`
- GitHub Actions CI workflow (`.github/workflows/ci.yml`) to run `gofmt` check, `go vet`, and `go test` on PRs and pushes

### Changed

- Godoc polish across all public types: `Node`, `Flow`, `Engine`, `Context`, `Result`, `KnowledgeUnit`, `KnowledgeStore`
- `Result` docs now clearly explain the priority of `Next` vs graph Links
- `Flow.Link` docs clarify that duplicate links are silently ignored (first call wins)
- `Engine.Run` docs clarify nil Context behaviour and concurrent safety
- Internal `executor.go` algorithm comment updated to match actual implementation
- Added more and clearer examples to package documentation to improve discoverability on pkg.go.dev

---

## [v0.1.0] — 2026-02-25

### Added

- Initial release
- `Node` — define reasoning units as plain Go functions via `NodeFunc`
- `Flow` — connect nodes into a weighted reasoning pipeline
- `Engine` — runs flows, returns `Result` with `Value` and `Confidence`
- `Context` — key-value map carrying data through flow execution
- `Result` — output of a node; `Value`, `Confidence`, `Next`
- `KnowledgeStore` — structured knowledge store, domain-scoped
- `KnowledgeUnit` — atomic knowledge unit with facts and weight
- `illygen.Knowledge(ctx)` — query the store from inside a NodeFunc
- `examples/intent/` — intent detection example (input → action)
- Frozen public API — nothing in `internal/` is accessible to users
