# Slice: Initial Regional Crowding Energy Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client renders and inspects that state.

This slice follows the new regional food-regeneration pressure rule. The larger world can now recover unevenly by neighborhood, which is a real step toward medium-term ecological consequence. But the main energy pressure from co-presence is still highly local and instantaneous: a circle either has enough nearby neighbors this tick to pay the crowding surcharge or it does not. That means dense regions are still not yet accumulating broader area-level energy pressure over time.

## Discovery Scope

Establish the smallest useful regional energy consequence:

- preserve the current larger seeded startup world
- preserve the current local crowding move-cost rule
- add deterministic extra energy pressure when a circle remains inside an already crowded regional cluster
- keep the slice focused on energy cost, not new movement heuristics or terrain systems

This slice does **not** attempt to implement:

- terrain or biome generation
- named regions or ownership
- new food mechanics
- new fight, reproduction, or continuity rules
- new client-only legibility work
- fully persistent heatmap or memory systems

## Why This Slice Next

The latest EGD identified the strongest next direction as medium-term regional ecological pressure:

- startup state is now less authored
- food recovery can now diverge by neighborhood
- but area-level energy consequence is still weak

The smallest coherent next step is to make crowded regions cost more than a single local tick surcharge:

- overused regions should become more expensive to inhabit
- circles should pay more for remaining in dense areas than for briefly passing through them
- larger world scale should begin to support visible settlement and displacement patterns

## Use-Case Contract

### Use Case

`ApplyRegionalCrowdingEnergyPressure`

### Primary Actor

The authoritative world simulation during ordinary tick advancement.

### Pre-conditions

- movement already costs energy
- a local crowding surcharge already exists
- larger world scale and regional food recovery now make area-level divergence meaningful

### Trigger

The simulation advances a tick for active circles in crowded parts of the world.

### Success Outcome

- circles still pay the existing local movement and crowding costs
- circles in more persistently crowded regions pay an additional deterministic regional energy cost
- crowded areas become more expensive over time than less pressured areas

### Failure Or Rejection Cases

- if the rule becomes nondeterministic, the slice fails
- if it silently changes movement direction rules, scope is exceeded
- if it changes fight, reproduction, or continuity semantics directly, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all energy costs.
2. The current local crowding surcharge remains in effect.
3. Regional crowding may add extra deterministic energy pressure beyond the local surcharge.
4. The rule must remain deterministic across reset and replay.
5. Existing movement, food, autonomy, fight, reproduction, and continuity semantics remain unchanged outside energy cost.

## Minimal Domain Concepts In Scope

- `Crowded Region`
- `Regional Energy Pressure`
- `Persistent Density Cost`
- `Deterministic Area Consequence`

## Bounded Interpretation

This slice chooses the smallest regionalization of the current crowding model:

- evaluate a circle's neighborhood against a broader area-level density condition
- add a simple extra energy penalty when that broader condition is met
- avoid introducing stored regional heatmaps, ownership, or pathfinding changes

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server-side energy-cost logic
- deterministic server and integration tests for regional-vs-local crowding cost behavior
- implementation notes to record the new regional pressure rule

## Build Guidance

- prefer one simple broader-radius rule over a generalized density system
- keep the additional cost modest and explainable in tests
- preserve existing narrow custom worlds as focused deterministic fixtures
- do not combine this with new visual cues or autonomy changes in the same slice

## Initial Test Plan

### Validation

- `go test ./...`
- `npm run test:contracts`

### Focused checks

- circles in locally crowded but regionally sparse situations keep the current local cost only
- circles in regionally dense situations pay the extra deterministic energy pressure
- the rule stays deterministic across reset
- existing movement and collision semantics remain unchanged

## Scenario Definition

Start the local server and open one browser client on the fullscreen viewport demo.

Scenario steps:

1. circles gather in one part of the larger world
2. that area becomes more energetically expensive than a sparse area
3. circles that remain in the dense region lose energy faster
4. larger-world neighborhoods begin to diverge not only in food recovery but also in inhabitation cost

## Done Criteria

- crowded regions create extra deterministic energy pressure beyond the current local surcharge
- the rule remains deterministic and readable
- larger-world neighborhoods become more differentiated over time
- existing runtime mechanics remain unchanged outside energy-cost calculation

## Out Of Scope Follow-Ups

- persistent regional heatmaps
- terrain systems
- new autonomy strategies tuned to named regions
- new UI-only regional overlays
