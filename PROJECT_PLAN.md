# TerminalizCrazy - Project Plan

## Projektziel

TerminalizCrazy ist ein modernes, AI-natives Terminal-Tool, das sich durch eingebaute Collaboration-Features, intelligente Session-Verwaltung und ein Plugin-System von bestehenden Terminals abhebt. Das Projekt startet als CLI-Tool und wird später zu einer vollständigen Web-App erweitert. Zielgruppe sind Entwickler, Tech-affine Nutzer und Teams, die produktiver im Terminal arbeiten wollen.

---

## Marktanalyse: Terminal-Landschaft 2025/2026

### Die beliebtesten Terminals

| Terminal | Sprache | Stärken | Schwächen |
|----------|---------|---------|-----------|
| **Ghostty** | Zig | Schnellster (2ms Latenz), GPU-Rendering, GLSL Shaders, 45k GitHub Stars | Neu, SSH-Kompatibilitätsprobleme |
| **WezTerm** | Rust | Feature-komplett, Lua-Config, Built-in Multiplexer | 300-400 MB RAM |
| **Alacritty** | Rust | Minimal, 30 MB RAM, schnellster Renderer | Keine Tabs/Splits, keine Ligatures |
| **Warp** | Rust | AI-Integration, Block-UI, IDE-ähnlich | Closed Source, 8ms Latenz |
| **Kitty** | C/Python | Image Protocol, Linux Power-User | Komplexe Config |
| **iTerm2** | Obj-C | macOS Standard, viele Features | Langsam, hoher CPU |

### Beliebteste Features (was Entwickler lieben)

1. **GPU-Beschleunigung** - Schnelles Rendering auch bei großen Outputs
2. **AI-Autocomplete** - Warp's Killer-Feature: Befehle vorschlagen, Fehler erklären
3. **Block-basierte Outputs** - Outputs als navigierbare Blöcke statt Textstream
4. **Smart Completions** - Kontextbewusste Tab-Completion mit Typo-Korrektur
5. **Built-in Multiplexer** - Tabs/Splits ohne tmux/zellij
6. **Programmierbare Config** - Lua (WezTerm) statt statischer YAML/TOML
7. **Session Persistence** - Sessions überleben Neustarts
8. **Cross-Platform** - Gleiche Erfahrung auf Win/Mac/Linux

### Pain Points (was fehlt / nervt)

| Problem | Betroffene Terminals | Opportunity |
|---------|---------------------|-------------|
| Keine Echtzeit-Collaboration | Alle außer sshx/Termius | **Hoch** |
| Komplexe Konfiguration | tmux, Kitty | **Mittel** |
| SSH-Kompatibilität | Ghostty (xterm-ghostty) | **Mittel** |
| Keine variable Textgröße | Alle | **Niedrig** |
| Secret-Leaks im Output | Alle | **Hoch** |
| Keine Workflow-Automation | Alle | **Hoch** |
| Kein Plugin-System (oder komplex) | Alacritty, Warp | **Mittel** |

---

## Differenzierungsstrategie: Was macht TerminalizCrazy einzigartig?

### Kern-Differenzierungen (MVP)

#### 1. **AI-Native Architecture**
Nicht nur Autocomplete, sondern volle Agent-Integration:
- Natürliche Sprache → Befehl
- Fehler automatisch erklären + Fix vorschlagen
- Kontext-bewusste Vorschläge basierend auf aktuellem Verzeichnis/Projekt
- Integration mit Claude Code, Aider, Gemini CLI

#### 2. **Real-Time Collaboration**
Session-Sharing wie Google Docs:
- Link teilen → andere sehen Terminal live
- Cursor der anderen Nutzer sichtbar
- "Driver/Navigator" Modus für Pair Programming
- End-to-End verschlüsselt (WebRTC + AES)

#### 3. **Smart Sessions**
Kontextbewusste Projekt-Sessions:
- Automatische Projekt-Erkennung (package.json, Cargo.toml, etc.)
- Session-spezifische Umgebungsvariablen
- Workflow-Templates pro Projekttyp
- Session-History mit Suche

#### 4. **Secret Guard**
Automatischer Schutz vor Secret-Leaks:
- Erkennung von API-Keys, Tokens, Passwörtern im Output
- Automatisches Maskieren vor Screen-Share/Logging
- Warnung bei versehentlichem Commit von Secrets

### Spätere Features (Post-MVP)

- **Plugin-System** (WASM-basiert wie Zellij)
- **Visual Debugging** - Structured Output Parsing
- **Workflow-Automation** - Wiederkehrende Tasks als Makros
- **Web-Version** - Gleiche Features im Browser
- **Team-Features** - Shared Snippets, Command Library

---

## Tech-Stack Empfehlung

