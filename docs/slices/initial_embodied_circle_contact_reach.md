# Slice: Initial Embodied Circle Contact Reach

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible encounter initiation that now relies more directly on embodied parent and attached-child geometry

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for interaction initiation.

This slice narrows one current shortcut in encounter detection by making circle-to-circle contact depend on visible parent-core and attached-child bodies rather than enlarged derived parent radius alone.

## Discovery Scope

Establish the smallest deterministic change that makes encounter initiation better match the current orbiting-child embodiment:

- parent-level fight and reproduction contact should be resolved from visible parent-core overlap and attached-child overlap, not from enlarged parent radius alone
- current child-triggered contact paths remain in force
- current fight, reproduction, continuity, steering, food targeting, and embodied food collection rules remain unchanged after contact is detected

This slice does **not** attempt to implement:

- removal of radius growth from fight winner selection
- removal of radius growth from reproduction or visibility
- new combat rules
- detached child combat or detached child reproduction
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now makes the energy loop more embodied by resolving food collection from visible parent-core and attached-child bodies. But encounter initiation still uses enlarged parent radius for parent-body contact, which means fights and reproduction can still start because of a hidden grown reach even when no visible body overlaps.

The model pressure is now:

- encounter initiation should become as embodied as food collection already is
- visible orbiting children should remain the main bridge from abstract child count toward contact geometry
- the repository should reduce one more transitional radius shortcut without changing downstream outcome logic

This slice is the narrowest next step because it:

- focuses only on interaction initiation geometry
- keeps current fight and reproduction resolution rules unchanged
- preserves existing child-triggered parent and child-to-child contact paths
- avoids a larger attempt to remove radius-derived abstraction from every system at once

## Use-Case Contract

### Use Case

`ResolveCircleContactFromEmbodiedBodies`

### Primary Actor

Any player or autonomous circle pair whose bodies may overlap during a simulation tick.

### Pre-conditions

- encounter detection is already authoritative and deterministic
- attached-child positions are already authoritative and deterministic per tick
- child-triggered contact already exists
- derived radius growth still exists as a current property

### Trigger

A simulation tick checks whether two circles should enter fight or reproduction contact.

### Success Outcome

- fight or reproduction starts only when visible parent-core or attached-child bodies overlap
- child-triggered contact remains possible and inspectable
- encounter initiation becomes more legible as an embodied interaction instead of a hidden enlarged parent reach effect

### Failure Or Rejection Cases

- if enlarged parent radius still silently determines contact, the embodied child model remains partially undermined
- if removing the shortcut breaks determinism or child-triggered contact, slice scope is exceeded
- if this slice also changes fight winner selection or reproduction payment, scope has drifted

## Main Business Rules

1. Contact detection remains authoritative server-side behavior.
2. Parent-body contact uses the visible parent-core body rather than enlarged derived radius.
3. Existing attached-child-to-parent and attached-child-to-attached-child contact paths remain valid.
4. Current fight and reproduction outcome rules remain unchanged after contact is detected.
5. Current steering rules remain unchanged.
6. The rule must remain deterministic for the same world state and tick.
7. Food collection, regeneration, and child-payment rules remain unchanged.
8. Player and autonomous circles follow the same embodied contact rule.

## Minimal Domain Concepts In Scope

- `Parent Core Body`
- `Attached Child Body`
- `Embodied Contact`
- `Fight Or Reproduction Trigger`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- encounter initiation checks visible parent-core and attached-child body overlap only
- parent radius may remain a visual and other-domain property for now, but it no longer silently enlarges contact reach
- child-triggered contact remains part of the same authoritative contact path
- no new fight, reproduction, or continuity semantics are introduced

This avoids the larger step of removing all radius-derived shortcuts while still moving one more core loop toward the embodied child model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- unchanged world snapshots
- visible parent and attached-child positions
- existing `contact_origin` once interaction occurs

Build should extend the contract only if embodied non-contact is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side contact detection that can distinguish parent-core and attached-child overlap from abstract radius reach
- deterministic tests that show enlarged derived radius alone no longer starts contact
- client rendering that remains sufficient to observe embodied contact

## Build Guidance

- prefer evolving the current encounter-detection path rather than adding a new subsystem
- preserve current fight and reproduction resolution logic
- keep steering, food, and continuity semantics unchanged
- avoid changing visual radius or winner ordering in this slice

## Initial Test Plan

### Server tests

- parent-core overlap still starts same-shape fight or different-shape reproduction
- attached-child overlap still starts same-shape fight or different-shape reproduction
- enlarged derived radius alone no longer starts contact when no visible body overlaps
- contact is still resolved exactly once per pair when multiple visible overlap paths exist

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives snapshots showing ordinary embodied contact by parent and child bodies
- the client no longer receives contact when only abstract derived parent radius would have explained it

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a pair advances into a near-contact state
2. contact is checked against visible parent-core and attached-child bodies
3. if a visible body overlaps, the ordinary fight or reproduction path begins
4. if only abstract enlarged parent radius would have reached, no interaction begins
5. once real contact occurs, the ordinary resolution path remains unchanged

## Done Criteria

- encounter initiation is resolved from visible parent and child bodies
- abstract derived parent radius no longer silently enlarges contact reach
- the rule is deterministic and documented
- current fight and reproduction resolution semantics remain unchanged
- tests cover embodied contact and non-contact cases

## Out Of Scope Follow-Ups

- removing radius from fight winner ordering
- removing radius from visual presentation
- changing reproduction resolution
- changing continuity rules
- detached child combat or reproduction
