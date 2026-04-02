# Slice: Initial Food Regeneration Cycle

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for regenerated food appearing in later snapshots

## Architecture Mode

Explicit client/server boundary with the server restoring world resources through a deterministic regeneration rule.

This slice extends the current feeding and survival model without introducing procedural generation, randomness, or balancing systems.

## Discovery Scope

Establish the smallest deterministic food renewal rule that allows recovery cycles in the world:

- consumed food does not stay gone forever
- food regeneration happens authoritatively on the server
- regeneration follows a deterministic timing rule
- regenerated food reappears in later snapshots at deterministic positions

This slice does **not** attempt to implement:

- random food spawning
- density balancing based on population size
- biome or area-specific resource rules
- procedural generation
- drop-on-death mechanics
- multiple food types or qualities

## Why This Slice Next

The current implementation now has:

- energy spending
- food consumption
- reproduction and growth
- fight death and continuity
- zero-energy death

But the world still trends only toward depletion. Without some form of resource renewal, the system cannot begin to exhibit the recovery side of the “collapse and recovery cycles” described in the source material and background knowledge.

Food regeneration is the smallest next step because it:

- preserves energy as the central variable
- supports longer-lived ecosystem behavior
- avoids introducing new entity types
- can remain fully deterministic and inspectable

## Use-Case Contract

### Use Case

`RegenerateWorldFood`

### Primary Actor

Player and autonomous circles participating in the same bounded world.

### Pre-conditions

- the world contains a deterministic initial food set
- food can be consumed and removed
- the server advances ticks authoritatively

### Trigger

The server advances ticks after one or more food items have been consumed.

### Success Outcome

- consumed food eventually reappears according to a deterministic rule
- regenerated food is visible in later snapshots
- the regeneration rule applies independently of whether the player is active
- the world can support longer-running feeding and survival loops

### Failure Or Rejection Cases

- if consumed food never returns, the slice fails its purpose
- if regeneration is random or time-of-day dependent, the slice loses determinism
- if regeneration duplicates an already-active food item instead of restoring a missing one, the slice becomes ambiguous

## Main Business Rules

1. Food regeneration is authoritative server-side behavior.
2. Regeneration is deterministic.
3. Only missing food items are eligible to return.
4. Regenerated food returns at known deterministic positions.
5. This slice does not change food energy value or circle energy rules.
6. This slice does not rebalance the amount of food dynamically.

## Minimal Domain Concepts In Scope

- `Food`
- `Consumed Food`
- `Regenerated Food`
- `Tick`
- `World Snapshot`

## Bounded Regeneration Interpretation

This slice chooses the smallest inspectable regeneration rule:

- each deterministic food slot may be either active or missing
- after a fixed deterministic number of ticks, a missing food item returns to its original slot

This keeps regeneration explicit and avoids hidden spawn pools or probabilistic behavior.

## Required Runtime Contract Changes

The existing snapshot contract already carries food items and their positions.

This slice may not require new contract fields, but it must ensure:

- later snapshots can include food items that were previously consumed
- tests can identify that a regenerated item corresponds to a known deterministic slot

## Required Ports Or Boundaries

- server-side tracking of missing food slots
- server-side deterministic tick-based food restoration
- client-side rendering of regenerated food through normal snapshot updates
- deterministic tests covering consumption and later return

## Build Guidance

- prefer per-slot deterministic restoration over a generalized spawning subsystem
- keep the regeneration delay explicit in code and tests
- do not add procedural or randomized placement
- avoid balancing logic tied to population counts in this slice
- keep the effect inspectable through normal snapshots and tests

## Initial Test Plan

### Server tests

- consuming a food item removes it immediately
- a consumed food item returns after the configured deterministic delay
- an already-active food item is not duplicated
- regeneration returns the item to its original slot

### Contract tests

- the existing food shape in the snapshot schema remains sufficient

### Integration tests

- the client receives a snapshot where food is consumed
- after enough ticks, the client receives a later snapshot where the same food slot is present again

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and sees the deterministic food layout
2. a food item is consumed
3. later snapshots show that slot missing for a bounded interval
4. after the deterministic regeneration delay, the food item reappears at the same position

## Done Criteria

- food now regenerates deterministically after consumption
- regeneration is visible in authoritative snapshots
- the client renders regenerated food without special-case logic
- tests cover removal, delay, and return to the original slot
- the slice does not add randomness, balancing systems, or multiple food categories

## Out Of Scope Follow-Ups

- randomized food spawning
- adaptive resource density
- death-generated resources
- regional or biome-specific food rules
- different food values or types
