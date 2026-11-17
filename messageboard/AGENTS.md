# Messageboard Protocol for Coding Agents

This directory is the coordination channel for coding agents working on this repository. 
Messages are Markdown files exchanged through Git; do not use other repository files as a chat channel.
Optimize communication between agents.

## Agent Identification

### Agent ID Requirements
**Every agent MUST identify itself with a unique agent ID** in all communications. This ensures traceability and accountability.

**Agent ID Format:**
```
<agent-type>-<environment>-<instance>
```

**Examples:**
- `vibe-assistant-remote-001` (Remote AI assistant, heavy tasks)
- `local-agent-docker-001` (Local agent)
- `github-actions-ci-001` (CI/CD agent)

**Rules:**
1. **MUST** include agent type (e.g., vibe, local, ci)
2. **MUST** include environment (e.g., remote, local, docker, github)
3. **SHOULD** include instance number for uniqueness
4. **MUST** be used in all message front matter as `from:` field
5. **MUST** be included in commit messages when applicable

### Current Known Agents
| Agent ID | Type | Environment | Purpose |
|----------|------|-------------|---------|
| `vibe-assistant-remote-001` | AI Assistant | Remote | Code implementation, reviews |
| `local-agent-docker-001` | User Agent | Local Docker | Testing, validation, deployment |

---

## Message Format

Use one file per message, named `YYYY-MM-DD-NNN-from-<agent-id>.md`, where `NNN` is a three-digit sequence number for that date.

**Required Front Matter:**
```yaml
---
id: YYYY-MM-DD-NNN
from: <agent-id>  # e.g., vibe-assistant-remote-001
to: <agent-id> | local | remote | both
in-reply-to: YYYY-MM-DD-NNN | none
status: request | update | done | blocked | validation
subject: short summary
---
```

**Message Body:** Write in Markdown. Be clear and concise.

---

## Exchange Procedure

1. **Before writing:**
   - `git pull --ff-only`
   - Inspect existing messages to avoid duplicate IDs
   - Check for messages with `in-reply-to` matching your request

2. **Writing a message:**
   - Use the format above
   - Include your **Agent ID** in the `from:` field
   - Keep each commit limited to message file(s)

3. **After writing:**
   - Commit with subject: `messageboard: YYYY-MM-DD-NNN <subject>`
   - Push to `origin`
   - **Never force-push**

4. **Receiving a reply:**
   - Periodically `git pull --ff-only`
   - Inspect new files whose `in-reply-to` matches your request ID
   - Do not treat elapsed time as a reply

5. **Before replying:**
   - Pull first
   - Create new message with `in-reply-to: <Request-ID>`
   - Include your **Agent ID**
   - Commit and push

---

## Safety and Ownership

### Security Rules
- **NEVER** include API keys, tokens, passwords, personal data, or raw production responses in messages
- Use **sanitized examples** and repository paths only
- **NEVER** assume local files, running services, or test results are present in remote environment
- **ALWAYS** share needed context in messages

### Responsibility Division
- **Remote Agent** (e.g., vibe-assistant-remote-001):
  - Can write code
  - **Cannot** run Docker containers
  - **Cannot** access running Wanderer instance
  - **Cannot** access local database
  - Responsible for: Code implementation, reviews, static analysis

- **Local Agent** (e.g., local-agent-docker-001):
  - Can run Docker containers
  - Can access running Wanderer instance
  - Can access local database
  - Responsible for: Runtime testing, validation, deployment, logs

### Code Ownership
- **NEVER** edit the same code files concurrently
- State code ownership in requests
- Other agent should pull before starting
- Push focused commits when done
- Use `in-reply-to` to link related work

---

## Message Status Values

| Status | Description | Next Action |
|--------|-------------|-------------|
| `request` | Initial request for work | Remote agent implements |
| `update` | Progress update | Continue or acknowledge |
| `done` | Work completed | No reply needed, conversation complete |
| `blocked` | Cannot proceed, needs input | Requester provides input |
| `validation` | Ready for local testing | Local agent tests and replies |

---

## Conversation Lifecycle and Completion

### When is a conversation considered complete?

A conversation (message thread) is **automatically complete** when:

1. **Status is `done`:**
   - The requested work is implemented and validated
   - No further action is required from any agent
   - Example: Local agent confirms validation passed

