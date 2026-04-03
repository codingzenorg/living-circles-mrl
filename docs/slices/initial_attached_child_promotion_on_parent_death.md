# Slice: Initial Attached Child Promotion On Parent Death

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible continuity after parent death

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether one attached child is promoted into the continuing active parent body when the current parent dies.

This slice reuses the current continuity idea but makes it more visibly embodied. It does not introduce child detachment during normal play or a full ancestry system.

## Discovery Scope

Establish the smallest deterministic rule that turns one attached child into the continuing active circle when a parent would otherwise die but still has children:

- when a parent dies and has attached children, one attached child is consumed as the continuation source
- the parent remains active as the continuing lineage carrier
- continuity becomes visibly grounded in attached-child state rather than only in a hidden counter decrement
- the current `lineage_id` and `generation` rules remain in place

This slice does **not** attempt to implement:

- free child detachment during ordinary life
- separate child lineages
- multiple simultaneous promotions
- mutation or inherited traits
- explicit corpse or decay entities
- removal of the current radius shortcut

## Why This Slice Next

The current implementation now makes attached children visible, useful in feeding, and protective in conflict. But when a parent dies and continuity occurs, the transition is still mostly an abstract counter rule rather than a visibly grounded one.

The model pressure is now:

- continuity should clearly come from an attached child that was actually present
- the bridge between visible children and lineage survival should be stronger
- death-and-continuation should feel like a transformation of the current visible state, not just a hidden numeric reset

This slice is the narrowest next step because it:

- preserves the existing replacement semantics while making their source explicit
- stays within the current active-circle model
- does not require full child autonomy or full lifecycle staging
- improves inspectability of lineage continuity

## Use-Case Contract

### Use Case

`PromoteAttachedChildOnParentDeath`

### Primary Actor

Any circle that would otherwise die while still owning at least one attached orbiting child.

### Pre-conditions

- the parent is reaching a current death path such as zero-energy collapse or a defeat path that still removes the parent
- the parent owns at least one attached orbiting child
- snapshots already expose attached-child state

### Trigger

The server resolves a death path for a parent that still owns attached children.

### Success Outcome

- one attached child is consumed as the source of continuity
- the parent remains the active lineage carrier with incremented generation
- later snapshots show one fewer attached child
- the current replacement energy and radius reset rules still apply unless explicitly changed by this slice

### Failure Or Rejection Cases

- if continuity still happens without reference to attached-child state, the slice fails its purpose
- if continuity duplicates the parent instead of consuming one child, inspectability is weakened
- if the player receives a different continuity rule from autonomous circles, fairness is weakened

## Main Business Rules

1. Parent death remains authoritative server-side behavior.
2. If a dying parent has at least one attached child, one attached child is consumed as the continuity source.
3. The active continuing circle keeps the same `lineage_id`.
4. The active continuing circle increments `generation` by one.
5. Player and autonomous circles follow the same promotion rule.
6. Existing replacement energy and radius reset rules may remain unchanged in this slice.
7. This slice does not require free child detachment during ordinary play.

## Minimal Domain Concepts In Scope

- `Parent Death`
- `Attached Child Promotion`
- `Lineage Continuity`
- `Generation Increment`
- `World Snapshot`

## Bounded Promotion Interpretation

This slice chooses the smallest inspectable interpretation:

- when continuity is triggered, consume exactly one currently attached child
- keep the parent body as the active circle representation
- use the attached-child loss as the visible evidence of where continuity came from

This avoids the larger step of converting an orbiting child into a separate promoted body with independent geometry, while still grounding continuity in the visible child model.

## Required Runtime Contract Changes

The current contract is structurally sufficient if promotion is expressed through:

- reduced `attached_children`
- reduced `children_count`
- incremented `generation`
- the existing parent snapshot state

Build should extend the contract only if one explicit continuity outcome marker is necessary for inspectability.

## Required Ports Or Boundaries

- server-side death resolution that explicitly consumes one attached child during continuity
- synchronized count, attached-child, and generation updates
- deterministic tests covering player and autonomous promotion paths
- client rendering that naturally shows the reduced attached-child count and continued active parent

## Build Guidance

- prefer reusing the current death and replacement entry points
- make the attached-child consumption explicit at those continuity points
- keep one-child consumption deterministic and inspectable
- do not introduce free child bodies or corpse entities in this slice
- preserve fairness across player and autonomous circles

## Initial Test Plan

### Server tests

- a zero-energy parent with attached children consumes one visible child and continues
- a defeated parent with attached children consumes one visible child and continues when the current defeat path allows continuity
- `children_count`, `attached_children`, `generation`, and energy remain synchronized after promotion
- a parent with no attached children still dies as before

### Contract tests

- the current snapshot schema remains sufficient unless one explicit continuity outcome is added

### Integration tests

- the client receives a continuity snapshot where one visible attached child is gone and the parent remains active
- the player's continuity in the demo is now visibly grounded in attached-child loss

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent circle reaches a death path while still owning attached children
2. one attached child is consumed as the continuity source
3. the parent remains active with incremented generation
4. the next snapshot shows one fewer orbiting child and continued lineage participation

## Done Criteria

- continuity after parent death explicitly consumes one attached child
- attached-child state, child count, and generation stay synchronized
- player and autonomous circles follow the same promotion rule
- a parent with no attached children still dies as before
- tests cover both continuity and no-child death paths

## Out Of Scope Follow-Ups

- free promotion into a newly positioned child body
- ordinary detachment during life
- corpse entities or decay
- mutation or inherited trait bundles
- removal of transitional radius shortcuts
