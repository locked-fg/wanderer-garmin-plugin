# Garmin Plugin - TODO List

## 🎯 **Current Status: Partially Working**

### ✅ **What's Working**
- Plugin loads without "Setup-Fehler" in Wanderer UI
- Plugin icon displays correctly (GA text)
- API key can be saved in database (encrypted by PocketBase)
- Token exchange endpoint works (returns JWT access token)
- Plugin capabilities and configuration are correct

### ❌ **What's NOT Working**
- **Token activation issue**: Access tokens return "Token is not active" when used
- Import/sync process cannot complete
- Activities are not being imported from Garmin

---

## 🔍 **Key Discoveries**

### 1. **Garmin Authentication Flow**
Garmin Personal API uses **two-step authentication**:
```
1. Exchange API Key for Access Token
   POST https://personal.connectapi.garmin.com/personalapitoken-service/personal-api/personal-api-key/token
   Authorization: ApiKey {api_key}
   
2. Use Access Token for API calls
   GET https://personal.connectapi.garmin.com/activitylist-service/activities/...
   Authorization: Bearer {access_token}
```

### 2. **Token Structure**
- API Key format: `GCDC_{uuid}.{uuid}.{base64}.{timestamp}`
- Access Token: JWT with RS256 signature
- JWT contains: scopes, client_id, garmin_guid, managed_status

### 3. **The Blocking Issue**
Access tokens have `managed_status: "DEFAULT"` instead of `"ACTIVE"`.
Garmin's API returns: `401 - Token is not active`

**Hypothesis**: Tokens need activation time (10-30 seconds) OR require additional activation step in Garmin's backend.

---

## 📋 **Action Items (Priority Order)**

### 🔥 **P0: Critical - Token Activation**
- [ ] Test with **fresh token** immediately after exchange (no delay)
- [ ] Verify if tokens need **10-30 second delay** before becoming active
- [ ] Check if `managed_status` changes from "DEFAULT" to "ACTIVE" over time
- [ ] Test token exchange from **within Docker container** (network/clock sync issues?)
- [ ] Verify **Garmin application status** in developer portal is fully activated

### 🔧 **P1: Plugin Code Updates**
Once token activation works:
- [ ] Modify `wanderer-garmin-plugin/main.go` to implement two-step auth:
  1. Exchange saved API key for access token
  2. Use access token for Garmin API calls
- [ ] Add error handling for token exchange failures
- [ ] Add token refresh logic (tokens have ~2 year expiration based on JWT)
- [ ] Test token caching/reuse for multiple API calls

### 🧪 **P2: Testing & Validation**
- [ ] Create test script to verify token exchange + API call in one flow
- [ ] Test with activities endpoint: `/activitylist-service/activities/search/activities`
- [ ] Test with user profile endpoint: `/userprofile-service/userprofile`
- [ ] Verify Meilisearch indexing of imported activities
- [ ] Test full import flow end-to-end

### 📝 **P3: Documentation**
- [ ] Document Garmin authentication flow in plugin README
- [ ] Update plugin manifest if any changes needed
- [ ] Document troubleshooting steps for token activation issues
- [ ] Add API key validation script

---

## 🔬 **Troubleshooting Notes**

### Error Messages
- `401 - Token is not active` = Token exists but not activated by Garmin
- `401 - invalid_token` = Token format/structure issue
- `400 - Authentication missing` = No token provided

### Working API Key
```
GCDC_62a8961e-0f89-448a-8780-8b89c2738095.69942ea9-687e-4c30-bfc9-7e6e7e9dcbee.gCH+ZSgvlGDSwhu+GtcdK4zWNMbXrJwFz5fmITSyxhI.3617965035
```

### Token Exchange Command (WORKS)
```bash
curl -X POST 'https://personal.connectapi.garmin.com/personalapitoken-service/personal-api/personal-api-key/token' \
  -H 'Authorization: ApiKey GCDC_62a8961e-0f89-448a-8780-8b89c2738095.69942ea9-687e-4c30-bfc9-7e6e7e9dcbee.gCH+ZSgvlGDSwhu+GtcdK4zWNMbXrJwFz5fmITSyxhI.3617965035'
```

### API Call Command (FAILS with "Token is not active")
```bash
curl -H 'Authorization: Bearer {access_token}' \
  'https://personal.connectapi.garmin.com/activitylist-service/activities/search/activities?start=0&limit=1'
```

---

## 📚 **Reference Material**

### Plugin Files
- `/opt/wanderer-dev/wanderer-garmin-plugin/plugin.json` - Plugin manifest
- `/opt/wanderer-dev/wanderer-garmin-plugin/main.go` - Plugin code
- `/opt/wanderer-dev/data/plugins/garmin/plugin.json` - Installed plugin config
- `/opt/wanderer-dev/data/plugins/garmin/plugin.wasm` - Compiled plugin
- `/opt/wanderer-dev/data/plugins/garmin/icon.svg` - Plugin icon

### Helper Scripts
- `/opt/wanderer-dev/wanderer-garmin-plugin/check_plugin_status.py` - Check plugin instance status
- `/opt/wanderer-dev/wanderer-garmin-plugin/check_and_trigger_import.py` - Trigger import
- `/opt/wanderer-dev/wanderer-garmin-plugin/trigger_sync.py` - Manual sync trigger
- `/opt/wanderer-dev/wanderer-garmin-plugin/clear_meilisearch.py` - Clear search index
- `/opt/wanderer-dev/wanderer-garmin-plugin/SCRIPTS.md` - Script documentation

### Related Issues
- Meilisearch ports: 7700 (HTTP), 7701 (Admin)
- PocketBase encryption: Keys are encrypted at rest
- Container: `wanderer-dev-db-1` (port 8091)
- Plugin instance: `z3ahsdr1i0pfty9`

---

## 🎯 **Success Criteria**

- [ ] Plugin loads without "Setup-Fehler"
- [ ] API key saves correctly in database
- [ ] Token exchange succeeds
- [ ] Access token becomes active
- [ ] Activities import successfully
- [ ] Activities appear in Wanderer UI
- [ ] Meilisearch indexes new activities

---

*Last updated: 2026-11-13*
*Status: Token activation issue identified, awaiting resolution*
