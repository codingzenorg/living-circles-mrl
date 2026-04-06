# Slice: Initial Food Pressure Play Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from interaction and crowding readability to resource readability. The simulation now contains food scarcity, deterministic regeneration, and energy-dependent survival, yet nearby food pressure is still mostly legible only after repeated play.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player infer where recovery opportunity is currently present and where nearby space is becoming resource-poor:

- visible food opportunity should become easier to notice at play speed
- nearby food scarcity should become easier to notice before the player fully depletes their energy
- the cues should stay grounded in current authoritative food positions and absence, not hidden timers
- the change should remain presentation-focused rather than adding new server rules

This slice does **not** attempt to implement:

- new food or regeneration semantics
- minimaps
- global resource heatmaps
- client-side prediction of food return timing
- new movement, fight, or reproduction rules

## Why This Slice Next

Recent slices established:

- larger world and population
- population-scaled starting food
- depletion-scaled food regeneration
- local crowding pressure and crowding-aware autonomy
- stronger shape, interaction, and crowding legibility

That means the simulation now has meaningful resource pressure, but the player still has to infer too much of it indirectly. The world shows food pieces clearly when they are present, yet the broader feeling of “this nearby area is rich enough to recover” versus “this area is running dry” is still weak.

The next pressure is therefore to make current food opportunity and food scarcity more readable without inventing hidden client-side resource logic.

## Use-Case Contract

### Use Case

`RenderImmediateFoodPressureMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already owns food placement, consumption, and regeneration
- the browser client already renders food items and existing shape/crowding cues
- the runtime contract already exposes the visible food state needed for bounded legibility

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- nearby food-rich space becomes easier to recognize as recovery opportunity
- locally sparse space becomes easier to read as resource pressure
- the player can better connect energy management to visible world conditions

### Failure Or Rejection Cases

- if the slice only adds more text, the legibility goal is missed
- if the client guesses hidden regeneration timing or future spawn state, scope is exceeded
- if food readability cues visually fight with existing shape or crowding cues, readability worsens

## Main Business Rules

1. The server remains authoritative for all food semantics.
2. The client may use existing authoritative snapshot data to improve visual legibility.
3. The slice should prefer direct visual cues over additional text.
4. Nearby recovery opportunity and nearby scarcity should become easier to notice.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Food Opportunity`
- `Food Scarcity`
- `Visible Recovery Cue`
- `Readable Snapshot`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of food readability:

- use local visual treatments around visible food-rich nearby areas or the player's immediate vicinity
- keep the cues grounded in current visible food positions and local absence
- avoid timers, predictions, or whole-map analytics overlays

This preserves the current immediate-play presentation while making one more ecosystem pressure easier to perceive.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- food positions
- circle positions
- player position

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen food-legibility cues

## Build Guidance

- prefer one or two strong local cues over a broad overlay system
- keep the cues visually compatible with existing danger/opportunity and crowding signals
- avoid implying precise food-regeneration forecasts
- keep attention centered on nearby actionable space rather than the full map

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- nearby food-rich areas should read more clearly as recovery opportunity
- nearby sparse areas should feel more obviously resource-poor
- the player should need less trial-and-error to read local food pressure

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with multiple active circles and food slots
2. some nearby neighborhoods retain visible food while others become sparse
3. the client renders stronger cues for local food opportunity and scarcity
4. the player can more easily read where recovery is available and where resource pressure is higher

## Done Criteria

- the client presents clearer food-pressure meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- minimaps
- global resource heatmaps
- client-side regeneration prediction
- new server-side food mechanics
- full HUD redesign
