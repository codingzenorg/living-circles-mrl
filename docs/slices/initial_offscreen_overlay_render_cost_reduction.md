# Slice: Initial Offscreen Overlay Render Cost Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure plus world and overlay render breakdowns
- current offscreen awareness overlays for nearby circles and foods at the viewport edge
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but reduces browser-side overlay draw cost by narrowing offscreen awareness work.

## Discovery Scope

Reduce the cost of offscreen-awareness overlays without changing gameplay semantics.

This slice should:

- preserve authoritative world state unchanged
- preserve transport behavior unchanged
- keep the main local orientation meaning intact
- reduce offscreen marker work where it is least necessary

This slice does **not** attempt to implement:

- broader overlay redesign
- gameplay changes
- transport changes
- minimap or camera changes

## Why This Slice Next

The client now exposes:

- total render pressure
- major render-family breakdown
- world and overlay subfamily breakdowns
- a bounded glow-cost reduction

That makes the next likely overlay optimization target narrower:

- glow work is already reduced
- offscreen awareness is still a repeated overlay family drawn every frame
- minimap and player heading already provide orientation support, so the offscreen path can be narrowed without removing orientation entirely

The next useful step is therefore a bounded offscreen-cost reduction rather than another broad overlay tweak.

## Use-Case Contract

### Use Case

`ReduceOffscreenOverlayRenderCost`

### Primary Actor

The browser client rendering the current viewport.

### Pre-conditions

- overlay render families are explicit
- minimap and player heading cues already exist
- offscreen awareness is currently one of the remaining repeated overlay families

### Trigger

The browser renders a normal active viewport frame.

### Success Outcome

- offscreen orientation remains usable
- repeated offscreen marker work is reduced
- render pressure drops without changing simulation behavior

### Failure Or Rejection Cases

- if nearby offscreen pressure becomes hard to read, the slice failed
- if the slice drifts into minimap or camera redesign, scope drifted
- if transport or gameplay change, it failed

## Main Business Rules

1. Gameplay and transport remain unchanged.
2. Offscreen orientation should remain usable through the combined viewport aids.
3. Offscreen overlay work should be reduced where the minimap and player heading already provide enough support.
4. Non-offscreen overlays remain unchanged.

## Minimal Domain Concepts In Scope

- `Offscreen Awareness Overlay`
- `Nearby Offscreen Pressure`
- `Viewport Edge Marker`
- `Orientation Support`
- `Client Draw Cost`

## Bounded Implementation Interpretation

This slice chooses the smallest useful offscreen optimization:

- keep the current offscreen cue vocabulary
- reduce marker count, qualifying distance, or duplicate marker work under one simple deterministic rule
- leave minimap, heading cue, recent effects, and glow behavior unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client rendering only
- implementation artifact updated with the selected offscreen rule

## Build Guidance

- do not change support panels, minimap, or transport
- prefer one simple deterministic offscreen reduction over several layered exceptions
- preserve the canvas as a readable play surface, not just a cheaper one

## Initial Test Plan

### Server or measurement tests

- existing server and contract validation only, because no server behavior changes are expected

### Contract tests

- unchanged beyond current validation

### Integration tests

- no new integration tests required unless the client bootstrap path changes

## Scenario Definition

Render the current viewport with:

1. unchanged minimap and heading support
2. reduced offscreen overlay work under one simple deterministic rule
3. render pressure observed through the existing client metrics

Record the resulting implementation note.

## Done Criteria

- offscreen overlays are lighter
- orientation remains readable
- gameplay and transport remain unchanged

## Out Of Scope Follow-Ups

- minimap redesign
- camera redesign
- transport work
- gameplay changes
