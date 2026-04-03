# Slice: Initial Lineage Identity Through Replacement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for lineage identity and replacement continuity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for continuity outcomes and lineage identity.

This slice does not add new autonomous intelligence or child entities. It makes the already-implemented replacement rule legible and inspectable.

## Discovery Scope

Establish the smallest explicit lineage model that gives replacement continuity observable meaning:

- every active circle carries a stable `lineage_id`
- every active circle carries a bounded `generation`
- replacement after defeat or zero-energy collapse preserves `lineage_id`
- replacement increments `generation`

This slice does **not** attempt to implement:

- separate spawned child entities
- parent-child trees
- inherited traits or mutation
- reproduction-created lineage branching
- ancestry history or replay timelines

## Why This Slice Next

The current implementation already supports:

- child accumulation
- growth from children
- replacement continuity after defeat
- replacement continuity after zero-energy collapse
- deterministic food regeneration for longer-running cycles

But continuity is still semantically opaque. A replaced circle remains active, yet neither the client nor the snapshot contract can show whether that active circle is the same line continuing or just a silent reset.

Making lineage identity explicit is the narrowest next step because it:

- directly matches the source language of `lineage` and `continuity`
- clarifies the meaning of child-based replacement without forcing child entities
- improves inspectability in snapshots and the demo
- stays deterministic and cheap

## Use-Case Contract

### Use Case

`ContinueLineageThroughReplacement`

### Primary Actor

Player and autonomous circles participating under the same authoritative continuity rules.

### Pre-conditions

- a circle has an active `lineage_id`
- replacement continuity is already possible when a defeated or collapsed circle has children
- the server publishes authoritative snapshots

### Trigger

A circle dies through fight resolution or zero-energy collapse and has at least one child available for replacement continuity.

### Success Outcome

- the active replacement remains part of the same lineage
- the replacement increments generation deterministically
- snapshots expose lineage continuity explicitly
- the browser demo can show that continuity happened rather than only that a circle still exists

### Failure Or Rejection Cases

- if replacement loses lineage identity, continuity is semantically empty
- if generation changes unpredictably, the rule loses inspectability
- if the player gets lineage semantics different from autonomous circles, fairness is weakened

## Main Business Rules

1. Every active circle belongs to exactly one lineage.
2. Initial circles start with a deterministic lineage identity.
3. Replacement continuity preserves lineage identity.
4. Replacement continuity increments generation by exactly one.
5. Fight and zero-energy replacement use the same lineage rule.
6. This slice does not define branching ancestry or inherited traits.

## Minimal Domain Concepts In Scope

- `Lineage`
- `Lineage Identity`
- `Generation`
- `Replacement Continuity`
- `World Snapshot`

## Bounded Lineage Interpretation

This slice chooses the smallest lineage model that makes existing continuity observable:

- `lineage_id` identifies the continuing line
- `generation` counts replacement continuity steps inside that line
- initial circles start at generation `0`
- a replacement circle in the same line becomes generation `1`, then `2`, and so on

This avoids ancestry graphs, child entities, and mutation systems while still making continuity explicit.

## Required Runtime Contract Changes

The snapshot contract should be extended so active circles expose:

- `lineage_id`
- `generation`

The contract should preserve existing circle identity fields and make clear that:

- stable circle `id` identifies the current active participant
- stable `lineage_id` identifies the continuing line across replacement

## Required Ports Or Boundaries

- server-side lineage assignment for initial circles
- server-side lineage preservation and generation increment during replacement
- client-side rendering or labeling of lineage continuity through ordinary snapshots
- deterministic tests covering both player and autonomous replacement continuity

## Build Guidance

- prefer one explicit lineage model shared by player and autonomous circles
- keep lineage state inside the authoritative server model
- expose lineage through the existing snapshot flow rather than new message types
- keep generation deterministic and monotonic
- do not introduce ancestry graphs or historical stores

## Initial Test Plan

### Server tests

- initial circles start with deterministic lineage identities and generation `0`
- autonomous replacement after fight preserves `lineage_id` and increments generation
- player replacement after fight preserves `lineage_id` and increments generation
- zero-energy replacement preserves `lineage_id` and increments generation

### Contract tests

- the snapshot schema explicitly includes `lineage_id` and `generation` for circles

### Integration tests

- the client receives snapshots where replacement continuity is visible through stable `lineage_id`
- generation changes only when replacement occurs

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and observes initial circles with generation `0`
2. a same-shape fight or zero-energy collapse triggers replacement continuity
3. a later snapshot shows the same `lineage_id` but a higher `generation`
4. the demo makes clear that the line continued rather than silently resetting

## Done Criteria

- active circles now expose lineage identity and generation
- replacement continuity preserves lineage identity and increments generation
- fight defeat and zero-energy collapse share the same lineage rule
- snapshots and the demo make continuity observable
- tests cover initial assignment and replacement updates

## Out Of Scope Follow-Ups

- separate child entities
- branching lineage trees
- inherited traits
- mutation
- ancestry history views
