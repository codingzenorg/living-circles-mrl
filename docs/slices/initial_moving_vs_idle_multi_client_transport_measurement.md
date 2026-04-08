# Slice: Initial Moving Versus Idle Multi-Client Transport Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- local websocket clients acting as deterministic observers or movers
- Go authoritative server with the current optimized transport path
- shared runtime contract unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps the current transport shape unchanged and refines the new multi-client measurement harness so it can distinguish between a mostly idle observer load and a load where one or more clients continuously send movement intent.

## Discovery Scope

Measure whether the remaining runtime pressure is mainly driven by:

- always-on local circle motion
- orientation/local-food refresh behavior under movement
- or simply the cost of serving several passive observers

This slice should:

- keep gameplay and transport unchanged
- reuse the current multi-client measurement harness
- compare idle multi-client load against a bounded moving-client load
- report how movement changes per-client and aggregate transport pressure

This slice does **not** attempt to implement:

- protocol redesign
- compression
- prediction
- transport throttling
- gameplay changes

## Why This Slice Next

The repository now has a deterministic multi-client baseline:

- `4` local clients
- explicit aggregate bytes/sec
- explicit per-client bytes/sec
- explicit inter-snapshot gap signal

That baseline is useful, but still too coarse to guide the next optimization confidently. It tells us what several clients cost in aggregate, but not whether the next pressure comes from movement-heavy local detail or from the mere presence of extra clients.

The new pressure is:

- separate passive-observer cost from active-play cost
- find out whether movement materially increases transport pressure under the current optimized design
- identify whether the next optimization target should be local moving-circle detail or something else

This is the narrowest next step because it:

- reuses the new measurement harness
- adds evidence instead of another speculative transport redesign
- keeps the result directly comparable to the existing multi-client baseline

## Use-Case Contract

### Use Case

`CompareIdleAndMovingMultiClientTransportPressure`

### Primary Actor

The repository maintainer evaluating whether current transport pressure is dominated by active play or by passive fanout.

### Pre-conditions

- a deterministic multi-client measurement harness already exists
- the current optimized transport path is already active
- the server still uses authoritative movement intents and snapshots

### Trigger

A bounded measurement run executes two deterministic scenarios:

1. several idle clients
2. the same multi-client setup with one or more clients continuously sending deterministic movement intent

### Success Outcome

- the repository has an explicit idle-versus-moving transport comparison
- the result shows whether active movement materially changes bytes/sec or cadence pressure
- the next optimization target becomes easier to justify

### Failure Or Rejection Cases

- if the run is not deterministic enough to compare meaningfully, the slice failed
- if the slice changes transport behavior instead of only measuring it, scope drifted
- if movement roles are too complex to explain or reproduce, scope exceeded

## Main Business Rules

1. This is a measurement slice, not a gameplay slice.
2. The current protocol shape and cadence remain unchanged.
3. Idle and moving scenarios must be deterministic and directly comparable.
4. The reported output must make the difference between passive and active transport pressure explicit.
5. The slice should help identify whether the next meaningful cost is moving local circles.

## Minimal Domain Concepts In Scope

- `Idle Client Load`
- `Moving Client Load`
- `Incremental Movement Pressure`
- `Deterministic Transport Comparison`

## Bounded Measurement Interpretation

This slice chooses the smallest useful extension of the existing harness:

- keep the same number of clients and same bounded window
- add one deterministic movement role
- compare aggregate bytes/sec, per-client bytes/sec, and gap signals between the two runs

It avoids moving into real browser profiling or complex gameplay scripting.

## Required Runtime Contract Changes

None.

## Required Ports Or Boundaries

- websocket measurement harness
- deterministic movement-intent sender for selected clients
- implementation artifact or measurement artifact recording both scenarios

## Build Guidance

- keep the client count fixed and explicit
- keep the movement intent simple and deterministic, such as constant directional input at the current sender cadence
- report the delta between idle and moving cases, not only the raw numbers
- keep the output easy to compare with future transport slices

## Initial Test Plan

### Server or measurement tests

- idle-versus-moving measurement remains deterministic enough to compare between repeated runs
- moving load is measurable as distinct from idle load
- reported aggregate and per-client figures remain coherent

### Contract tests

- none beyond the current contract validation

### Integration tests

- deterministic moving clients can connect, send movement intent, and keep receiving snapshots during the bounded measurement window

## Scenario Definition

Start one local server and a bounded set of local websocket clients.

Scenario steps:

1. measure the current idle multi-client case
2. measure the same case with one or more deterministic moving clients
3. compare aggregate bytes/sec, per-client bytes/sec, and inter-snapshot gap
4. record whether movement creates a materially stronger pressure signal

## Done Criteria

- the repository has an explicit idle-versus-moving multi-client transport comparison
- the result helps identify whether movement-heavy local detail is now the dominant remaining cost
- no gameplay semantics or protocol shape changed

## Out Of Scope Follow-Ups

- browser render profiling
- compression or protocol redesign
- movement prediction
- adaptive transport throttling
