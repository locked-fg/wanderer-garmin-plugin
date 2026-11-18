# Testing guide

This guide sets shared expectations for tests across the Wanderer project. Component-specific commands, fixtures, and assertions belong in that component's test documentation or skill and should link here rather than duplicate these principles.

## Choose the smallest useful test

- Start with focused unit or fixture tests. Use integration or end-to-end tests for behavior that depends on real service boundaries.
- Prefer small, deterministic fixtures that cover normal cases, error cases, and important boundaries. Avoid full datasets when a representative sample proves the behavior.
- Define the expected result before running the test, including what counts as failure, timeout, and partial completion.
- Run slower or stateful checks only when the change needs that evidence.

## Make evidence meaningful

- Verify behavior through the interface users or other components rely on, such as an API or UI. Use direct database or internal state inspection to diagnose gaps that interface does not expose.
- Record a baseline and compare before/after values when background work or existing data can affect results.
- A running status, nonzero count, or successful command exit is not sufficient by itself; assert the actual outcome relevant to the change.
- Distinguish completed and passed tests from samples, timeouts, failures, and runs stopped early.

## Control time and side effects

- Use bounded waits and explicit timeouts. Detect stalled work and avoid waiting for unrelated background tasks.
- Prefer isolated test data and disposable environments for tests that write or delete persistent state.
- Before destructive actions, make the target and expected impact clear and ensure the data can be recovered. Record persistent changes made by the test.
- Clean up only resources created by the test. Restore modified configuration and schedules even on failure or interruption, then verify that restoration.

## Protect credentials and diagnostics

- Keep credentials out of source, command history, logs, and reports. Avoid printing resolved configuration when it may contain secrets.
- Filter logs to the test being investigated and redact sensitive values before sharing output.
- Prefer supported APIs for user-visible state and configuration changes. Use direct database access only when the API lacks a needed capability, and avoid treating internal state alone as proof of externally visible behavior.

## Keep test tooling maintainable

- Document each helper's prerequisites, required credentials, read/write behavior, destructive effects, and verification method.
- A helper should verify that the intended operation actually happened, not merely set a status or issue a request.
- Make test procedures discoverable from the project development guide, and keep component-specific instructions linked to this shared guide.
- Record coverage gaps and limitations with the test results so partial evidence is not mistaken for a pass.
