# Slice: Initial Camera Deadzone Follow

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current viewport-based demo track. The browser now renders a bounded player-following viewport instead of shrinking the whole world into the stage. That fixes scale and presence, but the current camera still centers directly on the player whenever possible, which can make the view feel mechanically locked rather than comfortably staged.

## Discovery Scope

Establish the smallest useful improvement to camera feel:

- the viewport should still follow the player
- the camera should not need to recenter on every small movement
- world-edge clamping should still work
- the change should stay in client projection and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- cinematic camera smoothing
- zoom controls
- minimap systems
- alternative camera targets

## Why This Slice Next

The last slice established the correct presentation model:

- large authoritative world
- bounded viewport
- player-follow camera
- world-edge clamp

The next pressure is no longer scale. It is camera comfort. A strict center-lock follow works, but it can feel visually twitchy and over-responsive for small movement changes.

## Use-Case Contract

### Use Case

`RenderDeadzoneCameraFollow`

### Primary Actor

The player moving through the fullscreen browser demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- the runtime contract remains sufficient and unchanged
- the server remains authoritative for world-space coordinates

### Trigger

The player moves through the world while the browser renders the viewport.

### Success Outcome

- the camera still follows the player
- small movements inside a bounded inner region do not force constant recenters
- the viewport feels more comfortable without losing clarity
- world-edge clamping still behaves cleanly

### Failure Or Rejection Cases

- if the camera stops feeling reliably player-following, the slice fails
- if the slice adds cinematic smoothing or other new camera systems, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may adjust how the viewport follows the player without altering those semantics.
3. The camera should still be clearly player-following.
4. The camera should only recenter when the player leaves a bounded deadzone.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Player Follow Camera`
- `Camera Deadzone`
- `World Bounds Clamp`

## Bounded Interpretation

This slice chooses the smallest useful camera comfort rule:

- define a centered deadzone region inside the viewport
- keep the camera still while the player remains within that region
- move the camera only enough to bring the player back into the deadzone
- keep edge clamping unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport projection behavior

## Build Guidance

- prefer a simple deterministic deadzone over smoothing curves or interpolation systems
- keep the viewport logic easy to reason about
- avoid introducing zoom, minimap, or multi-target camera behavior in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the camera still follows the player
- small movements near the center do not constantly shift the whole viewport
- leaving the inner region causes the camera to move just enough to keep up
- the camera still clamps cleanly at world edges

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player enters the bounded viewport world
2. the player makes small movements near the current camera center
3. the camera stays more stable until the player leaves the deadzone
4. the player approaches world bounds and camera clamping still works

## Done Criteria

- camera deadzone follow is implemented
- the viewport feels more comfortable than strict center-lock
- world-edge clamping still works
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- cinematic smoothing
- zoom controls
- minimap
- alternative camera targets
