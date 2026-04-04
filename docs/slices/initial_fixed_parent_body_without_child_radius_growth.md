# Slice: Initial Fixed Parent Body Without Child Radius Growth

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for parent-circle snapshots that no longer use child count to enlarge the visible parent body

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for parent-body size.

This slice narrows the next remaining radius shortcut by making the visible parent body fixed again, while keeping child-based leverage in other already-explicit semantics.

## Discovery Scope

Establish the smallest deterministic change that makes visible parent geometry better match the current embodied child model:

- parent `radius` in snapshots should reflect the visible parent-core body, not child-derived growth
- attached children remain the visible embodiment of accumulated children
- current child-based fight power, replacement continuity, and child-payment semantics remain unchanged
- current food, contact, movement, orbit, reproduction, and steering rules remain unchanged

This slice does **not** attempt to implement:

- removal of `children_count`
- redesign of child-based fight power
- redesign of continuity or reproduction payment
- detached child entities
- new rendering styles beyond the fixed parent-body size

## Why This Slice Next

Recent slices already removed grown derived radius from the main hidden leverage paths:

- food collection no longer uses grown derived radius as silent reach
- circle contact no longer uses grown derived radius as silent encounter reach
- same-shape fight resolution no longer uses derived radius as a tie-break
- movement boundaries no longer use grown derived radius
- attached-child orbit distance no longer uses grown derived radius

That leaves one main remaining shortcut: parent `radius` still grows with child count and still defines the visible body rendered in snapshots and the browser. But the current embodied model already says attached children are the visible externalization of accumulated children. The model pressure is now to stop showing child accumulation twice: once as attached children and again as a larger parent body.

This slice is the narrowest next step because it:

- changes only the visible parent-body size
- preserves the existing attached-child model
- leaves child-based fight power, payment, and continuity intact
- removes a remaining visual shortcut without forcing a larger semantic rewrite

## Use-Case Contract

### Use Case

`RenderParentBodyAsFixedCore`

### Primary Actor

Any player or autonomous parent circle in an authoritative world snapshot.

### Pre-conditions

- attached children are already visible and authoritative
- parent circles still carry child-derived radius from the earlier growth slice
- most embodied mechanics already no longer depend on derived radius

### Trigger

A world snapshot is produced for rendering or downstream inspection.

### Success Outcome

- parent `radius` reflects the fixed visible parent-core body
- child accumulation remains visible through attached children rather than parent-body enlargement
- later snapshots still expose child counts and attached children for the same circles

### Failure Or Rejection Cases

- if parent `radius` still grows with child count, the embodied child model remains visually doubled
- if removing visible parent growth also changes fight power, payment, or continuity semantics, scope is exceeded
- if snapshot interpretation becomes ambiguous, inspectability weakens

## Main Business Rules

1. Parent-body size remains authoritative server-side behavior.
2. Parent `radius` in snapshots uses the fixed visible parent-core size, not child-derived growth.
3. Player and autonomous parent circles follow the same updated visible-body rule.
4. Attached children remain the visible embodiment of accumulated children.
5. Child-based fight power, replacement continuity, and reproduction payment remain unchanged in this slice.
6. Food, contact, movement, orbit distance, reproduction, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Parent Body`
- `Attached Child`
- `Visible Growth`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- the parent core stays visually fixed
- child accumulation is shown through attached children, not enlarged parent radius
- existing non-visual child leverage remains temporarily intact

This avoids a broader redesign of child power or lineage while making the visible geometry more coherent.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through ordinary snapshots:

- existing `radius`
- existing `children_count`
- existing attached-child positions

Build should extend the contract only if fixed parent-body size becomes too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side snapshot/body-size logic that keeps parent radius fixed
- deterministic tests that show child accumulation no longer enlarges the parent body
- client rendering that remains sufficient to inspect attached-child-based visible growth

## Build Guidance

- prefer evolving the current radius derivation and snapshot logic rather than adding a second visual-size field
- keep attached-child counts, positions, and motion unchanged
- preserve current embodied food, contact, movement, and orbit rules
- avoid coupling this slice to combat or reproduction redesign

## Initial Test Plan

### Server tests

- a circle with zero children uses the fixed parent-core radius
- a circle with one or more children still uses the same fixed parent-core radius
- snapshots still show attached children and children counts for the same circles

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client receives snapshots where parent-body size stays fixed while attached children remain visible

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent circle owns one or more attached children
2. the authoritative server produces a world snapshot
3. the parent body remains fixed at the visible core size
4. the attached children remain visible and continue to show accumulated-child state

## Done Criteria

- parent `radius` no longer grows from child count
- player and autonomous parent bodies use the same fixed visible size
- attached children remain visible and deterministic
- child-based fight power, payment, and continuity remain unchanged
- tests cover the removed visible-growth shortcut

## Out Of Scope Follow-Ups

- removing `children_count`
- redesigning child-based fight power to use only embodied child geometry
- redesigning replacement continuity
- changing attached-child rendering style
- detached child entities
