# Slice: Initial Player Motion Responsiveness Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from environmental and actor-readability to moment-to-moment control readability. The simulation now communicates more about danger, opportunity, crowding, food, and nearby autonomous motion, yet the player’s own immediate motion still remains visually plain.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player feel and read their own authoritative movement more clearly:

- the player’s current motion direction should become easier to perceive at play speed
- immediate changes in movement should feel more visually acknowledged without client-side prediction
- the cue should stay grounded in already-authoritative player position updates
- the change should remain presentation-focused rather than changing movement rules

This slice does **not** attempt to implement:

- client-side prediction
- interpolation systems
- camera follow redesign
- speed changes
- new input semantics
- new server-side movement logic

## Why This Slice Next

Recent slices established:

- stronger shape and interaction legibility
- stronger crowding-pressure legibility
- stronger food-pressure legibility
- stronger nearby autonomy-intent legibility

That means the world and nearby actors are becoming easier to read, but the player’s own authoritative motion still lacks immediate visual reinforcement. The EGD findings already flagged responsiveness and playability as not yet evaluated as experience qualities. A small client-only slice that makes movement changes easier to perceive is the narrowest next step before broader responsiveness evaluation.

## Use-Case Contract

### Use Case

`RenderImmediatePlayerMotionMeaning`

### Primary Actor

The player actively steering inside the running world.

### Pre-conditions

- the authoritative server already updates player position through the existing movement flow
- the browser client already renders the player distinctly from autonomous circles
- the runtime contract already exposes the player position needed for bounded motion-legibility cues

### Trigger

A world snapshot is rendered on the client while the player is active.

### Success Outcome

- current player movement direction becomes easier to perceive
- changes in steering feel more readable without changing authority or simulation rules
- the player can connect input to authoritative motion more easily during play

### Failure Or Rejection Cases

- if the slice adds client-side prediction or hides authority lag, scope is exceeded
- if the cue only adds more text, the legibility goal is missed
- if the cue visually conflicts with the current danger, crowding, and food layers, readability worsens

## Main Business Rules

1. The server remains authoritative for all movement state.
2. The client may use existing authoritative snapshot data to improve motion readability.
3. The slice should prefer direct visual cues over added text.
4. The player’s current movement direction should become easier to perceive.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Player Motion`
- `Authoritative Position Change`
- `Readable Direction Cue`
- `Immediate Steering Feedback`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of motion readability:

- use short-range directional emphasis, motion trail, or related local player cues derived from recent authoritative positions
- keep the cue tightly attached to the player rather than adding a broader UI system
- avoid implying predicted future position

This preserves the current authoritative client/server shape while making moment-to-moment movement easier to read.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- authoritative player position
- tick progression
- current world state

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen motion-legibility cues

## Build Guidance

- prefer one or two tight player-local cues over a larger UI addition
- keep the cue visually compatible with the existing dark-mode canvas language
- avoid masking authoritative movement behavior with prediction-like effects
- keep attention on readability of current motion, not future pathing

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- the player’s movement direction should be easier to read at a glance
- quick steering changes should feel more legible than before
- the cue should not overwhelm danger, crowding, or food readability

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world and begins moving
2. authoritative snapshots update the player position as usual
3. the client renders a stronger local cue for current player motion
4. the player can more easily perceive their own movement and steering changes without new server logic

## Done Criteria

- the client presents clearer player-motion meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- client-side prediction
- interpolation systems
- camera redesign
- new movement rules
- full HUD redesign
