# Slice: Initial Event-Driven Active Orientation Refresh

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using `active_local_detail`
- current reduced-cost active transport path with cached minimap orientation on the client
- passive observer path unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, active local detail, and passive observer behavior unchanged, but replaces the current fixed active-orientation cadence with a change-driven refresh policy plus fallback.

## Discovery Scope

Reduce the dominant remaining active payload family by making active orientation refreshes responsive to actual summary change instead of a short fixed timer.

This slice should:

- preserve current active local-detail cadence
- preserve current passive observer behavior unchanged
- preserve the current compact orientation-summary representation
- refresh active orientation when the compact whole-world summary materially changes, plus a slower fallback refresh

This slice does **not** attempt to implement:

- compression
- delta snapshots
- gameplay changes
- passive transport redesign
- client prediction or interpolation

## Why This Slice Next

The latest reassessment now makes the current active-path state explicit:

- full active payload: `3357` bytes
- without orientation support: `1200` bytes
- orientation support remains the dominant serialized active payload family even after the recent reductions

The current code path still refreshes active orientation on a fixed `3`-tick cadence. That means active clients can still receive compact whole-world orientation data simply because a short timer elapsed, even when the orientation summary has not materially changed.

The next pressure is therefore refresh relevance again, but now specifically for the active path:

- active local detail should remain high-cadence
- minimap orientation should remain usable
- active orientation refreshes should happen when needed, not just on a short repeating timer

## Use-Case Contract

### Use Case

`RefreshActiveOrientationOnlyWhenSummaryChanges`

### Primary Actor

The Go transport layer broadcasting authoritative state to active clients.

### Pre-conditions

- active clients already reuse cached orientation support between fresh snapshots
- active local detail already remains high-cadence
- current active orientation still refreshes on a fixed `3`-tick cadence

### Trigger

An active snapshot is prepared and the server decides whether the compact orientation summary is meaningfully different from the last active orientation refresh.

### Success Outcome

- active local play remains responsive
- active minimap orientation remains usable
- active orientation refreshes happen on summary change plus slower fallback instead of a fixed short cadence
- the active payload baseline drops below the current post-reduction reading

### Failure Or Rejection Cases

- if active local responsiveness changes, the slice failed
- if active minimap orientation can remain stale for too long, the slice failed
- if the slice becomes a broader transport redesign, scope drifted

## Main Business Rules

1. This is an active-orientation refresh-policy slice, not a gameplay slice.
2. Active local detail remains unchanged.
3. Passive observer behavior remains unchanged.
4. Active orientation refreshes should happen on material summary change plus deterministic fallback.
5. The active path must be remeasured after the change.

## Minimal Domain Concepts In Scope

- `Active Orientation Refresh`
- `Material Orientation Change`
- `Fallback Orientation Refresh`
- `Active Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful change:

- derive or reuse a deterministic compact-summary signature for the active path
- refresh active orientation when that signature changes
- otherwise skip the refresh until a slower fallback interval requires one

This avoids another protocol redesign while targeting the still-dominant active payload family directly.

## Required Runtime Contract Changes

None expected.

The existing `orientation_fresh` signal should remain sufficient.

## Required Ports Or Boundaries

- active snapshot assembly in the Go transport layer
- current browser caching of the last valid active orientation summary
- transport measurement helpers updated to record the new active baseline

## Build Guidance

- keep the compact orientation summary deterministic
- do not change active local player/autonomous/food detail cadence
- remeasure both one-active and two-active baselines after the change
- confirm active orientation remains plausibly usable with the existing client cache

## Initial Test Plan

### Server or measurement tests

- active orientation refreshes happen on compact-summary change
- active orientation refreshes are skipped when the compact summary is unchanged
- fallback refresh still occurs after a deterministic interval
- active transport cost drops below the current post-reduction baseline

### Contract tests

- current contract validation remains green without shape changes

### Integration tests

- active clients still receive responsive local detail
- active clients still retain usable orientation support across the optimized path

## Scenario Definition

Start the current server and measure the optimized active path.

Scenario steps:

1. build active snapshots under the new orientation refresh policy
2. remeasure the single-active baseline
3. remeasure the two-active baseline
4. confirm local active detail still behaves as before

## Done Criteria

- active local-detail behavior stays responsive
- active orientation refreshes are change-driven plus fallback
- active transport cost drops below the current post-reduction baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- compression
- delta encoding
- passive transport redesign
- generalized subscriptions
- gameplay changes
