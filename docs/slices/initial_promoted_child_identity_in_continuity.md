# Slice: Initial Promoted Child Identity In Continuity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where continuity remains inspectable through both promoted position and promoted child identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for death resolution and continuity.

This slice narrows one remaining continuity inspectability gap by exposing which attached child was promoted when continuity occurs.

## Discovery Scope

Establish the smallest deterministic change that makes promoted continuity identity explicit:

- when continuity occurs, the promoted child should remain deterministically selected as today
- the runtime should expose which child was consumed as the promotion source
- lineage, generation increment, replacement energy, and promoted-position behavior remain unchanged
- fight absorption, reproduction payment, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- multi-child promotion choice systems
- detached child entities during normal play
- redesign of continuity eligibility
- redesign of fight or reproduction rules
- new client-side effects beyond what is necessary to read the promoted child identity

## Why This Slice Next

Recent slices already made continuity much more embodied:

- attached children are the authoritative visible child state
- continuity explicitly consumes one attached child
- the continuing active parent now reappears at the promoted child's visible position

But the runtime still does not say which child was promoted. That weakens inspectability because the server knows the exact child identity, yet the client and tests can only infer promotion indirectly from position.

The next model pressure is therefore not a new gameplay rule. It is to make continuity identity explicit in the same way contact origin and child-paid reproduction were made explicit earlier.

This slice is the narrowest next step because it:

- changes only continuity inspectability
- preserves current continuity behavior
- avoids inventing new ancestry or promotion-choice rules
- gives build a deterministic contract extension instead of a larger semantic redesign

## Use-Case Contract

### Use Case

`ExposePromotedChildIdentityDuringContinuity`

### Primary Actor

Any player or autonomous parent circle that preserves continuity by consuming one attached child.

### Pre-conditions

- continuity already consumes exactly one attached child
- the promoted child is already selected deterministically
- continuity already repositions the active parent to the promoted child's visible position

### Trigger

A fight defeat or zero-energy collapse resolves into continuity.

### Success Outcome

- one attached child is consumed as before
- the continuing active parent keeps the same lineage and incremented generation
- the continuing active parent remains at the promoted child's visible position
- the authoritative outcome explicitly exposes the promoted child's identity

### Failure Or Rejection Cases

- if continuity still hides which child was promoted, inspectability remains incomplete
- if exposing promoted identity changes continuity behavior, scope is exceeded
- if the promoted child identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Continuity remains authoritative server-side behavior.
2. One attached child is still consumed as the continuity source.
3. The continuing active parent still takes the promoted child's current visible position.
4. The continuity outcome explicitly exposes the promoted child's identity.
5. Lineage preservation, generation increment, and replacement energy reset remain unchanged.
6. Fight absorption, reproduction payment, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Continuity`
- `Promoted Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- continuity still promotes exactly one child
- the same deterministic child is consumed as before
- snapshots expose that child's ID through the existing interaction channel or an equally small contract extension
- all broader continuity semantics remain unchanged

This avoids turning continuity into a larger lineage or ancestry system while still making the rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `death_promoted_child` says that continuity occurred
- parent position says where the promoted child was
- but the contract still does not say which attached child was consumed

Build should therefore make one minimal contract extension that exposes promoted child identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side continuity resolution that can surface the chosen promoted child ID
- deterministic tests that prove the exposed child ID matches the consumed child and promoted position
- client rendering or HUD output that remains sufficient to read the new continuity identity

## Build Guidance

- prefer extending the existing continuity interaction payload rather than creating a second continuity event system
- reuse the same deterministic child selection already used by continuity
- keep the contract addition minimal and explicit
- avoid adding speculative ancestry structures or historical event logs

## Initial Test Plan

### Server tests

- zero-energy player continuity exposes the promoted child ID
- zero-energy autonomous continuity exposes the promoted child ID
- the exposed child ID matches the child whose visible position became the continuing active position

### Contract tests

- the snapshot schema accepts the new promoted-child continuity field

### Integration tests

- the client receives a `death_promoted_child` snapshot that includes the promoted child identity

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent with attached children reaches a continuity path
2. one attached child is deterministically selected and consumed
3. the active continuing parent reappears at that child's visible position
4. the continuity outcome explicitly identifies that promoted child in the snapshot

## Done Criteria

- continuity still uses the promoted child's visible position
- the continuity outcome exposes the promoted child identity
- tests prove the exposed ID matches the consumed child and promoted position
- lineage, generation, and replacement energy remain unchanged

## Out Of Scope Follow-Ups

- promotion choice among multiple children
- detached child entities during normal play
- continuity mutation or inheritance systems
- redesigning replacement energy
- redesigning fight or reproduction rules
