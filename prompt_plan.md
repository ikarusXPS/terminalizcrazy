# TerminalizCrazy - Implementation Plan

> Roadmap for future development phases.

## Current Status

**Version:** Development
**Core Features:** Complete
**Test Coverage:** ~70%

### Implemented
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

## Phase 1: Testing & Stability

Priority: HIGH | Status: In Progress

- [ ] Increase test coverage to 80%+
- [ ] Fix flaky Windows workspace tests
- [ ] Add integration tests for AI providers
- [ ] Add E2E tests for collaboration flow
- [ ] Performance profiling and optimization

### Dependencies
- None (independent phase)

---

## Phase 2: AI Enhancements

Priority: MEDIUM | Status: Planned

- [ ] Add streaming support to Anthropic client
- [ ] Add streaming support to OpenAI client
- [ ] Implement context window management
- [ ] Add conversation summarization
- [ ] Support for custom system prompts
- [ ] Model-specific prompt optimization

### Dependencies
- Phase 1 (stability first)

---

## Phase 3: Agent Mode v2

Priority: MEDIUM | Status: Planned

- [ ] Parallel task execution
- [ ] Task dependency graph
- [ ] Rollback on failure
- [ ] Custom verification scripts
- [ ] Plan templates / saved workflows
- [ ] Agent memory across sessions

### Dependencies
- Phase 2 (AI enhancements)

---

## Phase 4: Collaboration v2

Priority: LOW | Status: Planned

- [ ] Public signaling server
- [ ] Room persistence
- [ ] File sharing
- [ ] Screen sharing mode
- [ ] Voice chat integration (optional)
- [ ] Permission levels (view-only, suggest, full)

### Dependencies
- Phase 1 (stability)

---

## Phase 5: Plugin Ecosystem

Priority: LOW | Status: Planned

- [ ] Plugin marketplace / registry
- [ ] Plugin SDK with examples
- [ ] Lua/JavaScript plugin support
- [ ] Plugin configuration UI
- [ ] Plugin auto-update

### Dependencies
- Phase 1 (stability)

---

## Phase 6: Distribution

Priority: MEDIUM | Status: Planned

- [ ] Homebrew formula
- [ ] APT/YUM packages
- [ ] Chocolatey package
- [ ] Docker image
- [ ] GitHub releases with goreleaser
- [ ] Auto-update mechanism

### Dependencies
- Phase 1 (stability)

---

## Backlog (Unprioritized)

- [ ] GPU acceleration for rendering
- [ ] SSH session support
- [ ] tmux integration
- [ ] Vim keybinding mode
- [ ] Custom keybinding configuration
- [ ] Session recording/playback
- [ ] Export session to HTML/Markdown
- [ ] Mobile companion app (view-only)
- [ ] VS Code extension integration
- [ ] MCP (Model Context Protocol) support

---

## Notes

### Coding Conventions
- Follow Go standard style
- Use `golangci-lint` before commits
- Conventional commit messages
- TDD for new features

### Testing Strategy
- Unit tests per package
- Integration tests for AI/storage
- E2E tests for user flows
- Race detection in CI

### Release Cadence
- Semantic versioning
- Changelog generation
- GitHub releases
