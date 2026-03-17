# AI-Integration

> KI-Features optimal nutzen

## Uebersicht

TerminalizCrazy integriert verschiedene KI-Anbieter fuer intelligente Befehlsvorschlaege. Die KI versteht natuerliche Sprache und generiert passende Terminal-Befehle.

---

## Unterstuetzte Anbieter

### Anthropic (Claude) - empfohlen

| Eigenschaft | Wert |
|-------------|------|
| Modell | Claude 3.5 Sonnet |
| Staerken | Code-Verstaendnis, Praezision |
| Setup | ANTHROPIC_API_KEY |

### OpenAI

| Eigenschaft | Wert |
|-------------|------|
| Modell | GPT-4 |
| Staerken | Breites Wissen |
| Setup | OPENAI_API_KEY |

### Ollama (lokal)

| Eigenschaft | Wert |
|-------------|------|
| Modell | CodeLlama, Llama3, etc. |
| Staerken | Privat, kostenlos |
| Setup | ollama_enabled = true |

---

## Anbieter einrichten

### Anthropic

1. Account erstellen: console.anthropic.com
2. API Key generieren
3. Umgebungsvariable setzen:

       export ANTHROPIC_API_KEY="sk-ant-api03-xxx"

### OpenAI

1. Account erstellen: platform.openai.com
2. API Key generieren
3. Umgebungsvariablen setzen:

       export OPENAI_API_KEY="sk-xxx"
       export AI_PROVIDER="openai"

### Ollama

1. Ollama installieren: ollama.ai
2. Modell herunterladen:

       ollama pull codellama

3. Konfiguration:

       export AI_PROVIDER="ollama"
       export OLLAMA_ENABLED="true"

---

## Effektive Anfragen

### Praezise formulieren

Schlecht:
    > mache was mit dateien

Gut:
    > Zeige alle Python-Dateien groesser als 1MB im aktuellen Verzeichnis

### Kontext nutzen

TerminalizCrazy erkennt automatisch:
- Projekttyp (Go, Node, Python, etc.)
- Git-Status
- Aktuelles Verzeichnis

Die KI beruecksichtigt diesen Kontext automatisch.

### Beispiele fuer gute Anfragen

    > Finde alle TODO-Kommentare im src-Verzeichnis
    > Zeige die letzten 10 Commits von Benutzer alice
    > Erstelle ein tar-Archiv aller .go-Dateien
    > Welche Prozesse verwenden Port 3000?

---

## Projektkontext

Die KI erkennt Ihren Projekttyp und passt Vorschlaege an:

| Projekttyp | Erkannte Dateien | Angepasste Vorschlaege |
|------------|------------------|------------------------|
| Go | go.mod | go build, go test |
| Node.js | package.json | npm, yarn |
| Python | requirements.txt | pip, python |
| Rust | Cargo.toml | cargo |
| Java | pom.xml | mvn |
| Ruby | Gemfile | bundle, gem |

### Beispiel

In einem Node.js-Projekt:

    > Installiere die Dependencies

Vorschlag: npm install (nicht pip install)

---

## Befehlserklaerungen

Fragen Sie nach Erklaerungen:

    > Was macht der Befehl "find . -name '*.txt' -mtime +30 -delete"?

Die KI erklaert jeden Teil:
- find . - Suche ab aktuellem Verzeichnis
- -name - Dateinamen-Muster
- -mtime +30 - Aelter als 30 Tage
- -delete - Gefundene Dateien loeschen

---

## Fehleranalyse

Wenn ein Befehl fehlschlaegt:

    > Der letzte Befehl ist fehlgeschlagen, hilf mir

Die KI analysiert:
- Den ausgefuehrten Befehl
- Die Fehlermeldung
- Moegliche Ursachen und Loesungen

---

## Limits und Kosten

### API-Limits

| Anbieter | Rate Limit | Token Limit |
|----------|-----------|-------------|
| Anthropic | 50 req/min | 100k tokens |
| OpenAI | 60 req/min | 8k tokens |
| Ollama | Unbegrenzt | Modellabhaengig |

### Kosten optimieren

1. **Praezise Anfragen** - Weniger Hin-und-Her
2. **Ollama fuer Entwicklung** - Lokal und kostenlos
3. **Claude fuer Produktion** - Beste Qualitaet

---

## Troubleshooting

### KI antwortet nicht

- Pruefen Sie den API-Key
- Pruefen Sie die Netzwerkverbindung
- Status-Anzeige sollte gruen sein

### Schlechte Vorschlaege

- Formulieren Sie spezifischer
- Geben Sie Kontext an
- Pruefen Sie den Projekttyp

### Hohe Latenz

- Ollama ist schneller fuer einfache Anfragen
- Cloud-APIs haben Netzwerk-Latenz
- Peak-Zeiten vermeiden

---

## Best Practices

1. **Beginnen Sie einfach**
   Testen Sie erst einfache Anfragen

2. **Nutzen Sie den Kontext**
   Bleiben Sie im relevanten Verzeichnis

3. **Bestaetigen Sie kritische Befehle**
   Verlassen Sie sich nicht blind auf KI

4. **Lernen Sie von Vorschlaegen**
   Die KI zeigt oft bessere Alternativen

---

## Siehe auch

- [Einstellungen](../referenz/einstellungen.md) - KI-Konfiguration
- [Agent-Modus](agent-modus.md) - Komplexe Aufgaben
- [Installation](../erste-schritte/installation.md) - API-Setup