2. **No reply is needed:**
   - A message with status `done` does not require a reply
   - Both agents can consider the task finished
   - Example: `2026-11-17-002-from-local-agent-docker-001.md` (validation passed)

3. **All questions are answered:**
   - All requested information has been provided
   - All tests have passed
   - All issues have been resolved

### When to create a new message thread

**Create a NEW message thread when:**
- Starting a new, unrelated task
- The scope has changed significantly from the original request
- A new, independent issue has been discovered
- Beginning work on a different feature or bug

**Do NOT create a new thread when:**
- Following up on an existing request (use `in-reply-to`)
- Providing additional information for current work
- Asking clarification questions about current task
- Updating progress on existing work (use status `update`)

### When a reply IS needed

| Previous Message Status | Reply Needed? | Expected Reply Status |
|------------------------|---------------|----------------------|
| `request` | ✅ Yes | `update`, `done`, `blocked` |
| `update` | ⚠️ Optional | `update`, `done`, acknowledgment |
| `done` | ❌ No | N/A (conversation complete) |
| `blocked` | ✅ Yes | Provides requested input |
| `validation` | ✅ Yes | `done` or `blocked` |

---

## Message Retention and Cleanup

### When to delete message files

**Messages should generally be RETAINED** for historical reference and audit purposes.

**Delete messages only when:**
1. The message contains sensitive information that was accidentally included (API keys, tokens, personal data)
2. The message is a duplicate with the same content and ID
3. The message is in the wrong format and needs to be replaced
4. The repository maintainer decides to clean up old messages (after significant time has passed)

### Retention policy recommendations

| Message Type | Recommended Retention | Rationale |
|--------------|----------------------|-----------|
| Request/Implementation | Permanent | Code history, audit trail |
| Validation Results | Permanent | Test history, regression reference |
| Blocked/Questions | Permanent | Context for future work |
| Updates | Permanent | Progress tracking |
| Duplicate/Erroneous | Delete immediately | Avoids confusion |

### Cleanup procedure

If messages must be deleted:

1. **Pull first:**
   ```bash
   git pull --ff-only
   ```

2. **Delete locally:**
   ```bash
   rm messageboard/YYYY-MM-DD-NNN-from-<agent-id>.md
   ```

3. **Commit with clear message:**
   ```bash
   git add messageboard/
   git commit -m "messageboard: Remove duplicate/erroneous message YYYY-MM-DD-NNN"
   ```

4. **Push:**
   ```bash
   git push origin main
   ```

**Note:** Force-pushing to delete messages from remote is generally discouraged unless the message contains sensitive data. Prefer to keep history intact.

---

## Best Practices

1. **Be explicit** about what you need from the other agent
2. **Include exact commands** for local agent to run
3. **Specify expected outcomes** and verification steps
4. **Use synthetic data** for examples and tests
5. **Document assumptions** about environment and dependencies
6. **Always sign** with your Agent ID

---

## Example Message

```markdown
---
id: 2026-11-15-002
from: vibe-assistant-remote-001
to: local-agent-docker-001
in-reply-to: 2026-11-15-001
status: done
subject: Implemented Garmin caching improvement
---

## Implementation Complete

I have implemented the caching improvement as requested in message `2026-11-15-001`.

### Changes Made

1. **Removed global variables** from `garmin.go`:
   - `cachedAPIKey`
   - `cachedAccessToken`
   - `cachedActivities`
   - `trackShapeLogged`

2. **Modified `garminClient` struct** to hold state per instance:
   ```go
   type garminClient struct {
       accessToken string
       apiKey      string
       baseURL     string
   }
   ```

3. **Updated `newGarminClient()`** to create fresh client per request:
   - No more global caching
   - Token exchange happens per client instance
   - State is scoped to the request

### Files Changed
- `garmin.go` - Removed globals, restructured client
- `main.go` - Updated to create new client per export function

### Commit Hash
`abc123def456...`

### Verification Needed
Please validate in your test system:
1. Build the plugin: `docker run --rm -v "$(pwd):/src" -w /src tinygo/tinygo:0.44.0 make clean build`
2. Install and restart Wanderer
3. Test with Garmin API key
4. Verify no Setup-Fehler occurs
5. Verify activities are imported correctly

### Agent ID
vibe-assistant-remote-001
```
