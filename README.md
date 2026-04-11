# TerminalizCrazy

**AI-native terminal with real-time collaboration, Agent Mode, and SecretGuard**

[![CI](https://github.com/ikarusXPS/terminalizcrazy/actions/workflows/ci.yml/badge.svg)](https://github.com/ikarusXPS/terminalizcrazy/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ikarusXPS/terminalizcrazy)](https://goreportcard.com/report/github.com/ikarusXPS/terminalizcrazy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ikarusXPS/terminalizcrazy)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/ikarusXPS/terminalizcrazy)](https://github.com/ikarusXPS/terminalizcrazy/releases)
[![Downloads](https://img.shields.io/github/downloads/ikarusXPS/terminalizcrazy/total)](https://github.com/ikarusXPS/terminalizcrazy/releases)

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   gemini   [abc123]                 |
+-------------------------------------------------------------+
|                                                             |
|  You: how to find large files?                              |
|                                                             |
|  AI: To find large files, use:                              |
|     find . -size +100M -type f                              |
|     This searches recursively for files over 100MB.         |
|                                                             |
|  $ find . -size +100M -type f                               |
|    ./videos/movie.mp4 (1.2GB)                               |
|    ./backup/archive.zip (500MB)                             |
|                                                             |
+-------------------------------------------------------------+
| > Ask anything...                                           |
+-------------------------------------------------------------+
| Enter: Send | Ctrl+E: Execute | Ctrl+Y: Copy | Esc: Quit    |
+-------------------------------------------------------------+
```

TerminalizCrazy is a modern, feature-complete AI terminal built with Go and the Charm.sh ecosystem. It combines powerful AI capabilities with real-time collaboration, making terminal work more productive and collaborative.

## Features

### AI Integration
- **Multi-Provider Support** - Gemini (default), Claude (Anthropic), GPT (OpenAI), and Ollama (local)
- **Natural Language Commands** - Describe what you want, get the right command
- **Error Explanation** - AI explains errors and suggests fixes
- **Project-Aware Context** - Recognizes 11 project types with framework detection

### Agent Mode
- **Autonomous Task Execution** - Multi-step plans with verification
- **Three Modes** - `off`, `suggest` (recommended), `auto`
- **Risk-Aware** - Confirms before risky operations
- **Plan Persistence** - Plans saved to SQLite for review

### Real-Time Collaboration
- **Session Sharing** - Share terminal sessions like Google Docs
- **E2E Encryption** - ECDH key exchange + AES-256-GCM
- **User Presence** - See who's connected with typing indicators
- **Share Codes** - Simple `xxxx-yyyy` format for joining

### Tabs, Splits & Workspaces

```
+-------------------------+-------------------------+
|  Chat                   |  Output                 |
+-------------------------+-------------------------+
| You: git status         | $ git status            |
|                         | On branch main          |
| AI: Running...          | Changes not staged:     |
|                         |   modified: app.go      |
+-------------------------+-------------------------+
| Tab: Next | Ctrl+W: Close | Ctrl+\: Split         |
+-------------------------+-------------------------+
```

- **Multi-Pane Layout** - Horizontal and vertical splits
- **Floating Panes** - Zellij-style floating windows
- **5 Layout Presets** - quad, tall, wide, stack, single
- **Workspace Persistence** - Layouts saved and restored

### Security
- **SecretGuard** - Auto-masks 7 secret types (AWS, GitHub, JWT, etc.)
- **Risk Assessment** - 4 levels (Low/Medium/High/Critical)
- **Confirmation Prompts** - Required for dangerous commands

### Customization
- **Theme System** - 5 built-in themes, custom YAML themes
- **Hot Reload** - Theme changes apply instantly
- **Plugin System** - 8 hook types, 4 built-in plugins
- **Workflow Templates** - 6 built-in, create your own

---

## Installation

### From Release (Recommended)

Download the latest release for your platform:

```bash
# macOS (Intel)
curl -L https://github.com/ikarusXPS/terminalizcrazy/releases/latest/download/terminalizcrazy_darwin_amd64.tar.gz | tar xz

# macOS (Apple Silicon)
curl -L https://github.com/ikarusXPS/terminalizcrazy/releases/latest/download/terminalizcrazy_darwin_arm64.tar.gz | tar xz

# Linux
curl -L https://github.com/ikarusXPS/terminalizcrazy/releases/latest/download/terminalizcrazy_linux_amd64.tar.gz | tar xz

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/ikarusXPS/terminalizcrazy/releases/latest/download/terminalizcrazy_windows_amd64.zip -OutFile terminalizcrazy.zip
Expand-Archive terminalizcrazy.zip
```

### Homebrew (macOS/Linux)

```bash
brew tap ikarusXPS/tap
brew install terminalizcrazy
```

### Scoop (Windows)

```bash
scoop bucket add terminalizcrazy https://github.com/ikarusXPS/scoop-bucket
scoop install terminalizcrazy
```

### From Source

```bash
git clone https://github.com/ikarusXPS/terminalizcrazy.git
cd terminalizcrazy
make build
./bin/terminalizcrazy
```

### Go Install

```bash
go install github.com/ikarusXPS/terminalizcrazy/cmd/terminalizcrazy@latest
```

---

## Quick Start

### 1. Set API Key

```bash
# Google Gemini (default)
export GEMINI_API_KEY="AIzaSy..."

# OR Anthropic (Claude)
export ANTHROPIC_API_KEY="sk-ant-api03-..."
export AI_PROVIDER="anthropic"

# OR OpenAI
export OPENAI_API_KEY="sk-..."
export AI_PROVIDER="openai"

# OR Ollama (local, no key needed)
export OLLAMA_ENABLED=true
export OLLAMA_MODEL=codellama
export AI_PROVIDER="ollama"
```

### 2. Run

```bash
terminalizcrazy
```

### 3. Ask Questions

```
> how to find large files over 100MB
```

AI suggests: `find . -size +100M -type f`

Press `Ctrl+E` to execute.

---

## Key Bindings

### Chat & Commands

| Key | Action |
|-----|--------|
| `Enter` | Send message to AI |
| `Ctrl+E` | Execute suggested command |
| `Ctrl+Y` | Copy command to clipboard |
| `Ctrl+L` | Clear chat |
| `Ctrl+R` | Show last command |
| `Up/Down` | Navigate history |
| `Esc` | Exit |

### Tabs & Panes

| Key | Action |
|-----|--------|
| `Ctrl+T` | New tab |
| `Ctrl+W` | Close pane |
| `Ctrl+\` | Vertical split |
| `Ctrl+-` | Horizontal split |
| `Alt+Arrow` | Navigate panes |
| `Ctrl+Z` | Toggle zoom |

### Collaboration

| Key | Action |
|-----|--------|
| `Ctrl+S` | Share session |
| `Ctrl+J` | Join session |
| `Ctrl+D` | Disconnect |

---

## Configuration

Config file: `~/.terminalizcrazy/config.toml`

```toml
# AI Provider
ai_provider = "gemini"     # gemini (default), anthropic, openai, ollama
gemini_model = "gemini-1.5-flash"  # gemini-1.5-flash, gemini-1.5-pro, gemini-2.0-flash-exp

# Agent Mode
agent_mode = "suggest"     # off, suggest, auto
agent_max_tasks = 10

# Ollama (Local AI)
ollama_enabled = false
ollama_url = "http://localhost:11434"
ollama_model = "codellama"

# Appearance
theme = "default"          # default, dracula, monokai, nord, solarized

[workspace]
default_layout = "quad"    # quad, tall, wide, stack, single
auto_save = true

[appearance]
theme_hot_reload = true
```

See [config.toml.example](config.toml.example) for all options.

---

## Agent Mode

Agent Mode enables autonomous multi-step task execution:

```
+-------------------------------------------------------------+
| AI Agent Mode: Plan created                                 |
+-------------------------------------------------------------+
|                                                             |
|  Plan: React TypeScript Setup (3 Tasks)                     |
|                                                             |
|  [1] ... npx create-react-app myapp --template typescript   |
|      Verification: myapp/ exists                            |
|                                                             |
|  [2] [ ] cd myapp && npm install eslint --save-dev          |
|      Verification: eslint in package.json                   |
|                                                             |
|  [3] [ ] npx eslint --init                                  |
|      Verification: .eslintrc.* exists                       |
|                                                             |
|  -----------------------------------------------------------+
|  [A]pprove  [R]eject  [M]odify Task  [S]kip Task            |
|                                                             |
+-------------------------------------------------------------+
```

### Modes

| Mode | Behavior |
|------|----------|
| `off` | No planning, single commands only |
| `suggest` | Creates plans, asks for approval (recommended) |
| `auto` | Executes LOW-risk commands automatically |

---

## Collaboration

Share your terminal session in real-time:

```
+-------------------------------------------------------------+
|  TerminalizCrazy   Sharing: a1b2-c3d4 (2 users)             |
+-------------------------------------------------------------+
|                                                             |
|  Alice: Let me check the logs                               |
|                                                             |
|  Bob: Sure, try journalctl -f                               |
|                                                             |
|  $ journalctl -f                                            |
|    [synced to all participants]                             |
|                                                             |
+-------------------------------------------------------------+
| E2E Encrypted | Ctrl+D: Disconnect                          |
+-------------------------------------------------------------+
```

**Host:**
1. Press `Ctrl+S`
2. Share the code: `a1b2-c3d4`

**Guest:**
1. Press `Ctrl+J`
2. Enter code: `a1b2-c3d4`

All messages, commands, and outputs are synchronized with E2E encryption.

---

## Project Detection

TerminalizCrazy automatically detects your project type:

| File | Type | Icon |
|------|------|------|
| `go.mod` | Go | 🐹 |
| `package.json` | Node.js | 📦 |
| `requirements.txt` | Python | 🐍 |
| `Cargo.toml` | Rust | 🦀 |
| `pom.xml` | Java | ☕ |
| `Dockerfile` | Docker | 🐳 |

The AI uses this context for better command suggestions.

---

## Documentation

- [Installation Guide](docs/erste-schritte/installation.md)
- [Quick Start](docs/erste-schritte/schnellstart.md)
- [Settings Reference](docs/referenz/einstellungen.md)
- [Keybindings](docs/referenz/tastenkuerzel.md)
- [Theme Customization](docs/referenz/themes.md)
- [Agent Mode Guide](docs/anleitungen/agent-modus.md)
- [Collaboration Guide](docs/anleitungen/zusammenarbeit.md)
- [Plugin Development](docs/anleitungen/plugins.md)

---

## Development

```bash
# Build
make build

# Test
make test

# Test with coverage
make test-coverage

# Lint
make lint

# Build all platforms
make build-all
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

---

## Comparison

| Feature | TerminalizCrazy | Warp | Wave |
|---------|-----------------|------|------|
| AI Commands | Yes | Yes | Yes |
| Agent Mode | Yes (3 modes) | Yes | No |
| Local AI (Ollama) | Yes | No | Yes |
| Real-Time Collab | Yes (E2E encrypted) | Paid | No |
| Tabs/Splits | Yes | Yes | No |
| Open Source | Yes | No | Yes |
| Windows Support | Yes | No | Yes |
| SecretGuard | Yes | No | No |

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Links

- [Documentation](docs/README.md)
- [Changelog](CHANGELOG.md)
- [Contributing](CONTRIBUTING.md)
- [Security Policy](SECURITY_CHECKLIST.md)

---

## Made With

Built with love using:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Charm.sh](https://charm.sh/) - CLI ecosystem
- [Go](https://go.dev/) - Programming language
