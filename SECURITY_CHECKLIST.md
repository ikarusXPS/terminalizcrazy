# Security Checklist - TerminalizCrazy

## Übersicht

Dieses Dokument definiert die Sicherheitsrichtlinien für TerminalizCrazy. Da das Tool mit sensiblen Daten (API-Keys, Terminal-Outputs, Session-Daten) arbeitet, ist Security by Design essentiell.

---

## 1. Secret Management

### .env Nutzung

```bash
# .env Datei (NIEMALS committen!)
ANTHROPIC_API_KEY=sk-ant-xxxxx
OPENAI_API_KEY=sk-xxxxx
ENCRYPTION_KEY=xxxxx
SIGNALING_SERVER_URL=wss://xxxxx
```

**Regeln:**
- [ ] `.env` ist in `.gitignore` eingetragen
- [ ] `.env.example` enthält nur Platzhalter (keine echten Werte)
- [ ] Secrets werden nur via `os.Getenv()` geladen
- [ ] Fehlende Secrets führen zu klarer Fehlermeldung (kein Crash)

### API-Key Handling

```go
// RICHTIG: Environment Variable
apiKey := os.Getenv("ANTHROPIC_API_KEY")
if apiKey == "" {
    return errors.New("ANTHROPIC_API_KEY not set")
}

// FALSCH: Hardcoded
apiKey := "sk-ant-xxxxx" // NIEMALS!
```

**Checkliste:**
- [ ] Keine API-Keys im Code
- [ ] Keine API-Keys in Logs
- [ ] Keine API-Keys in Error Messages
- [ ] API-Keys im Memory nach Nutzung nullen (wenn möglich)

---

## 2. Secret Guard Feature

TerminalizCrazy erkennt und maskiert Secrets automatisch:

### Erkennungsmuster

| Typ | Pattern | Beispiel |
|-----|---------|----------|
| AWS Access Key | `AKIA[0-9A-Z]{16}` | `AKIAIOSFODNN7EXAMPLE` |
| AWS Secret Key | `[A-Za-z0-9/+=]{40}` | Nach "aws_secret" |
| GitHub Token | `gh[pousr]_[A-Za-z0-9]{36,}` | `ghp_xxxx` |
| Anthropic Key | `sk-ant-[A-Za-z0-9-]+` | `sk-ant-api03-xxxx` |
| OpenAI Key | `sk-[A-Za-z0-9]{48,}` | `sk-xxxx` |
| Generic API Key | `[Aa]pi[_-]?[Kk]ey.*[=:]["']?[A-Za-z0-9]{20,}` | `api_key=xxxx` |
| Private Key | `-----BEGIN.*PRIVATE KEY-----` | RSA/EC Keys |
| JWT Token | `eyJ[A-Za-z0-9-_]+\.eyJ[A-Za-z0-9-_]+` | `eyJhbG...` |

### Maskierungsverhalten

```bash
# Original Output (GEFÄHRLICH)
export OPENAI_API_KEY=sk-proj-abc123xyz...

# Maskierter Output (SICHER)
export OPENAI_API_KEY=sk-proj-****...

# Log-Meldung
[SECRET_GUARD] Detected and masked: OpenAI API Key
```

**Checkliste:**
- [ ] Secrets werden vor Screen-Share maskiert
- [ ] Secrets werden vor Logging maskiert
- [ ] Secrets werden vor Collaboration-Übertragung maskiert
- [ ] User kann Maskierung temporär deaktivieren (mit Warnung)

---

## 3. Collaboration Security

### End-to-End Encryption

```
┌─────────────┐         ┌─────────────────┐         ┌─────────────┐
│   User A    │◄───────►│ Signaling Server│◄───────►│   User B    │
│             │  Signal │  (nur Metadata) │  Signal │             │
└──────┬──────┘         └─────────────────┘         └──────┬──────┘
       │                                                    │
       │              WebRTC (E2E Encrypted)                │
       └────────────────────────────────────────────────────┘
                    AES-256-GCM + Argon2
```

**Checkliste:**
- [ ] WebRTC mit DTLS-SRTP für Mediendaten
- [ ] Session-Keys via Argon2 abgeleitet
- [ ] AES-256-GCM für Datenverschlüsselung
- [ ] Signaling Server sieht keine Inhalte
- [ ] Session-Links sind kryptographisch zufällig
- [ ] Links verfallen nach X Stunden (konfigurierbar)

### Session Security

```bash
# Session-Link Format
terminalizcrazy://session/[SESSION_ID]#[ENCRYPTION_KEY]

# Beispiel (Key ist nach # = Fragment, wird nicht an Server gesendet)
terminalizcrazy://session/abc123#AES256KEY
```

**Regeln:**
- [ ] Session-ID ist UUID v4 (nicht erratbar)
- [ ] Encryption Key wird nur via Fragment übertragen
- [ ] Max. 10 Teilnehmer pro Session
- [ ] Session-Owner kann Teilnehmer kicken
- [ ] Rate Limiting auf Session-Join

---

## 4. Logging & Error Handling

### Was geloggt werden darf

```go
// RICHTIG: Keine sensiblen Daten
log.Info("Session started", "session_id", sessionID)
log.Error("API call failed", "status", resp.StatusCode)

// FALSCH: Sensitive Daten im Log
log.Info("API call", "api_key", apiKey)  // NIEMALS!
log.Error("Auth failed", "password", pw)  // NIEMALS!
```

### Log-Level Strategie

