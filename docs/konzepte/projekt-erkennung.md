# Projekt-Erkennung

> Intelligente Analyse Ihres Arbeitsverzeichnisses

## Uebersicht

TerminalizCrazy erkennt automatisch den Projekttyp basierend auf Dateien im Arbeitsverzeichnis. Diese Information wird genutzt fuer:

- Kontextbezogene KI-Vorschlaege
- Passende Befehlsempfehlungen
- Session-Namen
- Framework-spezifische Hilfe

---

## Erkannte Projekttypen

| Typ | Erkennungsdateien | Icon |
|-----|-------------------|------|
| Go | go.mod, go.sum | |
| Node.js | package.json | |
| Python | requirements.txt, setup.py, pyproject.toml | |
| Rust | Cargo.toml | |
| Java | pom.xml, build.gradle | |
| Ruby | Gemfile | |
| PHP | composer.json | |
| C/C++ | CMakeLists.txt, Makefile | |
| .NET | *.csproj, *.sln | |
| Docker | Dockerfile, docker-compose.yml | |
| Kubernetes | k8s/, kubernetes/ | |

---

## Framework-Erkennung

Zusaetzlich zum Projekttyp werden Frameworks erkannt:

### Node.js

| Framework | Erkennungsmerkmal |
|-----------|-------------------|
| React | react in dependencies |
| Vue | vue in dependencies |
| Angular | @angular/core in dependencies |
| Next.js | next in dependencies |
| Express | express in dependencies |
| NestJS | @nestjs/core in dependencies |

### Python

| Framework | Erkennungsmerkmal |
|-----------|-------------------|
| Django | django in requirements |
| Flask | flask in requirements |
| FastAPI | fastapi in requirements |
| Pytest | pytest in requirements |

### Go

| Framework | Erkennungsmerkmal |
|-----------|-------------------|
| Gin | gin-gonic/gin in go.mod |
| Echo | labstack/echo in go.mod |
| Fiber | gofiber/fiber in go.mod |

---

## Auswirkungen auf KI-Vorschlaege

### Ohne Projekterkennung

    > Installiere die Abhaengigkeiten
    
    Moegliche Antworten:
    - npm install
    - pip install -r requirements.txt
    - go mod download
    - (unklar welches)

### Mit Projekterkennung (Node.js)

    > Installiere die Abhaengigkeiten
    
    Antwort: npm install
    
    (Die KI weiss, dass es ein Node.js-Projekt ist)

### Mit Framework-Erkennung (Next.js)

    > Starte den Dev-Server
    
    Antwort: npm run dev
    
    (Die KI kennt den Next.js-Standard)

---

## Session-Namen

Erkannte Projekte beeinflussen den Session-Namen:

| Verzeichnis | Erkannt | Session-Name |
|-------------|---------|--------------|
| ~/projects/myapp | Node.js + React | myapp-react |
| ~/go/src/api | Go + Gin | api-go-gin |
| ~/repos/backend | Python + Django | backend-django |

---

## Technische Details

### Detector

    type Detector struct {
        workDir string
    }
    
    func (d *Detector) Detect() *Project {
        // Dateien pruefen
        // Projekttyp bestimmen
        // Framework erkennen
    }

### Project

    type Project struct {
        Name      string
        Type      ProjectType
        Framework string
        Path      string
    }

### Erkennungslogik

1. go.mod vorhanden? -> Go
2. package.json vorhanden? -> Node.js
3. requirements.txt vorhanden? -> Python
4. Cargo.toml vorhanden? -> Rust
5. ...

Bei Node.js zusaetzlich:
1. package.json lesen
2. Dependencies pruefen
3. Framework identifizieren

---

## KI-Kontext

Der Projektkontext wird an die KI gesendet:

    type RequestContext struct {
        CurrentDir       string
        ProjectName      string
        ProjectType      string
        ProjectFramework string
    }

Die KI kann dann passende Befehle vorschlagen.

---

## Grenzen

### Nicht erkannte Projekte

- Projekte ohne typische Markerdateien
- Polyglot-Projekte (mehrere Sprachen)
- Monorepos mit verschachtelter Struktur

### Workarounds

Bei nicht erkannten Projekten:
- Kontext explizit angeben: "In diesem Go-Projekt..."
- Standard-Befehle werden vorgeschlagen
- Keine negativen Auswirkungen

---

## Anpassung

Die Projekterkennung ist derzeit nicht konfigurierbar. Geplante Features:

- [ ] Manuelle Projekt-Deklaration
- [ ] Eigene Erkennungsregeln
- [ ] Monorepo-Unterstuetzung

---

## Beispiele

### Go-Projekt

    Verzeichnis: ~/projects/api
    Dateien: go.mod, go.sum, main.go
    
    Erkannt:
    - Typ: Go
    - Name: api
    - Session: api-go

### React-App

    Verzeichnis: ~/apps/dashboard
    Dateien: package.json (mit react dependency)
    
    Erkannt:
    - Typ: Node.js
    - Framework: React
    - Name: dashboard
    - Session: dashboard-react

### Django-Backend

    Verzeichnis: ~/backend/mysite
    Dateien: requirements.txt, manage.py
    
    Erkannt:
    - Typ: Python
    - Framework: Django
    - Name: mysite
    - Session: mysite-django

---

## Siehe auch

- [AI-Integration](../anleitungen/ai-integration.md) - Wie die KI Kontext nutzt
- [Einstellungen](../referenz/einstellungen.md) - Konfiguration
