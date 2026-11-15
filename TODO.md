# Garmin Plugin – Arbeitsstand

Stand: 2026-11-15

## Kurzstatus

Das Plugin baut und installiert erfolgreich. Die Garmin-Activities werden importiert. Ursache des bisherigen Fehlers war ein ungültiges `startedAt`-Format; die Zeit wird jetzt als RFC3339 ausgegeben. Der letzte Lauf legte 57 Garmin-Referenzen an; diese zeigen auf 48 Trail-Datensätze.

## Was bereits getan wurde

- Garmin-Antwortformat direkt geprüft und in `README.md` dokumentiert:
  - Token-Endpunkt: `200 text/plain`.
  - Aktivitätssuche: `200`, oberstes JSON-Element ist ein Array; `activityType` und `privacy` sind Objekte; `calories` ist numerisch; Zeitfelder sind Strings.
  - `/activity-service/activity/{id}/details` liefert ein JSON-Objekt mit `geoPolylineDTO`, `activityDetailMetrics` und `metricDescriptors`.
  - Der getestete Download `/download-service/files/activity/{id}` ist ein ZIP mit einer `.fit`-Datei, nicht mit GPX.
- Offiziellen Garmin-Client nach `garmin-personal-api-python/` geklont und geprüft. Er hat Wrapper für Aktivitätsübersicht, Workouts und ZIP-Download; für `/details` gibt es keinen eigenen Wrapper, aber der generische `get()`-Aufruf kann ihn anfragen.
- Parser an die Listenantwort angepasst (Array, Objekttypen, Dezimalwerte und Garmin-Zeitformate).
- Sessionweiten Token-Cache ergänzt, damit der API-Key nicht für jedes Detail neu getauscht wird.
- Import auf `geoPolylineDTO` und als Rückfall auf `activityDetailMetrics` umgestellt. Garmin lieferte bei einem geprüften Detail `3266` Polyline-Punkte und `249` Messwertzeilen.
- Ausführliche Handoff-Diagnose ergänzt: Summaries werden einzeln als JSON protokolliert; Detail-Logs enthalten Summary-Identität, Import-Metadaten, Trackpunkt-/GPX-Größen und einen GPX-Anfang. Zugangsdaten werden nicht protokolliert.
- TinyGo-Build (`0.44.0`) war erfolgreich; das Plugin wurde lokal installiert.
- Ein korrigierter Sync-Lauf verifizierte den Import: 57 Garmin-Aktivitätsreferenzen zeigen auf 48 Trails; Plugin-Status `configured`, `last_error={}`.
- Die zuvor ausführlichen Summary- und Import-Payload-Logs wurden reduziert. Ein erfolgreicher Aktivitätsimport protokolliert jetzt nur ID, normalisierten Startzeitpunkt, Trackpunktzahl und GPX-Größe; Seitenstatistiken und Fehlerdiagnosen bleiben erhalten. Summary-Inhalte, Metadaten und GPX-Vorschau werden nicht mehr im normalen Info-Log ausgegeben.

## Verlauf und Diagnose

### Früherer lokaler Diagnose-Lauf (2026-11-15, 07:40 UTC)

- Einmaliger Sync um 07:40 UTC; Cron danach wieder auf täglich 02:00 UTC gesetzt und DB-Service mit der normalen Konfiguration neu gestartet.
- Garmin-Auth und API-Aufruf erfolgreich; es kamen 50 und danach 7 Aktivitäten, alle eligible.
- Beim ersten Detail: 3266 `geoPolylineDTO`-Punkte; die erzeugte Plugin-Payload enthielt ebenfalls 3266 Punkte (GPX Base64: 265920 Bytes).
- Sync-Status wurde als `configured` abgeschlossen, `last_error` war leer (`{}`). Datenbank blieb bei 0 Trails und 0 `trail_external_reference`-Einträgen.
- Nach dem Diagnose-Lauf wurde der DB-Container neu erstellt; seine aktuellen Logs enthalten nur den Start um 07:43 und können den Sync nicht mehr belegen.
- Der damalige Laufzeitvergleich zeigte, dass `flomp/wanderer-db:v0.21.2` älter als der lokale Host-Quellstand war. Deshalb fehlten Host-Diagnosemarker. Ein Host-Image-Neubau war dafür jedoch nicht nötig: ausführliche Plugin-Logs machten den fehlerhaften `startedAt`-Wert sichtbar.
- `syncActivities` erzeugt Summaries mit Provider `garmin` und `externalId`; Garmin-Trackpunkte kamen bis zur Plugin-Payload. Der Fehler lag im `startedAt`-Feld der Detailantwort und wurde im Plugin behoben.

