# Themes

> Farbschemas anpassen und eigene Themes erstellen

## Übersicht

TerminalizCrazy unterstützt sowohl eingebaute als auch benutzerdefinierte Themes. Themes werden in YAML-Format definiert und können zur Laufzeit gewechselt werden.

---

## Eingebaute Themes

| Theme | Beschreibung |
|-------|--------------|
| `default` | Dunkles Theme mit lila Akzenten |
| `dracula` | Beliebtes dunkles Theme |
| `monokai` | Klassisches Editor-Theme |
| `nord` | Arktisch-inspiriertes Theme |

### Theme aktivieren

```toml
# In config.toml
theme = "dracula"

# Oder unter [appearance]
[appearance]
theme = "dracula"
```

---

## Benutzerdefinierte Themes

### Speicherort

Benutzerdefinierte Themes werden in `~/.terminalizcrazy/themes/` gespeichert:

```
~/.terminalizcrazy/
└── themes/
    ├── my-theme.yaml
    └── custom-dark.yaml
```

### Theme-Struktur

```yaml
# ~/.terminalizcrazy/themes/my-theme.yaml
name: "My Theme"
author: "Ihr Name"
version: "1.0.0"

colors:
  # Basis-Farben (erforderlich)
  background: "#1e1e2e"
  foreground: "#cdd6f4"
  primary: "#7D56F4"
  
  # Semantische Farben
  secondary: "#04B575"
  warning: "#FFAA00"
  error: "#FF6B6B"
  success: "#04B575"
  muted: "#888888"
  
  # Chat-Farben
  user_message: "#7D56F4"
  ai_message: "#04B575"
  system_message: "#888888"
  
  # UI-Elemente
  pane_border_focused: "#7D56F4"
  pane_border_unfocused: "#888888"
  tab_active: "#7D56F4"
  tab_inactive: "#888888"
  
  # Zusätzliche Farben
  selection: "#44475a"
  comment: "#888888"
  cyan: "#8be9fd"
  green: "#50fa7b"
  orange: "#ffb86c"
  pink: "#ff79c6"
  purple: "#bd93f9"
  red: "#ff5555"
  yellow: "#f1fa8c"
```

---

## Farbpalette

### Erforderliche Farben

| Farbe | Verwendung |
|-------|------------|
| `background` | Hintergrundfarbe |
| `foreground` | Standardtextfarbe |
| `primary` | Hauptakzentfarbe |

### Semantische Farben

| Farbe | Verwendung |
|-------|------------|
| `secondary` | Sekundäre Akzentfarbe |
| `warning` | Warnungen |
| `error` | Fehlermeldungen |
| `success` | Erfolgsmeldungen |
| `muted` | Gedämpfter/deaktivierter Text |

### Chat-Farben

| Farbe | Verwendung |
|-------|------------|
| `user_message` | Eigene Nachrichten |
| `ai_message` | KI-Antworten |
| `system_message` | System-Hinweise |

### UI-Farben

| Farbe | Verwendung |
|-------|------------|
| `pane_border_focused` | Rahmen des aktiven Panels |
| `pane_border_unfocused` | Rahmen inaktiver Panels |
| `tab_active` | Aktiver Tab |
| `tab_inactive` | Inaktive Tabs |
| `selection` | Ausgewählter Text |

---

## Default Theme (Referenz)

```yaml
name: "Default"
author: "TerminalizCrazy"
version: "1.0.0"

colors:
  background: "#1e1e2e"
  foreground: "#cdd6f4"
  primary: "#7D56F4"
  secondary: "#04B575"
  warning: "#FFAA00"
  error: "#FF6B6B"
  success: "#04B575"
  muted: "#888888"
  
  user_message: "#7D56F4"
  ai_message: "#04B575"
  system_message: "#888888"
  
  pane_border_focused: "#7D56F4"
  pane_border_unfocused: "#888888"
  tab_active: "#7D56F4"
  tab_inactive: "#888888"
  
  selection: "#44475a"
  comment: "#888888"
```

---

## Theme-Entwicklung

### Hot-Reload aktivieren

Für schnelles Testen während der Entwicklung:

```toml
[appearance]
theme_hot_reload = true
```

Änderungen an Theme-Dateien werden automatisch übernommen.

### Beispiel: Dracula-inspiriertes Theme

```yaml
# ~/.terminalizcrazy/themes/my-dracula.yaml
name: "My Dracula"
author: "Ihr Name"
version: "1.0.0"

colors:
  background: "#282a36"
  foreground: "#f8f8f2"
  primary: "#bd93f9"
  secondary: "#50fa7b"
  warning: "#ffb86c"
  error: "#ff5555"
  success: "#50fa7b"
  muted: "#6272a4"
  
  user_message: "#ff79c6"
  ai_message: "#50fa7b"
  system_message: "#6272a4"
  
  pane_border_focused: "#bd93f9"
  pane_border_unfocused: "#44475a"
  tab_active: "#bd93f9"
  tab_inactive: "#6272a4"
  
  selection: "#44475a"
  comment: "#6272a4"
  cyan: "#8be9fd"
  green: "#50fa7b"
  orange: "#ffb86c"
  pink: "#ff79c6"
  purple: "#bd93f9"
  red: "#ff5555"
  yellow: "#f1fa8c"
```

### Beispiel: Light Theme

```yaml
# ~/.terminalizcrazy/themes/light.yaml
name: "Light"
author: "TerminalizCrazy"
version: "1.0.0"

colors:
  background: "#ffffff"
  foreground: "#2e3440"
  primary: "#5e81ac"
  secondary: "#a3be8c"
  warning: "#ebcb8b"
  error: "#bf616a"
  success: "#a3be8c"
  muted: "#4c566a"
  
  user_message: "#5e81ac"
  ai_message: "#a3be8c"
  system_message: "#4c566a"
  
  pane_border_focused: "#5e81ac"
  pane_border_unfocused: "#d8dee9"
  tab_active: "#5e81ac"
  tab_inactive: "#4c566a"
  
  selection: "#eceff4"
  comment: "#4c566a"
```

---

## Farbformat

Farben werden als Hex-Codes angegeben:

| Format | Beispiel |
|--------|----------|
| 6-stellig | `#7D56F4` |
| 3-stellig | `#fff` (wird zu `#ffffff`) |

**Groß-/Kleinschreibung**: Beide funktionieren (`#7D56F4` = `#7d56f4`)

---

## Validierung

Ein Theme wird validiert beim Laden. Erforderliche Felder:

- `name` - Theme-Name
- `colors.background` - Hintergrundfarbe
- `colors.foreground` - Vordergrundfarbe
- `colors.primary` - Primärfarbe

Fehlende optionale Farben werden durch sinnvolle Standardwerte ersetzt.

---

## Troubleshooting

### Theme wird nicht geladen

1. Überprüfen Sie den Dateinamen (muss `.yaml` oder `.yml` sein)
2. Überprüfen Sie den Pfad (`~/.terminalizcrazy/themes/`)
3. Überprüfen Sie die YAML-Syntax

### Farben werden nicht angezeigt

1. Stellen Sie sicher, dass Ihr Terminal True Color unterstützt
2. Prüfen Sie die Hex-Farbcodes auf Gültigkeit
3. Aktivieren Sie `debug = true` für Details

---

## Siehe auch

- [Einstellungen](einstellungen.md) - Konfigurationsoptionen
- [Architektur](../konzepte/architektur.md) - System-Übersicht
