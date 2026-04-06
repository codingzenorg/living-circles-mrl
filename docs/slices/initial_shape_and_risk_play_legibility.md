# Slice: Initial Shape And Risk Play Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice shifts from ecosystem-rule expansion to player-facing legibility: the simulation now contains meaningful shape, fight, reproduction, food, and crowding semantics, but too much of that meaning still depends on reading dense labels and debug-like HUD text.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player infer danger and opportunity directly from the running world:

- same-shape conflict risk should become more immediately visible
- different-shape reproduction opportunity should become more immediately visible
- the player should be able to distinguish likely safe, risky, and promising nearby situations without relying mainly on text strings
- server-side world rules remain unchanged

This slice does **not** attempt to implement:

- new fight or reproduction semantics
- tutorial flows
- camera systems
- minimaps
- full UI redesign
- client-side prediction

## Why This Slice Next

Recent slices strengthened ecosystem validity through:

- larger world and population
- food-pressure rules
- local crowding pressure
- crowding-aware autonomy

That improves the world model, but EGD still flags a practical gap: a new player may only understand shape meaning, danger, and opportunity by reading the HUD instead of by perceiving the simulation directly.

The next pressure is therefore not another hidden server rule. It is to make the existing world semantics more legible at play speed.

This slice is the narrowest next step because it:

- changes presentation rather than world truth
- directly follows the EGD play-legibility finding
- preserves the authoritative model
- can make the richer simulation materially easier to evaluate and play

## Use-Case Contract

### Use Case

`RenderImmediateShapeAndRiskMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already exposes enough snapshot state to distinguish shape, energy, children, and current interaction outcomes
- the browser client already renders circles, food, and labels
- the world now contains richer local pressures and interactions than the early demo baseline

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- nearby same-shape conflict situations are visually more readable as danger or contest
- nearby different-shape situations are visually more readable as opportunity or blocked opportunity
- the player can infer more of the game state from the world view itself rather than mainly from text-heavy debug output

### Failure Or Rejection Cases

- if the slice only adds more text, the play-legibility goal is missed
- if the client invents its own world semantics instead of reflecting authoritative state, scope is exceeded
- if the rendering becomes visually noisy without clarifying meaning, readability worsens

## Main Business Rules

1. The server remains authoritative for all world and interaction semantics.
2. The client may use existing authoritative snapshot data to improve visual legibility.
3. The slice should prefer direct visual cues over more dense text output.
4. Same-shape conflict and different-shape reproduction meanings should become easier to distinguish at a glance.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Shape Meaning`
- `Danger`
- `Opportunity`
- `Visual Cue`
- `Readable Snapshot`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of legibility:

- use visual treatments such as outlines, indicators, or simple cueing around nearby circles and current interactions
- keep the rendering grounded in existing authoritative snapshot data
- avoid turning the client into a second rules engine

This avoids a broad UI redesign while still making the simulation more interpretable during ordinary play.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- shape
- energy
- attached children
- interaction outcomes

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen readability cues

## Build Guidance

- prefer a small number of strong visual cues over a large number of weak ones
- keep the cues consistent with the current dark-mode visual language
- avoid relying on textual overlays as the primary fix
- do not introduce client-side semantic guessing beyond bounded use of already authoritative fields

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- same-shape nearby circles should read more clearly as conflict pressure
- different-shape nearby circles should read more clearly as reproduction opportunity or blocked opportunity
- the player should need less HUD reading to understand what is happening nearby

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with nearby circles and food
2. the client renders authoritative state with stronger shape/risk legibility
3. the player can more quickly read same-shape threat, different-shape opportunity, and local pressure from the world view
4. the richer simulation becomes easier to evaluate without adding new server rules

## Done Criteria

- the client presents clearer shape and risk meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- tutorial systems
- minimaps
- camera redesign
- new server-side mechanics
- full HUD replacement
