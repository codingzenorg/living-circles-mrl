# Slice: Initial Support Panel Growth Bounds

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current support-area simplification track, but focuses on vertical growth rather than text density or panel arrangement. The player card is now stable and primary, and NPC plus encounter panels are side by side below it, but those lower panels can still grow with content and begin to compete with the canvas.

## Discovery Scope

Establish the smallest client-facing constraint that keeps support panels from visually taking over:

- NPC and encounter panels should stay bounded as their content grows
- the canvas should remain the dominant surface during ordinary play
- the player summary should remain stable above the lower support row
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- pagination or archival systems
- new panel types
- broader UI redesign beyond bounded support growth

## Why This Slice Next

Recent slices already:

- clarified panel hierarchy
- reduced legend density
- reduced support text density
- stabilized the player card position by moving it above the lower row

That means the next pressure is no longer where support information lives, but how far it is allowed to expand when the NPC list or recent encounter log becomes longer.

## Use-Case Contract

### Use Case

`RenderBoundedSupportPanels`

### Primary Actor

The player glancing between the canvas and the lower support row during ordinary play.

### Pre-conditions

- the browser client already renders a player card above NPC and encounter panels
- the canvas remains the primary interaction surface
- the runtime contract remains sufficient and unchanged

### Trigger

NPC summaries or recent encounter history grows during ordinary play.

### Success Outcome

- lower support panels remain useful without taking over the page
- the canvas remains visually primary
- support growth is handled without changing world semantics

### Failure Or Rejection Cases

- if bounded growth hides useful support information with no readable access, the slice fails
- if the change introduces new data systems instead of bounded presentation behavior, scope is exceeded
- if it changes world semantics or contract shape, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may bound support-panel growth without altering those semantics.
3. The canvas should remain visually primary.
4. The player card should remain stable and separate from lower support growth.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Bounded Panel Growth`
- `Canvas Primacy`
- `Lower Support Row`
- `Readable Overflow`

## Bounded Interpretation

This slice chooses the smallest useful growth constraint:

- cap the vertical growth of NPC and encounter panels
- preserve readability within those bounds
- avoid introducing new storage, filtering, or tutorial systems

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes support-panel growth behavior

## Build Guidance

- prefer simple bounded overflow behavior over new interaction systems
- keep the player card unaffected by lower-panel growth
- avoid reintroducing density or large explanatory chrome in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the player card remains fixed above the lower row
- longer NPC or encounter content does not keep pushing the page downward
- the canvas still reads as the dominant surface during ordinary play

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the page renders with the player card above the lower support row
2. NPC summaries and recent encounters accumulate during ordinary play
3. the lower panels stay bounded while remaining readable

## Done Criteria

- lower support panels are visually bounded
- the canvas remains primary
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- paging
- filtering
- panel collapse/expand systems
- new support panel types
