# Slice: Initial Autonomy Intent Play Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from world-state readability to actor-readability. The simulation now exposes shape meaning, crowding pressure, and food pressure more clearly, yet nearby autonomous circles still often look like they are moving for opaque reasons.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player infer why a nearby autonomous circle is currently moving the way it is:

- pursuit of food should become easier to distinguish from social pursuit
- retreat behavior should become easier to distinguish from ordinary movement
- the cues should stay grounded in already-visible world state and existing authoritative behavior
- the change should remain presentation-focused rather than adding new server rules

This slice does **not** attempt to implement:

- new autonomy rules
- explicit server-side intent publication
- path visualization across long horizons
- behavior trees or AI debug tooling
- new fight, reproduction, food, or crowding semantics

## Why This Slice Next

Recent slices established:

- larger world and higher population
- stronger interaction legibility
- stronger crowding-pressure legibility
- stronger food-pressure legibility

That means the world itself is becoming more readable, but the active participants are still less legible than the environment around them. A nearby autonomous circle may be chasing food, seeking a feasible partner, avoiding a threat, or escaping blocked reproduction pressure, yet those modes remain easy to confuse during ordinary play.

The next pressure is therefore to make current autonomous behavior more interpretable without inventing hidden client-side AI or adding new server mechanics.

## Use-Case Contract

### Use Case

`RenderImmediateAutonomyIntentMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already applies bounded food-seeking, interaction-seeking, and avoidance behavior
- the browser client already renders circles, food, shape-risk cues, crowding cues, and food-pressure cues
- the runtime contract already exposes the visible world state needed for bounded interpretation

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- nearby autonomous circles become easier to read as pursuing, approaching, or retreating
- the player can infer more of the live system from motion and local cues rather than only after repeated observation
- the richer ecosystem remains authoritative while becoming easier to evaluate

### Failure Or Rejection Cases

- if the slice only adds more text, the legibility goal is missed
- if the client invents hidden server intent rather than boundedly reflecting visible behavior, scope is exceeded
- if the new cues overload the existing shape/crowding/food cues, readability worsens

## Main Business Rules

1. The server remains authoritative for all autonomous behavior.
2. The client may use existing authoritative snapshot data to improve visual legibility.
3. The slice should prefer direct local visual cues over added text.
4. Retreat, food pursuit, and social pursuit should become easier to distinguish nearby.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Autonomy Intent`
- `Food Pursuit`
- `Social Pursuit`
- `Retreat`
- `Readable Motion Cue`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of autonomy readability:

- use lightweight local motion cues, aim cues, or short-range markers around nearby autonomous circles
- ground the cue in current visible geometry and current direction of movement
- avoid claiming exact server-side reasoning when only bounded motion inference is available

This preserves the current dark-mode immediate-play presentation while making autonomous behavior easier to parse.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- current circle positions
- current food positions
- current player position
- shape and nearby world geometry

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen autonomy-legibility cues

## Build Guidance

- prefer a small number of readable motion cues over full debug annotations
- keep the cues local to nearby autonomous circles
- stay visually compatible with the existing danger, crowding, and food cues
- avoid implying more certainty than the visible movement and world state justify

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- nearby retreating circles should feel easier to identify
- nearby circles moving toward visible food should read differently from circles closing on social targets
- the player should need less guesswork to interpret local autonomous motion

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with multiple autonomous circles and food
2. some nearby circles move toward food, some toward other circles, and some away from local pressure
3. the client renders stronger local cues for those nearby motion modes
4. the player can more easily read why nearby circles are moving the way they are

## Done Criteria

- the client presents clearer nearby autonomy-intent meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- explicit server-published AI intent
- long-horizon path previews
- minimaps
- new autonomy rules
- full HUD redesign
