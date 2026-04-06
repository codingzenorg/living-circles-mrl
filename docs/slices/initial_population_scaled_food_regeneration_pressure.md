# Slice: Initial Population Scaled Food Regeneration Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where food regeneration remains inspectable while adapting to the larger world baseline

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for resource renewal timing.

This slice introduces one bounded ecological pressure after the larger-world baseline: food regeneration should no longer be one static global delay regardless of world scale and population pressure.

## Discovery Scope

Establish the smallest deterministic resource-pressure rule that reacts to the new larger baseline:

- keep food slots deterministic and fixed in place
- keep regeneration authoritative and deterministic
- let regeneration timing scale from current world pressure rather than remain one unconditional constant
- keep movement, energy, fight, reproduction, continuity, child, and steering rules unchanged

This slice does **not** attempt to implement:

- procedural food spawning
- random regeneration timing
- biome or region logic
- crowding penalties
- population caps
- new autonomy systems
- long-term balancing across many world presets

## Why This Slice Next

The last slice increased world size, population, and starting food capacity. That made the world feel less staged, but it also made the current regeneration rule more obviously provisional.

Right now food always returns after one fixed delay, regardless of:

- how many circles are active
- how many food slots are currently missing
- whether the world is under depletion pressure or recovering

That fixed rule was appropriate while the world was a tiny demo. After scale-up, the next model pressure is to make resource recovery participate in ecosystem behavior instead of acting like a purely static timer.

This slice is the narrowest next step because it:

- changes one ecological pressure without redesigning combat, reproduction, or autonomy
- preserves deterministic testing
- makes depletion and recovery more meaningful in the expanded world
- directly follows from the EGD recommendation to strengthen ecosystem validity

## Use-Case Contract

### Use Case

`RegenerateFoodUnderPopulationPressure`

### Primary Actor

The authoritative server advancing world ticks and deciding when missing food slots become active again.

### Pre-conditions

- the world already tracks deterministic food slots and regeneration
- the expanded default world already contains a larger bounded map, larger population, and larger food set
- food consumption already removes slots temporarily and records when they became missing

### Trigger

A consumed food slot is eligible for eventual regeneration during subsequent simulation ticks.

### Success Outcome

- food still regenerates deterministically into its original slot
- regeneration timing now reflects a small documented pressure rule tied to the larger world baseline
- the larger world can exhibit more meaningful depletion and recovery rhythms than a fixed demo timer

### Failure Or Rejection Cases

- if regeneration stays effectively static regardless of the expanded world, the larger baseline remains mechanically thin
- if regeneration becomes random or hard to inspect, deterministic refinement weakens
- if the slice introduces a large balancing system instead of one bounded rule, scope is exceeded

## Main Business Rules

1. Food slots remain fixed and deterministic.
2. Regeneration remains authoritative server-side behavior.
3. Regeneration timing should respond to current world pressure through one bounded deterministic rule.
4. The rule should remain simple enough to explain in one short paragraph.
5. Movement, energy cost, fight, reproduction, continuity, child ownership, and autonomy remain unchanged.
6. Reset should recreate the same initial world state and the same regeneration behavior from tick zero.

## Minimal Domain Concepts In Scope

- `Food Slot`
- `Regeneration Delay`
- `Resource Pressure`
- `World Tick`

## Bounded Pressure Interpretation

This slice chooses the smallest useful interpretation of ecological pressure:

- regeneration timing can vary deterministically from current world pressure
- the variation is driven by one small rule, such as missing-slot count or active-population pressure
- the slot still returns to its original position with the same identity

This avoids turning the slice into a generalized balancing framework while still making food recovery more ecosystem-like.

## Required Runtime Contract Changes

The current contract is likely sufficient if the new behavior is visible through ordinary food disappearance and return timing.

Build should extend the contract only if a tiny inspectability field is truly needed to keep the rule understandable. Prefer not to add new contract surface by default.

## Required Ports Or Boundaries

- server-side food regeneration logic
- deterministic tests that prove the new delay rule from ordinary world state
- implementation notes that record the chosen pressure rule clearly

## Build Guidance

- prefer evolving the current regeneration function rather than introducing a separate subsystem
- choose one simple pressure signal, not a composite balancing formula
- keep the rule deterministic and inspectable
- avoid changing food placement, food energy, or movement rules in the same slice

## Initial Test Plan

### Server tests

- consumed food still regenerates into the same slot identity and position
- regeneration timing changes deterministically under the chosen pressure rule
- reset restores the same starting state and same regeneration behavior

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe deterministic recovery timing in the expanded world without any new protocol shape unless explicitly required

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. multiple food slots are consumed in the larger world
2. the server applies the new deterministic regeneration-pressure rule
3. depleted regions stay depleted long enough to matter
4. food returns predictably and allows recovery without collapsing into a static timer demo

## Done Criteria

- regeneration remains deterministic
- the rule is stronger than the previous one-size-fits-all delay
- the slice improves depletion/recovery behavior in the expanded world
- the current contract remains unchanged unless a small inspectability addition is clearly justified
- tests cover the chosen pressure rule

## Out Of Scope Follow-Ups

- random spawning
- region-specific ecology
- new autonomy logic
- crowding penalties
- fight or reproduction redesign
