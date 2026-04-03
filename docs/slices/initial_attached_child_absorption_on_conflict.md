# Slice: Initial Attached Child Absorption On Conflict

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible attached-child loss

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether overlap removes an attached child before removing a whole parent circle.

This slice keeps orbiting children attached and visible, but makes them vulnerable during hostile interactions. It does not replace the full current fight system or remove transitional radius semantics yet.

## Discovery Scope

Establish the smallest deterministic rule that makes attached children directly absorb hostile loss:

- in same-shape conflict, a circle that would otherwise lose may first lose one attached child instead of immediately losing the parent
- the loss of that child is visible in later snapshots
- the current parent-level fight and continuity rules remain in place after the absorbed loss
- the same visible child removal can still support the existing reproduction-payment and replacement-consumption rules

This slice does **not** attempt to implement:

- free child detachment
- child-vs-child targeting independent of the parent
- area or radius no longer depending on child count
- shape-specific combat trees
- probabilistic damage
- a generalized hit-point system

## Why This Slice Next

The current implementation made children visible as attached orbiters and already lets them be consumed abstractly for reproduction payment and replacement continuity. But hostile interaction still resolves almost entirely at the parent level.

The model pressure is now:

- visible child bodies should matter when conflict happens
- children should be loss-bearing before a whole parent disappears when that matches the current continuity idea
- the system should keep a direct bridge between visible children and existing `children_count`-based rules

This slice is the narrowest next step because it:

- gives orbiting children a direct survival role without replacing the whole combat model
- makes hostile interaction more legible in the demo
- stays deterministic and cheap
- preserves current fight winner selection while changing what the loss lands on first

## Use-Case Contract

### Use Case

`AbsorbConflictLossThroughAttachedChild`

### Primary Actor

Two same-shape circles entering a hostile authoritative overlap.

### Pre-conditions

- same-shape overlap currently resolves through the deterministic fight ordering
- a circle may own zero or more attached orbiting children
- snapshots already expose attached child state explicitly

### Trigger

The server resolves a same-shape hostile overlap during a simulation tick.

### Success Outcome

- the deterministic fight winner is still chosen
- if the loser owns at least one attached child, one attached child is removed first
- the losing parent remains active for now if the slice chooses that absorbed-loss path
- later snapshots show the reduced attached-child count visibly

### Failure Or Rejection Cases

- if hostile overlap still removes only abstract counters and not visible children, the slice fails its purpose
- if child loss is nondeterministic or visually opaque, inspectability is weakened
- if the player gets a hidden exemption from child loss, fairness is weakened

## Main Business Rules

1. Hostile conflict remains authoritative server-side behavior.
2. The current deterministic fight winner selection remains unchanged in this slice.
3. If the loser has at least one attached child, one attached child may absorb the conflict loss before the parent is removed or replaced.
4. Attached-child loss must stay visible and synchronized with `children_count`.
5. Player and autonomous circles follow the same child-loss rule.
6. This slice does not remove the current reproduction-payment or replacement-consumption semantics.
7. Radius growth and child-based fight power may remain transitional semantics for now.

## Minimal Domain Concepts In Scope

- `Hostile Overlap`
- `Fight Winner`
- `Absorbed Loss`
- `Attached Child Removal`
- `World Snapshot`

## Bounded Conflict Interpretation

This slice chooses the smallest inspectable interpretation:

- determine the fight winner using the current rule
- if the loser has at least one attached child, remove exactly one attached child
- keep the losing parent active on that tick after the child is removed
- only when no attached child exists does the current parent-level defeat path continue unchanged

This keeps hostile conflict simple while making visible children materially protective.

## Required Runtime Contract Changes

The current snapshot contract is structurally sufficient if attached-child loss is represented through:

- reduced `children_count`
- reduced `attached_children`
- the existing interaction object

Build should extend the contract only if one small explicit interaction outcome is needed to distinguish child absorption from parent defeat.

## Required Ports Or Boundaries

- server-side hostile-overlap resolution that can remove one attached child
- synchronized count and visible child updates
- client-side rendering of reduced orbiting children after conflict
- deterministic tests covering child-loss and fallback to current defeat behavior

## Build Guidance

- prefer reusing the current winner-selection rule
- add one explicit child-absorption step before full loser removal
- keep visible child removal synchronized with the authoritative `children_count`
- do not rewrite reproduction or detachment rules in this slice
- preserve inspectability in the demo and tests

## Initial Test Plan

### Server tests

- a losing circle with at least one attached child loses one child and remains active
- a losing circle with no attached children still follows the current defeat path
- visible attached-child count and `children_count` stay synchronized after conflict absorption
- player and autonomous circles both follow the same absorbed-loss rule

### Contract tests

- the current schema remains sufficient unless a small explicit conflict outcome field is added

### Integration tests

- the client receives a same-shape conflict snapshot after which one orbiting child is visibly gone
- the player can survive one hostile conflict through child absorption in the demo

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two same-shape circles enter hostile overlap
2. the deterministic winner is chosen as before
3. the loser loses one attached child instead of disappearing immediately when a child is available
4. later snapshots show one fewer orbiting child and continued parent participation

## Done Criteria

- same-shape conflict can remove an attached child visibly before removing a parent
- attached-child loss stays synchronized with `children_count`
- fallback to current parent-level defeat still works when no child is available
- player and autonomous circles follow the same rule
- tests cover absorbed loss and fallback defeat

## Out Of Scope Follow-Ups

- child-vs-child collision rules
- removal of radius growth as a transitional shortcut
- detachment after conflict
- explicit damage accumulation over multiple ticks
- special fight powers by shape or lineage
