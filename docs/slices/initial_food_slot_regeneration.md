# Slice: Initial Food Slot Regeneration

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for regenerated food appearing in later snapshots

## Architecture Mode

Explicit client/server boundary with the server restoring missing food through deterministic slot regeneration.

This slice extends the current feeding and survival loop without introducing randomness, procedural generation, or ecosystem balancing heuristics.

## Discovery Scope

Establish the smallest deterministic resource-renewal rule that lets the world recover from depletion:

- each food item belongs to a known deterministic slot
- when a food item is consumed, that slot becomes inactive
- after a fixed deterministic delay, the missing slot becomes active again
- regenerated food reappears in the same position and with the same identity

This slice does **not** attempt to implement:

- random food spawning
- adaptive density based on population or deaths
- different food categories or values
- regional spawning logic
- generated resources from fights or corpses
- continuous growth of food pools

## Why This Slice Next

The current implementation already supports:

- energy spending
- food consumption
- reproduction, growth, and continuity
- death from conflict or zero energy

But the world still only moves toward exhaustion because food never returns. That blocks the “collapse and recovery cycles” implied by the source material and background knowledge.

Slot-based food regeneration is the narrowest next step because it:

- preserves energy as the central survival currency
- supports longer-lived simulation behavior
- stays fully deterministic
- avoids opening broader balancing or procedural-generation questions

## Use-Case Contract

### Use Case

`RegenerateConsumedFoodSlot`

### Primary Actor

Player and autonomous circles participating in the same bounded world.

### Pre-conditions

- the world contains deterministic food slots
- food can be consumed and removed from active snapshots
- the server advances ticks authoritatively

### Trigger

The server advances enough ticks after a food slot has been consumed.

### Success Outcome

- the consumed food slot becomes active again after the deterministic delay
- regenerated food appears in later snapshots at the original position
- active food slots are never duplicated
- the feeding loop can continue without restarting the world

### Failure Or Rejection Cases

- if a consumed food slot never returns, the slice fails its purpose
- if regeneration occurs randomly rather than deterministically, the slice loses inspectability
- if regeneration duplicates an already-active slot, the rule is invalid

## Main Business Rules

1. Food regeneration is authoritative server-side behavior.
2. Food regeneration is deterministic.
3. Only inactive food slots are eligible to regenerate.
4. Regenerated food returns to its original slot identity and position.
5. Active food slots are never duplicated.
6. This slice does not change food energy gain or circle movement rules.

## Minimal Domain Concepts In Scope

- `Food Slot`
- `Consumed Food`
- `Regenerated Food`
- `Tick`
- `World Snapshot`

## Bounded Regeneration Interpretation

This slice chooses the smallest inspectable regeneration rule:

- each deterministic food item is treated as a reusable slot
- when consumed, the slot is marked inactive
- after a fixed deterministic number of ticks, the slot becomes active again

This avoids hidden spawn pools, probabilistic behavior, and balancing subsystems.

## Required Runtime Contract Changes

The existing snapshot contract already carries food identity and position.

This slice should preserve that contract and ensure:

- regenerated food reuses the same slot ID
- later snapshots can be compared directly against earlier snapshots to confirm regeneration

## Required Ports Or Boundaries

- server-side tracking of inactive food slots
- server-side deterministic regeneration timing
- client-side rendering of regenerated food through ordinary snapshots
- deterministic tests covering disappearance and return

## Build Guidance

- prefer one explicit regeneration delay in code and tests
- keep regeneration per-slot, not global or probabilistic
- do not add new message types for food changes
- avoid balancing rules tied to circle counts or deaths
- keep the effect inspectable through ordinary snapshots and tests

## Initial Test Plan

### Server tests

- consuming a food item removes it from active slots immediately
- a consumed food item returns after the deterministic delay
- the regenerated item has the same ID and position as the original slot
- active slots are not duplicated

### Contract tests

- the existing snapshot schema remains sufficient for regeneration

### Integration tests

- the client receives a snapshot where one food item is consumed
- after enough ticks, the client receives a later snapshot where that same food slot is present again

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and sees the initial deterministic food layout
2. a food item is consumed
3. later snapshots show the slot missing for a bounded interval
4. after the deterministic delay, the same food slot reappears in the same place

## Done Criteria

- consumed food now regenerates deterministically
- regeneration is visible in authoritative snapshots
- the client renders regenerated food without special logic
- tests cover removal, delay, and return to the original slot
- the slice does not add randomness, balancing systems, or multiple food categories

## Out Of Scope Follow-Ups

- randomized food spawning
- adaptive resource density
- death-generated resources
- region-specific food rules
- multiple food values or types
