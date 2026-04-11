# Slice: Initial Active Orientation Summary Compaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using the current reduced-cost active transport path
- compact minimap-orientation summaries already present in active snapshots
- passive observer path unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, cadence, and freshness policy unchanged, but reduces the payload weight of active orientation support itself by making the active minimap summary coarser than the current compact baseline.

## Discovery Scope

Reduce the dominant remaining active payload family without reusing the failed event-driven active-orientation approach.

This slice should:

- preserve the current active local-detail cadence
- preserve the current active orientation cadence and freshness behavior
- preserve the current passive observer path unchanged
- reduce the serialized size of active orientation support by sending a coarser active-only orientation summary

This slice does **not** attempt to implement:

- another refresh-policy change
- delta snapshots
- gameplay changes
- passive transport redesign
- client prediction or interpolation

## Why This Slice Next

The latest active-path evidence is now clearer:

- orientation support is still the dominant active payload family
- the attempted event-driven active-orientation refresh did not help; in the bounded movement window it refreshed every active snapshot and worsened the active baseline
- the repo therefore still needs an orientation-focused reduction, but not another refresh-policy experiment right away

The next clean move is to reduce the size of each active orientation refresh rather than changing when refreshes happen.

## Use-Case Contract

### Use Case

`CompactActiveOrientationSummary`

### Primary Actor

The Go transport layer broadcasting authoritative active-client snapshots.

### Pre-conditions

- active clients already receive compact orientation summaries
- current orientation support remains the dominant active serialized cost
- the current fixed active orientation cadence is still the stable baseline

### Trigger

An active snapshot is prepared and includes orientation support for the current cadence tick.

### Success Outcome

- active local play remains responsive
- active orientation support remains usable for the minimap
- each active orientation refresh is smaller than the current compact baseline
- the one-active and two-active baselines drop below the current stable readings

### Failure Or Rejection Cases

- if active minimap orientation becomes too vague to remain useful, the slice failed
- if the slice changes refresh policy again, scope drifted
- if gameplay semantics change, scope is exceeded

## Main Business Rules

1. This is an active-orientation payload-compaction slice, not a cadence slice.
2. Active local detail remains unchanged.
3. Active orientation cadence remains unchanged.
4. Passive observer behavior remains unchanged.
5. The active orientation summary may become coarser than the current compact baseline if minimap usability remains acceptable.

## Minimal Domain Concepts In Scope

- `Active Orientation Summary`
- `Minimap Usability`
- `Active Payload Size`
- `Coarse Whole-World Orientation`

## Bounded Optimization Interpretation

This slice chooses the smallest plausible orientation reduction after the failed event-driven attempt:

- keep the current refresh cadence
- keep the current contract shape recognizable
- make the active orientation summary itself coarser or lighter than the current compact form

## Required Runtime Contract Changes

Prefer none, or only the smallest contract-preserving change.

## Required Ports Or Boundaries

- active snapshot assembly in the Go transport layer
- client minimap rendering of the existing active orientation summary
- transport measurement helpers updated with the new stable baseline

## Build Guidance

- do not revisit refresh cadence in this slice
- target only the active orientation summary payload itself
- remeasure one-active and two-active baselines after the change
- keep active minimap orientation plausibly usable

## Initial Test Plan

### Server or measurement tests

- active transport measurement stays deterministic
- active orientation-support bytes drop below the current stable baseline
- one-active and two-active aggregate bytes drop below the current stable readings

### Contract tests

- current contract validation remains green, with only minimal updates if strictly required

### Integration tests

- active clients still receive responsive local detail
- active minimap orientation remains present and usable

## Scenario Definition

Start the current server and measure the optimized active orientation-summary path.

Scenario steps:

1. build one active snapshot under the new compaction rule
2. remeasure the current active payload breakdown
3. remeasure one-active and two-active active baselines
4. confirm active local detail and minimap orientation still behave as before

## Done Criteria

- active local-detail behavior stays responsive
- active orientation summary payload is smaller than the current compact baseline
- one-active and two-active active baselines drop below the current stable readings
- no gameplay semantics change

## Out Of Scope Follow-Ups

- event-driven active orientation refresh
- passive transport redesign
- delta encoding
- compression
- gameplay changes
