# Slice: Initial Recent Effect Overlay Render Cost Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure plus world and overlay render breakdowns
- current recent-effect afterglow overlays for fights, reproduction, and continuity outcomes
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but reduces browser-side overlay draw cost by narrowing recent-effect overlay work.

## Discovery Scope

Reduce the cost of recent-effect overlays without changing gameplay semantics.

This slice should:

- preserve authoritative world state unchanged
- preserve transport behavior unchanged
- keep the main meaning of recent fights, reproduction, and continuity visible
- reduce recent-effect draw work where it is least necessary

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
- bounded reductions for glow, offscreen, and cue overlay cost

That makes the next likely overlay optimization target narrower:

- recent effects are the remaining repeated overlay family
- they are useful for local legibility, but their current visual persistence may be stronger than necessary
- the next useful step is therefore a bounded recent-effect reduction rather than another broad overlay change

## Use-Case Contract

### Use Case

`ReduceRecentEffectOverlayRenderCost`

### Primary Actor

The browser client rendering the current viewport.

### Pre-conditions

- overlay render families are explicit
- recent effects already exist for resolved conflict, reproduction, and continuity outcomes
- glow, offscreen, and cue overlay work have already been narrowed

### Trigger

The browser renders a normal active viewport frame.

### Success Outcome

- recent outcomes remain visible
- repeated recent-effect draw work is reduced
- render pressure drops without changing simulation behavior

### Failure Or Rejection Cases

- if recent fights or reproduction outcomes become hard to notice, the slice failed
- if the slice drifts into transport, minimap, or camera redesign, scope drifted
- if gameplay or transport change, it failed

## Main Business Rules

1. Gameplay and transport remain unchanged.
2. Recent effect meaning should remain readable.
3. Recent-effect draw work should be reduced by narrowing persistence or effect radius rather than removing the cue vocabulary.
4. Non-recent-effect overlays remain unchanged.

## Minimal Domain Concepts In Scope

- `Recent Effect Overlay`
- `Fight Afterglow`
- `Reproduction Afterglow`
- `Continuity Afterglow`
- `Client Draw Cost`

## Bounded Implementation Interpretation

This slice chooses the smallest useful recent-effect optimization:

- keep the current recent-effect cue vocabulary
- reduce persistence or radius under one simple deterministic rule
- leave minimap, offscreen awareness, glows, and cue overlays unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client rendering only
- implementation artifact updated with the selected recent-effect rule

## Build Guidance

- do not change support panels, minimap, or transport
- prefer one simple deterministic recent-effect reduction over several layered exceptions
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

1. unchanged glow, offscreen, and cue overlays
2. reduced recent-effect overlay work under one simple deterministic rule
3. render pressure observed through the existing client metrics

Record the resulting implementation note.

## Done Criteria

- recent-effect overlays are lighter
- recent outcome meaning remains readable
- gameplay and transport remain unchanged

## Out Of Scope Follow-Ups

- minimap redesign
- camera redesign
- transport work
- gameplay changes
