# Slice: Initial Offscreen Food Awareness In Viewport Mode

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the viewport-based demo track. The browser now has a player-following camera, deadzone follow, minimap, offscreen circle awareness, and a local player heading cue. Those changes improve direction and nearby pressure awareness, but one immediate local gap remains: food just outside the current viewport is still easy to miss until it crosses into view, even though food is a first-order survival driver in the loop.

## Discovery Scope

Establish the smallest useful offscreen food-awareness aid for viewport mode:

- preserve the current player-following viewport
- preserve minimap and offscreen circle awareness as secondary aids
- add lightweight local awareness for nearby offscreen food
- keep the change in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- resource radar systems
- full object tracking
- new minimap interactions
- broader sensing systems beyond local offscreen food awareness

## Why This Slice Next

The current viewport mode already helps with:

- whole-world orientation through the minimap
- nearby offscreen circle pressure through edge cues
- local motion direction through the heading cue

But it still underplays one nearby survival input:

- food just outside the camera window can remain locally invisible
- the player may notice threats before recovery options
- viewport mode needs a small complementary cue for nearby offscreen food

## Use-Case Contract

### Use Case

`RenderOffscreenFoodAwareness`

### Primary Actor

The player moving through the viewport-based fullscreen demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- the runtime contract remains sufficient and unchanged
- food remains a visible authoritative world entity

### Trigger

Food lies near the current camera window but outside the visible area.

### Success Outcome

- the player gains lightweight awareness of nearby offscreen food
- the main viewport remains the dominant play surface
- food awareness complements rather than overwhelms existing offscreen circle cues

### Failure Or Rejection Cases

- if food cues become visually noisy or dominant, the slice fails
- if the slice turns into general resource tracking, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may add lightweight offscreen food cues without altering those semantics.
3. The main viewport remains the dominant play surface.
4. Food cues should stay local, secondary, and readable.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Offscreen Food Awareness`
- `Edge Cue`
- `Local Recovery Opportunity`

## Bounded Interpretation

This slice chooses the smallest useful offscreen-food mechanism:

- indicate nearby offscreen food at the viewport edge
- keep cues lightweight and local
- avoid turning the client into a generalized resource tracking overlay

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport-mode food-awareness cues

## Build Guidance

- prefer subtle edge cues over labels or large overlays
- keep the focus on nearby offscreen recovery opportunity, not distant tracking
- avoid making the viewport feel cluttered beside the existing offscreen circle cues

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- nearby offscreen food is easier to notice
- the viewport remains the dominant surface
- food cues stay lightweight and complementary to existing awareness aids

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player moves through the viewport-based world
2. food drifts just outside the visible camera window
3. the client shows lightweight local awareness for that nearby offscreen food
4. the player retains better recovery awareness without losing viewport focus

## Done Criteria

- viewport mode gains lightweight offscreen food awareness
- cues remain secondary and local
- the main viewport remains dominant
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- resource radar systems
- general tracking overlays
- interactive minimap
- server-side awareness metadata
