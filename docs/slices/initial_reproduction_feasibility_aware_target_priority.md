# Slice: Initial Reproduction Feasibility-Aware Target Priority

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now considers whether preferred reproduction targets are currently feasible

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous target selection.

This slice extends the current deterministic shape-aware interaction-seeking autonomy by making target choice sensitive to whether a preferred different-shape encounter could actually resolve through the current reproduction rules.

## Discovery Scope

Establish the smallest deterministic rule that keeps autonomous target choice aligned with current reproduction feasibility:

- when interaction seeking is active, autonomous circles may prefer different-shape targets only if current reproduction preconditions are satisfiable
- if a preferred different-shape target is currently infeasible, target choice can fall back to another eligible target under a documented rule
- the current movement, energy, contact, fight, and reproduction rules remain unchanged after contact is reached
- the existing food-priority, low-energy, and shape-aware rules remain part of the model

This slice does **not** attempt to implement:

- explicit avoidance of all blocked interactions
- multi-factor utility scoring
- long-horizon planning
- threat evaluation
- personality systems
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now makes shape matter in interaction-target ordering, which better aligns steering with collision meaning. But it still treats every different-shape target as equally attractive even when the current reproduction rule would be blocked by insufficient energy or missing child reserve.

The model pressure is now:

- pre-contact steering should align with current reproduction feasibility, not only with abstract shape preference
- autonomous pursuit should be less semantically wasteful
- energy should matter not only in whether reproduction succeeds, but also in whether reproduction is worth pursuing right now

This slice is the narrowest next step because it:

- changes only target ordering inside interaction-seeking mode
- keeps the existing reproduction and fight resolution rules unchanged
- reuses current energy and child-reserve concepts instead of inventing new AI concepts
- avoids adding avoidance or tactical combat systems

## Use-Case Contract

### Use Case

`PrioritizeInteractionTargetByReproductionFeasibility`

### Primary Actor

Any autonomous circle that is already in interaction-seeking mode.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- food-priority, low-energy, and shape-aware steering rules already exist
- current reproduction feasibility rules already exist
- current interaction target selection already exists

### Trigger

An autonomous circle selects an interaction target after food does not currently win priority.

### Success Outcome

- the autonomous circle may prefer a different-shape target only when the current reproduction rule is feasible for the relevant pair
- if no feasible preferred target exists, the documented fallback rule applies
- later snapshots show target pursuit that better matches actual downstream resolution possibilities

### Failure Or Rejection Cases

- if blocked reproduction targets are still always preferred, steering remains too disconnected from current rules
- if target ordering becomes overly conditional or opaque, inspectability weakens
- if feasibility checks become non-deterministic, reproducibility weakens

## Main Business Rules

1. Autonomous target selection remains authoritative server-side behavior.
2. Shape-aware preference may be narrowed by current reproduction feasibility.
3. The rule must remain deterministic for the same world state and tick.
4. Current food-priority and low-energy rules remain unchanged.
5. Once contact occurs, the current fight and reproduction rules remain unchanged.
6. This slice should use one simple documented feasibility interpretation rather than a large scoring model.
7. Player and autonomous targets remain part of the same candidate set once interaction seeking is active.

## Minimal Domain Concepts In Scope

- `Interaction-Seeking Mode`
- `Reproduction Feasibility`
- `Different-Shape Preference`
- `Fallback Target`
- `Deterministic Target Priority`
- `World Snapshot`

## Bounded Feasibility Interpretation

This slice chooses the smallest inspectable interpretation:

- once food does not currently win, autonomous circles prefer different-shape targets only when the current reproduction rule is feasible
- if no feasible different-shape target exists, they fall back to the nearest remaining eligible target
- no new contact or outcome semantics are introduced

This avoids the larger step of general strategic AI while still making steering better reflect the current executable reproduction model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing energy and child-count values
- existing interaction outcomes once contact occurs

Build should extend the contract only if steering provenance is necessary to explain the feasibility-driven choice.

## Required Ports Or Boundaries

- server-side interaction-target ordering that can evaluate current reproduction feasibility
- deterministic tie-breaking across player and autonomous targets
- tests that show feasible different-shape targets are preferred while infeasible ones are skipped
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current interaction target chooser rather than adding a new steering subsystem
- reuse the current reproduction threshold and child-reserve logic as the feasibility basis
- keep the rule to one clear preference-plus-fallback behavior
- preserve the existing low-energy and nearby-food priority rules
- keep movement speed, energy cost, and all downstream interaction semantics unchanged
- avoid inventing tactical categories or utility scores

## Initial Test Plan

### Server tests

- when a different-shape target is currently feasible for reproduction, it is preferred over a same-shape alternative
- when a different-shape target is currently infeasible, the chooser falls back to the nearest remaining eligible target
- the low-energy food-recovery rule still wins when it should
- tie-breaking remains deterministic across player and autonomous candidates

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle skip an infeasible different-shape target and pursue a valid fallback
- resulting encounters still resolve through the ordinary interaction HUD and world snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle enters interaction-seeking mode
2. more than one eligible target exists, including at least one different-shape target
3. the autonomous circle evaluates whether the preferred different-shape option is currently feasible under the existing reproduction rule
4. it chooses the feasible preferred target or falls back deterministically
5. contact and resolution follow the current fight or reproduction path

## Done Criteria

- reproduction feasibility influences autonomous target ordering
- the rule is deterministic and documented
- existing food-priority, low-energy, and downstream interaction semantics remain unchanged
- tests cover feasible preference and infeasible fallback behavior

## Out Of Scope Follow-Ups

- generalized blocked-interaction avoidance
- tactical combat selection
- personality-based target choice
- detached child autonomy
- removing current radius shortcuts
