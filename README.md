# Garmin Plugin für Wanderer

**Importiere deine Garmin-Aktivitäten und -Routen in Wanderer!**

Dieses Plugin verbindet [Wanderer](https://wanderer.to) mit der [Garmin Personal API](https://developer.garmin.com/gc-developer-program/personal/personal-api/), um deine Aktivitäten, Routen und zugehörige Daten (GPX-Tracks, Fotos, Metadaten) automatisch zu synchronisieren.

---

## 💡 Wichtig: Docker-Build mit TinyGo

**TinyGo ist NICHT im Standard-Wanderer-Container enthalten!**

Es gibt **zwei Methoden** das Plugin zu bauen – beide erfordern **keine Installation auf deinem Host**:

### 🔹 Methode 1: TinyGo im wanderer-db Container installieren

```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

Dann das Plugin im Container bauen:
```bash
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local"
```

### 🔹 Methode 2: Separater Dev-Container (empfohlen für Entwicklung)

Verwendet einen eigenständigen TinyGo-Container ohne den wanderer-db Container zu modifizieren:

```bash
cd wanderer-garmin-plugin
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build
```

Dann die Dateien nach `data/plugins/garmin/` kopieren (siehe [Entwicklung & Build-Prozess](#-entwicklung--build-prozess)).

**Beide Methoden funktionieren ohne Go-Installation auf deinem Host!** ✅

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

✅ **Aktivitäten importieren** – Alle deine Garmin-Aktivitäten werden als Wanderer-Trails importiert  
✅ **GPX-Tracks** – Vollständige Track-Daten mit Elevation, Zeitstempeln und Koordinaten  
✅ **Fotos** – Aktivitätsfotos werden mit Geotags importiert  
✅ **Metadaten** – Distanz, Dauer, Höhenmeter, Kalorien, Durchschnitts- und Maximalgeschwindigkeit  
✅ **Kategorien** – Automatische Zuordnung zu Wanderer-Kategorien (Laufen, Radfahren, Wandern, etc.)  
✅ **Privatsphäre** – Behält deine Garmin-Privatsphäre-Einstellungen bei  
✅ **Schwierigkeitsgrade** – easy, moderate, difficult  
✅ **Startdatum-Filter** – Importiere nur Aktivitäten ab einem bestimmten Datum  

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

### 0️⃣ TinyGo im Container installieren (einmalig!)

**TinyGo ist NICHT standardmäßig im Wanderer-Container enthalten.**

Führe diesen Befehl **einmalig** aus, um TinyGo und Make im Container zu installieren:

```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

✅ **Fertig!** TinyGo und Make sind jetzt im Container verfügbar.

---

### 🐳 Im Docker Container bauen

Da Wanderer im Docker läuft, ist dies die **einzige benötigte Methode** – du brauchst **nichts auf deinem Host** zu installieren!

#### Einzeiliger Befehl (kopieren & ausführen):
```bash
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local && echo '✅ Plugin installiert!'" && docker compose restart wanderer-db
```

#### Schritt-für-Schritt:
```bash
# 1. In den laufenden Container einsteigen
docker exec -it wanderer-db /bin/bash

# 2. In das Plugin-Verzeichnis wechseln
cd /app/plugins/garmin

# 3. Plugin bauen (TinyGo ist jetzt installiert!)
make build

# 4. Plugin installieren
make install-local

# 5. Container verlassen
exit

# 6. Wanderer neustarten
docker compose restart wanderer-db
```

✅ **Fertig!** Das Plugin ist jetzt installiert.

---

## ⚙️ Konfiguration

### Schritt 1: Plugin aktivieren

1. Öffne Wanderer in deinem Browser
2. Gehe zu **Einstellungen → Plugins**
3. Suche nach dem **Garmin Plugin**
4. Klicke auf **Hinzufügen** oder **Aktivieren**

### Schritt 2: API Key einrichten

1. Klicke auf das **⚙️ Symbol** neben dem Garmin Plugin
2. Gib deinen **Garmin Personal API Token** im **accessToken-Feld** ein:
   - ⚠️ **Wichtig:** Garmin verwendet einen API-Key (kein OAuth-Token)
   - Den Token hast du zuvor unter [https://connect.garmin.com/app/settings/personalApiToken](https://connect.garmin.com/app/settings/personalApiToken) erstellt
   - Der Token sieht so aus: `abc123def456...`
   - **Feldname in Wanderer:** `accessToken` (trotz des Namens hier deinen API-Key eintragen!)
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

Das Plugin synchronisiert automatisch in regelmäßigen Abständen (standardmäßig alle paar Stunden).

Du kannst die Synchronisationshäufigkeit in den **Plugin-Einstellungen** anpassen.

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
| **Routen** | Synchronisiere geplante Routen | ⚠️ Nein (noch experimentell) |
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

## 🔍 API Key Verifizierung & Debugging

Falls du einen "Setup-Fehler" siehst oder unsicher bist, ob dein API Key korrekt gespeichert wurde, kannst du dies mit den mitgelieferten Debugging-Skripten oder direkt in der Datenbank prüfen.

### **📁 Dedizierte Debugging-Skripte**

Dieses Repository enthält mehrere Python-Skripte für das Debugging. **Führe sie vom `wanderer-garmin-plugin/` Verzeichnis aus:**

| Skript | Zweck |
|--------|-------|
| `check_auth_fields.py` | Prüft welche Auth-Felder in der Datenbank gespeichert sind |
| `verify_key.py` | Verifiziert ob der Garmin API Key gespeichert ist |
| `check_and_trigger_import.py` | Kombiniert: Prüft API Key UND löst Import aus |
| `trigger_sync.py` | Löst manuell einen Sync/Import aus |
| `check_plugin_status.py` | Zeigt den Status der Plugin-Instanz an |
| `check_full_instance.py` | Zeigt alle Daten der Plugin-Instanz |
| `check_current_error.py` | Zeigt aktuelle Fehler der Plugin-Instanz |
| `clear_trails.py` | Löscht alle Trails (für Testzwecke) |

**💡 Tipp:** Für die meisten Fälle reicht:
```bash
cd wanderer-garmin-plugin
python3 check_and_trigger_import.py
```

Detaillierte Dokumentation aller Skripte: [SCRIPTS.md](./SCRIPTS.md)

### **Schnellprüfung (einzeilig):**
```bash
python3 -c "import sqlite3,json;conn=sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db');c=conn.cursor();c.execute('SELECT auth FROM plugin_instances WHERE plugin_id=\"garmin\"');r=c.fetchone();print('✅ API Key gespeichert' if r and 'apiKey' in json.loads(r[0]) else '❌ API Key NICHT gefunden')"
```

### **Detaillierte Prüfung:**
```bash
python3 << 'EOF'
import sqlite3, json
conn = sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db')
cursor = conn.cursor()
cursor.execute('SELECT id, plugin_id, auth, status FROM plugin_instances WHERE plugin_id = "garmin"')
instance = cursor.fetchone()
if instance:
    auth = json.loads(instance[2]) if instance[2] else {}
    print("=== Garmin Plugin Instance ===")
    print(f"ID: {instance[0]}")
    print(f"Status: {instance[3]}")
    print(f"Auth: {json.dumps(auth, indent=2)}")
    if 'apiKey' in auth:
        print(f"\n✅ API Key ist gespeichert (Länge: {len(auth['apiKey'])})")
        print(f"   Beginnt mit: {auth['apiKey'][:20]}...")
    elif 'accessToken' in auth:
        print(f"\n✅ Access Token ist gespeichert (Länge: {len(auth['accessToken'])})")
        print(f"   Beginnt mit: {auth['accessToken'][:20]}...")
    else:
        print("\n❌ Kein API Key oder Access Token in der Datenbank!")
else:
    print("❌ Keine Garmin Plugin Instanz in der Datenbank gefunden")
conn.close()
EOF
```

### **Letzten Fehler prüfen:**
```bash
python3 -c "import sqlite3,json;conn=sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db');c=conn.cursor();c.execute('SELECT last_error FROM plugin_instances WHERE plugin_id=\"garmin\"');r=c.fetchone();print(json.dumps(json.loads(r[0]),indent=2) if r[0] else '✅ Keine Fehler')"
```

### **Mögliche Ursachen für "Setup-Fehler":**
1. **❌ API Key nicht gespeichert** → Datenbank prüfen (siehe oben)
2. **❌ Ungültiger API Key** → Neue API Key bei Garmin erstellen
3. **❌ Falsche API Key** → Die bei Garmin gespeicherte Key stimmt nicht mit der in Wanderer überein
4. **❌ Plugin Manifest Mismatch** → Die `secretField` in plugin.json muss mit dem gespeicherten Feld übereinstimmen

**✅ FIXED:** Das plugin.json verwendet nun korrekt `secretField: "apiKey"` statt "accessToken".

**Quick Fix für ältere Installationen:**
Falls du eine ältere Version des Plugins hast, wo plugin.json noch `"accessToken"` erwartet:
- Lade die aktuelle Version herunter oder
- Ändere manuell in plugin.json: `"secretField": "apiKey"`

---

## ❓ Häufige Probleme & Lösungen

### 🔴 **"Authentication Failed" / "Authentifizierung fehlgeschlagen"**

**Ursache:**
- Ungültiger oder abgelaufener API Key
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
4. Teste deinen API Key mit einem Tool wie [Postman](https://www.postman.com/) oder [curl](https://curl.se/):
   ```bash
   curl -H "Authorization: Bearer DEIN_API_KEY" \
     "https://personal.connectapi.garmin.com/activitylist-service/activities/search/activities?start=0&limit=10"
   ```

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

### 🔴 **"make: command not found" oder "tinygo: command not found"**

**Ursache:**
- TinyGo oder Make sind nicht im Container installiert

**Lösung:**
**Installiere TinyGo und Make im Container:**
```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

Dann versuche den Build erneut:
```bash
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local"
```

---

## 📚 Technische Details

### API Endpunkte

Das Plugin nutzt folgende Garmin Personal API Endpunkte:

| Endpunkt | Zweck |
|----------|-------|
| `/activitylist-service/activities/search/activities` | Aktivitäten suchen und auflisten |
| `/activity-service/activity/{id}` | Aktivitätsdetails abrufen |
| `/activity-service/activity/{id}/workouts` | Workouts einer Aktivität |
| `/download-service/files/activity/{id}` | Aktivitätsdatei (ZIP mit GPX) herunterladen |

### Base URL

```
https://personal.connectapi.garmin.com
```

> ⚠️ **Wichtig:** Nicht `api.garmin.com` – das ist die alte API!

### Authentifizierung

- **Methode:** Bearer Token (Personal API Key)
- **Header:** `Authorization: Bearer <API_KEY>`
- **Kein OAuth** nötig – einfacher API Key
- **Konfiguration in Wanderer:** API-Key im `accessToken`-Feld eintragen

### Datenformat

- **Anfragen:** JSON oder Query-Parameter
- **Antworten:** JSON (Aktivitäten), ZIP (GPX-Dateien)
- **GPX:** Standard GPX 1.1 Format

---

## 🛠️ Entwicklung & Build-Prozess

### Docker-basiertes Building (empfohlen)

Da das Plugin in **TinyGo** geschrieben ist und **WASM** als Zielplattform hat, gibt es zwei Build-Methoden:

#### Methode 1: TinyGo im wanderer-db Container installieren

```bash
# TinyGo und Make im Container installieren
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make

# Plugin bauen und installieren
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local"
docker compose restart wanderer-db
```

#### Methode 2: Separater Dev-Container (ohne wanderer-db zu modifizieren)

Diese Methode verwendet einen eigenständigen TinyGo-Container und ist ideal für Entwicklung ohne die Produktionsumgebung zu beeinflussen.

```bash
# In das Plugin-Verzeichnis wechseln
cd /opt/wanderer-dev/wanderer-garmin-plugin

# Plugin mit TinyGo 0.44.0 bauen (enthält Go 1.27.3)
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build

# Plugin installieren
mkdir -p /opt/wanderer-dev/data/plugins/garmin
cp dist/garmin/plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/icon.svg

# Wanderer neustarten
docker compose restart wanderer-db
```

**Vorteile:**
- ✅ Keine Änderungen am wanderer-db Container
- ✅ Saubere Trennung von Build- und Runtime-Umgebung
- ✅ Einfacheres Debugging und Testing
- ✅ Keine Abhängigkeiten auf dem Host

### Makefile-Ziele

| Ziel | Beschreibung |
|------|--------------|
| `make build` | Baut das plugin.wasm |
| `make install` | Baut und kopiert Dateien in dist/garmin/ |
| `make install-local` | Installiert in ../../../data/plugins/garmin/ |
| `make clean` | Entfernt Build-Artefakte |
| `make package` | Erstellt ein tar.gz-Paket für Verteilung |

### Abhängigkeiten

- **TinyGo** ≥ 0.44.0 (empfohlen, enthält Go 1.27.3)
- **Go** ≥ 1.25.0 (durch TinyGo bereitgestellt)
- **Make** (für Build-Skripte)

### SDK-Integration

Das Plugin verwendet das **Wanderer Plugin SDK** aus dem Haupt-Repository. Die SDK-Dateien wurden lokal integriert:

```
wanderer-garmin-plugin/
├── sdk/
│   ├── go.mod
│   ├── types.go
│   ├── host_http.go
│   └── ...
```

Die `go.mod` Datei enthält einen **replace**-Eintrag:
```go
replace github.com/open-wanderer/wanderer/plugins/sdk => ./sdk
```

### Bugfixes und Code-Änderungen

Während des Build-Prozesses wurden folgende Korrekturen vorgenommen:

#### 1. Variable Shadowing in garmin.go (Zeilen 337 & 344)

**Problem:** Die Variable `track` überschattete den Typ `track`:
```go
var track track  // ❌ Variable name shadows type
```

**Lösung:** Umbenennung in `trackData`:
```go
var trackData track  // ✅ Different variable name
```

Alle Referenzen wurden entsprechend aktualisiert.

#### 2. InputJSON in main.go (Zeile 139)

**Problem:** Falsche Annahme über Return-Werte:
```go
data, err := pdk.InputJSON(input)  // ❌ Returns 1 value, not 2
```

**Lösung:** Korrekte Verwendung:
```go
err := pdk.InputJSON(&input)  // ✅ Correct: takes pointer, returns error
```

#### 3. Unused Import in main.go

**Problem:** Import ohne Verwendung:
```go
import "encoding/json"  // ❌ Not used
```

**Lösung:** Import entfernt.

#### 4. TinyGo Version & Flags

**Problem:** Die ursprünglichen Flags `-opt=2 -scheduler=none -gc=none` verursachten WASM-Linker-Fehler:
```
wasm-ld: error: main.lto..o: undefined symbol: runtime.alloc
```

**Lösung:** Aktualisierung auf TinyGo 0.44.0 und Entfernung der problematischen Flags.

---

## 🔧 Wichtige Hinweise für Entwickler

### Auth-Kontext: Nur EIN Kontext pro Plugin!

**Wichtig:** Wanderer erwartet **genau EINEN Auth-Kontext** pro Plugin. Das Garmin-Plugin hatte ursprünglich zwei Kontexte (`api_key` und `bearer_token`), was dazu führte, dass der Token nicht in der WandererDB gespeichert werden konnte.

**Lösung:** Nur `bearer_token` beibehalten, da:
- Der Garmin Personal API den API-Key als Bearer-Token verwendet
- Der Connector bereits `bearer_token` referenziert
- Alle anderen Plugins (Strava, Hammerhead, Komoot) verwenden ebenfalls nur EINEN Kontext

**Vergleich der Plugins:**

| Plugin | Auth-Kontexte | Token-Speicherung |
|--------|---------------|-------------------|
| Strava | 1 (`oauth_access_token`) | ✅ Funktioniert |
| Hammerhead | 1 (`provider_session`) | ✅ Funktioniert |
| Komoot | 1 (`provider_session`) | ✅ Funktioniert |
| Garmin | 1 (`bearer_token`) | ✅ Funktioniert |

### SDK-Typen verwenden

Verwende **SDK-Typen** anstelle von benutzerdefinierten Typen für bessere Konsistenz:

```go
// ✅ Richtig:
type listInput = sdk.ListInput
type detailInput = sdk.DetailInput
type trailSummary = sdk.TrailSummary

// ❌ Veraltet:
type listInput struct {
    Instance sdk.InstanceRef `json:"instance"`
    Auth     map[string]any  `json:"auth,omitempty"`
    // ...
}
```

### Referenz-Implementierung

Das **Strava-Plugin** ist eine gute Referenz für:
- Auth-Handling
- API-Aufrufe
- Fehlerbehandlung
- SDK-Nutzung

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

## 📝 Versionsgeschichte

### v0.1.0 (2024-11-13)
- **Erste stabile Version** des Garmin-Plugins
- Support für Garmin Personal API
- Import von Aktivitäten mit GPX-Tracks, Metadaten und Fotos
- Automatische Kategorie-Zuordnung

**Technische Änderungen:**
- Integration des Wanderer Plugin SDK
- TinyGo 0.44.0 Support (Go 1.27.3)
- WASM-Build ohne Host-Go-Installation
- Code-Fixes:
  - Variable Shadowing in garmin.go behoben
  - InputJSON-Verwendung in main.go korrigiert
  - Unused Import entfernt
  - TinyGo Flags für WASM-Build angepasst

---

## 📜 Lizenz

Dieses Plugin ist Teil von [Wanderer](https://github.com/open-wanderer/wanderer) und unterliegt der gleichen Lizenz.

---

## 🙏 Danke!

Vielen Dank, dass du das Garmin Plugin für Wanderer verwendest! 🎉

Falls du Fragen, Probleme oder Verbesserungsvorschläge hast, zögere nicht, diese im [Wanderer Repository](https://github.com/open-wanderer/wanderer) zu teilen.

**Happy Tracking!** 🚶‍♂️🚴‍♂️⛰️
