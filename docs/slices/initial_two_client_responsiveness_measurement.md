# Slice: Initial Two-Client Responsiveness Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in ordinary viewport mode
- Go authoritative server under concurrent local clients
- current transport and render instrumentation unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and rendering behavior unchanged, but turns the newly observed two-browser slowdown into an explicit bounded measurement artifact.

## Discovery Scope

Measure the current responsiveness impact of running two ordinary local clients against the same server.

This slice should:

- preserve current runtime behavior unchanged
- preserve transport shape unchanged
- preserve current render instrumentation unchanged
- make concurrent responsiveness more explicit than “it feels slow with a second browser”

This slice does **not** attempt to implement:

- transport redesign
- render redesign
- gameplay changes
- client prediction or interpolation

## Why This Slice Next

The repository now has stronger evidence than before:

- opening a second browser session caused visible slowdown
- the simulation itself felt slower
- player movement also felt slower in ordinary play

That means the remaining pressure is no longer hypothetical. It is now a real runtime concern under a very small concurrent-client count.

The next useful step is not another optimization by intuition. It is to make the two-client responsiveness problem explicit enough to choose the right next mitigation.

## Use-Case Contract

### Use Case

`MeasureTwoClientResponsiveness`

### Primary Actor

The repository’s MRL loop evaluating real concurrent responsiveness.

### Pre-conditions

- the browser and transport already expose bounded instrumentation
- a two-client slowdown was observed manually
- the next optimization should be evidence-backed

### Trigger

A bounded local measurement run exercises the system with two concurrent clients and records responsiveness-relevant signals.

### Success Outcome

- the two-client slowdown is expressed as explicit evidence instead of only anecdote
- the next loop can decide whether the pressure is mainly tick cadence, active fanout, browser rendering, or some combination
- no runtime behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the result still does not clarify concurrent responsiveness meaningfully, it failed
- if the slice drifts into mitigation instead of measurement, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Default runtime behavior remains unchanged.
3. The result should make two-client responsiveness more explicit than the current anecdotal observation.
4. The result should be simple enough to justify the next mitigation slice.

## Minimal Domain Concepts In Scope

- `Concurrent Responsiveness`
- `Two-Client Slowdown`
- `Tick Cadence Stability`
- `Per-Client Update Stability`
- `Observed Play Slowdown`

## Bounded Measurement Interpretation

This slice chooses the smallest useful concurrent-responsiveness measurement:

- use two concurrent clients against the current server
- record bounded responsiveness-relevant signals such as snapshot cadence stability, aggregate/per-client throughput, and any available browser-side render read
- summarize the likely dominant pressure source in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- local concurrent-client measurement only
- implementation artifact updated with the responsiveness findings

## Build Guidance

- do not change gameplay, transport, or current visuals
- prefer a concise bounded result over a large profiling system
- reuse the current server and client instrumentation where possible

## Initial Test Plan

### Server or measurement tests

- add or extend deterministic measurement helpers only if needed

### Contract tests

- unchanged beyond current validation

### Integration tests

- only if a new bounded measurement helper crosses the current transport boundary

## Scenario Definition

Run the current demo with two concurrent clients and record:

1. server-side snapshot/tick stability
2. aggregate and per-client update cost
3. any available browser-side render read
4. a concise conclusion about the likeliest current source of the slowdown

## Done Criteria

- the two-client slowdown is more explicit than the current anecdote
- the next mitigation slice can be chosen from evidence
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- transport redesign
- render redesign
- gameplay changes
- prediction or interpolation
