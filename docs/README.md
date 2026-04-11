# TerminalizCrazy Dokumentation

> KI-gestütztes Terminal für intelligente Befehlsausführung

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   ollama/gemma4   [session]         |
+-------------------------------------------------------------+
|  KI-native Terminal mit Agent Mode, Collaboration,          |
|  SecretGuard und Multi-Provider Support                     |
+-------------------------------------------------------------+
```

## Schnellnavigation

### Erste Schritte

| Dokument | Beschreibung |
|----------|--------------|
| [Installation](erste-schritte/installation.md) | System-Anforderungen und Installation |
| [Schnellstart](erste-schritte/schnellstart.md) | In 5 Minuten loslegen |
| [Tutorial](erste-schritte/tutorial.md) | Interaktives Tutorial (~12 Minuten) |

### Anleitungen

| Dokument | Beschreibung |
|----------|--------------|
| [AI-Integration](anleitungen/ai-integration.md) | KI-Features optimal nutzen |
| [Agent-Modus](anleitungen/agent-modus.md) | Komplexe Aufgaben automatisieren |
| [Zusammenarbeit](anleitungen/zusammenarbeit.md) | Sessions mit anderen teilen |
| [Workflows](anleitungen/workflows.md) | Wiederkehrende Aufgaben speichern |
| [Plugins](anleitungen/plugins.md) | Eigene Plugins entwickeln |

### Referenz

| Dokument | Beschreibung |
|----------|--------------|
| [Einstellungen](referenz/einstellungen.md) | Alle 40+ Konfigurationsoptionen |
| [Tastenkürzel](referenz/tastenkuerzel.md) | Vollständige Tastenbelegung |
| [Themes](referenz/themes.md) | Theme-Anpassung und -Erstellung |
| [Risikostufen](referenz/risikostufen.md) | Befehlssicherheit verstehen |

### Konzepte

| Dokument | Beschreibung |
|----------|--------------|
| [Architektur](konzepte/architektur.md) | System-Übersicht |
| [Secret Guard](konzepte/secret-guard.md) | Automatische Geheimnismaskierung |
| [Projekt-Erkennung](konzepte/projekt-erkennung.md) | Intelligente Projektanalyse |

---

## Was ist TerminalizCrazy?

TerminalizCrazy ist ein **KI-natives Terminal**, das natürliche Sprache in Terminal-Befehle übersetzt. Anstatt sich komplizierte Befehlssyntax zu merken, beschreiben Sie einfach, was Sie tun möchten.

### Hauptfunktionen

- **Multi-Provider Support** - Ollama/Gemma4 (Standard), Gemini, Claude, GPT
- **Agent-Modus** - Automatische Planung und Ausführung komplexer Aufgaben (Ctrl+A)
- **Live Model-Wechsel** - KI-Modell ohne Neustart wechseln (Ctrl+M)
- **Risikobewertung** - Warnung vor gefährlichen Befehlen (4 Stufen)
- **Secret Guard** - Automatische Maskierung von 7 Secret-Typen
- **Multi-Pane Layouts** - Horizontale/vertikale Splits, 5 Layout-Presets
- **Session-Persistenz** - Alle Gespräche werden gespeichert
- **Projekt-Erkennung** - 11 Projekttypen mit Framework-Erkennung
- **Zusammenarbeit** - E2E-verschlüsselte Sessions in Echtzeit teilen

### Beispiel

```
Sie: "zeige mir die größten dateien in diesem ordner"

TerminalizCrazy: Hier ist der Befehl:
  du -ah . | sort -rh | head -20

  Drücken Sie Ctrl+E zum Ausführen
```

---

## Schnellstart

```bash
# 1. Ollama mit Gemma4 starten (Standard, lokal)
ollama pull gemma4
ollama serve

# Oder Cloud-Provider:
# export GEMINI_API_KEY="AIzaSy..." && export AI_PROVIDER="gemini"

# 2. Starten
terminalizcrazy

# 3. Fragen stellen
> wie erstelle ich ein neues git repository?
```

Weitere Details: [Schnellstart-Anleitung](erste-schritte/schnellstart.md)

---

## Konfiguration

TerminalizCrazy wird über `~/.terminalizcrazy/config.toml` konfiguriert.

Vollständige Beispielkonfiguration: [`config.toml.example`](../config.toml.example)

Detaillierte Dokumentation: [Einstellungen](referenz/einstellungen.md)

---

## Support

- **Issues**: [GitHub Issues](https://github.com/terminalizcrazy/terminalizcrazy/issues)
- **Discussions**: [GitHub Discussions](https://github.com/terminalizcrazy/terminalizcrazy/discussions)

---

*Erstellt mit TerminalizCrazy*
