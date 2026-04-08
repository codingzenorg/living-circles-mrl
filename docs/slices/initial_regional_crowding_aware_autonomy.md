# Slice: Initial Regional Crowding Aware Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots show the consequences of regional-crowding-aware autonomous steering without requiring a new protocol surface by default

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice follows the new regional crowding-energy rule: autonomous circles should not keep steering as if broad dense regions were free once the world now penalizes them over time.

## Discovery Scope

Establish the smallest deterministic steering adjustment that respects broader regional crowding pressure:

- autonomous circles should prefer movement that exits or avoids denser regions when a similarly valid nearby alternative exists
- the rule should remain bounded and deterministic
- the existing local crowding-aware steering, food seeking, interaction seeking, threat avoidance, and blocked-reproduction avoidance should remain in place as the current steering layers
- movement, food placement, regeneration, fight, reproduction, continuity, and child ownership remain unchanged

This slice does **not** attempt to implement:

- pathfinding toward sparse zones
- strategic migration systems
- territory ownership
- player steering changes
- new client-side visualization
- a full rewrite of autonomous decision-making

## Why This Slice Next

The world now differentiates between:

- immediate local crowding cost
- broader regional density cost

That improves ecological consequence, but it creates a new coherence gap: autonomous steering already reacts to local crowding, while still behaving as if broader dense regions carry no extra ongoing downside.

If circles pay a regional density cost but only steer against immediate local crowding, the model remains only partly integrated:

- the energy model says dense regions matter over time
- the steering model still reacts mostly at the very local level

The next pressure is therefore to let autonomy respond minimally to regional pressure that already exists.

This slice is the narrowest next step because it:

- changes steering rather than outcome rules
- follows directly from the new regional crowding-energy behavior
- strengthens medium-term dispersal without inventing migration AI
- keeps the player under the same downstream world rules while only changing autonomous steering

## Use-Case Contract

### Use Case

`SteerAutonomousCircleAwayFromDenseRegion`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- regional crowding pressure already exists as an energy cost
- autonomous circles already use a layered steering stack
- the world already exposes enough nearby-circle geometry for deterministic regional reasoning

### Trigger

An autonomous circle determines its steering direction for the next tick.

### Success Outcome

- the autonomous circle avoids moving deeper into a denser region when the chosen bounded rule says a less-dense alternative should win
- the world gains a stronger route toward medium-term regional differentiation
- existing food, threat, interaction, and local crowding rules remain recognizable

### Failure Or Rejection Cases

- if autonomy ignores the regional pressure entirely, the new energy rule remains only half integrated
- if the slice replaces existing steering priorities wholesale, scope is exceeded
- if the steering rule becomes opaque or strategic, determinism and inspectability weaken

## Main Business Rules

1. Regional-crowding-aware steering remains authoritative server-side behavior.
2. The rule should use a bounded regional-density comparison, not a global population score.
3. The steering adjustment should remain deterministic.
4. Existing steering priorities should stay recognizable; this slice should not replace the autonomy model wholesale.
5. Player movement remains unchanged.
6. Food placement, regeneration, fight, reproduction, continuity, and child ownership remain unchanged.

## Minimal Domain Concepts In Scope

- `Autonomous Steering`
- `Regional Crowding`
- `Preferred Direction`
- `Deterministic Tie-Break`

## Bounded Steering Interpretation

This slice chooses the smallest useful interpretation of regional-crowding-aware autonomy:

- evaluate whether the currently preferred steering direction would move the circle into a denser nearby region
- allow one bounded less-dense alternative or repulsion adjustment to win when the regional difference is clear
- preserve the current steering stack rather than replacing it

This avoids turning the slice into route planning while still making autonomy more coherent with the ecological pressure model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the new behavior is visible through movement and resulting energy outcomes in ordinary snapshots.

Build should avoid new contract surface unless a tiny inspectability field becomes necessary to understand the steering rule during review.

## Required Ports Or Boundaries

- server-side autonomous steering logic
- deterministic tests that prove a regional-crowding-aware steering adjustment
- implementation notes that record the chosen regional rule and tie-break behavior

## Build Guidance

- prefer evolving the existing autonomy steering path rather than introducing another separate decision subsystem
- choose one simple regional-density adjustment
- keep the change compatible with the current local crowding-aware steering and other steering priorities
- avoid changing the regional crowding energy rule itself in the same slice

## Initial Test Plan

### Server tests

- an autonomous circle still follows the existing steering priorities when regional density is not meaningfully different
- an autonomous circle prefers a less-dense direction when a nearby alternative clearly reduces regional pressure
- the steering adjustment remains deterministic for the same world state

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe regional-crowding-aware steering through ordinary movement in snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle faces a movement option that leads deeper into a dense region
2. the server applies the bounded regional-crowding-aware steering rule
3. the circle chooses a less-dense direction when the rule clearly prefers it
4. the resulting world shows stronger medium-term dispersal without replacing the rest of autonomy

## Done Criteria

- autonomous steering responds in a bounded way to regional crowding
- existing steering priorities remain recognizable
- the rule is deterministic
- the rest of the simulation semantics remain unchanged
- tests cover both neutral and regional-crowding-aware steering cases

## Out Of Scope Follow-Ups

- pathfinding or migration AI
- player steering changes
- regional overlays
- territory systems
- fight or reproduction redesign