### Diagnose und Korrektur (2026-11-15, 08:21–08:28 UTC)

- Auf Wunsch wurden ausführliche Plugin-Logs ergänzt: jede Summary als JSON, die vom Host an die Detailfunktion übergebene Summary sowie das Import-Item mit Metadaten, Trackpunkt-/GPX-Größen und GPX-Anfang.
- Der erste Lauf bestätigte die Übergabe von 57 Summaries und GPX-Detailantworten, importierte aber 0 Trails. Die Detailantwort enthielt z. B. `startedAt: "2026-07-20 16:07:41T00:00:00Z"` – ungültiges RFC3339.
- Ursache: `startTimeGMT` wurde bei einem Leerzeichen-Zeitstempel fälschlich wie ein Datum ohne Uhrzeit behandelt und `T00:00:00Z` angehängt. Der Host erwartet für `startedAt` einen decodierbaren RFC3339-Zeitwert.
- `mapActivityToTrail` normalisiert `startTimeGMT` jetzt mit `parseActivityStartTime` und formatiert es in UTC als RFC3339; bei nicht parsebarem Zeitwert wird `startedAt` ausgelassen.
- Erneut gebaut/installiert und einen einzelnen Lauf ausgeführt. Logbeispiel: `startedAt: "2026-10-12T16:09:45Z"`, 3266 Trackpunkte, GPX Base64 265920 Bytes.
- Ergebnis: Garmin-Status `configured`, `last_error={}`, 57 unterschiedliche Garmin-Aktivitätsreferenzen und 48 Trails. Mehrere Garmin-Referenzen können auf denselben Trail zeigen; der Host kann Aktivitäten zusammenführen. Wanderer-Images wurden nicht verändert. Cron wieder auf täglich 02:00 UTC zurückgesetzt.
- Abgleich der 57 Referenzen abgeschlossen: `merge.enabled=true`; zehn Garmin-Referenzen zeigen auf denselben Trail (`Lenggries Express`), die übrigen 47 jeweils auf einen eigenen Trail. Das ergibt 48 Trails. Der Host-Merge-Code hängt externe Referenzen an den behaltenen Trail um und löscht den Quell-Trail; damit ist die Differenz 57→48 durch das aktivierte Auto-Merge erklärt.

### A. Lokalen Parser vorbereiten — erledigt

- Der Parser liegt unabhängig vom HTTP-Client unter `internal/garmintrack/`. Er bevorzugt `geoPolylineDTO.polyline` und fällt bei leerer Polyline auf Koordinaten in `activityDetailMetrics` zurück.
- Synthetische Fixtures liegen in `testdata/garmin_activity_details_polyline.json` und `testdata/garmin_activity_details_metrics_fallback.json`. Sie enthalten keine echten Garmin-Aktivitätsdaten oder Zugangsdaten.
- Parser-Tests decken Polyline, `valid:false`, Metrics-Fallback, `null`-Koordinaten und ungültiges JSON ab.
- Lokal gegen die Fixtures testen, ohne Garmin aufzurufen:
  `docker run --rm -v "$PWD:/workspace" -w /workspace/wanderer-garmin-plugin tinygo/tinygo:0.44.0 go test ./internal/garmintrack`
- TinyGo-Build und Installation mit dem neuen Paket wurden erfolgreich geprüft (`make clean install-local`).

### B. Build-Container wiederverwenden

