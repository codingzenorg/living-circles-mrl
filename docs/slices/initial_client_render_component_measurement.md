# Slice: Initial Client Render Component Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure metric
- current play-stage rendering stack including minimap, offscreen cues, glows, afterglow, and support UI
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and visual behavior unchanged, but makes client render cost more actionable by identifying which major render families dominate draw work.

## Discovery Scope

Break down the current browser render path into meaningful cost families.

This slice should:

- preserve current rendering behavior by default
- preserve transport behavior unchanged
- measure major client render families rather than only one rolling aggregate
- record a bounded comparative result that can guide the next client optimization

This slice does **not** attempt to implement:

- rendering optimization
- gameplay changes
- transport changes
- UI redesign

## Why This Slice Next

The client now exposes a live draw-time metric. That is useful, but it still leaves the next decision underdetermined:

- you can see render pressure
- but you still cannot tell whether the main cost comes from the minimap, local world drawing, offscreen cues, glows, or recent-effect overlays
- the next useful step is a bounded component breakdown, not another blind rendering tweak

## Use-Case Contract

### Use Case

`MeasureClientRenderComponents`

### Primary Actor

The browser client measuring its own draw-path cost.

### Pre-conditions

- the live render-pressure metric already exists
- the viewport client already has several layered visual systems
- the next client optimization should be evidence-backed

### Trigger

A bounded local render measurement run compares major render families under the current viewport experience.

### Success Outcome

- the dominant browser render families are explicit
- the next client optimization target can be chosen from evidence
- no runtime behavior changes in the default user path

### Failure Or Rejection Cases

- if the slice changes runtime behavior by default, it failed
- if the measurement is too noisy or too vague to guide the next step, it failed
- if the slice drifts into optimization, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Default runtime behavior remains unchanged.
3. The measurement should compare major render families, not tiny line-item drawing calls.
4. The result should be simple enough to support a clear next render slice.

## Minimal Domain Concepts In Scope

- `Render Component`
- `World Draw Cost`
- `Minimap Cost`
- `Overlay Cost`
- `Rolling Draw Baseline`

## Bounded Measurement Interpretation

This slice chooses the smallest useful client-side render breakdown:

- measure the draw path with one major render family omitted at a time, or equivalently instrument a few major sections
- compare those costs against the current full draw baseline
- record the dominant component families in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client instrumentation only
- implementation artifact updated with the render-component findings

## Build Guidance

- keep default behavior unchanged
- instrument a few major rendering sections only
- favor repeatable local comparison over profiler-style completeness

## Initial Test Plan

### Server or measurement tests

- any added client-side instrumentation must degrade safely when unavailable

### Contract tests

- none beyond the current contract validation, because no protocol shape changes are expected

### Integration tests

- none required unless the instrumentation crosses a server boundary

## Scenario Definition

Run the current viewport client and compare:

1. full render path baseline
2. render path without minimap timing contribution
3. render path without overlay timing contribution
4. render path without major world-decoration timing contribution

Record the comparative result.

## Done Criteria

- dominant client render families are explicit
- the next client optimization target is evidence-backed
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- render optimization
- transport optimization
- gameplay changes
- visual redesign
