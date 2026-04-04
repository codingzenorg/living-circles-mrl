# Slice: Initial Embodied Movement Boundary Without Radius

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for circle movement that now avoids derived radius as a hidden world-boundary clamp

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for movement and world-edge clamping.

This slice narrows one remaining radius shortcut in movement by making world bounds care about the visible parent body rather than grown derived radius.

## Discovery Scope

Establish the smallest deterministic change that makes movement boundaries better match the current embodied model:

- parent movement clamping should use the visible parent-core body rather than derived grown radius
- current attached-child orbit rendering remains unchanged
- current food, contact, fight, reproduction, continuity, and steering rules remain unchanged

This slice does **not** attempt to implement:

- removal of radius from rendering
- removal of radius from snapshot payloads
- changing attached-child orbit distance
- new wall-collision effects or damage
- client-side prediction or smoothing

## Why This Slice Next

Recent slices already removed derived radius from the main hidden leverage paths:

- food collection no longer uses grown derived radius as silent reach
- circle contact no longer uses grown derived radius as silent encounter reach
- same-shape fight resolution no longer uses derived radius as a tie-break

But movement still clamps the parent core against the world edge using grown derived radius. That means a parent with more children is pushed farther away from the wall even though contact and feeding now already use the visible parent body plus attached children. The model pressure is now to make movement boundaries follow the same embodied geometry rather than an older growth shortcut.

This slice is the narrowest next step because it:

- changes only parent-body boundary clamping
- keeps rendering and orbit layout stable
- preserves all current interaction and energy semantics
- removes one remaining hidden radius effect without broad visual redesign

## Use-Case Contract

### Use Case

`ClampParentMovementByVisibleBody`

### Primary Actor

Any active player or autonomous parent circle that moves near a world edge.

### Pre-conditions

- movement is already authoritative and deterministic
- parent circles still carry a derived radius for rendering and other transitional semantics
- embodied feeding and embodied contact are already in force

### Trigger

A moving parent circle reaches a world boundary.

### Success Outcome

- the parent core stays inside the world
- the boundary clamp uses the visible parent-body size rather than grown derived radius
- circles with more children can approach the wall just as closely as circles with fewer children, unless an attached child or later slice changes that explicitly

### Failure Or Rejection Cases

- if grown derived radius still silently decides how close a parent can get to the wall, movement remains less embodied than food and contact
- if this slice also changes orbit layout or rendering semantics, scope is exceeded
- if boundary behavior becomes non-deterministic, inspectability weakens

## Main Business Rules

1. Parent movement remains authoritative server-side behavior.
2. World-edge clamping uses the visible parent-core body, not derived grown radius.
3. Player and autonomous parent circles follow the same updated boundary rule.
4. Attached-child orbit positions may remain unchanged in this slice, even if they can visually extend beyond the parent body.
5. Food, contact, fight, reproduction, continuity, and steering remain unchanged.
6. Radius may remain in snapshots and rendering for now.

## Minimal Domain Concepts In Scope

- `Parent Body`
- `World Boundary`
- `Deterministic Movement Clamp`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- parent cores use a fixed visible body size for world-edge clamping
- grown derived radius no longer silently pushes parent cores away from walls
- attached children and visual radius stay transitional for now

This avoids a larger redesign of visual scale or orbit geometry while still removing one remaining hidden growth shortcut from movement.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through ordinary snapshots:

- existing parent positions
- existing radius and child count fields

Build should extend the contract only if the changed boundary behavior is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side movement clamping that uses visible parent-body size
- deterministic tests that show child-derived radius no longer changes how close the parent core can move to a wall
- client rendering that remains sufficient to observe the changed boundary position

## Build Guidance

- prefer evolving the current movement clamp helpers rather than adding a new movement system
- keep attached-child orbit layout unchanged in this slice
- preserve current food, contact, fight, reproduction, and continuity logic
- avoid coupling this slice to visual resizing or camera changes

## Initial Test Plan

### Server tests

- a base parent body still stays inside bounds
- a parent with children can move just as close to the wall as a parent without children
- autonomous circles also follow the same embodied boundary rule

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives boundary-clamped snapshots where parent position reflects visible-body clamping rather than grown derived radius

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. move a parent circle toward a world edge
2. the authoritative server clamps movement at the boundary
3. the clamp is based on the visible parent body, not on grown derived radius
4. later snapshots show the parent closer to the wall than it would have been under the old derived-radius rule

## Done Criteria

- parent movement no longer uses derived radius for world-edge clamping
- player and autonomous circles follow the same embodied boundary rule
- rendering and orbit layout remain unchanged
- food, contact, fight, reproduction, continuity, and steering remain unchanged
- tests cover the removed boundary shortcut

## Out Of Scope Follow-Ups

- removing radius from rendering
- changing attached-child orbit radius
- making attached children collide with world bounds independently
- wall damage or bouncing rules
- camera or viewport changes
