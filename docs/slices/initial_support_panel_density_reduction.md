# Slice: Initial Support Panel Density Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from reducing legend density to reducing support-panel density. The interface now places player state, NPC state, and recent encounters outside the canvas, which is good for the play surface, but those supporting panels are becoming their own dense stack.

## Discovery Scope

Establish the smallest client-facing simplification that makes the support area easier to scan:

- player, NPC, and encounter information should feel more clearly prioritized
- the supporting panels should take less visual effort to parse
- the change should stay in layout and presentation, not world semantics
- the canvas should remain the primary play surface

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- collapsible app chrome or settings systems
- tutorials or help drawers
- broad UI redesign beyond the current support area

## Why This Slice Next

Recent slices did the right thing by moving more information outside the canvas:

- player stats moved into a dedicated card
- NPC summaries moved into a dedicated panel
- recent encounters remained in a dedicated log
- the legend was reduced so the play surface is less crowded

That improves the canvas, but it also creates a new pressure: the support area is now the densest part of the screen. The next step is not more explanation, but clearer hierarchy and lower scan cost in that support area.

## Use-Case Contract

### Use Case

`RenderLowFrictionSupportPanels`

### Primary Actor

The player glancing between the canvas and the supporting information area during ordinary play.

### Pre-conditions

- the browser client already renders separate player, NPC, and encounter panels
- the canvas remains the primary interaction surface
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered for ordinary play.

### Success Outcome

- the support area becomes easier to scan quickly
- the player can find the most important supporting information with less effort
- the external information layout complements the canvas instead of competing with it

### Failure Or Rejection Cases

- if the change hides important information with no replacement, the slice fails
- if it adds new explanatory UI systems instead of simplifying the current ones, scope is exceeded
- if it changes world semantics or contract shape, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reorganize supporting presentation without altering those semantics.
3. The slice should prefer clearer hierarchy over more information.
4. The canvas should remain the primary play surface.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Support Panel`
- `Hierarchy`
- `Scanability`
- `Play Surface Priority`

## Bounded Interpretation

This slice chooses the smallest useful support-area simplification:

- reorganize or restyle the existing player, NPC, and encounter areas
- make priority clearer through spacing, grouping, or tighter layout
- avoid introducing new panel systems or settings

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client layout and presentation only
- implementation notes only if the slice materially changes how the current play UI is organized

## Build Guidance

- prefer hierarchy and grouping improvements over more ornament
- keep the support area compact enough that it feels secondary to the canvas
- avoid removing useful information unless the same meaning is preserved elsewhere

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the support area should feel easier to scan
- the most important player information should be easiest to find
- NPC and encounter sections should remain useful but visually secondary

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the world and sees the current canvas plus support area
2. the support area is reorganized to reduce scan cost
3. the player can read the needed external information more quickly while keeping the canvas as the primary focus

## Done Criteria

- support-area density is reduced
- visual hierarchy is clearer
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- tutorials
- help drawers
- settings systems
- new cue families
