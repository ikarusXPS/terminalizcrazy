# Risikostufen

> Befehlssicherheit und Schutzmaßnahmen verstehen

## Übersicht

TerminalizCrazy bewertet jeden Befehl automatisch auf potenzielle Risiken, bevor er ausgeführt wird. Dies schützt vor versehentlicher Ausführung gefährlicher Befehle.

---

## Die vier Risikostufen

### LOW (Niedrig)

| Eigenschaft | Wert |
|-------------|------|
| Farbe | Grün |
| Bestätigung | Keine |
| Auto-Modus | Wird ausgeführt |

Sichere Befehle, die nur Informationen lesen und nichts verändern.

**Beispiele**:
```bash
ls              # Dateien auflisten
cat file.txt    # Datei anzeigen
echo "hello"    # Text ausgeben
pwd             # Aktuelles Verzeichnis
whoami          # Benutzername
grep pattern    # Text suchen
```

---

### MEDIUM (Mittel)

| Eigenschaft | Wert |
|-------------|------|
| Farbe | Orange |
| Bestätigung | Plan-Genehmigung |
| Auto-Modus | Mit Warnung |

Befehle, die Dateien oder Systemzustand modifizieren, aber reversibel sind.

**Beispiele**:
```bash
mv file.txt dir/       # Datei verschieben
cp source dest         # Datei kopieren
mkdir new_folder       # Verzeichnis erstellen
touch file.txt         # Datei erstellen
git commit             # Git Commit
git push               # Git Push
npm install            # Pakete installieren
pip install package    # Python-Paket installieren
wget url               # Datei herunterladen
curl -o file url       # Datei herunterladen
```

---

### HIGH (Hoch)

| Eigenschaft | Wert |
|-------------|------|
| Farbe | Hellrot |
| Bestätigung | Einzelbestätigung erforderlich |
| Auto-Modus | Blockiert |

Destruktive Befehle, die Daten unwiederbringlich löschen oder verändern können.

**Beispiele**:
```bash
rm -rf directory       # Rekursives Löschen
rm -r *                # Alles löschen
rmdir directory        # Verzeichnis löschen
git reset --hard       # Änderungen verwerfen
git clean -fd          # Ungetrackte Dateien löschen
drop table             # Datenbank-Tabelle löschen
truncate table         # Tabelle leeren
npm uninstall          # Paket entfernen
```

---

### CRITICAL (Kritisch)

| Eigenschaft | Wert |
|-------------|------|
| Farbe | Rot |
| Bestätigung | Blockiert im Auto-Modus |
| Auto-Modus | Niemals automatisch |

Systembefehle, die schwerwiegende Auswirkungen haben können.

**Beispiele**:
```bash
sudo command           # Root-Rechte
su                     # Benutzer wechseln
chmod 777              # Gefährliche Berechtigungen
chown                  # Besitzer ändern
rm -rf /               # System löschen
dd if=/dev/zero        # Festplatte überschreiben
mkfs                   # Dateisystem formatieren
shutdown               # System herunterfahren
reboot                 # System neustarten
:(){ :|:& };:          # Fork-Bombe
> /dev/sda             # Direkter Festplattenzugriff
```

---

## Bestätigungsdialog

Bei Befehlen mit Risikostufe MEDIUM oder höher erscheint ein Bestätigungsdialog:

```
┌────────────────────────────────────────────────┐
│ HIGH: This command will delete or destroy data │
│                                                │
│ rm -rf ./build                                 │
│                                                │
│ Execute? [Y]es / [N]o                          │
└────────────────────────────────────────────────┘
```

**Tastenbelegung**:
| Taste | Aktion |
|-------|--------|
| `Y` | Befehl ausführen |
| `N` | Abbrechen |
| `Esc` | Abbrechen |

---

## Agent-Modus und Risiko

Das Verhalten hängt vom konfigurierten Agent-Modus ab:

| Risikostufe | `agent_mode = "off"` | `agent_mode = "suggest"` | `agent_mode = "auto"` |
|-------------|----------------------|--------------------------|------------------------|
| LOW | Immer fragen | Plan erstellen | Automatisch |
| MEDIUM | Immer fragen | Plan erstellen | Mit Warnung |
| HIGH | Immer fragen | Plan erstellen | Blockiert |
| CRITICAL | Immer fragen | Plan erstellen | Blockiert |

**Empfehlung**: Verwenden Sie `agent_mode = "suggest"` für die beste Balance zwischen Komfort und Sicherheit.

---

## Erkennungsmuster

### Kritische Befehle (Auszug)

```go
"sudo", "su ", "chmod 777", "chown", "mkfs",
"dd if=", ":(){ :|:& };:", "rm -rf /",
"format", "> /dev/", "shutdown", "reboot"
```

### Hochriskante Befehle (Auszug)

```go
"rm -rf", "rm -r", "rmdir", "del ", "erase",
"drop table", "drop database", "truncate",
"git reset --hard", "git clean -fd"
```

### Mittlere Risikobefehle (Auszug)

```go
"mv ", "cp ", "rename", "move",
"git push", "git commit", "git checkout",
"npm install", "pip install", "go install"
```

---

## Best Practices

### Do's

- Lesen Sie die Warnung sorgfältig durch
- Überprüfen Sie den Pfad bei `rm`-Befehlen
- Machen Sie Backups vor destruktiven Operationen
- Verwenden Sie `--dry-run` wenn verfügbar

### Don'ts

- Führen Sie nie `rm -rf /` aus
- Kopieren Sie keine Befehle blind aus dem Internet
- Ignorieren Sie keine CRITICAL-Warnungen
- Verwenden Sie `auto`-Modus nicht in Produktionsumgebungen

---

## Anpassung

Die Risikobewertung kann derzeit nicht angepasst werden. Die Patterns sind im Code definiert unter `internal/executor/executor.go`.

---

## Siehe auch

- [Agent-Modus](../anleitungen/agent-modus.md) - Automatisierung verstehen
- [Secret Guard](../konzepte/secret-guard.md) - Geheimnisschutz
- [Einstellungen](einstellungen.md) - Konfigurationsoptionen