1. Beim nächsten Build prüfen, ob das Image `tinygo/tinygo:0.44.0` lokal vorhanden ist (`docker image inspect`). Docker lädt ein vorhandenes Image nicht bei jedem `docker run` erneut herunter.
2. Verstehen: `--rm` entfernt den beendeten Container, nicht das Image. Ohne `--rm` bleibt ein gestoppter Container zurück; das spart Container-Neuanlage, aber ist für die Image-Downloads nicht nötig.
3. Für den einfachen Ablauf `--rm` beibehalten und das TinyGo-Image lokal behalten. Falls Container-Neuanlage tatsächlich stört, einmalig einen benannten Build-Container anlegen und für weitere Builds mit `docker start -ai <name>` wiederverwenden; Arbeitsverzeichnis und Bind-Mount müssen dabei auf `/opt/wanderer-dev/wanderer-garmin-plugin` zeigen.
4. Keine regelmäßigen `docker pull`-Aufrufe in den Build einbauen. Das Image nur bei Bedarf aktualisieren; danach sicherstellen, dass TinyGo-Version und Build-Ergebnis weiterhin stimmen.

## Nächste Schritte

1. Import ist in der GUI bestätigt; der Abgleich 57 Referenzen zu 48 Trails ist erledigt.
2. Erledigt: Diagnose-Logs reduziert und Plugin neu bauen/installieren. Bei künftigen Problemen zuerst die kompakte Aktivitätszeile und Fehlerlogs prüfen; eine GPX-Vorschau bei Bedarf gezielt temporär ergänzen.
3. Falls jede Garmin-Aktivität als eigener Trail gewünscht ist, `merge.enabled` für die Garmin-Instanz deaktivieren. Aktuell bleibt Auto-Merge aktiv und führt ähnliche Aktivitäten zusammen.

## Wiederholung des Garmin-Sync-Diagnoselaufs (nur falls neue API-Daten erforderlich sind)

1. In `docker-compose.yaml` steht wieder der tägliche Zeitplan `POCKETBASE_CRON_SYNC_SCHEDULE: 0 2 * * *`.
2. Für einen Diagnose-Lauf den Wert **einmalig** auf einen konkreten zukünftigen Zeitpunkt setzen, z. B. `"M H * * *"` (Minute/Stunde passend zur aktuellen UTC-Zeit), dann aus `/opt/wanderer-dev` `docker compose up -d --no-deps db` ausführen. **Keinen Minutentakt verwenden**: Ein Lauf braucht ungefähr 50 Sekunden; Minutentakt hat Syncs überlappen lassen und Garmin mit `429` ins Rate-Limit gebracht.
3. Den einzelnen Lauf abschließen lassen und in den Logs nach `Garmin: list capability output`, `Garmin: detail request received`, `Garmin: detail capability output payload` sowie Host-Importfehlern suchen.
4. Danach Cron sofort auf `0 2 * * *` zurückstellen und erneut `docker compose up -d --no-deps db` ausführen.
5. Datenbank prüfen: `trails` und `trail_external_reference` müssen Garmin-Einträge enthalten. Bei einem unerwarteten Ergebnis zuerst Summary, `startedAt` und Import-Payload in den Plugin-Logs prüfen; ein Host-Image-Neubau ist dafür nicht erforderlich.

## Relevante Dateien

- `garmin.go`, `types.go`, `plugin.json`: Parser, Detail-Trackpunkte, Manifest und ausführliche Handoff-Diagnose.
- `README.md`: dokumentiertes API-Antwortformat und Zeitformat-Hinweis.
- `garmin-personal-api-python/`: lokaler Clone des offiziellen Clients.
- `docker-compose.yaml`: aktueller Cron ist täglich um 02:00.

## Zugangsdaten

Ein API-Key wurde zuvor in einer Notiz abgelegt und in der Unterhaltung geteilt. Er wurde aus dieser TODO-Datei entfernt. Den betroffenen Key in Garmin Connect widerrufen und einen neuen verwenden. Zugangsdaten nicht in Logs, TODO-Dateien oder Commits speichern.
