# Slice: Initial Minimap Orientation In Viewport Mode

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the viewport-based demo track. The browser now renders a player-following camera with a deadzone, which improves scale and comfort, but it also removes the always-visible whole-world overview. That makes spatial orientation weaker: the player can feel the world is larger, but has less immediate sense of where the current viewport sits within it.

## Discovery Scope

Establish the smallest useful orientation aid for viewport mode:

- preserve the main player-follow viewport
- restore some sense of whole-world position
- keep the support chrome restrained
- keep the change in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- interactive maps
- fog of war
- zoomable map systems
- broader navigation UI beyond a small orientation aid

## Why This Slice Next

The viewport slice and deadzone follow improved presence and camera feel. The tradeoff is that the world is no longer visible all at once. The next pressure is therefore orientation:

- the player should still feel the larger world
- the current viewport should be locatable within that world
- orientation should return without giving up the denser viewport presentation

## Use-Case Contract

### Use Case

`RenderViewportOrientationAid`

### Primary Actor

The player moving through the viewport-based fullscreen demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- the runtime contract remains sufficient and unchanged
- the authoritative world remains larger than the current visible window

### Trigger

The page is rendered in viewport mode and the player moves through the world.

### Success Outcome

- the player can more easily infer where the current viewport sits within the larger world
- the main play surface remains dominant
- the added orientation aid stays small and secondary

### Failure Or Rejection Cases

- if the orientation aid competes with the play surface, the slice fails
- if the slice reintroduces a dominant whole-world overview, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may add a small orientation aid without altering those semantics.
3. The player-follow viewport remains the primary rendering mode.
4. The orientation aid should stay secondary to the play surface.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Orientation Aid`
- `Minimap`
- `World Position`

## Bounded Interpretation

This slice chooses the smallest useful orientation mechanism:

- one small passive minimap or equivalent whole-world locator
- current player and/or viewport position shown within world bounds
- no interaction, no zoom, no alternate control mode

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport-mode orientation

## Build Guidance

- keep the orientation aid visually small and secondary
- use authoritative world coordinates already present in the snapshot
- avoid reintroducing full-world rendering as the dominant mode

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the viewport remains the main play surface
- the player can tell where the current view sits within the larger world
- the orientation aid stays visually small and secondary

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player enters the bounded viewport world
2. the player moves through the world while the main viewport follows
3. a small orientation aid shows the player and/or viewport position within the larger world
4. the viewport remains primary while whole-world position becomes easier to understand

## Done Criteria

- viewport mode gains a small orientation aid
- the aid stays secondary to the main play surface
- whole-world position is easier to infer
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- interactive map systems
- zoomable minimap
- fog of war
- server-side navigation metadata
