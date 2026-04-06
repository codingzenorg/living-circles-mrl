# Slice: Initial Support Text Density Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current UI simplification direction, but shifts from support-panel structure to support-panel text density. The support area now has clearer grouping, yet the player and NPC summaries still carry fairly dense text strings for ordinary live play.

## Discovery Scope

Establish the smallest client-facing simplification that makes support text easier to scan:

- player-support text should read more quickly
- NPC rows should stay useful while becoming lighter-weight
- the support area should remain secondary to the canvas
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- collapsible info systems
- tutorials or onboarding
- broader UI redesign beyond support text presentation

## Why This Slice Next

Recent slices already did the heavier structural work:

- the canvas labels were simplified
- player information moved outside the canvas
- NPC summaries moved into their own panel
- recent encounters stayed in their own panel
- the support area gained a clearer two-column hierarchy

That means the next pressure is no longer about where information lives. It is about how much text each support block requires to parse during ordinary play.

## Use-Case Contract

### Use Case

`RenderLowFrictionSupportText`

### Primary Actor

The player glancing between the canvas and support panels during ordinary play.

### Pre-conditions

- the browser client already renders separate player, NPC, and encounter areas
- the canvas remains the primary interaction surface
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered for ordinary play.

### Success Outcome

- player and NPC text become easier to scan quickly
- the support area feels lighter-weight
- important support information remains available without reading long strings

### Failure Or Rejection Cases

- if the change hides meaningful state with no replacement, the slice fails
- if it adds new explanatory UI systems instead of simplifying the current ones, scope is exceeded
- if it changes world semantics or contract shape, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may simplify support text presentation without altering those semantics.
3. The slice should prefer compact readable text over verbose status lines.
4. The support area should remain secondary to the canvas.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Support Text`
- `Scanability`
- `Compact Status`
- `Secondary UI`

## Bounded Interpretation

This slice chooses the smallest useful text simplification:

- shorten support labels and values where meaning is still preserved
- reduce repeated wording in the player and NPC sections
- avoid removing the encounter log or primary player state

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes support-text presentation

## Build Guidance

- prefer concise text and clearer value emphasis over adding new visual ornament
- keep the encounter log readable, but focus primarily on the player and NPC support text
- avoid reintroducing density elsewhere in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the player panel should scan faster than before
- the NPC panel should remain useful but feel lighter
- the support area should read as compact support rather than a status dump

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the world and the support area renders under the canvas
2. the player glances between the canvas and the support area during ordinary play
3. the simplified support text reduces scan effort while preserving important state

## Done Criteria

- support text is more compact and readable
- support-area density is reduced further
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- tutorials
- settings systems
- new panels
- new cue families
