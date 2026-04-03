# Slice: Initial Autonomous To Autonomous Interaction Resolution

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible interaction outcomes between non-player circles

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether autonomous circles can interact with each other under the same fight and reproduction rules already used for player-involved encounters.

This slice extends the current simulation from a player-centered interaction graph toward a more fully shared ecosystem without introducing new runtime boundaries or a separate AI system.

## Discovery Scope

Establish the smallest deterministic rule that lets autonomous circles resolve encounters with each other:

- autonomous circles can detect overlap with other autonomous circles
- same-shape autonomous pairs resolve through the current fight path
- different-shape autonomous pairs resolve through the current reproduction path
- current child, energy, continuity, and contact-origin rules remain applicable where relevant

This slice does **not** attempt to implement:

- N-body interaction batching
- simultaneous multi-pair resolution in one tick beyond a deterministic first-resolved pair
- new autonomy strategies
- detached child autonomy
- removal of the current player-autonomous interaction path
- generalized population scaling

## Why This Slice Next

The current implementation has made orbiting children increasingly embodied, but the simulation still has one major structural asymmetry: interaction resolution is effectively player-centric. Autonomous circles coexist in the world, but they do not yet resolve encounters with each other under the same rules.

The model pressure is now:

- fairness through shared rules should apply across all circles, not only between player and autonomous circles
- the ecosystem should be capable of changing state without requiring player involvement
- food seeking, shape, children, and energy matter more when non-player circles can actually act on each other

This slice is the narrowest next step because it:

- reuses the current fight and reproduction outcomes
- keeps the contract small
- improves emergence without adding new entity types
- removes a major special-case asymmetry from the simulation

## Use-Case Contract

### Use Case

`ResolveAutonomousCircleInteraction`

### Primary Actor

Any autonomous circle that overlaps another autonomous circle during authoritative simulation.

### Pre-conditions

- at least two autonomous circles are active in the world
- overlap detection for player-involved encounters already exists
- current fight, reproduction, child, and continuity rules already exist

### Trigger

During authoritative overlap evaluation, one autonomous circle overlaps another autonomous circle by parent-body or currently valid child-based contact.

### Success Outcome

- the autonomous pair resolves through the same current same-shape or different-shape path used elsewhere
- the resulting world state is visible in later snapshots without requiring player involvement
- player and autonomous circles now participate in one more uniform rule set

### Failure Or Rejection Cases

- if autonomous circles still cannot resolve against each other, the ecosystem remains player-centered
- if autonomous-autonomous interaction uses a separate rule set, fairness and inspectability weaken
- if multiple autonomous overlaps create unstable or non-deterministic ordering, determinism weakens

## Main Business Rules

1. Autonomous circles may resolve interactions with other autonomous circles.
2. Same-shape autonomous pairs still resolve through the current fight path.
3. Different-shape autonomous pairs still resolve through the current reproduction path.
4. Current child absorption, child payment, reproduction gating, and continuity rules remain applicable where their preconditions are met.
5. Deterministic pair ordering is required when more than one autonomous-autonomous overlap exists in the same tick.
6. This slice does not require resolving more than one interaction per autonomous pair per overlap window.
7. Player-involved interaction rules remain unchanged in this slice.

## Minimal Domain Concepts In Scope

- `Autonomous Pair Contact`
- `Autonomous Fight Resolution`
- `Autonomous Reproduction Resolution`
- `Deterministic Pair Ordering`
- `World Snapshot`

## Bounded Interaction Interpretation

This slice chooses the smallest inspectable interpretation:

- extend the current interaction engine to include autonomous-autonomous pairs
- preserve existing fight and reproduction semantics after a pair is selected
- prefer deterministic pair ordering over attempting to resolve every possible overlap at once

This avoids the larger step of full many-body scheduling while still making the world less player-dependent.

## Required Runtime Contract Changes

The current contract is likely sufficient if autonomous-autonomous results are expressed through:

- the existing `interaction` object
- autonomous circle state changes already present in snapshots
- existing contact provenance fields when contact origin matters

Build should extend the contract only if autonomous-autonomous interaction cannot be understood from the current snapshot shape.

## Required Ports Or Boundaries

- server-side overlap detection that considers autonomous-autonomous pairs
- deterministic ordering when multiple autonomous pairs could resolve in one tick
- deterministic tests for same-shape and different-shape autonomous interactions
- client rendering that already exposes the resulting state changes

## Build Guidance

- prefer extending the current overlap-resolution entry points rather than creating a second simulation loop
- preserve current player-autonomous behavior
- choose one deterministic autonomous pair ordering rule and document it
- keep the current contact-origin model if it remains sufficient
- do not add new UI systems unless the result is unreadable in the demo

## Initial Test Plan

### Server tests

- two same-shape autonomous circles can resolve a fight without player involvement
- two different-shape autonomous circles can resolve reproduction without player involvement
- child-related fight and reproduction rules still apply for autonomous-autonomous encounters
- pair ordering is deterministic when more than one autonomous-autonomous overlap exists

### Contract tests

- the current snapshot schema remains sufficient unless build finds missing inspectability for autonomous-autonomous results

### Integration tests

- the client receives an autonomous-autonomous interaction outcome without needing player movement input
- resulting world changes are visible through ordinary snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two autonomous circles overlap without the player causing the encounter
2. the server resolves the encounter through the existing same-shape or different-shape path
3. the next snapshot shows the resulting autonomous world change

## Done Criteria

- autonomous circles can resolve interactions with other autonomous circles
- same-shape and different-shape semantics remain consistent with current rules
- deterministic pair ordering is explicit when multiple autonomous overlaps are possible
- player-involved behavior remains unchanged
- tests cover same-shape, different-shape, and ordering behavior

## Out Of Scope Follow-Ups

- simultaneous many-pair resolution per tick
- new steering policies for social targeting
- detached child autonomy
- population scaling beyond the current small-world model
- removing current radius shortcuts
