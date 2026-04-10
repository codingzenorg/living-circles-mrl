# Slice: Initial World Render Subcomponent Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure and render-family breakdown
- current play-stage world rendering path inside the canvas
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and visual behavior unchanged, but makes the broad `world` render family more actionable by splitting it into a few meaningful draw subfamilies.

## Discovery Scope

Break down the current canvas `world` draw family into a small number of useful subcomponents.

This slice should:

- preserve current runtime behavior by default
- preserve transport behavior unchanged
- keep the existing top-level render-family breakdown
- identify whether circles, foods, labels, or world framing dominate the current `world` draw bucket

This slice does **not** attempt to implement:

- render optimization
- gameplay changes
- transport changes
- UI redesign

## Why This Slice Next

The client now exposes two useful levels of render visibility:

- a rolling total render-pressure metric
- a rolling major-family breakdown across `world`, `overlay`, `support`, and `minimap`

That is better than before, but one meaningful gap remains:

- the `world` family still aggregates several very different canvas costs
- if the next client optimization ends up targeting world drawing, the current measurement is still too coarse
- the next useful step is to split only that broad bucket, not to jump directly into an optimization

## Use-Case Contract

### Use Case

`MeasureWorldRenderSubcomponents`

### Primary Actor

The browser client measuring its own canvas world-draw path.

### Pre-conditions

- live total render pressure already exists
- major render-family measurement already exists
- the next client optimization should still be evidence-backed

### Trigger

A bounded local render measurement run compares a few world-draw subfamilies under the current viewport experience.

### Success Outcome

- the dominant subfamilies inside the `world` draw bucket are explicit
- the next client optimization can target the right world-draw family, if that is still the main pressure
- no runtime behavior changes in the default user path

### Failure Or Rejection Cases

- if the slice changes runtime behavior by default, it failed
- if the measurement is too granular to guide the next step, scope drifted
- if the slice drifts into optimization, it failed

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Default runtime behavior remains unchanged.
3. The measurement should compare only a few meaningful world-draw families.
4. The result should stay simple enough to support one clear next render slice.

## Minimal Domain Concepts In Scope

- `World Render Subcomponent`
- `Circle Draw Cost`
- `Food Draw Cost`
- `Label Draw Cost`
- `World Frame Cost`

## Bounded Measurement Interpretation

This slice chooses the smallest useful subdivision of `world` draw cost:

- keep the current top-level families intact
- split `world` into a few major subfamilies such as frame/background, foods, circles/children, and labels
- record the dominant subfamilies in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client instrumentation only
- implementation artifact updated with the world-render findings

## Build Guidance

- keep default behavior unchanged
- instrument only a few major world-draw sections
- avoid profiler-style noise
- prefer a result that clearly informs one next optimization slice

## Initial Test Plan

### Server or measurement tests

- any added client-side instrumentation must degrade safely when unavailable

### Contract tests

- none beyond the current contract validation, because no protocol shape changes are expected

### Integration tests

- none required unless the instrumentation crosses a server boundary

## Scenario Definition

Run the current viewport client and compare:

1. full world draw baseline
2. world framing/background timing contribution
3. food draw timing contribution
4. circle-and-child draw timing contribution
5. label timing contribution

Record the comparative result.

## Done Criteria

- dominant world-render subfamilies are explicit
- the next client render optimization target is evidence-backed
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- render optimization
- transport optimization
- gameplay changes
- visual redesign