### Option A: Go + Charm.sh Ecosystem (Empfohlen für Start)

```
CLI Framework:    Bubble Tea (TUI) + Cobra (CLI)
UI Components:    Lip Gloss (Styling), Bubbles (Components)
AI Integration:   go-anthropic, go-openai
Collaboration:    WebRTC (pion/webrtc)
Config:           Viper + TOML
Testing:          go test + testify
```

**Warum Go?**
- Charm.sh Ecosystem ist das beste CLI/TUI Framework aktuell
- Einfache Cross-Compilation für Win/Mac/Linux
- Gute WebSocket/WebRTC Libraries für Collaboration
- Einfachere Lernkurve als Rust
- Später einfache Migration zu Web via WebAssembly (TinyGo)

### Option B: Rust + Ratatui

```
CLI Framework:    Ratatui (TUI) + Clap (CLI)
Terminal:         Crossterm
AI Integration:   async-openai, anthropic-rs
Collaboration:    webrtc-rs
Config:           config-rs + TOML
Testing:          cargo test + mockall
```

**Warum Rust?**
- Performance wie Alacritty/WezTerm
- Bewährtes Terminal-Ecosystem
- Memory Safety garantiert
- Steilere Lernkurve, aber langfristig robuster

### Option C: TypeScript + Ink (für schnellen Web-Übergang)

```
CLI Framework:    Ink (React für CLI)
UI Components:    ink-* packages
AI Integration:   @anthropic-ai/sdk, openai
Collaboration:    Socket.io
Config:           cosmiconfig
Testing:          Vitest
```

**Warum TypeScript?**
- Schnellste Entwicklung
- Gleicher Code für CLI und Web (mit Anpassungen)
- Größtes Ecosystem
- Weniger performant als Go/Rust

### Meine Empfehlung: **Go + Charm.sh**

Für ein CLI-first Projekt mit späterer Web-Erweiterung bietet Go den besten Kompromiss:
- Schnelle Entwicklung (vs. Rust)
- Echte Binary-Distribution (vs. TypeScript/Node)
- Hervorragendes TUI-Ecosystem (Charm.sh)
- Gute Performance (vs. TypeScript)

---

## Deployment-Strategie

### CLI Distribution

| Methode | Plattform | Priorität |
|---------|-----------|-----------|
| **GitHub Releases** | Alle | MVP |
| **Homebrew** | macOS/Linux | MVP |
| **Scoop/Winget** | Windows | MVP |
| **Go Install** | Alle (mit Go) | MVP |
| **Docker** | Alle | Post-MVP |
| **Snap/Flatpak** | Linux | Post-MVP |

### Web-Version (später)

| Option | Beschreibung |
|--------|--------------|
| **Vercel** | Ideal für Next.js Frontend |
| **Cloudflare Workers** | Edge-Funktionen für Collaboration |
| **Fly.io** | WebSocket-Server für Echtzeit |

### CI/CD Pipeline

```
GitHub Actions:
  - Build: Cross-compile für Win/Mac/Linux/ARM
  - Test: Unit + Integration Tests
  - Release: Automatische GitHub Releases + Homebrew Update
  - Security: Dependabot + Secret Scanning
```

---

## Feature-Liste

### MVP (Phase 1) ✅ COMPLETE

- [x] Basic CLI mit TUI Interface (Bubble Tea)
- [x] AI-Integration (Claude/OpenAI)
  - [x] Natürliche Sprache → Befehl
  - [x] Fehler-Erklärung
  - [x] Command-Autocomplete
- [x] Session Management
  - [x] Session erstellen/benennen/wechseln
  - [x] Session-Persistenz (SQLite)
- [x] Secret Guard (Basic)
  - [x] API-Key Pattern Detection (7 patterns)
  - [x] Output-Maskierung
- [x] Cross-Platform Build (Win/Mac/Linux)

### Phase 2 ✅ COMPLETE

- [x] Real-Time Collaboration
  - [x] Session-Sharing via Link (WebSocket)
  - [x] Live Cursor
  - [x] Chat
  - [x] E2E Encryption (ECDH + AES-256-GCM)
- [x] Smart Sessions
  - [x] Projekt-Erkennung (11 types)
  - [x] Workflow-Templates (6 built-in)
- [x] Erweiterte AI-Features
  - [x] Kontext-bewusste Vorschläge
  - [x] Multi-Model Support (Claude + OpenAI + Ollama)

### Phase 3 ✅ PARTIALLY COMPLETE

- [x] Plugin-System (Hook-based, 8 event types)
- [ ] Plugin-System (WASM) - Future upgrade
- [ ] Web-Version
- [ ] Team-Features
- [ ] Enterprise Features (SSO, Audit Logs)

### Phase 4 ✅ COMPLETE (Added)

