# Slice: Initial Food-Seeking Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract carrying the resulting autonomous movement through ordinary snapshots

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for non-player movement policy.

This slice does not add scripted roles or a general AI system. It gives autonomous circles one minimal environment-driven steering rule so the world produces more meaningful resource competition.

## Discovery Scope

Establish the smallest autonomous behavior that makes circles respond to food rather than only drifting on fixed rails:

- an autonomous circle may steer toward a reachable food target
- when no food target is available, it falls back to the current deterministic movement baseline
- the behavior remains cheap, deterministic, and inspectable

This slice does **not** attempt to implement:

- behavior trees
- pathfinding graphs
- flocking
- threat evasion
- mate seeking
- lineage-specific behavior

## Why This Slice Next

The current implementation already supports:

- energy spending and recovery
- food regeneration
- fight and reproduction resolution
- lineage continuity
- energy-gated reproduction with child reserve payment

But autonomous circles still behave mostly as moving fixtures. That weakens:

- resource competition
- encounter density
- the “system-driven ecosystem” identity
- the source claim that behavior should emerge from rules and environment rather than explicit AI roles

Minimal food-seeking autonomy is the narrowest next step because it:

- increases meaningful interaction without adding a full AI subsystem
- stays aligned with food and energy as the central loop
- makes the world feel less staged
- remains deterministic and cheap to test

## Use-Case Contract

### Use Case

`SteerAutonomousCircleTowardFood`

### Primary Actor

Autonomous circles participating under the same energy and collision rules as the player.

### Pre-conditions

- the world contains active food slots
- autonomous circles already move and spend energy each tick
- the server is authoritative for movement outcomes

### Trigger

The server advances a simulation tick for an autonomous circle.

### Success Outcome

- an autonomous circle selects a deterministic food target when one is relevant
- the circle steers toward that target using the existing movement rules
- food collection becomes more responsive to environment state rather than fixed direction alone

### Failure Or Rejection Cases

- if movement becomes random or opaque, the slice loses inspectability
- if food-seeking bypasses ordinary movement cost, fairness is weakened
- if the behavior requires heavy AI machinery, the slice exceeds scope

## Main Business Rules

1. Autonomous movement remains authoritative server-side behavior.
2. Food-seeking must be deterministic.
3. Food-seeking uses the same movement speed and energy cost as all other movement.
4. When no suitable food target exists, autonomous circles fall back to the simpler baseline policy.
5. This slice does not add threat avoidance or reproduction-seeking behavior.

## Minimal Domain Concepts In Scope

- `Autonomous Steering`
- `Food Target`
- `Deterministic Target Selection`
- `Movement Intent`
- `World Snapshot`

## Bounded Autonomy Interpretation

This slice chooses the smallest inspectable steering rule:

- an autonomous circle prefers the nearest active food by deterministic tie-break
- the target is recomputed from world state rather than stored as hidden AI memory
- movement still occurs one ordinary tick at a time

This avoids pathfinding, black-box AI, and role systems while still making behavior environment-driven.

## Required Runtime Contract Changes

No new message types are required.

The existing snapshot contract should remain sufficient because the visible result is ordinary circle position, energy, and food changes across ticks.

## Required Ports Or Boundaries

- server-side target selection for autonomous circles
- server-side steering logic using existing movement mechanics
- client-side rendering through ordinary snapshots
- deterministic tests covering target choice and food collection behavior

## Build Guidance

- keep target selection local and explicit
- prefer nearest-food selection with deterministic tie-breakers
- do not introduce probabilistic wandering in this slice
- preserve the current bounded-world rules and energy cost rules
- avoid storing complex internal AI state if the same result can be recomputed from current world state

## Initial Test Plan

### Server tests

- an autonomous circle steers toward the nearest active food when one is available
- when two foods are equally viable, target choice is deterministic
- food-seeking still consumes the normal movement energy
- autonomous food collection becomes a consequence of steering rather than only initial lane placement

### Contract tests

- the existing snapshot schema remains sufficient

### Integration tests

- the client receives snapshots where autonomous circles visibly steer toward food
- the client can observe autonomous food collection in a non-trivial path rather than only fixed-lane drift

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and sees food plus autonomous circles
2. an autonomous circle selects a nearby food target
3. later snapshots show the circle curving or redirecting toward the food
4. the food is collected through ordinary authoritative movement and feeding rules

## Done Criteria

- autonomous circles can steer toward food deterministically
- steering remains cheap and inspectable
- the existing movement and energy rules still apply
- tests cover target selection and resulting food collection
- the slice avoids introducing a general AI subsystem

## Out Of Scope Follow-Ups

- threat evasion
- mate seeking
- flocking
- pathfinding
- strategy by lineage or shape
