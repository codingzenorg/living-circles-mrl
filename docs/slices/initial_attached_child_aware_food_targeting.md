# Slice: Initial Attached-Child-Aware Food Targeting

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now treats attached-child proximity as part of food targeting

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic food-seeking rules by letting orbiting attached children influence which food slot counts as effectively nearest once food has already won priority.

## Discovery Scope

Establish the smallest deterministic rule that makes food-seeking reflect the current child-extended feeding model:

- when food-seeking is active, food nearness may be satisfied by a parent body or one of its attached children
- the current low-energy food-recovery, nearby-food priority, same-shape threat-avoidance, blocked-reproduction avoidance, reproduction-feasibility-aware preference, and fight-feasibility-aware fallback rules remain in force
- the current movement, energy, collection, regeneration, contact, fight, and reproduction rules remain unchanged after steering selects a direction

This slice does **not** attempt to implement:

- predictive path planning around future child orbit positions
- child-only detached foraging
- generalized influence maps
- food clustering heuristics
- client-side anticipation or prediction

## Why This Slice Next

The current implementation already lets attached children matter for:

- authoritative food collection
- authoritative contact triggering
- threat avoidance
- blocked-reproduction avoidance
- positive interaction seeking

But food targeting still mostly reasons from the parent-core position. That leaves a remaining asymmetry inside the core energy loop: attached children can physically collect food, yet autonomous circles still choose food targets as if only the parent center could reach them.

The model pressure is now:

- the same visible orbiting-child geometry should matter for food approach as well as for actual food collection
- autonomous feeding should feel less surprising when a child can collect a food slot that the steering logic did not treat as effectively near
- the energy loop should become more embodied without inventing a broader planning system

This slice is the narrowest next step because it:

- changes only the notion of food nearness inside the current food-targeting path
- keeps priority ordering and downstream food, contact, fight, and reproduction outcomes unchanged
- reuses the existing authoritative attached-child layout instead of inventing new entities or hidden scores
- avoids broader path-planning or resource-strategy work

## Use-Case Contract

### Use Case

`SelectFoodTargetWithAttachedChildAwareness`

### Primary Actor

Any autonomous circle selecting a food target after food has already won priority.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- attached-child layout is already authoritative and deterministic per tick
- attached children can already collect food on behalf of their parent
- current low-energy food priority and nearby-food priority rules already exist

### Trigger

An autonomous circle advances one simulation tick and selects a food target after the existing steering priorities say food should govern movement.

### Success Outcome

- a food slot may count as nearer because one attached child is currently closer to it than the parent body
- the acting autonomous circle steers toward the chosen effective food point under the existing priority rules
- later snapshots show food pursuit that better matches the already-implemented child-based food collection model

### Failure Or Rejection Cases

- if child-based collection exists but food targeting still ignores child geometry, the feeding loop remains semantically uneven
- if food choice depends on opaque scoring, inspectability weakens
- if the rule turns into predictive interception of future orbit positions, slice scope is exceeded

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Current low-energy food-recovery and nearby-food priority rules remain unchanged.
3. Current same-shape threat-avoidance and blocked-reproduction avoidance rules remain unchanged.
4. Current reproduction-feasibility-aware and fight-feasibility-aware social target rules remain unchanged.
5. Attached-child positions may contribute to which food slot counts as effectively nearest once food targeting is already active.
6. The rule must remain deterministic for the same world state and tick.
7. Once food is reached, the current collection, energy gain, and regeneration rules remain unchanged.
8. This slice should reuse the current authoritative child layout rather than introducing predictive extrapolation.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Effective Food Nearness`
- `Food-Seeking Mode`
- `Eligible Food Target`
- `Deterministic Child Layout`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- current food-priority rules stay exactly the same
- what changes is the distance basis for food slots once food targeting is active: the parent body or one of its attached children may define the nearest effective reach point
- once the target food slot is selected, steering still moves the parent body normally; no child-only motion is introduced
- no new collection, regeneration, or contact semantics are introduced

This avoids the larger step of predictive foraging while still aligning approach behavior with the current embodied feeding model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- already-exposed attached-child positions
- existing food positions
- existing parent energy changes after collection

Build should extend the contract only if attached-child-aware food pursuit is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side food target selection that can use current attached-child positions for distance checks
- deterministic tie-breaking across parent and child proximity paths
- tests that show child-aware food pursuit occurs before parent-core proximity alone would justify it
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current food-target selection path rather than creating a new steering subsystem
- reuse the current authoritative attached-child layout at the current tick
- keep the current food-priority thresholds and higher-priority avoidance rules unchanged
- preserve movement speed, energy cost, and downstream collection semantics
- avoid generalized influence maps, prediction, or resource planning

## Initial Test Plan

### Server tests

- a food slot whose effective distance is shorter through an attached child is selected ahead of a parent-core-nearer alternative
- a low-energy autonomous circle can choose the effectively nearer child-reachable food slot under the existing low-energy rule
- a food slot whose parent and child distances are all effectively farther does not win selection
- deterministic tie-breaking remains stable when multiple child-based effective food distances are equal

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle steer toward a food slot because of attached-child effective nearness before parent-core nearness alone would explain it
- the client also receives ordinary food pursuit when no attached-child position changes the effective food distance ordering

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle advances one tick in a state where food wins priority
2. multiple food slots are available
3. one food slot remains farther by parent-core distance alone
4. one attached child makes that food slot effectively nearer than the alternatives
5. the autonomous circle steers toward the food slot selected through effective child-aware reach
6. any resulting collection resolves through the ordinary food path

## Done Criteria

- attached children can influence food target selection
- the rule is deterministic and documented
- existing priority ordering and downstream collection semantics remain unchanged
- current contact, fight, and reproduction resolution rules remain unchanged
- tests cover child-aware food pursuit and non-pursuit cases

## Out Of Scope Follow-Ups

- predictive foraging
- generalized field-based steering
- detached child autonomy
- food clustering strategies
- child-only movement systems
