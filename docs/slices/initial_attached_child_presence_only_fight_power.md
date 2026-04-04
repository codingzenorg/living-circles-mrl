# Slice: Initial Attached-Child Presence Only Fight Power

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where same-shape fight outcomes depend on visible child presence rather than abstract child-count magnitude

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for fight resolution.

This slice narrows one remaining abstract child shortcut by changing same-shape fight power from raw child-count advantage to a smaller embodied rule based only on whether attached children are present.

## Discovery Scope

Establish the smallest deterministic change that makes child-based fight power less abstract:

- same-shape fights should still use energy first
- if energy ties, fight resolution should distinguish only between “has at least one attached child” and “has no attached children”
- exact ties after that remain deterministic
- current fight absorption, reproduction payment, continuity, feeding, contact, movement, orbit, and steering rules remain unchanged

This slice does **not** attempt to implement:

- per-child body hitboxes as independent damage sources
- child-position-based combat scoring
- redesign of reproduction payment
- redesign of continuity
- detached child entities

## Why This Slice Next

The model has already tightened most child-related shortcuts:

- attached children are the visible child state
- geometry and contract now follow attached-child bodies rather than mirrored child count
- visible parent growth is gone

But same-shape fight resolution still uses raw child-count magnitude as a direct combat advantage. That means the combat model still treats children as an abstract pile of strength rather than primarily as visible dependents that can already absorb one loss and participate in contact.

The next model pressure is to reduce that abstraction without opening a full combat redesign. The narrowest step is to let child presence matter, but stop letting larger raw child counts stack unlimited direct fight advantage.

This slice is the smallest coherent step because it:

- preserves energy as the main fight determinant
- keeps children relevant in fights
- removes magnitude-based child stacking from winner selection
- avoids inventing a new combat subsystem

## Use-Case Contract

### Use Case

`DetermineSameShapeFightWinnerFromEnergyAndChildPresence`

### Primary Actor

Any same-shape circle pair that reaches authoritative fight resolution.

### Pre-conditions

- same-shape fight resolution is already authoritative and deterministic
- attached children are already the authoritative child state
- current loser absorption and continuity rules already consume attached children explicitly

### Trigger

A same-shape contact resolves into a fight.

### Success Outcome

- higher energy still wins first
- if energy ties, a circle with one or more attached children beats a circle with none
- if both circles either have attached children or both do not, exact ties remain deterministic
- later snapshots show fight results that rely less on abstract count magnitude

### Failure Or Rejection Cases

- if raw child-count magnitude still directly decides the winner, combat remains more abstract than the embodied child model
- if this slice changes child-loss absorption or continuity semantics, scope is exceeded
- if winner ordering becomes opaque or non-deterministic, inspectability weakens

## Main Business Rules

1. Same-shape fight resolution remains authoritative server-side behavior.
2. Higher energy remains the first fight winner criterion.
3. Child presence remains a fight advantage only as a boolean distinction: has attached child versus does not.
4. Raw child-count magnitude no longer directly stacks fight power beyond that presence distinction.
5. Exact ties remain deterministic.
6. Fight absorption, reproduction payment, continuity, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Same-Shape Fight`
- `Energy Advantage`
- `Attached Child Presence`
- `Deterministic Exact Tie`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- energy first
- then child presence
- then deterministic exact tie

This reduces abstract child stacking while preserving the existing idea that children still help a circle in conflict.

## Required Runtime Contract Changes

The current contract is likely sufficient because:

- attached-child arrays are already visible
- fight outcomes are already visible

Build should extend the contract only if the reduced child-power rule is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side fight winner selection that reads child presence from attached children
- deterministic tests that show extra child count no longer stacks direct fight advantage
- client rendering that remains sufficient to observe the changed outcome

## Build Guidance

- prefer evolving the current fight winner helpers rather than adding a new combat layer
- preserve existing loser absorption and continuity paths
- keep current snapshot and rendering behavior unchanged
- avoid turning this slice into a full combat rebalance

## Initial Test Plan

### Server tests

- higher energy still beats lower energy regardless of child presence
- equal-energy fights let “has child” beat “no child”
- equal-energy fights where both sides have children no longer use larger child count to decide the winner
- exact ties after energy and child presence remain deterministic

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives a same-shape fight outcome where a larger child count no longer changes the winner once both sides already have at least one child

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two same-shape circles make authoritative contact
2. their energy is compared
3. if energy ties, only child presence versus absence is compared
4. if both have child presence or both do not, the ordinary deterministic exact tie rule decides
5. the ordinary fight outcome path continues from that winner

## Done Criteria

- same-shape fights no longer use raw child-count magnitude as direct fight power
- energy remains the first fight determinant
- child presence still matters in fights
- exact ties remain deterministic
- absorption, continuity, and payment semantics remain unchanged
- tests cover the removed child-stacking shortcut

## Out Of Scope Follow-Ups

- per-child combat geometry
- redesigning fight absorption
- redesigning reproduction payment
- redesigning continuity
- detached child entities
