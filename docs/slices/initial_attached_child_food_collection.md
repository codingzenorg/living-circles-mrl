# Slice: Initial Attached Child Food Collection

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible attached-child collection effects

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether an attached child can collect food on behalf of its parent.

This slice keeps children attached and orbiting. It does not add free child autonomy, direct child hunger, or independent child energy pools.

## Discovery Scope

Establish the smallest deterministic rule that makes attached orbiting children extend a parent's food collection reach:

- an attached child may overlap a food item and collect it for the parent
- the collected energy still goes to the parent circle
- food consumed by an attached child disappears and later regenerates through the existing food-slot rule
- visible orbiting children therefore matter in feeding, not only in reproduction and conflict

This slice does **not** attempt to implement:

- separate child energy values
- detached child foraging
- food prioritization by individual children
- different food values by collector type
- pathfinding or child-target steering
- removal of radius-based collection reach

## Why This Slice Next

The current implementation now gives attached children three roles:

- visible orbiting presence
- consumable reserve for reproduction payment and replacement continuity
- absorbable loss during hostile conflict

But food collection is still resolved only through the parent body radius. That leaves visible children mechanically absent from the feeding loop, even though the source model says children increase area and practical leverage.

This slice is the narrowest next step because it:

- makes orbiting children directly useful in a core ecosystem loop
- reinforces that children extend the parent's effective world presence
- stays compatible with the current transitional radius model
- remains deterministic and cheap to inspect

## Use-Case Contract

### Use Case

`CollectFoodThroughAttachedChild`

### Primary Actor

Any circle with one or more attached orbiting children participating in the shared world.

### Pre-conditions

- the world contains active food
- a parent circle may own one or more attached orbiting children
- snapshots already expose attached-child positions
- food collection currently restores energy to the collecting parent

### Trigger

An attached child overlaps a food item during a simulation tick.

### Success Outcome

- the overlapped food is consumed
- the parent gains the existing deterministic food energy amount
- later snapshots show the consumed food missing until regeneration
- player and autonomous circles follow the same child-based collection rule

### Failure Or Rejection Cases

- if attached-child overlap is visible but does not collect food, the slice fails its purpose
- if child-collected food gives energy to the wrong parent, fairness is weakened
- if child collection bypasses the existing regeneration rule, ecosystem consistency is weakened

## Main Business Rules

1. Food collection remains authoritative server-side behavior.
2. An attached child can collect food for its current parent.
3. Energy gained through attached-child collection is added to the parent, not to a separate child pool.
4. Player and autonomous circles follow the same child-based collection rule.
5. Existing parent-body food collection remains valid in this slice.
6. Food regeneration continues to follow the current deterministic slot rule.
7. This slice does not remove the current radius-growth shortcut.

## Minimal Domain Concepts In Scope

- `Attached Child Collector`
- `Parent Energy Recovery`
- `Food Slot Consumption`
- `World Snapshot`

## Bounded Collection Interpretation

This slice chooses the smallest inspectable interpretation:

- treat attached-child positions as additional authoritative overlap points for food collection
- if either a parent body or one of its attached children overlaps a food item, the food is consumed once
- the parent receives the current fixed energy gain

This keeps the feeding loop simple while making visible children materially useful.

## Required Runtime Contract Changes

The current contract is structurally sufficient if child-based food collection is expressed through:

- existing `foods`
- existing parent energy
- existing attached-child positions in snapshots

Build should extend the contract only if one explicit food-collection event becomes necessary for inspectability.

## Required Ports Or Boundaries

- server-side food-overlap detection that checks attached children as well as parent bodies
- synchronized food removal and parent energy gain
- client-side rendering that already makes attached-child positions visible
- deterministic tests covering player and autonomous child collection

## Build Guidance

- prefer reusing the existing food gain and regeneration rules
- add attached-child overlap checks before or alongside parent-body checks
- ensure one food item is consumed at most once per tick
- do not add independent child energy or AI in this slice
- keep the result inspectable through ordinary snapshots and tests

## Initial Test Plan

### Server tests

- an attached child can collect food that the parent body would not have reached
- collected energy is added to the parent
- a food item consumed by an attached child still follows the current regeneration rule
- autonomous circles can collect through attached children under the same rule

### Contract tests

- the current snapshot schema remains sufficient

### Integration tests

- the client can observe food disappearing when a visible attached child reaches it
- the parent's energy increases after that visible child-based collection

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a circle with an attached child moves near a food slot
2. the attached child reaches the food before the parent body does
3. the food disappears and the parent gains energy
4. the later snapshot shows the missing food and unchanged child attachment

## Done Criteria

- attached children can collect food on behalf of their parent
- parent energy gain and food removal stay deterministic
- player and autonomous circles follow the same rule
- food-slot regeneration still works after child-based collection
- tests cover collection reach and energy gain through attached children

## Out Of Scope Follow-Ups

- detached child foraging
- separate child metabolism
- child-target food preferences
- removal of parent-radius-based collection
- direct child growth or mutation from food
