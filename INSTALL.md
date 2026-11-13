# ⚡ Schnellstart: Garmin Plugin installieren & nutzen

**In 5 Minuten zum ersten Import!** 🚀

---

## 💡 Wichtig: Build-Optionen

Es gibt **zwei Methoden** das Plugin zu bauen:

### Methode 1: TinyGo im wanderer-db Container installieren (empfohlen für Produktion)
TinyGo direkt im laufenden Wanderer-Container installieren.

### Methode 2: Dev-Container mit TinyGo verwenden (empfohlen für Entwicklung)
Plugin in einem separaten Docker-Container bauen, ohne den Wanderer-Container zu modifizieren.

**Für beide Methoden: Nichts wird auf deinem Host installiert!** ✅

---

## 📋 Was du brauchst

✅ Wanderer Installation (im Docker)  
✅ Garmin Connect Account  
✅ [Garmin Personal API Token](https://connect.garmin.com/app/settings/personalApiToken)  
✅ Docker (hast du schon)  

---

## 🎯 Schritt-für-Schritt Anleitung

---

## 0️⃣ Build-Methode wählen

### Option A: TinyGo im wanderer-db Container installieren (einmalig!)

**Führe diesen Befehl nur EINMAL aus:**

```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

✅ **Fertig!** TinyGo und Make sind jetzt im Container verfügbar.

**Dann weiter mit [Schritt 2: Plugin installieren](#2️⃣-plugin-installieren-2-minuten)**

---

### Option B: Dev-Container mit TinyGo (ohne wanderer-db zu modifizieren)

Diese Methode baut das Plugin in einem separaten Container und kopiert die Dateien dann in deinen Daten-Ordner.

**Vorteile:**
- Keine Änderungen am wanderer-db Container
- Saubere Trennung von Build- und Runtime-Umgebung
- Einfacher zu pflegen und zu aktualisieren

**Build-Befehl:**

```bash
cd /opt/wanderer-dev/wanderer-garmin-plugin
docker run --rm -v $(pwd):/src -w /src tinygo/tinygo:0.44.0 make build
```

**Installieren:**

```bash
# Kopiere die gebauten Dateien in den Daten-Ordner
mkdir -p /opt/wanderer-dev/data/plugins/garmin
cp dist/garmin/plugin.json /opt/wanderer-dev/data/plugins/garmin/
cp dist/garmin/plugin.wasm /opt/wanderer-dev/data/plugins/garmin/
cp assets/icon.svg /opt/wanderer-dev/data/plugins/garmin/icon.svg

# Wanderer neustarten
docker compose restart wanderer-db
```

✅ **Fertig!** Das Plugin ist jetzt installiert.

**Dann weiter mit [Schritt 3: Plugin in Wanderer konfigurieren](#3️⃣-plugin-in-wanderer-konfigurieren-1-minute)**

---

## 1️⃣ Garmin API Key erstellen (2 Minuten)

1. **Öffne** [Garmin Connect Personal API Tokens](https://connect.garmin.com/app/settings/personalApiToken) in deinem Browser
2. **Logge dich ein** mit deinem Garmin-Konto
3. **Klicke auf** "Generate Token" (Token generieren)
4. **Kopiere** den generierten Token (z.B. `abc123def456...`)
   - ⚠️ **Wichtig:** Der Token wird nur **einmal** angezeigt – speichere ihn sicher!

✅ **Fertig!** Du hast jetzt deinen Garmin Personal API Key.

---

## 2️⃣ Plugin installieren (2 Minuten)

### 🐳 **Im Docker Container bauen**

```bash
# Alles in einem Befehl - kopieren und ausführen!
docker exec -it wanderer-db bash -c "cd /app/plugins/garmin && make build && make install-local && echo '✅ Plugin installiert!'" && docker compose restart wanderer-db
```

**Schritt-für-Schritt:**
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

## 3️⃣ Plugin in Wanderer konfigurieren (1 Minute)

1. **Öffne Wanderer** in deinem Browser: `http://localhost:8080` (oder deine Domain)
2. **Gehe zu:** Einstellungen → Plugins
3. **Suche nach:** "Garmin"
4. **Klicke auf:** ➕ Hinzufügen
5. **Konfiguriere:**
   - **accessToken:** Füge deinen Garmin API-Key ein (aus Schritt 1)
     - ⚠️ **Wichtig:** Feld heißt `accessToken`, aber du trägst hier deinen **API-Key** ein (nicht OAuth!)
   - **Startdatum:** (Optional) z.B. `2024-01-01` – importiert nur Aktivitäten ab diesem Datum
   - **Aktivitäten:** ✅ Aktiviert
   - **Routen:** ⚠️ Optional (noch experimentell)
6. **Speichern** und **Aktivieren**

✅ **Fertig!** Das Plugin ist jetzt aktiv.

---

## 4️⃣ Erste Synchronisation starten

1. **Klicke auf** "Jetzt synchronisieren" (oder warte auf die automatische Sync)
2. **Warte** – die Synchronisation kann je nach Anzahl der Aktivitäten 1-5 Minuten dauern
3. **Gehe zu** "Trails" – deine Garmin-Aktivitäten sollten jetzt angezeigt werden! 🎉

---

## 🎯 Was wird importiert?

Jede Aktivität wird als **Trail** in Wanderer importiert mit:

✅ **GPX-Track** – Komplette Route mit Koordinaten und Höhenprofil  
✅ **Metadaten** – Distanz, Dauer, Höhenmeter, Kalorien, Geschwindigkeit  
✅ **Fotos** – Aktivitätsfotos mit Geotags (falls vorhanden)  
✅ **Kategorie** – Automatische Zuordnung (Laufen, Radfahren, Wandern, etc.)  
✅ **Privatsphäre** – Behält deine Garmin-Einstellungen bei  

---

## 🔧 Tipps & Tricks

### 📅 Startdatum nutzen
- **Problem:** Du hast 1000+ Aktivitäten und die Sync dauert ewig
- **Lösung:** Setze ein Startdatum (z.B. `2024-01-01`) und importiere schrittweise:
  - 1. Sync: `2024-01-01`
  - 2. Sync: `2023-01-01`
  - 3. Sync: `2022-01-01`

### 🔄 Automatische Synchronisation
- Das Plugin synchronisiert **automatisch** alle paar Stunden
- Du kannst die Häufigkeit in den Plugin-Einstellungen anpassen

### 🔍 Manuelle Synchronisation
- Gehe zu **Einstellungen → Plugins → Garmin → Jetzt synchronisieren**

### 🗑️ Plugin deaktivieren
- Gehe zu **Einstellungen → Plugins → Garmin → Deaktivieren**
- Deine importierten Daten bleiben erhalten!

---

## ❌ Häufige Probleme & schnelle Lösungen

| Problem | Lösung |
|---------|--------|
| **"Authentifizierung fehlgeschlagen"** | API Key prüfen, neuen Token erstellen |
| **"Keine Aktivitäten gefunden"** | Startdatum prüfen (muss in der Vergangenheit liegen) |
| **"Plugin nicht sichtbar"** | Wanderer neustarten, Plugin-Dateien prüfen |
| **"Rate Limited"** | 1 Stunde warten, dann erneut versuchen |
| **"GPX fehlt"** | Aktivität hat keinen Track (z.B. Krafttraining) – nur Metadaten werden importiert |
| **"make: command not found"** | TinyGo und Make im Container installieren: `docker exec -it wanderer-db apt-get install -y tinygo make` |
| **"tinygo: command not found"** | TinyGo im Container installieren: `docker exec -it wanderer-db apt-get install -y tinygo` |
| **"Setup-Fehler" / API Key prüfen** | Datenbank prüfen: `python3 -c "import sqlite3,json;conn=sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db');c=conn.cursor();c.execute('SELECT auth FROM plugin_instances WHERE plugin_id=\"garmin\"');print(json.loads(c.fetchone()[0]) if c.fetchone() else 'Kein Eintrag')"` |
| **Plugin Status prüfen** | Datenbank prüfen: `python3 -c "import sqlite3;conn=sqlite3.connect('/opt/wanderer-dev/data/pb_data/data.db');c=conn.cursor();c.execute('SELECT id,status,enabled FROM plugin_instances WHERE plugin_id=\"garmin\"');print(c.fetchone())"` |

---

## 📚 Technische Notizen für Entwickler

### Code-Fixes während des Builds

Während der Entwicklung wurden folgende Korrekturen am Plugin-Code vorgenommen:

#### 1. Variable Shadowing in `garmin.go`
- **Problem:** Variable `track` überschattete Typ `track` (Zeilen 337 & 344)
- **Fix:** Umbenennung in `trackData`

#### 2. Falsche InputJSON-Verwendung in `main.go`
- **Problem:** `pdk.InputJSON(input)` gibt nur 1 Wert zurück, nicht 2
- **Fix:** Korrekte Verwendung mit Pointer: `pdk.InputJSON(&input)`

#### 3. Unused Import in `main.go`
- **Problem:** `encoding/json` wurde importiert aber nicht verwendet
- **Fix:** Import entfernt

#### 4. TinyGo Version & Flags
- **Problem:** Flags `-opt=2 -scheduler=none -gc=none` verursachten WASM-Linker-Fehler
- **Fix:** Aktualisierung auf TinyGo 0.44.0, problematische Flags entfernt

### SDK-Integration

Das Plugin verwendet das **Wanderer Plugin SDK**. Die SDK-Dateien wurden lokal unter `wanderer-garmin-plugin/sdk/` integriert und in `go.mod` per replace-Direktive eingebunden:

```go
replace github.com/open-wanderer/wanderer/plugins/sdk => ./sdk
```

Dies ermöglicht die Entwicklung ohne Abhängigkeit vom Haupt-Repository.

### Makefile-Anpassungen

- **TINYGO_VERSION:** Aktualisiert von 0.32.0 auf 0.44.0
- **TINYGO_FLAGS:** Leer gelassen (Flag-Entfernung nötig für erfolgreichen Build)

---

## 📞 Hilfe & Support

- **Dokumentation:** [README.md](./README.md)
- **Wanderer Community:** [Discussions](https://github.com/open-wanderer/wanderer/discussions)
- **Bug melden:** [Issues](https://github.com/open-wanderer/wanderer/issues)

---

## 🎉 Erfolg! 

Du hast das Garmin Plugin erfolgreich installiert und konfiguriert! 🎊

Jetzt kannst du alle deine Garmin-Aktivitäten in Wanderer verwalten, analysieren und mit anderen teilen.

**Viel Spaß beim Erkunden!** 🗺️⛰️🚴‍♂️
