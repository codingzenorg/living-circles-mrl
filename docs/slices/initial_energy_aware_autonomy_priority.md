# Slice: Initial Energy-Aware Autonomy Priority

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering shifts driven by energy condition

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic interaction-seeking autonomy by making energy level influence whether an autonomous circle prioritizes food or social contact.

## Discovery Scope

Establish the smallest deterministic rule that makes autonomous steering sensitive to energy condition:

- low-energy autonomous circles prefer food recovery over interaction seeking
- sufficiently stable autonomous circles may continue to seek interaction targets
- the current movement, energy cost, contact, fight, and reproduction rules remain unchanged after steering selects a target
- target choice remains deterministic and server-authoritative

This slice does **not** attempt to implement:

- tactical retreat from stronger circles
- explicit fear or aggression modes
- multi-step planning
- hunger decay beyond the current energy model
- separate species or personality systems
- client-side anticipation or prediction

## Why This Slice Next

Autonomous circles can now seek other circles, including the player, which makes the world more active. But that steering still ignores the central model variable: energy. A nearly starving circle currently pursues interaction just like a well-fed one whenever food is not inside the current priority distance.

The model pressure is now:

- energy should shape not only outcomes but also pre-contact behavior
- the ecosystem should show clearer self-preservation pressure before collapse
- autonomous motion should stay emergent without feeling strategically blind

This slice is the narrowest next step because it:

- changes only target-priority logic
- keeps the existing movement and outcome rules unchanged
- strengthens energy as the unifying variable
- avoids introducing explicit tactical AI

## Use-Case Contract

### Use Case

`PrioritizeAutonomousTargetByEnergyCondition`

### Primary Actor

Any autonomous circle selecting its next movement direction.

### Pre-conditions

- autonomous steering is already deterministic and server-authoritative
- food seeking already exists
- interaction-seeking target selection already exists
- current energy gain, energy cost, and interaction outcome rules already exist

### Trigger

An autonomous circle advances a simulation tick and determines its steering direction.

### Success Outcome

- the autonomous circle may prefer food when its energy is below a documented threshold
- once above that threshold, the existing interaction-seeking rule may apply again
- later snapshots show more legible shifts between self-preservation and social engagement

### Failure Or Rejection Cases

- if low-energy circles still chase interaction indiscriminately, energy remains too weak as a steering variable
- if steering flips through many hidden conditions, inspectability weakens
- if target choice becomes non-deterministic, reproducibility weakens

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Energy level may affect whether food or interaction is prioritized.
3. The rule must remain deterministic for the same world state and tick.
4. Food priority should become stronger, not weaker, for low-energy circles.
5. Current contact, fight, and reproduction rules remain unchanged after contact is reached.
6. This slice should use one simple documented energy threshold rather than a large behavior matrix.
7. Player and autonomous interaction targets remain governed by the current eligibility rule once interaction seeking is active.

## Minimal Domain Concepts In Scope

- `Autonomous Energy Condition`
- `Food Priority`
- `Interaction Priority`
- `Deterministic Threshold`
- `World Snapshot`

## Bounded Priority Interpretation

This slice chooses the smallest inspectable interpretation:

- below one energy threshold, autonomous circles prefer the nearest available food target
- at or above that threshold, they continue using the current interaction-seeking rule with the current food-priority exception
- no new outcome semantics are introduced

This avoids the larger step of full survival strategy while still making energy visibly shape behavior before collision.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing energy values
- existing interaction outcomes once contact occurs

Build should extend the contract only if steering provenance is necessary to explain the energy-driven shift.

## Required Ports Or Boundaries

- server-side autonomous target-priority logic that considers energy
- deterministic tests covering low-energy and stable-energy steering choices
- client rendering that remains sufficient to observe the changed motion

## Build Guidance

- prefer evolving the current autonomous-intent function rather than adding a second steering layer
- use one clear threshold and document it explicitly
- preserve the current food-priority distance rule unless the slice requires a small adjustment
- keep movement speed, energy cost, and all downstream interaction semantics unchanged
- avoid inventing named moods or tactical categories

## Initial Test Plan

### Server tests

- a low-energy autonomous circle prefers food over an otherwise eligible interaction target
- an adequately energized autonomous circle still pursues an interaction target under the current rule
- target choice remains deterministic across equal-distance cases
- the current food-priority-distance rule remains coherent alongside the new energy threshold

### Contract tests

- the current snapshot schema remains sufficient unless build adds steering provenance

### Integration tests

- the client receives snapshots showing a low-energy autonomous circle steering toward food instead of the player or another circle
- the client also receives snapshots showing a sufficiently energized autonomous circle still initiating interaction

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle with low energy selects a steering target
2. it prefers food recovery over social contact
3. a sufficiently energized autonomous circle still pursues an interaction target under the existing rule
4. the resulting motion is visible in ordinary snapshots

## Done Criteria

- energy level changes autonomous steering priority
- the rule is deterministic and documented
- current movement and interaction outcomes remain unchanged
- player movement is not required to observe the new behavior
- tests cover both low-energy recovery seeking and stable-energy interaction seeking

## Out Of Scope Follow-Ups

- threat avoidance
- multi-threshold behavior trees
- explicit aggression or fear systems
- detached child autonomy
- removing current radius shortcuts
