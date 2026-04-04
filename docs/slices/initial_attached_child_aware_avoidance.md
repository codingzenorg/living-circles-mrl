# Slice: Initial Attached-Child-Aware Avoidance

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now treats attached-child proximity as part of avoidance detection

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic avoidance rules by letting orbiting attached children participate in threat and blocked-reproduction detection before parent core bodies overlap.

## Discovery Scope

Establish the smallest deterministic rule that makes avoidance reflect the current child-triggered contact model:

- nearby attached-child positions may qualify a same-shape threat or blocked different-shape target even when parent cores are still farther apart
- the current movement, energy, contact, fight, reproduction, and food rules remain unchanged after steering selects a direction
- the existing low-energy, same-shape threat-avoidance, blocked-reproduction avoidance, reproduction-feasibility-aware, and fight-feasibility-aware rules remain part of the model

This slice does **not** attempt to implement:

- full field-based steering or influence maps
- collision prediction beyond current authoritative child layout
- detached child autonomy
- new child-only fight or child-only reproduction outcomes
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now has explicit negative steering for two social cases:

- nearby stronger same-shape threats
- nearby different-shape targets with blocked reproduction

But both avoidance rules still mainly reason from parent-core positions, while the simulation already allows authoritative interaction to start through attached-child-to-parent and attached-child-to-attached-child contact. That leaves a clear mismatch: child-triggered contact exists in the world model, but child-triggered avoidance pressure does not.

The model pressure is now:

- pre-contact steering should respect the same embodied orbiting-child geometry that already matters for authoritative contact
- visible orbiting children should influence danger and blocked-reproduction pressure before parent cores overlap
- the ecosystem should feel less surprising when contact can start from children but avoidance only notices parents late

This slice is the narrowest next step because it:

- only widens the detection basis for existing avoidance rules
- keeps downstream fight and reproduction outcomes unchanged
- reuses the existing authoritative orbiting-child layout instead of inventing new entities or prediction systems
- avoids a larger move into continuous steering fields

## Use-Case Contract

### Use Case

`DetectAvoidancePressureFromAttachedChildren`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- attached-child layout is already authoritative and deterministic per tick
- attached children can already trigger parent-level interaction contact
- current same-shape threat-avoidance and blocked-reproduction avoidance rules already exist

### Trigger

An autonomous circle advances one simulation tick and evaluates nearby avoidance pressure.

### Success Outcome

- an attached child may cause a same-shape threat to count as nearby before parent cores overlap
- an attached child may cause a blocked different-shape target to count as nearby before parent cores overlap
- later snapshots show avoidance that better matches the already-implemented child-triggered contact model

### Failure Or Rejection Cases

- if child-triggered contact remains possible before child-triggered avoidance, steering stays semantically behind the world model
- if avoidance begins depending on opaque spatial heuristics, inspectability weakens
- if the rule becomes a general prediction system, slice scope is exceeded

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Existing same-shape threat-avoidance and blocked-reproduction avoidance rules remain the only avoidance categories in scope.
3. Attached-child positions may contribute to whether a target counts as nearby for those existing categories.
4. The rule must remain deterministic for the same world state and tick.
5. The current low-energy food-recovery rule remains unchanged.
6. Once contact occurs, the current fight and reproduction rules remain unchanged.
7. This slice should reuse the current authoritative child layout rather than introducing predictive extrapolation.
8. Player and autonomous targets remain part of one candidate set for avoidance evaluation.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Child-Triggered Avoidance Pressure`
- `Nearby Threat`
- `Nearby Blocked Reproduction Target`
- `Deterministic Child Layout`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- existing avoidance categories stay exactly the same
- what changes is the notion of nearness: a parent body or one of its attached children may satisfy the current proximity test
- if neither parent nor attached children bring the target into the current avoidance window, the existing steering rules continue unchanged
- no new contact or outcome semantics are introduced

This avoids the larger step of field simulation while still aligning avoidance with the current embodied contact model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- already-exposed attached-child positions
- existing energy, child-count, and shape values

Build should extend the contract only if child-triggered avoidance is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side avoidance detection that can use current attached-child positions for proximity checks
- deterministic tie-breaking across parent and child proximity paths
- tests that show child-triggered avoidance occurs before parent-core overlap
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current avoidance detection path rather than creating a new steering subsystem
- reuse the current authoritative attached-child layout at the current tick
- keep the current avoidance categories and thresholds unless a minimal documented adjustment is required
- preserve low-energy, food-priority, and existing pursuit rules
- keep movement speed, energy cost, and downstream interaction semantics unchanged
- avoid generalized influence maps or prediction logic

## Initial Test Plan

### Server tests

- a same-shape threat whose attached child enters the avoidance window triggers retreat before parent-core overlap
- a blocked different-shape target whose attached child enters the avoidance window triggers retreat before parent-core overlap
- a target whose parent and children all remain outside the avoidance window does not trigger retreat
- deterministic tie-breaking remains stable when multiple child-based proximity paths are equally near

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing child-triggered retreat before parent bodies overlap
- the client also receives ordinary steering when child positions do not bring the target into the existing avoidance window

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle advances a tick
2. another circle remains outside the current avoidance window by parent-core distance alone
3. one attached child brings that target into effective avoidance proximity
4. the autonomous circle retreats according to the existing avoidance category
5. when no attached child creates that proximity, the current steering rules remain in control

## Done Criteria

- attached children can influence avoidance detection before parent-core overlap
- the rule is deterministic and documented
- existing avoidance categories remain unchanged
- current fight and reproduction resolution rules remain unchanged
- tests cover child-triggered avoidance and non-avoidance cases

## Out Of Scope Follow-Ups

- generalized influence-field steering
- detached child autonomy
- predictive path planning
- child-only fight systems
- child-only reproduction systems
