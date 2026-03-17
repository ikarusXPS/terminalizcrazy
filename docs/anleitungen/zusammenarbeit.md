# Zusammenarbeit

> Sessions in Echtzeit mit anderen teilen

## Uebersicht

TerminalizCrazy ermoeglicht die Echtzeit-Zusammenarbeit ueber verschluesselte WebSocket-Verbindungen. Mehrere Personen koennen dieselbe Session sehen und Befehle ausfuehren.

**Analogie**: Wie Google Docs, aber fuer Terminal-Sessions.

---

## Schnellstart

### Session teilen (Host)

1. Druecken Sie **Ctrl+S** in einer laufenden Session
2. Ein Share-Code wird angezeigt, z.B.: abcd-1234
3. Teilen Sie diesen Code mit anderen

### Session beitreten (Gast)

1. Druecken Sie **Ctrl+J**
2. Geben Sie den Share-Code ein
3. Sie sind verbunden und sehen alles in Echtzeit

### Verbindung beenden

Druecken Sie **Ctrl+D** um die Zusammenarbeit zu beenden.

---

## Funktionen

### Was wird geteilt?

| Element | Geteilt? |
|---------|----------|
| Chat-Nachrichten | Ja |
| KI-Antworten | Ja |
| Ausgefuehrte Befehle | Ja |
| Befehlsausgabe | Ja |
| Benutzer-Status | Ja |

### Was wird NICHT geteilt?

| Element | Grund |
|---------|-------|
| Lokale Dateien | Sicherheit |
| Umgebungsvariablen | Sicherheit |
| API-Schluessel | Werden maskiert |
| History vor Beitritt | Datenschutz |

---

## Sicherheit

### Ende-zu-Ende-Verschluesselung

Alle Nachrichten werden mit AES-256-GCM verschluesselt:

- ECDH-Schluesselaustausch beim Verbindungsaufbau
- Einmaliger Session-Schluessel pro Raum
- Server kann Inhalte nicht lesen

### Share-Code

- 8 Zeichen alphanumerisch (z.B. abcd-1234)
- Gueltig nur waehrend die Host-Session aktiv ist
- Kann nicht wiederverwendet werden

### Benutzeranzeige

In der Statusleiste sehen Sie:
- Anzahl verbundener Benutzer
- Benutzernamen bei Aktionen

---

## Anwendungsfaelle

### Pair Programming

Zwei Entwickler arbeiten gemeinsam an einem Problem:

1. Host teilt die Session
2. Beide sehen die KI-Vorschlaege
3. Beide koennen Befehle vorschlagen
4. Host fuehrt Befehle aus

### Support/Hilfe

Ein erfahrener Kollege hilft einem Anfaenger:

1. Anfaenger teilt seine Session
2. Experte sieht genau, was passiert
3. Experte kann Befehle vorschlagen
4. Anfaenger lernt durch Beobachtung

### Code Review

Team reviewt Aenderungen gemeinsam:

1. Entwickler teilt Session
2. Team sieht Git-Befehle und Diff-Ausgaben
3. Diskussion im Chat
4. Gemeinsame Entscheidungen

### Schulung

Trainer demonstriert Workflows:

1. Trainer teilt Session
2. Teilnehmer beobachten
3. Trainer erklaert jeden Schritt
4. Teilnehmer koennen Fragen stellen

---

## Tastenkuerzel

| Taste | Aktion |
|-------|--------|
| Ctrl+S | Session teilen (Host werden) |
| Ctrl+J | Session beitreten |
| Ctrl+D | Verbindung trennen |

---

## Nachrichten-Typen

### Chat

Normale Textnachrichten zwischen Teilnehmern:

    [Alice]: Hat jemand eine Idee?
    [Bob]: Versuche mal grep -r

### Befehl

Wenn jemand einen Befehl vorschlaegt:

    [Alice] suggested: grep -r "pattern" .

Der Host kann diesen mit Ctrl+E ausfuehren.

### Ausgabe

Befehlsausgaben werden automatisch geteilt:

    [Bob] output:
    file1.txt:pattern found
    file2.txt:pattern here too

### System

Benachrichtigungen ueber Verbindungsstatus:

    Alice joined the session
    Bob left the session

---

## Best Practices

### Fuer Hosts

1. **Kommunizieren Sie vor kritischen Befehlen**
   Warnen Sie, bevor Sie etwas Destruktives tun

2. **Halten Sie sensible Daten fern**
   SecretGuard maskiert, aber vermeiden Sie Risiken

3. **Beenden Sie die Session wenn fertig**
   Ctrl+D um sicher zu schliessen

### Fuer Gaeste

1. **Fragen Sie bevor Sie vorschlagen**
   Nicht ungefragt Befehle senden

2. **Respektieren Sie die Host-Kontrolle**
   Der Host entscheidet, was ausgefuehrt wird

3. **Verlassen Sie sauber**
   Ctrl+D statt einfach das Terminal zu schliessen

---

## Technische Details

### Protokoll

- WebSocket-Verbindung
- JSON-Nachrichten
- Heartbeat alle 30 Sekunden
- Automatische Wiederverbindung bei Unterbrechung

### Server

Der Collaboration-Server laeuft lokal auf Port 8765.

In Zukunft: Zentraler Server fuer entfernte Zusammenarbeit.

### Limits

- Max. 10 Teilnehmer pro Session
- Max. 1000 Nachrichten pro Session
- Session endet wenn Host disconnected

---

## Troubleshooting

### Verbindung schlaegt fehl

- Pruefen Sie die Netzwerkverbindung
- Pruefen Sie, ob der Share-Code korrekt ist
- Pruefen Sie, ob die Host-Session noch aktiv ist

### Nachrichten kommen nicht an

- Netzwerk-Latenz kann Verzoegerungen verursachen
- Bei Problemen: Neu verbinden

### Benutzer werden nicht angezeigt

- Aktualisierung kann verzoegert sein
- Druecken Sie Ctrl+L und verbinden Sie neu

---

## Siehe auch

- [Tastenkuerzel](../referenz/tastenkuerzel.md) - Alle Tasten
- [Schnellstart](../erste-schritte/schnellstart.md) - Grundlagen
