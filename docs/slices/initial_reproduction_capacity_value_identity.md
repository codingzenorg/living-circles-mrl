# Slice: Initial Reproduction Capacity Value Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where reproduction outcomes remain inspectable through explicit current capacity values

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction feasibility, payment, and child allocation.

This slice narrows one remaining reproduction inspectability gap by exposing each participant's current reproduction capacity value alongside the already-exposed blocked-capacity and payment identity fields.

## Discovery Scope

Establish the smallest deterministic change that makes current capacity explicit:

- reproduction feasibility should still be decided by the same authoritative rule as today
- the runtime should expose the current reproduction capacity value for both source and target participants
- blocked-capacity booleans, payment identity, created-child identity, ownership identity, redistribution kind, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of child reserve contribution
- historical energy logs
- alternate capacity formulas
- new client-side effects beyond what is necessary to read the values

## Why This Slice Next

Recent slices made reproduction outcomes much more inspectable:

- blocked reproduction now exposes which side or sides lacked capacity
- child-paid reproduction now exposes which side paid through a child and which concrete child was consumed
- successful reproduction now exposes created child identity, ownership identity, and redistribution kind

But the runtime still hides the actual current capacity values that those outcomes are based on. The client and tests can see whether a side was blocked, but they still cannot see the numeric state that made the rule succeed or fail. That weakens inspectability because the server already computes a concrete current capacity for each participant when deciding whether reproduction is feasible.

The next model pressure is therefore not to change reproduction semantics. It is to make the current capacity values explicit in the same bounded way that blockage, payment, and redistribution identity are now explicit.

This slice is the narrowest next step because it:

- changes only reproduction inspectability
- preserves the current feasibility, payment, creation, and redistribution rules
- avoids inventing staged resource systems or historical logs
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeReproductionCapacityValues`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction reaches reproduction evaluation.

### Pre-conditions

- reproduction feasibility is already decided by the current authoritative capacity rule
- blocked-capacity booleans are already exposed
- the server already has enough information to derive each side's current reproduction capacity

### Trigger

A different-shape interaction reaches the authoritative reproduction decision point.

### Success Outcome

- reproduction still resolves or blocks exactly as before
- the authoritative interaction outcome explicitly exposes the source-side and target-side current capacity values used by the decision

### Failure Or Rejection Cases

- if reproduction still hides the actual capacity values, inspectability remains incomplete
- if exposing the values changes feasibility or payment behavior, scope is exceeded
- if the values are not deterministic from the existing rule, inspectability weakens

## Main Business Rules

1. Current reproduction capacity remains an authoritative server-side concept.
2. The current feasibility rule remains unchanged.
3. The runtime explicitly exposes the current capacity value for both participants at the reproduction decision point.
4. Blocked-capacity booleans, payment identity, created-child identity, ownership identity, and redistribution kind remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Reproduction Capacity`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- reproduction still resolves or blocks exactly as it does today
- the same deterministic capacity rule is used as before
- snapshots expose the evaluated current capacity values only as outcome metadata
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a resource redesign while still making the decision basis legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- blocked reproduction already exposes which side failed
- successful reproduction already exposes payment and redistribution identity
- the server already knows the evaluated current capacity values
- but the contract still does not name those values explicitly

Build should therefore make one minimal contract extension that exposes source and target reproduction capacity values without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction decision path that can surface the evaluated current capacity values directly from the existing rule
- deterministic tests that prove the exposed values match the actual blocked or successful outcomes
- client rendering or HUD output that remains sufficient to read the values

## Build Guidance

- prefer extending the existing interaction payload rather than creating a separate reproduction diagnostic event
- reuse the same deterministic capacity calculation already used by reproduction feasibility
- keep the contract addition minimal and explicit
- avoid adding speculative energy history, reserve history, or alternative resource views

## Initial Test Plan

### Server tests

- blocked reproduction exposes the source and target current capacity values that explain the failure
- successful reproduction exposes the source and target current capacity values that explain the success
- child-paid reproduction exposes capacity values consistent with the current reserve rule

### Contract tests

- the snapshot schema accepts the new capacity-value fields

### Integration tests

- the client receives reproduction snapshots that include the source and target current capacity values

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches the authoritative reproduction decision point
2. the existing deterministic rule evaluates current capacity on both sides
3. reproduction still either resolves or blocks normally
4. the snapshot explicitly exposes the source-side and target-side current capacity values used by that decision

## Done Criteria

- reproduction still resolves or blocks with the same current gameplay behavior
- the outcome exposes source and target current capacity values
- tests prove the exposed values match the actual blocked or successful outcomes
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- redesign of reproduction thresholds or costs
- historical energy or reserve logs
- alternate reserve formulas
- client-side explanation systems beyond the minimal HUD output
