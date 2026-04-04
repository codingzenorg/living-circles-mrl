# Slice: Initial Fight Feasibility-Aware Target Priority

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now considers whether fallback same-shape pursuit is currently survivable

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous target selection.

This slice extends the current deterministic interaction-seeking autonomy by making same-shape fallback pursuit sensitive to the existing fight ordering rather than treating every remaining target as equally acceptable.

## Discovery Scope

Establish the smallest deterministic rule that keeps interaction-seeking autonomy aligned with current same-shape fight feasibility:

- once no preferred feasible reproduction target wins, autonomous circles may prefer same-shape targets only when the current fight rule says they would not lose
- if no non-losing same-shape target exists, the documented fallback rule applies
- the current movement, energy, contact, fight, and reproduction rules remain unchanged after contact is reached
- the existing food-priority, low-energy, shape-aware, and reproduction-feasibility-aware rules remain part of the model

This slice does **not** attempt to implement:

- generalized escape or evasion movement
- long-horizon tactical planning
- fear, aggression, or personality modes
- opponent prediction beyond the current deterministic fight ordering
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now makes autonomous steering sensitive to both energy condition and reproduction feasibility. But once feasible reproduction is unavailable, the chooser still treats same-shape fight targets as acceptable even when the current deterministic fight rule says the autonomous circle would lose.

The model pressure is now:

- pre-contact steering should align with the current fight semantics, not only with shape and reproduction semantics
- autonomous behavior should stop walking into clearly losing same-shape outcomes as a default fallback
- children and energy should matter before contact as well as during fight resolution

This slice is the narrowest next step because it:

- changes only fallback target ordering inside the current interaction-seeking mode
- keeps the current fight and reproduction resolution rules unchanged
- reuses the existing deterministic fight ordering rather than inventing a new AI model
- avoids introducing explicit fleeing or tactical combat systems

## Use-Case Contract

### Use Case

`PrioritizeInteractionTargetByFightFeasibility`

### Primary Actor

Any autonomous circle that is already in interaction-seeking mode.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- food-priority, low-energy, shape-aware, and reproduction-feasibility-aware steering rules already exist
- current same-shape fight ordering already exists
- current interaction target selection already exists

### Trigger

An autonomous circle selects an interaction target after food does not currently win priority.

### Success Outcome

- the autonomous circle may treat a same-shape target as a fallback pursuit only when current fight ordering says it would not lose
- if no non-losing same-shape fallback exists, the documented non-social fallback rule applies
- later snapshots show pursuit that is more consistent with the current executable fight semantics

### Failure Or Rejection Cases

- if clearly losing same-shape targets are still pursued by default, steering remains too disconnected from current fight rules
- if target ordering becomes opaque or multi-factor in an uninspectable way, reproducibility weakens
- if fallback behavior becomes non-deterministic, tests and explanation become harder

## Main Business Rules

1. Autonomous target selection remains authoritative server-side behavior.
2. Different-shape preference remains governed by the existing reproduction-feasibility-aware rule.
3. Same-shape fallback pursuit may be narrowed by the current deterministic fight ordering.
4. The rule must remain deterministic for the same world state and tick.
5. Current food-priority and low-energy rules remain unchanged.
6. Once contact occurs, the current fight and reproduction rules remain unchanged.
7. This slice should use one simple documented fight-feasibility interpretation rather than a large tactical scoring model.
8. Player and autonomous targets remain part of the same candidate set once interaction seeking is active.

## Minimal Domain Concepts In Scope

- `Interaction-Seeking Mode`
- `Fight Feasibility`
- `Same-Shape Fallback`
- `Non-Losing Target`
- `Deterministic Target Priority`
- `World Snapshot`

## Bounded Fight-Feasibility Interpretation

This slice chooses the smallest inspectable interpretation:

- after feasible different-shape targets are considered, same-shape targets are acceptable only when the current fight ordering says the autonomous circle would win or reach the existing exact-tie outcome that does not count as an immediate loss
- if no such same-shape target exists, the autonomous circle falls back to the nearest available food target, and only to baseline drift when no food target exists
- no new contact or outcome semantics are introduced

This avoids the larger step of explicit threat evasion while still making steering better reflect the current executable conflict model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing energy, child-count, and shape values
- existing interaction outcomes once contact occurs

Build should extend the contract only if steering provenance is necessary to explain why a losing same-shape target was skipped.

## Required Ports Or Boundaries

- server-side interaction-target ordering that can evaluate the current deterministic fight ordering before fallback pursuit
- deterministic tie-breaking across player and autonomous targets
- tests that show non-losing same-shape targets are pursued while clearly losing ones are skipped
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current interaction target chooser rather than adding a new steering subsystem
- reuse the current fight ordering as the feasibility basis
- keep the rule to one clear eligibility-plus-fallback behavior
- preserve the existing low-energy and nearby-food priority rules
- keep movement speed, energy cost, and all downstream interaction semantics unchanged
- avoid inventing explicit fear states or escape maneuvers

## Initial Test Plan

### Server tests

- when no feasible reproduction target exists, a same-shape target that would currently lose is skipped
- when no feasible reproduction target exists, a same-shape target that would currently win is pursued
- when all social targets are currently losing or blocked, the chooser falls back to food or baseline drift under a documented rule
- tie-breaking remains deterministic across player and autonomous candidates with equal current fight feasibility

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle skip a stronger same-shape target and continue toward the documented fallback
- the client also receives snapshots showing an autonomous circle pursue a weaker same-shape target when that fallback is currently survivable

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle enters interaction-seeking mode
2. no feasible different-shape reproduction target currently wins
3. at least one same-shape target exists
4. the autonomous circle evaluates whether that same-shape target would currently be survivable under the existing fight ordering
5. it pursues the non-losing target or falls back deterministically
6. any resulting contact resolves through the ordinary fight or reproduction path

## Done Criteria

- current fight feasibility influences fallback same-shape target ordering
- the rule is deterministic and documented
- existing food-priority, low-energy, shape-aware, and reproduction-feasibility-aware semantics remain coherent
- current fight and reproduction resolution rules remain unchanged
- tests cover non-losing pursuit and losing-target fallback behavior

## Out Of Scope Follow-Ups

- explicit retreat behavior
- multi-target tactical scoring
- personality-based aggression
- detached child autonomy
- removing current radius shortcuts
