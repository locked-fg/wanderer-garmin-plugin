# Schnellstart: Garmin Plugin installieren (Docker-only)

**In 5 Minuten zum ersten Import!**

---

## Was du brauchst

- Wanderer Installation (im Docker)
- Garmin Connect Account
- [Garmin Personal API Token](https://connect.garmin.com/app/settings/personalApiToken)
- Docker (bereits installiert)

---

## Schritt-fuer-Schritt Anleitung

### Schritt 1: Garmin API Key erstellen (2 Minuten)

1. Oeffne [Garmin Connect Personal API Tokens](https://connect.garmin.com/app/settings/personalApiToken) in deinem Browser
2. Logge dich ein mit deinem Garmin-Konto
3. Klicke auf "Generate Token" (Token generieren)
4. Kopiere den generierten Token (z.B. `GCDC_62a8961e-0f89-448a-8780-8b89c2738095...`)
   - Wichtig: Der Token wird nur einmal angezeigt - speichere ihn sicher!

Fertig! Du hast jetzt deinen Garmin Personal API Key.

---

### Schritt 2: Plugin bauen und installieren (2 Minuten)

**Alle Schritten werden in Docker-Containern ausgefuehrt - es wird NICHTS lokal installiert!**

**Einzeiliger Befehl (empfohlen - kopieren & ausfuehren):**
```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin && \
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build && \
mkdir -p /opt/wanderer-dev/data/plugins/garmin && \
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/ && \
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/ && \
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/ && \
cd /opt/wanderer-dev && docker compose restart db
```

**Oder Schritt-fuer-Schritt:**
```bash
# 1. In das Plugin-Verzeichnis wechseln
cd /opt/wanderer-dev/wanderer-garmin-plugin

# 2. Plugin mit TinyGo Docker-Container bauen
#    (TinyGo wird im Container ausgefuehrt, nicht lokal installiert!)
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build

# 3. Plugin-Dateien in das Wanderer data-Verzeichnis kopieren
mkdir -p /opt/wanderer-dev/data/plugins/garmin
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/

# 4. Wanderer neustarten
docker compose restart db
```

Fertig! Das Plugin ist jetzt installiert.

---

### Schritt 3: Plugin in Wanderer konfigurieren (1 Minute)

1. Oeffne Wanderer in deinem Browser: `http://localhost:8080` (oder deine Domain)
2. Gehe zu: Einstellungen -> Plugins
3. Suche nach: "Garmin"
4. Klicke auf: + Hinzufuegen
5. Konfiguriere:
   - apiKey: Fuege deinen Garmin API-Key ein (aus Schritt 1)
     - Wichtig: Feld heisst **apiKey** (nicht accessToken)
   - Startdatum: (Optional) z.B. 2024-01-01 - importiert nur Aktivitaeten ab diesem Datum
   - Aktiviaten: Aktiviert
   - Routen: Optional (noch experimentell)
6. Speichern und Aktivieren

Fertig! Das Plugin ist jetzt aktiv.

---

### Schritt 4: Erste Synchronisation starten

1. Klicke auf "Jetzt synchronisieren" (oder warte auf die automatische Sync)
2. Warte - die Synchronisation kann je nach Anzahl der Aktiviaten 1-5 Minuten dauern
3. Gehe zu "Trails" - deine Garmin-Aktiviaten sollten jetzt angezeigt werden!

---

## Was wird importiert?

Jede Aktivitaet wird als Trail in Wanderer importiert mit:

- GPX-Track - Komplette Route mit Koordinaten und Hoehenprofil
- Metadaten - Distanz, Dauer, Hoehenmeter, Kalorien, Geschwindigkeit
- Fotos - Aktivitaetsfotos mit Geotags (falls vorhanden)
- Kategorie - Automatische Zuordnung (Laufen, Radfahren, Wandern, etc.)
- Privatsphaere - Behaltet deine Garmin-Einstellungen bei

---

## Tipps & Tricks

### Startdatum nutzen
- Problem: Du hast 1000+ Aktiviaten und die Sync dauert ewig
- Loesung: Setze ein Startdatum (z.B. 2024-01-01) und importiere schrittweise:
  - 1. Sync: 2024-01-01
  - 2. Sync: 2023-01-01
  - 3. Sync: 2022-01-01

### Automatische Synchronisation
- Das Plugin synchronisiert automatisch alle paar Stunden
- Du kannst die Haeufigkeit in den Plugin-Einstellungen anpassen

### Manuelle Synchronisation
- Gehe zu Einstellungen -> Plugins -> Garmin -> Jetzt synchronisieren

### Plugin deaktivieren
- Gehe zu Einstellungen -> Plugins -> Garmin -> Deaktivieren
- Deine importierten Daten bleiben erhalten!

---

## Haeufige Probleme & schnelle Loesungen

| Problem | Loesung |
|---------|--------|
| "Authentifizierung fehlgeschlagen" | API Key pruefen, neuen Token erstellen |
| "Keine Aktiviaten gefunden" | Startdatum pruefen (muss in der Vergangenheit liegen) |
| "Plugin nicht sichtbar" | Wanderer neustarten, Plugin-Dateien pruefen |
| "Rate Limited" | 1 Stunde warten, dann erneut versuchen |
| "GPX fehlt" | Aktivitaet hat keinen Track (z.B. Krafttraining) - nur Metadaten werden importiert |
| Build-Fehler | Stelle sicher, dass du die korrekten Docker-Befehle verwendest (siehe Schritt 2) |

---

## Wichtig: Docker-only Build

**ACHTUNG:** Das Plugin wird **ausschliesslich** in Docker-Containern kompiliert. Es wird **NICHTS** auf deinem lokalen System installiert!

- TinyGo laeuft im Container `tinygo/tinygo:0.44.0`
- Make wird im selbem Container ausgefuehrt
- Alle Abhaengigkeiten sind im Container enthalten
- Dein System bleibt sauber - keine lokale Installation von Go, TinyGo oder Make noetig!

---

## Vollstaendiger Build & Install Befehl

```bash
# Alles in einem Befehl
cd /opt/wanderer-dev/wanderer-garmin-plugin && \
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make clean build && \
mkdir -p /opt/wanderer-dev/data/plugins/garmin && \
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/ && \
cp plugin.json /opt/wanderer-dev/data/plugins/garmin/ && \
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/ && \
docker compose -f /opt/wanderer-dev/docker-compose.yaml restart db
```

---

## Hilfe & Support

- Dokumentation: [README.md](./README.md)
- Wanderer Community: [Discussions](https://github.com/open-wanderer/wanderer/discussions)
- Bug melden: [Issues](https://github.com/open-wanderer/wanderer/issues)

---

## Erfolg!

Du hast das Garmin Plugin erfolgreich installiert und konfiguriert!

Jetzt kannst du alle deine Garmin-Aktiviaten in Wanderer verwalten, analysieren und mit anderen teilen.

Viel Spass beim Erkunden!
