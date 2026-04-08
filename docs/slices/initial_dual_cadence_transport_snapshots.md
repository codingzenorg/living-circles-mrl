# Slice: Initial Dual-Cadence Transport Snapshots

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a player-following viewport
- Go server with authoritative simulation state
- shared runtime contract already split between local play detail and lighter minimap/orientation summaries

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps the new viewport-interest transport shape, but stops treating every part of that transport as equally time-sensitive. The main play surface still needs quick local updates, while minimap-style orientation data changes more slowly and does not need to be resent at the same cadence.

## Discovery Scope

Reduce per-client transport cost further by introducing two bounded update cadences inside the existing viewport-era transport model:

- keep local full-detail play data at the current interactive cadence
- update distant orientation summaries less frequently
- preserve the current authoritative world model and gameplay semantics

This slice does **not** attempt to implement:

- delta compression
- binary encoding
- protocol versioning
- client-side prediction
- interpolation frameworks
- adaptive per-entity priority queues

## Why This Slice Next

The first thinning step already reduced the default expanded-world payload from about:

- `6487` bytes per snapshot
- to about `3907` bytes per snapshot

That is a meaningful reduction, but the transport is still sent at `10` snapshots/sec as one equally frequent package. The local viewport data justifies that cadence. The minimap-style orientation data usually does not.

The next pressure is therefore cadence coherence:

- local play detail is high-frequency
- distant orientation is lower-frequency
- the transport should reflect that distinction directly

This slice is the narrowest next step because it:

- builds on the existing viewport/local-vs-distant split
- reduces traffic without introducing compression or deltas
- keeps the client mental model simple

## Use-Case Contract

### Use Case

`SendHighCadenceLocalStateWithLowerCadenceOrientationState`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- the server already builds viewport-thinned transport snapshots
- the client already uses local full-detail arrays for play and lighter minimap summaries for orientation
- transport payload measurement already exists

### Trigger

The server produces the next authoritative update for a connected browser client.

### Success Outcome

- local play data remains responsive
- distant orientation data is refreshed often enough to keep the minimap useful
- average per-client transport cost drops below the current post-culling baseline
- gameplay semantics remain unchanged

### Failure Or Rejection Cases

- if local play becomes visibly stale, the slice failed its primary constraint
- if the client must now manage a complex replay or reconciliation model, scope drifted
- if the slice quietly changes simulation timing rather than transport cadence, scope is exceeded

## Main Business Rules

1. This is a transport-cadence slice, not a simulation-timing slice.
2. The local viewport detail path remains the high-priority stream.
3. Distant orientation data may update less frequently than local play data.
4. The client must continue to render the minimap and support orientation without inventing hidden world state.
5. The optimization should remain understandable and measurable.

## Minimal Domain Concepts In Scope

- `High-Cadence Local Snapshot`
- `Lower-Cadence Orientation Snapshot`
- `Average Per-Client Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful cadence interpretation:

- keep one authoritative websocket stream
- keep sending local play detail every tick
- send the heavier distant orientation summary only every N ticks, with the client reusing the last valid summary in between

This avoids jumping directly to deltas or multiple channels while still exploiting the different time sensitivity already present in the current transport shape.

## Required Runtime Contract Changes

Possibly small.

The existing transport shape may only need one bounded signal such as:

- whether minimap/orientation summary fields are present on this message

The slice should avoid a large protocol redesign.

## Required Ports Or Boundaries

- server-side transport snapshot assembly and broadcast path
- browser client handling of partial-orientation refreshes
- shared runtime contract definitions if presence/absence becomes explicit
- transport measurement tests comparing average cost before and after cadence splitting

## Build Guidance

- preserve current local viewport responsiveness first
- keep minimap behavior simple, with the client reusing the most recent valid orientation summary between refreshes
- measure average cost over a deterministic tick window rather than only single-message size
- avoid hidden client prediction for distant summaries

## Initial Test Plan

### Server tests

- local viewport detail is still present on every transport snapshot
- orientation summary appears only on the intended lower-frequency ticks
- average bytes/sec over a deterministic window is lower than the current post-culling baseline

### Contract tests

- the snapshot shape remains explicit about when orientation summary is or is not present

### Integration tests

- viewport play remains responsive
- the minimap continues to render using the most recent valid orientation summary

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect and receive high-cadence local viewport updates
2. observe that orientation summary refreshes less often than local play detail
3. verify that the minimap remains useful between orientation refresh ticks
4. measure average transport cost over a deterministic tick window

## Done Criteria

- local play detail still updates every tick
- orientation data updates at a lower explicit cadence
- average per-client transport cost drops below the current post-culling baseline
- the client still orients correctly in the larger world
- no gameplay semantics change

## Out Of Scope Follow-Ups

- delta snapshots
- compression
- binary transport
- protocol version negotiation
- client-side interpolation systems
