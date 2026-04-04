# Slice: Initial Embodied Fight Tie-Break Without Radius

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible fight outcomes that now avoids derived radius as a hidden same-shape tie-break

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for fight resolution.

This slice narrows one remaining shortcut in same-shape conflict by removing derived radius from fight winner selection once energy and child power are already considered.

## Discovery Scope

Establish the smallest deterministic change that makes fight resolution better match the current orbiting-child embodiment:

- same-shape fight winner selection should no longer depend on derived radius after energy and child count are compared
- the current energy-first and child-power ordering remains in force
- current contact, food, reproduction, continuity, steering, and visual radius rules remain unchanged

This slice does **not** attempt to implement:

- removal of radius from rendering
- removal of radius from movement boundary clamping
- new combat damage systems
- mutation or richer combat states
- client-side anticipation or prediction

## Why This Slice Next

The current implementation already tightened two major shortcuts:

- food collection no longer uses enlarged derived radius as silent reach
- circle contact no longer uses enlarged derived parent radius as silent encounter reach

But same-shape fight resolution still falls back to derived radius after energy and child count, which means one remaining hidden growth abstraction still decides combat even after child power became explicit. The model pressure is now to remove that last radius-based combat shortcut before tackling broader visual or movement implications.

This slice is the narrowest next step because it:

- changes only same-shape fight tie-breaking
- preserves the current energy-first and child-power ordering
- keeps contact, reproduction, and continuity behavior unchanged
- avoids a larger attempt to redefine size, rendering, or movement all at once

## Use-Case Contract

### Use Case

`DetermineSameShapeFightWinnerWithoutRadiusTieBreak`

### Primary Actor

Any same-shape circle pair that reaches authoritative fight resolution.

### Pre-conditions

- same-shape fight resolution is already authoritative and deterministic
- energy and child count are already explicit fight inputs
- derived radius still exists as a current visual and other-domain property

### Trigger

A same-shape contact resolves into a fight.

### Success Outcome

- higher energy still wins first
- if energy ties, higher child count still wins next
- if both energy and child count tie, winner selection no longer depends on derived radius
- later snapshots show fight results that are more grounded in visible child state and less in hidden radius-derived leverage

### Failure Or Rejection Cases

- if derived radius still silently decides fights after child power is explicit, combat remains more abstract than recent embodied slices
- if removing radius tie-break also changes contact or reproduction semantics, slice scope is exceeded
- if winner ordering becomes opaque or non-deterministic, inspectability weakens

## Main Business Rules

1. Same-shape fight resolution remains authoritative server-side behavior.
2. Higher energy remains the first fight winner criterion.
3. Higher child count remains the next fight winner criterion.
4. Derived radius no longer decides same-shape fights once energy and child count have tied.
5. Exact tie resolution must remain deterministic.
6. Contact detection, reproduction, continuity, food, and steering remain unchanged.
7. Visual radius may remain unchanged in this slice.
8. Player and autonomous circles follow the same updated fight ordering.

## Minimal Domain Concepts In Scope

- `Same-Shape Fight`
- `Energy Advantage`
- `Child Power`
- `Deterministic Exact Tie`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- energy first
- child count second
- no radius-based combat tie-break
- exact ties still resolve deterministically using the current final stable rule

This avoids a larger combat redesign while still removing one remaining hidden growth shortcut from fight resolution.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- existing energy, child count, and winner fields
- ordinary fight outcome snapshots

Build should extend the contract only if the new tie-break behavior is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side fight winner selection that no longer depends on derived radius after energy and child count
- deterministic tests that show equal-energy, equal-child fights no longer use radius to decide the winner
- client rendering that remains sufficient to observe the changed winner outcome

## Build Guidance

- prefer evolving the current fight winner helper rather than adding a new combat subsystem
- preserve current contact initiation and child-absorption behavior
- keep visual radius unchanged in this slice
- avoid introducing probabilistic or multi-phase combat

## Initial Test Plan

### Server tests

- higher energy still beats lower energy regardless of radius
- equal-energy fights still let higher child count win regardless of radius
- equal-energy, equal-child fights no longer use derived radius to choose the winner
- exact ties remain deterministic for player-versus-autonomous and autonomous-versus-autonomous cases

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives a same-shape fight outcome where previously larger radius would have won but now does not
- the client still receives unchanged outcomes when energy or child count already decides the fight

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two same-shape circles make authoritative contact
2. their energy is compared
3. if energy ties, their child counts are compared
4. if child counts also tie, winner selection no longer uses derived radius
5. the ordinary fight outcome path continues from that winner

## Done Criteria

- same-shape fights no longer use derived radius as a tie-break
- energy and child count remain deterministic fight inputs
- exact ties remain deterministic
- contact, reproduction, and continuity semantics remain unchanged
- tests cover the removed radius tie-break and unchanged higher-priority rules

## Out Of Scope Follow-Ups

- removing radius from rendering
- removing radius from movement boundaries
- redesigning combat into multi-step exchanges
- changing reproduction or continuity rules
- detached child combat
