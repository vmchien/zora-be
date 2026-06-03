# kafka consumer

Production-first Kafka consumer helper for Go built on franz-go (kgo).

## What this helper solves

1. **Crash hardening**
    - Recovers from panics in user handler (keeps process alive).
    - Bounded queues to avoid unbounded memory usage.
    - Deterministic shutdown order to avoid goroutine leaks and partial commits.

2. **Message loss minimization**
    - Supports **ManualCommit** mode (recommended), committing offsets **only after successful processing** (at-least-once).
    - Optional DLQ handling for poison messages. If DLQ succeeds, offsets can advance safely.

3. **Mode selection**
    - **Sync**: process records sequentially in the poll loop.
    - **Async**: process concurrently with a worker pool while preserving **per-partition ordering** (one in-flight per partition).

4. **Commit policy selection**
    - **ManualCommit** (recommended, default for Async).
    - **AutoCommit** (best-effort; Async + AutoCommit is risky and opt-in).

## Semantics (read this first)

- **ManualCommit** provides **at-least-once** semantics: records are committed only after handler success (or DLQ success).
    - You must make processing **idempotent**.
    - Duplicates can occur on crash/restart before commit.

- **AutoCommit** is **best-effort** and can drift toward at-most-once depending on timing.
    - **Async + AutoCommit may lose messages**: offsets can be committed while processing is still in-flight.
    - For that reason, Async + AutoCommit requires `AllowRiskyAsyncAutoCommit=true`.

## Install

```bash
go get github.com/twmb/franz-go@latest