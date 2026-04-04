# Slice: Initial Remove `children_count` From Runtime Contract

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract that uses attached children as the sole child representation

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for child state and the client deriving readable child quantity from attached children.

This slice narrows the remaining contract duplication by removing `children_count` from runtime snapshots now that attached children are already the sole authoritative child state internally.

## Discovery Scope

Establish the smallest deterministic change that makes the runtime contract match the current child model:

- `children_count` is removed from runtime snapshots
- the client derives child quantity from `attached_children`
- current gameplay behavior remains unchanged
- current attached-child positions, ownership, and motion remain unchanged

This slice does **not** attempt to implement:

- redesign of child-based fight power
- redesign of reproduction payment
- redesign of continuity
- detached child entities
- broader contract redesign beyond child-count duplication

## Why This Slice Next

The previous slice already made attached children the single authoritative child state inside the simulation. That means the runtime contract still carries one remaining duplication:

- `attached_children` as the real child state
- `children_count` as a mirrored convenience field

That duplication still costs contract surface and invites drift between what is visible and what is counted, even if the count is now derived correctly. The next model pressure is to make the contract itself tell the same story as the simulation: children are visible entities, and counts are derived by readers if needed.

This slice is the narrowest next step because it:

- changes only contract shape and client interpretation
- leaves gameplay rules untouched
- reduces one more mirrored representation
- keeps inspectability through visible child bodies

## Use-Case Contract

### Use Case

`RenderChildStateWithoutMirroredCountField`

### Primary Actor

The browser client and any contract consumer reading authoritative snapshots.

### Pre-conditions

- attached children are already the sole authoritative child state in the simulation
- snapshots still expose both `attached_children` and `children_count`
- the client currently has enough information to derive child quantity from attached children

### Trigger

A world snapshot is produced or consumed.

### Success Outcome

- snapshots expose child state only through attached-child arrays
- client-side labels or HUD elements derive counts from attached children where needed
- visible behavior remains unchanged

### Failure Or Rejection Cases

- if the contract still mirrors child state as both bodies and count, duplication remains
- if removing `children_count` makes the client or tests ambiguous, inspectability weakens
- if gameplay logic changes along with the contract change, scope is exceeded

## Main Business Rules

1. Attached-child arrays are the sole runtime representation of child state.
2. `children_count` is removed from authoritative snapshots.
3. Contract consumers derive child quantity from `attached_children`.
4. Player and autonomous circles follow the same updated contract rule.
5. Fight power, reproduction payment, continuity, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Derived Child Quantity`
- `Runtime Contract`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- contract child state is just attached children
- counts become reader-derived
- gameplay and visible geometry remain the same

This avoids a broader semantic rewrite while aligning the runtime contract with the current child model.

## Required Runtime Contract Changes

- remove `children_count` from snapshot payloads
- update schema/tests to treat attached-child arrays as sufficient
- update client usage to derive count from attached children

## Required Ports Or Boundaries

- server-side snapshot serialization without `children_count`
- contract tests updated to the reduced snapshot shape
- client rendering that derives readable child counts from attached children

## Build Guidance

- prefer deriving child quantity locally in the client rather than adding replacement count fields
- keep snapshot readability through labels/HUD logic rather than mirrored payload fields
- preserve current behaviors and outcomes
- avoid using this slice to redesign child semantics

## Initial Test Plan

### Server tests

- snapshots still expose attached children correctly for player and autonomous circles
- child-dependent outcomes remain unchanged after the contract reduction

### Contract tests

- schema no longer requires `children_count`
- attached-child arrays remain explicit

### Integration tests

- the client still renders and labels child state correctly without `children_count`

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent owns attached children
2. the authoritative server emits a snapshot without `children_count`
3. the client derives readable child quantity from `attached_children`
4. visible behavior remains unchanged

## Done Criteria

- `children_count` is removed from runtime snapshots
- client rendering derives child quantity from attached children
- gameplay behavior remains unchanged
- tests cover the reduced contract shape

## Out Of Scope Follow-Ups

- redesigning child-based fight power
- redesigning reproduction payment
- redesigning continuity
- detached child entities
- broader contract redesign beyond child-count duplication
