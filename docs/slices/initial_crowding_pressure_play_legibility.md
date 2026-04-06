# Slice: Initial Crowding Pressure Play Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from shape-meaning readability to ecosystem-pressure readability. The simulation now contains local crowding energy cost and crowding-aware autonomy, yet that pressure is still mostly invisible at play speed.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player infer when a local area is energetically costly because it is too crowded:

- nearby dense local clusters should become more legible as costly pressure
- the player should be able to notice that some spatial neighborhoods are expensive before reading HUD text
- current crowding semantics should remain authoritative and unchanged
- the change should stay presentation-focused rather than adding new server rules

This slice does **not** attempt to implement:

- new crowding rules
- minimaps
- camera systems
- pathfinding overlays
- general heatmaps for the entire world
- new fight or reproduction semantics

## Why This Slice Next

Recent slices established:

- larger world and population
- depletion-scaled food pressure
- local crowding energy pressure
- crowding-aware autonomous steering
- stronger shape and interaction legibility

That means the simulation now has a meaningful local survival cost that remains hard to perceive directly from the world view. The player can often see a cluster, but not whether that cluster is the kind of density that matters to the authoritative energy rule.

The next pressure is therefore to make existing crowding pressure more readable without turning the client into a second simulation.

## Use-Case Contract

### Use Case

`RenderImmediateCrowdingPressureMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already applies local crowding energy cost
- the browser client already renders circles, food, shape-risk cues, and a recent encounter log
- the runtime contract already exposes enough ordinary world state to estimate visible local clustering

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- costly local clustering becomes easier to notice directly on the canvas
- the player can better distinguish sparse and pressured nearby space
- the richer ecosystem model becomes easier to evaluate during ordinary play

### Failure Or Rejection Cases

- if the slice only adds more HUD text, the legibility goal is missed
- if the client invents hidden crowding semantics beyond bounded use of authoritative positions, scope is exceeded
- if the world view becomes noisy or cluttered, readability worsens

## Main Business Rules

1. The server remains authoritative for all crowding semantics.
2. The client may use existing authoritative snapshot data to improve visual legibility.
3. The slice should prefer direct visual cues over added textual explanation.
4. Dense local areas that are likely costly under the current rule should become easier to notice.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Crowding Pressure`
- `Dense Local Area`
- `Visible Cost Cue`
- `Readable Snapshot`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of crowding readability:

- use local visual treatments such as halos, soft fields, or emphasis around crowded nearby clusters
- keep the cues grounded in current visible circle positions rather than hidden future simulation
- avoid full-world overlays or analytics-style visualization

This preserves the current dark-mode, immediate-play presentation while making one important ecosystem rule easier to perceive.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- circle positions
- player position
- current world bounds

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen crowding-legibility cues

## Build Guidance

- prefer one or two strong visual cues over a broad overlay system
- keep the cues local to nearby relevant space rather than the full world
- stay consistent with the current dark-mode canvas language
- avoid implying more precision than the current authoritative crowding rule actually has

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- locally dense areas near the player should read more clearly as costly zones
- sparse nearby areas should remain visually calmer
- the player should need less HUD reading to understand local crowding pressure

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with multiple active circles and food
2. local crowding develops around some neighborhoods
3. the client renders stronger cues for costly nearby clustering
4. the player can more easily avoid or interpret dense pressure zones during play

## Done Criteria

- the client presents clearer crowding-pressure meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- minimaps
- camera redesign
- global density maps
- new server-side crowding mechanics
- full HUD redesign
