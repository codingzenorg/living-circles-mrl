# Slice: Initial Local Transport Precision Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with a pixel-based 2D canvas viewport
- Go server with authoritative simulation state
- shared runtime contract already optimized by viewport culling, dual cadence, and compact minimap summaries

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps the current transport shape and cadence strategy, but reduces how much numeric precision is sent for local play data. It does not alter simulation precision internally. It only changes the serialized transport representation where the client already renders against a pixel grid.

## Discovery Scope

Reduce per-client transport cost further by trimming unnecessary numeric precision in the high-cadence local payload.

This slice should:

- keep the current local-vs-orientation transport structure
- keep the current cadences
- keep authoritative server-side simulation precision unchanged
- reduce the serialized precision of local transport numbers to the smallest useful display precision

This slice does **not** attempt to implement:

- binary encoding
- compression
- delta snapshots
- client-side prediction
- changes to simulation math or gameplay thresholds

## Why This Slice Next

The transport path has already become much cheaper:

- full snapshot baseline: about `64870` bytes/sec per client
- viewport-culling baseline: about `39070` bytes/sec
- dual-cadence compact-summary baseline: about `15854` bytes/sec

At this point, the main remaining chatty path is the high-cadence local payload, which still serializes:

- floating-point world positions
- floating-point energy values
- floating-point radii and child positions

But the client currently renders onto a pixel-based canvas and mostly displays whole-number values in support panels. That means some transport precision is likely no longer buying real user-visible value.

This slice is the narrowest next step because it:

- keeps the current transport structure intact
- avoids the complexity of deltas or compression
- directly targets the still-frequent local payload

## Use-Case Contract

### Use Case

`SendDisplay-SufficientLocalNumericPrecision`

### Primary Actor

The Go server broadcasting authoritative state to one browser client.

### Pre-conditions

- local viewport detail already arrives every tick
- the browser already renders on a pixel-based canvas
- transport measurement already exists for the current optimized baseline

### Trigger

The server serializes the next local-detail transport snapshot.

### Success Outcome

- the client still renders local play correctly
- the client does not lose meaningful readability or motion clarity
- local payload size drops below the current compact-summary baseline
- authoritative simulation precision remains unchanged inside the server

### Failure Or Rejection Cases

- if visible motion becomes jittery or misleading, the precision was reduced too far
- if the slice changes simulation math instead of serialized representation, scope is exceeded
- if the slice requires version negotiation or major compatibility logic, scope drifted

## Main Business Rules

1. This is a transport-representation slice, not a simulation-math slice.
2. The server remains authoritative at full internal precision.
3. Only serialized local transport numbers may be reduced.
4. Precision should remain sufficient for current viewport rendering and panel readability.
5. The optimization should be measurable against the current compact-summary baseline.

## Minimal Domain Concepts In Scope

- `Display-Sufficient Local Precision`
- `Authoritative Internal Precision`
- `Serialized Local Payload Cost`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful precision interpretation:

- round or trim local transport values to the smallest precision the current client actually uses
- keep the minimap and orientation summaries unchanged in this slice
- keep interaction semantics and server-side calculations unchanged

This avoids jumping into more complex protocol changes while still reducing the highest-frequency payload.

## Required Runtime Contract Changes

Likely no schema-shape change.

The contract fields may stay the same, with only the serialized numeric precision changing.

## Required Ports Or Boundaries

- server-side transport serialization of local snapshot numbers
- browser client rendering validation against reduced local precision
- transport measurement tests comparing pre- and post-reduction local payload cost

## Build Guidance

- reduce precision at the transport boundary, not inside simulation state
- prefer deterministic rounding rules
- keep the local viewport path visually stable first
- measure the effect on the high-cadence local payload and overall average cost

## Initial Test Plan

### Server tests

- reduced-precision transport serialization remains deterministic
- local payload size drops below the current compact-summary baseline

### Contract tests

- schema shape remains explicit and unchanged

### Integration tests

- local viewport play remains readable and responsive
- movement and interaction rendering do not lose meaningful clarity

## Scenario Definition

Start a local server with the current viewport-mode client.

Scenario steps:

1. connect and receive local transport snapshots with reduced numeric precision
2. verify that local movement and interaction rendering remain visually correct
3. compare transport measurements against the current compact-summary baseline

## Done Criteria

- local transport numbers are sent with lower but still sufficient precision
- local play remains readable
- measured transport cost drops below the current compact-summary baseline
- no gameplay semantics change

## Out Of Scope Follow-Ups

- deltas
- compression
- binary transport
- simulation math changes
- prediction or interpolation systems
