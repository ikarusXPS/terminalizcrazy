# Installation

> System-Anforderungen und Installationsanleitung

## System-Anforderungen

### Minimum

| Komponente | Anforderung |
|------------|-------------|
| Betriebssystem | Windows 10+, macOS 10.15+, Linux |
| Terminal | True Color Unterstuetzung empfohlen |
| Speicher | 50 MB Festplatte |
| RAM | 256 MB |

### Empfohlen

| Komponente | Anforderung |
|------------|-------------|
| Terminal | Windows Terminal, iTerm2, Kitty, Alacritty |
| Schriftart | Nerd Font (fuer Icons) |

---

## Installation

### Option 1: Vorkompilierte Binary (empfohlen)

#### Windows

1. Laden Sie die neueste Version von GitHub Releases herunter:
   terminalizcrazy-windows-amd64.exe

2. Verschieben Sie die Datei in Ihren PATH:
   
   mkdir C:\Program Files\TerminalizCrazy
   move terminalizcrazy-windows-amd64.exe "C:\Program Files\TerminalizCrazy	erminalizcrazy.exe"

3. Fuegen Sie den Ordner zum PATH hinzu (System > Umgebungsvariablen)

#### macOS

1. Laden Sie die neueste Version herunter:
   curl -LO https://github.com/terminalizcrazy/terminalizcrazy/releases/latest/download/terminalizcrazy-darwin-arm64

2. Machen Sie die Datei ausfuehrbar und verschieben Sie sie:
   chmod +x terminalizcrazy-darwin-arm64
   sudo mv terminalizcrazy-darwin-arm64 /usr/local/bin/terminalizcrazy

#### Linux

1. Laden Sie die neueste Version herunter:
   curl -LO https://github.com/terminalizcrazy/terminalizcrazy/releases/latest/download/terminalizcrazy-linux-amd64

2. Machen Sie die Datei ausfuehrbar und verschieben Sie sie:
   chmod +x terminalizcrazy-linux-amd64
   sudo mv terminalizcrazy-linux-amd64 /usr/local/bin/terminalizcrazy

---

### Option 2: Aus Quellcode kompilieren

Voraussetzungen:
- Go 1.21 oder neuer

1. Repository klonen:
   git clone https://github.com/terminalizcrazy/terminalizcrazy
   cd terminalizcrazy

2. Kompilieren:
   go build -o terminalizcrazy ./cmd/terminalizcrazy

3. Optional: In PATH verschieben:
   sudo mv terminalizcrazy /usr/local/bin/

---

## Konfiguration

### API-Schluessel einrichten

TerminalizCrazy benoetigt einen KI-API-Schluessel. Waehlen Sie einen Anbieter:

#### Anthropic (Claude) - empfohlen

1. Erstellen Sie einen Account bei console.anthropic.com
2. Generieren Sie einen API-Schluessel
3. Setzen Sie die Umgebungsvariable:

   export ANTHROPIC_API_KEY="sk-ant-api03-xxxxx"

   Fuer dauerhafte Nutzung zu ~/.bashrc oder ~/.zshrc hinzufuegen.

#### OpenAI

1. Erstellen Sie einen Account bei platform.openai.com
2. Generieren Sie einen API-Schluessel
3. Setzen Sie die Umgebungsvariablen:

   export OPENAI_API_KEY="sk-xxxxx"
   export AI_PROVIDER="openai"

#### Ollama (lokal, kostenlos)

1. Installieren Sie Ollama von ollama.ai
2. Laden Sie ein Modell herunter:

   ollama pull codellama

3. Setzen Sie die Umgebungsvariable:

   export OLLAMA_ENABLED="true"
   export AI_PROVIDER="ollama"

---

## Konfigurationsdatei

Erstellen Sie optional eine Konfigurationsdatei:

mkdir -p ~/.terminalizcrazy
cp config.toml.example ~/.terminalizcrazy/config.toml

Bearbeiten Sie ~/.terminalizcrazy/config.toml nach Bedarf.

Siehe [Einstellungen](../referenz/einstellungen.md) fuer alle Optionen.

---

## Erster Start

1. Starten Sie TerminalizCrazy:

   terminalizcrazy

2. Bei erfolgreichem Start sehen Sie:
   - Gruene Statusanzeige fuer KI-Verbindung
   - Eingabeaufforderung

3. Testen Sie mit einer einfachen Anfrage:

   > Wie liste ich alle Dateien auf?

---

## Troubleshooting

### KI nicht verbunden

- Pruefen Sie, ob ANTHROPIC_API_KEY oder OPENAI_API_KEY gesetzt ist
- Pruefen Sie die Netzwerkverbindung
- Bei Ollama: Stellen Sie sicher, dass ollama serve laeuft

### Terminal zeigt komische Zeichen

- Aktualisieren Sie auf ein Terminal mit True Color Support
- Installieren Sie eine Nerd Font

### Befehl nicht gefunden

- Stellen Sie sicher, dass das Binary im PATH ist
- Unter Windows: Starten Sie das Terminal neu nach PATH-Aenderungen

---

## Naechste Schritte

- [Schnellstart](schnellstart.md) - In 5 Minuten loslegen
- [Tutorial](tutorial.md) - Interaktives Einfuehrungs-Tutorial
- [Einstellungen](../referenz/einstellungen.md) - Konfigurationsoptionen
