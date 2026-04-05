# Slice: Initial Created Child Ownership Identity In Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where successful reproduction remains inspectable through explicit created-child ownership identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution and child allocation.

This slice narrows one remaining reproduction inspectability gap by exposing which side received each newly created child during a successful reproduction outcome.

## Discovery Scope

Establish the smallest deterministic change that makes created-child ownership explicit:

- when reproduction resolves successfully, the same authoritative child-creation and redistribution path should remain exactly as it is today
- the runtime should expose which created child IDs were allocated to the source side and which were allocated to the target side
- creation identity, feasibility, blocked-capacity identity, payment-child identity, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of redistribution
- alternate child-creation counts
- historical reproduction logs
- new client-side effects beyond what is necessary to read created-child ownership

## Why This Slice Next

Recent slices made successful reproduction much more inspectable:

- successful reproduction now exposes which new child IDs were created
- child-paid reproduction now exposes which concrete child was consumed as payment
- blocked reproduction now exposes which side lacked enough current capacity

But successful reproduction still stops one step short of the redistribution result itself. The runtime can now say which child IDs were created, yet it still does not say which side actually received each created child. That weakens inspectability because the authoritative server already knows exactly how the deterministic redistribution rule assigned those new children, while the client and tests still need to infer ownership by diffing attached-child sets.

The next model pressure is therefore not to change reproduction semantics. It is to make created-child ownership explicit in the same bounded way creation identity and payment identity are now explicit.

This slice is the narrowest next step because it:

- changes only successful reproduction inspectability
- preserves the current feasibility rule, payment rule, creation rule, and redistribution rule
- avoids inventing staged reproduction or ancestry systems
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeCreatedChildOwnershipDuringReproduction`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction resolves as `reproduce_resolved` or `reproduce_paid_child`.

### Pre-conditions

- successful reproduction already creates new children deterministically
- created-child identity is already exposed
- redistribution across the pair already follows the current deterministic rule

### Trigger

A different-shape interaction resolves successfully as reproduction.

### Success Outcome

- reproduction still resolves as before
- the authoritative interaction outcome explicitly exposes which created child IDs were allocated to the source side and which were allocated to the target side

### Failure Or Rejection Cases

- if successful reproduction still hides which side received each created child, inspectability remains incomplete
- if exposing created-child ownership changes allocation behavior, scope is exceeded
- if ownership identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Child allocation remains authoritative server-side behavior.
2. Successful reproduction still allocates children through the existing deterministic redistribution rule.
3. The runtime explicitly exposes which created child IDs were allocated to the source side and which were allocated to the target side.
4. Creation identity, feasibility, blocked-capacity identity, cost, and payment-child identity remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Created Child`
- `Created Child Ownership`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- successful reproduction still resolves exactly as it does today
- the same deterministic child-creation and redistribution path is used as before
- snapshots expose created child ownership only for successful reproduction outcomes
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction redesign while still making the current redistribution rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- successful reproduction already exposes the created child IDs
- post-resolution attached-child arrays already show the resulting owners
- but the contract still does not say which side received each created child

Build should therefore make one minimal contract extension that exposes created-child ownership without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction path that can surface created-child ownership directly from the existing redistribution rule
- deterministic tests that prove the exposed ownership matches the actual newly attached children on each side
- client rendering or HUD output that remains sufficient to read created-child ownership

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second reproduction-event system
- reuse the same deterministic allocation path already used by successful reproduction
- keep the contract addition minimal and explicit
- avoid adding speculative staging, event logs, or ancestry structures

## Initial Test Plan

### Server tests

- energy-paid reproduction exposes which created children were allocated to source and target
- child-paid reproduction also exposes which created children were allocated to source and target
- the exposed ownership matches the actual newly attached children after resolution
- blocked reproduction exposes no created-child ownership

### Contract tests

- the snapshot schema accepts the new created-child ownership field or fields

### Integration tests

- the client receives successful reproduction snapshots that include created-child ownership identity

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches a successful reproduction outcome
2. the existing deterministic rule creates and redistributes new attached children
3. reproduction still resolves normally
4. the snapshot explicitly identifies which created child IDs were allocated to the source side and which were allocated to the target side

## Done Criteria

- successful reproduction still resolves with the same current gameplay behavior
- the outcome exposes created-child ownership identity
- tests prove the exposed ownership matches the actual newly attached children on each side
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction thresholds or costs
- staged reproduction systems
- blocked-interaction history logs
- ancestry or event-log redesign
