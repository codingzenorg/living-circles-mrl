# Slice: Initial Offscreen Edge Awareness In Viewport Mode

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the viewport-based demo track. The browser now renders a player-following camera with a deadzone and a small minimap, which improves scale, comfort, and orientation. The remaining pressure is immediate situational awareness: once the world is no longer fully visible, nearby dangers or opportunities just outside the current viewport can disappear too abruptly.

## Discovery Scope

Establish the smallest useful offscreen-awareness aid for viewport mode:

- preserve the main player-follow viewport
- preserve the minimap as a secondary orientation aid
- make nearby offscreen world pressure more legible at the viewport edges
- keep the change in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- radar systems
- full object tracking
- new minimap interactions
- broader sensing systems beyond local offscreen awareness

## Why This Slice Next

The viewport slice and minimap solved scale and orientation, but they do not fully solve immediate local awareness:

- the minimap is good for whole-world position
- it is weaker for quick nearby danger/opportunity reading during movement
- the player still benefits from some local indication when meaningful entities are just outside the visible viewport

## Use-Case Contract

### Use Case

`RenderOffscreenEdgeAwareness`

### Primary Actor

The player moving through the viewport-based fullscreen demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- the runtime contract remains sufficient and unchanged
- the minimap already provides a secondary whole-world locator

### Trigger

Meaningful world entities are near the current viewport but outside the visible area.

### Success Outcome

- the player gains lightweight awareness of nearby offscreen pressure
- the main viewport remains the primary play surface
- the minimap remains secondary and complementary

### Failure Or Rejection Cases

- if edge cues become visually noisy or dominant, the slice fails
- if the slice turns into general object tracking, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may add local offscreen-awareness cues without altering those semantics.
3. The main viewport remains the dominant play surface.
4. Offscreen cues should stay lightweight and local.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Offscreen Awareness`
- `Edge Cue`
- `Local Pressure`

## Bounded Interpretation

This slice chooses the smallest useful offscreen-awareness mechanism:

- indicate nearby offscreen entities or pressures at the viewport edge
- keep cues lightweight and local
- avoid turning the client into a generalized tracking overlay

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport-mode awareness cues

## Build Guidance

- prefer subtle edge cues over labels or large overlays
- keep the focus on nearby offscreen meaning, not distant whole-world tracking
- avoid making the viewport feel cluttered

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the player can notice nearby meaningful offscreen pressure more easily
- the viewport remains the dominant surface
- cues stay lightweight and do not overwhelm the scene

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player moves through the viewport-based world
2. meaningful entities drift just outside the visible camera window
3. the client shows lightweight edge awareness for that nearby offscreen pressure
4. the player retains better local awareness without losing the viewport presentation

## Done Criteria

- viewport mode gains lightweight offscreen edge awareness
- cues remain secondary and local
- the main viewport remains dominant
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- radar systems
- general object tracking
- interactive minimap
- server-side awareness metadata
