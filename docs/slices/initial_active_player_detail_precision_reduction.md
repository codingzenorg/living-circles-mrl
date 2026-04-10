# Slice: Initial Active Player Detail Precision Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active-client player detail serialization
- current active local-detail transport after autonomous cadence reduction
- browser player rendering under authoritative snapshots

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and cadence unchanged. It only reduces the wire precision of player-centric active detail to display-sufficient values where that precision is still more exact than the current viewport visibly needs.

## Discovery Scope

Trim one more active payload family without touching active cadence, passive behavior, or the meaning of the current contract.

This slice should:

- preserve gameplay rules unchanged
- preserve active cadence unchanged
- preserve passive transport behavior unchanged
- keep the current `autonomous_fresh` active-path behavior unchanged
- reduce only the serialized precision of player-centric active detail

This slice does **not** attempt to implement:

- a new transport shape
- cadence changes
- compression
- gameplay changes
- client prediction or interpolation

## Why This Slice Next

The last active-path slice removed one bounded repeated cost:

- under the current bounded `300ms` window, the two-active baseline dropped from `17884` bytes to `17008`
- timing stayed effectively flat near one tick
- the single-client active path stayed effectively intact

That is a good reduction, but the active path is still the clearest remaining transport pressure under concurrent play. The next clean step is not another cadence change by default. It is to trim one more payload family while leaving the refreshed active shape intact.

Player detail is the cleanest next target:

- it remains present on every active snapshot
- it still carries high-precision numeric state for positions and attached-child layout
- the current viewport rendering likely does not need every serialized decimal carried across the wire

## Use-Case Contract

### Use Case

`ReduceActivePlayerDetailPrecision`

### Primary Actor

The authoritative server serializing active player snapshots for browser clients.

### Pre-conditions

- active local autonomous cadence reduction is already built
- active player cadence and semantics are currently stable
- active payload fanout remains the clearest remaining pressure

### Trigger

An active snapshot is serialized for a browser client.

### Success Outcome

- active player behavior and readability remain intact
- player-centric payload cost drops modestly without shape changes
- the current two-active baseline drops below the latest `17008` byte reading if the precision cut is meaningful

### Failure Or Rejection Cases

- if player readability or responsiveness degrades, the slice failed
- if the slice expands into a transport redesign, scope drifted
- if the precision cut is too small to matter and cannot be justified, the slice failed

## Main Business Rules

1. Gameplay and active transport cadence remain unchanged.
2. The server remains authoritative.
3. Only serialized precision changes; internal simulation precision stays unchanged.
4. Player-centric detail must stay display-sufficient for the current viewport and overlays.

## Minimal Domain Concepts In Scope

- `Active Player Detail`
- `Serialized Precision`
- `Display-Sufficient Precision`
- `Active Payload Fanout`

## Bounded Implementation Interpretation

This slice chooses the smallest useful mitigation:

- keep the current active snapshot shape
- reduce only the serialized precision of player-centric active detail where safe
- measure the resulting effect against the current active baseline

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server active snapshot builder
- measurement helpers
- implementation artifact updated with the resulting baseline

## Build Guidance

- prefer explicit rounding rules to ad hoc formatting changes
- preserve visual readability and active play feel
- validate against the current one-active and two-active baselines

## Initial Test Plan

### Server or measurement tests

- prove player detail still serializes deterministically
- prove the active baseline stays below the current post-autonomous-cadence reading if the reduction is meaningful

### Contract tests

- unchanged beyond current validation

### Integration tests

- update only if current player rendering assumptions require a focused proof

## Scenario Definition

Run the current deterministic active harness with:

1. one active client baseline
2. two active clients baseline
3. reduced serialized precision for player-centric active detail
4. compare the resulting aggregate payload against the current post-autonomous-cadence baseline

## Done Criteria

- player-centric active precision is reduced without semantic drift
- active readability remains intact
- the active baseline is measurably lower or the repo learns that this target is not worth pursuing

## Out Of Scope Follow-Ups

- transport redesign
- compression
- gameplay changes
- passive transport changes
