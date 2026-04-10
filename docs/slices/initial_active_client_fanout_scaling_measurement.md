# Slice: Initial Active Client Fanout Scaling Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using `active_local_detail`
- passive observer path already optimized and measured
- Go authoritative server with the current transport policies

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and protocol shape unchanged, but measures how the remaining high-cadence active path scales once several clients are simultaneously steering instead of leaving the next optimization target to intuition.

## Discovery Scope

Measure the current transport pressure of several active clients.

This slice should:

- preserve the current transport protocol and cadence policies
- preserve the current passive observer path unchanged
- run deterministic multi-client measurements with several moving clients
- make aggregate bytes, per-client bytes, and cadence pressure explicit for active-client fanout

This slice does **not** attempt to implement:

- compression
- delta encoding
- new subscriptions
- gameplay changes
- another transport optimization yet

## Why This Slice Next

The current transport loop now has a clear shape:

- passive observer fanout has been reduced materially
- single active clients still use the full local-detail path
- mixed active/passive runs now show the active path dominates the remaining bytes

That makes the next pressure explicit:

- the main remaining transport question is no longer calm passive observers
- it is whether multiple active clients scale acceptably on the current local-detail path
- the next useful step is measurement, not another shape change

## Use-Case Contract

### Use Case

`MeasureActiveClientFanoutScaling`

### Primary Actor

The Go transport layer measuring its own current deterministic websocket output under active-client load.

### Pre-conditions

- the server already supports deterministic multi-client measurement
- active clients already receive the current full local-detail path
- passive observers already have an optimized lower-cost path

### Trigger

A bounded local measurement run is started with several simultaneous moving clients.

### Success Outcome

- aggregate bytes/sec for multiple active clients is explicit
- per-client bytes/sec stays interpretable
- max inter-snapshot gap under active fanout is explicit
- future transport work can be motivated from measured active pressure rather than assumption

### Failure Or Rejection Cases

- if gameplay or transport shape changes, the slice failed
- if the measurement cannot run deterministically, the slice failed
- if the slice drifts into implementing optimization instead of measurement, scope drifted

## Main Business Rules

1. This is a measurement slice, not a gameplay slice.
2. Active-client transport behavior remains unchanged.
3. Passive observer behavior remains unchanged.
4. The measurement must compare at least several explicit active-client counts.
5. The reported metrics must include aggregate bytes, per-client bytes, and cadence pressure.

## Minimal Domain Concepts In Scope

- `Active Client`
- `Active Fanout Scaling`
- `Aggregate Transport Pressure`
- `Per-Client Pressure`
- `Inter-Snapshot Gap`

## Bounded Measurement Interpretation

This slice chooses the smallest useful active-path measurement:

- run deterministic active-client fanout ladders such as `1 / 2 / 4`
- keep each client moving under the same bounded pattern
- record aggregate bytes, per-client bytes, snapshot counts, and max inter-snapshot gap

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- deterministic transport measurement harness
- transport implementation artifact with the measured active ladder

## Build Guidance

- reuse the existing multi-client measurement harness
- keep the moving-client pattern deterministic
- do not modify protocol shape or cadence policy in this slice
- make the final measurement easy to compare with the existing passive and mixed baselines

## Initial Test Plan

### Server or measurement tests

- active fanout measurement is deterministic across repeated runs
- aggregate bytes increase with active-client count
- per-client bytes remain positive and interpretable
- max inter-snapshot gap stays bounded under the measured active ladder

### Contract tests

- none beyond the current contract validation, because shape should remain unchanged

### Integration tests

- none required beyond current websocket behavior, unless the measurement harness needs new integration proof

## Scenario Definition

Start one local server and run deterministic multi-client transport measurement with several simultaneously moving clients.

Scenario steps:

1. measure one moving client
2. measure a small active fanout ladder such as `2` and `4` moving clients
3. compare aggregate bytes, per-client bytes, and max inter-snapshot gap across the ladder
4. compare those results with the current passive and mixed baselines

## Done Criteria

- active-client fanout scaling is explicitly measured
- the current main remaining transport pressure is clearer
- no gameplay or protocol shape changed

## Out Of Scope Follow-Ups

- compression
- delta encoding
- active-path optimization
- passive-path redesign
- gameplay changes
