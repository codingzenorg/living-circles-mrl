# Slice: Initial Autonomous Circle Participation

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for snapshots that include more than one circle

## Architecture Mode

Explicit client/server boundary with server authority over all circle behavior, including non-player movement and food consumption.

This slice extends the existing feeding loop without introducing circle-to-circle interaction yet.

## Discovery Scope

Establish the smallest multi-circle world that proves the player is participating in a shared rule system rather than occupying a privileged solitary loop:

- the server hosts one player-controlled circle and one autonomous non-player circle
- both circles exist in the same bounded world with the same food resources
- the autonomous circle follows a deterministic movement policy
- the autonomous circle can consume food and recover energy
- the client renders both the player circle and the autonomous circle

This slice does **not** attempt to implement:

- collisions between circles
- fight semantics
- reproduction semantics
- child accumulation
- death and continuity
- random autonomous behavior
- multiple autonomous circles
- food respawn

## Why This Slice Next

The extracted model says:

- the player participates as one circle
- all entities follow the same rules
- the game is a system-driven ecosystem rather than a purely player-centric arena

The current implementation still treats the player as the only circle in existence. Adding one deterministic autonomous circle is the smallest step that begins to validate fairness, shared resource pressure, and non-player participation without forcing unresolved combat or reproduction semantics.

## Use-Case Contract

### Use Case

`RunSharedFeedingSession`

### Primary Actor

Player controlling one circle while one autonomous circle participates in the same world.

### Pre-conditions

- a server process can host one bounded world
- the world contains one player-controlled circle, one autonomous circle, and a deterministic initial set of food items
- a client can open a WebSocket connection to the server

### Trigger

The server advances ticks while the player sends movement intent and the autonomous circle follows its deterministic movement policy.

### Success Outcome

- world snapshots include the player circle and the autonomous circle
- the autonomous circle moves without direct player input
- the autonomous circle can consume food and recover energy
- both circles remain governed by the same movement and energy principles

### Failure Or Rejection Cases

- if the autonomous circle has no movement policy, it should remain stationary rather than inventing ad hoc behavior
- if food is consumed by one circle, it is no longer available to the other
- if either circle reaches zero energy, movement behavior must remain deterministic and explicit for that circle

## Main Business Rules

1. The server remains the source of truth for both player and non-player circles.
2. The world contains at least one autonomous non-player circle in addition to the player circle.
3. The autonomous circle follows a deterministic movement policy defined by this slice.
4. The autonomous circle consumes energy when it moves.
5. The autonomous circle can consume food and recover energy under the same energy cap as the player circle.
6. Food is shared across circles; once consumed, it is removed from the world for all participants.
7. This slice still does not define any circle-to-circle interaction outcome.

## Minimal Domain Concepts In Scope

- `World`
- `Player Circle`
- `Autonomous Circle`
- `Food`
- `Energy`
- `Deterministic Movement Policy`
- `World Snapshot`

## Required Runtime Contract Changes

The server-to-client `world_snapshot` contract should evolve from a single `player` record to a shape that can express:

- the player-controlled circle
- one or more non-player circles
- enough information for the client to distinguish which circle is the player

The build step may use either:

- a `player` field plus a separate `circles` or `autonomous_circles` collection
- or a unified `circles` collection with an explicit role marker

Whichever shape is chosen must be explicit, stable, and tested.

The client-to-server `movement_intent` contract does not need to change in this slice.

## Required Ports Or Boundaries

- server-side deterministic movement policy for one autonomous circle
- server-side world advancement that applies the same energy and food rules to more than one circle
- shared contract definition for multiple circles in a snapshot
- client-side rendering boundary for non-player circles
- deterministic tests covering non-player participation and shared food pressure

## Build Guidance

- keep the autonomous movement policy simple and deterministic, such as a fixed direction sequence or bounded patrol pattern
- do not introduce pathfinding, steering systems, or AI layers
- keep food competition implicit through shared resources rather than through direct combat
- avoid restructuring the entire codebase into a generalized ECS or similar abstraction
- make the player/non-player distinction explicit in the contract so the client remains simple

## Initial Test Plan

### Server tests

- a new world contains one player circle and one autonomous circle
- the autonomous circle moves on server ticks without player input
- the autonomous circle spends energy when it moves
- the autonomous circle can consume food and recover energy
- food consumed by the autonomous circle disappears from later snapshots

### Contract tests

- the snapshot schema expresses both the player-controlled circle and non-player circles
- non-player circles expose the fields the client needs to render them

### Integration tests

- the client receives snapshots containing more than one circle
- after several ticks, the autonomous circle position changes
- when the autonomous circle reaches food, a later snapshot reflects both food removal and energy recovery

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot with the player circle, one autonomous circle, and visible food
2. the autonomous circle begins moving according to its deterministic policy
3. the player moves independently under the existing movement rules
4. the autonomous circle reaches a food item and consumes it
5. the client receives a later snapshot showing the autonomous circle in a new position, one fewer food item, and recovered energy for that circle

## Done Criteria

- the server hosts one autonomous non-player circle alongside the player circle
- snapshots expose the autonomous circle in a stable contract shape
- the client renders the autonomous circle distinctly from the player
- the autonomous circle moves deterministically and consumes energy
- the autonomous circle can consume shared food and recover energy
- tests cover multi-circle snapshots, deterministic autonomous movement, and shared food consumption
- the slice still does not define fight, reproduction, child accumulation, or death behavior

## Out Of Scope Follow-Ups

- multiple autonomous circles
- circle-to-circle collisions
- same-shape fight resolution
- different-shape reproduction
- children, growth, and continuity
- random movement policies or emergent steering
- death from starvation or defeat
