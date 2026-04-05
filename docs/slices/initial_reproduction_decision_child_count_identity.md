# Slice: Initial Reproduction Decision Child Count Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where reproduction outcomes remain inspectable through explicit decision-time child counts

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction feasibility, payment, and child allocation.

This slice narrows one remaining reproduction inspectability gap by exposing how many attached children each side had at the authoritative reproduction decision point.

## Discovery Scope

Establish the smallest deterministic change that makes decision-time child counts explicit:

- reproduction feasibility should still be decided by the same authoritative rule as today
- the runtime should expose, for source and target sides, the attached-child count at the reproduction decision point
- capacity totals, energy/reserve components, threshold, cost, blocked-capacity booleans, payment identity, created-child identity, ownership identity, redistribution kind, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of child reserve contribution
- historical child-count logs across multiple ticks
- alternate child-accounting rules
- new client-side effects beyond what is necessary to read the decision-time counts

## Why This Slice Next

Recent slices made reproduction outcomes much more inspectable:

- blocked-capacity identity is explicit
- child-payment identity and concrete paid child IDs are explicit
- created-child identity, ownership identity, and redistribution kind are explicit
- total capacity values are explicit
- threshold and cost constants are explicit
- energy and reserve components of current capacity are explicit

But one important ambiguity remains. After successful reproduction, the current attached-child arrays already include both consumed-payment effects and newly created children. That means the post-outcome snapshot no longer directly tells you how many children each side had when the authoritative reproduction decision was made.

The next model pressure is therefore not to change reproduction semantics. It is to expose the decision-time child counts explicitly in the same bounded way that payment identity and capacity composition are now explicit.

This slice is the narrowest next step because it:

- changes only reproduction inspectability
- preserves the current feasibility, payment, creation, and redistribution rules
- avoids inventing historical event logs
- gives build a minimal contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeReproductionDecisionChildCounts`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction reaches reproduction evaluation.

### Pre-conditions

- reproduction feasibility is already decided by the current authoritative rule
- capacity totals, constants, and capacity components are already exposed
- the server already knows how many attached children each side has before payment and redistribution are applied

### Trigger

A different-shape interaction reaches the authoritative reproduction decision point.

### Success Outcome

- reproduction still resolves or blocks exactly as before
- the authoritative interaction outcome explicitly exposes the source-side and target-side attached-child counts at the decision point

### Failure Or Rejection Cases

- if reproduction still hides decision-time child counts, inspectability remains incomplete
- if exposing the counts changes feasibility or payment behavior, scope is exceeded
- if the counts are not deterministic from the existing rule, inspectability weakens

## Main Business Rules

1. Decision-time child count remains an authoritative server-side fact.
2. The current feasibility and payment rules remain unchanged.
3. The runtime explicitly exposes the source-side and target-side attached-child counts at the reproduction decision point.
4. Capacity totals, capacity components, threshold, cost, blocked-capacity identity, payment identity, created-child identity, ownership identity, and redistribution kind remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Decision-Time Child Count`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- reproduction still resolves or blocks exactly as it does today
- the same deterministic decision path is used as before
- snapshots expose decision-time child counts only as outcome metadata
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction-history redesign while still making the decision context fully legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- it already exposes capacity totals, threshold, cost, and capacity components
- the server already knows the decision-time child counts
- but the post-outcome attached-child arrays do not directly preserve those counts once payment and creation are applied

Build should therefore make one minimal contract extension that exposes source-side and target-side decision-time child counts without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction decision path that can surface the decision-time child counts directly from the existing rule
- deterministic tests that prove the exposed counts remain stable even when payment and creation change post-outcome attached-child arrays
- client rendering or HUD output that remains sufficient to read those counts

## Build Guidance

- prefer extending the existing interaction payload rather than creating a separate decision-history event
- reuse the same pre-payment state already available in the server decision path
- keep the contract addition minimal and explicit
- avoid adding speculative multi-tick history or event logs

## Initial Test Plan

### Server tests

- blocked reproduction exposes decision-time child counts for both sides
- child-paid reproduction exposes a target-side decision-time child count that still shows the pre-payment child availability
- ordinary energy-paid reproduction exposes zero decision-time child counts when neither side had children

### Contract tests

- the snapshot schema accepts the new decision-time child-count fields

### Integration tests

- the client receives reproduction snapshots that include the decision-time child counts

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches the authoritative reproduction decision point
2. the existing deterministic rule evaluates the pair using their current attached-child counts
3. reproduction still either resolves or blocks normally
4. the snapshot explicitly exposes the source-side and target-side child counts from that decision point, even if post-outcome attached-child arrays later differ

## Done Criteria

- reproduction still resolves or blocks with the same current gameplay behavior
- the outcome exposes the decision-time child counts
- tests prove those counts are preserved even when payment or creation changes the post-outcome child arrays
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- redesign of reproduction thresholds or costs
- redesign of reserve contribution semantics
- multi-tick decision history
- client-side explanation systems beyond the minimal HUD output
