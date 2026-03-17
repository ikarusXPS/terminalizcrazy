package tutorial

// WelcomeMessage is the initial tutorial welcome text
const WelcomeMessage = """
Willkommen bei TerminalizCrazy!

Dieses kurze Tutorial zeigt Ihnen die wichtigsten Funktionen.
Sie koennen das Tutorial jederzeit mit Esc ueberspringen.

Druecken Sie Enter um zu beginnen...
"""

// CompletionMessage is shown when the tutorial is completed
const CompletionMessage = """
Tutorial abgeschlossen!

Sie kennen jetzt die wichtigsten Funktionen:
- KI-Befehlsvorschlaege
- Sichere Befehlsausfuehrung
- History-Navigation
- Clipboard-Integration
- Session-Zusammenarbeit

Viel Erfolg mit TerminalizCrazy!
"""

// WelcomeDescription describes the welcome step
const WelcomeDescription = """
TerminalizCrazy ist ein KI-gestuetztes Terminal.
Beschreiben Sie was Sie tun moechten und erhalten Sie den passenden Befehl.
"""

// FirstQuestionDescription describes the first question step
const FirstQuestionDescription = """
Stellen Sie Ihre erste Frage an die KI.
Tippen Sie eine Anfrage in natuerlicher Sprache, z.B. "list files".
Die KI wird einen passenden Terminal-Befehl vorschlagen.
"""

// ExecuteCommandDescription describes the execute command step
const ExecuteCommandDescription = """
Die KI hat einen Befehl vorgeschlagen.
Mit Ctrl+E fuehren Sie den Befehl aus.
Sichere Befehle werden sofort ausgefuehrt.
"""

// RiskConfirmationDescription describes the risk confirmation step
const RiskConfirmationDescription = """
Bei riskanten Befehlen erscheint eine Bestaetigung.
LOW-Risk: Sofortige Ausfuehrung
MEDIUM-Risk: Plan-Genehmigung
HIGH-Risk: Einzelbestaetigung erforderlich
CRITICAL: Wird im Auto-Modus blockiert
"""

// HistoryNavigationDescription describes the history navigation step
const HistoryNavigationDescription = """
Mit den Pfeiltasten navigieren Sie durch Ihre History.
Pfeil hoch: Vorherige Eingaben
Pfeil runter: Neuere Eingaben
Die History bleibt ueber Sessions erhalten.
"""

// ClipboardDescription describes the clipboard step
const ClipboardDescription = """
Mit Ctrl+Y kopieren Sie den letzten Befehl ins Clipboard.
So koennen Sie Befehle in andere Anwendungen einfuegen
oder anpassen bevor Sie sie ausfuehren.
"""

// SessionSharingDescription describes the session sharing step
const SessionSharingDescription = """
Mit Ctrl+S starten Sie die Zusammenarbeit.
Ein Share-Code wird generiert (z.B. abcd-1234).
Andere koennen mit Ctrl+J beitreten.
Alle sehen Nachrichten und Befehle in Echtzeit.
"""

// SessionRestoreDescription describes the session restore step
const SessionRestoreDescription = """
Beim Start koennen Sie eine vorherige Session laden.
Mit den Pfeiltasten navigieren Sie durch verfuegbare Sessions.
Sessions im gleichen Verzeichnis werden mit einem Stern markiert.
"""

// CompleteDescription describes the completion step
const CompleteDescription = """
Sie sind bereit!

Wichtige Tastenkuerzel:
Ctrl+E  - Befehl ausfuehren
Ctrl+Y  - Befehl kopieren
Ctrl+S  - Session teilen
Ctrl+J  - Session beitreten
Ctrl+L  - Chat leeren
Esc     - Beenden

Stellen Sie jetzt Ihre erste echte Frage!
"""

// SkipMessage is shown when the user skips the tutorial
const SkipMessage = "Tutorial uebersprungen. Tippen Sie /help fuer Hilfe."

// ProgressFormat is the format string for progress display
const ProgressFormat = "Schritt %d von %d"
