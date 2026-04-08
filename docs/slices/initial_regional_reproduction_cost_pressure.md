# Slice: Initial Regional Reproduction Cost Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots show the consequences of region-sensitive reproduction cost without requiring a new protocol surface by default

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for reproduction resolution.

This slice follows the recent regional ecology work: neighborhoods now differ in regeneration timing, crowding cost, steering consequence, and food yield. The next bounded step is to make population growth itself feel more expensive in stripped regions.

## Discovery Scope

Establish the smallest deterministic rule that ties successful reproduction cost to regional depletion pressure:

- reproduction still requires the current threshold and still uses the existing cost model as a baseline
- ownership, child payment, redistribution, and continuity rules remain unchanged
- the new change is only that successful reproduction may cost more in more locally depleted regions

This slice does **not** attempt to implement:

- fertility maps
- seasonal or biome-dependent reproduction rules
- mutation or trait inheritance
- region-specific fight rules
- new child ownership semantics
- player-only reproduction penalties

## Why This Slice Next

Regional scarcity now affects:

- when food returns
- how much a collected food helps
- how expensive dense regions are to inhabit
- how autonomous circles steer around regional pressure

But successful reproduction still costs the same everywhere once the current threshold is met. That leaves a remaining ecological gap:

- stripped regions are harsher to survive in
- but they are not yet harsher places to expand population

The next pressure is therefore to make regional scarcity influence not only survival and movement, but also the cost of producing new continuity units.

This slice is the narrowest next step because it:

- deepens the regional resource model without creating a new subsystem
- acts on an already-authoritative reproduction path
- can create stronger long-horizon population differentiation with one bounded deterministic rule

## Use-Case Contract

### Use Case

`ResolveReproductionWithRegionalCostPressure`

### Primary Actor

Any same-tick pair that successfully resolves different-shape reproduction.

### Pre-conditions

- reproduction feasibility, child payment, and distribution already exist
- regional depletion can already be counted deterministically around food slots
- reproduction resolution already has a deterministic cost step

### Trigger

A different-shape pair meets the current reproduction conditions and reproduction resolves successfully.

### Success Outcome

- reproduction still succeeds under the current threshold rule
- child creation and ownership still follow the existing rule
- the successful reproduction pays a higher deterministic cost when it occurs in a more regionally depleted neighborhood

### Failure Or Rejection Cases

- if the rule changes feasibility and payment semantics at the same time, scope is exceeded
- if the regional pressure becomes opaque, reproduction becomes harder to reason about
- if only one actor type is affected, shared-rule fairness weakens

## Main Business Rules

1. Regional reproduction cost pressure is authoritative server-side behavior.
2. The rule should be based on nearby missing food slots around the reproduction location, not on a global depletion score.
3. The cost increase should remain bounded and deterministic.
4. Player and autonomous circles follow the same regional reproduction-cost rule.
5. Reproduction threshold, child payment semantics, ownership, distribution, and continuity remain unchanged.
6. Food placement, regeneration timing, fight, and steering rules remain unchanged.

## Minimal Domain Concepts In Scope

- `Reproduction`
- `Reproduction Cost`
- `Regional Depletion`
- `Deterministic Cost Increase`

## Bounded Cost Interpretation

This slice chooses the smallest useful interpretation of regional reproduction pressure:

- start from the current deterministic reproduction cost
- increase that cost by one bounded regional penalty when reproduction occurs in a more depleted neighborhood
- keep the result deterministic and easy to inspect

This avoids rewriting the reproduction model while still making regional scarcity affect population growth directly.

## Required Runtime Contract Changes

The current contract is likely sufficient if the changed cost is visible through ordinary energy outcomes.

Build should avoid new contract surface unless a tiny inspectability field becomes necessary during review.

## Required Ports Or Boundaries

- server-side reproduction resolution
- deterministic tests that prove different successful reproduction costs in healthier versus more depleted regions
- implementation notes that record the chosen regional cost rule

## Build Guidance

- prefer reusing the same regional depletion neighborhood already used by regional food pressure
- keep the added cost small and bounded
- do not change the reproduction threshold or payment mechanism in the same slice
- avoid adding random fertility behavior

## Initial Test Plan

### Server tests

- successful reproduction in a healthier region still pays the ordinary baseline cost
- successful reproduction in a more depleted region pays a higher deterministic cost
- ownership and child-payment behavior remain unchanged

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe the regional reproduction-cost difference through ordinary snapshot energy changes

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. one different-shape pair reproduces in a healthier region
2. another otherwise similar reproduction resolves in a more depleted region
3. the authoritative world charges the higher reproduction cost in the more depleted region
4. regional scarcity therefore affects both survival and population expansion

## Done Criteria

- regional scarcity can increase successful reproduction cost in a bounded deterministic way
- player and autonomous circles follow the same rule
- reproduction threshold, payment semantics, and ownership remain unchanged
- tests cover healthier-region and depleted-region reproduction outcomes

## Out Of Scope Follow-Ups

- fertility zones
- seasonal reproduction
- mutation systems
- region-specific child distribution rules
- lineage-specific reproduction advantages
