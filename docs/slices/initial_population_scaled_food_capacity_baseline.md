# Slice: Initial Population Scaled Food Capacity Baseline

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where the larger world baseline can expose food abundance that scales coherently with the starting population

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for initial world composition.

This slice introduces one bounded ecosystem-baseline rule after the larger world and regeneration-pressure slices: the default food-capacity baseline should be derived from the initialized population profile rather than remain a hand-picked fixed count.

## Discovery Scope

Establish the smallest deterministic rule that ties initial food abundance to initial population scale:

- keep food slots deterministic and fixed in place
- keep initial world creation deterministic
- derive default starting food capacity from the initial active population profile
- keep regeneration, movement, energy, fight, reproduction, continuity, child, and steering rules unchanged

This slice does **not** attempt to implement:

- random food seeding
- procedural map generation
- per-region ecology
- adaptive runtime food spawning
- player-configurable difficulty
- long-term balancing across many scenario presets

## Why This Slice Next

The current world is now larger and regeneration pressure is more ecosystem-like, but the initial food baseline is still effectively a curated hand choice. That means the world starts from a larger stage, yet the initial abundance still depends on fixed authored numbers instead of an explicit relation to population scale.

The next pressure is therefore not another local mechanic. It is to make the initial resource baseline itself more model-shaped:

- larger starting populations should imply a larger starting food baseline
- smaller custom worlds should still be able to stay sparse and deterministic
- initialization should become easier to reason about as a rule rather than as a remembered set of tuned constants

This slice is the narrowest next step because it:

- changes only initial resource abundance logic
- keeps the world deterministic
- complements the new regeneration-pressure rule without redesigning runtime ecology
- directly improves the ecosystem baseline rather than adding another metadata detail

## Use-Case Contract

### Use Case

`StartPopulationScaledFoodWorld`

### Primary Actor

The authoritative server constructing a new initial world or reset state.

### Pre-conditions

- the server already determines the initial active population deterministically
- the world already supports deterministic fixed food slots
- the larger default world baseline already exists

### Trigger

A new world or reset state is created.

### Success Outcome

- the starting food baseline is derived from the initial population profile using one documented deterministic rule
- the larger default world no longer depends on a purely hand-picked starting food count
- custom narrower worlds can still keep a smaller deterministic resource baseline

### Failure Or Rejection Cases

- if initial food abundance remains a fixed authored count, the larger ecosystem baseline stays partly arbitrary
- if initial food abundance becomes random or hard to inspect, deterministic refinement weakens
- if the slice starts changing runtime regeneration or movement behavior, scope is exceeded

## Main Business Rules

1. Initial food abundance remains authoritative server-side behavior.
2. Initial food slots remain deterministic and fixed in place.
3. The default initial food baseline should scale from the initialized population profile through one small deterministic rule.
4. The rule should remain simple enough to explain briefly and test directly.
5. Regeneration, movement, energy, fight, reproduction, continuity, child ownership, and steering remain unchanged.
6. Reset should recreate the same population-scaled initial abundance for the same world profile.

## Minimal Domain Concepts In Scope

- `Initial Population`
- `Initial Food Capacity`
- `Deterministic World Baseline`
- `Reset`

## Bounded Baseline Interpretation

This slice chooses the smallest useful interpretation of population-scaled abundance:

- derive the default starting food count from the starting population profile
- keep slot identity and positions deterministic
- avoid runtime balancing or dynamic spawn rules

This avoids turning the slice into a general ecology system while still making the initial world less arbitrarily tuned.

## Required Runtime Contract Changes

The current contract is likely sufficient because:

- the client already receives the complete food array
- world size and population are already visible in ordinary snapshots

Build should avoid new contract surface unless a minimal inspectability field becomes clearly necessary.

## Required Ports Or Boundaries

- server-side world initialization logic
- deterministic tests that prove the chosen food-capacity rule
- implementation notes that explain the new initialization rule

## Build Guidance

- prefer evolving the current default food-slot construction logic rather than introducing a separate subsystem
- choose one simple population-based derivation rule
- keep custom narrow configs deterministic and practical for focused tests
- avoid changing regeneration timing in the same slice

## Initial Test Plan

### Server tests

- the expanded default world derives a larger starting food baseline from its larger initial population
- a narrower custom world still starts with a smaller deterministic food baseline
- reset recreates the same population-scaled initial food set

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives the population-scaled initial food baseline in ordinary snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the server creates a world from its deterministic initial population profile
2. initial food abundance is derived from that profile rather than hard-coded as a remembered constant
3. the larger default world begins with a richer but still bounded resource baseline
4. reset recreates the same starting abundance deterministically

## Done Criteria

- initial food abundance is derived from starting population scale
- the default world keeps a richer baseline than narrow custom worlds
- initialization remains deterministic
- runtime regeneration and other core rules remain unchanged
- tests cover the new initialization rule

## Out Of Scope Follow-Ups

- dynamic runtime food spawning
- biome or region logic
- random map generation
- new autonomy rules
- new fight or reproduction semantics