- [x] Agent Mode (3 modes: off/suggest/auto)
- [x] Tabs/Splits (PaneManager)
- [x] Floating Panes
- [x] Theme System (5 themes, YAML, hot-reload)
- [x] Workspace Management (5 layouts)
- [x] Local AI (Ollama integration)

---

## Technische Architektur

```
┌─────────────────────────────────────────────────────────────┐
│                     TerminalizCrazy                         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   CLI/TUI   │  │  AI Engine  │  │  Collaboration      │ │
│  │  (Bubble    │  │  (Claude/   │  │  (WebRTC/           │ │
│  │   Tea)      │  │   OpenAI)   │  │   WebSocket)        │ │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
│         │                │                     │            │
│  ┌──────┴────────────────┴─────────────────────┴──────────┐│
│  │                    Core Engine                          ││
│  │  ┌────────────┐ ┌────────────┐ ┌────────────────────┐  ││
│  │  │  Session   │ │  Secret    │ │  Plugin System     │  ││
│  │  │  Manager   │ │  Guard     │ │  (WASM Runtime)    │  ││
│  │  └────────────┘ └────────────┘ └────────────────────┘  ││
│  └─────────────────────────────────────────────────────────┘│
│                              │                              │
│  ┌───────────────────────────┴────────────────────────────┐│
│  │                    Storage Layer                        ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌────────────────┐  ││
│  │  │   Config    │  │  Sessions   │  │   History      │  ││
│  │  │   (TOML)    │  │  (SQLite)   │  │   (SQLite)     │  ││
│  │  └─────────────┘  └─────────────┘  └────────────────┘  ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    External Services                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │  Claude API │  │  OpenAI API │  │  Signaling Server   │ │
│  │             │  │             │  │  (Collaboration)    │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## Status: Feature-Complete (Phases 1-4)

**Implementiert** (~23,670 LoC):
- Go + Charm.sh (Bubble Tea) Tech-Stack ✅
- AI: Claude + OpenAI + Ollama ✅
- Agent Mode mit Plan-Verifikation ✅
- Real-Time Collaboration mit E2E Encryption ✅
- Tabs/Splits/Floating Panes ✅
- Plugin-System (Hook-based) ✅
- Theme-System + Workspace Management ✅

## Nächste Schritte

1. **GPU Acceleration** - Performance-Verbesserung (Zig/Rust Rendering Layer)
2. **Web-Version** - Browser-basierter Zugriff
3. **WASM Plugins** - Upgrade von Hooks zu WASM
4. **Community Launch** - GitHub Release, Homebrew, Scoop

---

## Quellen

### Terminal-Vergleiche
- [Best Terminal Emulators 2026 - Scopir](https://scopir.com/posts/best-terminal-emulators-developers-2026/)
- [Modern Terminals: Ghostty, WezTerm, Alacritty - Calmops](https://calmops.com/tools/modern-terminal-emulators-2026-ghostty-wezterm-alacritty/)
- [Terminal Showdown - CodeMiner42](https://blog.codeminer42.com/modern-terminals-alacritty-kitty-and-ghostty/)

### AI CLI Tools
- [12 CLI Tools Redefining Developer Workflows - Qodo](https://www.qodo.ai/blog/best-cli-tools/)
- [Top 5 Agentic CLI Tools - KDnuggets](https://www.kdnuggets.com/top-5-agentic-coding-cli-tools)
- [AI Terminal Coding Tools - Augment Code](https://www.augmentcode.com/guides/ai-terminal-coding-tools-that-actually-work-in-2025)

### Warp Features
- [Why Developers Should Try Warp - DEV Community](https://dev.to/trantn/why-developers-should-try-warp-the-terminal-that-boosts-your-productivity-8i9)
- [Warp All Features](https://www.warp.dev/all-features)

### Terminal Multiplexer
- [Tmux vs Zellij - TmuxAI](https://tmuxai.dev/tmux-vs-zellij/)
- [Zellij vs Tmux Complete Comparison - Medium](https://rrmartins.medium.com/zellij-vs-tmux-complete-comparison-or-almost-8e5b57d234ae)

### Collaboration
- [sshx - Collaborative Terminal Sharing](https://sshx.io/)
- [Termius Multiplayer - Termius Blog](https://termius.com/blog/from-chaos-to-clarity-reimagining-real-time-collaboration-in-the-terminal)
- [CoScreen Terminal - Datadog](https://www.datadoghq.com/blog/datadog-coscreen-collaborative-terminal-pair-programming/)

### Tech-Stack
- [Ghostty Terminal Built with Zig - Calmops](https://calmops.com/programming/ghostty-terminal-zig/)
- [CLI Apps in Rust - Rust CLI Book](https://rust-cli.github.io/book/index.html)
