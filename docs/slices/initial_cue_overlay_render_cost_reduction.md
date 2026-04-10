# Slice: Initial Cue Overlay Render Cost Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure plus world and overlay render breakdowns
- current local cue overlays including lineage links, intent cues, and player heading cue
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but reduces browser-side overlay draw cost by narrowing cue-overlay work.

## Discovery Scope

Reduce the cost of cue-style overlays without changing gameplay semantics.

This slice should:

- preserve authoritative world state unchanged
- preserve transport behavior unchanged
- keep the main local meaning of heading, intent, and lineage cues intact
- reduce cue draw work where it is least necessary

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
- bounded reductions for glow and offscreen overlay cost

That makes the next likely overlay optimization target narrower:

- recent effects are already short-lived and event-bounded
- cue overlays still draw repeatedly every frame
- cue work can likely be narrowed without removing the main local play meaning

The next useful step is therefore a bounded cue-cost reduction rather than another broad overlay tweak.

## Use-Case Contract

### Use Case

`ReduceCueOverlayRenderCost`

### Primary Actor

The browser client rendering the current viewport.

### Pre-conditions

- overlay render families are explicit
- heading, intent, and lineage cues already exist
- glow and offscreen overlay work have already been narrowed

### Trigger

The browser renders a normal active viewport frame.

### Success Outcome

- cue meaning remains readable
- repeated cue draw work is reduced
- render pressure drops without changing simulation behavior

### Failure Or Rejection Cases

- if the player loses useful local motion or lineage meaning, the slice failed
- if the slice drifts into minimap, camera, or transport redesign, scope drifted
- if gameplay or transport change, it failed

## Main Business Rules

1. Gameplay and transport remain unchanged.
2. Local cue meaning should remain readable.
3. Cue draw work should be reduced where support is already duplicated by nearby geometry or support panels.
4. Non-cue overlays remain unchanged.

## Minimal Domain Concepts In Scope

- `Cue Overlay`
- `Player Heading Cue`
- `Autonomy Intent Cue`
- `Lineage Link`
- `Client Draw Cost`

## Bounded Implementation Interpretation

This slice chooses the smallest useful cue optimization:

- keep the current cue vocabulary
- reduce cue frequency, qualifying distance, or duplicate cue work under one simple deterministic rule
- leave minimap, offscreen awareness, recent effects, and glow behavior unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client rendering only
- implementation artifact updated with the selected cue rule

## Build Guidance

- do not change support panels, minimap, or transport
- prefer one simple deterministic cue reduction over several layered exceptions
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

1. unchanged glow, offscreen, and recent-effect overlays
2. reduced cue overlay work under one simple deterministic rule
3. render pressure observed through the existing client metrics

Record the resulting implementation note.

## Done Criteria

- cue overlays are lighter
- local meaning remains readable
- gameplay and transport remain unchanged

## Out Of Scope Follow-Ups

- minimap redesign
- camera redesign
- transport work
- gameplay changes
