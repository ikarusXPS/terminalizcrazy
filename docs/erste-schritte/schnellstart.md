# Quick Start

> Get started in 5 minutes

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   ollama/gemma4   [session-id]      |
+-------------------------------------------------------------+
|                                                             |
|  Welcome! Ask your first question...                        |
|                                                             |
+-------------------------------------------------------------+
| > _                                                         |
+-------------------------------------------------------------+
| Enter: Send | Ctrl+E: Execute | Esc: Quit                   |
+-------------------------------------------------------------+
```

## Prerequisites

- TerminalizCrazy installed (see [Installation](installation.md))
- Ollama with Gemma4 (default) or cloud provider configured

---

## 1. Start

Open a terminal and start TerminalizCrazy:

    terminalizcrazy

You should see:
- Header with version and AI status (green = connected)
- Input field at the bottom

---

## 2. First Request

Enter a question in natural language:

    > How do I find the largest files in this folder?

```
+-------------------------------------------------------------+
|  You: How do I find the largest files?                      |
|                                                             |
|  AI: Here is the appropriate command:                       |
|                                                             |
|     find . -type f -exec du -h {} + | sort -rh | head -20   |
|                                                             |
|     This command finds all files, sorts them by size,       |
|     and shows the 20 largest.                               |
|                                                             |
|  Press Ctrl+E to execute                                    |
+-------------------------------------------------------------+
```

TerminalizCrazy responds with:
- The appropriate terminal command
- A brief explanation
- Hint: "Press Ctrl+E to execute"

---

## 3. Execute Command

Press **Ctrl+E** to execute the suggested command.

For risky commands, a confirmation appears:

```
+-------------------------------------------------------------+
|  MEDIUM: This command will modify files                     |
|                                                             |
|  find . -type f -exec du -h {} + | sort -rh | head -20      |
|                                                             |
|  Execute? [Y]es / [N]o                                      |
+-------------------------------------------------------------+
```

- Press **Y** for Yes
- Press **N** for No

The result is displayed in the chat.

---

## 4. Copy Command

Press **Ctrl+Y** to copy the last command to the clipboard.

Useful when you want to modify the command or use it in another terminal.

---

## 5. Browse History

Press the **Up arrow** to browse through previous inputs.

History is preserved across sessions.

---

## Important Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Enter | Send message |
| Ctrl+E | Execute last command |
| Ctrl+Y | Copy command |
| Ctrl+L | Clear chat |
| Up/Down arrows | Browse history |
| Esc | Quit |

Full list: [Keyboard Shortcuts](../referenz/tastenkuerzel.md)

---

## Examples

### File Management

    > Show all files larger than 100MB
    > Delete all .tmp files
    > Find duplicate files

### Git

    > Show the last 5 commits
    > Create a new branch called feature/login
    > What did I change today?

### System Information

    > How much disk space is left?
    > Which processes are using the most RAM?
    > Show my IP address

### Development

    > Install all npm dependencies
    > Run the tests
    > Start the development server

---

## Sessions

TerminalizCrazy saves your conversations automatically.

On the next start you can:
1. Continue an existing session
2. Start a new session

Sessions are organized by project (based on your working directory).

---

## Collaboration

Share your session with others:

```
+-------------------------------------------------------------+
|  Share Session                                              |
|                                                             |
|  Share Code: ABCD-1234                                      |
|                                                             |
|  Share this code with other participants.                   |
|  They can join with Ctrl+J.                                 |
|                                                             |
|  [Connected: 2 participants]                                |
+-------------------------------------------------------------+
```

1. Press **Ctrl+S** to share
2. A share code is displayed (e.g., ABCD-1234)
3. Others press **Ctrl+J** and enter the code

Both see all messages and commands in real-time (E2E encrypted).

---

## Next Steps

- [Tutorial](tutorial.md) - Detailed interactive tutorial
- [Agent Mode](../anleitungen/agent-modus.md) - Automate complex tasks
- [Settings](../referenz/einstellungen.md) - Customization options
