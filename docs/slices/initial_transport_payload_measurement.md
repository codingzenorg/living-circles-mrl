# Slice: Initial Transport Payload Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract measured as it exists today, without protocol thinning in this slice

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for snapshot production and the client remaining a consumer of those snapshots.

This slice does not optimize the protocol yet. It makes current transport cost explicit so later protocol evolution can be driven by measured pressure rather than guesswork.

## Discovery Scope

Establish the smallest deterministic measurement slice that exposes current snapshot cost under representative current worlds:

- measure snapshot payload size under current full-snapshot broadcasting
- measure broadcast rate and resulting bytes per second
- include the default expanded world and at least one larger-world scenario
- keep the existing message schema and runtime behavior unchanged

This slice does **not** attempt to implement:

- delta snapshots
- viewport-based interest culling
- compression
- message batching
- protocol versioning
- client rendering optimization

## Why This Slice Next

The current protocol remains full and chatty. That was acceptable while the model was still small and rapidly changing, but the world is now:

- larger
- denser
- more inspectable
- more dependent on continuous authoritative snapshots

There is now credible transport pressure, but optimization should not begin from intuition alone. The next pressure is to make the current cost explicit:

- how large are ordinary snapshots?
- how much do they grow with expanded population and map scale?
- what is the current bytes-per-second cost per client?

This slice is the narrowest next step because it:

- produces evidence without prematurely freezing a protocol strategy
- keeps current semantics and contract shape intact
- creates a solid baseline for later optimization slices

## Use-Case Contract

### Use Case

`MeasureCurrentSnapshotTransportCost`

### Primary Actor

The repository itself during build-time validation and local evaluation.

### Pre-conditions

- the server already emits authoritative world snapshots
- the current protocol is explicit and deterministic enough to measure directly
- the default expanded world and larger-world scenarios already exist

### Trigger

A measurement runner, focused test, or deterministic artifact captures snapshot sizes and broadcast cost under current protocol behavior.

### Success Outcome

- the repository gains explicit baseline numbers for:
  - snapshot byte size
  - snapshots per second or tick cadence
  - approximate bytes per second per client
- results are reproducible enough to guide later optimization slices
- no gameplay semantics change

### Failure Or Rejection Cases

- if the slice starts thinning or changing the protocol, scope is exceeded
- if the measurements are too ad hoc or not reproducible, they will not guide later refinement well
- if only one tiny world is measured, the pressure remains underdescribed

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. The current protocol should be measured as-is.
3. The measurement should include at least:
   - the current default expanded world
   - a larger-world scenario relevant to the new map scale direction
4. The output should be reproducible enough to compare later optimization work against it.
5. The slice should avoid inventing production-like monitoring infrastructure unless the lab actually needs it.

## Minimal Domain Concepts In Scope

- `World Snapshot`
- `Payload Size`
- `Broadcast Rate`
- `Per-Client Transport Cost`

## Bounded Measurement Interpretation

This slice chooses the smallest useful interpretation of protocol measurement:

- capture serialized snapshot size directly from current authoritative snapshots
- combine that with known tick cadence to estimate per-client transport load
- record the results in repository artifacts or deterministic tests

This avoids speculative transport redesign while still making the protocol cost visible.

## Required Runtime Contract Changes

None by default.

The point of the slice is to measure the current contract without changing it.

## Required Ports Or Boundaries

- server-side snapshot serialization or transport path where payload size can be measured
- deterministic test or runner support for representative world states
- implementation or evaluation artifact that records the resulting baseline numbers

## Build Guidance

- prefer a focused deterministic measurement helper or test over ad hoc manual observation
- measure serialized payload size using the actual snapshot format
- capture at least one idle/default-expanded snapshot and one larger-world snapshot
- keep the output easy to compare against future optimization slices

## Initial Test Plan

### Server tests

- a deterministic measurement helper reports snapshot byte size for representative world states
- the current default expanded world and a larger-world scenario can both be measured repeatably

### Contract tests

- the current schema remains unchanged

### Integration tests

- optional only if needed to confirm measured transport shape matches the real websocket payload path

## Scenario Definition

Start a local server and run the measurement flow against representative world states.

Scenario steps:

1. measure serialized snapshot size for the current default expanded world
2. measure serialized snapshot size for a larger-world scenario
3. combine those with current broadcast cadence into an approximate per-client transport cost
4. record the result for later optimization slices

## Done Criteria

- current snapshot payload size is explicit
- current per-client transport cost is explicit enough to guide future protocol work
- at least two representative world states are measured
- no protocol or gameplay behavior changes

## Out Of Scope Follow-Ups

- delta snapshot design
- viewport-based protocol thinning
- compression
- protocol versioning
- rollout or compatibility strategy
