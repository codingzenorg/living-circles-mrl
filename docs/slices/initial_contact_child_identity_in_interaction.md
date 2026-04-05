# Slice: Initial Contact Child Identity In Interaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where child-triggered contact remains inspectable through explicit child identity

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for contact detection and interaction resolution.

This slice narrows one remaining contact inspectability gap by exposing which attached child or children actually triggered pair contact when interaction begins through attached-child geometry.

## Discovery Scope

Establish the smallest deterministic change that makes child-triggered contact identity explicit:

- when interaction contact begins through attached-child geometry, the same authoritative contact path should remain exactly as it is today
- the runtime should expose which side's attached child participated in contact, and which concrete child ID was involved
- contact origin, winner selection, reproduction payment, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new interaction rules
- alternate contact precedence rules
- historical contact logs
- redesign of contact origin categories
- new client-side effects beyond what is necessary to read contact-child identity

## Why This Slice Next

Recent slices made child-dependent outcomes much more inspectable:

- `death_promoted_child` now exposes the promoted child identity
- `fight_absorbed_child` now exposes the absorbed child identity
- `reproduce_paid_child` now exposes which side paid through a child

But child-triggered contact still stops one step short. The runtime already says whether contact came from `parent_body` or `attached_child`, yet it still hides which attached child actually touched first. That weakens inspectability because the server already knows the exact participating child identity while the client and tests can only infer it indirectly from positions.

The next model pressure is therefore not to change contact semantics. It is to make child-triggered contact identity explicit in the same way other child-dependent outcomes are now explicit.

This slice is the narrowest next step because it:

- changes only contact inspectability
- preserves current contact and interaction behavior
- avoids inventing event history or child-level combat systems
- gives build a deterministic contract extension instead of a larger semantic redesign

## Use-Case Contract

### Use Case

`ExposeContactChildIdentityDuringInteraction`

### Primary Actor

Any player or autonomous circle pair whose interaction begins through one or more attached-child bodies.

### Pre-conditions

- contact detection already distinguishes `parent_body` from `attached_child`
- attached-child-to-parent and attached-child-to-attached-child paths already exist
- the current interaction payload already exposes pair identity and contact origin

### Trigger

An interaction begins and `contact_origin` is `attached_child`.

### Success Outcome

- contact still resolves through the same authoritative path as before
- the authoritative interaction outcome explicitly exposes which source-side and/or target-side attached child participated in the triggering contact

### Failure Or Rejection Cases

- if `attached_child` contact still hides which child participated, inspectability remains incomplete
- if exposing contact-child identity changes interaction timing or outcome, scope is exceeded
- if the exposed child identity can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Contact detection remains authoritative server-side behavior.
2. `contact_origin` remains the current coarse contact category.
3. When `contact_origin` is `attached_child`, the runtime explicitly exposes which source-side and/or target-side attached child participated in the triggering contact.
4. When contact is parent-body only, no child identity is exposed.
5. Fight resolution, reproduction resolution, child payment, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Contact`
- `Interaction`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- interaction triggering still works exactly as it does today
- the same deterministic contact geometry is used as before
- snapshots expose the participating source-side and target-side child IDs only when attached-child contact is actually involved
- all broader contact and interaction semantics remain unchanged

This avoids turning contact into a larger event-history system while still making embodied contact legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `contact_origin` says whether attached-child contact was involved
- pair identity is already visible
- but the contract still does not say which attached child or children actually participated

Build should therefore make one minimal contract extension that exposes child-contact identity without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side contact detection path that can surface the participating child IDs
- deterministic tests that prove the exposed IDs match the contact geometry used for the triggering interaction
- client rendering or HUD output that remains sufficient to read the new contact-child identity

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second contact-event system
- reuse the same deterministic geometry and pair-selection rules already used by interaction triggering
- keep the contract addition minimal and explicit
- avoid adding speculative contact-history or replay structures

## Initial Test Plan

### Server tests

- source-child-to-target-parent contact exposes the source-side child ID only
- source-parent-to-target-child contact exposes the target-side child ID only
- child-to-child contact exposes both participating child IDs
- parent-body-only contact exposes no child-contact identity

### Contract tests

- the snapshot schema accepts the new contact-child identity fields

### Integration tests

- the client receives an interaction snapshot with attached-child contact identity when the trigger path uses attached-child geometry

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an interaction begins through attached-child geometry
2. the existing authoritative contact path selects the pair as usual
3. the snapshot still exposes the same pair and `contact_origin`
4. the snapshot now also explicitly identifies the participating source-side and/or target-side child IDs

## Done Criteria

- child-triggered interaction still resolves with the same current gameplay behavior
- attached-child contact exposes the participating child identity or identities
- tests prove the exposed IDs match the actual triggering contact geometry
- broader contact, fight, reproduction, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new contact precedence rules
- child-level combat or reproduction systems
- historical contact logs
- replay or event-stream redesign
