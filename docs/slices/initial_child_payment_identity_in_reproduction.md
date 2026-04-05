# Slice: Initial Child Payment Identity In Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where child-paid reproduction remains inspectable through explicit payer identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution.

This slice narrows one remaining reproduction inspectability gap by exposing which participant paid reproduction cost through an attached child when `reproduce_paid_child` occurs.

## Discovery Scope

Establish the smallest deterministic change that makes child-payment identity explicit:

- when `reproduce_paid_child` occurs, the child-payment path should remain exactly as it is today
- the runtime should expose which participant or participants consumed a child as reserve payment
- reproduction outcome, child redistribution, continuity, fight, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction rules
- new payment thresholds or costs
- alternate reserve strategies
- redesign of child redistribution
- new client-side effects beyond what is necessary to read payer identity

## Why This Slice Next

Recent slices already made child-dependent outcomes much more inspectable:

- `death_promoted_child` now exposes the promoted child identity
- `fight_absorbed_child` now exposes the absorbed child identity
- `reproduce_paid_child` already distinguishes child-funded reproduction from ordinary energy-funded reproduction

But reproduction still hides who actually paid with a child. That weakens inspectability because the server knows whether the player, the opponent, or both participants consumed a child as reserve payment, while the runtime contract only says that at least one child payment happened.

The next model pressure is therefore not to change reproduction semantics. It is to make child-payment identity explicit in the same way promotion and absorption identity were made explicit.

This slice is the narrowest next step because it:

- changes only reproduction inspectability
- preserves the current threshold, cost, and redistribution rules
- avoids inventing new fertility or resource systems
- gives build a deterministic contract extension instead of a larger semantic redesign

## Use-Case Contract

### Use Case

`ExposeChildPaymentIdentityDuringReproduction`

### Primary Actor

Any player or autonomous circle that successfully reproduces by consuming one attached child as reserve payment.

### Pre-conditions

- reproduction already supports `reproduce_paid_child`
- one participant may already consume one attached child as reserve payment
- the current interaction payload already exposes the pair-level reproduction outcome

### Trigger

A different-shape interaction resolves into `reproduce_paid_child`.

### Success Outcome

- reproduction still resolves as before
- child redistribution still resolves as before
- the authoritative reproduction outcome explicitly exposes which participant or participants paid through a child

### Failure Or Rejection Cases

- if `reproduce_paid_child` still hides who paid with a child, inspectability remains incomplete
- if exposing payer identity changes reproduction behavior, scope is exceeded
- if payer identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Reproduction resolution remains authoritative server-side behavior.
2. `reproduce_paid_child` still means at least one participant used one attached child as reserve payment.
3. The reproduction outcome explicitly exposes which participant or participants paid through a child.
4. Threshold, cost, redistribution, and later child state remain unchanged.
5. Fight, continuity, feeding, contact, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Reproduction`
- `Child Payment`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- reproduction still resolves exactly as it does today
- the same deterministic payer path is used as before
- snapshots expose which participant or participants consumed a child for payment through the existing interaction channel or an equally small contract extension
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction redesign while still making the rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `reproduce_paid_child` says that child-funded payment occurred
- the pair identity is already visible
- but the contract still does not say which participant or participants consumed a child

Build should therefore make one minimal contract extension that exposes child-payment identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction path that can surface which participant paid through a child
- deterministic tests that prove the exposed identity matches the consumed child-payment path
- client rendering or HUD output that remains sufficient to read the child-payment identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second reproduction-event system
- reuse the same deterministic payment path already used by `reproduce_paid_child`
- keep the contract addition minimal and explicit
- avoid adding speculative fertility, mutation, or accounting systems

## Initial Test Plan

### Server tests

- player-paid reproduction exposes that the player paid with a child
- autonomous-paid reproduction exposes that the autonomous participant paid with a child
- when both sides can pay with energy, no child-payment identity is emitted beyond the existing ordinary reproduction path

### Contract tests

- the snapshot schema accepts the new child-payment identity field

### Integration tests

- the client receives a `reproduce_paid_child` snapshot that includes which participant paid through a child

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape reproduction path reaches `reproduce_paid_child`
2. one participant deterministically consumes one attached child as reserve payment
3. reproduction still resolves normally
4. the snapshot explicitly identifies which participant paid through a child

## Done Criteria

- `reproduce_paid_child` still resolves with the same current gameplay behavior
- the reproduction outcome exposes child-payment identity
- tests prove the exposed identity matches the actual payer path
- broader reproduction, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction rules
- alternate payment strategies
- fertility or mutation systems
- ancestry or historical event logs
