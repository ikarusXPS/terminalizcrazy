# Schnellstart

> In 5 Minuten loslegen

```
+-------------------------------------------------------------+
|  TerminalizCrazy v0.1.0   ollama/gemma4   [session-id]      |
+-------------------------------------------------------------+
|                                                             |
|  Willkommen! Stellen Sie Ihre erste Frage...                |
|                                                             |
+-------------------------------------------------------------+
| > _                                                         |
+-------------------------------------------------------------+
| Enter: Senden | Ctrl+E: Ausfuehren | Esc: Beenden           |
+-------------------------------------------------------------+
```

## Voraussetzungen

- TerminalizCrazy installiert (siehe [Installation](installation.md))
- Ollama mit Gemma4 (Standard) oder Cloud-Provider konfiguriert

---

## 1. Starten

Oeffnen Sie ein Terminal und starten Sie TerminalizCrazy:

    terminalizcrazy

Sie sollten sehen:
- Header mit Version und KI-Status (gruen = verbunden)
- Eingabefeld am unteren Rand

---

## 2. Erste Anfrage

Geben Sie eine Frage in natuerlicher Sprache ein:

    > Wie finde ich die groessten Dateien in diesem Ordner?

```
+-------------------------------------------------------------+
|  You: Wie finde ich die groessten Dateien?                  |
|                                                             |
|  AI: Hier ist der passende Befehl:                          |
|                                                             |
|     find . -type f -exec du -h {} + | sort -rh | head -20   |
|                                                             |
|     Dieser Befehl findet alle Dateien, sortiert sie nach    |
|     Groesse und zeigt die 20 groessten an.                  |
|                                                             |
|  Press Ctrl+E to execute                                    |
+-------------------------------------------------------------+
```

TerminalizCrazy antwortet mit:
- Dem passenden Terminal-Befehl
- Einer kurzen Erklaerung
- Hinweis: "Press Ctrl+E to execute"

---

## 3. Befehl ausfuehren

Druecken Sie **Ctrl+E** um den vorgeschlagenen Befehl auszufuehren.

Bei riskanten Befehlen erscheint eine Bestaetigung:

```
+-------------------------------------------------------------+
|  MEDIUM: This command will modify files                     |
|                                                             |
|  find . -type f -exec du -h {} + | sort -rh | head -20      |
|                                                             |
|  Execute? [Y]es / [N]o                                      |
+-------------------------------------------------------------+
```

- Druecken Sie **Y** fuer Ja
- Druecken Sie **N** fuer Nein

Das Ergebnis wird im Chat angezeigt.

---

## 4. Befehl kopieren

Druecken Sie **Ctrl+Y** um den letzten Befehl ins Clipboard zu kopieren.

Nuetzlich, wenn Sie den Befehl anpassen oder in einem anderen Terminal verwenden moechten.

---

## 5. Historie durchsuchen

Druecken Sie die **Pfeiltaste nach oben** um vorherige Eingaben durchzublaettern.

Die Historie bleibt ueber Sessions hinweg erhalten.

---

## Wichtige Tastenkuerzel

| Taste | Aktion |
|-------|--------|
| Enter | Nachricht senden |
| Ctrl+E | Letzten Befehl ausfuehren |
| Ctrl+Y | Befehl kopieren |
| Ctrl+L | Chat leeren |
| Pfeil hoch/runter | Historie durchsuchen |
| Esc | Beenden |

Vollstaendige Liste: [Tastenkuerzel](../referenz/tastenkuerzel.md)

---

## Beispiele

### Dateiverwaltung

    > Zeige alle Dateien groesser als 100MB
    > Loesche alle .tmp Dateien
    > Finde doppelte Dateien

### Git

    > Zeige die letzten 5 Commits
    > Erstelle einen neuen Branch namens feature/login
    > Was habe ich heute geaendert?

### Systeminformationen

    > Wie viel Speicherplatz ist noch frei?
    > Welche Prozesse verbrauchen am meisten RAM?
    > Zeige meine IP-Adresse

### Entwicklung

    > Installiere alle npm Dependencies
    > Fuehre die Tests aus
    > Starte den Development Server

---

## Sessions

TerminalizCrazy speichert Ihre Gespraeche automatisch.

Beim naechsten Start koennen Sie:
1. Eine bestehende Session fortsetzen
2. Eine neue Session starten

Sessions werden nach Projekt organisiert (basierend auf Ihrem Arbeitsverzeichnis).

---

## Zusammenarbeit

Teilen Sie Ihre Session mit anderen:

```
+-------------------------------------------------------------+
|  Session teilen                                             |
|                                                             |
|  Share-Code: ABCD-1234                                      |
|                                                             |
|  Teilen Sie diesen Code mit anderen Teilnehmern.            |
|  Sie koennen mit Ctrl+J beitreten.                          |
|                                                             |
|  [Verbunden: 2 Teilnehmer]                                  |
+-------------------------------------------------------------+
```

1. Druecken Sie **Ctrl+S** um zu teilen
2. Ein Share-Code wird angezeigt (z.B. ABCD-1234)
3. Andere druecken **Ctrl+J** und geben den Code ein

Beide sehen alle Nachrichten und Befehle in Echtzeit (E2E-verschluesselt).

---

## Naechste Schritte

- [Tutorial](tutorial.md) - Ausfuehrliches interaktives Tutorial
- [Agent-Modus](../anleitungen/agent-modus.md) - Komplexe Aufgaben automatisieren
- [Einstellungen](../referenz/einstellungen.md) - Anpassungsmoeglichkeiten
