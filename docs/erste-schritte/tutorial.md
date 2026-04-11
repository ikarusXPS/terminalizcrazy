# Interaktives Tutorial

> Schritt-fuer-Schritt Einfuehrung in TerminalizCrazy (~12 Minuten)

## Inhaltsverzeichnis

| Schritt | Thema | Dauer |
|---------|-------|-------|
| 1 | [Installation pruefen](#schritt-1-installation-pruefen) | 2 min |
| 2 | [Erster Start](#schritt-2-erster-start) | 1 min |
| 3 | [Erste Frage stellen](#schritt-3-erste-frage-stellen) | 2 min |
| 4 | [Befehl ausfuehren](#schritt-4-befehl-ausfuehren) | 1 min |
| 5 | [Multi-Pane nutzen](#schritt-5-multi-pane-nutzen) | 2 min |
| 6 | [Agent Mode ausprobieren](#schritt-6-agent-mode-ausprobieren) | 3 min |
| 7 | [Theme wechseln](#schritt-7-theme-wechseln) | 1 min |

## Lernziele

Nach diesem Tutorial werden Sie:

- KI-Befehlsvorschlaege verstehen und nutzen
- Befehle sicher ausfuehren koennen
- Multi-Pane Layouts effektiv einsetzen
- Den Agent-Modus fuer komplexe Aufgaben nutzen
- Sessions teilen und wiederherstellen
- Themes anpassen

---

## Schritt 1: Installation pruefen

Pruefen Sie, ob TerminalizCrazy korrekt installiert ist:

```bash
# API-Key pruefen (Gemini ist Standard)
echo $GEMINI_API_KEY

# Oder einen der anderen Provider
echo $ANTHROPIC_API_KEY
echo $OPENAI_API_KEY
```

Falls noch kein API-Key gesetzt ist:

```bash
# Gemini (kostenlos, empfohlen fuer den Start)
export GEMINI_API_KEY="AIzaSy..."
```

---

## Schritt 2: Erster Start

Starten Sie TerminalizCrazy:

```bash
terminalizcrazy
```

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   gemini   [abc123]                 |
+-------------------------------------------------------------+
|                                                             |
|  Willkommen bei TerminalizCrazy!                            |
|  Stellen Sie eine Frage in natuerlicher Sprache.            |
|                                                             |
+-------------------------------------------------------------+
| > _                                                         |
+-------------------------------------------------------------+
| Enter: Senden | Ctrl+E: Ausfuehren | Ctrl+Y: Kopieren       |
+-------------------------------------------------------------+
```

Sie sehen:
- Den Header mit Projektname und Version
- KI-Provider und Status (gemini = verbunden)
- Session-ID in eckigen Klammern
- Eingabefeld am unteren Rand

Hilfreiche Tastenkuerzel werden in der Fusszeile angezeigt.

---

## Schritt 3: Erste Frage stellen

Stellen Sie Ihre erste Frage. Tippen Sie:

    list files

Druecken Sie **Enter**.

```
+-------------------------------------------------------------+
|  You: list files                                            |
|                                                             |
|  AI: Hier ist der Befehl zum Auflisten von Dateien:         |
|                                                             |
|     ls -la                                                  |
|                                                             |
|     -l zeigt Details (Rechte, Groesse, Datum)               |
|     -a zeigt auch versteckte Dateien                        |
|                                                             |
|  Press Ctrl+E to execute                                    |
+-------------------------------------------------------------+
```

Die KI antwortet mit einem passenden Befehl und Erklaerung.

---

## Schritt 4: Befehl ausfuehren

Druecken Sie **Ctrl+E** um den vorgeschlagenen Befehl auszufuehren.

Da `ls` ein sicherer Befehl ist (Risikostufe: LOW), wird er sofort ausgefuehrt:

```
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

Das Ergebnis erscheint im Chat mit Exit-Code.

---

## Schritt 5: Multi-Pane nutzen

TerminalizCrazy unterstuetzt mehrere Panes gleichzeitig.

### Pane teilen

Druecken Sie **Ctrl+\\** fuer einen vertikalen Split:

```
+-------------------------+-------------------------+
|  Chat                   |  Neuer Pane             |
+-------------------------+-------------------------+
| You: list files         |                         |
|                         | > _                     |
| AI: ls -la              |                         |
+-------------------------+-------------------------+
```

Oder **Ctrl+-** fuer einen horizontalen Split.

### Zwischen Panes wechseln

- **Alt+Pfeiltasten**: Navigieren zwischen Panes
- **Ctrl+Z**: Pane zoomen (Vollbild)
- **Ctrl+W**: Aktuellen Pane schliessen

### Layout-Presets

Starten Sie mit einem vorgefertigten Layout:

```bash
terminalizcrazy --layout quad     # 4 Panes
terminalizcrazy --layout tall     # 2 Panes vertikal
terminalizcrazy --layout wide     # 2 Panes horizontal
```

---

## Schritt 6: Agent Mode ausprobieren

Der Agent Mode plant und fuehrt komplexe Aufgaben automatisch aus.

### Agent Mode aktivieren

Druecken Sie **Ctrl+A** um den Modus zu wechseln:
- `off` -> `suggest` -> `auto` -> `off`

Empfohlen: `suggest` (plant, fragt vor Ausfuehrung)

### Komplexe Aufgabe stellen

```
> Set up a new React project with TypeScript and ESLint
```

Der Agent erstellt einen Plan:

```
+-------------------------------------------------------------+
| Agent Mode: Plan erstellt                                   |
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

- **A**: Plan genehmigen und ausfuehren
- **R**: Plan ablehnen
- **M**: Einzelne Tasks anpassen
- **S**: Task ueberspringen

---

## Schritt 7: Theme wechseln

TerminalizCrazy bietet 5 eingebaute Themes.

### Theme in config.toml setzen

```toml
# ~/.terminalizcrazy/config.toml
theme = "dracula"   # default, dracula, monokai, nord, solarized
```

### Verfuegbare Themes

| Theme | Beschreibung |
|-------|--------------|
| `default` | Standard-Farben |
| `dracula` | Dunkles Theme mit Lila-Akzenten |
| `monokai` | Klassisches Entwickler-Theme |
| `nord` | Nordische Farbpalette |
| `solarized` | Augenfreundliche Farben |

### Hot Reload

Themes werden automatisch neu geladen wenn die Datei geaendert wird:

```toml
[appearance]
theme_hot_reload = true
```

---

## Zusammenfassung

Sie haben gelernt:

### Grundlagen

| Taste | Aktion |
|-------|--------|
| Enter | Nachricht senden |
| Ctrl+E | Befehl ausfuehren |
| Ctrl+Y | Befehl kopieren |
| Pfeil hoch/runter | Historie durchsuchen |
| Esc | Beenden |

### Panes & Layouts

| Taste | Aktion |
|-------|--------|
| Ctrl+\\ | Vertikaler Split |
| Ctrl+- | Horizontaler Split |
| Alt+Pfeiltasten | Zwischen Panes wechseln |
| Ctrl+Z | Pane zoomen |
| Ctrl+W | Pane schliessen |

### Agent & Collaboration

| Taste | Aktion |
|-------|--------|
| Ctrl+A | Agent Mode umschalten |
| Ctrl+M | Model Selector oeffnen |
| Ctrl+S | Session teilen |
| Ctrl+J | Session beitreten |
| Ctrl+D | Zusammenarbeit beenden |

---

## Naechste Schritte

Jetzt sind Sie bereit fuer fortgeschrittene Funktionen:

- [Agent-Modus](../anleitungen/agent-modus.md) - Komplexe Aufgaben automatisieren
- [Workflows](../anleitungen/workflows.md) - Wiederkehrende Aufgaben speichern
- [Zusammenarbeit](../anleitungen/zusammenarbeit.md) - Team-Features
- [Einstellungen](../referenz/einstellungen.md) - Alle Konfigurationsoptionen

---

## Tipps fuer den Alltag

### Praezise Fragen stellen

```
Schlecht: "mache etwas mit dateien"
Gut:      "zeige alle Python-Dateien groesser als 1MB"
```

### Kontext nutzen

TerminalizCrazy erkennt Ihr Projekt automatisch:
- In einem Git-Repo: Git-spezifische Vorschlaege
- In einem Node.js-Projekt: npm-Befehle
- In einem Python-Projekt: pip/python-Befehle
- In einem Go-Projekt: go-Befehle

### Bei Unsicherheit: Nachfragen

```
> Was macht der Befehl "tar -xvzf archive.tar.gz"?
```

Die KI erklaert jeden Teil des Befehls.

### Model wechseln

Druecken Sie **Ctrl+M** um zwischen KI-Modellen zu wechseln:
- Gemini Flash (schnell, kosteneffizient)
- Gemini Pro (leistungsfaehiger)
- Claude (Anthropic)
- GPT-4 (OpenAI)

---

*Viel Erfolg mit TerminalizCrazy!*
