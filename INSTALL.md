# ⚡ Schnellstart: Garmin Plugin installieren & nutzen

**In 5 Minuten zum ersten Import!** 🚀

---

## 💡 Wichtig: TinyGo im Container

**TinyGo ist NICHT im Standard-Wanderer-Container enthalten!**

Da du Wanderer im Docker betreibst und **Nichts auf deinem Host installieren** möchtest, musst du **TinyGo einmalig im Container installieren**. Danach kannst du das Plugin bauen.

---

## 📋 Was du brauchst

✅ Wanderer Installation (im Docker)  
✅ Garmin Connect Account  
✅ [Garmin Personal API Token](https://connect.garmin.com/app/settings/personalApiToken)  
✅ Docker (hast du schon)  

---

## 🎯 Schritt-für-Schritt Anleitung

---

## 0️⃣ TinyGo im Container installieren (einmalig!)

**Führe diesen Befehl nur EINMAL aus:**

```bash
docker exec -it wanderer-db apt-get update && apt-get install -y tinygo make
```

✅ **Fertig!** TinyGo und Make sind jetzt im Container verfügbar.

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
   - **API Key:** Füge deinen Garmin Token ein (aus Schritt 1)
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
