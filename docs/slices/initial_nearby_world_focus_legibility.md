# Slice: Initial Nearby World Focus Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from individual cues to spatial focus. The simulation now communicates much more about local danger, food, crowding, motion, continuity, and recent outcomes, yet the larger world can still feel visually busy because near and far activity are rendered with similar emphasis.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player focus on nearby actionable space:

- nearby space around the player should feel more visually primary
- distant non-actionable activity should feel more backgrounded without disappearing
- the cue should stay grounded in current visible geometry and current player position
- the change should remain presentation-focused rather than changing world rules

This slice does **not** attempt to implement:

- camera zoom systems
- minimaps
- hiding distant entities completely
- culling or gameplay-affecting visibility rules
- new server-side proximity mechanics

## Why This Slice Next

Recent slices established:

- stronger shape and interaction legibility
- stronger crowding-pressure legibility
- stronger food-pressure legibility
- stronger autonomy-intent legibility
- stronger player-motion legibility
- stronger lineage-continuity legibility
- brief recent-event afterglow
- cleaner out-of-canvas player and NPC information

That means the game now exposes much more useful information, but the canvas also carries more simultaneous signals. In the expanded world, readability will improve further if the player can more easily distinguish nearby relevant space from distant background motion.

The next pressure is therefore to improve local focus, not by adding more semantic cues, but by subtly organizing the scene around the player’s current neighborhood.

## Use-Case Contract

### Use Case

`RenderNearbyWorldAsPrimaryContext`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the browser client already renders a larger world with multiple active circles, food, and several bounded legibility cues
- the runtime contract already exposes the player position and all visible entities
- the server remains authoritative for all game state

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- nearby actionable space becomes easier to scan
- distant background activity remains visible but less dominant
- the current legibility cues become easier to use together during ordinary play

### Failure Or Rejection Cases

- if distant world state disappears entirely, scope is exceeded
- if the slice adds new world rules or hidden gameplay constraints, scope is exceeded
- if the focus treatment visually fights with current danger, crowding, food, motion, lineage, or afterglow cues, readability worsens

## Main Business Rules

1. The server remains authoritative for all world state.
2. The client may use existing authoritative snapshot data to improve focus and readability.
3. The slice should prefer a subtle spatial emphasis rather than new dense overlays.
4. Nearby actionable space should become easier to scan than distant space.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Nearby Context`
- `Distant Background`
- `Readable Focus`
- `Player-Centered World Legibility`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of focus:

- use player-centered emphasis such as subtle falloff, stronger nearby contrast, or restrained background de-emphasis
- keep all entities visible
- avoid introducing a minimap, fog of war, or visibility gameplay rule

This preserves the current immediate-play presentation while making the expanded world easier to read.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- player position
- circle positions
- food positions
- current world bounds

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen focus cues

## Build Guidance

- prefer one restrained focus treatment over multiple competing effects
- keep the effect subtle enough that it organizes the scene instead of overpowering it
- preserve visibility of the full world
- stay visually compatible with the current dark-mode canvas language

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- nearby circles and food should feel easier to scan than distant background motion
- distant world state should still remain visible
- the effect should support, not replace, the existing danger/food/crowding/intent/lineage/event cues

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with multiple active circles and food across a wider map
2. nearby and distant activity coexist in the same frame
3. the client renders a subtle player-centered focus treatment
4. the player can more easily read nearby actionable space without losing global awareness

## Done Criteria

- the client presents clearer nearby-world focus
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- minimaps
- camera zoom systems
- fog of war
- server-side visibility rules
- full HUD redesign
