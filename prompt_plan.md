# TerminalizCrazy - Implementation Plan

> Roadmap for future development phases with prioritized tasks.

---

## Current Status

| Metric | Value |
|--------|-------|
| Version | Development (pre-1.0) |
| Core Features | Complete |
| Test Coverage | ~70% |
| Open Issues | See GitHub Issues |

### Implemented Features
- [x] Multi-provider AI (Ollama, Gemini, Anthropic, OpenAI)
- [x] Agent mode with task planning
- [x] Real-time collaboration with E2E encryption
- [x] Multi-pane TUI with layouts
- [x] Plugin system with hooks
- [x] Workspace persistence
- [x] Secret guard
- [x] Theme system with hot-reload
- [x] Project type detection
- [x] GDPR compliance (data retention)

---

## Priority Legend

| Priority | Meaning | Timeline |
|----------|---------|----------|
| P0 | Critical / Blocker | This sprint |
| P1 | High / Important | Next sprint |
| P2 | Medium / Nice-to-have | Future |
| P3 | Low / Backlog | When time permits |

| Effort | Meaning |
|--------|---------|
| XS | < 1 hour |
| S | 1-4 hours |
| M | 1-2 days |
| L | 3-5 days |
| XL | 1+ week |

---

## Phase 1: Testing & Stability (Current)

> **Goal:** Reach 80% test coverage, fix flaky tests, ensure CI green

| Priority | Task | Effort | Owner | Status |
|----------|------|--------|-------|--------|
| P0 | Fix flaky Windows workspace tests | M | - | 🔴 Blocked |
| P0 | Add race detection to CI | S | - | ✅ Done |
| P1 | Unit tests for `internal/ai/` (target: 85%) | L | - | 🟡 In Progress |
| P1 | Unit tests for `internal/tui/` (target: 80%) | L | - | 🟡 In Progress |
| P1 | Unit tests for `internal/collab/` (target: 80%) | M | - | ⚪ Pending |
| P1 | Integration test: Ollama provider | M | - | ⚪ Pending |
| P2 | Integration test: collaboration 2-user | L | - | ⚪ Pending |
| P2 | Performance profiling (startup < 500ms) | M | - | ⚪ Pending |
| P2 | Memory profiling (< 100MB steady state) | S | - | ⚪ Pending |

### Success Criteria
- [ ] `go test -race ./...` passes
- [ ] `go test -cover ./...` shows 80%+
- [ ] CI pipeline green on all platforms
- [ ] No flaky tests in last 10 runs

### Risks
| Risk | Mitigation |
|------|------------|
| Windows tests require specific env | Add Windows CI runner with fixtures |
| Ollama not available in CI | Use mock client for unit tests |

---

## Phase 2: AI Enhancements

> **Goal:** Full streaming support, better context management

| Priority | Task | Effort | Depends On | Status |
|----------|------|--------|------------|--------|
| P1 | Add streaming to Anthropic client | M | Phase 1 | ⚪ Pending |
| P1 | Add streaming to OpenAI client | M | Phase 1 | ⚪ Pending |
| P1 | Context window tracking (show tokens used) | S | - | ⚪ Pending |
| P2 | Conversation summarization (auto-compact) | L | Context tracking | ⚪ Pending |
| P2 | Custom system prompts per session | M | - | ⚪ Pending |
| P2 | Model-specific prompt templates | M | - | ⚪ Pending |
| P3 | Token cost estimation display | S | Context tracking | ⚪ Pending |
| P3 | Response caching (same prompt = cached) | M | - | ⚪ Pending |

### Success Criteria
- [ ] All 4 providers support streaming
- [ ] Token count visible in UI
- [ ] Context auto-summarizes at 80% capacity

### Technical Notes
- Anthropic SDK supports streaming via `WithStreaming()`
- OpenAI SDK supports streaming via `CreateChatCompletionStream()`
- Context window: Ollama varies, Gemini 1M, Claude 200K, GPT-4 128K

---

## Phase 3: Agent Mode v2

> **Goal:** More powerful autonomous task execution

| Priority | Task | Effort | Depends On | Status |
|----------|------|--------|------------|--------|
| P1 | Task dependency graph (DAG) | L | Phase 1 | ⚪ Pending |
| P1 | Parallel task execution | L | DAG | ⚪ Pending |
| P2 | Rollback on failure (undo last N tasks) | L | - | ⚪ Pending |
| P2 | Plan templates (save/load workflows) | M | - | ⚪ Pending |
| P2 | Custom verification scripts | M | - | ⚪ Pending |
| P3 | Agent memory (learn from past sessions) | XL | - | ⚪ Pending |
| P3 | Natural language plan editing | L | - | ⚪ Pending |

### Success Criteria
- [ ] Can run 3 independent tasks in parallel
- [ ] Failed task triggers rollback prompt
- [ ] User can save plan as reusable template

### Technical Design
```
Plan DAG Example:
  task-1 (init) ──┬── task-2 (build) ──┬── task-4 (deploy)
                  └── task-3 (test) ───┘
```

---

## Phase 4: Distribution

> **Goal:** Easy installation on all platforms

