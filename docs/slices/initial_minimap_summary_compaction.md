# Slice: Initial Minimap Summary Compaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a player-following viewport and passive minimap
- Go server with authoritative simulation state
- shared runtime contract already split into high-cadence local detail and lower-cadence orientation refreshes

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps the current dual-cadence transport model, but reduces the cost of the slower orientation refreshes themselves. It does not introduce delta encoding or compression. It only makes the minimap summaries more compact.

## Discovery Scope

Reduce transport cost further by shrinking what the minimap/orientation refresh actually needs to describe.

This slice should:

- keep local viewport play detail unchanged
- keep the lower-cadence orientation refresh model
- send less data per orientation refresh
- preserve enough whole-world orientation for the current minimap

This slice does **not** attempt to implement:

- delta snapshots
- binary encoding
- compression
- zoom-aware minimap redesign
- server-side aggregation that changes gameplay semantics

## Why This Slice Next

The current transport path now averages about `17096` bytes/sec per client under the default expanded world, which is a major improvement from the original full-snapshot baseline. But the slower orientation refreshes still resend:

- every world food position
- every world autonomous circle position

That is still richer than the current minimap actually needs. The minimap mainly needs:

- approximate whole-world activity placement
- player position
- viewport window
- broad shape of food opportunity and remote population

The next pressure is therefore summary compactness, not cadence or compression.

This is the narrowest next step because it:

- preserves the current dual-cadence model
- avoids a larger protocol redesign
- should further reduce average bytes/sec without touching gameplay

## Use-Case Contract

### Use Case

`SendCompactOrientationSummaryForMinimap`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- local viewport detail already arrives every tick
- orientation data already refreshes at a lower cadence
- the current minimap consumes the latest valid orientation summary

### Trigger

The server emits a lower-cadence orientation refresh message.

### Success Outcome

- the minimap still orients the player within the larger world
- orientation refresh payload is smaller than the current full-summary refresh
- local viewport play remains unchanged
- average per-client transport cost drops below the current dual-cadence baseline

### Failure Or Rejection Cases

- if the minimap becomes too vague to orient the player, the slice failed
- if the slice starts requiring interpolation or prediction, scope drifted
- if the slice changes simulation behavior rather than orientation representation, scope is exceeded

## Main Business Rules

1. This is a compact-summary slice, not a gameplay slice.
2. Local viewport detail remains unchanged.
3. Orientation refreshes remain lower-cadence.
4. The minimap only needs enough detail to preserve orientation and broad world distribution.
5. Compaction should be deterministic and measurable.

## Minimal Domain Concepts In Scope

- `Compact Orientation Summary`
- `Minimap-Relevant World Summary`
- `Average Per-Client Transport Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful compaction interpretation:

- stop sending full per-entity minimap summaries when a coarser summary would do
- prefer coarse bins, counts, or reduced precision over exact remote world replay
- keep the minimap visually useful without pretending it is a second full-resolution world view

This avoids jumping into deltas or compression while still cutting unnecessary refresh detail.

## Required Runtime Contract Changes

Likely yes, but bounded.

The current minimap summary fields will likely need to evolve toward a more compact orientation representation. The slice should avoid turning that into a general protocol redesign.

## Required Ports Or Boundaries

- server-side orientation summary builder
- shared runtime contract definitions for the compact minimap summary
- browser client minimap consumption of the new compact summary
- transport measurement tests comparing current dual-cadence cost against the compacted version

## Build Guidance

- preserve the current local viewport transport untouched
- compact only the lower-cadence orientation payload
- prefer deterministic spatial bucketing or reduced precision over free-form heuristics
- measure the effect over the same deterministic cadence window already established

## Initial Test Plan

### Server tests

- compact orientation summaries remain deterministic
- average cost over a cadence window drops below the current dual-cadence baseline

### Contract tests

- the compact orientation summary remains explicit and parseable

### Integration tests

- the minimap still shows usable whole-world orientation
- the client continues to render local play correctly while consuming the compact summary

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect and receive normal local viewport updates
2. receive compact lower-cadence orientation refreshes
3. verify that the minimap still orients the player within the larger world
4. compare average transport cost against the current dual-cadence baseline

## Done Criteria

- local viewport play is unchanged
- orientation refreshes are more compact than the current dual-cadence summary
- minimap orientation remains useful
- average per-client transport cost drops below the current dual-cadence baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- delta snapshots
- compression
- binary transport
- adaptive per-entity prioritization
- client-side prediction
