# Installation

> System requirements and installation guide

## System Requirements

### Minimum

| Component | Requirement |
|-----------|-------------|
| Operating System | Windows 10+, macOS 10.15+, Linux |
| Terminal | True Color support recommended |
| Storage | 50 MB disk space |
| RAM | 256 MB |

### Recommended

| Component | Requirement |
|-----------|-------------|
| Terminal | Windows Terminal, iTerm2, Kitty, Alacritty |
| Font | Nerd Font (for icons) |

---

## Installation

### Option 1: Pre-compiled Binary (recommended)

#### Windows

1. Download the latest version from GitHub Releases:
   terminalizcrazy-windows-amd64.exe

2. Move the file to your PATH:

   mkdir C:\Program Files\TerminalizCrazy
   move terminalizcrazy-windows-amd64.exe "C:\Program Files\TerminalizCrazy\terminalizcrazy.exe"

3. Add the folder to PATH (System > Environment Variables)

#### macOS

1. Download the latest version:
   curl -LO https://github.com/terminalizcrazy/terminalizcrazy/releases/latest/download/terminalizcrazy-darwin-arm64

2. Make the file executable and move it:
   chmod +x terminalizcrazy-darwin-arm64
   sudo mv terminalizcrazy-darwin-arm64 /usr/local/bin/terminalizcrazy

#### Linux

1. Download the latest version:
   curl -LO https://github.com/terminalizcrazy/terminalizcrazy/releases/latest/download/terminalizcrazy-linux-amd64

2. Make the file executable and move it:
   chmod +x terminalizcrazy-linux-amd64
   sudo mv terminalizcrazy-linux-amd64 /usr/local/bin/terminalizcrazy

---

### Option 2: Compile from Source

Prerequisites:
- Go 1.21 or newer

1. Clone the repository:
   git clone https://github.com/terminalizcrazy/terminalizcrazy
   cd terminalizcrazy

2. Compile:
   go build -o terminalizcrazy ./cmd/terminalizcrazy

3. Optional: Move to PATH:
   sudo mv terminalizcrazy /usr/local/bin/

---

## Configuration

### Set Up AI Provider

TerminalizCrazy uses Ollama with Gemma4 by default (local, free).

#### Ollama with Gemma4 - Default

1. Install Ollama from ollama.ai
2. Download Gemma4:

   ollama pull gemma4

3. Start Ollama:

   ollama serve

   No further configuration required - Ollama is the default provider.

#### Google Gemini (Cloud Alternative)

1. Create an account at aistudio.google.com
2. Generate an API key
3. Set the environment variables:

   export GEMINI_API_KEY="AIzaSyxxxxx"
   export AI_PROVIDER="gemini"

#### Anthropic Claude (Cloud Alternative)

1. Create an account at console.anthropic.com
2. Generate an API key
3. Set the environment variables:

   export ANTHROPIC_API_KEY="sk-ant-api03-xxxxx"
   export AI_PROVIDER="anthropic"

#### OpenAI (Cloud Alternative)

1. Create an account at platform.openai.com
2. Generate an API key
3. Set the environment variables:

   export OPENAI_API_KEY="sk-xxxxx"
   export AI_PROVIDER="openai"

---

## Configuration File

Optionally create a configuration file:

mkdir -p ~/.terminalizcrazy
cp config.toml.example ~/.terminalizcrazy/config.toml

Edit ~/.terminalizcrazy/config.toml as needed.

See [Settings](../referenz/einstellungen.md) for all options.

---

## First Start

1. Start TerminalizCrazy:

   terminalizcrazy

2. On successful start you will see:
   - Green status indicator for AI connection
   - Input prompt

3. Test with a simple request:

   > How do I list all files?

---

## Troubleshooting

### AI Not Connected

- Check if GEMINI_API_KEY, ANTHROPIC_API_KEY, or OPENAI_API_KEY is set
- Check the network connection
- For Ollama: Make sure ollama serve is running

### Terminal Shows Strange Characters

- Update to a terminal with True Color support
- Install a Nerd Font

### Command Not Found

- Make sure the binary is in PATH
- On Windows: Restart the terminal after PATH changes

---

## Next Steps

- [Quick Start](schnellstart.md) - Get started in 5 minutes
- [Tutorial](tutorial.md) - Interactive introduction tutorial
- [Settings](../referenz/einstellungen.md) - Configuration options
