# Slice: Initial Absorbed Child Identity In Fight

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where absorbed fight loss remains inspectable through explicit child identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for fight resolution.

This slice narrows one remaining fight inspectability gap by exposing which attached child absorbed a same-shape loss when `fight_absorbed_child` occurs.

## Discovery Scope

Establish the smallest deterministic change that makes absorbed-child identity explicit:

- when `fight_absorbed_child` occurs, the absorbed child should remain deterministically selected as today
- the runtime should expose which child was consumed as the absorption source
- fight winner selection, child-loss behavior, continuity, reproduction, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new combat rules
- multi-child absorption strategies
- redesign of winner selection
- redesign of continuity or reproduction behavior
- new client-side effects beyond what is necessary to read absorbed-child identity

## Why This Slice Next

Recent slices already made child-based loss much more inspectable:

- attached children are the authoritative visible child state
- `fight_absorbed_child` explicitly distinguishes child absorption from full defeat
- continuity now exposes which promoted child was consumed

But fight absorption still hides which child was lost. That weakens inspectability because the server knows exactly which attached child was consumed, yet the client and tests can only infer the loss indirectly from the remaining child set and orbit positions.

The next model pressure is therefore not to change combat semantics. It is to make absorbed-child identity explicit in the same way continuity identity was made explicit.

This slice is the narrowest next step because it:

- changes only fight inspectability
- preserves the current fight rule and child-loss behavior
- avoids inventing new combat or ancestry systems
- gives build a deterministic contract extension instead of a larger semantic redesign

## Use-Case Contract

### Use Case

`ExposeAbsorbedChildIdentityDuringFight`

### Primary Actor

Any player or autonomous circle that remains active by losing one attached child during same-shape conflict.

### Pre-conditions

- same-shape fight resolution already supports `fight_absorbed_child`
- one attached child is already consumed deterministically from the loser
- the current interaction payload already exposes fight outcome identity at the parent level

### Trigger

A same-shape fight resolves into `fight_absorbed_child`.

### Success Outcome

- one attached child is consumed as before
- the losing parent remains active as before
- the authoritative fight outcome explicitly exposes the absorbed child's identity

### Failure Or Rejection Cases

- if `fight_absorbed_child` still hides which child was consumed, inspectability remains incomplete
- if exposing absorbed-child identity changes fight behavior, scope is exceeded
- if absorbed child identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Fight resolution remains authoritative server-side behavior.
2. `fight_absorbed_child` still consumes exactly one attached child from the loser.
3. The fight outcome explicitly exposes the absorbed child's identity.
4. Winner selection, loser identity, and parent activity outcome remain unchanged.
5. Continuity, reproduction, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Fight`
- `Absorbed Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- same-shape fight still resolves exactly as it does today
- the same deterministic child is consumed as before
- snapshots expose that absorbed child ID through the existing interaction channel or an equally small contract extension
- all broader fight semantics remain unchanged

This avoids turning the slice into a combat redesign while still making the rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `fight_absorbed_child` says that child absorption occurred
- loser identity is already visible
- but the contract still does not say which attached child was consumed

Build should therefore make one minimal contract extension that exposes absorbed child identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side fight-resolution path that can surface the chosen absorbed child ID
- deterministic tests that prove the exposed child ID matches the consumed child
- client rendering or HUD output that remains sufficient to read the absorbed-child identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second fight-event system
- reuse the same deterministic child selection already used by `fight_absorbed_child`
- keep the contract addition minimal and explicit
- avoid adding speculative damage systems or historical fight logs

## Initial Test Plan

### Server tests

- player-side child absorption exposes the absorbed child ID
- autonomous-side child absorption exposes the absorbed child ID
- the exposed child ID matches the child removed from the losing parent's attached-child set

### Contract tests

- the snapshot schema accepts the new absorbed-child fight field

### Integration tests

- the client receives a `fight_absorbed_child` snapshot that includes the absorbed-child identity

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a same-shape fight reaches the child-absorption path
2. one attached child is deterministically selected and consumed from the loser
3. the loser remains active
4. the fight outcome explicitly identifies the absorbed child in the snapshot

## Done Criteria

- `fight_absorbed_child` still consumes one child and keeps the loser active
- the fight outcome exposes the absorbed child identity
- tests prove the exposed ID matches the consumed child
- broader fight, continuity, and reproduction behavior remain unchanged

## Out Of Scope Follow-Ups

- new combat rules
- alternate child absorption strategies
- damage accumulation systems
- ancestry or historical event logs
