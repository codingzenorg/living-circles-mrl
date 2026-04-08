# Slice: Initial Multi-Client Transport Pressure Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser-like websocket clients simulated locally for measurement
- Go authoritative server and current optimized transport path
- shared runtime contract unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice does not change gameplay, transport shape, or cadence policy. It adds a deterministic measurement harness that observes what the current optimized server actually costs when several clients are connected at once.

## Discovery Scope

Measure stronger runtime pressure before deciding whether more transport optimization is justified.

This slice should:

- connect several local websocket clients to one server
- exercise the current transport path under deterministic conditions
- record aggregate and per-client transport cost
- record whether the server can still keep up with its intended tick cadence

This slice does **not** attempt to implement:

- protocol redesign
- compression
- binary transport
- simulation throttling
- gameplay changes

## Why This Slice Next

Recent slices made the protocol much less chatty, but the current pressure is now evidentiary rather than structural. The project has enough optimization work in place that another redesign should be motivated by measured runtime pressure, not by intuition alone.

The new pressure is:

- determine whether the current optimized transport path is already good enough for several simultaneous clients
- determine whether remaining pressure comes more from bytes sent or from tick/broadcast timing
- create a repeatable baseline so later transport changes can be judged against something stronger than one-client measurements

This is the narrowest next step because it:

- leaves current protocol behavior untouched
- turns current lag concerns into explicit evidence
- can guide whether the next loop should stay on transport or return to ecological/play work

## Use-Case Contract

### Use Case

`MeasureMultiClientTransportPressure`

### Primary Actor

The repository maintainer evaluating the current websocket transport path under deterministic local load.

### Pre-conditions

- the server already emits optimized viewport-mode transport snapshots
- measurement helpers already exist for single-snapshot payload cost
- the current world and tick cadence are deterministic enough for repeatable local runs

### Trigger

A deterministic measurement run starts one server and multiple local websocket clients, then records transport and tick-pressure observations over a bounded window.

### Success Outcome

- the repository has a deterministic multi-client transport baseline
- aggregate bytes/sec and per-client bytes/sec are explicit
- the run reports whether the server still stays near its intended cadence
- later optimization slices can be justified by evidence instead of intuition

### Failure Or Rejection Cases

- if the run is not deterministic enough to compare between revisions, the slice failed
- if the slice changes protocol behavior or gameplay, scope drifted
- if the measurement couples itself to browser rendering or external network conditions, scope is exceeded

## Main Business Rules

1. This is a measurement slice, not a gameplay slice.
2. The current protocol shape and cadence remain unchanged.
3. Multiple clients should be simulated deterministically.
4. The output must separate aggregate transport cost from per-client cost.
5. The output must say something about server cadence pressure, not only payload size.

## Minimal Domain Concepts In Scope

- `Per-Client Transport Cost`
- `Aggregate Transport Cost`
- `Observed Tick Pressure`
- `Deterministic Multi-Client Baseline`

## Bounded Measurement Interpretation

This slice chooses the smallest useful pressure harness:

- one local server
- a bounded number of deterministic websocket clients
- a short deterministic measurement window
- explicit metrics such as:
  - bytes per client
  - aggregate bytes
  - approximate bytes/sec
  - observed message count
  - observed tick-delay or broadcast-lag signal

It avoids moving into browser profiling, distributed load testing, or full production realism.

## Required Runtime Contract Changes

None.

The slice should measure the current transport boundary as it already exists.

## Required Ports Or Boundaries

- websocket transport test harness
- deterministic client behavior for the measurement run
- implementation artifact or measurement artifact recording the results

## Build Guidance

- keep the client count bounded and explicit, for example `4`, `8`, or another fixed small set
- use deterministic client roles such as idle observer and moving client if multiple behaviors are included
- record both per-client and aggregate cost
- record whether the server still maintains the intended cadence closely enough to trust the transport path
- keep the result easy to compare with later transport slices

## Initial Test Plan

### Server or measurement tests

- multi-client measurement remains deterministic across repeated runs
- aggregate bytes exceed one-client bytes in the expected direction
- reported per-client and aggregate numbers remain coherent
- the measurement captures a cadence-pressure signal or explicit lack of one

### Contract tests

- none beyond the current contract validation, because protocol shape is unchanged

### Integration tests

- deterministic websocket clients can connect simultaneously and receive snapshots over the bounded measurement window

## Scenario Definition

Start one local server and a bounded set of local websocket clients.

Scenario steps:

1. connect several clients
2. optionally apply one deterministic movement pattern to a subset of them
3. collect snapshot counts, bytes, and timing observations over a short fixed window
4. compute per-client and aggregate transport estimates
5. record whether cadence pressure is visible

## Done Criteria

- the repository has a deterministic multi-client transport-pressure baseline
- per-client and aggregate transport cost are explicit
- cadence pressure is explicitly reported
- no gameplay semantics or protocol shape changed

## Out Of Scope Follow-Ups

- browser render profiling
- internet-latency simulation
- compression or protocol redesign
- multi-process load generation
