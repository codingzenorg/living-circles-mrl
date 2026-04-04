# Slice: Initial Threat-Avoidance Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now includes one bounded avoidance response to nearby stronger same-shape threats

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic autonomy rules by adding one narrow avoidance behavior when a nearby same-shape encounter would currently be losing.

## Discovery Scope

Establish the smallest deterministic threat response that keeps autonomous motion aligned with the existing fight model:

- when a nearby same-shape circle would currently beat the acting autonomous circle under the existing fight ordering, the acting circle may steer away instead of merely declining pursuit
- this avoidance should apply only inside one documented proximity window
- the current movement, energy, contact, fight, reproduction, and food rules remain unchanged after steering selects a direction
- the existing food-priority, low-energy, reproduction-feasibility-aware, and fight-feasibility-aware rules remain part of the model

This slice does **not** attempt to implement:

- generalized flocking or pathfinding
- long-horizon escape planning
- threat memory
- aggression personalities
- explicit coordination among autonomous circles
- client-side anticipation or prediction

## Why This Slice Next

The current implementation can now skip infeasible reproduction targets and skip same-shape targets that would currently lose. That is a cleaner chooser, but it still leaves one gap: a stronger same-shape circle can be nearby and socially relevant while the weaker circle merely falls back to food or drift. The model still lacks an explicit response to immediate hostile pressure.

The model pressure is now:

- pre-contact steering should acknowledge not only which fights are worth pursuing, but also which nearby fights are worth avoiding
- energy, children, and shape should influence motion before contact in both positive and negative directions
- the ecosystem should feel less passive when hostile asymmetry is already visible

This slice is the narrowest next step because it:

- adds one documented avoidance response instead of a larger tactical behavior system
- reuses the current deterministic fight ordering as the threat basis
- keeps downstream collision, fight, and reproduction semantics unchanged
- avoids broad AI concepts like fear states or planning trees

## Use-Case Contract

### Use Case

`AvoidNearbyLosingSameShapeThreat`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- current food-priority, low-energy, reproduction-feasibility-aware, and fight-feasibility-aware rules already exist
- current same-shape fight ordering already exists
- current movement intent selection already exists

### Trigger

An autonomous circle advances one simulation tick and evaluates its steering direction.

### Success Outcome

- if a nearby same-shape circle would currently beat the acting autonomous circle, the acting circle may steer away from that threat under one deterministic rule
- if no such nearby threat exists, current steering rules continue to apply
- later snapshots show more legible hostile asymmetry before contact occurs

### Failure Or Rejection Cases

- if weaker circles still drift into immediate stronger same-shape threats without any explicit response, steering remains too disconnected from the fight model
- if avoidance becomes a large tactical system, inspectability weakens
- if different threat cases resolve through hidden scoring, deterministic explanation becomes harder

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Threat avoidance may only respond to same-shape circles that would currently win under the existing deterministic fight ordering.
3. The rule must remain deterministic for the same world state and tick.
4. The current low-energy food-recovery rule remains unchanged.
5. The current reproduction-feasibility-aware preference remains unchanged.
6. Threat avoidance should use one documented proximity threshold rather than an open-ended tactical model.
7. Once contact occurs, the current fight and reproduction rules remain unchanged.
8. Player and autonomous circles remain part of one candidate set for evaluating nearby same-shape threats.

## Minimal Domain Concepts In Scope

- `Nearby Threat`
- `Same-Shape Losing Fight`
- `Avoidance Direction`
- `Deterministic Proximity Threshold`
- `World Snapshot`

## Bounded Avoidance Interpretation

This slice chooses the smallest inspectable interpretation:

- if a same-shape player or autonomous circle is within one documented threat radius and would currently beat the acting autonomous circle, that threat may override ordinary interaction pursuit
- avoidance means steering directly away from the nearest qualifying threat
- if no qualifying threat exists, the current food and interaction rules continue unchanged
- no new contact or outcome semantics are introduced

This avoids the larger step of generalized evasion while still making the current fight model visible in motion before collision.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing energy, child-count, and shape values
- existing interaction outcomes once contact occurs

Build should extend the contract only if the avoidance decision is too subtle to infer from ordinary snapshots.

## Required Ports Or Boundaries

- server-side steering logic that can evaluate nearby same-shape threats before ordinary pursuit
- deterministic threat selection across player and autonomous candidates
- tests that show nearby stronger same-shape threats cause retreat while non-threatening circles do not
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current autonomous-intent and interaction-target selection path rather than creating a separate AI subsystem
- reuse the existing fight ordering as the threat basis
- use one explicit distance threshold and document it
- preserve the existing low-energy, food-priority, reproduction-feasibility-aware, and fight-feasibility-aware rules
- keep movement speed, energy cost, and downstream interaction semantics unchanged
- avoid named moods, threat scores, or memory systems

## Initial Test Plan

### Server tests

- a nearby stronger same-shape player causes an autonomous circle to steer away instead of toward food or interaction
- a nearby stronger same-shape autonomous circle causes the same avoidance response
- a stronger same-shape circle outside the threat radius does not override the ordinary steering rule
- a weaker same-shape circle does not trigger avoidance

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing an autonomous circle retreat from a nearby stronger same-shape threat
- the client also receives ordinary behavior when the same target is outside the threat radius or not stronger

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle advances a tick
2. a nearby same-shape target exists that would currently win the fight
3. the autonomous circle evaluates the threat inside the documented radius
4. it steers away from that target
5. when the threat condition is not met, the current food or interaction rules remain in control

## Done Criteria

- nearby stronger same-shape threats can change autonomous steering direction
- the rule is deterministic and documented
- current food-priority and interaction-selection semantics remain coherent
- current fight and reproduction resolution rules remain unchanged
- tests cover both avoidance and non-avoidance cases

## Out Of Scope Follow-Ups

- generalized retreat from different-shape blocked reproduction
- long-horizon escape planning
- coordinated swarming behavior
- memory of past threats
- removing current radius shortcuts
