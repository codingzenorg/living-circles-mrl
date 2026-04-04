# Slice: Initial Embodied Food Collection Reach

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible food collection that now relies more directly on embodied parent and attached-child geometry

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for food collection outcomes.

This slice narrows one current shortcut in the feeding model by making food collection rely on visible bodies rather than abstract accumulated radius reach alone.

## Discovery Scope

Establish the smallest deterministic change that makes food collection better match the current orbiting-child embodiment:

- food collection should be resolved from parent-core overlap and attached-child overlap, not from enlarged parent radius alone
- current food targeting, food regeneration, energy gain, movement, and attached-child collection rules remain in force
- current fight, reproduction, continuity, and social steering rules remain unchanged

This slice does **not** attempt to implement:

- removal of radius growth from combat or visuals
- a general removal of all count-based shortcuts
- new food types or nutrition systems
- detached child foraging
- client-side anticipation or prediction

## Why This Slice Next

The current implementation has already made attached children matter for:

- actual food collection
- food targeting
- contact initiation
- avoidance
- positive interaction seeking
- conflict absorption
- continuity
- reproduction payment visibility

But food collection still also benefits from the abstract radius shortcut derived from `children_count`. That leaves the energy loop partially split between embodied visible bodies and a hidden reach expansion. The model pressure is now to tighten one part of that shortcut without disturbing the broader system.

This slice is the narrowest next step because it:

- focuses only on food collection resolution
- keeps current food targeting and child-based collection behavior
- reduces one embodied-versus-abstract mismatch without changing combat or reproduction
- avoids a larger attempt to remove all radius-based shortcuts at once

## Use-Case Contract

### Use Case

`ResolveFoodCollectionFromEmbodiedBodies`

### Primary Actor

Any active player or autonomous circle that overlaps a food slot through its parent body or attached children.

### Pre-conditions

- food collection is already authoritative and deterministic
- attached-child positions are already authoritative and deterministic per tick
- attached children can already collect food on behalf of their parent
- radius growth still exists as a current derived property

### Trigger

A simulation tick checks whether a food slot should be consumed.

### Success Outcome

- a food slot is consumed only when it overlaps the parent core body or one attached child body
- attached-child collection remains possible and visible
- food collection becomes more legible as an embodied interaction instead of a hidden enlarged reach effect

### Failure Or Rejection Cases

- if enlarged radius still silently determines food reach, the embodied child model remains partially undermined
- if removing the shortcut breaks determinism or double-consumes food, slice scope is exceeded
- if this slice also changes fight or reproduction reach, scope has drifted

## Main Business Rules

1. Food collection remains authoritative server-side behavior.
2. Food is consumed through overlap with visible parent-core or attached-child bodies.
3. Current attached-child food collection behavior remains valid.
4. Current food energy gain and regeneration rules remain unchanged.
5. Current food targeting rules remain unchanged.
6. The rule must remain deterministic for the same world state and tick.
7. Fight, reproduction, continuity, and social steering remain unchanged.
8. Player and autonomous circles follow the same embodied food collection rule.

## Minimal Domain Concepts In Scope

- `Food Collection`
- `Parent Core Body`
- `Attached Child Body`
- `Embodied Reach`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- food collection checks visible body overlap only
- parent radius may remain a visual and other-domain property for now, but it no longer silently enlarges food reach
- attached-child collection remains part of the same authoritative collection path
- no new event stream or food provenance model is introduced

This avoids the larger step of removing all radius-derived shortcuts while still moving one core loop toward the embodied child model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- unchanged food snapshots
- visible parent and attached-child positions
- ordinary energy changes after collection

Build should extend the contract only if embodied collection is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side food collection resolution that can distinguish parent-core and attached-child overlap from abstract radius reach
- deterministic tests that show a large derived radius alone no longer consumes food
- client rendering that remains sufficient to observe embodied collection

## Build Guidance

- prefer evolving the current food-overlap resolution path rather than adding a new subsystem
- preserve the current food targeting, energy gain, and regeneration rules
- keep movement, fight, reproduction, and continuity semantics unchanged
- avoid changing visual radius or combat reach in this slice

## Initial Test Plan

### Server tests

- a parent core body overlapping food still consumes it
- an attached child overlapping food still consumes it
- a circle with enlarged derived radius but no body overlap does not consume the food
- food is still consumed exactly once when parent and child overlap the same slot

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives snapshots showing ordinary embodied food collection by parent and child bodies
- the client no longer receives collection when only the abstract derived radius would have explained it

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a circle advances into a state near a food slot
2. the food slot is checked against visible parent and attached-child bodies
3. if either visible body overlaps, the slot is consumed
4. if only abstract enlarged radius would have reached the slot, the slot remains
5. energy and regeneration continue through the ordinary path

## Done Criteria

- food collection is resolved from visible parent and child bodies
- abstract derived radius no longer silently enlarges food reach
- the rule is deterministic and documented
- current food targeting and downstream food semantics remain unchanged
- tests cover embodied collection and non-collection cases

## Out Of Scope Follow-Ups

- removing radius shortcuts from combat
- removing radius shortcuts from reproduction
- changing rendered parent size
- new food systems
- detached child foraging
