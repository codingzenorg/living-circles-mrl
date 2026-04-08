# Slice: Initial Event-Driven Local Food Transport

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a player-following viewport and cached orientation summaries
- Go server with authoritative simulation state
- shared runtime contract already optimized by viewport culling, dual cadence, compact minimap summaries, reduced local precision, and event-driven orientation refresh

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps high-cadence local circle transport intact, but stops resending unchanged local food detail on every tick. Instead, local viewport food detail should refresh when the visible food set materially changes or when a slower fallback refresh requires revalidation.

## Discovery Scope

Reduce high-cadence transport cost further by making local viewport food transmission event-driven.

This slice should:

- keep player and local autonomous-circle detail at the current interactive cadence
- keep current orientation-summary transport unchanged
- refresh local food detail when the visible viewport food set materially changes
- retain a deterministic fallback refresh so the client is never left indefinitely stale

This slice does **not** attempt to implement:

- entity deltas
- circle-state caching
- compression
- binary transport
- gameplay changes

## Why This Slice Next

The transport path is much lighter than the original baseline, and orientation refreshes are now driven by summary change. The next visible redundancy is local food detail: foods in the viewport often stay unchanged across many ticks, but they are still resent at the same cadence as moving circles.

The next pressure is therefore high-cadence local redundancy:

- local circle motion is still worth sending every tick
- local food detail usually changes only on collection, regeneration, or viewport transition
- the client already has a proven cache-and-reuse pattern for stale orientation data

This is the narrowest next step because it:

- targets a largely static local payload
- preserves current local circle responsiveness
- reuses an existing client-side caching concept instead of redesigning the protocol

## Use-Case Contract

### Use Case

`RefreshLocalFoodOnlyWhenViewportFoodChanges`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- local viewport circle detail already arrives every tick
- the client can already reuse stale orientation summaries
- viewport food membership can be derived deterministically at the transport boundary

### Trigger

The server emits a transport snapshot and decides whether the current visible local food set materially differs from the last transmitted local food set.

### Success Outcome

- local circle play remains fully responsive
- food still appears authoritative and trustworthy in the viewport
- unchanged local food detail is omitted on most quiet ticks
- average per-client transport cost drops below the current event-driven orientation baseline

### Failure Or Rejection Cases

- if food can remain stale for too long, the slice failed its local trust guarantee
- if circle responsiveness is coupled to food refresh cadence, scope drifted
- if the slice forces a wider delta-snapshot redesign, scope is exceeded

## Main Business Rules

1. This is a transport policy slice, not a gameplay slice.
2. Local circle detail remains unchanged.
3. Local viewport food detail refreshes when the visible food set materially changes.
4. A deterministic fallback refresh cadence must still exist for local food.
5. The optimization must remain measurable over deterministic windows.

## Minimal Domain Concepts In Scope

- `Material Local Food Change`
- `Fallback Local Food Refresh`
- `Average Per-Client Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful event-driven interpretation:

- derive a deterministic signature for the visible local food set
- send local food detail when that signature changes
- otherwise omit local food detail unless a slower fallback interval is reached
- let the client reuse the last valid local food detail between refreshes

This avoids moving to generalized entity deltas while still preventing redundant local food payloads.

## Required Runtime Contract Changes

Probably yes.

The current contract always sends `foods` as a present array. This slice will likely need:

- a freshness signal for local food detail
- the ability to omit or null local food detail on stale ticks

## Required Ports Or Boundaries

- server-side local food signature and refresh decision
- transport measurement tests over deterministic windows
- browser client reuse of the last valid local food detail
- shared contract schema for stale-food ticks

## Build Guidance

- keep visible-food signature derivation deterministic
- compare the current visible local food set against the last transmitted one
- retain a slower fallback refresh so quiet local viewports still revalidate
- preserve current local circle cadence and shape
- measure average cost over a deterministic tick window and compare it to the current event-driven orientation baseline

## Initial Test Plan

### Server tests

- local food refresh happens when the visible food set changes
- local food refresh is skipped when the visible food set is unchanged
- fallback local food refresh still occurs after a deterministic interval
- average transport cost drops below the current event-driven orientation baseline

### Contract tests

- stale local-food ticks remain explicit and parseable

### Integration tests

- local circle play remains unchanged
- the client continues to keep a valid visible food set between refreshes

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect and receive local viewport circle detail every tick
2. observe that local food detail refreshes when visible food changes
3. observe that it is skipped when visible food stays unchanged
4. verify that a slower fallback refresh still occurs
5. compare average transport cost to the current event-driven orientation baseline

## Done Criteria

- local circle play remains unchanged
- local food detail is driven by visible-food change plus a slower fallback interval
- average per-client transport cost drops below the current event-driven orientation baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- generalized entity deltas
- circle-state caching
- compression
- binary transport
- multiplayer-specific interest management
