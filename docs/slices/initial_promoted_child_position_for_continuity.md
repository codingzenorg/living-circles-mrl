# Slice: Initial Promoted Child Position For Continuity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where continuity remains visible through the promoted child’s last embodied position

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for death resolution and continuity.

This slice narrows one remaining continuity abstraction by making the continuing active parent emerge from the promoted child’s visible position instead of simply staying at the former parent-body center.

## Discovery Scope

Establish the smallest deterministic change that makes continuity more embodied:

- when continuity occurs, one attached child is still consumed
- the continuing active parent should take the promoted child’s current visible position
- lineage, generation increment, energy reset, and current continuity eligibility remain unchanged
- fight absorption, reproduction payment, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- detached child entities during normal play
- multi-child promotion choice systems
- continuity mutation or inheritance systems
- redesign of replacement energy
- redesign of fight or reproduction rules

## Why This Slice Next

Recent slices already made children much more embodied:

- attached children are the visible and authoritative child state
- continuity explicitly consumes one attached child
- the runtime contract now centers child bodies instead of mirrored counts

But continuity still resolves from the old parent body position. That means “a child replaces the dead parent” is only half embodied: a child disappears, but the surviving body does not actually emerge from that child’s visible place in the world.

The next model pressure is to make continuity use the promoted child’s last visible position as the new active-circle position. That keeps the rule bounded while making continuity read as child promotion instead of silent parent persistence.

This slice is the narrowest next step because it:

- changes only the position semantics of continuity
- preserves lineage, generation, and energy reset behavior
- keeps the one-child promotion rule intact
- improves embodied continuity without opening a larger ancestry system

## Use-Case Contract

### Use Case

`ContinueLineFromPromotedChildPosition`

### Primary Actor

Any player or autonomous parent circle that would die while still owning at least one attached child.

### Pre-conditions

- continuity is already possible when a dying parent has an attached child
- one attached child is already consumed as the continuity source
- attached-child positions are already visible and deterministic at the current tick

### Trigger

A fight defeat or zero-energy collapse resolves into continuity.

### Success Outcome

- one attached child is consumed as before
- the continuing active parent keeps the same lineage and incremented generation
- the continuing active parent takes the promoted child’s current visible position
- later snapshots make continuity look like child promotion rather than stationary parent persistence

### Failure Or Rejection Cases

- if continuity still keeps the parent centered on its former body position, the promotion remains only partially embodied
- if this slice changes continuity eligibility or energy reset, scope is exceeded
- if promotion position becomes ambiguous or non-deterministic, inspectability weakens

## Main Business Rules

1. Continuity remains authoritative server-side behavior.
2. One attached child is still consumed as the continuity source.
3. The continuing active parent takes the promoted child’s current visible position.
4. Lineage preservation and generation increment remain unchanged.
5. Replacement energy reset remains unchanged.
6. Fight absorption, reproduction payment, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Continuity`
- `Promoted Position`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- continuity still promotes exactly one child
- that child’s last visible position becomes the continuing active position
- all other continuity semantics remain the same

This avoids a larger redesign of continuity while making the existing rule more physically legible.

## Required Runtime Contract Changes

The current contract is likely sufficient because:

- attached-child positions are already visible
- lineage and generation are already visible
- continuity outcomes are already inspectable through `death_promoted_child`

Build should extend the contract only if promoted-position continuity is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side continuity resolution that can capture the promoted child’s visible position
- deterministic tests that show continuity repositions the active parent to that child position
- client rendering that remains sufficient to observe the new continuity placement

## Build Guidance

- prefer evolving the current continuity helpers rather than adding a second death system
- capture the promoted child position from the same deterministic orbit layout already used in snapshots
- preserve current lineage, generation, and energy semantics
- avoid turning this slice into a generalized spawn system

## Initial Test Plan

### Server tests

- zero-energy continuity repositions the player to the promoted child’s visible position
- zero-energy continuity repositions an autonomous circle to the promoted child’s visible position
- fight-defeat continuity uses the same promoted-position rule

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives a continuity snapshot where the surviving circle appears at the promoted child’s former position

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent with attached children reaches a death path
2. one attached child is selected and consumed as the continuity source
3. the active continuing parent reappears at that child’s current visible position
4. later snapshots show the line continuing from the promoted child position

## Done Criteria

- continuity uses the promoted child’s visible position
- lineage, generation, and energy reset stay unchanged
- fight and zero-energy continuity follow the same promoted-position rule
- tests cover the new embodied continuity placement

## Out Of Scope Follow-Ups

- detached child entities during normal play
- promotion choice among multiple children
- mutation or inheritance systems
- redesigning replacement energy
- redesigning fight or reproduction rules
