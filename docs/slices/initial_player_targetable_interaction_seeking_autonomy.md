# Slice: Initial Player-Targetable Interaction Seeking Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that may now choose the player as an interaction target

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous target selection.

This slice extends the current deterministic interaction-seeking autonomy so the player can become part of the same target-selection world rather than remaining outside that steering rule.

## Discovery Scope

Establish the smallest deterministic rule that lets autonomous circles choose the player as an interaction target:

- the player becomes an eligible interaction target under a bounded steering rule
- target selection remains deterministic and server-authoritative
- current movement, energy, contact, fight, and reproduction rules remain unchanged after contact is reached
- food priority remains part of the model instead of being discarded

This slice does **not** attempt to implement:

- explicit threat or fear logic
- tactical avoidance of stronger circles
- player-specific aggression systems
- client-side anticipation or prediction
- personality systems
- separate pursuit rules for fight versus reproduction

## Why This Slice Next

The current implementation lets autonomous circles actively create encounters with each other, which makes the ecosystem more self-propelled. But the player is still excluded from that steering layer. That leaves one more structural asymmetry in a game whose model hypothesis says the player should be one participant under the same rules rather than a privileged outsider.

The model pressure is now:

- the player should be targetable under the same broad steering world as other circles
- the ecosystem should be able to engage the player without requiring the player to initiate everything
- fairness-through-shared-rules is stronger when the player is not exempt from non-player targeting

This slice is the narrowest next step because it:

- changes only autonomous target eligibility
- keeps the current fight and reproduction outcomes unchanged
- preserves deterministic steering
- improves parity without adding new entity types or strategy layers

## Use-Case Contract

### Use Case

`SteerAutonomousCircleTowardPlayerWhenEligible`

### Primary Actor

Any autonomous circle that is selecting its next movement direction.

### Pre-conditions

- autonomous interaction-seeking steering already exists
- the player is active in the world
- current player-autonomous interaction rules already exist

### Trigger

An autonomous circle advances a simulation tick and determines its steering direction.

### Success Outcome

- the autonomous circle may select the player as an interaction target under the bounded steering rule
- later snapshots can show autonomous circles initiating encounters with the player without requiring player movement to start them
- downstream fight and reproduction semantics remain unchanged once contact occurs

### Failure Or Rejection Cases

- if the player remains excluded from target selection, the world stays only partially shared
- if player targeting becomes hand-scripted or role-like, the model drifts away from emergence
- if player targeting is non-deterministic or overly privileged, inspectability and fairness weaken

## Main Business Rules

1. Autonomous target selection remains authoritative server-side behavior.
2. The player may become an eligible interaction target when active in the world.
3. Target choice must remain deterministic for the same world state and tick.
4. Food priority remains part of the model in this slice.
5. Once contact occurs, the current player-autonomous fight and reproduction rules remain unchanged.
6. This slice should not introduce player-only aggression mechanics.
7. Build should choose one bounded deterministic rule for when the player is eligible versus when other autonomous circles remain preferred.

## Minimal Domain Concepts In Scope

- `Autonomous Target Selection`
- `Player Eligibility`
- `Deterministic Nearest Target`
- `Interaction Opportunity`
- `World Snapshot`

## Bounded Targeting Interpretation

This slice chooses the smallest inspectable interpretation:

- the player is simply one more eligible circle in the current interaction-seeking target set under a documented rule
- target choice continues using deterministic distance and tie-breaking
- no new combat or reproduction semantics are introduced

This avoids the larger step of explicit hostility systems while still making the player more honestly part of the ecosystem.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing interaction outcomes once contact happens
- current rendering and HUD state

Build should extend the contract only if target-selection provenance is necessary to understand the new motion.

## Required Ports Or Boundaries

- server-side autonomous target-selection logic that may include the player
- deterministic tie-breaking across player and autonomous candidates
- tests that show autonomous circles can now initiate player encounters without player movement
- client rendering that remains sufficient to observe the resulting motion

## Build Guidance

- prefer evolving the current interaction-seeking target function rather than adding a separate pursuit layer
- preserve the current food-priority rule
- choose one narrow deterministic eligibility rule for including the player
- keep movement speed, energy cost, and outcome rules unchanged
- avoid inventing special hostility labels or behavior modes

## Initial Test Plan

### Server tests

- an autonomous circle can seek the player and create a same-shape encounter without player movement
- an autonomous circle can seek the player and create a different-shape encounter without player movement
- target choice is deterministic when both the player and another autonomous circle are eligible
- food priority still wins when nearby food satisfies the current threshold rule

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle initiating an encounter with the player without player movement input
- resulting player-autonomous encounters remain visible through the ordinary canvas and interaction HUD

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle selects the player as an eligible interaction target
2. it steers toward the player under ordinary movement rules
3. contact occurs and the existing player-autonomous interaction logic resolves
4. the resulting state change appears in ordinary snapshots without the player causing the encounter

## Done Criteria

- autonomous circles can deterministically target the player under a bounded rule
- food priority remains coherent
- current player-autonomous fight and reproduction semantics remain unchanged after contact
- player movement is not required to observe the new behavior
- tests cover player-targeting and deterministic choice

## Out Of Scope Follow-Ups

- avoidance behavior
- player-specific aggression tuning
- explicit threat scoring
- detached child autonomy
- removing current radius shortcuts
