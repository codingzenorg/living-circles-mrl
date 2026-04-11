# Slice: Initial Active Local Autonomous Summary Compaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using the current reduced-cost active transport path
- current coarser active orientation-summary path already built
- passive observer path unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, cadence, and orientation behavior unchanged, but reduces the next repeated active payload family by making active local autonomous detail lighter while keeping nearby NPC awareness usable.

## Discovery Scope

Reduce active payload fanout further without revisiting the already-optimized orientation path or the already-rejected local-food compaction path.

This slice should:

- preserve the current active local-detail cadence
- preserve the current active orientation-summary path unchanged
- preserve the current passive observer path unchanged
- reduce the size of active local autonomous detail on fresh autonomous ticks

This slice does **not** attempt to implement:

- another orientation optimization
- another cadence-policy change
- gameplay changes
- passive transport redesign
- prediction or interpolation

## Why This Slice Next

The latest active-path evidence now says:

- orientation support is still dominant, but materially smaller than before
- local-food compaction was tested as the next target and produced no measurable reduction in the active baseline
- the next worthwhile active target should therefore move to another repeated payload family instead of squeezing a no-op path

Local autonomous detail is the best bounded next target:

- it is still present on fresh autonomous ticks
- it is a repeated local family with more structure than local food detail
- nearby NPC awareness matters for active play, but the wire representation may still be heavier than the viewport visibly needs

## Use-Case Contract

### Use Case

`CompactActiveLocalAutonomousDetail`

### Primary Actor

The Go transport layer broadcasting authoritative active-client snapshots.

### Pre-conditions

- active orientation compaction already reduced the dominant active family materially
- local-food compaction has already been attempted and shown to be a non-win
- local autonomous detail still appears on fresh autonomous ticks

### Trigger

An active snapshot is prepared on a tick where local autonomous detail is fresh.

### Success Outcome

- active local play remains responsive
- nearby NPC awareness remains usable
- local autonomous detail is lighter on the wire than the current baseline
- one-active and two-active baselines drop below the current stable readings

### Failure Or Rejection Cases

- if nearby NPC awareness becomes too vague for active play, the slice failed
- if the slice changes cadence or orientation behavior, scope drifted
- if gameplay semantics change, scope is exceeded

## Main Business Rules

1. This is an active local-autonomous payload-compaction slice, not a cadence slice.
2. Active local detail remains present.
3. Active orientation behavior remains unchanged.
4. Passive observer behavior remains unchanged.
5. Nearby autonomous-circle information may be serialized more coarsely if active play readability remains acceptable.

## Minimal Domain Concepts In Scope

- `Active Local Autonomous Detail`
- `Nearby NPC Awareness`
- `Fresh Autonomous Tick`
- `Active Payload Size`

## Bounded Optimization Interpretation

This slice chooses the smallest plausible next reduction:

- keep the current active autonomous refresh policy
- keep the current contract shape recognizable
- make active local autonomous detail itself coarser or lighter than the current local baseline

## Required Runtime Contract Changes

Prefer none, or only the smallest contract-preserving change.

## Required Ports Or Boundaries

- active snapshot assembly in the Go transport layer
- browser rendering of nearby autonomous circles
- transport measurement helpers updated with the new active baseline

## Build Guidance

- do not revisit orientation behavior in this slice
- do not revisit autonomous refresh cadence in this slice
- target only the payload weight of active local autonomous detail
- remeasure the active payload breakdown and one-active/two-active baselines after the change

## Initial Test Plan

### Server or measurement tests

- active transport measurement stays deterministic
- active local autonomous detail contributes less payload than before on fresh autonomous ticks
- one-active and two-active aggregate bytes drop below the current stable readings

### Contract tests

- current contract validation remains green, with only minimal updates if strictly required

### Integration tests

- active clients still receive usable nearby NPC awareness
- active local play remains responsive

## Scenario Definition

Start the current server and measure the optimized active local-autonomous path.

Scenario steps:

1. build one active snapshot under the new local-autonomous compaction rule
2. remeasure the active payload breakdown
3. remeasure one-active and two-active active baselines
4. confirm active nearby NPC awareness still behaves as before

## Done Criteria

- active local-detail behavior stays responsive
- active local autonomous payload is smaller than the current baseline
- one-active and two-active active baselines drop below the current stable readings
- no gameplay semantics change

## Out Of Scope Follow-Ups

- another orientation optimization
- cadence-policy changes
- passive transport redesign
- gameplay changes
