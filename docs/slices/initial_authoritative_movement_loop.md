# Slice: Initial Authoritative Movement Loop

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for WebSocket messages between client and server

## Architecture Mode

Explicit client/server boundary with server authority over simulation state.

This slice intentionally exercises the extracted product shape instead of continuing with the repository's current `python_ddd_monolith` example default.

## Discovery Scope

Establish the smallest end-to-end playable loop that proves the intended architecture and preserves one core semantic pressure:

- the browser renders a circle in a bounded world
- the player can send movement intent
- the Go server advances the world authoritatively
- movement changes position and consumes energy

This slice does **not** attempt to implement:

- collision resolution
- fight semantics
- reproduction semantics
- food consumption
- child accumulation
- death and continuity
- multiplayer balance
- local interpolation beyond what is necessary to render snapshots

## Why This Slice Next

The extracted materials create immediate pressure on architecture:

- the intended product is not a single-runtime monolith
- the system is conceived as a JavaScript canvas client plus authoritative Go server over WebSocket
- movement and energy are the smallest domain behaviors that already matter to the game's identity

Without validating this runtime shape first, later slices risk refining the wrong code structure and the wrong authority model.

## Use-Case Contract

### Use Case

`RunAuthoritativeMovementSession`

### Primary Actor

Player controlling one circle from the browser client.

### Pre-conditions

- a server process can host one bounded world
- a client can open a WebSocket connection to that server
- one player-controlled circle exists in the world at session start

### Trigger

The player connects and provides movement intent.

### Success Outcome

- the client receives authoritative world snapshots from the server
- the player's circle position changes in response to movement intent
- the player's circle energy decreases when movement is applied
- the client renders the current server state on a 2D canvas

### Failure Or Rejection Cases

- if no movement intent is provided, the server keeps the circle stationary
- if the circle reaches a world boundary, movement cannot place it outside the world
- if the connection is not established, no local simulation should pretend the session is active

## Main Business Rules

1. The server is the source of truth for world state.
2. The client may capture input and render snapshots, but it does not decide final position or energy.
3. The world is bounded; positions must remain inside the map.
4. Each successful movement step consumes energy.
5. Energy does not recover in this slice.
6. One player-controlled circle is sufficient for this slice; autonomous circles are out of scope.
7. The slice should prefer deterministic tick progression so movement and energy behavior can be replayed in tests.

## Minimal Domain Concepts In Scope

- `World`
- `Player Circle`
- `Position`
- `Energy`
- `Movement Intent`
- `Tick`
- `World Snapshot`

## Required Runtime Contract

The slice must define an explicit contract for:

- client-to-server movement intent
- server-to-client world snapshot

The contract only needs fields required for this slice, such as:

- player direction or movement vector
- world bounds
- player circle position
- player circle energy
- tick or sequence identifier

The build step may choose the simplest serialization that works cleanly across JavaScript and Go, but the contract must be explicit and testable.

## Required Ports Or Boundaries

- WebSocket transport boundary between client and server
- server-side simulation loop or tick driver
- client-side rendering boundary for canvas drawing
- shared contract definition usable by both runtimes
- deterministic clock or tick control for tests where needed

## Build Guidance

- begin by updating the implementation-facing architecture documents and decisions to reflect that this slice uses `polyglot_client_server`
- keep the server-side model small and behavior-first
- keep the client thin: input capture, connection handling, and rendering only
- avoid introducing speculative abstractions for future combat, reproduction, or multiplayer behavior
- choose straightforward message naming over generic protocol frameworks

## Initial Test Plan

### Server tests

- moving a circle updates position according to movement intent
- movement reduces energy
- movement cannot place a circle outside world bounds
- idle ticks do not move the circle

### Contract tests

- the world snapshot shape is stable and parseable by the client
- the movement intent shape is stable and parseable by the server

### Integration tests

- a client can connect to the server and receive an initial world snapshot
- after sending movement intent, the next snapshot reflects changed position and reduced energy

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial world snapshot with one player circle
2. the player provides movement intent in one direction
3. the server advances at least one tick
4. the client receives a new authoritative snapshot
5. the rendered circle position changes and displayed energy is lower than before

## Done Criteria

- the repository has an explicit multi-runtime layout consistent with `polyglot_client_server`
- a Go server can host one bounded world with one player circle
- a JavaScript canvas client can connect over WebSocket and render authoritative snapshots
- movement intent flows from client to server
- server snapshots flow from server to client
- movement consumes energy
- deterministic tests cover the server rules and the runtime contract
- no combat, reproduction, or food behavior is smuggled into this slice without a new refinement pass

## Out Of Scope Follow-Ups

- food spawning and consumption
- autonomous non-player circles
- collision detection
- same-shape fight resolution
- different-shape reproduction
- child accumulation and growth
- death, replacement, and lineage
- interpolation and prediction refinements
- multi-player sessions
