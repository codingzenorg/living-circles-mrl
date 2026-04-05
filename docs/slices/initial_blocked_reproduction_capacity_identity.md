# Slice: Initial Blocked Reproduction Capacity Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where blocked reproduction remains inspectable through explicit capacity-failure identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction feasibility and outcome reporting.

This slice narrows one remaining reproduction inspectability gap by exposing which participant or participants lacked enough current capacity when a different-shape interaction resolves as `reproduce_blocked_energy`.

## Discovery Scope

Establish the smallest deterministic change that makes blocked reproduction identity explicit:

- when reproduction is blocked, the same current feasibility rule should remain exactly as it is today
- the runtime should expose which side or sides failed the current reproduction-capacity check
- reproduction cost, child payment, child redistribution, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- alternate reserve strategies
- new blocked-reproduction avoidance rules
- history logs for failed interactions
- new client-side effects beyond what is necessary to read failure identity

## Why This Slice Next

Recent slices made the reproduction path much more inspectable:

- `reproduce_paid_child` now exposes which side paid through a child
- child-triggered contact now exposes which child or children started the interaction
- continuity and fight absorption already expose which child was consumed

But blocked reproduction still stops one step short. The runtime can say that reproduction was blocked, yet it still hides which participant lacked enough current capacity. That weakens inspectability because the authoritative server already knows whether the source side, target side, or both failed the feasibility rule, while the client and tests can only infer it indirectly from energy and attached-child state.

The next model pressure is therefore not to change reproduction semantics. It is to make blocked reproduction identity explicit in the same bounded way other child- and reproduction-dependent outcomes are now explicit.

This slice is the narrowest next step because it:

- changes only blocked reproduction inspectability
- preserves the current feasibility threshold, cost, and reserve rules
- avoids inventing new failure-recovery or strategy systems
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeBlockedReproductionCapacityIdentity`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction reaches `reproduce_blocked_energy`.

### Pre-conditions

- reproduction feasibility is already determined authoritatively on the server
- the current feasibility rule already depends on current energy plus the child-reserve path
- the interaction payload already exposes pair identity and the blocked reproduction kind

### Trigger

A different-shape interaction resolves as `reproduce_blocked_energy`.

### Success Outcome

- blocked reproduction still resolves as before
- the authoritative interaction outcome explicitly exposes which participant or participants lacked enough current reproduction capacity

### Failure Or Rejection Cases

- if `reproduce_blocked_energy` still hides which side failed the feasibility rule, inspectability remains incomplete
- if exposing blocked-capacity identity changes whether reproduction succeeds or fails, scope is exceeded
- if failure identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Reproduction feasibility remains authoritative server-side behavior.
2. `reproduce_blocked_energy` still means the current reproduction rule prevented resolution.
3. The blocked reproduction outcome explicitly exposes which source-side and/or target-side participant lacked enough current capacity.
4. Threshold, cost, reserve payment, and child redistribution rules remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Blocked Reproduction`
- `Capacity Failure`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- blocked reproduction still occurs exactly as it does today
- the same deterministic feasibility rule is used as before
- snapshots expose whether the source side, target side, or both failed the capacity check
- all broader reproduction semantics remain unchanged

This avoids turning blocked reproduction into a strategy or recovery system while still making the current rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `reproduce_blocked_energy` says that reproduction failed
- pair identity is already visible
- but the contract still does not say which side or sides failed the capacity check

Build should therefore make one minimal contract extension that exposes blocked-capacity identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction-feasibility path that can surface which side failed
- deterministic tests that prove the exposed identity matches the actual feasibility check
- client rendering or HUD output that remains sufficient to read blocked-capacity identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second failure-event system
- reuse the same deterministic feasibility rule already used by `reproduce_blocked_energy`
- keep the contract addition minimal and explicit
- avoid adding speculative cooldowns, debt systems, or retry-state tracking

## Initial Test Plan

### Server tests

- source-only blocked reproduction exposes that the source side failed the capacity check
- target-only blocked reproduction exposes that the target side failed the capacity check
- both-sides blocked reproduction exposes that both sides failed the capacity check
- successful reproduction paths expose no blocked-capacity identity

### Contract tests

- the snapshot schema accepts the new blocked-capacity identity fields

### Integration tests

- the client receives a `reproduce_blocked_energy` snapshot that includes which side or sides lacked enough current capacity

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches the blocked reproduction path
2. the current feasibility rule determines that one side, the other side, or both sides lack enough capacity
3. reproduction remains blocked as before
4. the snapshot explicitly identifies which side or sides failed the capacity check

## Done Criteria

- blocked reproduction still resolves with the same current gameplay behavior
- the blocked outcome exposes capacity-failure identity
- tests prove the exposed identity matches the actual feasibility check
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction thresholds or costs
- retry, debt, or cooldown systems
- blocked-interaction history logs
- redesign of avoidance behavior
