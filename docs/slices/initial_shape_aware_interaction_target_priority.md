# Slice: Initial Shape-Aware Interaction Target Priority

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now distinguishes interaction targets by shape outcome

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous target selection.

This slice extends the current deterministic interaction-seeking autonomy by making shape meaning influence which interaction target is preferred once a circle is in interaction-seeking mode.

## Discovery Scope

Establish the smallest deterministic rule that makes autonomous target choice aware of shape meaning:

- when interaction seeking is active, autonomous circles may prefer one shape outcome over another
- target selection remains deterministic and server-authoritative
- the current movement, energy, contact, fight, and reproduction rules remain unchanged after contact is reached
- the existing food-priority and low-energy rules remain part of the model

This slice does **not** attempt to implement:

- explicit combat strategy
- multi-factor utility scoring
- threat evaluation
- personality systems
- separate mate-seeking and hunt-seeking state machines
- client-side anticipation or prediction

## Why This Slice Next

The current implementation now lets autonomous circles choose between food and social interaction based on energy. But once a circle is ready to seek interaction, its target choice is still shape-blind. That leaves an important semantic gap because shape is the variable that changes interaction meaning.

The model pressure is now:

- shape should influence pre-contact behavior, not only post-contact resolution
- autonomous steering should feel more aligned with the existing fight vs reproduction semantics
- the ecosystem should remain simple while becoming less arbitrary

This slice is the narrowest next step because it:

- changes only interaction-target ordering
- keeps the current movement and outcome rules unchanged
- uses a domain concept already central to the model
- avoids introducing explicit tactical AI

## Use-Case Contract

### Use Case

`PrioritizeInteractionTargetByShapeMeaning`

### Primary Actor

Any autonomous circle that is already in interaction-seeking mode.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- food-priority and low-energy steering rules already exist
- interaction-seeking target selection already exists
- current shape-based interaction outcomes already exist

### Trigger

An autonomous circle selects an interaction target after food does not currently win priority.

### Success Outcome

- the autonomous circle may prefer a target whose shape implies the documented interaction outcome
- later snapshots show more legible pursuit behavior that matches existing shape semantics
- downstream fight and reproduction rules remain unchanged once contact occurs

### Failure Or Rejection Cases

- if shape remains irrelevant to target choice, steering stays too disconnected from interaction meaning
- if target choice uses too many hidden conditions, inspectability weakens
- if target choice becomes non-deterministic, reproducibility weakens

## Main Business Rules

1. Autonomous target selection remains authoritative server-side behavior.
2. Shape meaning may affect interaction-target priority once food does not currently win.
3. The rule must remain deterministic for the same world state and tick.
4. Current food-priority distance and low-energy food-recovery rules remain unchanged.
5. Once contact occurs, the current fight and reproduction rules remain unchanged.
6. This slice should choose one simple documented priority rule rather than a large scoring model.
7. Player and autonomous targets remain part of the same candidate set once interaction seeking is active.

## Minimal Domain Concepts In Scope

- `Interaction-Seeking Mode`
- `Shape Meaning`
- `Fight Candidate`
- `Reproduction Candidate`
- `Deterministic Target Priority`
- `World Snapshot`

## Bounded Priority Interpretation

This slice chooses the smallest inspectable interpretation:

- once food does not currently win, autonomous circles prefer the nearest target of one documented shape outcome
- if no preferred-outcome target exists, they fall back to the nearest remaining eligible target
- no new contact or outcome semantics are introduced

This avoids the larger step of explicit survival strategy while still making shape matter earlier in the behavior chain.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing shape labels and interaction outcomes
- the current HUD and rendering

Build should extend the contract only if steering provenance is necessary to explain the new shape-driven target choice.

## Required Ports Or Boundaries

- server-side interaction-target ordering that can distinguish preferred shape outcomes
- deterministic tie-breaking across player and autonomous targets
- tests that show shape-aware pursuit without changing downstream rules
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current interaction target chooser rather than adding another steering layer
- keep the rule to one clear priority, for example preferring different-shape reproduction opportunities before same-shape fight opportunities
- preserve the existing low-energy and nearby-food priority rules
- keep movement speed, energy cost, and all downstream interaction semantics unchanged
- avoid inventing named moods or role categories

## Initial Test Plan

### Server tests

- when interaction seeking is active, an autonomous circle prefers the documented shape outcome target over an equally valid alternative
- if no preferred-outcome target exists, the current nearest-target fallback still applies
- the low-energy food-recovery rule still wins when it should
- tie-breaking remains deterministic across player and autonomous candidates

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing a shape-aware autonomous target choice before contact
- resulting encounters still resolve through the ordinary interaction HUD and world snapshots

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle enters interaction-seeking mode
2. more than one eligible circle target exists with different shape outcomes
3. the autonomous circle chooses the target matching the documented shape priority rule
4. contact and resolution follow the current fight or reproduction path

## Done Criteria

- shape meaning influences autonomous target ordering
- the rule is deterministic and documented
- food-priority and low-energy rules remain coherent
- downstream interaction semantics remain unchanged
- tests cover preferred-outcome choice and fallback behavior

## Out Of Scope Follow-Ups

- threat avoidance
- explicit aggression systems
- personality-based target choice
- detached child autonomy
- removing current radius shortcuts
