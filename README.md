# Garmin Plugin für Wanderer

**Importiere deine Garmin-Aktivitäten und -Routen in Wanderer!**

Dieses Plugin verbindet [Wanderer](https://wanderer.to) mit der [Garmin Personal API](https://developer.garmin.com/gc-developer-program/personal/personal-api/), um deine Aktivitäten, Routen und zugehörige Daten (GPX-Tracks, Fotos, Metadaten) automatisch zu synchronisieren.

---

## 💡 Wichtig: Docker-Build mit TinyGo

**TinyGo ist NICHT im Standard-Wanderer-Container enthalten!**

Da du Wanderer im Docker betreibst und **Nichts auf deinem Host installieren** möchtest, musst du **TinyGo einmalig im Container installieren**:

```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

Danach kannst du das Plugin im Container bauen – **ohne etwas auf deinem Host zu installieren!**

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

## 📜 Lizenz

Dieses Plugin ist Teil von [Wanderer](https://github.com/open-wanderer/wanderer) und unterliegt der gleichen Lizenz.

---

## 🙏 Danke!

Vielen Dank, dass du das Garmin Plugin für Wanderer verwendest! 🎉

Falls du Fragen, Probleme oder Verbesserungsvorschläge hast, zögere nicht, diese im [Wanderer Repository](https://github.com/open-wanderer/wanderer) zu teilen.

**Happy Tracking!** 🚶‍♂️🚴‍♂️⛰️
