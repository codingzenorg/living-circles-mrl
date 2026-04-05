# Slice: Initial Reproduction Distribution Kind Identity

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where successful reproduction remains inspectable through explicit redistribution kind

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution and child allocation.

This slice narrows one remaining successful-reproduction inspectability gap by exposing which deterministic redistribution kind was selected when new children were allocated across the pair.

## Discovery Scope

Establish the smallest deterministic change that makes redistribution kind explicit:

- when reproduction resolves successfully, the same authoritative child-creation and redistribution path should remain exactly as it is today
- the runtime should expose whether the successful result allocated children as `source_only`, `split`, or `target_only`
- created-child identity, ownership identity, feasibility, blocked-capacity identity, payment-child identity, contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new reproduction thresholds or costs
- redesign of redistribution
- alternate child-creation counts
- historical reproduction logs
- new client-side effects beyond what is necessary to read redistribution kind

## Why This Slice Next

Recent slices made successful reproduction much more inspectable:

- successful reproduction now exposes which child IDs were created
- successful reproduction now exposes which side received each created child
- child-paid reproduction now exposes which concrete child was consumed as payment

But the deterministic redistribution rule itself still remains implicit. The runtime shows the created-child IDs and their ownership, yet it still does not explicitly say which redistribution case was selected. That weakens inspectability because the authoritative server already knows whether the result was source-only, split, or target-only, while the client and tests still need to infer that from the per-side ownership lists.

The next model pressure is therefore not to change reproduction semantics. It is to make the chosen redistribution kind explicit in the same bounded way ownership and creation identity are now explicit.

This slice is the narrowest next step because it:

- changes only successful reproduction inspectability
- preserves the current feasibility rule, payment rule, creation rule, and redistribution rule
- avoids inventing staged reproduction or ancestry systems
- gives build a deterministic contract extension instead of a broader reproduction redesign

## Use-Case Contract

### Use Case

`ExposeReproductionDistributionKind`

### Primary Actor

Any player or autonomous circle pair whose different-shape interaction resolves as `reproduce_resolved` or `reproduce_paid_child`.

### Pre-conditions

- successful reproduction already creates children deterministically
- created-child identity and ownership identity are already exposed
- redistribution across the pair already follows the current deterministic rule

### Trigger

A different-shape interaction resolves successfully as reproduction.

### Success Outcome

- reproduction still resolves as before
- the authoritative interaction outcome explicitly exposes which redistribution kind was selected

### Failure Or Rejection Cases

- if successful reproduction still hides the chosen redistribution kind, inspectability remains incomplete
- if exposing redistribution kind changes allocation behavior, scope is exceeded
- if redistribution kind can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Child allocation remains authoritative server-side behavior.
2. Successful reproduction still allocates children through the existing deterministic redistribution rule.
3. The runtime explicitly exposes which redistribution kind was selected.
4. Created-child identity, ownership identity, feasibility, blocked-capacity identity, cost, and payment-child identity remain unchanged.
5. Contact identity, fight, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Reproduction Distribution Kind`
- `Created Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- successful reproduction still resolves exactly as it does today
- the same deterministic child-creation and redistribution path is used as before
- snapshots expose the selected redistribution kind only for successful reproduction outcomes
- all broader reproduction semantics remain unchanged

This avoids turning the slice into a reproduction redesign while still making the allocation rule legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- successful reproduction already exposes created child IDs and per-side ownership
- the server already knows which redistribution case was selected
- but the contract still does not name that case explicitly

Build should therefore make one minimal contract extension that exposes redistribution kind without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side reproduction path that can surface the selected redistribution kind directly from the existing allocation rule
- deterministic tests that prove the exposed kind matches the actual created-child ownership on each side
- client rendering or HUD output that remains sufficient to read redistribution kind

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second reproduction-event system
- reuse the same deterministic allocation path already used by successful reproduction
- keep the contract addition minimal and explicit
- avoid adding speculative staging, event logs, or ancestry structures

## Initial Test Plan

### Server tests

- source-only successful reproduction exposes `source_only`
- split successful reproduction exposes `split`
- target-only successful reproduction exposes `target_only`
- blocked reproduction exposes no redistribution kind

### Contract tests

- the snapshot schema accepts the new redistribution kind field

### Integration tests

- the client receives successful reproduction snapshots that include the redistribution kind

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape interaction reaches a successful reproduction outcome
2. the existing deterministic rule creates and redistributes new attached children
3. reproduction still resolves normally
4. the snapshot explicitly identifies whether the redistribution was `source_only`, `split`, or `target_only`

## Done Criteria

- successful reproduction still resolves with the same current gameplay behavior
- the outcome exposes redistribution kind
- tests prove the exposed kind matches the actual created-child ownership on each side
- broader reproduction, contact, fight, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new reproduction thresholds or costs
- staged reproduction systems
- blocked-interaction history logs
- ancestry or event-log redesign
