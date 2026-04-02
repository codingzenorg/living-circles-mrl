# Slice: Initial Child Growth To Radius

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for child accumulation and current circle size

## Architecture Mode

Explicit client/server boundary with the server deriving visible circle growth from accumulated children.

This slice extends the current reproduction accumulation model without yet adding spawned child entities, lineage replacement, or a richer combat framework.

## Discovery Scope

Establish the smallest deterministic growth rule that gives child accumulation visible leverage:

- a circle's accumulated children affect its radius
- the increased radius is reflected in snapshots and rendering
- larger radius improves food collection reach
- larger radius can already matter in fights through the existing radius tie-breaker

This slice does **not** attempt to implement:

- separate child entities moving in the world
- new fight formulas beyond the current deterministic rule
- explicit attack or defense stats
- replacement of dead circles by children
- lineage trees, inheritance, or mutation
- shrinking or spending children

## Why This Slice Next

The source material says children increase area, improve food collection, and increase power in fights. The current implementation now records `children_count`, but that state is still inert.

Turning children into radius growth is the narrowest next step because it:

- gives reproduction a visible consequence
- stays faithful to "increases area"
- improves food collection without inventing new collection rules
- influences existing fights without rewriting the fight system
- keeps continuity and replacement semantics deferred

## Use-Case Contract

### Use Case

`ApplyChildGrowthLeverage`

### Primary Actor

Player controlling one circle while autonomous circles participate under the same growth rules.

### Pre-conditions

- a server process can host one bounded world
- a circle may accumulate one or more children through resolved reproduction
- the world snapshot already carries radius and child count

### Trigger

The server advances a tick after a circle has accumulated children.

### Success Outcome

- the circle's radius increases deterministically from its child count
- later snapshots expose the increased radius
- food overlap and fight tie-breaks reflect the larger radius automatically
- player and autonomous circles follow the same growth rule

### Failure Or Rejection Cases

- if a circle has zero children, it should keep the base radius
- if growth requires hidden mutable state separate from child count, this slice should reject that design and derive radius directly
- if growth makes circles exceed world bounds, position clamping should still keep them inside the world

## Main Business Rules

1. Child accumulation must have a visible physical consequence in this slice.
2. Radius is derived deterministically from `children_count`.
3. The same growth rule applies to player and autonomous circles.
4. Growth improves leverage indirectly through existing overlap logic rather than through a new combat subsystem.
5. This slice does not consume, transfer, or reset children.
6. This slice does not yet define whether child growth is linear, multiplicative, or capped in the long term; it only needs one explicit deterministic starting rule.

## Minimal Domain Concepts In Scope

- `Child Count`
- `Radius Growth`
- `Food Collection Reach`
- `Fight Leverage`
- `World Snapshot`

## Bounded Growth Interpretation

The source material says children increase area and fight power, but it does not specify the exact formula.

This slice therefore chooses the smallest inspectable interpretation:

- each circle has a base radius
- each accumulated child increases radius by a fixed deterministic amount
- no other stat changes are introduced

That keeps the effect legible and makes later refinement possible if nonlinear growth or caps become necessary.

## Required Runtime Contract Changes

The snapshot contract already carries `radius` and `children_count`.

This slice may not require a new field, but it must ensure:

- `radius` now reflects child-based growth rather than staying static
- the client renders the larger radius directly from authoritative snapshots
- tests verify the relationship between child count and radius

## Required Ports Or Boundaries

- server-side deterministic radius derivation from child count
- server-side collision and food collection using the grown radius
- client-side rendering that naturally reflects larger circles
- deterministic tests covering radius growth and resulting leverage

## Build Guidance

- prefer deriving radius from `children_count` rather than storing a second mutable growth counter
- use one explicit fixed growth increment per child
- avoid rewriting food or fight systems when the existing overlap and tie-break behavior can express the leverage
- keep the effect inspectable in snapshots and tests
- do not add caps, decay, or balancing rules yet unless required for determinism

## Initial Test Plan

### Server tests

- a circle with zero children keeps the base radius
- a circle with one or more children has a larger radius
- child-based radius growth lets a circle reach food sooner or from farther away
- child-based radius growth can break an otherwise equal-energy fight through the existing radius tie-breaker
- player and autonomous circles follow the same growth rule

### Contract tests

- the snapshot schema remains sufficient with `children_count` plus `radius`
- no extra contract field is required for the first growth interpretation

### Integration tests

- after reproduction, later snapshots show increased radius for the participating circles
- the browser client receives and renders the larger circles without local inference

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot with base-sized circles
2. the player triggers different-shape reproduction
3. a later snapshot shows increased child counts
4. a later snapshot also shows increased radius for the circles that accumulated children
5. the larger radius becomes visible in the canvas and meaningful in subsequent food or fight interactions

## Done Criteria

- child accumulation causes deterministic radius growth
- grown radius is visible in snapshots and the browser client
- food collection and fight leverage can reflect the grown radius without separate special-case rules
- tests cover radius derivation and at least one leverage consequence
- the slice does not implement spawned children, continuity replacement, mutation, or a new combat framework

## Out Of Scope Follow-Ups

- caps or decay for child-based growth
- explicit combat power formulas tied to children
- child spending mechanics
- replacement of dead circles by children
- separate child entities
- lineage inheritance behavior
