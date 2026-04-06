# Slice: Initial Seeded Expanded Autonomous Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client renders and inspects that state.

This slice follows the larger-world seeded-food baseline. The default expanded world is now larger, more populated, and food no longer starts from a visibly hand-authored layout. But the autonomous circles still do. Their startup positions remain a small fixed pattern around the center, which makes the world feel only partially de-authored.

## Discovery Scope

Establish the smallest useful follow-up to the seeded food layout:

- preserve the larger expanded world baseline
- preserve the larger default autonomous count
- replace the current hand-authored expanded autonomous placement with deterministic seeded placement
- keep player startup, reset, and focused narrow test setups deterministic

This slice does **not** attempt to implement:

- stochastic runtime movement
- random reproduction or fight rules
- new autonomy semantics
- procedural terrain
- changing food regeneration or startup food seeding again
- changing the player start position

## Why This Slice Next

The last slice made food placement feel less arranged, but the world still starts with an obviously curated autonomous layout:

- autonomous circles still cluster around a few authored offsets
- the larger viewport now exposes more empty authored space
- the world startup still feels half synthetic and half world-like

The next pressure is to make the expanded startup world more uniformly de-authored while preserving deterministic reset behavior.

## Use-Case Contract

### Use Case

`InitializeSeededExpandedAutonomousLayout`

### Primary Actor

The developer or player starting or resetting the default demo world.

### Pre-conditions

- the server still owns world initialization and reset
- expanded startup food already uses deterministic seeded layout
- deterministic startup remains required for testability and evaluation

### Trigger

The default expanded demo world is created or reset.

### Success Outcome

- expanded autonomous circles start from a deterministic seeded layout rather than a visibly hand-authored pattern
- startup reset reproduces the same authoritative autonomous placement
- the player still starts in a stable readable position

### Failure Or Rejection Cases

- if reset no longer reproduces the same startup state, the slice fails
- if narrow custom worlds lose explicit placement control, the slice fails
- if the slice changes runtime movement or interaction semantics, scope is exceeded

## Main Business Rules

1. The server remains authoritative for initialization and reset.
2. Expanded autonomous startup placement may become more random-looking only through deterministic seeded generation.
3. The player start remains explicit and stable.
4. Reset must rebuild the same authoritative startup layout for a given configuration.
5. Existing runtime mechanics remain unchanged after initialization.

## Minimal Domain Concepts In Scope

- `Expanded Autonomous Baseline`
- `Seeded Autonomous Layout`
- `Deterministic Reset`
- `Stable Player Spawn`

## Bounded Interpretation

This slice chooses the smallest coherent de-authoring step:

- keep the existing expanded population size
- generate expanded autonomous startup positions from a deterministic seeded rule
- keep narrow custom worlds and targeted tests explicit rather than generated

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server initialization and reset behavior
- server and integration tests for deterministic expanded startup/reset expectations
- implementation notes to record the new autonomous seeding rule

## Build Guidance

- use deterministic seeded placement, not uncontrolled randomness
- avoid overlaps with the player spawn, world edges, and immediate unreadable clustering
- keep shape assignment and current starting energy/children semantics unchanged
- do not fold additional autonomy behavior changes into this slice

## Initial Test Plan

### Validation

- `go test ./...`
- `npm run test:contracts`

### Focused checks

- expanded autonomous placement is deterministic across reset
- all initial autonomous circles stay inside bounds
- startup spacing remains readable in the viewport and minimap
- narrow custom worlds retain explicit targeted placement

## Scenario Definition

Start the local server and open one browser client on the fullscreen viewport demo.

Scenario steps:

1. the default expanded world starts in the larger seeded-food baseline
2. autonomous circles appear in a deterministic seeded layout rather than a small authored pattern
3. reset reproduces the same startup arrangement again
4. the world feels less arranged before any runtime behavior begins

## Done Criteria

- expanded autonomous startup placement is seeded and deterministic
- reset reproduces that layout exactly
- player spawn remains stable and readable
- narrow custom worlds remain practical for focused tests
- existing runtime mechanics remain unchanged outside initialization/reset

## Out Of Scope Follow-Ups

- runtime autonomous randomness
- changing autonomy heuristics
- procedural terrain
- changing food regeneration
