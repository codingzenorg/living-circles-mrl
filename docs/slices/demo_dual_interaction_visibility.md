# Slice: Demo Dual Interaction Visibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for multiple autonomous circles and their current interaction outcomes

## Architecture Mode

Explicit client/server boundary with the server hosting a small deterministic world that demonstrates both implemented interaction paths during normal manual play.

This slice is demo-oriented, but it must still preserve the same bounded, deterministic server authority model.

## Discovery Scope

Establish the smallest default world setup that makes both currently implemented interaction behaviors visible in the live demo:

- one player-controlled circle remains in the world
- one autonomous circle shares the player’s shape and can resolve into a same-shape fight
- one autonomous circle has a different shape and can still produce a `reproduce_candidate` overlap
- the client can render all circles distinctly enough to make the demo readable

This slice does **not** attempt to implement:

- actual reproduction outcomes
- additional fight rules
- child accumulation
- continuity or replacement
- random spawning or procedural scene setup
- an arbitrary number of autonomous circles

## Why This Slice Next

The current implementation contains both:

- same-shape fight resolution
- different-shape interaction classification

But the default live demo only makes one path obvious at a time, depending on the chosen shapes in the initial world. That reduces the value of the browser as a review surface.

This slice makes the current model more legible by default without introducing new domain semantics.

## Use-Case Contract

### Use Case

`RunVisibleDualInteractionDemo`

### Primary Actor

Player controlling one circle while two deterministic autonomous circles share the world.

### Pre-conditions

- a server process can host one bounded world
- the world contains one player-controlled circle
- the world contains one same-shape autonomous circle and one different-shape autonomous circle
- the initial placement and movement policy are deterministic
- a client can open a WebSocket connection to the server

### Trigger

The server advances ticks while the player moves through the world and autonomous circles follow deterministic policies.

### Success Outcome

- snapshots contain both autonomous circles
- manual play can naturally expose a same-shape fight path
- manual play can naturally expose a different-shape `reproduce_candidate` path
- the demo remains deterministic and inspectable

### Failure Or Rejection Cases

- if both autonomous circles use the same shape, the demo loses the reproduction-candidate visibility goal
- if both autonomous circles use different shapes from the player, the demo loses the same-shape fight visibility goal
- if the world becomes unreadable because circles overlap immediately at start, the initial placement should be adjusted rather than hiding the issue in the client

## Main Business Rules

1. The default world contains exactly one player-controlled circle, one same-shape autonomous circle, and one different-shape autonomous circle.
2. All circles still follow the same movement and energy rules.
3. Same-shape overlap still resolves as a deterministic fight.
4. Different-shape overlap still remains classification-only as `reproduce_candidate`.
5. The slice is about visibility and reviewability, not new combat or reproduction semantics.

## Minimal Domain Concepts In Scope

- `Player Circle`
- `Autonomous Circle`
- `Shape`
- `Fight Resolution`
- `Reproduce Candidate`
- `Default Demo World`
- `World Snapshot`

## Required Runtime Contract Changes

The existing snapshot contract already supports multiple autonomous circles. This slice may not need a new contract shape, but it must ensure:

- multiple autonomous circles appear in default snapshots
- the client can distinguish same-shape and different-shape autonomous circles clearly enough for review
- any active or resolved interaction still maps unambiguously to the involved circle identifiers

## Required Ports Or Boundaries

- server-side deterministic world initialization for two autonomous circles
- deterministic movement setup that allows both interaction paths to be observed during manual play
- client-side rendering distinctions for the additional autonomous circle
- tests covering default world composition and demo-path visibility assumptions

## Build Guidance

- keep the change in world initialization and client visibility, not in core domain semantics
- prefer exactly two autonomous circles for now
- make the same-shape autonomous circle easy to recognize in the client
- keep different-shape interaction unresolved, exactly as in the current implementation
- do not turn the demo setup into a configurable scenario framework yet

## Initial Test Plan

### Server tests

- the default world contains two autonomous circles
- one autonomous circle matches the player shape
- one autonomous circle differs from the player shape
- same-shape overlap still resolves as a fight
- different-shape overlap still remains classification-only

### Contract tests

- the snapshot schema still covers multiple autonomous circles without further contract changes

### Integration tests

- the initial snapshot contains two autonomous circles
- one autonomous circle can be identified as same-shape relative to the player
- one autonomous circle can be identified as different-shape relative to the player
- at least one same-shape fight path and one different-shape classification path remain reachable in deterministic tests

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot with one player circle and two autonomous circles
2. one autonomous circle visibly matches the player’s shape
3. one autonomous circle visibly differs from the player’s shape
4. moving toward the same-shape circle can trigger a fight
5. moving toward the different-shape circle can trigger a `reproduce_candidate` classification

## Done Criteria

- the default live demo includes one same-shape and one different-shape autonomous circle
- the client renders the world clearly enough for a reviewer to tell which interaction path is being exercised
- same-shape fight behavior still works
- different-shape classification-only behavior still works
- tests cover the default world composition and both visible interaction paths
- the slice does not add reproduction resolution, children, continuity, or new combat semantics

## Out Of Scope Follow-Ups

- actual reproduction resolution
- more than two autonomous circles
- scenario selection UI
- configurable debug presets
- growth, children, or lineage behavior
