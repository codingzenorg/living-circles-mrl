# Slice: Initial Overlay Render Subcomponent Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure, render-family, and world-subcomponent measurement
- current overlay draw path including afterglow, crowding, food glows, offscreen cues, lineage links, and heading/intent cues
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and visual behavior unchanged, but makes the broad `overlay` render family more actionable by splitting it into a few meaningful draw subfamilies.

## Discovery Scope

Break down the current canvas `overlay` draw family into a small number of useful subcomponents.

This slice should:

- preserve current runtime behavior by default
- preserve transport behavior unchanged
- keep the existing top-level render-family breakdown
- identify whether glow fields, offscreen awareness, recent effects, or lineage/intent cues dominate the current `overlay` draw bucket

This slice does **not** attempt to implement:

- render optimization
- gameplay changes
- transport changes
- UI redesign

## Why This Slice Next

The client now exposes:

- total render pressure
- major render-family breakdown
- world-subcomponent breakdown
- a bounded label optimization

That improves the picture, but one broad family remains unresolved:

- `overlay` still aggregates several different visual systems
- if overlays are still materially expensive, the current measurement is too coarse to choose the next client optimization confidently
- the next useful step is to split only that family, not to jump into another blind tweak

## Use-Case Contract

### Use Case

`MeasureOverlayRenderSubcomponents`

### Primary Actor

The browser client measuring its own canvas overlay-draw path.

### Pre-conditions

- live total render pressure already exists
- major render-family measurement already exists
- the next client optimization should remain evidence-backed

### Trigger

A bounded local render measurement run compares a few overlay-draw subfamilies under the current viewport experience.

### Success Outcome

- the dominant subfamilies inside the `overlay` draw bucket are explicit
- the next client optimization can target the right overlay family, if that remains the main pressure
- no runtime behavior changes in the default user path

### Failure Or Rejection Cases

- if the slice changes runtime behavior by default, it failed
- if the measurement becomes too granular to guide the next step, scope drifted
- if the slice drifts into optimization, it failed

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Default runtime behavior remains unchanged.
3. The measurement should compare only a few meaningful overlay families.
4. The result should stay simple enough to support one clear next render slice.

## Minimal Domain Concepts In Scope

- `Overlay Render Subcomponent`
- `Glow Field Cost`
- `Offscreen Awareness Cost`
- `Recent Effect Cost`
- `Lineage And Intent Cue Cost`

## Bounded Measurement Interpretation

This slice chooses the smallest useful subdivision of `overlay` draw cost:

- keep the current top-level families intact
- split `overlay` into a few major subfamilies such as ambient glows, offscreen awareness, recent effects, and lineage/intent/heading cues
- record the dominant subfamilies in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client instrumentation only
- implementation artifact updated with the overlay-render findings

## Build Guidance

- keep default behavior unchanged
- instrument only a few major overlay sections
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

1. full overlay baseline
2. glow-field timing contribution
3. offscreen-awareness timing contribution
4. recent-effect timing contribution
5. lineage and intent timing contribution

Record the comparative result.

## Done Criteria

- dominant overlay-render subfamilies are explicit
- the next client render optimization target is evidence-backed
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- render optimization
- transport optimization
- gameplay changes
- visual redesign
