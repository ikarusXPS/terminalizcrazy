# Einstellungen

> Vollstaendige Referenz aller Konfigurationsoptionen

## Ueberblick

TerminalizCrazy wird ueber ~/.terminalizcrazy/config.toml konfiguriert. Alle Einstellungen koennen auch ueber Umgebungsvariablen gesetzt werden.

**Beispielkonfiguration**: [config.toml.example](../../config.toml.example)

---

## KI-Anbieter

### ai_provider

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | anthropic |
| Optionen | anthropic, openai, ollama |
| Umgebungsvariable | AI_PROVIDER |

Waehlt den KI-Anbieter fuer Befehlsvorschlaege und Chat.

### anthropic_api_key

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | (leer) |
| Umgebungsvariable | ANTHROPIC_API_KEY |

API-Schluessel fuer Claude (Anthropic). Erforderlich wenn ai_provider = anthropic.

Format: sk-ant-api03-xxxxx

Sicherheitshinweis: Niemals in Versionskontrolle committen\!

### openai_api_key

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | (leer) |
| Umgebungsvariable | OPENAI_API_KEY |

API-Schluessel fuer OpenAI. Erforderlich wenn ai_provider = openai.

---

## Ollama (Lokale KI)

### ollama_enabled

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | false |
| Umgebungsvariable | OLLAMA_ENABLED |

Aktiviert Ollama als KI-Backend.

### ollama_url

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | http://localhost:11434 |
| Umgebungsvariable | OLLAMA_URL |

URL des Ollama-Servers.

### ollama_model

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | codellama |
| Umgebungsvariable | OLLAMA_MODEL |

Zu verwendendes Ollama-Modell.

Empfohlene Modelle:
- codellama - Optimal fuer Code
- llama3 - Allgemeine Aufgaben
- mistral - Schnell und effizient

---

## Agent-Modus

### agent_mode

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | suggest |
| Optionen | off, suggest, auto |
| Umgebungsvariable | AGENT_MODE |

- off: Agent deaktiviert
- suggest: Erstellt Plaene, fuehrt nichts automatisch aus (empfohlen)
- auto: Fuehrt sichere Befehle automatisch aus

### agent_max_tasks

| Eigenschaft | Wert |
|-------------|------|
| Typ | int |
| Standard | 10 |
| Bereich | 1-50 |
| Umgebungsvariable | AGENT_MAX_TASKS |

Maximale Anzahl von Tasks in einem Plan.

---

## Benutzeroberflaeche

### theme

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | default |

Farbschema der Anwendung. Siehe [Themes](themes.md).

### show_welcome

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

Zeigt beim Start Projektinformationen.

---

## Erscheinung [appearance]

### appearance.transparency

| Eigenschaft | Wert |
|-------------|------|
| Typ | float |
| Standard | 1.0 |
| Bereich | 0.0 - 1.0 |

Fenster-Transparenz.

### appearance.enable_animations

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

Aktiviert Animationen.

### appearance.theme_hot_reload

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

Laedt Theme-Aenderungen automatisch neu.

---

## Panel-Einstellungen [pane]

### pane.border_style

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | rounded |
| Optionen | rounded, normal, double, hidden |

### pane.inactive_opacity

| Eigenschaft | Wert |
|-------------|------|
| Typ | float |
| Standard | 0.8 |
| Bereich | 0.0 - 1.0 |

Deckkraft inaktiver Panels.

### pane.show_pane_titles

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

### pane.min_pane_width / min_pane_height

| Eigenschaft | Wert |
|-------------|------|
| Typ | int |
| Standard | 20 / 5 |

Minimale Panel-Groesse.

---

## Workspace-Einstellungen [workspace]

### workspace.default_layout

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | quad |
| Optionen | quad, tall, wide, stack, single |

- quad: 2x2 Raster
- tall: Hauptpanel + 2 seitliche
- wide: Oberes + 2 untere
- stack: 4 vertikal gestapelt
- single: Einzelnes Panel

### workspace.auto_save

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

### workspace.restore_on_startup

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

---

## Sicherheit

### secret_guard_enabled

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | true |

Aktiviert automatische Geheimnismaskierung.

Erkannte Typen: AWS Keys, GitHub Tokens, JWT, API Keys, Private Keys

Siehe [Secret Guard](../konzepte/secret-guard.md).

---

## Speicher und Historie

### history_limit

| Eigenschaft | Wert |
|-------------|------|
| Typ | int |
| Standard | 1000 |
| Bereich | 100-10000 |

---

## Debugging

### debug

| Eigenschaft | Wert |
|-------------|------|
| Typ | bool |
| Standard | false |
| Umgebungsvariable | DEBUG |

### log_level

| Eigenschaft | Wert |
|-------------|------|
| Typ | string |
| Standard | info |
| Optionen | debug, info, warn, error |

---

## Umgebungsvariablen

| Variable | Einstellung |
|----------|-------------|
| ANTHROPIC_API_KEY | anthropic_api_key |
| OPENAI_API_KEY | openai_api_key |
| AI_PROVIDER | ai_provider |
| OLLAMA_ENABLED | ollama_enabled |
| OLLAMA_URL | ollama_url |
| OLLAMA_MODEL | ollama_model |
| AGENT_MODE | agent_mode |
| DEBUG | debug |
| LOG_LEVEL | log_level |

Umgebungsvariablen haben Vorrang vor config.toml.

---

## Siehe auch

- [config.toml.example](../../config.toml.example)
- [Tastenkuerzel](tastenkuerzel.md)
- [Themes](themes.md)
