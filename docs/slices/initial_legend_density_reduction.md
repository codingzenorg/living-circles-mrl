# Slice: Initial Legend Density Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from adding meaning cues to reducing presentation density. The canvas, player card, NPC panel, and encounter log now carry more of the important state explicitly, while the legend has grown into a long glossary that competes with the play surface.

## Discovery Scope

Establish the smallest client-facing simplification that keeps essential legend guidance while reducing visual and cognitive overhead:

- the legend should become shorter and easier to scan
- the most important cue families should remain represented
- the canvas and side panels should carry more of the explanatory load than the legend
- the change should remain presentation-focused rather than altering any world rule

This slice does **not** attempt to implement:

- tutorial systems
- collapsible documentation panels
- onboarding flows
- help modals
- new game mechanics

## Why This Slice Next

Recent slices added several useful visual systems:

- danger and opportunity cues
- crowding-pressure cues
- food-pressure cues
- autonomy-intent cues
- lineage-continuity cues
- recent-event afterglow
- external player and NPC information panels

That made the world more legible, but it also increased legend density. The result is a presentation imbalance: the supporting explanation area is starting to feel busier than the simplified play surface itself.

The next pressure is therefore to simplify the legend so the UI stays readable as a whole.

## Use-Case Contract

### Use Case

`RenderCompactLegendForLivePlay`

### Primary Actor

The player glancing at the interface during ordinary play.

### Pre-conditions

- the browser client already renders multiple bounded cue families on the canvas
- the interface already includes separate player, NPC, and encounter areas
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered for ordinary play.

### Success Outcome

- the legend becomes easier to scan quickly
- the interface feels less overloaded
- the player can still understand the main cue families without reading a dense glossary

### Failure Or Rejection Cases

- if the legend removal makes the cue system cryptic, the slice fails
- if the slice adds new UI systems instead of simplifying the current one, scope is exceeded
- if the change alters world rules or contract shape, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may simplify explanatory UI without altering those semantics.
3. The slice should prefer fewer, stronger legend cues over a long cue inventory.
4. The interface should remain understandable without carrying every detail in the legend.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Legend`
- `Cue Family`
- `Interface Density`
- `Readable Play Surface`

## Bounded Interpretation

This slice chooses the smallest useful legend simplification:

- reduce the legend to the most important cue families
- rely on the current player card, NPC panel, encounter log, and in-canvas naming to carry the rest
- avoid adding tutorial or documentation systems in the same slice

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client layout and presentation only
- implementation notes only if the slice materially changes how the current play UI is organized

## Build Guidance

- prefer consolidation over removal-for-its-own-sake
- keep the most important cue families visible
- avoid re-expanding the UI elsewhere in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the legend should become easier to scan
- the overall interface should feel less crowded
- the canvas and panels should remain sufficient for ordinary play understanding

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the world and sees the current play layout
2. the simplified legend presents only the core cue families
3. the rest of the understanding comes from the world view and panels rather than a dense legend row

## Done Criteria

- the legend is shorter and more readable
- interface density is reduced
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- tutorials
- help drawers
- onboarding
- new cue families
