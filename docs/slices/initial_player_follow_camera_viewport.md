# Slice: Initial Player Follow Camera Viewport

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice changes the fullscreen demo from a scaled whole-world overview into a player-following viewport. The authoritative world remains large and unchanged, but the browser stops trying to fit the whole world into the visible stage and instead renders a bounded camera window into that world.

## Discovery Scope

Establish the smallest useful camera/viewport model for the current demo:

- the world remains authoritative and unchanged on the server
- the client renders a bounded viewport into the world rather than scaling the full world down
- the camera follows the player and clamps at world edges
- the presentation stays eye-catching without introducing broader camera systems

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- zoom controls
- minimap systems
- cinematic camera motion
- arbitrary camera targets

## Why This Slice Next

Recent fullscreen refinement improved layout, but the current browser still scales the whole world into the stage. That weakens visual density and undercuts the “large living world” feel the fullscreen work was moving toward.

The new pressure is not more UI polish. It is presentation mode:

- whole-world scaling is good for overview, but weak for presence
- 1:1 rendering with a camera viewport is better for eye-candy and motion readability
- the world is now large enough that a viewport model is more coherent than shrinking everything to fit

## Use-Case Contract

### Use Case

`RenderPlayerFollowViewport`

### Primary Actor

The player moving through the fullscreen browser demo.

### Pre-conditions

- the server still publishes authoritative world-space coordinates
- the browser already renders the world on canvas
- the runtime contract remains sufficient and unchanged

### Trigger

The player enters the world and the browser renders the simulation.

### Success Outcome

- the canvas shows a player-following window into the world
- world rendering stays 1:1 inside that viewport
- the camera remains bounded to the world edges
- the experience feels more spatially immediate than a scaled full-world view

### Failure Or Rejection Cases

- if the browser still shrinks the whole world to fit, the slice fails
- if off-screen handling becomes confusing with no bounded camera rule, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may change how much of the world is visible at once without altering those semantics.
3. The browser should render a bounded viewport into the world rather than the whole world scaled down.
4. The camera should follow the player while respecting world bounds.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Player Follow Camera`
- `World Bounds Clamp`
- `Client Projection`

## Bounded Interpretation

This slice chooses the smallest useful camera model:

- one viewport sized by the canvas element
- camera centered on the player when possible
- world-edge clamp
- existing cues rendered only when their world positions fall within the current viewport

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes rendering projection and viewport behavior

## Build Guidance

- keep world coordinates authoritative and unchanged
- introduce the smallest possible projection layer from world space to viewport space
- avoid adding zoom, minimap, or alternative camera modes in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the browser no longer shows the whole world shrunk to fit
- moving the player causes the viewport to follow
- approaching world bounds causes the camera to clamp cleanly
- the result feels more spatially immediate than the previous scaled overview

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player enters the world and sees a bounded viewport rather than the full world
2. the player moves through the world
3. the camera follows the player and clamps at world bounds
4. the viewport preserves a denser, more eye-catching world presentation

## Done Criteria

- browser rendering uses a bounded player-following viewport
- world rendering stays 1:1 inside that viewport
- world-edge camera clamping works
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- zoom controls
- minimap
- cinematic smoothing
- alternative camera targets
