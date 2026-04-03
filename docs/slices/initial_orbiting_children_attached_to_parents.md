# Slice: Initial Orbiting Children Attached To Parents

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for attached orbiting child state

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for child ownership, orbit state, and post-reproduction child distribution.

This slice reinterprets children as attached dependents instead of immediately free autonomous circles. It does not yet introduce child detachment, mutation, or a full ancestry system.

## Discovery Scope

Establish the smallest deterministic rule that makes children visible as orbiting dependents of a parent circle:

- each child is represented as an orbiting attached child around one parent circle
- successful reproduction still creates child units, but those child units are distributed between the two parents rather than becoming an immediately independent free circle
- orbiting children stay attached to their current parent
- the current `children_count`, radius growth, and child-based fight leverage may remain in place as transitional semantics for now

This slice does **not** attempt to implement:

- later detachment of children into free autonomous circles
- mutation or inherited traits
- juvenile growth stages
- lineage trees beyond the current lineage identity
- probabilistic simulation that breaks deterministic testing
- full replacement of the current count/radius model in one step

## Why This Slice Next

The current implementation made reproduction visible by spawning one free active circle. That was a useful intermediate step, but it conflicts with the stated product idea that children orbit their parents.

The model pressure is now:

- children should be visibly attached to a parent rather than immediately behaving like ordinary autonomous circles
- children should still matter for power, consumption, and continuity
- the system should preserve inspectability and deterministic tests even though the intended distribution between parents is described as random

This slice is the narrowest next step because it:

- makes the child embodiment closer to the intended game feel
- avoids a premature full rewrite of fight, food, and continuity semantics
- keeps the current child-count-based business rules available while the visible embodiment changes
- makes later decisions about eating, fighting, and consuming children more concrete

## Use-Case Contract

### Use Case

`AttachOrbitingChildrenAfterReproduction`

### Primary Actor

Two different-shape circles reproducing under the current authoritative reproduction rules.

### Pre-conditions

- reproduction succeeds under the current energy and child-reserve rules
- the server owns authoritative world state
- each active circle may host zero or more attached orbiting children
- snapshots can expose attached child state explicitly

### Trigger

The server resolves a successful different-shape reproduction during a simulation tick.

### Success Outcome

- new child units are assigned to one of the two parent circles
- each assigned child appears in later snapshots as an orbiting attached child of that parent
- orbiting children move with the parent through deterministic orbit rules
- the parent remains the main active circle; the child does not become an independent autonomous circle

### Failure Or Rejection Cases

- if children still appear only as hidden counters, the slice fails its purpose
- if children still appear as free autonomous circles immediately after reproduction, the slice stays misaligned with the intended model
- if child distribution is opaque or irreproducible, inspectability is weakened
- if the player gets a special ownership rule different from autonomous circles, fairness is weakened

## Main Business Rules

1. Children are authoritative server-side dependents attached to a parent circle.
2. Orbiting children do not detach in this slice.
3. Successful reproduction creates child units that are distributed between the two participating parents.
4. The intended distribution may feel random, but build must keep it deterministic and inspectable.
5. Orbiting children remain available to the current child-based consumption and leverage rules.
6. Radius growth and child-based fight power may remain active alongside visible orbiting children in this transitional slice.
7. Player and autonomous circles follow the same child ownership and orbit rules.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Child Ownership`
- `Orbit State`
- `Post-Reproduction Child Distribution`
- `Child Count`
- `World Snapshot`

## Bounded Orbit Interpretation

This slice chooses the smallest inspectable interpretation:

- a child belongs to exactly one parent at a time
- a child is rendered as a small visible orbiting body around that parent
- the orbit is deterministic from authoritative state
- child distribution after reproduction is authoritative and reproducible, even if the intended game feel is “random”
- existing child counts remain the source of truth for current power, radius, replacement, and consumption rules until a later slice deliberately removes those shortcuts

This avoids full child autonomy while still making children visible as part of the parent's state.

## Determinism Tension

The intended product idea says children produced by reproduction should be randomly distributed between the parents.

This repository still favors deterministic and inspectable execution. So build should not introduce unbounded runtime randomness. Instead, it should choose one explicit authoritative distribution rule that feels random at the game level but is reproducible in tests, for example by deriving the assignment from the reproduction context, parent IDs, or tick state.

The important point for this slice is not perfect statistical randomness. It is that ownership is not always symmetric and that the resulting distribution is visible and inspectable.

## Required Runtime Contract Changes

The current snapshot contract is no longer sufficient if children are attached orbiters rather than only counts or free circles.

This slice should add explicit child attachment state, for example:

- attached child identifiers
- parent identifier
- orbit position or orbit slot information
- enough data for the client to render visible orbiting children deterministically

The contract should stay small and explicit rather than creating a generic entity graph.

## Required Ports Or Boundaries

- server-side child ownership and orbit derivation
- server-side authoritative post-reproduction child distribution
- client-side rendering of orbiting children around their parent
- deterministic tests covering ownership, distribution, and orbit visibility

## Build Guidance

- treat the current free-spawned child behavior as provisional and replace it with attached orbiting children
- prefer deriving visible orbit positions from small explicit parent-child state rather than storing large mutable motion histories
- keep one explicit authoritative child distribution rule
- preserve the current `children_count`-based business rules for now unless removing one of them is necessary to keep the model coherent
- do not add detachment, mutation, or separate child AI in this slice

## Initial Test Plan

### Server tests

- successful reproduction creates attached orbiting children rather than free autonomous circles
- child ownership after reproduction is deterministic and inspectable
- attached children remain associated with the same parent across later ticks
- child counts still line up with the visible attached children under the chosen transitional rule
- blocked reproduction does not create attached children

### Contract tests

- the snapshot schema exposes enough attached-child state for deterministic rendering
- no generic or open-ended entity graph is required

### Integration tests

- the client receives snapshots showing visible orbiting children around a parent
- after reproduction, child ownership is visible in ordinary browser play
- the demo no longer needs free autonomous spawned children to show reproduction output

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two different-shape circles reproduce successfully
2. the next authoritative snapshot shows child units attached to the participating parents
3. later snapshots show those children orbiting their parent instead of acting as free circles
4. the current parent-level radius and power semantics remain visible for now

## Done Criteria

- successful reproduction now creates visible attached orbiting children instead of immediately free autonomous circles
- child ownership is authoritative, deterministic, and inspectable
- orbiting children remain attached to their parent in later ticks
- the client renders attached orbiting children clearly in the demo
- current child-count-based power, consumption, and continuity rules still work unless intentionally superseded
- tests cover ownership, orbit visibility, and blocked reproduction

## Out Of Scope Follow-Ups

- child detachment into free circles
- direct child-target collision rules separate from parent semantics
- removal of transitional radius growth
- replacement of all child-count shortcuts with fully embodied child-body mechanics
- mutation, inheritance bundles, or ancestry visualization
