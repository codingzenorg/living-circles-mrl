# Slice: Initial Larger Population With Seeded Food Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client renders and inspects that state.

This slice returns to ecosystem scale after the recent viewport-mode presentation work. The current default world is larger than the original baseline, but it still starts from a relatively curated configuration: `1200x900`, `5` autonomous circles, and a food layout taken from a fixed ordered slot list around the center. That remains useful for deterministic testing, but it underplays the broader-world pressure found in the recent lightweight EGD.

## Discovery Scope

Establish the next bounded population-scale baseline:

- expand the default large-world baseline again
- increase default autonomous population in that expanded world
- replace the current hand-placed expanded food layout with seeded pseudo-random slot generation
- keep initialization, reset, and tests deterministic

This slice does **not** attempt to implement:

- runtime random spawning
- stochastic simulation behavior after initialization
- new food mechanics
- new fight, reproduction, or continuity rules
- procedural terrain or biome systems
- minimap or viewport rule changes

## Why This Slice Next

The current large-world baseline solved the first EGD complaint that the world was too small, but it still stays more curated than ecological:

- the default space is still modest for the newer viewport-based presentation
- autonomous population remains small
- food distribution is still visibly authored rather than world-like

The next pressure is to make the default startup world feel larger and less arranged while preserving deterministic authoritative startup.

## Use-Case Contract

### Use Case

`InitializeExpandedPopulationWorld`

### Primary Actor

The developer or player starting or resetting the default demo world.

### Pre-conditions

- the server still owns world initialization and reset
- the default expanded-world path is already separate from narrow custom test setups
- deterministic startup remains required for testability and evaluation

### Trigger

The default expanded demo world is created or reset.

### Success Outcome

- the world starts in a larger bounded space
- more autonomous circles are active by default
- food slots are distributed by a deterministic seeded layout rather than a visibly hand-authored list
- reset reproduces the same authoritative initial state

### Failure Or Rejection Cases

- if initialization becomes nondeterministic between resets, the slice fails
- if custom narrow worlds are forced into the larger baseline, the slice fails
- if the slice changes runtime food behavior after initialization, scope is exceeded

## Main Business Rules

1. The server remains authoritative for initialization and reset.
2. The expanded default world may become larger and more populated.
3. Food placement may become more random-looking only through deterministic seeded generation.
4. Reset must restore the same authoritative initial state for a given configuration.
5. Existing runtime mechanics remain unchanged after initialization.

## Minimal Domain Concepts In Scope

- `Expanded Default World`
- `Initial Population Baseline`
- `Seeded Food Layout`
- `Deterministic Reset`

## Bounded Interpretation

This slice chooses the smallest world-scale shift that matches the user pressure:

- increase expanded-world dimensions
- increase expanded autonomous count
- generate initial food slots from a deterministic seeded layout within safe world bounds
- preserve the narrow custom-world path for focused tests

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server initialization and reset behavior
- server and integration tests for deterministic startup/reset expectations
- implementation notes to record the new baseline and seeding rule

## Build Guidance

- keep pseudo-random food placement deterministic from a fixed seed or fixed generation rule
- preserve safe spacing from world edges and initial entities
- avoid combining this with new food-regeneration or autonomy semantics
- keep the test path readable by allowing focused custom worlds to stay narrow and explicit

## Initial Test Plan

### Validation

- `go test ./...`
- `npm run test:contracts`

### Focused checks

- default expanded world dimensions increase as intended
- default expanded autonomous count increases as intended
- seeded food layout is deterministic across reset
- generated food remains inside bounds and readable in the viewport/minimap
- narrow custom worlds retain their existing smaller deterministic behavior

## Scenario Definition

Start the local server and open one browser client on the fullscreen viewport demo.

Scenario steps:

1. the default world starts with a larger bounded space than the current expanded baseline
2. more autonomous circles are active immediately
3. food appears less hand-authored and more world-like at startup
4. reset reproduces the same startup world again

## Done Criteria

- the default expanded world is larger
- the default expanded population is higher
- initial food layout is seeded and deterministic rather than manually arranged
- narrow custom worlds remain practical for focused tests
- existing runtime mechanics remain unchanged outside initialization/reset

## Out Of Scope Follow-Ups

- runtime random food spawning
- changing regeneration behavior
- procedural terrain
- new autonomy or combat rules
