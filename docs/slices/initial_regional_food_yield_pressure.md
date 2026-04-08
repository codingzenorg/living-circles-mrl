# Slice: Initial Regional Food Yield Pressure

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots show the consequences of region-sensitive food yield without requiring a new protocol surface by default

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for feeding outcomes.

This slice follows the recent regional ecology work: neighborhoods already recover at different rates and dense regions already cost more energy. The next bounded step is to make food gathered from heavily stripped regions less rewarding than food gathered from healthier regions.

## Discovery Scope

Establish the smallest deterministic rule that ties food energy gain to regional depletion pressure:

- food still exists as fixed slots with fixed identities and positions
- collection still removes exactly one visible food slot
- regeneration timing remains the current deterministic rule
- the new change is only that energy gained from a collected food slot may be reduced when that slot sits in a region with stronger nearby depletion

This slice does **not** attempt to implement:

- variable food size
- terrain or biome systems
- permanent regional fertility maps
- agent memory of good or bad regions
- new reproduction or fight rules
- player-only bonuses or penalties

## Why This Slice Next

The world now supports:

- regional recovery divergence
- local and regional crowding cost
- regional-crowding-aware autonomous steering

That improves medium-term ecological consequence, but regional scarcity still mostly acts through timing and movement cost. Once food does appear, a collected food slot still pays the same reward everywhere.

That leaves a coherence gap:

- stripped regions take longer to recover
- dense regions cost more to inhabit
- but recovered food still has identical value regardless of surrounding depletion

The next pressure is therefore to make regional scarcity affect not only *when* food returns, but also *how much* it helps once collected.

This slice is the narrowest next step because it:

- deepens the regional resource model without adding new entity types
- works through the existing energy loop rather than through new combat or lineage rules
- can produce stronger long-horizon divergence with one deterministic adjustment

## Use-Case Contract

### Use Case

`ResolveFoodCollectionWithRegionalYieldPressure`

### Primary Actor

Any active circle that collects a food slot.

### Pre-conditions

- food collection and energy gain already exist
- regional food regeneration pressure already exists
- nearby missing food slots can already be counted deterministically

### Trigger

A player or autonomous circle collects a food slot.

### Success Outcome

- the circle still consumes exactly one food slot
- the circle still gains energy immediately
- the energy gain is reduced when the collected slot belongs to a more regionally depleted neighborhood under the chosen bounded rule

### Failure Or Rejection Cases

- if the rule is too broad or opaque, feeding becomes hard to reason about
- if the slice changes regeneration timing or slot placement again, scope is exceeded
- if only one actor type is affected, shared-rule fairness weakens

## Main Business Rules

1. Food yield pressure is authoritative server-side behavior.
2. Food yield remains deterministic for the same world state.
3. The pressure should be based on nearby missing food slots around the collected slot, not on a global depletion score.
4. Player and autonomous circles follow the same collection-yield rule.
5. Food slot identity, position, and regeneration timing remain unchanged.
6. Fight, reproduction, continuity, child ownership, and steering rules remain unchanged.

## Minimal Domain Concepts In Scope

- `Food Collection`
- `Energy Gain`
- `Regional Depletion`
- `Deterministic Yield Reduction`

## Bounded Yield Interpretation

This slice chooses the smallest useful interpretation of regional food yield pressure:

- start from the existing deterministic food energy value
- reduce that gain by one bounded regional penalty when nearby food slots are also currently missing
- keep the result deterministic and easy to inspect

This avoids turning food into a continuously variable economy while still making stripped regions less immediately restorative.

## Required Runtime Contract Changes

The current contract is likely sufficient if the changed yield is visible through ordinary energy outcomes.

Build should avoid new contract surface unless a tiny inspectability field becomes necessary during review.

## Required Ports Or Boundaries

- server-side food collection resolution
- deterministic tests that prove different yield outcomes for otherwise similar collections in healthier versus more depleted regions
- implementation notes that record the chosen regional yield rule

## Build Guidance

- prefer reusing the existing regional depletion calculation around food slots
- keep the penalty small and bounded
- do not change slot regeneration timing in the same slice
- do not add random food value variation

## Initial Test Plan

### Server tests

- a collected food slot in a healthy region still yields the ordinary baseline gain
- a collected food slot in a more depleted region yields a smaller deterministic gain
- collection remains deterministic for the same world state

### Contract tests

- the current snapshot schema remains sufficient unless build adds a minimal inspectability field

### Integration tests

- the client can observe the regional yield difference through ordinary snapshot energy changes

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. one circle collects food in a region with low nearby depletion
2. another otherwise similar collection happens in a more depleted region
3. the authoritative world grants the lower-yield recovery in the more depleted region
4. regional scarcity therefore affects both recovery timing and payoff

## Done Criteria

- regional scarcity can reduce food payoff in a bounded deterministic way
- player and autonomous circles follow the same rule
- regeneration timing and slot identity remain unchanged
- tests cover healthy-region and depleted-region collection outcomes

## Out Of Scope Follow-Ups

- biome systems
- food size variation
- food quality memories or agent learning
- lineage-specific feeding advantages
- terrain-driven ecological zones
