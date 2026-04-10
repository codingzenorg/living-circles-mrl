# Runtime Evidence

## Date

2026-04-10

## Source

Direct local user observation during manual demo use.

## Observation

Opening a second browser session against the local demo caused visible slowdown:

- the simulation became slower overall
- player movement felt slower
- the slowdown was noticeable from ordinary play, not only from instrumentation

## Why It Matters

This is stronger than a hypothetical optimization concern because it is already observable in normal use with only a small number of concurrent clients.

The pressure is no longer:

- "maybe transport is too chatty"
- "maybe render work is too heavy"

It is now:

- real multi-client responsiveness degrades under ordinary local use
- the next loop should treat concurrent responsiveness as explicit runtime evidence

## Immediate Implication

The next bounded slice should focus on concurrent responsiveness under multiple clients rather than continuing speculative client-side micro-optimizations by default.
