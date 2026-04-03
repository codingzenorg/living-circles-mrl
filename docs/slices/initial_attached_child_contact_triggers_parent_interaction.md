# Slice: Initial Attached Child Contact Triggers Parent Interaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible interaction resolution driven by attached-child contact

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether a visible orbiting child touching another parent body can trigger the same parent-level interaction semantics already used for direct parent overlap.

This slice extends the current orbiting-child model into contact detection without introducing free child autonomy or a separate child-to-child combat system.

## Discovery Scope

Establish the smallest deterministic rule that lets visible attached children matter directly in collision initiation:

- if an attached child overlaps another parent body, that contact can trigger the owning parents' existing interaction classification
- same-shape contact still trends toward fight under the current fight rules
- different-shape contact still trends toward reproduction under the current reproduction rules
- interaction outcomes remain parent-level outcomes in this slice

This slice does **not** attempt to implement:

- child-to-child interactions as a separate combat domain
- free child detachment
- child-only lineages
- new damage systems
- removal of current parent-radius overlap rules
- replacement of current transitional `children_count` shortcuts

## Why This Slice Next

Orbiting children are now visible and mechanically meaningful in feeding, payment, conflict absorption, and continuity. But they still do not participate in how contact begins. That leaves a gap between the visible orbital model and the authoritative collision model.

The model pressure is now:

- visible child bodies should matter in interaction initiation, not only after an interaction is already underway
- orbiting children should help explain why a parent can reach or threaten another circle beyond its core body
- the repository needs an executable bridge from radius shortcuts toward a more embodied orbiting-child interaction model

This slice is the narrowest next step because it:

- keeps the existing parent-level fight and reproduction rules
- reuses the existing interaction contract
- makes orbiting children affect an already-central system boundary: contact
- avoids inventing a separate child entity lifecycle

## Use-Case Contract

### Use Case

`TriggerParentInteractionFromAttachedChildContact`

### Primary Actor

Any circle that owns at least one attached orbiting child.

### Pre-conditions

- snapshots already expose attached-child positions
- parent-level interaction rules already exist for same-shape and different-shape contact
- a parent or one of its attached children can physically overlap another parent body

### Trigger

During authoritative overlap evaluation, an attached child of one parent overlaps the body of another parent.

### Success Outcome

- the owning parents enter the same parent-level interaction path that ordinary parent-body contact would have triggered
- visible attached-child positions become part of why an interaction happened
- the interaction object remains sufficient unless build finds a strong need for contact provenance

### Failure Or Rejection Cases

- if attached children remain irrelevant to contact initiation, the orbiting model stays mechanically partial
- if child contact triggers a separate incompatible rule set, the model becomes harder to reason about
- if parent and child contact can double-trigger the same interaction in one tick, inspectability weakens

## Main Business Rules

1. Parent-level contact detection remains authoritative server-side behavior.
2. An attached child overlapping another parent body may trigger the same parent-level interaction classification as direct parent overlap.
3. Same-shape child-originated contact still resolves through the current fight path.
4. Different-shape child-originated contact still resolves through the current reproduction path.
5. One parent pair may still resolve at most one interaction per relevant overlap window.
6. Build should avoid duplicate triggering when parent body and attached child both overlap the same target in the same tick.
7. This slice keeps interaction outcomes parent-level rather than creating child-specific winner or loser entities.

## Minimal Domain Concepts In Scope

- `Attached Child Contact`
- `Parent-Level Interaction Trigger`
- `Same-Shape Fight Path`
- `Different-Shape Reproduction Path`
- `World Snapshot`

## Bounded Contact Interpretation

This slice chooses the smallest inspectable interpretation:

- treat an attached child as an extension of its parent's contact reach for interaction triggering only
- preserve existing parent-level interaction resolution after contact is detected
- allow build to keep the current interaction object unless contact provenance is necessary to understand the resulting behavior

This avoids the larger step of turning attached children into fully independent combatants while still making the visible orbit model materially affect interaction start conditions.

## Required Runtime Contract Changes

The current contract is likely sufficient if the effect is visible through:

- earlier or otherwise explainable interaction resolution
- attached-child positions already present in the snapshot
- the existing interaction object

Build should extend the contract only if contact provenance is necessary to distinguish child-originated interaction from parent-body interaction.

## Required Ports Or Boundaries

- server-side overlap detection that considers attached-child positions against other parent bodies
- deterministic pair-level de-duplication so one contact path does not trigger multiple outcomes in a tick
- deterministic tests covering same-shape and different-shape child-originated contact
- client rendering that remains unchanged except for naturally showing the visible child positions that explain the contact

## Build Guidance

- prefer extending the current overlap-detection entry points rather than introducing a second interaction pipeline
- preserve the existing fight, reproduction, and overlap-window rules after contact detection succeeds
- keep player and autonomous circles under the same contact-trigger rule
- do not remove the current parent-body overlap path in this slice
- do not add special child-only combat states

## Initial Test Plan

### Server tests

- an attached child of a same-shape parent can trigger the current fight path against another parent body
- an attached child of a different-shape parent can trigger the current reproduction path against another parent body
- one pair still resolves at most one interaction during a continuous overlap window even when child contact is involved
- parent-body and child-body overlap together do not create duplicate interaction resolutions

### Contract tests

- the current snapshot schema remains sufficient unless build adds explicit contact provenance

### Integration tests

- the client receives a snapshot where interaction resolution occurs while visible attached-child positions explain the contact path
- default rendering remains enough to inspect that orbiting children helped trigger interaction

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. a parent's orbiting child reaches another parent body before the two parent cores overlap
2. the server resolves the resulting parent-level interaction
3. the next snapshot shows the ordinary interaction outcome while the visible child positions explain why contact occurred

## Done Criteria

- attached-child contact can trigger parent-level interaction resolution
- same-shape and different-shape semantics remain unchanged after contact is detected
- one pair does not double-trigger from combined parent and child overlap in the same tick
- player and autonomous circles follow the same rule
- tests cover same-shape, different-shape, and de-duplication behavior

## Out Of Scope Follow-Ups

- separate child-to-child combat
- detached child autonomy
- replacing current parent-radius shortcuts
- provenance-rich interaction history
- child-specific lineage branching