| Priority | Task | Effort | Depends On | Status |
|----------|------|--------|------------|--------|
| P1 | GitHub releases with goreleaser | M | Phase 1 | ⚪ Pending |
| P1 | Homebrew formula (macOS/Linux) | S | goreleaser | ⚪ Pending |
| P2 | Chocolatey package (Windows) | M | goreleaser | ⚪ Pending |
| P2 | Docker image | S | goreleaser | ⚪ Pending |
| P2 | APT repository (Debian/Ubuntu) | M | goreleaser | ⚪ Pending |
| P3 | AUR package (Arch Linux) | S | goreleaser | ⚪ Pending |
| P3 | Auto-update mechanism | L | - | ⚪ Pending |

### Success Criteria
- [ ] `brew install terminalizcrazy` works
- [ ] `choco install terminalizcrazy` works
- [ ] `docker run terminalizcrazy` works
- [ ] GitHub release has binaries for all platforms

### goreleaser Config
```yaml
# Already partially configured in .goreleaser.yaml
builds:
  - goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
```

---

## Phase 5: Collaboration v2

> **Goal:** Production-ready collaboration features

| Priority | Task | Effort | Depends On | Status |
|----------|------|--------|------------|--------|
| P2 | Public signaling server | L | Phase 1 | ⚪ Pending |
| P2 | Permission levels (view/suggest/full) | M | - | ⚪ Pending |
| P2 | Room persistence (rejoin after disconnect) | M | - | ⚪ Pending |
| P3 | File sharing in session | L | - | ⚪ Pending |
| P3 | Increase max users (5 → 10) | S | - | ⚪ Pending |
| P3 | Session recording | L | - | ⚪ Pending |

### Success Criteria
- [ ] Users can reconnect to same room after network drop
- [ ] View-only users cannot execute commands
- [ ] Public server handles 100 concurrent rooms

---

## Phase 6: Plugin Ecosystem

> **Goal:** Third-party plugin support

| Priority | Task | Effort | Depends On | Status |
|----------|------|--------|------------|--------|
| P2 | Plugin SDK documentation | M | Phase 1 | ⚪ Pending |
| P2 | Example plugins (3-5 examples) | M | SDK docs | ⚪ Pending |
| P3 | Plugin manifest.yaml schema | S | - | ⚪ Pending |
| P3 | Plugin enable/disable UI | M | - | ⚪ Pending |
| P3 | Plugin marketplace (registry) | XL | - | ⚪ Pending |
| P3 | Lua plugin runtime | XL | - | ⚪ Pending |

### Success Criteria
- [ ] Developer can create plugin following docs
- [ ] Plugin can be installed from URL
- [ ] Plugin config editable without code changes

---

## Backlog (Prioritized)

### Must Have (P1-P2) - Next 3 Months
| Task | Effort | Rationale |
|------|--------|-----------|
| MCP (Model Context Protocol) support | L | Industry standard for AI tools |
| SSH session support | L | Remote server management |
| Custom keybinding configuration | M | Power user request |
| Export session to Markdown | S | Documentation use case |

### Should Have (P2-P3) - Next 6 Months
| Task | Effort | Rationale |
|------|--------|-----------|
| Vim keybinding mode | M | Developer preference |
| tmux integration | M | Workflow integration |
| Session recording/playback | L | Training/demo use case |
| VS Code extension | L | IDE integration |

### Could Have (P3) - Future
| Task | Effort | Rationale |
|------|--------|-----------|
| GPU acceleration for rendering | XL | Performance edge case |
| Mobile companion app | XL | View-only monitoring |
| Voice chat integration | XL | Nice-to-have |
| High contrast themes | S | Accessibility |

### Won't Have (Descoped)
| Task | Reason |
|------|--------|
| Browser-based version | Focus on native terminal |
| Windows CMD support | Too limited, use Windows Terminal |

---

## Sprint Planning Template

### Sprint N (2 weeks)

**Goal:** [One sentence goal]

**Capacity:** [X story points / Y hours]

| Task | Priority | Effort | Assignee |
|------|----------|--------|----------|
| Task 1 | P0 | S | - |
| Task 2 | P1 | M | - |
| Task 3 | P1 | M | - |

**Definition of Done:**
- [ ] Code complete
- [ ] Tests passing (80%+ coverage)
- [ ] Documentation updated
- [ ] PR reviewed and merged

---

## Milestones

| Milestone | Target | Key Deliverables |
|-----------|--------|------------------|
| v0.9.0 | Phase 1 complete | 80% coverage, CI green |
| v1.0.0 | Phase 2 + 4 complete | Full streaming, GitHub releases |
| v1.1.0 | Phase 3 complete | Agent v2 with parallel tasks |
| v1.2.0 | Phase 5 + 6 complete | Plugin ecosystem, collab v2 |

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-04 | Default to Ollama | Local-first, no API key friction |
| 2026-04 | SQLite for storage | Embedded, zero-config |
| 2026-04 | Bubble Tea for TUI | Best Go TUI framework |

---

## Notes

### Coding Conventions
- Follow Go standard style
- Use `golangci-lint` before commits
- Conventional commit messages
- TDD for new features

### Testing Strategy
- Unit tests per package (target: 80%)
- Integration tests for AI/storage
- E2E tests for critical user flows
- Race detection in CI (`go test -race`)

### Release Cadence
- Semantic versioning (MAJOR.MINOR.PATCH)
- Changelog generation (keep-a-changelog format)
- GitHub releases with release notes
- Pre-release tags for testing (v1.0.0-rc.1)
