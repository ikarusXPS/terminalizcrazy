# Plugins

> Eigene Plugins entwickeln und verwenden

## Uebersicht

Das Plugin-System ermoeglicht die Erweiterung von TerminalizCrazy durch Hooks. Plugins koennen Befehle modifizieren, Ausgaben filtern und zusaetzliche Funktionen bereitstellen.

---

## Plugin-Konzept

### Hooks

Plugins reagieren auf Events durch Hooks:

| Hook | Beschreibung | Ausfuehrungszeitpunkt |
|------|--------------|----------------------|
| pre_command | Vor Befehlsausfuehrung | Kann Befehl modifizieren oder blockieren |
| post_command | Nach Befehlsausfuehrung | Erhaelt Ausgabe und Exit-Code |
| pre_ai | Vor KI-Anfrage | Kann Prompt modifizieren |
| post_ai | Nach KI-Antwort | Kann Antwort modifizieren |

### Prioritaet

Plugins werden nach Prioritaet sortiert ausgefuehrt:
- Niedrigere Zahl = hoehere Prioritaet
- Standard: 100

---

## Eingebaute Plugins

### SafetyPlugin (Prioritaet: 1)

Blockiert gefaehrliche Befehle.

Funktionen:
- Erkennt kritische Befehle (rm -rf /, etc.)
- Zeigt Warnungen bei riskanten Operationen
- Kann Ausfuehrung blockieren

### AliasPlugin (Prioritaet: 10)

Ersetzt Kurzbefehle durch vollstaendige Befehle.

Standard-Aliase:
- ll -> ls -la
- gs -> git status
- gp -> git push
- gc -> git commit

### TimestampPlugin (Prioritaet: 50)

Fuegt Zeitstempel zur Ausgabe hinzu.

### HistoryLoggerPlugin (Prioritaet: 100)

Speichert Befehlshistorie.

---

## Plugin-Entwicklung

### Grundstruktur

Plugins implementieren das Plugin-Interface:

    type Plugin interface {
        Name() string
        Priority() int
        PreCommand(cmd string) (string, error)
        PostCommand(cmd string, result *Result) error
        PreAI(prompt string) string
        PostAI(response string) string
    }

### Beispiel: Einfaches Plugin

    package myplugin
    
    type MyPlugin struct{}
    
    func (p *MyPlugin) Name() string {
        return "my-plugin"
    }
    
    func (p *MyPlugin) Priority() int {
        return 50
    }
    
    func (p *MyPlugin) PreCommand(cmd string) (string, error) {
        // Befehl vor Ausfuehrung modifizieren
        log.Printf("Executing: %s", cmd)
        return cmd, nil  // unveraendert zurueckgeben
    }
    
    func (p *MyPlugin) PostCommand(cmd string, result *Result) error {
        // Nach Ausfuehrung
        if \!result.Success {
            log.Printf("Command failed: %s", cmd)
        }
        return nil
    }
    
    func (p *MyPlugin) PreAI(prompt string) string {
        return prompt
    }
    
    func (p *MyPlugin) PostAI(response string) string {
        return response
    }

### Plugin registrieren

    func init() {
        plugins.Register(&MyPlugin{})
    }

---

## Anwendungsfaelle

### Befehl-Aliase

    func (p *AliasPlugin) PreCommand(cmd string) (string, error) {
        aliases := map[string]string{
            "ll": "ls -la",
            "gs": "git status",
        }
        
        for alias, full := range aliases {
            if strings.HasPrefix(cmd, alias) {
                return strings.Replace(cmd, alias, full, 1), nil
            }
        }
        return cmd, nil
    }

### Befehl-Validierung

    func (p *ValidatorPlugin) PreCommand(cmd string) (string, error) {
        if strings.Contains(cmd, "DROP TABLE") {
            return "", errors.New("DROP TABLE not allowed")
        }
        return cmd, nil
    }

### Ausgabe-Filterung

    func (p *FilterPlugin) PostCommand(cmd string, result *Result) error {
        // Sensitive Daten aus Ausgabe entfernen
        result.Output = regexp.MustCompile(
            ,
        ).ReplaceAllString(result.Output, "password=***")
        return nil
    }

### KI-Kontext erweitern

    func (p *ContextPlugin) PreAI(prompt string) string {
        // Aktuelles Verzeichnis zum Kontext hinzufuegen
        pwd, _ := os.Getwd()
        return fmt.Sprintf("[In %s] %s", pwd, prompt)
    }

### KI-Antwort formatieren

    func (p *FormatterPlugin) PostAI(response string) string {
        // Befehle hervorheben
        return highlightCommands(response)
    }

---

## Best Practices

### Prioritaeten

| Bereich | Prioritaet |
|---------|-----------|
| Sicherheit | 1-10 |
| Aliase/Transformation | 10-50 |
| Logging/Monitoring | 50-100 |
| Nachbearbeitung | 100+ |

### Fehlerbehandlung

- Fehler in PreCommand blockieren die Ausfuehrung
- Fehler in PostCommand werden geloggt, aber nicht weitergereicht
- Immer graceful degradieren

### Performance

- Keine blockierenden Operationen in Hooks
- Async-Operationen fuer lange Tasks
- Caching wo moeglich

### Testing

    func TestMyPlugin_PreCommand(t *testing.T) {
        p := &MyPlugin{}
        
        result, err := p.PreCommand("test command")
        
        assert.NoError(t, err)
        assert.Equal(t, "test command", result)
    }

---

## Plugin-Konfiguration

Plugins koennen ueber config.toml konfiguriert werden:

    [plugins]
    
    [plugins.safety]
    enabled = true
    block_critical = true
    
    [plugins.alias]
    enabled = true
    custom_aliases = { "d" = "docker", "k" = "kubectl" }
    
    [plugins.timestamp]
    enabled = true
    format = "15:04:05"

---

## Debugging

Debug-Modus aktivieren:

    debug = true
    log_level = "debug"

Alle Plugin-Aktionen werden geloggt:

    [DEBUG] plugin:safety PreCommand: ls -la
    [DEBUG] plugin:alias PreCommand: ll -> ls -la
    [DEBUG] plugin:timestamp PostCommand: added timestamp

---

## Siehe auch

- [Einstellungen](../referenz/einstellungen.md) - Plugin-Konfiguration
- [Architektur](../konzepte/architektur.md) - System-Uebersicht
