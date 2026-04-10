# Slice: Initial Two Active Tick Broadcast Pressure Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- Go authoritative server tick loop under two active browsers
- websocket broadcast fanout under concurrent active local-detail transport
- current two-client responsiveness measurement path

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged. It only measures whether the remaining two-browser slowdown is mainly server tick/broadcast pressure under two active clients.

## Discovery Scope

Turn the remaining active-path slowdown into explicit evidence before choosing another mitigation.

This slice should:

- preserve gameplay rules unchanged
- preserve transport shape and cadence unchanged
- keep the new idle-intent suppression baseline intact
- measure server-side timing pressure under two active browsers
- distinguish snapshot-byte cost from tick/broadcast timing cost

This slice does **not** attempt to implement:

- transport redesign
- client render redesign
- prediction or interpolation
- gameplay changes

## Why This Slice Next

The latest reassessment closed one important ambiguity:

- one active browser plus one idle browser now measures `12276` aggregate bytes over `300ms`
- two active browsers measure `17884` aggregate bytes over the same window
- the second idle browser now truly reaches the passive path

That means neutral idle intent churn is no longer the main explanation for the observed slowdown. The remaining pressure is now more plausibly the active path itself. But that still leaves an important uncertainty:

- is the remaining problem mostly transport bytes?
- or is it the server tick plus websocket broadcast work under concurrent active clients?

The next useful step is therefore not another optimization guess. It is to measure tick/broadcast pressure explicitly for the two-active-client case.

## Use-Case Contract

### Use Case

`MeasureTwoActiveTickBroadcastPressure`

### Primary Actor

The authoritative server under concurrent active-browser load.

### Pre-conditions

- idle-intent suppression is already built
- post-fix two-client responsiveness reassessment is already recorded
- two-active-client transport is the clearest remaining pressure path

### Trigger

Two active browsers drive the server concurrently.

### Success Outcome

- the repo has explicit timing evidence for the two-active-client case
- the next mitigation can target the real dominant pressure
- the repo can distinguish byte cost from tick/broadcast timing pressure

### Failure Or Rejection Cases

- if the slice drifts into implementation changes, scope drifted
- if it records only bytes and not timing pressure, it is incomplete

## Main Business Rules

1. Runtime behavior remains unchanged in this slice.
2. The current post-idle-intent active-path state is the baseline under measurement.
3. Measurement should focus on at least:
   - two active clients
   - bounded tick/broadcast timing pressure
   - relation between observed tick cadence and expected cadence
4. The result should clarify whether the next mitigation should target:
   - active payload size
   - active broadcast scheduling
   - server tick workload

## Minimal Domain Concepts In Scope

- `Two Active Browsers`
- `Tick Pressure`
- `Broadcast Pressure`
- `Expected Tick Cadence`
- `Active Path Baseline`

## Bounded Implementation Interpretation

This slice chooses the smallest useful next step:

- extend the current deterministic measurement path with timing-oriented evidence
- avoid changing transport or gameplay behavior
- record the result in implementation evidence

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server-side measurement helpers
- runtime evidence path
- implementation artifact updated with the measured result

## Build Guidance

- prefer deterministic bounded timing measures over subjective interpretation
- reuse existing multi-client measurement helpers where possible
- keep the evidence legible enough to guide the next mitigation slice directly

## Initial Test Plan

### Server or measurement tests

- prove the two-active timing measurement is deterministic
- prove it reports bounded timing pressure against the expected tick cadence

### Contract tests

- unchanged beyond current validation

### Integration tests

- reuse current transport integration unless a focused timing-sensitive proof is needed

## Scenario Definition

Run the deterministic two-active-client harness and record:

1. aggregate snapshots and bytes
2. expected tick cadence
3. bounded timing or inter-snapshot pressure signal
4. the resulting read on whether active slowdown is primarily timing pressure or payload pressure

## Done Criteria

- two-active-client timing pressure is explicitly recorded
- the next mitigation target is clearer than “active path feels slower”
- the repo avoids optimizing the wrong part of the active path

## Out Of Scope Follow-Ups

- transport redesign
- render redesign
- gameplay changes
- prediction or interpolation
