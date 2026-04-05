# Slice: Initial Reproduction Rule Constant Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where reproduction outcomes remain inspectable through explicit rule constants

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction feasibility, payment, and child allocation.

This slice narrows one remaining reproduction inspectability gap by exposing the reproduction threshold and payment cost values that the current capacity fields are measured against.

## Discovery Scope

Establish the smallest deterministic change that makes the governing rule constants explicit:

- reproduction feasibility should still be decided by the same authoritative rule as today
- the runtime should expose the current reproduction threshold and reproduction cost values alongside the existing capacity and blocked/payment identity fields
- capacity values, blocked-capacity booleans, payment identity, created-child identity, ownership identity, redistribution kind, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of child reserve contribution
- alternate balancing rules
- historical rule-change logs
- new client-side effects beyond what is necessary to read the constants

## Why This Slice Next

Recent slices made reproduction outcomes much more inspectable:

- blocked reproduction now exposes which side or sides lacked capacity
- child-paid reproduction exposes which side paid through a child and which concrete child was consumed
- successful reproduction exposes created-child identity, ownership identity, and redistribution kind
- reproduction outcomes now expose the evaluated current capacity values on each side

But the runtime still leaves the governing rule constants implicit. The client and tests can see the evaluated capacities, yet they still need repository knowledge to know what threshold those values were compared against and what cost was actually paid on success.

The next model pressure is therefore not to change reproduction semantics. It is to expose the current threshold and cost values explicitly in the same bounded way that capacity values and identities are now explicit.

This slice is the narrowest next step because it:

- changes only reproduction inspectability
- preserves the current feasibility, payment, creation, and redistribution rules
- avoids inventing balancing systems or dynamic tuning
- gives build a minimal contract extension instead of a broader rule redesign

## Use-Case Contract

### Use Case

`ExposeReproductionRuleConstants`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction reaches reproduction evaluation.

### Pre-conditions

- reproduction feasibility is already decided by the current authoritative capacity rule
- current capacity values are already exposed
- the threshold and cost constants already exist in the server model

### Trigger

A different-shape interaction reaches the authoritative reproduction decision point.

### Success Outcome

- reproduction still resolves or blocks exactly as before
- the authoritative interaction outcome explicitly exposes the current threshold and cost constants used by that decision

### Failure Or Rejection Cases

- if reproduction still hides the threshold or cost constants, inspectability remains incomplete
- if exposing the constants changes feasibility or payment behavior, scope is exceeded
- if the constants vary non-deterministically, inspectability weakens

## Main Business Rules

1. Reproduction threshold remains an authoritative server-side constant.
2. Reproduction payment cost remains an authoritative server-side constant.
3. The current feasibility and payment rules remain unchanged.
4. The runtime explicitly exposes the threshold and cost constants used by the current decision.
5. Capacity values, blocked-capacity identity, payment identity, created-child identity, ownership identity, and redistribution kind remain unchanged.
6. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Reproduction Capacity`
- `Reproduction Threshold`
- `Reproduction Cost`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- reproduction still resolves or blocks exactly as it does today
- the same deterministic threshold and cost constants are used as before
- snapshots expose those constants only as outcome metadata
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a balancing redesign while still making the rule basis fully legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- it already exposes blocked-capacity booleans and current capacity values
- the server already knows the threshold and cost constants used by the current rule
- but the contract still does not name those constants explicitly

Build should therefore make one minimal contract extension that exposes reproduction threshold and cost values without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction decision path that can surface the threshold and cost constants directly from the existing rule
- deterministic tests that prove the exposed constants match the existing reproduction rule
- client rendering or HUD output that remains sufficient to read them

## Build Guidance

- prefer extending the existing interaction payload rather than creating a separate rule-description event
- reuse the same constant values already used by the server rule
- keep the contract addition minimal and explicit
- avoid adding speculative balancing metadata or version histories

## Initial Test Plan

### Server tests

- blocked reproduction exposes the current threshold and cost constants
- successful reproduction exposes the same threshold and cost constants
- child-paid reproduction still reports the same rule constants

### Contract tests

- the snapshot schema accepts the new threshold and cost fields

### Integration tests

- the client receives reproduction snapshots that include the threshold and cost constants

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches the authoritative reproduction decision point
2. the existing deterministic rule evaluates current capacity against the current threshold and current payment cost
3. reproduction still either resolves or blocks normally
4. the snapshot explicitly exposes the threshold and cost constants used by that decision

## Done Criteria

- reproduction still resolves or blocks with the same current gameplay behavior
- the outcome exposes the governing threshold and cost constants
- tests prove the exposed constants match the existing server rule
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- redesign of reproduction thresholds or costs
- dynamic tuning or balancing systems
- historical rule-change logs
- client-side explanation systems beyond the minimal HUD output
