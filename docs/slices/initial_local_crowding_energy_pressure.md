# Slice: Initial Local Crowding Energy Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots can show the consequences of denser local populations without adding a new protocol surface by default

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for movement cost and world-state consequences.

This slice introduces one bounded ecosystem pressure after world expansion and food-pressure work: local crowding should increase survival pressure instead of letting denser worlds behave like the same mechanics spread across a bigger map.

## Discovery Scope

Establish the smallest deterministic crowding rule that affects medium-term outcomes:

- circles should experience additional energy pressure when too many other active circles are nearby
- the rule should be local rather than global
- the rule should remain deterministic and inspectable
- movement, food placement, regeneration, fight, reproduction, continuity, child ownership, and steering rules remain otherwise unchanged

This slice does **not** attempt to implement:

- collision damage over time
- soft-body physics
- flocking behavior
- explicit territory systems
- camera or UI redesign
- broad economic balancing

## Why This Slice Next

The world is now larger, the default population is larger, startup food support is more rule-shaped, and regeneration reacts to depletion. But one major ecosystem pressure is still missing: density itself has little direct cost until a discrete event such as food collection, fight, or reproduction occurs.

That means increased population still risks reading as “more tokens on the board” rather than as a world where local concentration has consequences.

The next pressure is therefore to make crowding itself matter in a minimal way:

- dense areas should become more energetically expensive than sparse areas
- local clustering should create a real survival tradeoff even before explicit collision resolution
- the larger world should gain another route toward collapse, dispersal, and recovery dynamics

This slice is the narrowest next step because it:

- changes one world-pressure variable
- keeps the rest of the rule set intact
- strengthens ecosystem validity without inventing a new entity type or subsystem
- gives future EGD a stronger basis for assessing dominance, collapse, and recovery

## Use-Case Contract

### Use Case

`ApplyLocalCrowdingEnergyPressure`

### Primary Actor

The authoritative server advancing one simulation tick for active circles.

### Pre-conditions

- circles already spend energy for movement
- the world already tracks active player and autonomous circles in shared space
- the larger default world already produces more simultaneous local interactions

### Trigger

A simulation tick advances for active circles in the shared world.

### Success Outcome

- circles in locally crowded neighborhoods lose additional energy through one deterministic rule
- sparse areas remain cheaper than crowded areas
- the world exhibits stronger pressure toward dispersal, depletion, and selective survival

### Failure Or Rejection Cases

- if crowding remains cost-free until explicit collision or food scarcity, the larger world stays mechanically thin
- if the crowding rule becomes global rather than local, it weakens spatial meaning
- if the slice changes fight or reproduction semantics directly, scope is exceeded

## Main Business Rules

1. Crowding pressure remains authoritative server-side behavior.
2. The crowding rule should be based on local nearby-circle presence, not on total global population alone.
3. The rule should add energy pressure through one small deterministic formula or threshold.
4. Player and autonomous circles should be subject to the same crowding rule.
5. Existing movement cost remains in place; this slice only adds a bounded local pressure.
6. Food placement, regeneration, fight, reproduction, continuity, child ownership, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Local Neighborhood`
- `Crowding Pressure`
- `Additional Energy Cost`
- `World Tick`

## Bounded Pressure Interpretation

This slice chooses the smallest useful interpretation of crowding:

- count nearby active circles within one documented radius
- apply a small additional energy cost when local density crosses a simple threshold
- keep the rule uniform across player and autonomous circles

This avoids turning the slice into a full territory or physics system while still making density meaningful.

## Required Runtime Contract Changes

The current contract is likely sufficient if the new pressure is visible through ordinary energy changes and resulting survival outcomes.

Build should avoid extending the runtime contract unless a tiny inspectability field becomes necessary to keep the rule understandable during review.

## Required Ports Or Boundaries

- server-side tick logic for local crowding evaluation
- deterministic tests that prove the additional energy pressure
- implementation notes that record the chosen neighborhood radius and cost rule

## Build Guidance

- prefer adding the rule near existing energy-consumption logic rather than creating a separate subsystem
- choose one simple neighborhood rule
- keep the added energy pressure small but meaningful
- avoid changing steering or food rules in the same slice

## Initial Test Plan

### Server tests

- a circle with no nearby neighbors does not pay the extra crowding cost
- a circle in a locally crowded area pays the extra crowding cost deterministically
- player and autonomous circles are both subject to the same rule

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe crowding-driven energy consequences through ordinary snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. several circles occupy a locally dense region
2. the server applies the crowding pressure during tick advancement
3. circles in that region lose energy faster than isolated circles
4. the resulting world exhibits more meaningful local survival pressure without changing the rest of the rule set

## Done Criteria

- local crowding causes a deterministic additional energy cost
- the rule applies equally to player and autonomous circles
- sparse regions remain cheaper than dense regions
- the rest of the mechanics remain unchanged
- tests cover both crowded and uncrowded cases

## Out Of Scope Follow-Ups

- crowding-aware steering
- territory systems
- continuous collision damage
- UI overlays for density
- fight or reproduction redesign
