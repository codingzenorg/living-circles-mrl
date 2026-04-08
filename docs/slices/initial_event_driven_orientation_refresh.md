# Slice: Initial Event-Driven Orientation Refresh

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a player-following viewport and passive minimap
- Go server with authoritative simulation state
- shared runtime contract already optimized by viewport culling, dual cadence, compact minimap summaries, and reduced local precision

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps the current compact summary and local-detail transport shape, but stops refreshing orientation data on a purely fixed cadence. Instead, orientation refreshes should happen when the compact whole-world summary has materially changed or when a much slower safety cadence requires a refresh.

## Discovery Scope

Reduce average transport cost further by making orientation refreshes responsive to summary change instead of only to a fixed timer.

This slice should:

- keep the local viewport detail cadence unchanged
- keep the compact minimap summary representation
- refresh orientation when the summary materially changes
- retain a deterministic low-frequency safety refresh so the client is never left indefinitely stale

This slice does **not** attempt to implement:

- delta snapshots
- prediction
- compression
- multiple websocket channels
- gameplay-driven simulation throttling

## Why This Slice Next

The transport path is already much lighter than the original baseline, but the minimap/orientation refresh still happens on a fixed schedule. That means the server can still send a compact whole-world summary even when that summary is effectively unchanged since the last refresh.

The next pressure is therefore refresh relevance:

- local play detail should remain frequent
- orientation refreshes should be driven by actual summary change, not just the passage of a few ticks

This is the narrowest next step because it:

- builds directly on the existing compact-summary transport
- avoids a larger protocol redesign
- targets avoidable refreshes rather than already-optimized local payloads

## Use-Case Contract

### Use Case

`RefreshOrientationOnlyWhenSummaryChanges`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- local viewport detail already arrives every tick
- compact minimap summaries already exist
- the client already reuses the last valid orientation summary

### Trigger

The server emits a transport snapshot and decides whether the compact orientation summary is meaningfully different from the last refresh.

### Success Outcome

- local play remains fully responsive
- the minimap still remains orienting and trustworthy
- orientation refreshes happen when needed, not merely on a fixed short cadence
- average per-client transport cost drops below the current reduced-precision compact-summary baseline

### Failure Or Rejection Cases

- if the minimap can remain stale for too long, the slice failed its orientation guarantee
- if the slice requires a large state-sync redesign, scope drifted
- if gameplay semantics change, scope is exceeded

## Main Business Rules

1. This is an orientation-refresh policy slice, not a gameplay slice.
2. Local viewport detail remains unchanged.
3. Compact orientation summaries refresh when they materially change.
4. A deterministic fallback refresh cadence must still exist.
5. The optimization must remain measurable over deterministic windows.

## Minimal Domain Concepts In Scope

- `Material Orientation Change`
- `Fallback Orientation Refresh`
- `Average Per-Client Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful event-driven interpretation:

- derive a deterministic compact-summary signature
- send orientation refreshes when that signature changes
- otherwise omit the refresh unless a slower fallback interval is reached

This avoids jumping to deltas or complex subscription logic while still preventing redundant minimap refreshes.

## Required Runtime Contract Changes

Possibly none.

The current `orientation_fresh` signal may already be sufficient if the server changes only the refresh policy.

## Required Ports Or Boundaries

- server-side orientation summary builder and refresh decision
- transport measurement tests over deterministic windows
- browser client reuse of the last valid summary, which should already remain valid

## Build Guidance

- keep the compact summary builder deterministic
- compare the current compact summary against the last transmitted one
- retain a slower fallback refresh so long-lived quiet worlds still revalidate
- measure average cost over a deterministic tick window and compare it to the current baseline

## Initial Test Plan

### Server tests

- orientation refreshes happen on summary change
- orientation refreshes are skipped when the summary is unchanged
- fallback refresh still occurs after a deterministic interval
- average transport cost drops below the current reduced-precision compact-summary baseline

### Contract tests

- the current contract shape remains explicit and parseable

### Integration tests

- local viewport play remains unchanged
- the client continues to keep a valid minimap summary between refreshes

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect and receive local viewport updates every tick
2. observe that compact orientation refreshes occur when the whole-world summary changes
3. observe that they are skipped when the summary remains unchanged
4. verify that a slower fallback refresh still occurs
5. compare average transport cost to the current reduced-precision compact-summary baseline

## Done Criteria

- local play remains unchanged
- orientation refreshes are driven by summary change plus a slower fallback interval
- average per-client transport cost drops below the current reduced-precision compact-summary baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- delta snapshots
- compression
- binary transport
- adaptive per-entity priority systems
- multiplayer-specific interest management
