# Slice: Initial Population Scale World Expansion

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary world snapshots can now carry a larger bounded world and a modestly larger active population

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world initialization, autonomous participation, food placement, and snapshot broadcasting.

This slice increases the playable world space and the default active population so the simulation can move from a pairwise rules demo toward a small but more ecosystem-like world.

## Discovery Scope

Establish the smallest deterministic scale-up that materially changes ordinary play:

- enlarge the bounded world from the current demo-sized map
- increase the number of autonomous circles beyond the current two-circle default
- increase food availability enough to support the larger initial population
- keep initialization deterministic and inspectable
- keep the current movement, energy, fight, reproduction, continuity, child, and steering rules unchanged

This slice does **not** attempt to implement:

- randomized world generation
- unbounded population growth
- new autonomy subsystems
- crowding penalties
- mutation or inherited variation
- client-side camera, zoom, or minimap systems
- balancing for production-scale performance

## Why This Slice Next

The current EGD result identified the main gap clearly: the build proves that the rules exist, but the world is still too small and curated to validate ecosystem behavior.

Right now the implementation deliberately keeps:

- one player
- exactly two autonomous circles
- deterministic initial placements
- a very small number of deterministic food slots

That is appropriate for earlier slices, but it keeps the simulation legible mostly as a mechanics demonstrator rather than as a shared world with visible population pressure.

The next model pressure is therefore not to add another local mechanic. It is to create enough simultaneous participants and enough space for:

- more than one encounter path to exist at once
- recovery and depletion to happen in parallel
- dominance, survival, and avoidance to become more observable at the world level

This slice is the narrowest next step because it:

- changes the world baseline rather than the rule set
- keeps determinism intact
- makes existing mechanics interact under denser conditions
- directly addresses the strongest EGD finding without jumping to large new systems

## Use-Case Contract

### Use Case

`StartExpandedPopulationWorld`

### Primary Actor

The local player entering the default world, plus the authoritative server responsible for constructing the initial simulation state.

### Pre-conditions

- the current world model already supports one player, autonomous circles, food slots, and bounded space
- autonomous circles already participate under the current shared rules
- the client already renders ordinary snapshots with multiple circles and food items

### Trigger

A new session starts or the current demo world is reset.

### Success Outcome

- the server initializes a larger bounded world
- the default population includes one player plus several autonomous circles rather than only two
- food capacity is increased enough that the world does not immediately collapse into empty-space drift
- snapshots expose the larger world and the expanded initial state through the existing contract shape

### Failure Or Rejection Cases

- if the world becomes larger but population stays nearly unchanged, encounter density may become worse rather than better
- if population increases without corresponding food support, the slice may become a starvation demo rather than an ecosystem step
- if initialization becomes random or hard to inspect, deterministic refinement weakens
- if the slice adds new rule systems instead of only scaling the current baseline, scope is exceeded

## Main Business Rules

1. The server remains authoritative for world size, initial population, and food initialization.
2. The default world should include one player and a modest deterministic set of autonomous circles, not just two.
3. Food initialization should scale with the larger initial population.
4. Initialization remains deterministic for the same build and reset event.
5. Existing movement, energy, fight, reproduction, continuity, child, and steering rules remain unchanged.
6. The slice should aim for a small ecosystem baseline, not a production-scale population simulation.
7. The client should remain able to render the larger world without requiring a new camera system in this slice.

## Minimal Domain Concepts In Scope

- `World Size`
- `Initial Population`
- `Food Capacity`
- `Deterministic Initialization`
- `Shared World Snapshot`

## Bounded Scale Interpretation

This slice chooses the smallest useful interpretation of “increase the world space and the population”:

- enlarge the world to a clearly larger bounded map
- expand the default autonomous population to a small multi-actor set
- scale food slots to match that larger baseline
- preserve deterministic spawn layout rather than introducing procedural generation

This avoids turning the slice into a balancing system, procedural world generator, or camera redesign.

## Required Runtime Contract Changes

The current contract is likely sufficient because:

- `world.width` and `world.height` already exist
- the snapshot already supports multiple autonomous circles and foods

Build should only extend the contract if a minimal summary field becomes necessary for inspectability, such as current population counts. Otherwise the slice should prefer using the existing snapshot structure.

## Required Ports Or Boundaries

- server-side world/session initialization logic
- deterministic initial autonomous-circle and food placement for the larger world
- tests that prove expanded initialization stays deterministic and bounded
- client rendering that remains readable on the larger map without introducing authority leaks

## Build Guidance

- prefer extending the current world/session config rather than inventing a second initialization path
- choose a modest population increase that still keeps tests and live inspection practical
- increase food support together with population rather than treating them independently
- preserve one player-centered entry point while avoiding a player-privileged rule change
- avoid adding new ecology mechanics in the same slice

## Initial Test Plan

### Server tests

- a fresh world initializes with the larger configured dimensions
- a fresh world initializes with the expanded deterministic autonomous population
- a fresh world initializes with a scaled deterministic food set
- reset recreates the same expanded initial world deterministically
- all initialized circles and foods remain inside the enlarged world bounds

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal population summary field

### Integration tests

- the client receives a snapshot for the larger world with the expanded population and food set
- reset reproduces the same larger baseline over WebSocket and HTTP

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the server initializes a larger bounded world
2. multiple autonomous circles appear under the same current rules
3. food is distributed at a scale appropriate for that larger baseline
4. the player can observe more than one simultaneous movement and interaction opportunity without changing the underlying mechanics

## Done Criteria

- the default world is materially larger than the current demo baseline
- the default active population includes several autonomous circles rather than only two
- food capacity is increased alongside the larger world and population
- initialization and reset remain deterministic
- existing mechanics remain unchanged outside scale and initialization
- tests cover expanded initialization and determinism

## Out Of Scope Follow-Ups

- procedural world generation
- new ecological pressures such as scarcity tuning or crowding
- minimap, camera, or zoom systems
- performance work for large populations
- new fight or reproduction semantics
