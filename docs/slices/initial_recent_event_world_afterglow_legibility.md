# Slice: Initial Recent Event World Afterglow Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from persistent state readability to recent outcome readability. The simulation now communicates more about danger, resource pressure, autonomy, motion, and lineage, yet fight, reproduction, and continuity outcomes still disappear from the world view too quickly unless the player catches them in the exact moment or reads the encounter log.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player notice where a meaningful recent event just happened:

- same-shape fight outcomes should linger briefly in-world
- reproduction outcomes should linger briefly in-world
- continuity-preserving promotion should linger briefly in-world
- the cue should stay grounded in already-authoritative interaction outcomes
- the change should remain presentation-focused rather than adding new server rules

This slice does **not** attempt to implement:

- server-side event history
- replay systems
- permanent map markings
- combat effects systems
- new interaction rules
- new persistence or lineage mechanics

## Why This Slice Next

Recent slices established:

- stronger shape and interaction legibility
- stronger crowding-pressure legibility
- stronger food-pressure legibility
- stronger autonomy-intent legibility
- stronger player-motion legibility
- stronger lineage-continuity legibility

That means the running world is easier to read in the present, but one practical gap remains: important events still fade from player understanding the moment they resolve. The encounter log helps, but it is detached from world position. The next pressure is therefore to let major recent outcomes remain briefly visible where they happened.

## Use-Case Contract

### Use Case

`RenderRecentWorldOutcomeMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already exposes current interaction outcomes in snapshots
- the browser client already renders the current world state and a recent encounter log
- the runtime contract already contains enough outcome identity to distinguish recent fight, reproduction, and continuity events

### Trigger

A world snapshot with an active or resolved interaction is rendered on the client.

### Success Outcome

- recent important outcomes remain briefly visible in-world after they occur
- the player can connect the recent event log to concrete world locations
- ordinary play becomes easier to interpret without changing server authority

### Failure Or Rejection Cases

- if the slice only adds more text, the legibility goal is missed
- if the client invents hidden outcome history beyond bounded recent local memory, scope is exceeded
- if the visual aftermath effects become noisy or obscure ordinary world readability, readability worsens

## Main Business Rules

1. The server remains authoritative for all interaction outcomes.
2. The client may use current authoritative interaction snapshots to keep a very short-lived local aftermath cue.
3. The slice should prefer direct in-world cues over more HUD text.
4. Fight, reproduction, and continuity outcomes should remain visible a little longer in the scene.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Recent Outcome`
- `Fight Afterglow`
- `Reproduction Afterglow`
- `Continuity Afterglow`
- `Readable World Memory`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of recent outcome readability:

- keep a very short-lived client-local memory of the last few authoritative interaction outcomes
- render restrained location-tied afterglow cues where those outcomes occurred
- avoid turning the client into a replay timeline or full event history

This preserves the current dark-mode, immediate-play presentation while making important recent outcomes easier to understand.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- current interaction kind
- source and target identifiers
- current circle positions
- current continuity outcomes

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen aftermath-legibility cues

## Build Guidance

- prefer one short-lived in-world cue family over a broad effects system
- keep cues visually distinct but restrained
- tie each cue to current or very recent authoritative world positions
- avoid cluttering the scene during busy moments

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- recent fights should remain readable for a brief moment where they occurred
- recent reproduction should remain readable for a brief moment where it occurred
- continuity-preserving promotion should remain readable for a brief moment where it occurred
- the effect should help interpretation without obscuring ongoing play

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world and interacts with nearby circles
2. some fights, reproductions, or continuity-preserving outcomes resolve
3. the client renders brief recent-event afterglow cues at those world locations
4. the player can more easily connect recent outcomes to concrete places in the scene

## Done Criteria

- the client presents clearer recent-outcome meaning in-world
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- replay systems
- permanent world annotations
- full combat effects
- event timelines
- new server-side event-history mechanics
