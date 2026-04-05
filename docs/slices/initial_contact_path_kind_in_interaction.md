# Slice: Initial Contact Path Kind In Interaction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where child-triggered interaction remains inspectable through explicit contact path kind

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for contact detection and interaction resolution.

This slice narrows one remaining contact inspectability gap by exposing whether an interaction was triggered by `child_to_parent` or `child_to_child` geometry when `contact_origin` is `attached_child`.

## Discovery Scope

Establish the smallest deterministic change that makes contact path kind explicit:

- when interaction contact begins through attached-child geometry, the same authoritative contact path should remain exactly as it is today
- the runtime should expose whether the triggering path was source-child-to-target-parent, source-parent-to-target-child, or child-to-child
- participating child identity, contact origin, winner selection, reproduction payment, continuity, feeding, movement, orbit, and steering remain unchanged

This slice does **not** attempt to implement:

- new interaction rules
- alternate contact precedence rules
- historical contact logs
- redesign of contact origin categories
- new client-side effects beyond what is necessary to read contact path kind

## Why This Slice Next

Recent slices made child-triggered interaction much more inspectable:

- `contact_origin` already distinguishes `attached_child` from `parent_body`
- child-triggered interaction now exposes the participating source-side and/or target-side child IDs
- reproduction outcomes now expose payment, creation, and ownership details

But contact triggering still stops one step short. The runtime can say which child IDs participated, yet it still does not say whether the actual trigger was child-to-parent or child-to-child. That weakens inspectability because the authoritative server already knows which geometric path fired first, while the client and tests still need to infer the path indirectly from which side has child IDs populated.

The next model pressure is therefore not to change contact semantics. It is to make the triggering contact path explicit in the same bounded way child identities are now explicit.

This slice is the narrowest next step because it:

- changes only contact inspectability
- preserves current contact precedence and interaction behavior
- avoids inventing event history or child-level combat systems
- gives build a deterministic contract extension instead of a larger semantic redesign

## Use-Case Contract

### Use Case

`ExposeContactPathKindDuringInteraction`

### Primary Actor

Any player or autonomous circle pair whose interaction begins through attached-child geometry.

### Pre-conditions

- contact detection already distinguishes `parent_body` from `attached_child`
- attached-child identity is already exposed for source and target sides
- attached-child-to-parent and attached-child-to-attached-child paths already exist

### Trigger

An interaction begins and `contact_origin` is `attached_child`.

### Success Outcome

- contact still resolves through the same authoritative path as before
- the authoritative interaction outcome explicitly exposes whether the triggering path was `source_child_to_target_parent`, `source_parent_to_target_child`, or `child_to_child`

### Failure Or Rejection Cases

- if `attached_child` contact still hides the trigger path kind, inspectability remains incomplete
- if exposing contact path kind changes interaction timing or outcome, scope is exceeded
- if the exposed path kind can vary non-deterministically, inspectability weakens

## Main Business Rules

1. Contact detection remains authoritative server-side behavior.
2. `contact_origin` remains the current coarse contact category.
3. When `contact_origin` is `attached_child`, the runtime explicitly exposes the triggering path kind.
4. When contact is parent-body only, no child contact path kind is exposed.
5. Participating child IDs remain unchanged and still reflect the current contact path.
6. Fight resolution, reproduction resolution, child payment, continuity, feeding, movement, orbit, and steering remain unchanged.

## Minimal Domain Concepts In Scope

- `Attached Child`
- `Contact`
- `Contact Path Kind`
- `Interaction`
- `World Snapshot`

## Bounded Interpretation

This slice chooses the smallest inspectable interpretation:

- interaction triggering still works exactly as it does today
- the same deterministic contact geometry and precedence are used as before
- snapshots expose the triggering path kind only when attached-child contact is actually involved
- all broader contact and interaction semantics remain unchanged

This avoids turning contact into a larger event-history system while still making embodied triggering legible.

## Required Runtime Contract Changes

The current contract is no longer fully sufficient because:

- `contact_origin` says that attached-child contact was involved
- child identity already says which participating child IDs were present
- but the contract still does not say whether the trigger path was child-to-parent or child-to-child

Build should therefore make one minimal contract extension that exposes contact path kind without redesigning the full snapshot shape.

## Required Ports Or Boundaries

- server-side contact detection path that can surface the chosen trigger path kind
- deterministic tests that prove the exposed path kind matches the geometry used for the triggering interaction
- client rendering or HUD output that remains sufficient to read the new contact path kind

## Build Guidance

- prefer extending the existing interaction payload rather than creating a second contact-event system
- reuse the same deterministic geometry and precedence rules already used by interaction triggering
- keep the contract addition minimal and explicit
- avoid adding speculative contact-history or replay structures

## Initial Test Plan

### Server tests

- source-child-to-target-parent contact exposes `source_child_to_target_parent`
- source-parent-to-target-child contact exposes `source_parent_to_target_child`
- child-to-child contact exposes `child_to_child`
- parent-body-only contact exposes no contact path kind

### Contract tests

- the snapshot schema accepts the new contact path kind field

### Integration tests

- the client receives an interaction snapshot with the correct contact path kind when attached-child geometry triggered the interaction

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. an interaction begins through attached-child geometry
2. the existing authoritative contact path selects the pair as usual
3. the snapshot still exposes the same pair, `contact_origin`, and child IDs
4. the snapshot now also explicitly identifies the triggering path kind

## Done Criteria

- child-triggered interaction still resolves with the same current gameplay behavior
- attached-child contact exposes the triggering path kind
- tests prove the exposed path kind matches the actual triggering geometry
- broader contact, fight, reproduction, and continuity behavior remain unchanged

## Out Of Scope Follow-Ups

- new contact precedence rules
- child-level combat or reproduction systems
- historical contact logs
- replay or event-stream redesign
