# Slice: Initial Food Collection And Energy Recovery

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for world snapshots that now include food entities

## Architecture Mode

Explicit client/server boundary with server authority over movement, food presence, and energy recovery.

This slice extends the existing authoritative movement loop without changing the chosen runtime topology.

## Discovery Scope

Establish the smallest deterministic feeding loop that proves Living Circles can recover energy through interaction with world resources:

- the server hosts a bounded world containing one player circle and a deterministic set of food items
- the client renders food items as part of the authoritative world snapshot
- when the player circle overlaps a food item, the server consumes that food item
- consuming food restores player energy

This slice does **not** attempt to implement:

- randomly generated food placement
- food respawn
- autonomous circles
- circle-to-circle collisions
- fight semantics
- reproduction semantics
- child accumulation
- death and continuity
- scoreboards or broader progression systems

## Why This Slice Next

The current implementation proves `move -> spend energy`, but the extracted model says the core life cycle is:

- move
- eat
- interact
- grow
- die

Without food and energy recovery, the game cannot yet express the simplest survival loop. Food collection is the smallest next step that closes the energy cycle without dragging in unresolved combat or lineage rules.

## Use-Case Contract

### Use Case

`RunMovementAndFeedingSession`

### Primary Actor

Player controlling one circle from the browser client.

### Pre-conditions

- a server process can host one bounded world
- one player-controlled circle exists in the world at session start
- the world contains a deterministic initial set of food items
- a client can open a WebSocket connection to the server

### Trigger

The player moves through the world and overlaps a food item.

### Success Outcome

- the client receives world snapshots that include food items
- when the player overlaps a food item, the item disappears from later snapshots
- the player's energy increases after consuming the food item
- movement continues to consume energy as in the first slice

### Failure Or Rejection Cases

- if the player does not overlap a food item, no energy recovery occurs
- if a food item has already been consumed, it cannot be consumed again
- if the player is at full energy, the implementation may either cap energy at the defined maximum or still clamp to the maximum after recovery, but it must be deterministic and explicit

## Main Business Rules

1. Food is an authoritative server-side world resource.
2. Food presence is part of the world snapshot contract.
3. A food item is consumed when the player circle overlaps it.
4. Consuming food restores player energy.
5. Energy must remain bounded by an explicit maximum value.
6. A consumed food item is removed from later snapshots in this slice.
7. Food placement for this slice must be deterministic so tests and scenarios are reproducible.

## Minimal Domain Concepts In Scope

- `World`
- `Player Circle`
- `Food`
- `Position`
- `Energy`
- `World Snapshot`
- `Food Consumption`

## Required Runtime Contract Changes

The server-to-client `world_snapshot` contract must now include a `foods` collection with enough information for the client to render each food item, such as:

- food item identifier
- x and y position
- radius or renderable size

The client-to-server `movement_intent` contract does not need to change in this slice.

## Required Ports Or Boundaries

- server-side world initialization that includes deterministic food placement
- server-side food consumption logic during tick advancement
- shared contract definition for food items inside snapshots
- client-side rendering boundary for food entities
- deterministic tests covering food presence and consumption

## Build Guidance

- keep food as a simple world resource, not as a generalized entity framework
- prefer a fixed initial food layout over random spawning
- keep the energy cap explicit in code and tests
- do not introduce respawn timers or procedural generation yet
- extend the existing snapshot contract rather than inventing a second data channel

## Initial Test Plan

### Server tests

- a new world contains deterministic food items
- overlapping a food item removes it from the world
- consuming food increases energy
- energy recovery cannot exceed the configured maximum

### Contract tests

- the world snapshot schema includes a stable `foods` collection
- food items expose the fields the client needs to render them

### Integration tests

- the client receives an initial snapshot containing food items
- after movement into a food item, a later snapshot shows higher energy and one fewer food item

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial world snapshot with one player circle and visible food items
2. the player moves toward a nearby food item
3. the server advances enough ticks for the player to overlap the food item
4. the client receives a later snapshot with that food item removed
5. the displayed energy is higher than it would have been from movement alone

## Done Criteria

- the server initializes a deterministic set of food items
- the shared snapshot contract includes food items
- the JavaScript client renders food items on the canvas
- overlapping food causes authoritative consumption on the server
- food consumption restores energy and respects a deterministic energy cap
- tests cover food initialization, food consumption, energy recovery, and updated snapshot shape
- the slice does not add collision, combat, reproduction, or food respawn behavior

## Out Of Scope Follow-Ups

- random or procedural food spawning
- food respawn after consumption
- autonomous foraging circles
- circle-to-circle interactions
- growth effects from food or children
- death from starvation
