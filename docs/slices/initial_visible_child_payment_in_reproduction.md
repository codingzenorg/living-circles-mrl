# Slice: Initial Visible Child Payment In Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible reproduction outcomes that now make child-paid reproduction inspectable

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution.

This slice extends the current deterministic reproduction rule by making child-based payment visibly grounded in the attached-child model rather than leaving it as mostly hidden bookkeeping.

## Discovery Scope

Establish the smallest deterministic rule that makes child-paid reproduction visibly consistent with the current orbiting-child model:

- when a circle pays reproduction with a child reserve instead of ordinary energy, that child loss should be visible in the authoritative result
- current reproduction feasibility, reproduction cost, child redistribution, child continuity, growth, and orbiting-child rules remain in force
- the current movement, food, contact, fight, and death rules remain unchanged

This slice does **not** attempt to implement:

- a separate staged reproduction animation system
- detached child entities during payment
- mutation or inheritance mechanics
- generalized resource accounting beyond the current reproduction cost rule
- client-side anticipation or prediction

## Why This Slice Next

The current implementation already lets a child act as a temporary reserve unit when energy is below the reproduction cost. That matches the intended gameplay direction, but the visible result is still semantically weak: attached children already matter for conflict absorption, continuity, contact, pursuit, avoidance, and food collection, while reproduction payment by child can still feel like an invisible swap inside the final count change.

The model pressure is now:

- child-paid reproduction should be visibly grounded in the same orbiting-child model used everywhere else
- the system should better distinguish “paid fully with energy” from “consumed one child to reproduce”
- reproduction should remain deterministic without opening a larger lineage or animation subsystem

This slice is the narrowest next step because it:

- changes only inspectability and visible accounting around the existing reproduction rule
- keeps the current reproduction feasibility and child redistribution semantics intact
- reuses the existing attached-child representation instead of inventing new entities
- avoids broader ancestry, mutation, or animation work

## Use-Case Contract

### Use Case

`ResolveReproductionWithVisibleChildPayment`

### Primary Actor

Any circle pair that reaches a different-shape reproduction outcome when at least one participant must pay with a child reserve.

### Pre-conditions

- reproduction resolution is already deterministic and server-authoritative
- attached children are already visible and authoritative
- child reserve payment is already part of the current reproduction rule
- current reproduction outcome snapshots already expose child counts and attached-child arrays

### Trigger

A different-shape overlap resolves as reproduction and at least one participant lacks enough direct energy to pay the reproduction cost.

### Success Outcome

- the participant that used child reserve visibly loses one attached child as part of the reproduction outcome
- the reproduction result remains deterministic and inspectable in ordinary snapshots
- later snapshots make child-paid reproduction legible as distinct from energy-only reproduction

### Failure Or Rejection Cases

- if child-paid reproduction still looks identical to energy-only reproduction, the orbiting-child model remains uneven
- if build introduces a second independent child accounting system, coherence weakens
- if the rule changes reproduction feasibility or redistribution semantics silently, slice scope is exceeded

## Main Business Rules

1. Reproduction remains authoritative server-side behavior.
2. Current reproduction feasibility and cost rules remain unchanged.
3. If a participant pays through child reserve, that payment must be visible in attached-child state after resolution.
4. Current deterministic child redistribution across participants remains unchanged unless a minimal explicit adjustment is required to preserve visibility.
5. The rule must remain deterministic for the same world state and tick.
6. Current fight, food, movement, and death rules remain unchanged.
7. This slice should reuse the current attached-child representation rather than inventing detached child-payment entities.
8. Player and autonomous circles must follow the same child-payment visibility rule.

## Minimal Domain Concepts In Scope

- `Reproduction Payment`
- `Child Reserve`
- `Visible Child Loss`
- `Reproduction Outcome`
- `Attached Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- the current reproduction feasibility threshold and cost stay exactly the same
- what changes is that child-based payment becomes visibly grounded in the authoritative attached-child state after resolution
- build may use one explicit resolved outcome distinction if ordinary child-count changes are too ambiguous
- no new ancestry, mutation, or detached-child semantics are introduced

This avoids the larger step of full reproduction staging while still aligning the payment rule with the embodied child model.

## Required Runtime Contract Changes

The current contract may be sufficient if the effect is visible through:

- changed attached-child arrays
- changed child counts
- existing interaction outcome data

Build should extend the contract only if child-paid reproduction remains too ambiguous to distinguish from energy-only reproduction.

## Required Ports Or Boundaries

- server-side reproduction resolution that can surface visible child payment without breaking deterministic redistribution
- deterministic tests that show child-paid reproduction differs observably from energy-only reproduction
- client rendering that remains sufficient to observe the changed child state

## Build Guidance

- prefer evolving the current reproduction resolution path rather than adding a separate staging subsystem
- preserve the current reproduction threshold, cost, and deterministic redistribution rule unless a minimal documented adjustment is required
- keep movement, energy drain, and downstream lineage rules unchanged
- avoid new animation channels, event streams, or detached child placeholders unless the current snapshot shape proves insufficient

## Initial Test Plan

### Server tests

- a reproduction paid with ordinary energy preserves the expected attached-child result under the current rule
- a reproduction paid through child reserve visibly consumes one attached child from the paying participant
- two otherwise similar reproduction cases differ observably when one uses child payment and the other does not
- deterministic redistribution remains stable after visible child payment

### Contract tests

- the current snapshot schema remains sufficient unless build adds one explicit reproduction outcome distinction

### Integration tests

- the client receives snapshots showing visible child loss when reproduction is paid through child reserve
- the client also receives ordinary reproduction snapshots when both participants pay fully with energy

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a different-shape pair reaches reproduction contact
2. one participant has enough reproduction capacity only because it has an attached child available for reserve payment
3. reproduction resolves
4. the paying participant visibly loses one attached child as part of the result
5. the ordinary redistribution result still appears under the current deterministic rule

## Done Criteria

- child-paid reproduction is visibly distinguishable from energy-only reproduction
- the rule is deterministic and documented
- current reproduction feasibility and cost semantics remain unchanged
- existing fight, food, and death rules remain unchanged
- tests cover child-paid and energy-only reproduction paths

## Out Of Scope Follow-Ups

- reproduction animations
- detached child-payment entities
- mutation and inheritance systems
- multi-step incubation
- removal of current transitional count-based shortcuts
