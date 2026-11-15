# Project notes

- This directory is the Garmin plugin source repository within the Wanderer development checkout.
- Keep implementation, tests, fixtures, and setup instructions in sync.
- Do not store API keys, access tokens, or other credentials in source files, documentation, logs, or commits.
- Build with the TinyGo Docker image and Makefile as documented in `DEVELOPMENT.md`.
- Optional Python tools live in `scripts/`; see `SCRIPTS.md` before using diagnostic or cleanup tools.
