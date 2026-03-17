# Workflows

> Wiederkehrende Aufgaben speichern und wiederverwenden

## Uebersicht

Workflows sind gespeicherte Befehlssequenzen, die Sie immer wieder verwenden koennen. Sie sparen Zeit bei repetitiven Aufgaben.

**Analogie**: Wie Makros in Excel, aber fuer Terminal-Befehle.

---

## Konzept

### Was ist ein Workflow?

Ein Workflow besteht aus:
- **Name**: Identifikation (z.B. "deploy-production")
- **Beschreibung**: Was macht dieser Workflow?
- **Schritte**: Sequenz von Befehlen
- **Variablen**: Anpassbare Parameter

### Workflow-Typen

| Typ | Beschreibung |
|-----|--------------|
| Sequential | Befehle nacheinander |
| Conditional | Mit Bedingungen |
| Parametrized | Mit Variablen |

---

## Eingebaute Workflows

### git-feature

Erstellt einen Feature-Branch mit Standard-Setup:

    Schritte:
    1. git checkout main
    2. git pull origin main
    3. git checkout -b feature/{name}

### npm-publish

Publiziert ein npm-Paket:

    Schritte:
    1. npm test
    2. npm version {type}
    3. npm publish

### docker-deploy

Baut und deployed ein Docker-Image:

    Schritte:
    1. docker build -t {image}:{tag} .
    2. docker push {registry}/{image}:{tag}

---

## Workflows verwenden

### Workflow ausfuehren

    > Fuehre den deploy-production Workflow aus

Der Agent laedt den Workflow und zeigt den Plan:

    Workflow: deploy-production
    
    [1] npm test
    [2] npm run build
    [3] docker build -t app:latest .
    [4] docker push registry/app:latest
    [5] kubectl apply -f k8s/

### Mit Parametern

    > Fuehre git-feature mit name=login aus

Der Workflow ersetzt {name} durch "login":

    [1] git checkout -b feature/login

---

## Eigene Workflows erstellen

### YAML-Format

Workflows werden als YAML gespeichert:

    # ~/.terminalizcrazy/workflows/my-workflow.yaml
    name: my-workflow
    description: Mein erster Workflow
    
    variables:
      - name: project
        description: Projektname
        default: myapp
    
    steps:
      - name: Verzeichnis erstellen
        command: mkdir -p {project}
      
      - name: Initialisieren
        command: cd {project} && npm init -y
      
      - name: Dependencies installieren
        command: cd {project} && npm install express

### Speicherort

Workflows werden gespeichert in:

    ~/.terminalizcrazy/workflows/

### Workflow-Struktur

    name: string          # Eindeutiger Name
    description: string   # Beschreibung
    
    variables:            # Optionale Parameter
      - name: string
        description: string
        default: string   # Optional
        required: bool    # Optional, Standard: false
    
    steps:                # Befehlssequenz
      - name: string      # Schrittname
        command: string   # Auszufuehrender Befehl
        condition: string # Optional: Bedingung
        on_error: string  # Optional: continue/stop

---

## Fortgeschrittene Funktionen

### Bedingungen

    steps:
      - name: Build wenn noetig
        command: npm run build
        condition: "test -f package.json"

### Fehlerbehandlung

    steps:
      - name: Tests
        command: npm test
        on_error: stop  # Workflow abbrechen bei Fehler
      
      - name: Lint
        command: npm run lint
        on_error: continue  # Weitermachen trotz Fehler

### Umgebungsvariablen

    steps:
      - name: Mit Umgebung
        command: echo 
        env:
          MY_VAR: "Wert"

---

## Beispiele

### Frontend Deploy

    name: frontend-deploy
    description: Frontend bauen und deployen
    
    variables:
      - name: env
        description: Umgebung (staging/production)
        default: staging
    
    steps:
      - name: Dependencies
        command: npm ci
      
      - name: Build
        command: npm run build:{env}
      
      - name: Deploy
        command: aws s3 sync dist/ s3://bucket-{env}/
      
      - name: Cache invalidieren
        command: aws cloudfront create-invalidation --distribution-id XXX --paths "/*"

### Database Backup

    name: db-backup
    description: Datenbank-Backup erstellen
    
    variables:
      - name: database
        required: true
    
    steps:
      - name: Backup erstellen
        command: pg_dump {database} > backup-20260317.sql
      
      - name: Komprimieren
        command: gzip backup-20260317.sql
      
      - name: Zu S3 hochladen
        command: aws s3 cp backup-*.sql.gz s3://backups/

### Release erstellen

    name: release
    description: Neue Version releasen
    
    variables:
      - name: version
        required: true
    
    steps:
      - name: Tests
        command: npm test
        on_error: stop
      
      - name: Version setzen
        command: npm version {version}
      
      - name: Changelog generieren
        command: conventional-changelog -p angular -i CHANGELOG.md -s
      
      - name: Committen
        command: git add . && git commit -m "Release {version}"
      
      - name: Tag erstellen
        command: git tag v{version}
      
      - name: Pushen
        command: git push && git push --tags

---

## Best Practices

### Benennung

- Kurze, beschreibende Namen
- Kebab-case: "deploy-production", nicht "DeployProduction"
- Praefix fuer Kategorien: "git-", "docker-", "npm-"

### Variablen

- Sinnvolle Standardwerte setzen
- Erforderliche Parameter als required markieren
- Beschreibungen fuer Klarheit

### Fehlerbehandlung

- Kritische Schritte mit on_error: stop
- Optionale Schritte mit on_error: continue
- Am Ende: Aufraeumen

### Dokumentation

- description immer ausfuellen
- Kommentare fuer komplexe Befehle
- README in workflows/ fuer Team-Workflows

---

## Workflow-Verwaltung

### Workflows auflisten

    > Zeige alle Workflows

### Workflow-Details

    > Beschreibe den Workflow deploy-production

### Workflow loeschen

    > Loesche den Workflow test-workflow

---

## Siehe auch

- [Agent-Modus](agent-modus.md) - Automatische Planung
- [Einstellungen](../referenz/einstellungen.md) - Konfiguration
