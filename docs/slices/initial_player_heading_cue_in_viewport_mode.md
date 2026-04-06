# Slice: Initial Player Heading Cue In Viewport Mode

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the viewport-based demo track. The browser now has a player-following camera, deadzone follow, a small minimap, and lightweight offscreen edge cues. Those additions improve orientation and nearby awareness, but the player still lacks a simple local cue for where their current movement is pulling the viewport attention. The next pressure is not more global information, but more immediate local heading clarity.

## Discovery Scope

Establish the smallest useful heading cue for viewport mode:

- preserve the current player-following viewport
- preserve the minimap and edge cues as secondary aids
- add a lightweight local indication of the player's current movement heading
- keep the change in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- path prediction systems
- waypoint systems
- new camera targets
- broader navigation UI beyond a local heading cue

## Why This Slice Next

The viewport, deadzone, minimap, and edge-awareness slices collectively solve scale and orientation better than the earlier whole-world view. The remaining gap is immediate local directionality:

- the player can see where they are in the larger world
- the player can sense nearby offscreen pressure
- but the viewport still lacks a simple local heading cue tied to current movement

## Use-Case Contract

### Use Case

`RenderPlayerHeadingCue`

### Primary Actor

The player moving through the viewport-based fullscreen demo.

### Pre-conditions

- the browser already renders a bounded player-following viewport
- the runtime contract remains sufficient and unchanged
- player movement is already represented by authoritative position changes

### Trigger

The player moves through the world.

### Success Outcome

- the player gains a lightweight local sense of heading
- the cue stays tied to immediate movement rather than whole-world systems
- the viewport remains the dominant play surface

### Failure Or Rejection Cases

- if the cue becomes visually noisy or dominant, the slice fails
- if the slice introduces path prediction or navigation systems, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics and coordinates.
2. The client may add a lightweight heading cue without altering those semantics.
3. The viewport remains the dominant play surface.
4. The cue should stay local and immediately readable.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Viewport`
- `Player Heading`
- `Local Direction Cue`
- `Immediate Movement Readability`

## Bounded Interpretation

This slice chooses the smallest useful heading mechanism:

- infer heading from recent authoritative player motion
- render one lightweight local cue near the player
- avoid prediction, waypoints, or broader navigation overlays

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes viewport-mode direction cues

## Build Guidance

- prefer a small, elegant local cue over labels or large overlays
- keep the cue tied to actual recent motion, not a separate predicted path
- avoid making the viewport feel more cluttered

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the player can read local heading more easily while moving
- the cue stays subtle and local
- the viewport remains the primary focus

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the player moves through the viewport-based world
2. the client renders a lightweight heading cue near the player
3. the player can read local movement direction more easily without relying on minimap or edge cues

## Done Criteria

- viewport mode gains a lightweight player heading cue
- the cue stays local and secondary
- the main viewport remains dominant
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- path prediction
- waypoints
- new camera targets
- server-side direction metadata
