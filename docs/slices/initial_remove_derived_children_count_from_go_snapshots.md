# Slice: Initial Remove Derived Children Count From Go Snapshots

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where attached children remain the only child representation

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for simulation and continuity.

This slice narrows one remaining internal duplication by removing the mirrored `ChildrenCount` field from Go snapshot-facing structs and making child quantity derive from attached children everywhere.

## Discovery Scope

Establish the smallest deterministic change that completes the attached-child single-source-of-truth move:

- Go snapshot-facing circle structs should no longer carry a separate `ChildrenCount` field
- tests and local server-side consumers should derive child quantity from `attached_children`
- current gameplay behavior, contracts, continuity, fight, reproduction, feeding, contact, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new gameplay rules
- new client-side child behaviors
- redesign of continuity or reproduction semantics
- ancestry or event-history systems
- any new contract fields

## Why This Slice Next

Recent slices already made attached children the authoritative child state:

- attached children became the simulation source of truth
- the runtime JSON contract removed `children_count`
- the client already derives readable child quantity from `attached_children`

But the Go-side snapshot and test path still carries a mirrored `ChildrenCount` field as a convenience. That means child state is still represented twice inside the server/test boundary, even though one of those representations is already derived.

The next model pressure is not to change gameplay. It is to complete the representational cleanup so attached children are the only child state everywhere, not only on the wire.

This slice is the narrowest next step because it:

- changes only representation, not behavior
- preserves the existing authoritative child mechanics
- reduces internal drift between the runtime contract and the Go-side snapshot model
- gives build a deterministic cleanup target with straightforward test updates

## Use-Case Contract

### Use Case

`UseAttachedChildrenAsOnlyChildRepresentation`

### Primary Actor

The authoritative server and its deterministic test consumers.

### Pre-conditions

- attached children already are the authoritative child state in simulation
- the runtime contract already omits `children_count`
- Go-side tests still rely on a mirrored derived `ChildrenCount` field

### Trigger

A snapshot is produced or consumed inside Go-side tests or other local server-side code.

### Success Outcome

- Go snapshot-facing structs no longer expose a mirrored `ChildrenCount` field
- all Go-side child quantity reads derive from `attached_children`
- gameplay behavior and snapshot payloads remain unchanged

### Failure Or Rejection Cases

- if build changes gameplay semantics while removing the mirrored field, scope is exceeded
- if child quantity becomes ambiguous for tests or debug consumers, inspectability weakens
- if the JSON runtime contract changes again, this slice is no longer bounded cleanup

## Main Business Rules

1. Attached children remain the sole authoritative child state.
2. Child quantity in Go-side snapshots should be derived from attached children rather than mirrored in a separate field.
3. Current gameplay behavior remains unchanged.
4. The existing runtime contract remains unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `World Snapshot`
- `Derived Child Quantity`

## Bounded Interpretation

This slice chooses the smallest coherent cleanup:

- remove the mirrored Go-side `ChildrenCount` field from snapshot-facing circle structs
- update tests and local readers to use `len(attached_children)`
- leave all other simulation semantics untouched

This avoids turning the slice into a broader refactor or semantic redesign while still finishing the single-source-of-truth move.

## Required Runtime Contract Changes

None expected.

The current runtime contract already expresses child state only through attached children. This slice should keep that contract stable.

## Required Ports Or Boundaries

- Go-side snapshot structs and any helper code that still mirrors child quantity
- deterministic server and integration tests that currently read `ChildrenCount`
- implementation notes that still describe child quantity as duplicated inside Go-side snapshots

## Build Guidance

- prefer removing the mirrored field rather than keeping it and adding more wrappers
- update tests to read child quantity from attached-children length directly
- avoid redesigning snapshot shape or adding replacement child-count helpers back under another name
- keep the scope inside the Go-side snapshot/test boundary

## Initial Test Plan

### Server tests

- existing continuity, fight, and reproduction tests continue to pass when child quantity is derived from `attached_children`
- focused checks use attached-child length directly instead of a mirrored count field

### Contract tests

- the current snapshot schema remains unchanged

### Integration tests

- websocket snapshot tests continue to read child quantity from attached children successfully

## Scenario Definition

Start a local server and run the existing deterministic test suite.

Scenario steps:

1. snapshots are produced for movement, feeding, fight, reproduction, and continuity paths
2. Go-side readers derive child quantity from attached children only
3. all current behavior remains unchanged

## Done Criteria

- Go-side snapshot-facing structs no longer mirror `ChildrenCount`
- tests and local readers derive child quantity from attached children
- gameplay behavior remains unchanged
- contract and integration validation remain green

## Out Of Scope Follow-Ups

- new gameplay rules for children
- continuity or reproduction redesign
- ancestry or historical event logs
- new client-side visual effects
