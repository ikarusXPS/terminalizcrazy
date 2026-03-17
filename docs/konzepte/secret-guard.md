# Secret Guard

> Automatische Erkennung und Maskierung von Geheimnissen

## Uebersicht

SecretGuard ist ein Sicherheitsmechanismus, der automatisch sensible Daten in der Terminal-Ausgabe erkennt und maskiert. Dies verhindert versehentliche Exposition von API-Keys, Passwoertern und anderen Geheimnissen.

**Analogie**: Wie ein Zensor, der sensible Stellen in Dokumenten schwaerzt, bevor sie angezeigt werden.

---

## Erkannte Secret-Typen

| Typ | Beispiel | Maskiert als |
|-----|----------|--------------|
| AWS Access Key | AKIAIOSFODNN7EXAMPLE | AKIA************MPLE |
| GitHub Token | ghp_xxxxxxxxxxxx | ghp_************xxxx |
| Anthropic Key | sk-ant-api03-xxxxx | sk-a************xxxx |
| OpenAI Key | sk-xxxxxxxxxxxxxxxx | sk-x************xxxx |
| JWT Token | eyJhbGciOiJIUzI1... | eyJa************... |
| Private Key | -----BEGIN PRIVATE KEY----- | (blockiert) |
| Generic API Key | api_key=abc123xyz | api_key=**** |

---

## Funktionsweise

### 1. Scanning

Jede Ausgabe wird mit Regex-Mustern gescannt:

    // AWS Access Key
    AKIA[0-9A-Z]{16}
    
    // GitHub Token
    gh[pousr]_[A-Za-z0-9]{36,}
    
    // JWT
    eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*

### 2. Maskierung

Erkannte Secrets werden maskiert:
- Erste 4 Zeichen bleiben sichtbar
- Mittelteil wird durch **** ersetzt
- Letzte 4 Zeichen bleiben sichtbar

Beispiel:

    Vorher:  sk-ant-api03-abcdefghijklmnop
    Nachher: sk-a********************mnop

### 3. Benachrichtigung

Bei erkannten Secrets wird eine Warnung angezeigt:

    [SecretGuard] 2 secrets detected and masked

---

## Konfiguration

### Aktivieren/Deaktivieren

    # config.toml
    secret_guard_enabled = true

**Empfehlung**: Immer aktiviert lassen.

### Wann deaktivieren?

Nur in sehr speziellen Faellen:
- Debugging von API-Verbindungen
- Lokale Entwicklung mit Test-Keys
- Temporaer fuer Diagnose

**Danach sofort wieder aktivieren.**

---

## Erkennungsmuster

### AWS Access Key

    Muster: AKIA[0-9A-Z]{16}
    Beispiel: AKIAIOSFODNN7EXAMPLE

### GitHub Token (klassisch und fine-grained)

    Muster: gh[pousr]_[A-Za-z0-9]{36,}
    Beispiele:
      ghp_xxxx... (Personal Access Token)
      gho_xxxx... (OAuth Token)
      ghs_xxxx... (Server Token)
      ghr_xxxx... (Refresh Token)

### Anthropic API Key

    Muster: sk-ant-[A-Za-z0-9-]{20,}
    Beispiel: sk-ant-api03-xxxxxxxxxx

### OpenAI API Key

    Muster: sk-[A-Za-z0-9]{32,}
    Beispiel: sk-xxxxxxxxxxxxxxxxxxxxxxxx

### JWT Token

    Muster: eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*
    Beispiel: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM...

### Private Key (PEM)

    Muster: -----BEGIN[A-Z ]+PRIVATE KEY-----
    Beispiel: -----BEGIN RSA PRIVATE KEY-----

### Generische API Keys

    Muster: (api[_-]?key|apikey|api_secret|access_token)[=:]["']?([A-Za-z0-9_-]{20,})["']?
    Beispiele:
      api_key=abcdefghijklmnop
      API-KEY: xyz123456789

---

## Was wird NICHT maskiert?

- Normale Texte und Befehle
- Oeffentliche Schluessel (nur Private Keys)
- Kurze Strings (unter 8 Zeichen)
- UUIDs und IDs (keine typischen Secret-Muster)

---

## Best Practices

### 1. Niemals Secrets commiten

    # .gitignore
    .env
    *.pem
    credentials.json

### 2. Umgebungsvariablen nutzen

    # Statt
    curl -H "Authorization: Bearer sk-xxx"
    
    # Besser
    curl -H "Authorization: Bearer "

### 3. SecretGuard als letzte Verteidigung

SecretGuard ist eine Sicherheitsschicht, aber nicht die einzige:
- Secrets in Env-Vars speichern
- Secret Manager verwenden
- Regelmaessig Secrets rotieren

---

## Troubleshooting

### Secret wird nicht maskiert

- Prufen Sie, ob SecretGuard aktiv ist
- Das Muster koennte nicht erkannt werden
- Melden Sie neue Muster als Feature Request

### Zu viel wird maskiert

- Falsch-Positive sind selten
- Debug-Modus aktivieren fuer Details
- Bei Problemen: secret_guard_enabled = false (temporaer)

### Performance

SecretGuard ist optimiert:
- Regex wird einmal kompiliert
- Nur Text-Ausgaben werden gescannt
- Minimaler Overhead

---

## Technische Details

### Implementation

    type Guard struct {
        enabled bool
    }
    
    func (g *Guard) Scan(text string) []Detection {
        // Pattern matching
    }
    
    func (g *Guard) Mask(text string) string {
        // Replace secrets with masked versions
    }

### Maskierungslogik

    func maskSecret(secret string) string {
        if len(secret) <= 8 {
            return "****"
        }
        prefix := secret[:4]
        suffix := secret[len(secret)-4:]
        middle := strings.Repeat("*", min(len(secret)-8, 20))
        return prefix + middle + suffix
    }

---

## Siehe auch

- [Einstellungen](../referenz/einstellungen.md) - Konfiguration
- [Risikostufen](../referenz/risikostufen.md) - Befehlssicherheit
