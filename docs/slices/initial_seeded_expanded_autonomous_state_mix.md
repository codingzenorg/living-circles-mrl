# Slice: Initial Seeded Expanded Autonomous State Mix

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client renders and inspects that state.

This slice follows the seeded expanded autonomous layout. The larger default world now starts with more circles in deterministic seeded positions, and food already starts from a deterministic seeded layout. But the additional expanded autonomous circles still carry a fixed authored state mix: shape, starting energy, and child presence are still hand-picked per ID. That keeps the startup world partially staged even after spatial de-authoring.

## Discovery Scope

Establish the smallest useful follow-up to seeded expanded autonomous placement:

- preserve the larger expanded world baseline
- preserve seeded expanded autonomous positions
- replace the fixed authored state mix for the additional expanded autonomous circles with a deterministic seeded startup mix
- keep the player and the first explicitly-configurable autonomous circles stable

This slice does **not** attempt to implement:

- stochastic runtime behavior
- new fight, reproduction, or continuity rules
- changing the player's startup state
- changing narrow custom-world explicit setups
- procedural factions or biome systems
- new UI-only legibility work

## Why This Slice Next

The startup world is now less authored in space, but not yet in state:

- extra expanded circles still alternate shapes and energies by hard-coded choice
- startup opportunities and risks are still partly pre-scripted
- the world can feel spatially organic while semantically arranged

The next pressure is to make startup state feel more world-like without giving up deterministic reset.

## Use-Case Contract

### Use Case

`InitializeSeededExpandedAutonomousStateMix`

### Primary Actor

The developer or player starting or resetting the default demo world.

### Pre-conditions

- the expanded default world already uses seeded autonomous placement and seeded food layout
- deterministic startup remains required for testability and evaluation
- the player and first explicitly-configurable circles still define the core demo anchor

### Trigger

The default expanded demo world is created or reset.

### Success Outcome

- additional expanded autonomous circles start from a deterministic seeded mix of shape and starting energy
- reset reproduces the same state mix exactly
- the startup world feels less pre-scripted without changing runtime rules

### Failure Or Rejection Cases

- if reset no longer reproduces the same state mix, the slice fails
- if narrow custom worlds lose explicit startup control, the slice fails
- if runtime fight/reproduction/autonomy semantics change beyond startup state, scope is exceeded

## Main Business Rules

1. The server remains authoritative for initialization and reset.
2. Expanded autonomous startup state may become more random-looking only through deterministic seeded generation.
3. The player and first explicitly-configurable circles remain stable startup anchors.
4. Reset must rebuild the same authoritative startup state mix for a given configuration.
5. Existing runtime mechanics remain unchanged after initialization.

## Minimal Domain Concepts In Scope

- `Expanded Autonomous State Mix`
- `Seeded Startup State`
- `Deterministic Reset`
- `Stable Demo Anchor`

## Bounded Interpretation

This slice chooses the smallest coherent startup-state de-authoring step:

- keep the current expanded population size
- keep the current seeded placement model
- derive the additional expanded circles' startup shape and energy from a deterministic seeded rule
- leave player state and narrow custom worlds explicit

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- server initialization and reset behavior
- server and integration tests for deterministic expanded startup/reset expectations
- implementation notes to record the new startup-state seeding rule

## Build Guidance

- use deterministic seeded generation, not uncontrolled randomness
- keep state values within the current semantic range so no other rule needs to change
- preserve readability by avoiding an all-identical or chaotic startup mix
- do not combine this with new runtime behavior rules

## Initial Test Plan

### Validation

- `go test ./...`
- `npm run test:contracts`

### Focused checks

- expanded startup shape/energy mix is deterministic across reset
- the current demo anchor circles remain stable
- narrow custom worlds retain explicit startup state control
- the expanded startup world no longer depends on authored extra-circle state values

## Scenario Definition

Start the local server and open one browser client on the fullscreen viewport demo.

Scenario steps:

1. the default expanded world starts in the larger seeded baseline
2. the additional expanded autonomous circles appear with a deterministic but less authored state mix
3. reset reproduces the same mix again
4. the startup world feels less pre-scripted before runtime interaction begins

## Done Criteria

- additional expanded autonomous circles use a deterministic seeded startup state mix
- reset reproduces that state mix exactly
- demo anchor circles remain stable
- narrow custom worlds remain practical for focused tests
- existing runtime mechanics remain unchanged outside initialization/reset

## Out Of Scope Follow-Ups

- runtime stochastic behavior
- new autonomy heuristics
- changing the player's startup state
- procedural factions or terrain
