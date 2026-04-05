# Slice: Initial Reproduction Capacity Component Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where reproduction outcomes remain inspectable through explicit capacity components

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction feasibility, payment, and child allocation.

This slice narrows one remaining reproduction inspectability gap by exposing how each side's reported reproduction capacity is composed from direct energy versus child-reserve contribution.

## Discovery Scope

Establish the smallest deterministic change that makes capacity composition explicit:

- reproduction feasibility should still be decided by the same authoritative rule as today
- the runtime should expose, for source and target sides, how much of the reported current capacity comes from direct energy and how much comes from child reserve contribution
- threshold, cost, total capacity values, blocked-capacity booleans, payment identity, created-child identity, ownership identity, redistribution kind, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of child reserve contribution
- alternate reserve formulas
- historical contribution logs
- new client-side effects beyond what is necessary to read the component values

## Why This Slice Next

Recent slices made reproduction outcomes much more inspectable:

- blocked-capacity identity is explicit
- child-payment identity and concrete paid child IDs are explicit
- created-child identity, ownership identity, and redistribution kind are explicit
- current capacity values are explicit
- threshold and cost constants are explicit

But one important ambiguity remains. The runtime now tells you the total current capacity and the rule constants, yet it still does not say how much of that capacity came from direct energy versus child reserve. That matters because child reserve is the distinctive part of the current reproduction rule, and right now the client and tests still need to infer that split from other fields.

The next model pressure is therefore not to change reproduction semantics. It is to expose the composition of current capacity in the same bounded way that totals, thresholds, and payment identities are now explicit.

This slice is the narrowest next step because it:

- changes only reproduction inspectability
- preserves the current feasibility, payment, creation, and redistribution rules
- avoids inventing alternative reserve mechanics
- gives build a minimal contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeReproductionCapacityComponents`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction reaches reproduction evaluation.

### Pre-conditions

- reproduction feasibility is already decided by the current authoritative capacity rule
- current capacity values, threshold, and cost are already exposed
- the server already knows whether a side's current capacity includes child reserve contribution

### Trigger

A different-shape interaction reaches the authoritative reproduction decision point.

### Success Outcome

- reproduction still resolves or blocks exactly as before
- the authoritative interaction outcome explicitly exposes the energy contribution and child-reserve contribution that compose each side's current capacity

### Failure Or Rejection Cases

- if reproduction still hides how total capacity is composed, inspectability remains incomplete
- if exposing capacity composition changes feasibility or payment behavior, scope is exceeded
- if the component values are not deterministic from the existing rule, inspectability weakens

## Main Business Rules

1. Current reproduction capacity remains an authoritative server-side concept.
2. Direct energy contribution remains the current energy available at the decision point.
3. Child-reserve contribution remains the current reserve contribution implied by the existing rule.
4. The current feasibility and payment rules remain unchanged.
5. The runtime explicitly exposes the energy and child-reserve components for both source and target participants.
6. Threshold, cost, total capacity values, blocked-capacity identity, payment identity, created-child identity, ownership identity, and redistribution kind remain unchanged.
7. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Reproduction Capacity`
- `Energy Contribution`
- `Child Reserve Contribution`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- reproduction still resolves or blocks exactly as it does today
- the same deterministic capacity and reserve rule is used as before
- snapshots expose energy and reserve contribution values only as outcome metadata
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reserve redesign while still making the decision basis fully legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- it already exposes total capacity values, threshold, and cost
- the server already knows whether child reserve contributes to those totals
- but the contract still does not expose the component split directly

Build should therefore make one minimal contract extension that exposes source-side and target-side energy and reserve contributions without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction decision path that can surface the current energy and reserve contribution values directly from the existing rule
- deterministic tests that prove the exposed components add up to the already-exposed total capacity values
- client rendering or HUD output that remains sufficient to read the component values

## Build Guidance

- prefer extending the existing interaction payload rather than creating a separate reproduction-diagnostic event
- reuse the same deterministic capacity calculation already used by the server rule
- keep the contract addition minimal and explicit
- avoid adding speculative balancing metadata or alternate reserve interpretations

## Initial Test Plan

### Server tests

- blocked reproduction exposes energy and reserve components consistent with the reported total capacity
- child-paid reproduction exposes a target-side reserve contribution consistent with the current rule
- energy-only successful reproduction exposes zero reserve contribution on both sides

### Contract tests

- the snapshot schema accepts the new capacity-component fields

### Integration tests

- the client receives reproduction snapshots that include the energy and reserve contribution fields

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches the authoritative reproduction decision point
2. the existing deterministic rule evaluates current capacity from direct energy and optional child reserve contribution
3. reproduction still either resolves or blocks normally
4. the snapshot explicitly exposes the energy and reserve contribution values used to form each side's current capacity

## Done Criteria

- reproduction still resolves or blocks with the same current gameplay behavior
- the outcome exposes the capacity component values
- tests prove the components add up to the already-exposed total capacity values
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- redesign of reproduction thresholds or costs
- redesign of reserve contribution semantics
- dynamic balancing systems
- historical contribution logs
