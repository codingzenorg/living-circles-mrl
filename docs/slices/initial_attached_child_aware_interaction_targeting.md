# Slice: Initial Attached-Child-Aware Interaction Targeting

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now treats attached-child proximity as part of positive interaction seeking

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic interaction-seeking rules by letting orbiting attached children influence which social target counts as nearest once food and avoidance rules no longer dominate.

## Discovery Scope

Establish the smallest deterministic rule that makes positive pursuit reflect the current child-triggered contact model:

- when interaction-seeking mode is active, target nearness may be satisfied by a parent body or one of its attached children
- the current low-energy food-recovery, same-shape threat-avoidance, blocked-reproduction avoidance, reproduction-feasibility-aware preference, and fight-feasibility-aware fallback rules remain in force
- the current movement, energy, contact, fight, reproduction, and food rules remain unchanged after steering selects a direction

This slice does **not** attempt to implement:

- predictive intercept steering
- child-only pursuit or child-only outcomes
- detached child autonomy
- generalized influence maps
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now uses attached-child geometry for:

- authoritative contact triggering
- threat avoidance
- blocked-reproduction avoidance

But ordinary positive interaction seeking still mostly reasons from parent-core positions. That leaves another asymmetry in the embodied model: attached children can start contact and can trigger retreat, yet they still do not help make a target feel socially near enough to pursue.

The model pressure is now:

- the same visible orbiting-child geometry should matter for positive pursuit as well as for contact and avoidance
- autonomous steering should feel less surprising when contact can be initiated through children but attraction still notices only parent centers
- pursuit should become more semantically aligned without inventing a larger tactical system

This slice is the narrowest next step because it:

- changes only the notion of social nearness inside the current interaction-seeking path
- keeps avoidance priority and downstream interaction outcomes unchanged
- reuses the existing authoritative attached-child layout instead of inventing new entities or hidden scores
- avoids broader prediction or path-planning work

## Use-Case Contract

### Use Case

`SelectInteractionTargetWithAttachedChildAwareness`

### Primary Actor

Any autonomous circle that is already in interaction-seeking mode.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- attached-child layout is already authoritative and deterministic per tick
- attached children can already trigger parent-level interaction contact
- attached children can already influence avoidance detection
- current reproduction-feasibility-aware and fight-feasibility-aware target rules already exist

### Trigger

An autonomous circle advances one simulation tick and selects a positive social target after food and avoidance priorities are evaluated.

### Success Outcome

- a target may count as nearer because one of its attached children is currently closer than its parent body
- the acting autonomous circle steers toward the chosen effective contact point under the existing priority rules
- later snapshots show pursuit that better matches the already-implemented child-triggered contact model

### Failure Or Rejection Cases

- if child-triggered contact and avoidance exist but pursuit still ignores child geometry, steering remains semantically uneven
- if target selection depends on opaque scoring, inspectability weakens
- if the rule turns into predictive interception, slice scope is exceeded

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Current low-energy food-recovery, same-shape threat-avoidance, and blocked-reproduction avoidance rules remain unchanged.
3. Current reproduction-feasibility-aware and fight-feasibility-aware social target eligibility rules remain unchanged.
4. Attached-child positions may contribute to which eligible social target counts as nearest.
5. The rule must remain deterministic for the same world state and tick.
6. Once contact occurs, the current fight and reproduction rules remain unchanged.
7. This slice should reuse the current authoritative child layout rather than introducing predictive extrapolation.
8. Player and autonomous targets remain part of one candidate set for positive interaction seeking.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Effective Social Nearness`
- `Interaction-Seeking Mode`
- `Eligible Interaction Target`
- `Deterministic Child Layout`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- current target eligibility rules stay exactly the same
- what changes is the distance basis for already-eligible social targets: the parent body or one of its attached children may define the nearest effective contact point
- once the target is selected, steering goes toward that effective point rather than toward the parent center alone
- no new contact categories or outcome semantics are introduced

This avoids the larger step of predictive steering while still aligning pursuit with the current embodied contact model.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- already-exposed attached-child positions
- existing energy, child-count, and shape values
- existing `contact_origin` once interaction occurs

Build should extend the contract only if attached-child-aware pursuit is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side interaction target selection that can use current attached-child positions for distance checks
- deterministic tie-breaking across parent and child proximity paths
- tests that show child-aware pursuit occurs before parent-core proximity alone would justify it
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current interaction-target chooser rather than creating a new steering subsystem
- reuse the current authoritative attached-child layout at the current tick
- keep the current social target eligibility and priority rules unchanged
- preserve low-energy, food-priority, and existing avoidance rules
- keep movement speed, energy cost, and downstream interaction semantics unchanged
- avoid generalized influence maps, prediction, or intercept logic

## Initial Test Plan

### Server tests

- an eligible different-shape target whose attached child is nearer than its parent body is selected ahead of a farther effective target
- an eligible same-shape fallback target whose attached child is effectively nearer is selected under the current fight-feasibility rule
- a target whose parent and children are all effectively farther does not win selection
- deterministic tie-breaking remains stable when multiple child-based effective distances are equal

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle steer toward a target because of attached-child effective nearness before parent-core nearness alone would explain it
- the client also receives ordinary pursuit when no attached-child position changes the effective social distance ordering

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle enters interaction-seeking mode
2. multiple eligible social targets exist under the current rules
3. one target remains farther by parent-core distance alone
4. one attached child makes that target effectively nearer than the alternatives
5. the autonomous circle steers toward that effective contact point
6. any resulting contact resolves through the ordinary fight or reproduction path

## Done Criteria

- attached children can influence positive interaction target selection
- the rule is deterministic and documented
- existing avoidance categories and social eligibility rules remain unchanged
- current fight and reproduction resolution rules remain unchanged
- tests cover child-aware pursuit and non-pursuit cases

## Out Of Scope Follow-Ups

- predictive interception
- generalized field-based steering
- detached child autonomy
- child-only fight systems
- child-only reproduction systems
