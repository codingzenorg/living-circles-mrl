# Slice: Initial Remove Dead Derived Radius State

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where parent radius remains the fixed visible body size

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for simulation state.

This slice narrows one remaining internal shortcut by removing dead derived-radius state that no longer changes behavior.

## Discovery Scope

Establish the smallest deterministic cleanup that matches the already-implemented embodied model:

- parent radius should remain the fixed visible body size
- server state should stop pretending radius is still derived from child count
- dead helpers and constants related to child-driven parent radius growth should be removed or collapsed
- gameplay behavior, contracts, continuity, fight, reproduction, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new growth behavior
- new fight or reproduction rules
- new child semantics
- new client visuals
- any runtime contract change

## Why This Slice Next

Recent slices already removed child-derived parent size from the actual behavior:

- food collection no longer uses grown parent radius
- contact initiation no longer uses grown parent radius
- movement boundaries no longer use grown parent radius
- orbit distance no longer uses grown parent radius
- parent snapshots and rendering now keep a fixed visible body size

But the server still carries dead derived-radius structure:

- `derivedRadius(...)` still exists even though it now always returns the fixed size
- `DefaultChildRadiusGain` still exists even though it no longer affects behavior
- child-sync helpers still read as if radius is recalculated from child count

The next model pressure is not a new game rule. It is to remove this stale internal model so the code says what the current simulation actually does.

This slice is the narrowest next step because it:

- changes only internal representation and naming
- preserves the current embodied behavior exactly
- reduces confusion for future refinement
- gives build a deterministic cleanup target with straightforward regression coverage

## Use-Case Contract

### Use Case

`KeepParentRadiusAsFixedVisibleBody`

### Primary Actor

The authoritative server and its local implementation artifacts.

### Pre-conditions

- parent radius is already behaviorally fixed at the visible core size
- derived-radius helpers no longer change any runtime outcome
- attached children already carry the visible embodiment of accumulated children

### Trigger

A world is initialized or child state changes during reproduction, absorption, payment, or continuity.

### Success Outcome

- parent radius stays fixed at the visible core size
- no dead derived-radius helper or child-radius-growth constant remains in the active simulation path
- gameplay behavior and snapshot payloads remain unchanged

### Failure Or Rejection Cases

- if build changes gameplay while removing dead state, scope is exceeded
- if parent radius becomes ambiguous again, embodied geometry weakens
- if child growth is reintroduced implicitly through a renamed helper, the cleanup is incomplete

## Main Business Rules

1. Parent radius remains the fixed visible body size.
2. Attached children remain the visible embodiment of accumulated children.
3. Removing dead derived-radius state must not change gameplay behavior.
4. The runtime contract remains unchanged.

## Minimal Domain Concepts In Scope

- `Parent Radius`
- `Attached Child`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest coherent cleanup:

- remove or collapse dead derived-radius helpers and constants
- make initialization and child-sync logic say explicitly that parent radius is fixed
- leave all gameplay rules unchanged

This avoids turning the slice into a new growth design while making the current embodied model honest in code.

## Required Runtime Contract Changes

None expected.

Parent radius is already fixed in the runtime contract, so this slice should keep the wire shape stable.

## Required Ports Or Boundaries

- Go simulation state and helper functions that still reference derived radius
- deterministic tests that currently imply fixed parent radius
- implementation notes that still describe child-derived radius growth as active behavior

## Build Guidance

- prefer removing dead helpers rather than keeping no-op abstractions
- update tests only where they still describe child-driven parent radius as live behavior
- keep the cleanup inside the server/state boundary
- do not reintroduce a shadow “future growth” abstraction during this slice

## Initial Test Plan

### Server tests

- existing fixed-parent-radius tests continue to pass
- child-related state changes do not alter parent radius

### Contract tests

- the current snapshot schema remains unchanged

### Integration tests

- websocket snapshot tests continue to observe fixed parent radius

## Scenario Definition

Start a local server and run the existing deterministic test suite.

Scenario steps:

1. worlds initialize with player and autonomous parent circles
2. child-affecting events occur through reproduction, absorption, payment, and continuity
3. parent radius remains fixed and no derived-radius state is needed to explain the behavior

## Done Criteria

- dead derived-radius helpers or constants are removed or collapsed
- parent radius remains fixed in all current behavior
- tests and implementation notes align with the fixed-radius model
- contract and integration validation remain green

## Out Of Scope Follow-Ups

- reintroducing area growth through a new embodied mechanism
- new fight or reproduction semantics
- new client visual effects
- redesigning child power or continuity
