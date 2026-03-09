# Change: Improve Conversation Discovery UX

## Why
Shiftlog already stores useful metadata and search results, but several obvious workflows are still awkward:
- Web resume ignores the stored agent and always launches Claude
- The web UI cannot search stored conversations
- `shiftlog list` exposes only a minimal text dump
- `shiftlog search` cannot emit JSON or jump directly into a matching conversation
- Branch overview counts can silently undercount long-lived branches

## What Changes
- Make web resume launch the agent associated with the stored conversation
- Add `/api/search` and a simple web search flow
- Extend `shiftlog list` with filtering, limits, and JSON output
- Extend `shiftlog search` with JSON output, top-result handoff, and highlighted matches
- Count branch conversations across full branch history and add a branch-history load-more control in the web UI

## Impact
- Affected specs: `cli`, `web-visualization`
- Affected code: `cmd/list.go`, `cmd/search.go`, `cmd/show.go`, `internal/git/repo.go`, `internal/web/handlers.go`, `internal/web/server.go`, `internal/web/static/index.html`
- Affected tests: acceptance CLI tests and web handler tests
