# Slice: Initial Child Payment Child Identity In Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where child-paid reproduction remains inspectable through explicit consumed-child identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution and child payment.

This slice narrows one remaining reproduction inspectability gap by exposing which concrete attached child was consumed when a participant pays reproduction cost through child reserve.

## Discovery Scope

Establish the smallest deterministic change that makes child-payment identity fully explicit:

- when `reproduce_paid_child` occurs, the same authoritative payment path should remain exactly as it is today
- the runtime should expose which concrete source-side and/or target-side attached child was consumed as payment
- reproduction feasibility, blocked-capacity identity, contact identity, child redistribution, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- alternate reserve strategies
- redesign of redistribution
- historical reproduction logs
- new client-side effects beyond what is necessary to read consumed-child identity

## Why This Slice Next

Recent slices made the reproduction path much more inspectable:

- `reproduce_paid_child` now exposes which side paid through a child
- `reproduce_blocked_energy` now exposes which side lacked enough current capacity
- child-triggered contact now exposes which child or children started the interaction

But child-paid reproduction still stops one step short. The runtime can now say which side paid with a child, yet it still hides which concrete child was consumed on that side. That weakens inspectability because the authoritative server already knows the exact child identity used for payment, while the client and tests can only infer it indirectly from the changed attached-child set.

The next model pressure is therefore not to change reproduction semantics. It is to make child-payment identity explicit in the same bounded way promotion, fight absorption, and contact identity are already explicit.

This slice is the narrowest next step because it:

- changes only child-payment inspectability
- preserves the current feasibility threshold, cost, reserve rule, and redistribution rule
- avoids inventing staged reproduction or ancestry systems
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeConsumedChildIdentityDuringReproductionPayment`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction resolves as `reproduce_paid_child`.

### Pre-conditions

- reproduction already supports child-reserve payment
- `reproduce_paid_child` already identifies which side or sides paid through a child
- the current payment path already consumes one deterministic attached child from each paying side

### Trigger

A different-shape interaction resolves as `reproduce_paid_child`.

### Success Outcome

- reproduction still resolves as before
- the authoritative interaction outcome explicitly exposes which concrete source-side and/or target-side child was consumed as payment

### Failure Or Rejection Cases

- if `reproduce_paid_child` still hides which child was consumed, inspectability remains incomplete
- if exposing consumed-child identity changes payment behavior, scope is exceeded
- if consumed-child identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Reproduction payment remains authoritative server-side behavior.
2. `reproduce_paid_child` still means at least one side paid with one attached child.
3. The runtime explicitly exposes which concrete source-side and/or target-side child was consumed as payment.
4. Reproduction feasibility, blocked-capacity identity, cost, reserve rule, and child redistribution remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Reproduction Payment`
- `Consumed Payment Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- `reproduce_paid_child` still resolves exactly as it does today
- the same deterministic child-consumption path is used as before
- snapshots expose the consumed source-side and/or target-side child IDs only when payment actually used a child
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction redesign while still making the current rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `reproduce_paid_child` already says which side paid through a child
- pair identity is already visible
- but the contract still does not say which concrete child was consumed

Build should therefore make one minimal contract extension that exposes consumed payment-child identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction payment path that can surface the consumed child IDs
- deterministic tests that prove the exposed IDs match the actual consumed children
- client rendering or HUD output that remains sufficient to read the consumed-child identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second reproduction-event system
- reuse the same deterministic child-selection path already used by child payment
- keep the contract addition minimal and explicit
- avoid adding speculative staging, event logs, or ancestry structures

## Initial Test Plan

### Server tests

- source-side child-paid reproduction exposes the consumed source child ID only
- target-side child-paid reproduction exposes the consumed target child ID only
- when both sides pay through a child, both consumed child IDs are exposed
- energy-only reproduction exposes no payment-child identity

### Contract tests

- the snapshot schema accepts the new consumed payment-child identity fields

### Integration tests

- the client receives a `reproduce_paid_child` snapshot that includes the consumed source-side and/or target-side child IDs

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches `reproduce_paid_child`
2. one side, the other side, or both sides deterministically consume one attached child as payment
3. reproduction still resolves normally
4. the snapshot explicitly identifies the consumed payment child or children

## Done Criteria

- child-paid reproduction still resolves with the same current gameplay behavior
- the outcome exposes the consumed payment-child identity or identities
- tests prove the exposed IDs match the actual consumed children
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction thresholds or costs
- staged reproduction systems
- blocked-interaction history logs
- ancestry or event-log redesign
