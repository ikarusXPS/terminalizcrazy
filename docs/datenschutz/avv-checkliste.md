# Auftragsverarbeitungsverträge (AVV) - Checkliste

> Erforderliche Vereinbarungen mit AI-Anbietern gemäß DSGVO Art. 28

## Übersicht

TerminalizCrazy überträgt Nutzerdaten an externe AI-Anbieter. Für den
produktiven Einsatz in der EU sind Auftragsverarbeitungsverträge (AVVs)
erforderlich.

---

## Anbieter-Status

| Anbieter | Standort | DPA verfügbar | SCCs | Status |
|----------|----------|---------------|------|--------|
| Google (Gemini) | USA | ✅ Ja | ✅ Integriert | [DPA anfordern](#google-gemini) |
| Anthropic (Claude) | USA | ✅ Ja | ✅ Integriert | [DPA anfordern](#anthropic) |
| OpenAI | USA | ✅ Ja | ✅ Integriert | [DPA anfordern](#openai) |
| Ollama | Lokal | N/A | N/A | Kein AVV nötig |

---

## Google Gemini

### DPA anfordern

1. Google Cloud Console öffnen: https://console.cloud.google.com
2. "Datenschutz" / "Data Processing Terms" aufrufen
3. Google Cloud DPA akzeptieren (inkl. SCCs)

### Relevante Dokumente

- [Google Cloud Data Processing Terms](https://cloud.google.com/terms/data-processing-terms)
- [Google AI Terms of Service](https://ai.google.dev/terms)
- [EU Model Contract Clauses](https://cloud.google.com/terms/eu-model-contract-clause)

### Übermittelte Daten

| Datentyp | Zweck | Rechtsgrundlage |
|----------|-------|-----------------|
| Chat-Nachrichten | AI-Antworten generieren | Art. 6 1b (Vertrag) |
| Projektkontext | Kontextuelle Antworten | Art. 6 1b (Vertrag) |
| Fehler-Nachrichten | Erklärungen generieren | Art. 6 1b (Vertrag) |

---

## Anthropic

### DPA anfordern

1. Anthropic Console öffnen: https://console.anthropic.com
2. Settings > Legal > Data Processing Agreement
3. DPA elektronisch unterzeichnen

Alternativ: E-Mail an privacy@anthropic.com

### Relevante Dokumente

- [Anthropic Privacy Policy](https://www.anthropic.com/privacy)
- [API Terms of Service](https://www.anthropic.com/legal/terms)
- [Data Processing Addendum](https://www.anthropic.com/legal/dpa)

### Übermittelte Daten

| Datentyp | Zweck | Rechtsgrundlage |
|----------|-------|-----------------|
| Chat-Nachrichten | AI-Antworten generieren | Art. 6 1b (Vertrag) |
| System-Prompts | Kontext bereitstellen | Art. 6 1b (Vertrag) |
| Fehler-Logs | Debugging/Support | Art. 6 1f (Interesse) |

### Anthropic Datenverarbeitung

- API-Anfragen werden **nicht** für Training verwendet (API-Nutzung)
- Logs werden 30 Tage aufbewahrt (Standard)
- Zero Data Retention (ZDR) für Enterprise verfügbar

---

## OpenAI

### DPA anfordern

1. OpenAI Platform öffnen: https://platform.openai.com
2. Settings > Data Controls > Data Processing Addendum
3. DPA akzeptieren

### Relevante Dokumente

- [OpenAI Privacy Policy](https://openai.com/policies/privacy-policy)
- [API Data Usage Policies](https://openai.com/policies/api-data-usage-policies)
- [Data Processing Addendum](https://openai.com/policies/data-processing-addendum)

### Übermittelte Daten

| Datentyp | Zweck | Rechtsgrundlage |
|----------|-------|-----------------|
| Chat-Nachrichten | AI-Antworten generieren | Art. 6 1b (Vertrag) |
| Kontext-Daten | Bessere Antworten | Art. 6 1b (Vertrag) |

### OpenAI Datenverarbeitung

- API-Daten werden **nicht** für Training verwendet (seit März 2023)
- 30 Tage Log-Retention (Standard)
- Zero Data Retention auf Anfrage verfügbar

---

## Ollama (Lokal)

Bei Nutzung von Ollama werden **keine Daten** an externe Server übertragen.

- Alle Verarbeitung erfolgt lokal
- Kein AVV erforderlich
- Empfohlen für maximale Datensicherheit

### Konfiguration

```toml
ai_provider = "ollama"
ollama_enabled = true
ollama_model = "codellama"
```

---

## AVV-Checkliste (Art. 28 DSGVO)

Jeder AVV muss folgende Punkte regeln:

### Pflichtinhalte

- [ ] **Gegenstand und Dauer** der Verarbeitung
- [ ] **Art und Zweck** der Verarbeitung
- [ ] **Art der personenbezogenen Daten**
- [ ] **Kategorien betroffener Personen**
- [ ] **Pflichten und Rechte** des Verantwortlichen

### Auftragsverarbeiter-Pflichten

- [ ] Verarbeitung nur auf **dokumentierte Weisung**
- [ ] **Vertraulichkeit** der verarbeitenden Personen
- [ ] Ergreifung aller **Sicherheitsmaßnahmen** (Art. 32)
- [ ] Bedingungen für **Unterauftragnehmer**
- [ ] **Unterstützung** bei Betroffenenrechten
- [ ] **Unterstützung** bei Sicherheitsvorfällen
- [ ] **Löschung/Rückgabe** nach Vertragsende
- [ ] **Nachweispflichten** und Audits

### Internationale Übermittlung

- [ ] **Standardvertragsklauseln (SCCs)** bei USA-Transfer
- [ ] **Ergänzende Maßnahmen** dokumentiert
- [ ] **Transfer Impact Assessment** durchgeführt

---

## Implementierung in TerminalizCrazy

### 1. Datenschutzerklärung aktualisieren

Fügen Sie folgende Passage zur Datenschutzerklärung hinzu:

```markdown
## Externe Dienste

Diese Anwendung nutzt KI-Dienste zur Generierung von Befehlsvorschlägen:

- **Google Gemini** (Google LLC, USA)
- **Anthropic Claude** (Anthropic PBC, USA)
- **OpenAI** (OpenAI LP, USA)

Die Datenübermittlung erfolgt auf Grundlage von Art. 6 Abs. 1 lit. b DSGVO
(Vertragserfüllung) sowie Standardvertragsklauseln (SCCs) gemäß Art. 46
Abs. 2 lit. c DSGVO.

Für lokale Verarbeitung ohne Datenübermittlung steht Ollama zur Verfügung.
```

### 2. Consent-Mechanismus (optional)

Für erhöhte Compliance:

```go
// Beim ersten Start
if !userHasConsented() {
    showConsentDialog(
        "Diese App sendet Ihre Eingaben an externe AI-Dienste. " +
        "Für lokale Verarbeitung wählen Sie Ollama.",
    )
}
```

### 3. Dokumentation aufbewahren

Speichern Sie alle unterzeichneten AVVs unter:
`~/.terminalizcrazy/legal/`

---

## Wartung

| Aufgabe | Häufigkeit |
|---------|------------|
| AVV-Status prüfen | Jährlich |
| Anbieter-ToS prüfen | Bei Änderungen |
| Transfer Impact Assessment | Bei neuen Anbietern |
| Datenschutzerklärung aktualisieren | Bei Änderungen |

---

## Kontakte

| Anbieter | Datenschutz-Kontakt |
|----------|---------------------|
| Google | privacy@google.com |
| Anthropic | privacy@anthropic.com |
| OpenAI | privacy@openai.com |

---

## Siehe auch

- [DSGVO-Analyse](../konzepte/dsgvo-analyse.md)
- [Einstellungen](../referenz/einstellungen.md)
- [AI-Integration](../anleitungen/ai-integration.md)
