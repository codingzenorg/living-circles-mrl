# Slice: Initial Shape Identity And Collision Classification

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for circle shape and current collision classification

## Architecture Mode

Explicit client/server boundary with the server classifying circle-to-circle overlap and publishing the resulting interaction state.

This slice adds semantic structure to circle interaction without yet choosing combat or reproduction outcomes.

## Discovery Scope

Establish the smallest interaction slice that makes Living Circles’ defining shape rule explicit and observable:

- every circle has a stable shape identity
- the server detects overlap between the player circle and the autonomous circle
- when circles overlap, the server classifies the overlap as either `fight_candidate` or `reproduce_candidate`
- the client renders shape identity and the current classified interaction state

This slice does **not** attempt to implement:

- fight resolution
- reproduction outcomes
- child accumulation
- circle death
- circle replacement or continuity
- multiple autonomous circles
- random shape assignment

## Why This Slice Next

The extracted model says:

- circles have a shape or type
- same shape trends toward fight
- different shape trends toward reproduction

The current implementation has shared food pressure but still cannot express the domain’s core interaction distinction. Adding explicit shape identity and collision classification is the smallest next step that prepares the model for later combat or reproduction slices without silently inventing their outcomes now.

## Use-Case Contract

### Use Case

`RunClassifiedInteractionSession`

### Primary Actor

Player controlling one circle while one autonomous circle participates in the same world.

### Pre-conditions

- a server process can host one bounded world
- the world contains one player-controlled circle and one autonomous circle
- each circle has a deterministic shape identity
- a client can open a WebSocket connection to the server

### Trigger

The server advances ticks until the two circles overlap.

### Success Outcome

- snapshots expose each circle’s shape identity
- when circles overlap, snapshots expose a current collision classification
- same-shape overlap is classified as `fight_candidate`
- different-shape overlap is classified as `reproduce_candidate`
- no fight or reproduction side effects occur in this slice

### Failure Or Rejection Cases

- if circles do not overlap, no interaction classification is active
- if shape identity is missing for either circle, the server should not invent a hidden default at runtime; the world must define shapes explicitly
- if circles overlap repeatedly across ticks, the classification may remain active while overlap persists, but it must do so deterministically

## Main Business Rules

1. Every circle has an explicit shape identity.
2. Shape identity is authoritative server-side state and part of the snapshot contract.
3. Overlap between circles is classified server-side.
4. Same-shape overlap becomes `fight_candidate`.
5. Different-shape overlap becomes `reproduce_candidate`.
6. This slice publishes classification only; it does not resolve fight or reproduction.
7. Existing movement, food, and energy rules remain intact.

## Minimal Domain Concepts In Scope

- `Shape`
- `Circle`
- `Circle Overlap`
- `Interaction Classification`
- `fight_candidate`
- `reproduce_candidate`
- `World Snapshot`

## Required Runtime Contract Changes

The snapshot contract should now include:

- shape identity for the player circle
- shape identity for autonomous circles
- an explicit interaction field or collection describing the current classified overlap state

The contract only needs to express the current active classification for this slice, for example:

- the involved circle identifiers
- the classification kind
- whether the interaction is currently active

The client-to-server `movement_intent` contract does not need to change.

## Required Ports Or Boundaries

- server-side deterministic shape assignment for the current circles
- server-side overlap detection between circles
- shared contract definition for shape identity and current interaction classification
- client-side rendering boundary for shape distinction and active classification display
- deterministic tests covering overlap classification logic

## Build Guidance

- keep shape identity simple and explicit, such as a small fixed set of strings
- prefer a single current interaction classification over a generalized event stream
- keep overlap detection geometric and deterministic
- do not sneak in damage, child spawning, or growth effects
- make the client rendering distinct enough that a reviewer can see shape difference without reading logs

## Initial Test Plan

### Server tests

- a new world assigns deterministic shapes to the player and autonomous circle
- overlapping same-shape circles produce `fight_candidate`
- overlapping different-shape circles produce `reproduce_candidate`
- non-overlapping circles produce no active classification

### Contract tests

- the snapshot schema includes shape identity for all circles
- the snapshot schema includes explicit interaction classification fields

### Integration tests

- the client receives snapshots that expose both circle shapes
- after overlap, a later snapshot includes an active interaction classification
- the classification kind matches the configured shapes of the overlapping circles

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot containing shape identity for both circles
2. the player moves into the autonomous circle, or the autonomous circle reaches the player depending on the deterministic setup
3. the server detects the overlap
4. the client receives a later snapshot with an active interaction classification
5. the client can distinguish whether the overlap is a `fight_candidate` or a `reproduce_candidate`

## Done Criteria

- both circles have explicit shape identity
- snapshots expose shape identity and current interaction classification
- the server classifies same-shape overlap as `fight_candidate`
- the server classifies different-shape overlap as `reproduce_candidate`
- the client renders the two circles distinctly enough to reflect shape identity
- tests cover shape assignment, overlap detection, and classification logic
- the slice does not implement combat damage, reproduction, children, growth, or death

## Out Of Scope Follow-Ups

- fight resolution rules
- reproduction resolution rules
- child spawning or accumulation
- growth and power changes
- death and continuity semantics
- multiple simultaneous circle interactions
