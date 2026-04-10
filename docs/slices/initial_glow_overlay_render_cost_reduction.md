# Slice: Initial Glow Overlay Render Cost Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure plus world and overlay render breakdowns
- current glow-style overlays including crowding zones, food opportunity glows, and similar broad ambient fields
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but reduces browser-side overlay draw cost by narrowing the heaviest ambient glow work.

## Discovery Scope

Reduce the cost of glow-field overlays without changing gameplay semantics.

This slice should:

- preserve authoritative world state unchanged
- preserve transport behavior unchanged
- keep the strongest local cue meanings intact
- reduce the breadth or frequency of ambient glow drawing where it is least necessary

This slice does **not** attempt to implement:

- broader overlay redesign
- gameplay changes
- transport changes
- minimap or camera changes

## Why This Slice Next

The client now exposes:

- total render pressure
- major render-family breakdown
- world-render subfamilies
- overlay-render subfamilies

That makes the next likely render optimization target much narrower:

- overlays are now measurable as glows, offscreen awareness, recent effects, and cues
- glow fields are the broadest ambient overlay work and are usually more expensive than simple lines or markers
- the next useful step is a bounded glow-cost reduction, not another generic overlay tweak

## Use-Case Contract

### Use Case

`ReduceGlowOverlayRenderCost`

### Primary Actor

The browser client rendering the current viewport.

### Pre-conditions

- overlay render families are now explicit
- ambient glow fields already exist for crowding and food readability
- the next client optimization should still preserve local legibility

### Trigger

The browser renders a normal active viewport frame.

### Success Outcome

- local cue meaning stays readable
- broad ambient glow work is reduced
- render pressure drops without changing simulation behavior

### Failure Or Rejection Cases

- if local crowding or food opportunity becomes hard to read, the slice failed
- if the slice drifts into broad visual redesign, scope drifted
- if transport or gameplay change, it failed

## Main Business Rules

1. Gameplay and transport remain unchanged.
2. The strongest local crowding and food cues must remain readable.
3. Glow work should be reduced where it is broad ambient decoration rather than essential immediate guidance.
4. Existing support panels and non-glow cues remain unchanged.

## Minimal Domain Concepts In Scope

- `Glow Overlay`
- `Crowding Glow`
- `Food Opportunity Glow`
- `Local Cue Readability`
- `Client Draw Cost`

## Bounded Implementation Interpretation

This slice chooses the smallest useful glow optimization:

- keep the existing cue vocabulary
- reduce glow radius, frequency, or qualifying conditions for the broadest ambient fields
- leave recent effects, offscreen awareness, lineage, and intent cues unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client rendering only
- implementation artifact updated with the selected glow rule

## Build Guidance

- do not change support panels or transport
- prefer one simple deterministic glow reduction over several layered exceptions
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

1. unchanged non-glow overlays
2. reduced ambient glow work under one simple deterministic rule
3. render pressure observed through the existing client metrics

Record the resulting implementation note.

## Done Criteria

- glow overlays are lighter
- local cue meaning remains readable
- gameplay and transport remain unchanged

## Out Of Scope Follow-Ups

- broader overlay redesign
- transport work
- gameplay changes
- minimap redesign
