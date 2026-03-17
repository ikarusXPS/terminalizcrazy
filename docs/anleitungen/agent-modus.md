# Agent-Modus

> Komplexe Aufgaben automatisch planen und ausfuehren

## Uebersicht

Der Agent-Modus ermoeglicht es TerminalizCrazy, komplexe Aufgaben in mehrere Schritte aufzuteilen und systematisch abzuarbeiten.

**Analogie**: Stellen Sie sich einen persoenlichen Assistenten vor, der nicht nur eine einzelne Aufgabe erledigt, sondern einen ganzen Arbeitsplan erstellt und Schritt fuer Schritt abarbeitet.

---

## Die drei Modi

### off - Agent deaktiviert

    agent_mode = "off"

- Nur einfache Befehlsvorschlaege
- Jeder Befehl muss einzeln bestaetigt werden
- Keine automatische Planung

**Geeignet fuer**: Lernende, sicherheitskritische Umgebungen

### suggest - Planung mit Bestaetigung (empfohlen)

    agent_mode = "suggest"

- Erstellt ausfuehrliche Plaene fuer komplexe Aufgaben
- Zeigt alle geplanten Schritte vor der Ausfuehrung
- Sie bestaetigen jeden Schritt oder den gesamten Plan

**Geeignet fuer**: Die meisten Anwendungsfaelle

### auto - Automatische Ausfuehrung

    agent_mode = "auto"

- Fuehrt sichere Befehle (LOW-Risk) automatisch aus
- Fragt nur bei riskanten Befehlen nach
- CRITICAL-Befehle werden immer blockiert

**Geeignet fuer**: Erfahrene Nutzer, automatisierte Workflows

---

## Wie funktioniert Planung?

### 1. Komplexe Aufgabe stellen

    > Erstelle ein neues React-Projekt mit TypeScript und ESLint

### 2. Plan wird erstellt

Der Agent analysiert die Aufgabe und erstellt einen Plan:

    Plan: React-Projekt mit TypeScript erstellen
    
    [1] npx create-react-app myapp --template typescript
        Verifikation: Verzeichnis myapp existiert
    
    [2] cd myapp && npm install eslint --save-dev
        Verifikation: package.json enthaelt eslint
    
    [3] npx eslint --init
        Verifikation: .eslintrc.* existiert

### 3. Plan bestaetigen

Je nach Modus:
- **suggest**: Gesamten Plan oder einzelne Schritte bestaetigen
- **auto**: Sichere Schritte werden automatisch ausgefuehrt

### 4. Ausfuehrung mit Verifikation

Nach jedem Schritt prueft der Agent:
- Exit-Code (Erfolg/Fehler)
- Erwartete Ausgabe vorhanden
- Dateien/Verzeichnisse wie erwartet

---

## Task-Verifikation

Jeder Task kann Verifikationskriterien haben:

### Exit Code

    Verifikation: exit_code = 0

Der Befehl muss erfolgreich sein (Exit Code 0).

### Ausgabe enthaelt

    Verifikation: output_contains = "Success"

Die Ausgabe muss den Text enthalten.

### Datei existiert

    Verifikation: file_exists = "package.json"

Die angegebene Datei muss existieren.

### Befehl ausfuehren

    Verifikation: run_command = "test -d node_modules"

Ein Verifikationsbefehl wird ausgefuehrt.

---

## Beispiele

### Projekt einrichten

    > Richte ein Python-Projekt mit virtualenv und pytest ein

Plan:
1. python -m venv venv (Verifikation: venv/ existiert)
2. source venv/bin/activate (Verifikation: VIRTUAL_ENV gesetzt)
3. pip install pytest (Verifikation: pytest in pip list)
4. mkdir tests && touch tests/__init__.py

### Git-Workflow

    > Erstelle einen Feature-Branch, committe die Aenderungen und pushe

Plan:
1. git checkout -b feature/new-feature
2. git add .
3. git commit -m "Add new feature"
4. git push -u origin feature/new-feature

### Deployment

    > Baue das Docker-Image und deploye auf Production

Plan:
1. docker build -t myapp:latest .
2. docker tag myapp:latest registry/myapp:latest
3. docker push registry/myapp:latest
4. kubectl set image deployment/myapp myapp=registry/myapp:latest

---

## Sicherheit

### Maximale Tasks

    agent_max_tasks = 10

Begrenzt die Anzahl von Tasks in einem Plan.

- Niedrigere Werte = konservativer
- Hoehere Werte = mehr Automatisierung

### Risikobewertung

Jeder Task wird einzeln bewertet:

| Risiko | Im suggest-Modus | Im auto-Modus |
|--------|------------------|---------------|
| LOW | Im Plan enthalten | Automatisch |
| MEDIUM | Im Plan enthalten | Mit Warnung |
| HIGH | Einzelbestaetigung | Blockiert |
| CRITICAL | Blockiert | Blockiert |

### Best Practices

1. **Starten Sie mit suggest-Modus**
   Verstehen Sie die Plaene bevor Sie zu auto wechseln

2. **Pruefen Sie jeden Plan**
   Besonders bei Loeschungen oder Systemoperationen

3. **Verwenden Sie auto nur in bekannten Umgebungen**
   Entwicklungs-VMs, Container, etc.

4. **Setzen Sie sinnvolle Limits**
   agent_max_tasks = 5 fuer kritische Systeme

---

## Konfiguration

    # config.toml
    
    # Agent-Modus
    agent_mode = "suggest"
    
    # Maximale Tasks pro Plan
    agent_max_tasks = 10

Siehe auch: [Einstellungen](../referenz/einstellungen.md)

---

## Troubleshooting

### Plan wird nicht erstellt

- Prufen Sie, ob die KI-Verbindung aktiv ist (gruener Status)
- Formulieren Sie die Aufgabe spezifischer

### Tasks schlagen fehl

- Pruefen Sie die Fehlermeldung im Output
- Moeglicherweise fehlen Berechtigungen oder Voraussetzungen

### Zu viele Bestaetigung

- Erhoehen Sie agent_max_tasks
- Wechseln Sie zu auto-Modus (nur wenn sicher)

---

## Siehe auch

- [Risikostufen](../referenz/risikostufen.md) - Befehlssicherheit verstehen
- [Workflows](workflows.md) - Wiederverwendbare Plaene speichern
- [Einstellungen](../referenz/einstellungen.md) - Konfigurationsoptionen
