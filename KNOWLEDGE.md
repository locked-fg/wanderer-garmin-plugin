# Project Knowledge Base

This document captures key insights, decisions, and lessons learned during the development of the Garmin plugin for Wanderer. It serves as a reference for future development and troubleshooting.

---

## 🔧 Garmin API Specifics

### Authentication and Token Management

| Aspect | Details | Source |
|--------|---------|--------|
| **Base URL** | `https://personal.connectapi.garmin.com` | API Documentation |
| **Token Exchange Endpoint** | `/personalapitoken-service/personal-api/personal-api-key/token` | `74f39d0` |
| **Token Exchange Header** | `Authorization: ApiKey <api_key>` (not Bearer!) | `74f39d0` |
| **API Request Header** | `Authorization: Bearer <access_token>` | `74f39d0` |
| **Token Activation Time** | **10-30 seconds** after exchange - tokens have `managed_status: "DEFAULT"` initially | `74f39d0`, `95d2cef` |

### Activity Data Endpoints

| Endpoint | Purpose | Response Format |
|----------|---------|----------------|
| `/activitylist-service/activity` | List activities | Paginated JSON |
| `/activity-service/activity/{id}/details` | Activity details | JSON with polyline or metrics |
| `/download-service/activity/{id}` | Download activity | GPX format |

### Response Formats

- **Token Response:** JWT with `access_token`, `expires_in`, `token_type`, and `managed_status` fields
- **Activity List:** Contains `activityID`, `activityName`, `startTime`, `duration`, `distance`, etc.
- **Activity Detail:** Supports two formats:
  - **Polyline format:** `geoPolylineDTO` with encoded coordinates
  - **Metrics fallback:** `activityDetailMetrics` with structured data
- **Activity IDs:** Can collide between different Garmin accounts → never cache globally!

---

## 🐛 Problems & Solutions

### Critical Issues

| Problem | Root Cause | Solution | Commit |
|---------|------------|----------|--------|
| **Setup-Fehler** | Global variables caused cross-account state leakage | Removed all global variables (`cachedAPIKey`, `cachedAccessToken`, `cachedActivities`, `trackShapeLogged`) | `95d2cef` |
| **HTTP 406 NotAcceptableException** | Single connector couldn't handle both `ApiKey` and `Bearer` headers | Created separate connectors: `token_exchange` (no auth) and `api` (with auth) | `74f39d0` |
| **Token not saved in WandererDB** | Wanderer expects only ONE auth context, had two (`api_key` + `bearer_token`) | Consolidated to single `bearer_token` auth context | `765c89f` |

### Configuration Issues

| Problem | Root Cause | Solution | Commit |
|---------|------------|----------|--------|
| **Invalid field name** | Used `secretFields` (plural) instead of `secretField` (singular) | Changed to `secretField` in plugin.json | `ea5592a` |
| **Local plugin install failed** | Makefile paths and INSTALL.md instructions were incorrect | Updated Makefile and INSTALL.md for local development | `3bee8ad` |
| **Activity start times inconsistent** | Different time formats not handled uniformly | Normalized time parsing for `StartTime` and `StartTimeLocal` | `07a0f6a` |
| **Activity response formats not supported** | Parser only handled one format | Added support for both polyline and metrics fallback formats | `996998a` |

---

## 🏗️ Architecture Decisions

### State Management

| Decision | Rationale | Commit |
|----------|-----------|--------|
| **No global variables for caching** | Prevents cross-account data leakage when activity IDs collide; WASM instance lifetime is unpredictable | `95d2cef` |
| **Client per request** | Ensures state isolation between different accounts and sync operations; follows Strava plugin pattern | `95d2cef` |
| **Track shape logging per instance** | Moved `trackShapeLogged` to `garminClient` struct to prevent suppression across requests | `95d2cef` |

### Authentication Strategy

| Decision | Rationale | Commit |
|----------|-----------|--------|
| **Two connectors** | Token exchange requires `ApiKey` header, API calls require `Bearer` header - cannot use same connector | `74f39d0` |
| **Single auth context** | Wanderer expects exactly one auth context per plugin; Garmin API key is provided as `apiKey` field | `765c89f` |
| **Bearer token type** | Matches Wanderer's expected auth context types; API key is treated as bearer token | `765c89f` |

### Code Quality

| Decision | Rationale | Commit |
|----------|-----------|--------|
| **Squash commits before merging** | Keeps main history clean and readable; combines implementation steps into logical changes | `5959f2e` |
| **Always use feature branches** | Main remains stable and functional; enables parallel work by multiple agents | `d3b0316` |

---

## 📁 Documentation

### Development Guidelines

| Document | Purpose | Commit |
|----------|---------|--------|
| **DEVELOPMENT.md** | Build instructions, project structure, test procedures | `d3b0316` |
| **Git Workflow section** | Branching strategy, merge procedures, squash requirements | `d3b0316`, `5959f2e` |

### Agent Coordination

| Document | Purpose | Commit |
|----------|---------|--------|
| **messageboard/AGENTS.md** | Agent identification, message format, exchange procedure | `0d2bb45`, `73c18cf`, `46fb003` |
| **Conversation Lifecycle** | When conversations are complete, when replies are needed | `73c18cf` |
| **Message Retention Policy** | When to keep/delete message files (Git history is sufficient) | `46fb003` |

---

## 💡 Lessons Learned

### Workflow

- **Always use feature branches** - Main should always be in a working state, especially when multiple agents collaborate
- **Squash commits before merging** - Keeps main history clean and makes it easier to understand changes
- **Pull before starting work** - Avoids duplicate work and ensures you have the latest changes
- **Test thoroughly before merging** - Local validation prevents broken main branch

### Garmin API

- **Token activation delay is real** - Garmin Personal API tokens need 10-30 seconds to activate after exchange; plan for this in your code
- **Header types matter** - Token exchange uses `ApiKey` header, API calls use `Bearer` header - they cannot be combined
- **Activity IDs are not unique across accounts** - Never cache activities globally; always fetch per request
- **Response formats vary** - Garmin API can return data in different formats (polyline vs. metrics); handle both gracefully

### Wanderer Integration

- **Wanderer expects one auth context** - Multiple auth contexts cause token saving issues
- **Field names are exact** - `secretField` (singular) works, `secretFields` (plural) does not
- **Plugin manifest must be precise** - Small errors in plugin.json can cause Setup-Fehler

### Code Quality

- **Avoid global state in WASM** - WASM instance lifetime is tied to host usage; global variables can persist across unrelated operations
- **Follow existing patterns** - Strava plugin provides a good reference for architecture and state management
- **Document decisions** - Future maintainers (and agents) need to understand why certain choices were made

---

## 📚 References

- **Strava Plugin:** Reference implementation for plugin architecture and state management
- **Wanderer SDK:** Check `plugins/sdk/` for type definitions and helper functions
- **Garmin Personal API Documentation:** Official API reference for endpoint details and requirements
