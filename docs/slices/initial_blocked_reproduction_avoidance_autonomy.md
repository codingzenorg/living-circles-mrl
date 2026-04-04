# Slice: Initial Blocked-Reproduction Avoidance Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now includes one bounded avoidance response to nearby different-shape targets when reproduction is currently blocked

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic autonomy rules by adding one narrow avoidance behavior when a nearby different-shape encounter is currently non-viable for reproduction.

## Discovery Scope

Establish the smallest deterministic blocked-reproduction response that keeps autonomous motion aligned with the existing reproduction model:

- when a nearby different-shape circle cannot currently satisfy the existing reproduction feasibility rule with the acting autonomous circle, the acting circle may steer away instead of treating that target as merely non-preferred
- this avoidance should apply only inside one documented proximity window
- the current movement, energy, contact, fight, reproduction, and food rules remain unchanged after steering selects a direction
- the existing low-energy, threat-avoidance, reproduction-feasibility-aware, and fight-feasibility-aware rules remain part of the model

This slice does **not** attempt to implement:

- generalized avoidance of every blocked social interaction
- long-horizon mate selection or planning
- fertility or personality systems
- memory of previously blocked pairs
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now expresses both positive and negative same-shape behavior before contact: circles can pursue survivable fights and retreat from nearby stronger threats. But different-shape behavior is still asymmetric. When reproduction is currently blocked, circles only stop preferring that target; they do not yet show an explicit negative response to nearby non-viable reproduction opportunities.

The model pressure is now:

- pre-contact steering should express blocked reproduction as clearly as it already expresses losing fight threat
- energy and child reserve should shape motion before contact in both attraction and avoidance paths
- the ecosystem should feel less semantically ambiguous when different-shape proximity is present but reproduction cannot succeed

This slice is the narrowest next step because it:

- adds one documented avoidance response instead of a broad social AI system
- reuses the existing reproduction feasibility rule as the sole basis for blocked reproduction
- keeps downstream collision, fight, and reproduction semantics unchanged
- avoids introducing fertility models, planning, or personality abstractions

## Use-Case Contract

### Use Case

`AvoidNearbyBlockedReproductionTarget`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- current low-energy, threat-avoidance, reproduction-feasibility-aware, and fight-feasibility-aware rules already exist
- current reproduction feasibility rule already exists
- current movement intent selection already exists

### Trigger

An autonomous circle advances one simulation tick and evaluates its steering direction.

### Success Outcome

- if a nearby different-shape player or autonomous circle cannot currently reproduce with the acting autonomous circle, the acting circle may steer away under one deterministic rule
- if no such nearby blocked target exists, current steering rules continue to apply
- later snapshots show more legible negative reproduction pressure before contact occurs

### Failure Or Rejection Cases

- if blocked reproduction remains visible only as absence of pursuit, the motion model stays more ambiguous than the fight side
- if avoidance becomes a broad tactical layer, inspectability weakens
- if blocked different-shape handling depends on hidden scores, deterministic explanation becomes harder

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Blocked-reproduction avoidance may only respond to different-shape circles that currently fail the existing reproduction feasibility rule.
3. The rule must remain deterministic for the same world state and tick.
4. The current low-energy food-recovery rule remains unchanged.
5. The current same-shape threat-avoidance rule remains unchanged.
6. Blocked-reproduction avoidance should use one documented proximity threshold rather than an open-ended social model.
7. Once contact occurs, the current fight and reproduction rules remain unchanged.
8. Player and autonomous circles remain part of one candidate set for evaluating nearby blocked different-shape targets.

## Minimal Domain Concepts In Scope

- `Nearby Blocked Reproduction Target`
- `Different-Shape Infeasibility`
- `Avoidance Direction`
- `Deterministic Proximity Threshold`
- `World Snapshot`

## Bounded Avoidance Interpretation

This slice chooses the smallest inspectable interpretation:

- if a different-shape player or autonomous circle is within one documented blocked-reproduction radius and the pair currently fails the existing reproduction feasibility rule, that blocked target may override ordinary interaction pursuit
- avoidance means steering directly away from the nearest qualifying blocked different-shape target
- if no qualifying blocked target exists, the current food and interaction rules continue unchanged
- no new contact or outcome semantics are introduced

This avoids the larger step of generalized social planning while still making the current reproduction model visible in motion before contact.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing energy, child-count, and shape values
- existing interaction outcomes once contact occurs

Build should extend the contract only if the avoidance decision is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side steering logic that can evaluate nearby blocked different-shape targets before ordinary pursuit
- deterministic target selection across player and autonomous candidates
- tests that show nearby blocked different-shape targets cause retreat while viable different-shape targets do not
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current autonomous-intent and target-selection path rather than creating a separate AI subsystem
- reuse the existing reproduction feasibility rule as the sole blocked-reproduction basis
- use one explicit distance threshold and document it
- preserve the existing low-energy, same-shape threat-avoidance, reproduction-feasibility-aware, and fight-feasibility-aware rules
- keep movement speed, energy cost, and downstream interaction semantics unchanged
- avoid moods, utility scores, or memory systems

## Initial Test Plan

### Server tests

- a nearby different-shape player that currently fails reproduction feasibility causes an autonomous circle to retreat
- a nearby different-shape autonomous circle that currently fails reproduction feasibility causes the same avoidance response
- a nearby feasible different-shape target does not trigger blocked-reproduction avoidance
- a blocked different-shape target outside the documented radius does not override ordinary steering

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle retreat from a nearby blocked different-shape target
- the client also receives ordinary pursuit when the different-shape target is feasible or outside the blocked-reproduction radius

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle advances a tick
2. a nearby different-shape target exists
3. the pair currently fails the reproduction feasibility rule
4. the autonomous circle evaluates the blocked target inside the documented radius
5. it steers away from that target
6. when the blocked condition is not met, the current food or interaction rules remain in control

## Done Criteria

- nearby blocked different-shape targets can change autonomous steering direction
- the rule is deterministic and documented
- current food-priority, same-shape threat-avoidance, and interaction-selection semantics remain coherent
- current fight and reproduction resolution rules remain unchanged
- tests cover both blocked-target avoidance and non-avoidance cases

## Out Of Scope Follow-Ups

- generalized avoidance of all blocked interactions
- long-horizon mate search
- explicit fertility systems
- memory of blocked pairs
- removing current radius shortcuts
