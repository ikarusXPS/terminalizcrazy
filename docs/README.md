# TerminalizCrazy Documentation

> AI-powered terminal for intelligent command execution

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   ollama/gemma4   [session]         |
+-------------------------------------------------------------+
|  AI-native terminal with Agent Mode, Collaboration,         |
|  SecretGuard and Multi-Provider Support                     |
+-------------------------------------------------------------+
```

## Quick Navigation

### Getting Started

| Document | Description |
|----------|-------------|
| [Installation](erste-schritte/installation.md) | System requirements and installation |
| [Quick Start](erste-schritte/schnellstart.md) | Get started in 5 minutes |
| [Tutorial](erste-schritte/tutorial.md) | Interactive tutorial (~12 minutes) |

### Guides

| Document | Description |
|----------|-------------|
| [AI Integration](anleitungen/ai-integration.md) | Use AI features optimally |
| [Agent Mode](anleitungen/agent-modus.md) | Automate complex tasks |
| [Collaboration](anleitungen/zusammenarbeit.md) | Share sessions with others |
| [Workflows](anleitungen/workflows.md) | Save recurring tasks |
| [Plugins](anleitungen/plugins.md) | Develop your own plugins |

### Reference

| Document | Description |
|----------|-------------|
| [Settings](referenz/einstellungen.md) | All 40+ configuration options |
| [Keyboard Shortcuts](referenz/tastenkuerzel.md) | Complete key bindings |
| [Themes](referenz/themes.md) | Theme customization and creation |
| [Risk Levels](referenz/risikostufen.md) | Understanding command safety |

### Concepts

| Document | Description |
|----------|-------------|
| [Architecture](konzepte/architektur.md) | System overview |
| [Secret Guard](konzepte/secret-guard.md) | Automatic secret masking |
| [Project Detection](konzepte/projekt-erkennung.md) | Intelligent project analysis |

---

## What is TerminalizCrazy?

TerminalizCrazy is an **AI-native terminal** that translates natural language into terminal commands. Instead of memorizing complex command syntax, simply describe what you want to do.

### Key Features

- **Multi-Provider Support** - Ollama/Gemma4 (default), Gemini, Claude, GPT
- **Agent Mode** - Automatic planning and execution of complex tasks (Ctrl+A)
- **Live Model Switching** - Change AI model without restart (Ctrl+M)
- **Risk Assessment** - Warning for dangerous commands (4 levels)
- **Secret Guard** - Automatic masking of 7 secret types
- **Multi-Pane Layouts** - Horizontal/vertical splits, 5 layout presets
- **Session Persistence** - All conversations are saved
- **Project Detection** - 11 project types with framework detection
- **Collaboration** - Share E2E-encrypted sessions in real-time

### Example

```
You: "show me the largest files in this folder"

TerminalizCrazy: Here is the command:
  du -ah . | sort -rh | head -20

  Press Ctrl+E to execute
```

---

## Quick Start

```bash
# 1. Start Ollama with Gemma4 (default, local)
ollama pull gemma4
ollama serve

# Or cloud provider:
# export GEMINI_API_KEY="AIzaSy..." && export AI_PROVIDER="gemini"

# 2. Start
terminalizcrazy

# 3. Ask questions
> how do I create a new git repository?
```

More details: [Quick Start Guide](erste-schritte/schnellstart.md)

---

## Configuration

TerminalizCrazy is configured via `~/.terminalizcrazy/config.toml`.

Full example configuration: [`config.toml.example`](../config.toml.example)

Detailed documentation: [Settings](referenz/einstellungen.md)

---

## Support

- **Issues**: [GitHub Issues](https://github.com/terminalizcrazy/terminalizcrazy/issues)
- **Discussions**: [GitHub Discussions](https://github.com/terminalizcrazy/terminalizcrazy/discussions)

---

*Created with TerminalizCrazy*
