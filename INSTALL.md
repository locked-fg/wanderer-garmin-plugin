# ⚡ Schnellstart: Garmin Plugin installieren & nutzen

**In 5 Minuten zum ersten Import!**

---

## 📋 Was du brauchst

✅ Laufende Wanderer-Installation (Docker)
✅ Garmin Connect Account
✅ [Garmin Personal API Key](https://connect.garmin.com/app/settings/personalApiToken) (`GCDC_...`)
✅ Das fertige Plugin-Paket (`dist/garmin/` in diesem Repository – `plugin.json`, `plugin.wasm`, `icon.svg`)

> **Kein TinyGo nötig:** Das Plugin ist bereits gebaut (`dist/garmin/plugin.wasm`, mit Go 1.25.14 + TinyGo 0.41.1 – derselben Version wie der Wanderer-CI-Build). Du musst nur die 3 Dateien kopieren.

---

## 1️⃣ Garmin API Key erstellen (2 Minuten)

1. **Öffne** [Garmin Connect Personal API Tokens](https://connect.garmin.com/app/settings/personalApiToken)
2. **Logge dich ein** mit deinem Garmin-Konto
3. **Klicke auf** „Generate Token"
4. **Kopiere** den Key (wird nur einmal angezeigt – sicher speichern!)

> ⚠️ **Wichtig (live verifiziert):** Der `GCDC_...`-Wert ist ein **API Key, kein Bearer-Token**. Direkt als `Authorization: Bearer GCDC_...` gesendet, antwortet Garmin mit `401 invalid_token, "Token is not active"`. Das Plugin tauscht den Key daher automatisch gegen ein **24h gültiges JWT** aus (`POST /personalapitoken-service/personal-api/personal-api-key/token` mit `Authorization: ApiKey <key>`) und verwendet danach `Authorization: Bearer <jwt>`. Du musst nichts weiter konfigurieren.

---

## 2️⃣ Plugin installieren (1 Minute)

Das Plugin landet im Plugin-Verzeichnis von Wanderer: `data/plugins/garmin/` (im Container: `/data/plugins/garmin/`).

**Variante A – Dateien direkt auf dem Host kopieren** (empfohlen, falls `data/` ein Bind-Mount ist):

```bash
# Im Repository-Verzeichnis des Plugins:
mkdir -p /pfad/zu/deinem/wanderer/data/plugins/garmin
cp dist/garmin/plugin.json dist/garmin/plugin.wasm dist/garmin/icon.svg /pfad/zu/deinem/wanderer/data/plugins/garmin/

docker compose restart wanderer-db
```

**Variante B – In den laufenden Container kopieren:**

```bash
docker exec wanderer-db mkdir -p /data/plugins/garmin
docker cp dist/garmin/plugin.json wanderer-db:/data/plugins/garmin/plugin.json
docker cp dist/garmin/plugin.wasm wanderer-db:/data/plugins/garmin/plugin.wasm
docker cp dist/garmin/icon.svg  wanderer-db:/data/plugins/garmin/icon.svg

docker compose restart wanderer-db
```

✅ **Fertig!** Nach dem Neustart erscheint das Plugin in Wanderer unter **Einstellungen → Plugins**.

---

## 3️⃣ Plugin in Wanderer konfigurieren (1 Minute)

1. **Öffne Wanderer** im Browser (z.B. `http://localhost:8080`)
2. **Gehe zu:** Einstellungen → Plugins
3. **Suche nach:** „Garmin" → **Hinzufügen**
4. **Konfiguriere:**
   - **API Key** (Feld `apiKey`): Deinen Garmin Personal API Key aus Schritt 1 (`GCDC_...`) – das Plugin tauscht ihn automatisch gegen ein 24h-Access-Token
   - **Startdatum:** (optional) z.B. `2024-01-01` – importiert nur Aktivitäten ab diesem Datum (Filter über `startDate`-Query-Param + clientseitig über `beginTimestamp`)
5. **Speichern** und **Aktivieren**

---

## 4️⃣ Erste Synchronisation starten

1. Klicke auf **„Jetzt synchronisieren"** (oder warte auf die automatische Sync)
2. Warte 1–5 Minuten (je nach Anzahl der Aktivitäten)
3. Gehe zu **„Trails"** – deine Garmin-Aktivitäten sollten jetzt angezeigt werden

---

## 🎯 Was wird importiert?

Jede Aktivität wird als **Trail** importiert mit:

✅ **GPX-Track** – aus `geoPolylineDTO.polyline` (`/activity-service/activity/{id}/details`) rekonstruiert: Koordinaten, Höhe und Zeitstempel pro Punkt
✅ **Metadaten** – Distanz, Dauer (moving/elapsed), Höhenmeter, Kalorien, Geschwindigkeiten, HF, Ort (`summaryDTO`)
✅ **Kategorie** – automatische Zuordnung aus `activityTypeDTO.typeKey` (Laufen, MTB, Wandern, etc.)
✅ **Privatsphäre** – aus `accessControlRuleDTO.typeKey` übernommen
❌ **Fotos** – werden derzeit nicht importiert (kein Endpunkt verifiziert)

---

## 🔨 Alternative: Aus dem Quellcode bauen

Nur nötig, wenn du das Plugin selbst kompilieren möchtest (verifizierte Rezeptur):

- **Go 1.25.14 + TinyGo 0.41.1** auf dem Build-Host (nicht im Container)
- Das Plugin muss im Wanderer-Repository unter `plugins/garmin` liegen (wegen `replace ../sdk` in `go.mod`) – ein standalone Checkout baut nicht
- Build: `tinygo build -target=wasi -scheduler=none -no-debug` (CI-kompatibel), Manifest via `go run .../sdk/cmd/manifestcheck` (validiert gegen `plugins/schema/plugin.schema.json`)

```bash
git clone https://github.com/open-wanderer/wanderer.git
cp -r wanderer-garmin-plugin wanderer/plugins/garmin
cd wanderer/plugins/garmin
make build           # -> dist/garmin/{plugin.json,plugin.wasm,icon.svg}
make install-local   # -> ../../data/plugins/garmin/
```

Danach wie in Schritt 2 beschrieben installieren und Wanderer neu starten.

---

## 🔧 Tipps & Tricks

### 📅 Startdatum nutzen
- Bei 1000+ Aktivitäten: Startdatum setzen (z.B. `2024-01-01`) und schrittweise importieren (erst 2024, dann 2023, …)

### 🔄 Automatische Synchronisation
- Das Plugin synchronisiert automatisch in regelmäßigen Abständen
- Häufigkeit in den Plugin-Einstellungen anpassbar
- Hinweis: Pro Sync wird der API Key gegen ein neues 24h-JWT getauscht

### 🔍 Manuelle Synchronisation
- **Einstellungen → Plugins → Garmin → Jetzt synchronisieren**

### 🗑️ Plugin deaktivieren
- **Einstellungen → Plugins → Garmin → Deaktivieren**
- Bereits importierte Daten bleiben erhalten

---

## ❌ Häufige Probleme & schnelle Lösungen

| Problem | Lösung |
|---------|--------|
| **„Authentifizierung fehlgeschlagen" / `401 Token is not active`** | Der `GCDC_...`-Key ist revoked/abgelaufen – neuen Key unter [connect.garmin.com/app/settings/personalApiToken](https://connect.garmin.com/app/settings/personalApiToken) erstellen (ohne Leerzeichen kopieren). Das Plugin macht den Token-Austausch automatisch; schlägt der aus, ist der Key das Problem |
| **„Keine Aktivitäten gefunden"** | Startdatum prüfen (muss in der Vergangenheit liegen); sicherstellen, dass Aktivitäten im Garmin-Konto existieren |
| **„Plugin nicht sichtbar"** | Prüfen, ob die 3 Dateien unter `data/plugins/garmin/` liegen (`plugin.json`, `plugin.wasm`, `icon.svg`) und Wanderer **neu starten** |
| **„Rate Limited"** | 1 Stunde warten; Sync-Häufigkeit reduzieren; Startdatum setzen |
| **„GPX fehlt"** | Aktivität hat keinen Track (z.B. Krafttraining) – nur Metadaten werden importiert |
| **Build-Fehler `../sdk` nicht gefunden** | Plugin muss unter `plugins/garmin` **innerhalb** des wanderer-Repos liegen, nicht standalone |

---

## 🔬 Technische Details (live verifiziert gegen `https://personal.connectapi.garmin.com`)

Alle Angaben wurden mit einem echten Personal API Token gegen die Live-API getestet (der Test-Key wurde danach **revoked**):

| Endpunkt | Verhalten |
|----------|-----------|
| `POST /personalapitoken-service/personal-api/personal-api-key/token` | `Authorization: ApiKey GCDC_...`, JSON-Body `{}` → 200, Body ist das **JWT (24h)** |
| `GET /activitylist-service/activities/search/activities?start=0&limit=5&startDate=2026-10-01` | 200, **nacktes JSON-Array**; Felder: `activityId`, `activityName`, `startTimeGMT` (`"2006-01-02 15:04:05"`), `beginTimestamp` (ms), `distance`, `duration`, `elevationGain/Loss`, `averageSpeed`, `maxSpeed`, `calories`, `startLatitude/Longitude`, `locationName`, `hasImages` |
| `GET /activity-service/activity/{id}` | 200; Metriken unter `summaryDTO` (`startTimeGMT` hier als `"2026-10-12T16:09:45.0"`), Typ unter `activityTypeDTO.typeKey`, Privatsphäre unter `accessControlRuleDTO.typeKey`, `metadataDTO.fileFormat.formatKey` |
| `GET /activity-service/activity/{id}/details` | 200; Track unter `geoPolylineDTO.polyline` = `[{lat, lon, altitude, time(ms)}, …]` → GPX |
| `GET /download-service/files/activity/{id}` | 200, aber **ZIP mit FIT-Datei** (`*_ACTIVITY.fit`), **kein GPX** – wird nicht verwendet |

**Gotchas:** `activityType`/`privacy` sind Objekte `{typeId, typeKey}`, keine Strings; Zeitstempel kommen in zwei Formaten (`"2006-01-02 15:04:05"` in der Liste, `"2026-10-12T16:09:45.0"` in `summaryDTO`) und werden nach RFC3339 UTC normalisiert; `startDate` filtert tagesgenau, zusätzlich clientseitig über `beginTimestamp`.

**Live gegen Docker-Wanderer verifiziert (Plugin-Load + Sync + Fehlerpfad):**

- Das Token-Endpunkt antwortet mit `Content-Type: text/plain` (der JWT steht nackt im Body) – `text/plain` muss deshalb in `permissions.downloads.contentTypes` des Manifests erlaubt sein, sonst blockt die Host-Policy den Request (`content type "text/plain" is not allowed by manifest permissions`).
- Fehler müssen über den **Extism-Fehlerkanal** gemeldet werden (`pdk.SetErrorString` mit JSON `{code, message}` bei Rückgabewert ≠ 0), nicht über `pdk.OutputJSON` – sonst sieht der Host nur den generischen `plugin_error: call <export> failed with code 1`.
- `sdk.HostRequest` liefert **jeden** HTTP-Status als normalen Rückgabewert; der Plugin-Code muss `resp.Status` selbst prüfen (revoked/deleted Key → `404 {"error":"NotFoundException","message":"Personal API Key with ID … does not exist."}` am Token-Endpunkt, als `auth_failed` im Instanz-`last_error` sichtbar).

**Build-Artefakt:** `dist/garmin/plugin.wasm` (~295 KB) mit den Exports `list_routes_v1`, `list_activities_v1`, `get_route_detail_v1`, `get_activity_detail_v1`; `plugin.json` valide gegen `plugins/schema/plugin.schema.json` und via `manifestcheck` geprüft.

**Bekannte Einschränkungen:** Fotos und Routen werden nicht importiert (kein Endpunkt in der Personal API verifiziert); der API Key wird pro Sync gegen ein neues 24h-JWT getauscht.

---

## 📞 Hilfe & Support

- **Dokumentation:** [README.md](./README.md)
- **Wanderer Community:** [Discussions](https://github.com/open-wanderer/wanderer/discussions)
- **Bug melden:** [Issues](https://github.com/open-wanderer/wanderer/issues)
