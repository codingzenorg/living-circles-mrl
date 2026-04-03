# Slice: Initial Attached Child To Child Contact Triggers Parent Interaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible parent-level interaction triggered by orbiting child-to-child contact

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether two attached orbiting children touching each other can trigger the same parent-level interaction semantics already used for parent-body and child-to-parent contact.

This slice extends the current embodied child-contact model without introducing a separate child combat subsystem or free child autonomy.

## Discovery Scope

Establish the smallest deterministic rule that makes two orbiting child bodies count as meaningful contact between their owning parents:

- if an attached child of one parent overlaps an attached child of another parent, that contact can trigger the owning parents' current interaction classification
- same-shape parent pairs still trend toward the current fight path
- different-shape parent pairs still trend toward the current reproduction path
- interaction outcomes remain parent-level outcomes in this slice

This slice does **not** attempt to implement:

- child-only winners or losers
- direct child removal caused solely by child-to-child contact
- detached child autonomy
- new damage or stun systems
- removal of current parent-body or child-to-parent contact rules
- removal of current transitional `children_count` shortcuts

## Why This Slice Next

The previous slice made orbiting children matter in contact initiation against parent bodies. That was a necessary step, but it still leaves a visible inconsistency: two orbiting children can visibly pass through each other without authoritative meaning even though they are now otherwise treated as part of encounter initiation.

The model pressure is now:

- visible child swarms should matter when they meet each other, not only when they touch a parent core
- the orbiting-child model should feel more spatially coherent
- the repository needs a clearer bridge from abstract radius shortcuts toward visibly embodied contact

This slice is the narrowest next step because it:

- preserves the current parent-level fight and reproduction rules
- reuses the new contact-provenance contract
- extends one already-established interpretation of children as contact-capable bodies
- avoids inventing a separate child lifecycle

## Use-Case Contract

### Use Case

`TriggerParentInteractionFromAttachedChildToChildContact`

### Primary Actor

Any circle that owns at least one attached orbiting child and comes into child-to-child contact with another parent's attached child.

### Pre-conditions

- snapshots already expose attached-child positions
- parent-level interaction rules already exist for same-shape and different-shape parent pairs
- both parents own at least one attached child

### Trigger

During authoritative overlap evaluation, an attached child of one parent overlaps an attached child of another parent.

### Success Outcome

- the owning parents enter the same parent-level interaction path that ordinary parent-body contact would have triggered
- child-to-child meeting becomes a valid authoritative reason for the interaction
- snapshots remain inspectable through visible child positions and contact provenance

### Failure Or Rejection Cases

- if child-to-child contact still has no meaning, the orbiting-child model remains spatially inconsistent
- if child-to-child contact creates a separate incompatible rule set, the model becomes harder to reason about
- if one pair can trigger more than one interaction in the same tick through multiple child bodies, inspectability weakens

## Main Business Rules

1. Parent-level interaction resolution remains authoritative server-side behavior.
2. Attached-child-to-attached-child overlap may trigger the same parent-level interaction classification as parent-body or child-to-parent overlap.
3. Same-shape parent pairs still resolve through the current fight path.
4. Different-shape parent pairs still resolve through the current reproduction path.
5. One parent pair may still resolve at most one interaction per relevant overlap window.
6. Build should preserve deterministic de-duplication when parent-body, child-to-parent, and child-to-child contact all occur in the same tick.
7. This slice keeps interaction outcomes parent-level rather than introducing child-specific outcome entities.

## Minimal Domain Concepts In Scope

- `Attached Child To Child Contact`
- `Parent-Level Interaction Trigger`
- `Same-Shape Fight Path`
- `Different-Shape Reproduction Path`
- `World Snapshot`

## Bounded Contact Interpretation

This slice chooses the smallest inspectable interpretation:

- treat attached children as valid contact bodies against each other for interaction triggering
- preserve the existing parent-level interaction outcomes after contact is detected
- reuse the current `contact_origin` vocabulary if it is still sufficient for inspectability

This avoids the larger step of creating direct child combat while still making the visible orbiting swarm matter more fully in encounter initiation.

## Required Runtime Contract Changes

The current contract is likely sufficient if build can continue using:

- attached-child positions already present in the snapshot
- the current parent-level interaction object
- the existing `contact_origin` field, unless child-to-child provenance requires a more specific value

Build should extend the contract only if `attached_child` is too coarse to distinguish child-to-parent from child-to-child contact for inspectability.

## Required Ports Or Boundaries

- server-side overlap detection that considers child-to-child contact across parent pairs
- deterministic pair-level de-duplication across all valid contact paths
- deterministic tests covering same-shape and different-shape child-to-child-triggered interaction
- client rendering that remains sufficient to inspect child-to-child contact

## Build Guidance

- prefer extending the current contact-detection entry points rather than creating a third interaction pipeline
- preserve the current fight, reproduction, and overlap-window rules after contact detection succeeds
- keep player and autonomous circles under the same child-to-child contact rule
- keep `contact_origin` stable unless a finer-grained value is truly necessary
- do not introduce special child-only combat states in this slice

## Initial Test Plan

### Server tests

- same-shape parent interaction can be triggered by child-to-child contact
- different-shape parent interaction can be triggered by child-to-child contact
- one pair still resolves at most one interaction during a continuous overlap window even when multiple child bodies overlap
- combined parent-body, child-to-parent, and child-to-child contact still resolves only one interaction per pair per tick

### Contract tests

- the current snapshot schema remains sufficient unless build needs more specific provenance than `attached_child`

### Integration tests

- the client receives a snapshot where visible child swarms trigger a parent-level interaction
- current rendering remains sufficient to inspect that child-to-child contact contributed to the interaction

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two parents approach such that their orbiting children meet before or while the parent cores remain apart
2. the server resolves the resulting parent-level interaction
3. the next snapshot shows the ordinary interaction outcome while visible orbiting children explain the contact path

## Done Criteria

- attached-child-to-attached-child contact can trigger parent-level interaction resolution
- same-shape and different-shape semantics remain unchanged after contact is detected
- one pair does not double-trigger from combined contact paths in the same tick
- player and autonomous circles follow the same rule
- tests cover same-shape, different-shape, and de-duplication behavior

## Out Of Scope Follow-Ups

- child-only combat outcomes
- detached child autonomy
- child-specific lineage branching
- removing the current parent-radius shortcuts
- full contact history or replay traces
