# Slice: Initial Crowding Aware Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots show the consequences of crowding-aware autonomous steering without requiring a new protocol surface by default

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice follows the new local crowding-energy rule: autonomous circles should not keep steering as if dense local clustering were free once the world now penalizes it.

## Discovery Scope

Establish the smallest deterministic steering adjustment that respects local crowding pressure:

- autonomous circles should avoid steering deeper into locally crowded space when a similarly valid less-crowded option exists
- the rule should remain bounded and deterministic
- food seeking, interaction seeking, threat avoidance, and blocked-reproduction avoidance should remain in place as the existing steering layers
- movement, food placement, regeneration, fight, reproduction, continuity, and child ownership remain unchanged

This slice does **not** attempt to implement:

- full pathfinding
- flocking or boids behavior
- strategic group formation
- explicit territory seeking
- player steering changes
- crowding visualization UI

## Why This Slice Next

The world now applies an additional local energy cost in dense areas. That improves ecosystem pressure, but it also creates a new coherence gap: autonomous steering still behaves as if dense local concentration had no direct downside.

If circles pay a crowding cost but do not react to it at all, the system risks feeling artificially split:

- the energy model says density matters
- the steering model still acts like density is irrelevant except after the cost is already paid

The next pressure is therefore to let autonomy respond minimally to the world pressure that already exists.

This slice is the narrowest next step because it:

- changes steering rather than core outcome rules
- directly follows the new crowding-pressure behavior
- strengthens ecosystem plausibility without inventing a general AI subsystem
- keeps the player under the same downstream world rules while only changing autonomous steering

## Use-Case Contract

### Use Case

`SteerAutonomousCircleAwayFromExcessCrowding`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- local crowding pressure already exists as an energy cost
- autonomous circles already choose among food, interaction, avoidance, and fallback steering paths
- the world already exposes enough nearby-circle geometry for deterministic local reasoning

### Trigger

An autonomous circle determines its steering direction for the next tick.

### Success Outcome

- the autonomous circle avoids moving deeper into a locally crowded area when the chosen bounded rule says a less-crowded direction should win
- the world gains a stronger route toward dispersal and spatial differentiation
- existing food, threat, and interaction rules remain recognizable

### Failure Or Rejection Cases

- if autonomy keeps steering into crowding as if there were no local pressure, the new energy rule remains only half integrated
- if the slice replaces existing steering priorities entirely, scope is exceeded
- if the steering rule becomes opaque or over-strategic, determinism and inspectability weaken

## Main Business Rules

1. Crowding-aware steering remains authoritative server-side behavior.
2. The rule should be based on local nearby-circle density, not on a global population score.
3. The steering adjustment should remain bounded and deterministic.
4. Existing steering priorities should stay recognizable; this slice should not replace the autonomy model wholesale.
5. Player movement remains unchanged.
6. Food placement, regeneration, fight, reproduction, continuity, and child ownership remain unchanged.

## Minimal Domain Concepts In Scope

- `Autonomous Steering`
- `Local Crowding`
- `Preferred Direction`
- `Deterministic Tie-Break`

## Bounded Steering Interpretation

This slice chooses the smallest useful interpretation of crowding-aware autonomy:

- evaluate whether the currently preferred steering direction would move the circle into a denser local neighborhood
- allow one bounded less-crowded alternative or repulsion adjustment to win when the density difference is clear
- preserve the current steering stack rather than replacing it

This avoids turning the slice into pathfinding or social AI while still making the autonomy model more coherent with the pressure model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the new behavior is visible through movement and resulting energy outcomes in ordinary snapshots.

Build should avoid new contract surface unless a tiny inspectability field becomes necessary to understand the steering rule during review.

## Required Ports Or Boundaries

- server-side autonomous steering logic
- deterministic tests that prove a crowding-aware steering adjustment
- implementation notes that record the chosen local rule and tie-break behavior

## Build Guidance

- prefer evolving the existing autonomous-intent logic rather than introducing a second decision subsystem
- choose one simple local-density adjustment
- keep the change compatible with the current food and interaction priorities
- avoid changing the crowding energy rule itself in the same slice

## Initial Test Plan

### Server tests

- an autonomous circle still follows the existing steering priorities when local density is not meaningfully different
- an autonomous circle prefers a less-crowded direction when a nearby alternative clearly reduces local crowding
- the steering adjustment remains deterministic for the same world state

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe crowding-aware steering through ordinary movement in snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle faces a locally crowded movement option
2. the server applies the bounded crowding-aware steering rule
3. the circle chooses a less-crowded direction when the rule clearly prefers it
4. the resulting world shows stronger dispersal behavior without replacing the rest of autonomy

## Done Criteria

- autonomous steering responds in a bounded way to local crowding
- existing steering priorities remain recognizable
- the rule is deterministic
- the rest of the simulation semantics remain unchanged
- tests cover both neutral and crowding-aware steering cases

## Out Of Scope Follow-Ups

- full pathfinding
- flocking/group behavior
- player steering changes
- crowding overlays
- fight or reproduction redesign
