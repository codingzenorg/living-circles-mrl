# Slice: Initial Same-Shape Fight Resolution

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for current interaction classification and fight outcome state

## Architecture Mode

Explicit client/server boundary with the server resolving same-shape overlap into a deterministic fight outcome.

This slice extends the current interaction-classification behavior without introducing reproduction or lineage semantics.

## Discovery Scope

Establish the smallest deterministic fight loop for same-shape overlap:

- circles still carry explicit shape identity
- same-shape overlap is no longer only classified, it resolves as a fight
- the fight winner is chosen deterministically from existing state
- the losing circle is removed from the active world snapshot

This slice does **not** attempt to implement:

- different-shape reproduction outcomes
- child spawning or accumulation
- continuity or replacement after defeat
- randomness in combat
- damage-over-time or multi-tick combat
- multiplayer scoreboards

## Why This Slice Next

The current implementation can say that a same-shape overlap is a `fight_candidate`, but it cannot yet do anything with that fact. The extracted model explicitly says same shape should lead toward fight, and fight is simpler to resolve next than reproduction because it does not require children, lineage, or continuity semantics.

This slice keeps the change bounded by:

- resolving only same-shape interaction
- leaving different-shape interaction as classification-only for now
- using deterministic winner selection from already available circle state

## Use-Case Contract

### Use Case

`ResolveSameShapeFight`

### Primary Actor

Player controlling one circle while one autonomous circle participates in the same world.

### Pre-conditions

- a server process can host one bounded world
- the world contains one player-controlled circle and one autonomous circle
- both circles have explicit shape identity
- the circles overlap
- the overlap is classified as `fight_candidate`

### Trigger

The server advances a tick in which two same-shape circles overlap.

### Success Outcome

- the server resolves the fight deterministically on that tick
- one circle remains active in the world
- the losing circle is removed from later snapshots
- the fight outcome is visible in the snapshot contract

### Failure Or Rejection Cases

- if circles do not overlap, no fight resolution occurs
- if shapes differ, this slice must not resolve the interaction as a fight
- if the deterministic winner rule produces a tie, the slice must define one explicit tie-breaker rather than leaving the outcome ambiguous

## Main Business Rules

1. Only same-shape overlap resolves as a fight in this slice.
2. Different-shape overlap remains `reproduce_candidate` only; it does not resolve yet.
3. Fight resolution happens authoritatively on the server.
4. The winner is chosen deterministically from explicit state.
5. The losing circle is removed from active world participation.
6. The fight outcome is reflected in later snapshots.
7. This slice does not introduce child creation, replacement, or continuity after defeat.

## Minimal Domain Concepts In Scope

- `Shape`
- `Fight Candidate`
- `Fight Resolution`
- `Winner`
- `Loser`
- `Active Circle`
- `World Snapshot`

## Deterministic Winner Rule

The build step must choose one explicit deterministic fight rule and test it clearly.

Preferred priority order:

1. higher energy wins
2. if energy is equal, larger radius wins
3. if still equal, player circle wins only if that tie-breaker is declared explicitly in the implementation artifact; otherwise choose a fixed lexical ID tie-breaker

The chosen rule must remain simple, inspectable, and fully deterministic.

## Required Runtime Contract Changes

The snapshot contract should expose enough information for the client to understand that a fight resolved, for example:

- the current active circles after resolution
- an optional current or last interaction outcome field indicating `fight_resolved`
- the identifiers of the winner and loser if an outcome field is included

The client-to-server `movement_intent` contract does not need to change.

## Required Ports Or Boundaries

- server-side deterministic fight resolution for same-shape overlap
- server-side world update that removes the losing circle from future participation
- shared contract definition for a resolved interaction outcome
- client-side rendering boundary that reflects when only one circle remains
- deterministic tests covering winner selection and post-fight snapshots

## Build Guidance

- keep the fight as a one-tick authoritative resolution, not a health system
- prefer removal of the losing circle over a more complex disabled state
- do not award children, score, growth, or inheritance yet
- keep the outcome visible in the snapshot contract rather than burying it in logs
- avoid introducing a generalized combat framework

## Initial Test Plan

### Server tests

- same-shape overlap resolves to one winner and one loser
- higher-energy circle wins
- deterministic tie-breaker applies when energies are equal
- different-shape overlap does not resolve as a fight
- the losing circle is absent from later snapshots

### Contract tests

- the snapshot schema includes explicit resolved-interaction fields when applicable
- the contract still expresses different-shape overlap separately from resolved same-shape fights

### Integration tests

- the client receives snapshots showing same-shape circles moving into overlap
- after overlap, a later snapshot shows only the winning circle still active
- the interaction outcome in the snapshot identifies the fight as resolved

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives an initial snapshot with one player circle and one autonomous circle
2. both circles have the same shape
3. the circles move into overlap
4. the server resolves the fight deterministically
5. the client receives a later snapshot showing the winner still active, the loser removed, and the fight outcome visible

## Done Criteria

- same-shape overlap resolves as a deterministic fight
- the server removes the losing circle from future snapshots
- different-shape overlap remains unresolved reproduction classification only
- the snapshot contract exposes the fight result clearly enough for the client to render or inspect
- tests cover winner selection, loser removal, and contract shape
- the slice does not implement reproduction, children, continuity, or growth effects

## Out Of Scope Follow-Ups

- different-shape reproduction resolution
- child accumulation and lineage
- growth and power changes after victory
- death replacement behavior
- multiple simultaneous fights
- persistent scoring or ranking
