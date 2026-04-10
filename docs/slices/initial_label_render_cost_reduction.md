# Slice: Initial Label Render Cost Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure, render-family, and world-subcomponent measurement
- current player and NPC name labeling on the play surface
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but reduces browser-side draw cost by narrowing label rendering to what is most useful for play.

## Discovery Scope

Reduce the cost of world labels without changing gameplay semantics.

This slice should:

- preserve authoritative world state unchanged
- preserve support-panel information unchanged
- keep the player name visible on the play surface
- reduce broad repeated label drawing for non-player circles

This slice does **not** attempt to implement:

- transport optimization
- gameplay changes
- a broader visual redesign
- new camera or minimap behavior

## Why This Slice Next

The client now exposes:

- total render pressure
- major render-family breakdown
- world-subcomponent breakdown

That makes one likely next optimization target clear:

- labels are now separately measurable
- they are a pure client-side presentation cost
- their current repeated drawing is less important to moment-to-moment play than circles, foods, and local cues

The next useful step is therefore a bounded label-cost reduction rather than a broader world-draw redesign.

## Use-Case Contract

### Use Case

`ReduceLabelRenderCost`

### Primary Actor

The browser client rendering the current viewport.

### Pre-conditions

- label cost is now explicit in the render breakdown
- support panels already carry detailed identity and state outside the canvas
- the canvas should prioritize motion and local play readability

### Trigger

The browser renders a normal active viewport frame.

### Success Outcome

- the player label remains visible on the play surface
- non-player labels are reduced to a more selective rule
- render pressure drops without changing simulation behavior

### Failure Or Rejection Cases

- if the player loses useful self-identification on the play surface, the slice failed
- if support information becomes ambiguous, the slice failed
- if the slice drifts into broader UI redesign, scope drifted

## Main Business Rules

1. Gameplay and transport remain unchanged.
2. The player name must remain visible on the canvas.
3. Non-player labels should render only when they are likely to help local play.
4. Support panels remain the main source of detailed entity information outside the canvas.

## Minimal Domain Concepts In Scope

- `Player Label`
- `Selective NPC Label`
- `Play-Surface Readability`
- `Client Draw Cost`

## Bounded Implementation Interpretation

This slice chooses the smallest useful label optimization:

- keep the player label always visible
- render non-player labels only under a bounded nearby or highlighted condition
- leave all authoritative world and support data unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client rendering only
- implementation artifact updated with the selected label rule

## Build Guidance

- do not change support panels or transport
- prefer one simple deterministic label rule over several special cases
- keep the play surface easy to read during ordinary movement

## Initial Test Plan

### Server or measurement tests

- existing server and contract validation only, because no server behavior changes are expected

### Contract tests

- unchanged beyond current validation

### Integration tests

- no new integration tests required unless the client bootstrap path changes

## Scenario Definition

Render the current viewport with:

1. player label always visible
2. reduced NPC label drawing under a single bounded rule
3. render pressure observed through the existing client metrics

Record the resulting implementation note.

## Done Criteria

- canvas labels are lighter
- player identity remains obvious on the play surface
- support panels keep detailed identity/state available
- gameplay and transport remain unchanged

## Out Of Scope Follow-Ups

- broader canvas redesign
- transport work
- gameplay changes
- minimap redesign
