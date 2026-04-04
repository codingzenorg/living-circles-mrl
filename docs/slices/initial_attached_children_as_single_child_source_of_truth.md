# Slice: Initial Attached Children As Single Child Source Of Truth

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract that may still expose `children_count`, but no longer relies on a separately authoritative stored child-count state

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for child state.

This slice narrows the remaining child-state duplication by making attached children the single authoritative source for child-dependent rules, while keeping the current snapshot shape readable.

## Discovery Scope

Establish the smallest deterministic change that makes child semantics less redundant:

- attached children become the sole authoritative child state in the simulation
- child-dependent rules derive their count from attached children rather than a separate stored `children_count` authority
- snapshots may still expose `children_count` as a derived convenience field
- current fight power, reproduction payment, continuity, feeding, contact, movement, orbit, and steering behaviors remain unchanged

This slice does **not** attempt to implement:

- removal of `children_count` from the snapshot contract
- redesign of child-based fight power
- redesign of reproduction payment semantics
- redesign of continuity semantics
- detached child entities

## Why This Slice Next

The current embodied model has advanced far enough that attached children are now:

- visible in snapshots and rendering
- consumed for payment, absorption, and continuity
- used in feeding, contact, avoidance, and targeting

But the simulation still carries a separate stored `children_count` authority and many rules still read that duplicated value directly. That means the model continues to split “what a child is” between visible child entities and mirrored bookkeeping state.

The next model pressure is to make attached children the real authoritative child state and treat `children_count` as a derived view rather than as a second source of truth.

This slice is the narrowest next step because it:

- changes state authority rather than visible behavior
- reduces drift risk between visible children and rule inputs
- preserves the current contract shape and semantics
- prepares later slices to narrow remaining child-count shortcuts without carrying duplicated state

## Use-Case Contract

### Use Case

`DeriveChildDependentRulesFromAttachedChildren`

### Primary Actor

Any player or autonomous parent circle that owns zero or more attached children.

### Pre-conditions

- attached children are already visible and authoritative enough to be consumed and positioned
- the simulation still stores `children_count` alongside attached children
- child-dependent rules already expect deterministic child-related behavior

### Trigger

A child-dependent rule is evaluated or a snapshot is produced.

### Success Outcome

- child-dependent rule evaluation derives child quantity from attached children
- snapshots may still expose `children_count`, but it is only a derived convenience field
- visible behavior remains unchanged for fights, reproduction, continuity, and other child-dependent outcomes

### Failure Or Rejection Cases

- if `children_count` remains a separate authority, visible child state and rule inputs can still drift apart
- if this slice changes user-visible child semantics, scope is exceeded
- if the contract stops exposing readable child counts without deliberate decision, inspectability weakens

## Main Business Rules

1. Attached children are the authoritative child state for each parent.
2. Child-dependent rules derive child quantity from attached children.
3. `children_count` may remain in snapshots as a derived value.
4. Player and autonomous circles follow the same child-state rule.
5. Fight power, reproduction payment, continuity, feeding, contact, movement, orbit, and steering remain behaviorally unchanged in this slice.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Derived Child Count`
- `Single Source Of Truth`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- attached children are authoritative
- `children_count` becomes derived
- the visible and behavioral model should stay the same

This avoids a larger semantics rewrite while removing one remaining duplication point from the model.

## Required Runtime Contract Changes

The current contract is likely sufficient if:

- `children_count` remains present as a derived convenience field
- attached-child arrays remain present and authoritative

Build should change the contract only if deriving the count internally creates ambiguity in snapshots.

## Required Ports Or Boundaries

- server-side child-dependent rules that derive child quantity from attached children
- deterministic tests that prove child counts remain synchronized because they are derived, not separately authoritative
- client rendering that can continue to use the existing snapshot shape unchanged

## Build Guidance

- prefer removing direct dependence on stored child-count authority rather than adding more synchronization code
- keep snapshot readability by continuing to expose `children_count` if useful
- preserve existing behaviors and outcomes
- avoid turning this slice into a broader redesign of child mechanics

## Initial Test Plan

### Server tests

- child-dependent rules still behave the same when attached children exist
- derived `children_count` matches attached-child length in snapshots
- no behavior changes occur in payment, absorption, or continuity scenarios

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client still receives stable `children_count` plus attached-child arrays with no visible behavior change

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent owns attached children
2. the authoritative server evaluates a child-dependent rule
3. the rule derives child quantity from attached children
4. later snapshots still expose the same readable child count and visible child bodies

## Done Criteria

- attached children are the single authoritative child state
- `children_count` is derived rather than separately authoritative
- current visible and gameplay behaviors remain unchanged
- tests cover derivation and unchanged child-dependent outcomes

## Out Of Scope Follow-Ups

- removing `children_count` from the contract
- redesigning child-based fight power
- redesigning reproduction payment
- redesigning replacement continuity
- detached child entities
