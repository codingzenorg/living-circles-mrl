# Slice: Initial Client Render Pressure Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with minimap, offscreen cues, recent effects, and support panels
- current reduced-cost active transport path
- Go authoritative server with current tick cadence

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but measures whether the next likely pressure has shifted from websocket payload to browser-side rendering and draw cost.

## Discovery Scope

Measure the client render path under the current larger-world viewport model.

This slice should:

- preserve the current transport behavior unchanged
- preserve the current client behavior unchanged
- make draw or render pressure explicit under ordinary active movement
- record whether client-side rendering is now a plausible next bottleneck

This slice does **not** attempt to implement:

- further transport optimization
- rendering optimization
- gameplay changes
- UI redesign

## Why This Slice Next

The transport loop is now in a better state:

- passive observer fanout is much cheaper
- active transport cost is lower than before
- active orientation usability is now explicitly measured

That changes the next uncertainty:

- transport is no longer the only plausible source of roughness
- the viewport client now does substantial drawing work: play stage, minimap, offscreen cues, glows, logs, and support UI
- the next useful step is to measure the client render path before optimizing further

## Use-Case Contract

### Use Case

`MeasureClientRenderPressure`

### Primary Actor

The repository’s MRL loop evaluating the browser client’s current rendering path.

### Pre-conditions

- the viewport client is already feature-rich
- the transport path has already been trimmed and measured
- the next likely bottleneck may now be render cost rather than bytes

### Trigger

A bounded local measurement run records browser-side render pressure during ordinary active movement.

### Success Outcome

- client render pressure is explicit enough to guide the next step
- the next loop can choose between more transport work and client rendering work from evidence
- no behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the result does not clarify render pressure meaningfully, it failed
- if the slice drifts into implementing optimization, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Transport behavior remains unchanged.
3. Client behavior remains unchanged.
4. The result must make browser-side render pressure more explicit than it is now.

## Minimal Domain Concepts In Scope

- `Client Render Pressure`
- `Frame Work`
- `Viewport Draw Cost`
- `Ordinary Movement Window`

## Bounded Measurement Interpretation

This slice chooses the smallest useful render measurement:

- instrument the current draw loop or a narrow client-side render path metric
- record bounded render pressure during ordinary active movement
- keep the metric simple enough to compare later after any client optimization

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client instrumentation only
- implementation artifact updated with the measurement result

## Build Guidance

- do not change gameplay or transport logic
- keep the metric deterministic or at least repeatable enough for local comparison
- prefer one or two clear client render metrics over a broad profiler-style dump

## Initial Test Plan

### Server or measurement tests

- any added client-side instrumentation should degrade safely when measurement is unavailable

### Contract tests

- none beyond the current contract validation, because runtime shape remains unchanged

### Integration tests

- none required unless the measurement crosses a server boundary

## Scenario Definition

Run the current viewport client under ordinary movement and record:

1. one bounded render-pressure metric
2. one bounded movement window
3. the observed result in the implementation artifact

## Done Criteria

- render pressure is more explicit than before
- the next optimization direction can be chosen from evidence
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- transport optimization
- render optimization
- UI redesign
- gameplay changes