| Level | Inhalt | Sensitive Daten |
|-------|--------|-----------------|
| DEBUG | Detaillierter Flow | NIE |
| INFO | Wichtige Events | NIE |
| WARN | Potentielle Probleme | NIE |
| ERROR | Fehler + Context | NIE |

### Error Messages

```go
// RICHTIG: Generisch für User, Detail für Log
if err != nil {
    log.Error("AI request failed", "error", err)
    return errors.New("AI service temporarily unavailable")
}

// FALSCH: Leaking Details
return fmt.Errorf("API call to %s failed with key %s", url, apiKey)
```

**Checkliste:**
- [ ] Keine Stack Traces an User
- [ ] Keine internen Pfade an User
- [ ] Keine API-Endpunkte an User
- [ ] Error-Codes für Debugging (z.B. "Error TC-1234")

---

## 5. Input Validation

### Command Injection Prevention

```go
// RICHTIG: Keine Shell-Interpretation
cmd := exec.Command("git", "status")

// FALSCH: Shell-Injection möglich
cmd := exec.Command("sh", "-c", userInput)
```

### User Input Sanitization

```go
// Vor Verarbeitung von User-Input
func sanitizeInput(input string) string {
    // Null bytes entfernen
    input = strings.ReplaceAll(input, "\x00", "")
    // Control characters entfernen (außer Tab, Newline)
    // Max length enforcing
    if len(input) > MAX_INPUT_LENGTH {
        input = input[:MAX_INPUT_LENGTH]
    }
    return input
}
```

**Checkliste:**
- [ ] Keine direkte Shell-Ausführung von User-Input
- [ ] Max Input Length definiert
- [ ] Special Characters escaped
- [ ] Path Traversal verhindert (`../`)

---

## 6. Storage Security

### Lokale Daten

```
~/.terminalizcrazy/
├── config.toml          # Nicht-sensible Config
├── sessions.db          # SQLite, verschlüsselt at-rest
├── history.db           # SQLite, verschlüsselt at-rest
└── .secrets             # OS Keychain oder verschlüsselt
```

### Encryption at Rest

```go
// Secrets im OS Keychain speichern
import "github.com/keybase/go-keychain"

// Fallback: Verschlüsselte Datei
// Key derivation via PBKDF2/Argon2 aus Machine-ID
```

**Checkliste:**
- [ ] Sensible Daten im OS Keychain (wenn verfügbar)
- [ ] Fallback: AES-256 verschlüsselt
- [ ] DB-Dateien haben 0600 Permissions
- [ ] Keine sensiblen Daten in Plaintext

---

## 7. Dependency Security

### Go Module Security

```bash
# Regelmäßig ausführen
go list -m all | nancy sleuth
govulncheck ./...

# In CI/CD
- name: Security Scan
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

### Supply Chain

**Checkliste:**
- [ ] Dependencies minimieren
- [ ] Nur vertrauenswürdige Packages
- [ ] Dependabot aktiviert
- [ ] Lock-File committet (go.sum)
- [ ] Regelmäßige Audits (govulncheck)

---

## 8. Network Security

### TLS/HTTPS

```go
// Nur HTTPS für API-Calls
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    },
}
```

### Certificate Pinning (Optional für High-Security)

```go
// Für kritische Verbindungen (AI APIs)
pinnedCerts := []string{"sha256/xxxxx"}
```

**Checkliste:**
- [ ] TLS 1.2 Minimum
- [ ] Keine HTTP Fallbacks
- [ ] Certificate Validation aktiv
- [ ] Keine selbstsignierten Certs akzeptieren

---

## 9. Build & Release Security

### Reproducible Builds

```bash
# Build-Info einbetten
go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT"

# Checksums veröffentlichen
sha256sum terminalizcrazy-* > checksums.txt
```

### Code Signing (Später)

```bash
# macOS
codesign --sign "Developer ID" terminalizcrazy

# Windows
signtool sign /f cert.pfx terminalizcrazy.exe
```

**Checkliste:**
- [ ] Builds sind reproduzierbar
- [ ] SHA256 Checksums für alle Releases
- [ ] Signed Releases (später)
- [ ] SBOM (Software Bill of Materials) generieren

---

## 10. Hardening-Maßnahmen

### Memory Safety

```go
// Secrets aus Memory löschen nach Nutzung
defer func() {
    for i := range apiKey {
        apiKey[i] = 0
    }
}()
```

### Rate Limiting

```go
// API Calls rate-limiten
limiter := rate.NewLimiter(rate.Every(time.Second), 10)
if !limiter.Allow() {
    return errors.New("rate limit exceeded")
}
```

### Timeout Enforcement

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

---

## Security Review Checklist (Pre-Commit)

Vor jedem Commit prüfen:

- [ ] `git diff` enthält keine Secrets
- [ ] Keine neuen Dependencies ohne Review
- [ ] Error Messages leaken keine internen Details
- [ ] Input Validation für neue Endpoints
- [ ] Tests für Security-relevanten Code

---

## Incident Response

### Bei Secret-Leak

1. **Sofort**: Secret rotieren (neuen Key generieren)
2. **Audit**: Logs auf unauthorisierte Nutzung prüfen
3. **Cleanup**: Git History bereinigen (BFG Repo-Cleaner)
4. **Review**: Wie konnte es passieren?

### Kontakt

Security Issues an: [security@example.com]
(Wird später eingerichtet)

---

## Quellen & Best Practices

- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [Datadog Sensitive Data Scanner](https://www.datadoghq.com/blog/datadog-coscreen-collaborative-terminal-pair-programming/)
