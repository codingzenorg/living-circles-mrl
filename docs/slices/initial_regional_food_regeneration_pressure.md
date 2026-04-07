# Slice: Initial Regional Food Regeneration Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client renders and inspects that state.

This slice follows the recent startup-world de-authoring work. The default world now begins from a larger, less staged state, but the main EGD gap is no longer startup plausibility. It is medium-term ecological consequence. Food already regenerates more slowly under global missing-slot pressure, yet that pressure is still world-wide and uniform. A depleted local area does not become meaningfully distinct from a healthy one if the total missing-slot count happens to match.

## Discovery Scope

Establish the smallest useful regional ecological pressure:

- preserve the current larger seeded startup world
- preserve the existing food-slot regeneration model
- add deterministic local regeneration slowdown based on nearby missing food slots
- keep the slice focused on food recovery pressure, not broader biome or terrain systems

This slice does **not** attempt to implement:

- random runtime food spawning
- terrain or biome generation
- new autonomy rules
- new fight, reproduction, or continuity rules
- new client-only legibility systems
- region ownership or faction systems

## Why This Slice Next

The latest EGD identified the main remaining gap as medium-term ecological consequence:

- the world starts plausibly
- startup state is less authored
- but different areas of the world still do not diverge strongly enough over time

The smallest coherent next step is to make resource recovery more local:

- a region that has been stripped of food should recover more slowly than a region that has not
- local depletion should matter even when global missing-slot count is similar
- larger world scale should begin to create meaningfully different neighborhoods

## Use-Case Contract

### Use Case

`RegenerateFoodUnderRegionalPressure`

### Primary Actor

The authoritative world simulation during ordinary tick advancement.

### Pre-conditions

- food remains modeled as deterministic slots
- consumed food still returns to its original slot
- global missing-slot tracking already exists
- larger default worlds now make regional divergence meaningful

### Trigger

The simulation evaluates whether missing food slots should regenerate on a tick.

### Success Outcome

- food still regenerates deterministically into the same slot
- a slot in a locally depleted area takes longer to return than a slot in a less depleted area
- different regions of the world can recover at different rates without changing other rules

### Failure Or Rejection Cases

- if regeneration becomes nondeterministic, the slice fails
- if slots start moving or respawning elsewhere, scope is exceeded
- if the rule changes food consumption, autonomy, or combat semantics directly, scope is exceeded

## Main Business Rules

1. The server remains authoritative for food regeneration timing.
2. Food still regenerates only into its original slot with the same ID.
3. Regeneration delay may increase under local missing-slot pressure, not just global pressure.
4. The rule must remain deterministic across reset and replay.
5. Existing movement, feeding, autonomy, fight, reproduction, and continuity semantics remain unchanged.

## Minimal Domain Concepts In Scope

- `Food Slot`
- `Regional Pressure`
- `Local Depletion`
- `Deterministic Recovery Delay`

## Bounded Interpretation

This slice chooses the smallest regionalization of the current food model:

- define local neighborhood pressure around a food slot
- let nearby missing slots extend that slot's regeneration delay
- preserve current slot identity and position
- avoid introducing terrain, ownership, or adaptive spawning

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server-side regeneration logic
- deterministic server and integration tests for local-vs-global recovery behavior
- implementation notes to record the local regeneration rule

## Build Guidance

- prefer one simple local-radius rule over a generalized regional subsystem
- keep the rule readable enough that tests can prove why one slot regenerates sooner than another
- preserve narrow custom worlds as practical deterministic test fixtures
- do not combine this with new visual cues or autonomy changes in the same slice

## Initial Test Plan

### Validation

- `go test ./...`
- `npm run test:contracts`

### Focused checks

- slots in less depleted neighborhoods regenerate sooner than slots in more depleted neighborhoods
- regeneration stays deterministic across reset
- slot identity and position remain unchanged
- existing narrow-world regeneration tests remain understandable

## Scenario Definition

Start the local server and open one browser client on the fullscreen viewport demo.

Scenario steps:

1. food is consumed in multiple parts of the larger world
2. one local area becomes more stripped than another
3. food in the heavily depleted area returns more slowly
4. different neighborhoods begin to diverge in resource recovery pace

## Done Criteria

- food regeneration becomes region-sensitive rather than only globally pressure-sensitive
- the rule remains deterministic and slot-based
- larger-world neighborhoods can recover at different rates
- existing runtime mechanics remain unchanged outside regeneration timing

## Out Of Scope Follow-Ups

- dynamic terrain
- food relocation
- runtime random spawning
- new autonomy strategies based on named regions
