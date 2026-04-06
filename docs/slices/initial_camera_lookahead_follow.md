# Slice: Initial Camera Lookahead Follow

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the viewport-based demo track. The browser now has a player-following camera, deadzone follow, a minimap, lightweight offscreen circle and food awareness, and a local heading cue. Those additions improve orientation and nearby awareness, but the camera still sits too neutrally around the player once movement intent is already clear. The next pressure is not more information, but a camera that gives a little more room in the direction the player is actually moving.

## Discovery Scope

Establish the smallest useful camera lookahead for viewport mode:

- preserve the current player-following viewport
- preserve deadzone behavior and world-edge clamping
- add a lightweight directional lookahead tied to recent authoritative player motion
- keep the change in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- cinematic camera systems
- zoom changes
- prediction beyond recent authoritative movement
- broader navigation UI beyond camera behavior

## Why This Slice Next

The current viewport mode already helps with:

- whole-world orientation through the minimap
- nearby offscreen awareness through edge cues
- local directionality through the heading cue

But camera feel still underplays movement:

- the player can read where they are going
- the camera still gives equal space in all directions
- the viewport could feel more alive if it slightly favored the player’s heading

## Use-Case Contract

### Use Case

`RenderCameraLookaheadFollow`

### Primary Actor

The player moving through the viewport-based fullscreen demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- deadzone follow is already active
- the runtime contract remains sufficient and unchanged
- player movement is already represented by authoritative position changes

### Trigger

The player moves through the world with a readable heading.

### Success Outcome

- the camera gives a small amount of extra forward room in the player’s direction of travel
- the viewport feels more responsive without becoming floaty or cinematic
- world-edge clamping and the deadzone remain intact

### Failure Or Rejection Cases

- if the camera becomes unstable or overly dramatic, the slice fails
- if the slice introduces prediction-heavy or cinematic camera behavior, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may add a lightweight directional camera lookahead without altering those semantics.
3. Deadzone follow and world-edge clamping remain in effect.
4. Lookahead should stay small, readable, and tied to recent authoritative motion.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Camera Lookahead`
- `Player Heading`
- `Forward Room`

## Bounded Interpretation

This slice chooses the smallest useful lookahead mechanism:

- derive a small camera offset from recent authoritative player motion
- preserve the current deadzone model rather than replacing it
- avoid cinematic smoothing, zooming, or prediction-heavy tracking

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport-mode camera behavior

## Build Guidance

- prefer subtle forward bias over dramatic camera shifts
- keep the player readable and central enough for control clarity
- avoid making the viewport feel disconnected from the authoritative player position

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- movement gives the player slightly more room in the direction of travel
- the camera still feels stable near walls and during small adjustments
- the viewport remains the dominant play surface

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player moves through the viewport-based world
2. the camera follows with a small directional lookahead based on recent authoritative motion
3. the player experiences slightly more forward room without losing control clarity

## Done Criteria

- viewport mode gains a lightweight directional camera lookahead
- deadzone follow and world-edge clamping still work
- the viewport feels more responsive without becoming visually noisy
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- cinematic cameras
- zoom controls
- prediction-heavy tracking
- server-side camera metadata
