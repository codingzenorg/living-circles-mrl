# Slice: Initial Interaction Seeking Autonomy

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible autonomous steering that now prefers potential circle interaction targets in bounded cases

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for autonomous steering decisions.

This slice extends the current deterministic food-seeking autonomy into a minimal social steering rule so the ecosystem can actively create encounters instead of only resolving them when geometry happens to line up.

## Discovery Scope

Establish the smallest deterministic rule that lets autonomous circles deliberately seek other circles:

- autonomous circles can identify a nearest potential interaction target
- steering remains deterministic and server-authoritative
- the current movement, energy, contact, fight, and reproduction rules remain unchanged after contact is reached
- food seeking remains part of the model instead of being discarded

This slice does **not** attempt to implement:

- full behavior trees
- tactical avoidance
- threat evaluation
- separate mate-seeking versus fight-seeking subsystems
- player prediction or client-side simulation
- multi-goal optimization

## Why This Slice Next

The simulation can now resolve autonomous-autonomous encounters, which removes a major player-centered asymmetry. But autonomous circles still mostly create those encounters only by chance because their steering policy is still primarily food-seeking plus fallback drift.

The model pressure is now:

- autonomous circles should help create the ecosystem interactions the simulation now knows how to resolve
- emergence is stronger when encounters can arise from the system’s own steering, not only from initial placement or player action
- the project should keep avoiding explicit scripted roles while still letting non-player circles behave more like participants than drifting props

This slice is the narrowest next step because it:

- changes only steering, not outcome rules
- keeps autonomy deterministic
- increases encounter frequency without requiring new entity types
- makes the autonomous-autonomous slice materially visible in ordinary play

## Use-Case Contract

### Use Case

`SteerAutonomousCircleTowardInteractionOpportunity`

### Primary Actor

Any autonomous circle that is choosing its next movement direction.

### Pre-conditions

- autonomous circles already move under deterministic server authority
- food seeking already exists
- current interaction rules already exist for player-autonomous and autonomous-autonomous pairs

### Trigger

An autonomous circle advances a simulation tick and determines its steering direction.

### Success Outcome

- the autonomous circle may steer toward a nearest potential circle interaction target when the selected rule says to prefer interaction over drift or food
- later snapshots show more naturally produced encounters without changing outcome semantics
- player and non-player circles remain under the same downstream interaction rules once contact occurs

### Failure Or Rejection Cases

- if steering remains effectively food-only, the new autonomous-autonomous interaction engine stays underused
- if steering becomes too hand-scripted or role-like, the model drifts away from emergence
- if autonomous target choice becomes non-deterministic, inspectability weakens

## Main Business Rules

1. Autonomous steering remains authoritative server-side behavior.
2. Steering must remain deterministic for the same world state and tick.
3. Food seeking remains a valid steering concern in this slice.
4. Build should choose one bounded rule for when interaction seeking overrides drift or food seeking.
5. Once contact occurs, the current fight and reproduction rules remain unchanged.
6. Player and autonomous circles remain eligible interaction targets under the same steering rule unless build documents a narrower scope.
7. This slice should avoid explicit role assignment or hard-coded personality types.

## Minimal Domain Concepts In Scope

- `Autonomous Steering`
- `Interaction Opportunity`
- `Nearest Target`
- `Deterministic Target Selection`
- `World Snapshot`

## Bounded Autonomy Interpretation

This slice chooses the smallest inspectable interpretation:

- an autonomous circle may select a nearest eligible circle as its steering target under a documented condition
- target selection uses deterministic ordering and tie-breaking
- no new combat or reproduction semantics are introduced

This avoids the larger step of full agent logic while still making the ecosystem more self-propelled.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- ordinary movement changes in snapshots
- existing interaction outcomes once contact happens
- current labels and debug information in the client

Build should extend the contract only if steering provenance is necessary to understand the new motion.

## Required Ports Or Boundaries

- server-side autonomous target-selection logic
- deterministic tie-breaking for target choice
- tests that show autonomous circles now create encounters without player movement
- client rendering that remains sufficient to observe the resulting motion

## Build Guidance

- prefer evolving the current autonomous-intent function rather than adding a separate AI layer
- keep food seeking in the model and document clearly when interaction seeking takes priority
- choose a minimal deterministic target-selection rule
- preserve current movement speed, energy cost, and contact handling
- avoid inventing rich strategy labels or personality systems

## Initial Test Plan

### Server tests

- an autonomous circle can steer toward another circle and create a same-shape encounter without player movement
- an autonomous circle can steer toward another circle and create a different-shape encounter without player movement
- target selection is deterministic when more than one eligible circle exists
- the current food-seeking rule remains coherent under the new interaction-seeking priority rule

### Contract tests

- the current snapshot schema remains sufficient unless build exposes steering provenance

### Integration tests

- the client receives snapshots showing autonomous circles creating an interaction without player movement input
- resulting encounters remain visible through the ordinary canvas and interaction HUD

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an autonomous circle selects a nearby circle interaction target
2. it steers toward that target under ordinary movement rules
3. contact occurs and the existing fight or reproduction logic resolves
4. the resulting state change appears in ordinary snapshots without the player causing it

## Done Criteria

- autonomous circles can deliberately steer toward interaction opportunities
- target selection is deterministic
- current fight and reproduction semantics remain unchanged after contact
- player movement is not required to observe the new behavior
- tests cover encounter creation and target-selection determinism

## Out Of Scope Follow-Ups

- avoidance or escape behavior
- multi-objective planning
- explicit personality types
- detached child autonomy
- removing current radius shortcuts
