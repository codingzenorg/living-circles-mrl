# Slice: Initial Embodied Child Orbit Distance Without Radius

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for attached-child positions that now avoids derived parent radius as a hidden orbit-distance input

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for attached-child layout.

This slice narrows one remaining radius shortcut in the visible child model by making orbit distance depend on the visible parent body rather than grown derived radius.

## Discovery Scope

Establish the smallest deterministic change that makes attached-child geometry better match the current embodied model:

- attached-child orbit distance should use the visible parent-core body rather than grown derived radius
- attached-child orbit identity, motion, and ordering remain deterministic
- current food, contact, fight, reproduction, continuity, movement, and steering rules remain unchanged

This slice does **not** attempt to implement:

- removal of radius from rendering
- removal of radius from snapshots
- removal of child-count-based visual size growth
- independent child collisions against world bounds
- promotion of attached children into free circles

## Why This Slice Next

Recent slices already removed derived radius from several hidden leverage paths:

- food collection no longer uses grown derived radius as silent reach
- circle contact no longer uses grown derived radius as silent encounter reach
- same-shape fight resolution no longer uses derived radius as a tie-break
- movement boundaries no longer use grown derived radius as a hidden wall clamp

But attached-child layout still uses the grown parent radius to decide how far children orbit from the parent core. That means a parent with more children still pushes all visible child bodies farther outward through the old growth abstraction, even though other embodied slices have already shifted feeding, contact, fights, and movement toward visible-body geometry.

This slice is the narrowest next step because it:

- changes only attached-child orbit distance
- keeps attached-child existence, ownership, and motion intact
- preserves all current interaction and energy semantics
- removes one more hidden radius effect without forcing a broader redesign of rendering or growth

## Use-Case Contract

### Use Case

`LayOutAttachedChildrenByVisibleParentBody`

### Primary Actor

Any parent circle that owns one or more attached children.

### Pre-conditions

- attached children are already visible and authoritative
- child orbit positions are already deterministic
- parent circles still carry a derived radius for rendering and transitional growth semantics

### Trigger

A world snapshot is produced for a parent circle with attached children.

### Success Outcome

- attached-child positions are laid out from the visible parent-core body plus the existing orbit gap
- child count no longer silently pushes orbit distance outward through grown derived radius
- later snapshots still show deterministic orbit motion for the same attached children

### Failure Or Rejection Cases

- if grown derived radius still silently decides orbit distance, the visible child model remains partly abstract
- if this slice changes contact, feeding, or continuity semantics, scope is exceeded
- if orbit layout becomes unstable or non-deterministic, inspectability weakens

## Main Business Rules

1. Attached-child layout remains authoritative server-side behavior.
2. Orbit distance uses the visible parent-core body, not grown derived radius.
3. Orbit identity, slot assignment, and motion remain deterministic.
4. Player and autonomous parent circles follow the same updated orbit-distance rule.
5. Food, contact, fight, reproduction, continuity, movement, and steering remain unchanged.
6. Parent `radius` may remain in snapshots and rendering for now.

## Minimal Domain Concepts In Scope

- `Parent Body`
- `Attached Child`
- `Orbit Distance`
- `Deterministic Snapshot Layout`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- attached children orbit from the fixed visible parent-core body
- the existing orbit gap remains unchanged
- visible child positions no longer expand outward just because derived parent radius grew

This avoids a larger redesign of growth or rendering while still making the orbiting-child model more embodied.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through ordinary snapshots:

- existing parent positions
- existing parent radius
- existing attached-child positions

Build should extend the contract only if the changed orbit distance is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side attached-child layout that uses the visible parent-core body for orbit distance
- deterministic tests that show child-derived parent radius no longer changes orbit distance
- client rendering that remains sufficient to observe the tighter orbit layout

## Build Guidance

- prefer evolving the current attached-child layout helper rather than introducing a new child subsystem
- keep orbit motion deterministic and tick-driven
- preserve current embodied food, contact, and movement rules
- avoid coupling this slice to any rendering redesign or removal of radius from snapshots

## Initial Test Plan

### Server tests

- a parent with children still exposes deterministic attached-child positions
- a parent with more children does not push orbit distance outward solely through grown derived radius
- autonomous parents follow the same orbit-distance rule

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives snapshots where attached-child positions remain dynamic but are anchored to visible parent-body distance rather than grown radius

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent owns one or more attached children
2. the authoritative server computes attached-child positions for a snapshot
3. orbit distance is based on the visible parent body plus the existing orbit gap
4. later snapshots still show deterministic orbit motion without derived-radius orbit expansion

## Done Criteria

- attached-child orbit distance no longer uses grown derived parent radius
- player and autonomous parents follow the same embodied orbit rule
- deterministic orbit motion remains intact
- food, contact, fight, reproduction, continuity, movement, and steering remain unchanged
- tests cover the removed orbit-distance shortcut

## Out Of Scope Follow-Ups

- removing radius from rendering
- removing child-count-based parent size growth
- making attached children obey independent world-boundary clamps
- redesigning reproduction distribution or continuity
- turning attached children into free active circles
