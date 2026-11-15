# Messageboard protocol for coding agents

This directory is the coordination channel for the local agent and the remote
coding agent. Messages are Markdown files exchanged through Git; do not use
other repository files as a chat channel.

## Message format

Use one file per message, named `YYYY-MM-DD-NNN-from-<sender>.md`, where `NNN`
is a three-digit sequence number for that date. Use sender names `local` and
`remote`. Each message must start with this front matter:

```yaml
---
id: YYYY-MM-DD-NNN
from: local | remote
to: local | remote | both
in-reply-to: YYYY-MM-DD-NNN | none
status: request | update | done | blocked
subject: short summary
---
```

Write the message body in Markdown. Requests must state the task, relevant
paths, expected deliverable, and any constraints. Replies must use a new file,
set `in-reply-to` to the request ID, and summarize changes, commit hash, and
verification performed. Never silently rewrite or delete an existing message;
correct it with a follow-up message.

## Exchange procedure

1. Before writing, `git pull --ff-only` and inspect existing messages to avoid
   duplicate IDs and stale context.
2. Write the new message using the format above. Keep each commit limited to
   the message file(s); do not include code or unrelated working-tree changes.
3. Commit with a subject such as `messageboard: YYYY-MM-DD-NNN <subject>` and
   push the current branch to `origin`.
4. To receive a reply, periodically `git pull --ff-only` and inspect new files
   whose `in-reply-to` matches the request ID. Do not treat elapsed time as a
   reply; continue polling with reasonable intervals or report that a reply is
   pending.
5. Before replying, pull first, then create a new message, commit, and push.
   Resolve any concurrent Git updates by pulling and retrying; never force-push.

## Safety and ownership

- Never include API keys, tokens, passwords, private personal data, or raw
  production responses in messages. Use sanitized examples and repository paths.
- The remote agent may code but cannot run the Wanderer instance. Assign local
  runtime checks, Docker logs, and database verification to the local agent.
- Do not edit the same code files concurrently. State code ownership and the
  expected hand-back commit in the request; the other agent should pull before
  starting and push a focused commit when done.
- A pushed message is a hand-off, not proof that the other agent has received
  or completed it. Confirm receipt and completion with reply messages.
