# Slice: Initial Active Local Food Summary Compaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using the current reduced-cost active transport path
- current coarser active orientation-summary path already built
- passive observer path unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, active cadence, and active orientation behavior unchanged, but reduces the next most relevant active payload family by making active local food detail lighter while keeping local food awareness usable.

## Discovery Scope

Reduce active payload fanout further without revisiting the already-optimized orientation path.

This slice should:

- preserve the current active local-detail cadence
- preserve the current active orientation-summary path unchanged
- preserve the current passive observer path unchanged
- reduce the size of active local food detail on fresh food ticks

This slice does **not** attempt to implement:

- another orientation optimization
- another cadence-policy change
- gameplay changes
- passive transport redesign
- prediction or interpolation

## Why This Slice Next

The latest stable active-path evidence now shows:

- active orientation support is still dominant, but materially smaller than before
- one-active aggregate bytes dropped from `8942` to `7524`
- two-active aggregate bytes dropped from `17884` to `13980`

That means the last orientation-focused slice succeeded. The next pressure is no longer “keep squeezing orientation first at all costs.” The cleaner next step is to move to the next materially repeated active payload family.

Local food detail is the best bounded next target:

- it is still present in active local-detail snapshots
- it is already change-driven, so reducing its per-refresh payload is cleaner than revisiting cadence again
- food awareness matters locally, but active rendering likely does not need full float-detail precision for every nearby food on the wire

## Use-Case Contract

### Use Case

`CompactActiveLocalFoodDetail`

### Primary Actor

The Go transport layer broadcasting authoritative active-client snapshots.

### Pre-conditions

- active orientation compaction already reduced the dominant active family materially
- local food detail still appears on fresh food refreshes
- local active play still depends on usable nearby food awareness

### Trigger

An active snapshot is prepared on a tick where local food detail is fresh.

### Success Outcome

- active local play remains responsive
- nearby food awareness remains usable
- local food detail is lighter on the wire than the current baseline
- one-active and two-active baselines drop below the current stable readings

### Failure Or Rejection Cases

- if nearby food awareness becomes too vague for active play, the slice failed
- if the slice changes cadence or orientation behavior, scope drifted
- if gameplay semantics change, scope is exceeded

## Main Business Rules

1. This is an active local-food payload-compaction slice, not a cadence slice.
2. Active local detail remains present.
3. Active orientation behavior remains unchanged.
4. Passive observer behavior remains unchanged.
5. Nearby food information may be serialized more coarsely if local food awareness remains usable.

## Minimal Domain Concepts In Scope

- `Active Local Food Detail`
- `Nearby Food Awareness`
- `Fresh Food Tick`
- `Active Payload Size`

## Bounded Optimization Interpretation

This slice chooses the smallest plausible next reduction:

- keep the current active food refresh policy
- keep the current contract shape recognizable
- make active local food detail itself coarser or lighter than the current local baseline

## Required Runtime Contract Changes

Prefer none, or only the smallest contract-preserving change.

## Required Ports Or Boundaries

- active snapshot assembly in the Go transport layer
- browser rendering of local food detail
- transport measurement helpers updated with the new active baseline

## Build Guidance

- do not revisit orientation behavior in this slice
- do not revisit food refresh cadence in this slice
- target only the payload weight of active local food detail
- remeasure the active payload breakdown and one-active/two-active baselines after the change

## Initial Test Plan

### Server or measurement tests

- active transport measurement stays deterministic
- active local food detail contributes less payload than before on fresh food ticks
- one-active and two-active aggregate bytes drop below the current stable readings

### Contract tests

- current contract validation remains green, with only minimal updates if strictly required

### Integration tests

- active clients still receive usable local food awareness
- active local play remains responsive

## Scenario Definition

Start the current server and measure the optimized active local-food path.

Scenario steps:

1. build one active snapshot under the new local-food compaction rule
2. remeasure the active payload breakdown
3. remeasure one-active and two-active active baselines
4. confirm active local food awareness still behaves as before

## Done Criteria

- active local-detail behavior stays responsive
- active local food payload is smaller than the current baseline
- one-active and two-active active baselines drop below the current stable readings
- no gameplay semantics change

## Out Of Scope Follow-Ups

- another orientation optimization
- cadence-policy changes
- passive transport redesign
- gameplay changes
