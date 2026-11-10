# Garmin Plugin für Wanderer

**Importiere deine Garmin-Aktivitäten und -Routen in Wanderer!**

Dieses Plugin verbindet [Wanderer](https://wanderer.to) mit der [Garmin Personal API](https://developer.garmin.com/gc-developer-program/personal/personal-api/), um deine Aktivitäten und zugehörige Daten (GPX-Tracks, Metadaten) automatisch zu synchronisieren. Fotos und geplante Routen werden derzeit nicht importiert.

---

## 💡 Fertiges Plugin-Paket: Kein TinyGo nötig!

Das Plugin ist bereits gebaut (`dist/garmin/` in diesem Repository: `plugin.json`, `plugin.wasm`, `icon.svg` – gebaut mit TinyGo 0.41.1, derselben Version wie der Wanderer-CI-Build). Für die Installation musst du **nur die 3 Dateien in das Plugin-Verzeichnis von Wanderer kopieren** – ganz ohne Build im Container. Details siehe [Installation](#-installation).

---

## 📋 Inhaltsverzeichnis

- [🎯 Funktionen](#-funktionen)
- [📥 Voraussetzungen](#-voraussetzungen)
- [🔧 Installation](#-installation)
- [⚙️ Konfiguration](#️-konfiguration)
- [🚀 Nutzung](#-nutzung)
- [🔄 Synchronisationsoptionen](#-synchronisationsoptionen)
- [🎨 Unterstützte Aktivitätstypen](#-unterstützte-aktivitätstypen)
- [❓ Häufige Probleme & Lösungen](#-häufige-probleme--lösungen)
- [📚 Technische Details](#-technische-details)
- [🤝 Mitwirken](#-mitwirken)
- [📜 Lizenz](#-lizenz)

---

## 🎯 Funktionen

✅ **Aktivitäten importieren** – Alle deine Garmin-Aktivitäten werden als Wanderer-Trails importiert (end-to-end im Docker-Wanderer verifiziert)  
✅ **GPX-Tracks** – Vollständige Track-Daten mit Elevation, Zeitstempeln und Koordinaten (aus `geoPolylineDTO.polyline` rekonstruiert)  
❌ **Fotos** – werden derzeit nicht importiert (kein Endpunkt in der Personal API verifiziert)  
✅ **Metadaten** – Distanz, Dauer, Höhenmeter, Kalorien, Durchschnitts- und Maximalgeschwindigkeit  
✅ **Kategorien** – Automatische Zuordnung zu Wanderer-Kategorien (Laufen, Radfahren, Wandern, etc.)  
✅ **Privatsphäre** – Behält deine Garmin-Privatsphäre-Einstellungen bei  
✅ **Schwierigkeitsgrade** – easy, moderate, difficult  
✅ **Startdatum-Filter** – Importiere nur Aktivitäten ab einem bestimmten Datum  
✅ **Inkrementell** – Bereits importierte Aktivitäten werden bei Folge-Syncs übersprungen (Dedup über `trail_external_reference`)  

---

## 📥 Voraussetzungen

Bevor du das Plugin installierst, stelle sicher, dass du folgende Dinge hast:

### 1. **Wanderer Installation**
- Eine laufende Wanderer-Instanz im Docker
- Zugriff auf das Wanderer-Verzeichnis
- Plugin-System aktiviert

### 2. **Garmin Account**
- Ein [Garmin Connect](https://connect.garmin.com) Konto
- Mindestens eine Aktivität in deinem Garmin-Konto

### 3. **Garmin Personal API Key**
- Ein gültiger **Personal API Token** von Garmin
- 🔗 [Hier erstellen](https://connect.garmin.com/app/settings/personalApiToken)

### 4. **Docker**
- Docker muss auf deinem System installiert sein
- Dein Wanderer-Container muss laufen

---

## 🔧 Installation

### 🐳 Fertiges Plugin installieren (empfohlen – kein Build nötig!)

Kopiere die 3 Dateien aus `dist/garmin/` in das Plugin-Verzeichnis von Wanderer (`data/plugins/garmin/`, im Container `/data/plugins/garmin/`):

```bash
# Variante A: direkt auf dem Host (data/ als Bind-Mount)
mkdir -p /pfad/zu/deinem/wanderer/data/plugins/garmin
cp dist/garmin/plugin.json dist/garmin/plugin.wasm dist/garmin/icon.svg /pfad/zu/deinem/wanderer/data/plugins/garmin/
docker compose restart wanderer-db

# Variante B: in den laufenden Container kopieren
docker exec wanderer-db mkdir -p /data/plugins/garmin
docker cp dist/garmin/plugin.json wanderer-db:/data/plugins/garmin/plugin.json
docker cp dist/garmin/plugin.wasm wanderer-db:/data/plugins/garmin/plugin.wasm
docker cp dist/garmin/icon.svg  wanderer-db:/data/plugins/garmin/icon.svg
docker compose restart wanderer-db
```

✅ **Fertig!** Das Plugin ist jetzt installiert und erscheint nach dem Neustart unter **Einstellungen → Plugins**.

---

### 🔨 Aus dem Quellcode bauen (optional)

Nur nötig, um das Plugin selbst zu kompilieren. Dafür brauchst du **Go 1.25** und **TinyGo 0.41.1** auf dem Build-Host, und das Plugin muss im Wanderer-Repository unter `plugins/garmin` liegen (wegen `replace ../sdk` in `go.mod`):

```bash
git clone https://github.com/open-wanderer/wanderer.git
cp -r wanderer-garmin-plugin wanderer/plugins/garmin
cd wanderer/plugins/garmin
make build           # -> dist/garmin/{plugin.json,plugin.wasm,icon.svg}
make install-local   # -> ../../data/plugins/garmin/
```

---

## ⚙️ Konfiguration

### Schritt 1: Plugin aktivieren

1. Öffne Wanderer in deinem Browser
2. Gehe zu **Einstellungen → Plugins**
3. Suche nach dem **Garmin Plugin**
4. Klicke auf **Hinzufügen** oder **Aktivieren**

### Schritt 2: API Key einrichten

1. Klicke auf das **⚙️ Symbol** neben dem Garmin Plugin
2. Gib deinen **Garmin Personal API Key** ein (Feld `apiKey`):
   - Den Key hast du zuvor unter [https://connect.garmin.com/app/settings/personalApiToken](https://connect.garmin.com/app/settings/personalApiToken) erstellt
   - Der Key sieht so aus: `GCDC_...`
   - **Hinweis:** Das Plugin tauscht den Key automatisch gegen ein 24h gültiges Access Token (`POST /personalapitoken-service/...` mit `Authorization: ApiKey <key>`) und verwendet dieses als Bearer-Token. Du musst nichts weiter tun.
3. (Optional) Setze ein **Startdatum**:
   - Nur Aktivitäten ab diesem Datum werden importiert
   - Nützlich, um alte Aktivitäten schrittweise zu synchronisieren
   - Format: `YYYY-MM-DD` (z.B. `2024-01-01`)

### Schritt 3: Synchronisation aktivieren

1. Aktiviere den Schalter **Aktiviert**
2. Wähle aus, was synchronisiert werden soll:
   - ✅ **Aktivitäten** (empfohlen)
   - ⚠️ **Routen** (derzeit experimentell)
3. Klicke auf **Speichern**

---

## 🚀 Nutzung

### Erste Synchronisation

Nach der Konfiguration:

1. Klicke auf **Jetzt synchronisieren** (oder warte auf die automatische Synchronisation)
2. Wanderer lädt deine Garmin-Aktivitäten herunter
3. Die Aktivitäten erscheinen in deiner **Trail-Liste**

### Automatische Synchronisation

Der Sync läuft als Cron im `wanderer-db`-Container – standardmäßig **nächtlich (`0 2 * * *`)**, änderbar über `POCKETBASE_CRON_SYNC_SCHEDULE` in der `.env` der Wanderer-Installation.

⚠️ **Nicht überlappen lassen:** Syncs werden nicht serialisiert; ein Cron-Intervall, das kürzer ist als eine Sync dauert, führt zu doppelten Importversuchen (`external_id … Value must be unique`). Das Standardintervall ist sicher.

### Manuelle Synchronisation

Um manuell zu synchronisieren:

1. Gehe zu **Einstellungen → Plugins**
2. Finde das Garmin Plugin
3. Klicke auf **Jetzt synchronisieren**

---

## 🔄 Synchronisationsoptionen

| Option | Beschreibung | Empfohlen |
|--------|-------------|-----------|
| **Startdatum** | Importiere nur Aktivitäten ab diesem Datum | ✅ Ja |
| **Aktivitäten** | Synchronisiere completed Aktivitäten | ✅ Ja |
| **Routen** | Geplante Routen – von der Personal API **nicht unterstützt** (das Plugin meldet `not_supported`) | ❌ Nein |
| **Privatsphäre** | Behalte Garmin-Privatsphäre-Einstellungen bei | ✅ Ja |

### Tipps für große Aktivitätsmengen

Wenn du viele Aktivitäten hast (z.B. 1000+):

1. **Setze ein Startdatum** (z.B. `2024-01-01`) um nur aktuelle Aktivitäten zu importieren
2. **Deaktiviere das Plugin nach dem ersten Import** und aktiviere es später wieder für neue Aktivitäten
3. **Wiederhole den Import** mit verschiedenen Startdaten, um alle Aktivitäten schrittweise zu importieren

Beispiel:
- Erster Import: Startdatum `2024-01-01`
- Zweiter Import: Startdatum `2023-01-01`
- usw.

---

## 🎨 Unterstützte Aktivitätstypen

Das Plugin mappt Garmin-Aktivitätstypen automatisch zu Wanderer-Kategorien:

| Garmin Typ | Wanderer Kategorie | Subkategorie |
|------------|-------------------|--------------|
| `running` | Running | Road |
| `trail_running` | Running | Trail |
| `cycling` | Biking | Road |
| `mountain_biking` | Biking | MTB |
| `gravel_cycling` | Biking | Gravel |
| `e_bike_mountain` | Biking | E-Bike |
| `e_bike_road` | Biking | E-Bike |
| `hiking` | Hiking | - |
| `walking` | Walking | - |
| `swimming` | Other | - |
| `skiing` | Skiing | - |
| `alpine_skiing` | Skiing | Alpine |
| `backcountry_skiing` | Skiing | Backcountry |
| `cross_country_skiing` | Skiing | Cross-country |
| `snowboarding` | Skiing | - |
| `kayaking` | Canoeing | - |
| `rowing` | Canoeing | - |
| `paddle_boarding` | Canoeing | - |
| `climbing` | Climbing | - |
| `rock_climbing` | Climbing | - |
| `ice_climbing` | Climbing | - |

> ⚠️ **Hinweis:** Falls ein Aktivitätstyp nicht in der Liste ist, wird er als "Other" importiert.

---

## ❓ Häufige Probleme & Lösungen

### 🔴 **"Authentication Failed" / "Authentifizierung fehlgeschlagen"**

**Ursache:**
- Ungültiger, abgelaufener oder bei Garmin **gelöschter** API Key (Token-Endpunkt antwortet dann mit `404 NotFoundException` statt `401`)
- API Key wurde nicht korrekt kopiert

**Lösung:**
1. Überprüfe, dass dein API Key korrekt ist
2. Erstelle einen **neuen API Key** unter [https://connect.garmin.com/app/settings/personalApiToken](https://connect.garmin.com/app/settings/personalApiToken)
3. Kopiere den Key **ohne Leerzeichen** am Anfang oder Ende
4. Füge den Key erneut in den Plugin-Einstellungen ein

---

### 🔴 **"No activities found" / "Keine Aktivitäten gefunden"**

**Ursache:**
- Startdatum ist zu neu (in der Zukunft)
- Keine Aktivitäten in deinem Garmin-Konto
- API Key hat keine Berechtigungen

**Lösung:**
1. Überprüfe dein **Startdatum** – es sollte in der Vergangenheit liegen
2. Prüfe, ob du Aktivitäten in deinem [Garmin Connect](https://connect.garmin.com) hast
3. Erstelle eine Test-Aktivität und versuche erneut zu synchronisieren
4. Teste deinen API Key mit [curl](https://curl.se/) – **zweistufig**, der `GCDC_...`-Key ist kein Bearer-Token:
   ```bash
   # 1) Key gegen 24h-JWT tauschen
   JWT=$(curl -s -X POST \
     -H "Authorization: ApiKey DEIN_GCDC_KEY" \
     -H "Content-Type: application/json" -d '{}' \
     "https://personal.connectapi.garmin.com/personalapitoken-service/personal-api/personal-api-key/token")
   # 2) JWT als Bearer verwenden
   curl -H "Authorization: Bearer $JWT" \
     "https://personal.connectapi.garmin.com/activitylist-service/activities/search/activities?start=0&limit=10"
   ```
   Antwortet Schritt 1 mit `404 NotFoundException` ("Personal API Key with ID … does not exist"), wurde der Key bei Garmin gelöscht – neuen erstellen.

---

### 🔴 **"Rate Limited" / "Zu viele Anfragen"**

**Ursache:**
- Garmin hat Rate Limits für die API
- Zu viele Synchronisationsversuche in kurzer Zeit

**Lösung:**
1. Warte **1 Stunde** und versuche es erneut
2. Reduziere die **Synchronisationshäufigkeit** in den Plugin-Einstellungen
3. Verwende ein **Startdatum**, um weniger Aktivitäten auf einmal zu importieren

---

### 🔴 **"Plugin nicht sichtbar" nach der Installation**

**Ursache:**
- Wanderer wurde nicht neu gestartet
- Plugin-Dateien sind im falschen Verzeichnis

**Lösung:**
1. **Starte Wanderer neu** (wichtig!)
2. Überprüfe, ob die Plugin-Dateien im richtigen Verzeichnis sind:
   ```
   data/plugins/garmin/
   ├── plugin.json
   ├── plugin.wasm
   └── icon.svg
   ```
3. Überprüfe die Dateiberechtigungen:
   ```bash
   docker exec -it wanderer-db chmod -R 755 /data/plugins/garmin
   ```

---

### 🔴 **"GPX-Datei fehlt" oder "Track nicht verfügbar"**

**Ursache:**
- Die Aktivität hat keinen GPX-Track
- Garmin hat die Aktivitätsdatei nicht gespeichert

**Lösung:**
1. Überprüfe in Garmin Connect, ob die Aktivität einen Track hat
2. Manche Aktivitäten (z.B. Krafttraining) haben keinen GPX-Track – das ist normal
3. Das Plugin importiert trotzdem die Metadaten (Distanz, Dauer, etc.)

---

### 🔴 **Build-Fehler beim Selbstkompilieren ("tinygo: command not found" / "../sdk" nicht gefunden)**

**Ursache:**
- TinyGo 0.41.1 ist nicht auf dem Build-Host installiert, oder das Plugin liegt nicht innerhalb des wanderer-Repos

**Lösung:**
- Nutze das fertige Paket aus `dist/garmin/` (kein Build nötig), oder baue auf einem Host mit Go 1.25 + TinyGo 0.41.1, wobei das Plugin unter `plugins/garmin` im wanderer-Repository liegen muss:
```bash
git clone https://github.com/open-wanderer/wanderer.git
cp -r wanderer-garmin-plugin wanderer/plugins/garmin
cd wanderer/plugins/garmin && make build && make install-local
```

---

## 📚 Technische Details

### API Endpunkte

Das Plugin nutzt folgende Garmin Personal API Endpunkte (verifiziert gegen die Live-API):

| Endpunkt | Zweck |
|----------|-------|
| `POST /personalapitoken-service/personal-api/personal-api-key/token` | API Key gegen 24h-JWT tauschen (`Authorization: ApiKey <key>`) |
| `/activitylist-service/activities/search/activities` | Aktivitäten suchen/auflisten (`start`, `limit`, `startDate`) |
| `/activity-service/activity/{id}` | Aktivitätsdetails (summaryDTO, activityTypeDTO, accessControlRuleDTO) |
| `/activity-service/activity/{id}/details` | Track als `geoPolylineDTO.polyline` (wird zu GPX gebaut) |

### Base URL

```
https://personal.connectapi.garmin.com
```

### Authentifizierung

- Der Personal API Key (`GCDC_...`) aus den Garmin-Einstellungen ist **kein** Bearer-Token
- Das Plugin tauscht ihn per Token-Endpunkt gegen ein **24h gültiges JWT** und nutzt danach `Authorization: Bearer <jwt>`
- Der GPX-Download (`/download-service/...`) liefert nur ein ZIP mit **FIT**-Datei – deshalb wird der Track aus `geoPolylineDTO` als GPX rekonstruiert
- Gelöschter Key: Token-Endpunkt antwortet `404 NotFoundException` (nicht 401)

### Host-Vertrag (Wanderer-Plugin-System, live im Docker getestet)

- **Fehler** müssen über den Extism-Fehlerkanal gemeldet werden: `pdk.SetErrorString` mit JSON `{"code": ..., "message": ...}` bei Rückgabewert ≠ 0 – nicht über `pdk.OutputJSON`, sonst sieht der Host nur `plugin_error: call <export> failed with code 1`
- **HTTP-Status selbst prüfen:** `sdk.HostRequest` liefert jeden Status als normalen Rückgabewert; das Plugin prüft `resp.Status` und meldet `auth_failed`/`provider_unavailable` mit Status + Body-Ausschnitt
- **Content-Types:** Das Token-Endpunkt antwortet `text/plain` (JWT nackt im Body) – `text/plain` muss in `permissions.downloads.contentTypes` des Manifests erlaubt sein, sonst blockt die Host-Policy den Request
- **Verifiziert end-to-end:** Plugin-Discovery (`GET /plugins` → `garmin`, `available`), Instanz-Anlage per PocketBase-API, Cron-Sync (`POCKETBASE_CRON_SYNC_SCHEDULE`, Standard nächtlich) mit Import echter Aktivitäten als Trails inkl. GPX, Dedup bei Folge-Syncs

### Datenformat

- **Anfragen:** JSON oder Query-Parameter
- **Antworten:** JSON (Aktivitäten), ZIP (GPX-Dateien)
- **GPX:** Standard GPX 1.1 Format

---

## 🤝 Mitwirken

Du möchtest das Plugin verbessern? Super! Hier sind ein paar Ideen:

### 🐛 Bugs melden
- Öffne ein Issue im [Wanderer Repository](https://github.com/open-wanderer/wanderer)
- Beschreibe das Problem so detailliert wie möglich
- Füge Screenshots oder Logs hinzu

### 💡 Feature-Anfragen
- Du vermisst eine Funktion? Lass es uns wissen!
- Beispiele: Routen-Import, spezielle Aktivitätstypen, etc.

### 🔨 Code beitragen
1. Fork das [Wanderer Repository](https://github.com/open-wanderer/wanderer)
2. Erstelle einen neuen Branch: `git checkout -b feature/garmin-xy`
3. Implementiere deine Änderungen
4. Teste gründlich
5. Erstelle einen Pull Request

### 🧪 Testing
- Teste mit verschiedenen Garmin-Aktivitätstypen
- Teste mit großen Datenmengen
- Teste Edge Cases (leere Aktivitäten, fehlende Daten, etc.)

---

## 📜 Lizenz

Dieses Plugin ist Teil von [Wanderer](https://github.com/open-wanderer/wanderer) und unterliegt der gleichen Lizenz.

---

## 🙏 Danke!

Vielen Dank, dass du das Garmin Plugin für Wanderer verwendest! 🎉

Falls du Fragen, Probleme oder Verbesserungsvorschläge hast, zögere nicht, diese im [Wanderer Repository](https://github.com/open-wanderer/wanderer) zu teilen.

**Happy Tracking!** 🚶‍♂️🚴‍♂️⛰️
