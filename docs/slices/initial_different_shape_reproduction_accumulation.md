# Slice: Initial Different-Shape Reproduction Accumulation

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for reproduction outcomes and child accumulation state

## Architecture Mode

Explicit client/server boundary with the server resolving different-shape overlap into a deterministic reproduction outcome.

This slice extends the current classification-only reproduction path without committing yet to spawned child entities, lineage replacement, or growth-driven combat changes.

## Discovery Scope

Establish the smallest deterministic reproduction loop for different-shape overlap:

- different-shape overlap is no longer only classified
- a successful reproduction outcome is resolved authoritatively on the server
- the participating circles accumulate explicit child state
- the resulting child accumulation is visible in later snapshots

This slice does **not** attempt to implement:

- spawned child circles as separate world entities
- lineage replacement after death
- inherited traits or mutation
- area growth or combat power changes from children
- reproduction energy cost beyond existing movement and food rules
- same-tick multi-pair reproduction systems

## Why This Slice Next

The current implementation already makes both interaction paths visible:

- same-shape overlap resolves into a fight
- different-shape overlap remains a `reproduce_candidate`

That leaves a major asymmetry in the model. The source material says circles interact to fight or reproduce, and it explicitly says circles can accumulate children. Child accumulation is the narrowest supported next step because it:

- resolves the current reproduction placeholder into an actual outcome
- preserves fairness between player and autonomous circles
- avoids prematurely deciding whether children must appear as separate entities
- keeps continuity and lineage open for later refinement

## Use-Case Contract

### Use Case

`ResolveDifferentShapeReproduction`

### Primary Actor

Player controlling one circle while at least one autonomous circle participates in the same world.

### Pre-conditions

- a server process can host one bounded world
- the world contains at least two active circles
- the overlapping circles have different shapes
- the overlap is classified as a reproduction path
- both circles are still active participants in the world

### Trigger

The server advances a tick in which two different-shape circles overlap.

### Success Outcome

- the server resolves a deterministic reproduction outcome on that tick
- each participating circle gains explicit child accumulation state
- both circles remain active in the world
- later snapshots expose that reproduction occurred and how many children each circle currently holds

### Failure Or Rejection Cases

- if shapes match, this slice must not resolve the interaction as reproduction
- if circles are not overlapping, no reproduction occurs
- if the same pair remains continuously overlapped across multiple ticks, the build must avoid awarding repeated child accumulation every tick without an intervening separation

## Main Business Rules

1. Only different-shape overlap resolves as reproduction in this slice.
2. Same-shape overlap still follows the existing fight resolution path.
3. Reproduction resolves authoritatively on the server.
4. A successful reproduction awards exactly one child accumulation unit to each participating circle.
5. Reproduction does not remove either participant from the world.
6. A circle pair may reproduce at most once while continuously overlapping; another reproduction requires at least one separating tick.
7. This slice records child accumulation only. It does not yet turn children into growth, power, replacement, or separate child entities.

## Minimal Domain Concepts In Scope

- `Different-Shape Reproduction`
- `Child Accumulation`
- `Parent Circle`
- `Active Circle`
- `Resolved Interaction`
- `World Snapshot`

## Bounded Interpretation Of Children

The source material supports that circles can accumulate children and that children later matter for area, fight power, and continuity. It does **not** yet force one concrete representation for children.

This slice therefore chooses the smallest explicit interpretation:

- each circle carries a `children_count`
- reproduction increments `children_count`
- no separate child entity is spawned yet
- no growth or replacement effect is applied yet

This interpretation is intentionally provisional, but it is explicit, inspectable, and reversible if later refinement requires a richer child model.

## Required Runtime Contract Changes

The snapshot contract should expose enough state for the client to understand that reproduction resolved, for example:

- a `children_count` on active circles
- a resolved interaction kind such as `reproduce_resolved`
- the identifiers of the participating circles

The client-to-server `movement_intent` contract does not need to change.

## Required Ports Or Boundaries

- server-side deterministic reproduction resolution for different-shape overlap
- server-side memory of whether a pair has already reproduced during a continuous overlap
- shared contract updates for child accumulation state and resolved reproduction outcomes
- client-side rendering boundary that makes child accumulation inspectable
- deterministic tests covering single-award reproduction behavior

## Build Guidance

- keep reproduction as a one-tick authoritative resolution, not a long animation system
- prefer explicit `children_count` over introducing a generalized lineage subsystem
- keep the pair-cooldown rule minimal: one reproduction per continuous overlap is enough
- make the result inspectable in snapshots rather than hiding it in server logs
- do not add growth, replacement, or mutation behavior yet

## Initial Test Plan

### Server tests

- different-shape overlap resolves into reproduction rather than classification-only
- both participating circles gain one child accumulation unit
- same-shape overlap still resolves as a fight
- a continuously overlapping pair does not gain children repeatedly every tick
- a pair can reproduce again only after separating and overlapping again

### Contract tests

- the snapshot schema includes child accumulation on circles
- the interaction schema distinguishes `reproduce_resolved` from `fight_resolved`

### Integration tests

- the client receives snapshots showing different-shape overlap resolving into reproduction
- a later snapshot exposes increased child counts for the participating circles
- repeated ticks during the same continuous overlap do not keep increasing child counts

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot with one player circle and two autonomous circles
2. the player moves into overlap with the different-shape autonomous circle
3. the server resolves the overlap as reproduction
4. both participating circles remain active
5. the client receives a later snapshot showing a resolved reproduction outcome and increased child counts

## Done Criteria

- different-shape overlap resolves into a deterministic reproduction outcome
- both participating circles accumulate one child count
- the same pair does not accumulate children repeatedly while continuously overlapping
- the snapshot contract exposes child counts and reproduction resolution clearly enough for the client to inspect
- same-shape fight behavior still works
- tests cover child accumulation, no-repeat overlap behavior, and contract shape
- the slice does not implement spawned children, growth effects, continuity replacement, or mutation

## Out Of Scope Follow-Ups

- turning children into area growth
- turning children into combat advantage
- replacing dead circles with children
- spawned child entities with their own movement
- lineage trees or inheritance
- reproduction cost balancing
