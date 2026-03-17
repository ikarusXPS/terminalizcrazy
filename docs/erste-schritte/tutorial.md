# Tutorial

> Schritt-fuer-Schritt Einfuehrung in TerminalizCrazy

## Uebersicht

Dieses Tutorial fuehrt Sie durch alle wichtigen Funktionen von TerminalizCrazy. Am Ende werden Sie:

- KI-Befehlsvorschlaege verstehen und nutzen
- Befehle sicher ausfuehren koennen
- Die Historie effektiv nutzen
- Sessions teilen und wiederherstellen
- Den Agent-Modus kennen

Dauer: ca. 15-20 Minuten

---

## Schritt 1: Willkommen

Starten Sie TerminalizCrazy:

    terminalizcrazy

Sie sehen:
- Den Header mit Projektname und Version
- KI-Status (gruen = verbunden)
- Session-ID in eckigen Klammern
- Eingabefeld am unteren Rand

Hilfreiche Tastenkuerzel werden in der Fusszeile angezeigt.

---

## Schritt 2: Erste Anfrage

Stellen Sie Ihre erste Frage. Tippen Sie:

    list files

Druecken Sie **Enter**.

Die KI antwortet mit einem passenden Befehl, z.B.:

    ls -la

Unter dem Befehl steht: "Press Ctrl+E to execute"

---

## Schritt 3: Befehl ausfuehren

Druecken Sie **Ctrl+E** um den vorgeschlagenen Befehl auszufuehren.

Da ls ein sicherer Befehl ist, wird er sofort ausgefuehrt.

Das Ergebnis erscheint im Chat mit Zeitstempel.

---

## Schritt 4: Riskante Befehle

Fragen Sie nach einem riskanteren Befehl:

    delete all .log files

Die KI schlaegt vor:

    rm -f *.log

Druecken Sie **Ctrl+E**.

Jetzt erscheint ein Bestaetigungsdialog:

    HIGH: This command will delete or destroy data
    
    rm -f *.log
    
    Execute? [Y]es / [N]o

- Druecken Sie **Y** um auszufuehren
- Druecken Sie **N** oder **Esc** um abzubrechen

Fuer dieses Tutorial: Druecken Sie **N**.

---

## Schritt 5: Historie nutzen

Druecken Sie die **Pfeiltaste nach oben**.

Ihre vorherige Eingabe "list files" erscheint.

Druecken Sie erneut nach oben fuer aeltere Eingaben.

Mit **Pfeiltaste nach unten** gehen Sie wieder vorwaerts.

Druecken Sie **Esc** um die Historie-Navigation zu beenden.

---

## Schritt 6: Clipboard verwenden

Druecken Sie **Ctrl+Y**.

Der letzte vorgeschlagene Befehl wird ins System-Clipboard kopiert.

Eine Bestaetigung erscheint: "Copied to clipboard: rm -f *.log"

Sie koennen den Befehl nun in andere Anwendungen einfuegen.

---

## Schritt 7: Session teilen (optional)

Wenn Sie mit anderen zusammenarbeiten moechten:

1. Druecken Sie **Ctrl+S**
2. Ein Share-Code wird angezeigt (z.B. abcd-1234)
3. Teilen Sie diesen Code mit anderen

Andere Teilnehmer:
1. Druecken **Ctrl+J**
2. Geben den Code ein
3. Sehen alle Nachrichten in Echtzeit

Druecken Sie **Ctrl+D** um die Zusammenarbeit zu beenden.

---

## Schritt 8: Letzten Befehl anzeigen

Druecken Sie **Ctrl+R**.

Der letzte vorgeschlagene Befehl wird im Chat angezeigt, ohne ihn auszufuehren.

Nuetzlich wenn Sie den Befehl nach laengerer Unterhaltung wiederfinden moechten.

---

## Schritt 9: Chat leeren

Druecken Sie **Ctrl+L**.

Der Chat-Verlauf wird geleert.

Die Session bleibt erhalten, nur die Anzeige wird zurueckgesetzt.

---

## Schritt 10: Session wechseln

Beenden Sie mit **Esc** und starten Sie neu:

    terminalizcrazy

Beim Start sehen Sie die Session-Auswahl:
- Pfeil hoch/runter zum Navigieren
- Enter zum Auswaehlen
- N fuer neue Session

Waehlen Sie Ihre vorherige Session, um den Verlauf wiederherzustellen.

---

## Zusammenfassung

Sie haben gelernt:

| Taste | Aktion |
|-------|--------|
| Enter | Nachricht senden |
| Ctrl+E | Befehl ausfuehren |
| Ctrl+Y | Befehl kopieren |
| Ctrl+R | Letzten Befehl anzeigen |
| Ctrl+L | Chat leeren |
| Ctrl+S | Session teilen |
| Ctrl+J | Session beitreten |
| Ctrl+D | Zusammenarbeit beenden |
| Pfeil hoch/runter | Historie durchsuchen |
| Esc | Beenden |

---

## Naechste Schritte

Jetzt sind Sie bereit fuer fortgeschrittene Funktionen:

- [Agent-Modus](../anleitungen/agent-modus.md) - Komplexe Aufgaben automatisieren
- [Workflows](../anleitungen/workflows.md) - Wiederkehrende Aufgaben speichern
- [Zusammenarbeit](../anleitungen/zusammenarbeit.md) - Team-Features
- [Einstellungen](../referenz/einstellungen.md) - Anpassungsmoeglichkeiten

---

## Tipps fuer den Alltag

### Praezise Fragen stellen

Schlecht: "mache etwas mit dateien"
Gut: "zeige alle Python-Dateien groesser als 1MB"

### Kontext nutzen

TerminalizCrazy erkennt Ihr Projekt automatisch:
- In einem Git-Repo: Git-spezifische Vorschlaege
- In einem Node.js-Projekt: npm-Befehle
- In einem Python-Projekt: pip/python-Befehle

### Bei Unsicherheit: Nachfragen

    > Was macht der Befehl "tar -xvzf archive.tar.gz"?

Die KI erklaert jeden Teil des Befehls.

---

*Viel Erfolg mit TerminalizCrazy\!*
