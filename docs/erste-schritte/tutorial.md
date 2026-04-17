# Interactive Tutorial

> Step-by-step introduction to TerminalizCrazy (~12 minutes)

## Table of Contents

| Step | Topic | Duration |
| ---- | ----- | -------- |
| 1 | [Check Installation](#step-1-check-installation) | 2 min |
| 2 | [First Launch](#step-2-first-launch) | 1 min |
| 3 | [Ask Your First Question](#step-3-ask-your-first-question) | 2 min |
| 4 | [Execute a Command](#step-4-execute-a-command) | 1 min |
| 5 | [Using Multi-Pane](#step-5-using-multi-pane) | 2 min |
| 6 | [Try Agent Mode](#step-6-try-agent-mode) | 3 min |
| 7 | [Change Theme](#step-7-change-theme) | 1 min |

## Learning Goals

After this tutorial you will be able to:

- Understand and use AI command suggestions
- Execute commands safely
- Use multi-pane layouts effectively
- Use Agent mode for complex tasks
- Share and restore sessions
- Customize themes

---

## Step 1: Check Installation

Verify that TerminalizCrazy is correctly installed:

```bash
# Check Ollama (default provider)
ollama list

# Gemma4 should be listed
# If not:
ollama pull gemma4
ollama serve
```

For cloud providers (optional):

```bash
# Gemini
export GEMINI_API_KEY="AIzaSy..." && export AI_PROVIDER="gemini"
```

---

## Step 2: First Launch

Start TerminalizCrazy:

```bash
terminalizcrazy
```

```plaintext
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   ollama/gemma4   [abc123]          |
+-------------------------------------------------------------+
|                                                             |
|  Welcome to TerminalizCrazy!                                |
|  Ask a question in natural language.                        |
|                                                             |
+-------------------------------------------------------------+
| > _                                                         |
+-------------------------------------------------------------+
| Enter: Send | Ctrl+E: Execute | Ctrl+Y: Copy                |
+-------------------------------------------------------------+
```

You will see:

- Header with project name and version
- AI provider and status (ollama/gemma4 = connected)
- Session ID in square brackets
- Input field at the bottom

Helpful keyboard shortcuts are shown in the footer.

---

## Step 3: Ask Your First Question

Ask your first question. Type:

```plaintext
list files
```

Press **Enter**.

```plaintext
+-------------------------------------------------------------+
|  You: list files                                            |
|                                                             |
|  AI: Here is the command to list files:                     |
|                                                             |
|     ls -la                                                  |
|                                                             |
|     -l shows details (permissions, size, date)              |
|     -a shows hidden files too                               |
|                                                             |
|  Press Ctrl+E to execute                                    |
+-------------------------------------------------------------+
```

The AI responds with an appropriate command and explanation.

---

## Step 4: Execute a Command

Press **Ctrl+E** to execute the suggested command.

Since `ls` is a safe command (risk level: LOW), it executes immediately:

```plaintext
+-------------------------------------------------------------+
|  $ ls -la                                                   |
|  total 64                                                   |
|  drwxr-xr-x  8 user  staff   256 Apr 11 10:30 .             |
|  drwxr-xr-x  5 user  staff   160 Apr 10 09:15 ..            |
|  -rw-r--r--  1 user  staff  1234 Apr 11 10:30 main.go       |
|  -rw-r--r--  1 user  staff   567 Apr 11 10:25 go.mod        |
|                                                             |
|  [Exit code: 0]                                             |
+-------------------------------------------------------------+
```

The result appears in the chat with the exit code.

---

## Step 5: Using Multi-Pane

TerminalizCrazy supports multiple panes simultaneously.

### Split Pane

Press **Ctrl+\\** for a vertical split:

```plaintext
+-------------------------+-------------------------+
|  Chat                   |  New Pane               |
+-------------------------+-------------------------+
| You: list files         |                         |
|                         | > _                     |
| AI: ls -la              |                         |
+-------------------------+-------------------------+
```

Or **Ctrl+-** for a horizontal split.

### Navigate Between Panes

- **Alt+Arrow keys**: Navigate between panes
- **Ctrl+Z**: Zoom pane (fullscreen)
- **Ctrl+W**: Close current pane

### Layout Presets

Start with a predefined layout:

```bash
terminalizcrazy --layout quad     # 4 panes
terminalizcrazy --layout tall     # 2 panes vertical
terminalizcrazy --layout wide     # 2 panes horizontal
```

---

## Step 6: Try Agent Mode

Agent Mode plans and executes complex tasks automatically.

### Activate Agent Mode

Press **Ctrl+A** to cycle through modes:

- `off` -> `suggest` -> `auto` -> `off`

Recommended: `suggest` (plans, asks before execution)

### Ask a Complex Task

```plaintext
> Set up a new React project with TypeScript and ESLint
```

The agent creates a plan:

```plaintext
+-------------------------------------------------------------+
| Agent Mode: Plan created                                    |
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
+-------------------------------------------------------------+
| [A]pprove  [R]eject  [M]odify  [S]kip                       |
+-------------------------------------------------------------+
```

- **A**: Approve and execute plan
- **R**: Reject plan
- **M**: Modify individual tasks
- **S**: Skip task

---

## Step 7: Change Theme

TerminalizCrazy offers 5 built-in themes.

### Set Theme in config.toml

```toml
# ~/.terminalizcrazy/config.toml
theme = "dracula"   # default, dracula, monokai, nord, solarized
```

### Available Themes

| Theme | Description |
| ----- | ----------- |
| `default` | Default colors |
| `dracula` | Dark theme with purple accents |
| `monokai` | Classic developer theme |
| `nord` | Nordic color palette |
| `solarized` | Eye-friendly colors |

### Hot Reload

Themes are automatically reloaded when the file changes:

```toml
[appearance]
theme_hot_reload = true
```

---

## Summary

You have learned:

### Basics

| Key | Action |
| --- | ------ |
| Enter | Send message |
| Ctrl+E | Execute command |
| Ctrl+Y | Copy command |
| Up/Down arrows | Browse history |
| Esc | Quit |

### Panes & Layouts

| Key | Action |
| --- | ------ |
| Ctrl+\\ | Vertical split |
| Ctrl+- | Horizontal split |
| Alt+Arrow keys | Switch between panes |
| Ctrl+Z | Zoom pane |
| Ctrl+W | Close pane |

### Agent & Collaboration

| Key | Action |
| --- | ------ |
| Ctrl+A | Toggle agent mode |
| Ctrl+M | Open model selector |
| Ctrl+S | Share session |
| Ctrl+J | Join session |
| Ctrl+D | End collaboration |

---

## Next Steps

Now you're ready for advanced features:

- [Agent Mode](../anleitungen/agent-modus.md) - Automate complex tasks
- [Workflows](../anleitungen/workflows.md) - Save recurring tasks
- [Collaboration](../anleitungen/zusammenarbeit.md) - Team features
- [Settings](../referenz/einstellungen.md) - All configuration options

---

## Tips for Daily Use

### Ask Precise Questions

```plaintext
Bad:  "do something with files"
Good: "show all Python files larger than 1MB"
```

### Use Context

TerminalizCrazy automatically detects your project:

- In a Git repo: Git-specific suggestions
- In a Node.js project: npm commands
- In a Python project: pip/python commands
- In a Go project: go commands

### When Uncertain: Ask

```plaintext
> What does the command "tar -xvzf archive.tar.gz" do?
```

The AI explains each part of the command.

### Switch Models

Press **Ctrl+M** to switch between AI models:

- Gemma4 (default, local via Ollama)
- Gemma4:e4b (compact variant)
- Gemini Flash/Pro (cloud)
- Claude (Anthropic, cloud)
- GPT-4 (OpenAI, cloud)

---

*Good luck with TerminalizCrazy!*
