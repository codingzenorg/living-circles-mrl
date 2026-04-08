# Slice: Initial Viewport Interest-Culling Transport

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a player-following viewport
- Go server with authoritative simulation state
- shared runtime contract still based on authoritative snapshots, but thinned for viewport-era transport pressure

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice introduces the first bounded transport optimization after measurement. It does not attempt generic compression, delta encoding, or protocol versioning. It uses the already-established viewport presentation model to reduce how much remote world state is sent on each tick.

## Discovery Scope

Reduce per-client snapshot payload by limiting full-detail world state to the area that is relevant to the current player viewport.

This slice should:

- keep the player and world bounds fully authoritative
- keep nearby circles and foods fully represented
- stop sending distant world entities at full detail every tick
- preserve enough whole-world orientation for the minimap to remain useful

This slice does **not** attempt to implement:

- binary encoding
- compression
- delta snapshots
- protocol version negotiation
- client-side prediction
- multiplayer interest management between independent players

## Why This Slice Next

The current measurement baseline made the transport pressure explicit:

- default expanded world snapshot: about `6487` bytes
- larger deterministic scenario snapshot: about `12539` bytes
- both are broadcast as full snapshots at `10` snapshots/sec

That is acceptable for local lab work, but the current payload cost now scales directly with world population and map breadth. The client has also already shifted into a viewport model, which means it no longer needs the entire world at equal detail just to render the main play surface.

The next pressure is therefore not a generic optimization. It is a coherence move:

- the client sees a local viewport
- the transport still sends the whole world in full detail

This slice is the narrowest next step because it:

- follows an already accepted presentation model
- should produce a meaningful payload reduction
- avoids a larger protocol redesign

## Use-Case Contract

### Use Case

`SendViewportRelevantAuthoritativeSnapshot`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- the server knows the player's authoritative position
- the browser already renders a bounded player-following viewport
- the current protocol measurement baseline exists

### Trigger

The server produces the next authoritative snapshot for a connected browser client.

### Success Outcome

- the player still receives authoritative local world state needed for play
- nearby circles and foods remain fully represented
- distant world entities are no longer sent at the same detail cadence
- the client retains enough whole-world orientation to keep the minimap useful
- per-client payload size decreases measurably from the current full-snapshot baseline

### Failure Or Rejection Cases

- if the slice breaks minimap orientation entirely, scope is incomplete
- if the slice requires versioning or complex fallbacks, scope has drifted
- if it changes gameplay semantics instead of transport representation, scope is exceeded

## Main Business Rules

1. This is a transport-thinning slice, not a gameplay slice.
2. The player, world bounds, and current interaction outcome remain authoritative.
3. Entities inside or near the current viewport should remain fully represented.
4. Distant world entities should no longer consume the same full-detail payload on every tick.
5. Whole-world orientation should remain available in a lighter form sufficient for the minimap.

## Minimal Domain Concepts In Scope

- `Viewport-Relevant Snapshot`
- `Nearby Entity Detail`
- `Distant Entity Summary`
- `Per-Client Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful transport-thinning interpretation:

- keep the existing authoritative snapshot shape as the conceptual base
- derive a client-facing transport snapshot that sends full detail only for viewport-relevant entities
- represent the larger world outside that region with lighter orientation data suitable for the minimap

This avoids jumping directly to deltas or compression while still making transport cost meaningfully responsive to player locality.

## Required Runtime Contract Changes

Likely yes.

The current full `world_snapshot` contract is too blunt for this slice. Build should expect a bounded contract evolution such as:

- a full-detail local entity set
- a lighter distant-world summary for minimap/orientation

But the slice should avoid an expansive redesign beyond what the viewport model strictly needs.

## Required Ports Or Boundaries

- server-side snapshot assembly at the transport boundary
- shared runtime contract definitions for local-vs-distant world data
- browser client consumption of the thinned authoritative snapshot
- tests that prove both play-surface correctness and payload reduction

## Build Guidance

- prefer transport-boundary shaping over changing simulation storage
- keep the authoritative world model intact in simulation
- derive viewport relevance from the current player-centered viewport plus a small margin
- preserve minimap usefulness with a lighter distant-world representation
- compare payload size against the current measurement baseline

## Initial Test Plan

### Server tests

- nearby entities remain present in full detail
- distant entities are omitted from the full-detail transport set
- the transport helper reports a smaller payload than the current full snapshot baseline for representative expanded worlds

### Contract tests

- the new transport shape remains explicit and parseable

### Integration tests

- the browser still renders local play correctly in viewport mode
- the minimap still receives enough whole-world information to orient the player

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect with the current player-follow camera active
2. receive a transport snapshot with full local detail and lighter distant-world orientation
3. verify that local play remains correct in the viewport
4. verify that the minimap still shows where the viewport sits in the larger world
5. verify that payload size is lower than the measured full-snapshot baseline

## Done Criteria

- the client no longer receives the whole large world at equal detail every tick
- the viewport play surface still has the authoritative local state it needs
- minimap orientation still works
- measured payload size is lower than the current baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- compression
- delta snapshots
- protocol versioning
- adaptive cadence reduction
- prediction or interpolation
