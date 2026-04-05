# Slice: Initial Created Child Identity In Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where successful reproduction remains inspectable through explicit created-child identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution and child creation.

This slice narrows one remaining reproduction inspectability gap by exposing which new attached child or children were created during a successful reproduction outcome.

## Discovery Scope

Establish the smallest deterministic change that makes reproduction creation identity explicit:

- when reproduction resolves successfully, the same authoritative child-creation and redistribution path should remain exactly as it is today
- the runtime should expose which concrete new child IDs were created by the reproduction outcome
- reproduction feasibility, blocked-capacity identity, payment-child identity, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of redistribution
- alternate child-creation counts
- historical reproduction logs
- new client-side effects beyond what is necessary to read created-child identity

## Why This Slice Next

Recent slices made the reproduction path much more inspectable:

- `reproduce_paid_child` now exposes which side paid through a child
- child-paid reproduction now exposes which concrete payment child was consumed
- `reproduce_blocked_energy` now exposes which side lacked enough current capacity

But successful reproduction still stops short of the created-child identity itself. The runtime shows the post-resolution attached-child state, yet it still does not say which concrete child IDs were newly created by the reproduction outcome. That weakens inspectability because the authoritative server already knows exactly which new child IDs it allocated, while the client and tests still need to infer them indirectly by diffing attached-child sets.

The next model pressure is therefore not to change reproduction semantics. It is to make created-child identity explicit in the same bounded way consumed-child identity and blocked-capacity identity are now explicit.

This slice is the narrowest next step because it:

- changes only successful reproduction inspectability
- preserves the current feasibility rule, payment rule, and redistribution rule
- avoids inventing staged reproduction or ancestry systems
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeCreatedChildIdentityDuringReproduction`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction resolves as `reproduce_resolved` or `reproduce_paid_child`.

### Pre-conditions

- successful reproduction already creates new attached children deterministically
- redistribution across the pair already follows the current deterministic rule
- the current interaction payload already exposes pair identity and the successful reproduction kind

### Trigger

A different-shape interaction resolves successfully as reproduction.

### Success Outcome

- reproduction still resolves as before
- the authoritative interaction outcome explicitly exposes which concrete child IDs were newly created by the reproduction result

### Failure Or Rejection Cases

- if successful reproduction still hides which new children were created, inspectability remains incomplete
- if exposing created-child identity changes creation or redistribution behavior, scope is exceeded
- if created-child identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Child creation remains authoritative server-side behavior.
2. Successful reproduction still creates children through the existing deterministic rule.
3. The runtime explicitly exposes which new child IDs were created by that reproduction outcome.
4. Feasibility, blocked-capacity identity, cost, payment-child identity, and redistribution remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Created Child`
- `Attached Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- successful reproduction still resolves exactly as it does today
- the same deterministic child-creation and redistribution path is used as before
- snapshots expose the created child IDs only for successful reproduction outcomes
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction redesign while still making the current rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- successful reproduction already says that new children were produced
- post-resolution attached-child arrays already show the resulting state
- but the contract still does not say which concrete child IDs were newly created by that outcome

Build should therefore make one minimal contract extension that exposes created-child identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction path that can surface the created child IDs
- deterministic tests that prove the exposed IDs match the actual created children
- client rendering or HUD output that remains sufficient to read created-child identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second reproduction-event system
- reuse the same deterministic child-allocation path already used by successful reproduction
- keep the contract addition minimal and explicit
- avoid adding speculative staging, event logs, or ancestry structures

## Initial Test Plan

### Server tests

- energy-paid reproduction exposes the created child IDs
- child-paid reproduction also exposes the created child IDs
- the exposed IDs match the actual newly attached children after resolution
- blocked reproduction exposes no created-child identity

### Contract tests

- the snapshot schema accepts the new created-child identity field or fields

### Integration tests

- the client receives successful reproduction snapshots that include the created child identity

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches a successful reproduction outcome
2. the existing deterministic rule creates and redistributes new attached children
3. reproduction still resolves normally
4. the snapshot explicitly identifies which child IDs were newly created

## Done Criteria

- successful reproduction still resolves with the same current gameplay behavior
- the outcome exposes the created child identity or identities
- tests prove the exposed IDs match the actual created children
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction thresholds or costs
- staged reproduction systems
- blocked-interaction history logs
- ancestry or event-log redesign
