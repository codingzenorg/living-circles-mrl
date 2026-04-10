# Slice: Initial Active Orientation Support Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using `active_local_detail`
- active payload component measurement already recorded
- Go authoritative server with the current local-detail transport

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and active local circle detail unchanged, but reduces the cost of active orientation support now that the measurements show it is the dominant serialized active payload family.

## Discovery Scope

Lower active transport cost by targeting the dominant active payload family.

This slice should:

- preserve the current active-client local-detail path
- preserve the current passive observer path unchanged
- reduce active orientation-support bytes
- keep the current runtime contract shape recognizable

This slice does **not** attempt to implement:

- compression
- delta encoding
- gameplay changes
- active local-detail removal
- passive transport changes

## Why This Slice Next

The latest active payload breakdown now shows:

- full active payload: `3333` bytes
- without local autonomous detail: `2946`
- without local food detail: `3105`
- without player detail: `3097`
- without orientation support: `1176`

That makes the next pressure explicit:

- orientation support is the dominant active serialized cost
- optimizing player, food, or local autonomous detail first would likely target the wrong area
- the next bounded step is to reduce active orientation-support traffic without harming active responsiveness

## Use-Case Contract

### Use Case

`ReduceActiveOrientationSupportCost`

### Primary Actor

The Go transport layer broadcasting authoritative active-client snapshots.

### Pre-conditions

- active clients already receive the current `active_local_detail` path
- active payload component measurement already identified orientation support as the dominant active serialized cost
- passive observer transport already has its own optimized path

### Trigger

An active snapshot is prepared for transport and the server decides how much whole-world orientation support to include on that tick.

### Success Outcome

- active local play remains responsive
- active snapshots carry less orientation-support payload on average
- the runtime contract remains recognizable
- measured active bytes/sec drop below the current active baseline

### Failure Or Rejection Cases

- if local active responsiveness changes, the slice failed
- if the slice removes necessary orientation support entirely, the slice failed
- if the slice turns into a generalized transport redesign, scope drifted

## Main Business Rules

1. This is an active-transport optimization slice, not a gameplay slice.
2. Active local detail remains high-cadence.
3. Passive observer behavior remains unchanged.
4. Orientation support for active clients may be reduced in cadence or compaction, but must remain available.
5. The new active transport cost should be remeasured against the current active baseline.

## Minimal Domain Concepts In Scope

- `Active Orientation Support`
- `Active Local Detail`
- `Orientation Refresh`
- `Active Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest evidence-backed optimization:

- keep the active snapshot shape recognizable
- reduce the frequency or payload weight of active orientation support only
- leave player detail, local autonomous detail, and local food detail unchanged

## Required Runtime Contract Changes

Prefer none, or the smallest possible contract-preserving adjustment.

## Required Ports Or Boundaries

- active snapshot assembly in the Go transport layer
- client handling of any active orientation staleness if needed
- transport measurement helpers updated to record the new active baseline

## Build Guidance

- optimize only the active orientation-support portion
- preserve active movement readability and local responsiveness
- keep the transport contract recognizable
- remeasure both single active and active fanout baselines

## Initial Test Plan

### Server or measurement tests

- active transport measurement stays deterministic
- active payload bytes drop below the current baseline
- active fanout scaling is remeasured after the change
- active cadence pressure remains bounded

### Contract tests

- current contract validation remains green, with only minimal updates if strictly required

### Integration tests

- active client still receives responsive local detail
- active client still retains usable orientation support across the optimized path

## Scenario Definition

Start one local server and measure the optimized active path.

Scenario steps:

1. build one active snapshot under the optimized orientation policy
2. remeasure single active-client bytes/sec
3. remeasure the `1 / 2 / 4` active fanout ladder
4. confirm local active play still behaves as before

## Done Criteria

- active active-local-detail behavior stays responsive
- active transport cost drops from the current baseline
- the dominant active cost family is reduced directly
- no gameplay semantics changed

## Out Of Scope Follow-Ups

- compression
- delta encoding
- passive transport redesign
- generalized subscriptions
- gameplay changes
