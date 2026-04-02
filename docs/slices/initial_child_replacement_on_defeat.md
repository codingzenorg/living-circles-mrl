# Slice: Initial Child Replacement On Defeat

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for post-defeat continuity outcomes

## Architecture Mode

Explicit client/server boundary with the server resolving one bounded continuity rule after defeat.

This slice extends the current child accumulation and growth model without introducing full lineage trees, inheritance variability, or spawned child entities.

## Discovery Scope

Establish the smallest deterministic continuity rule supported by the source material:

- when a circle loses a fight and has no children, it still disappears
- when a circle loses a fight and has at least one child, one child replaces it immediately
- the replacement remains one active circle in the world, not a new swarm
- the replacement outcome is visible in later snapshots

This slice does **not** attempt to implement:

- multiple replacement candidates
- separate child entities existing before replacement
- inherited mutations or mixed traits
- replacement after non-fight death causes
- lineage trees, ancestry history, or family IDs
- continuity rules for reproduction outside defeat

## Why This Slice Next

The current implementation now supports:

- child accumulation through reproduction
- visible growth from accumulated children
- deterministic defeat through same-shape fights

But death is still purely terminal. The source material explicitly says a dead circle may disappear or be replaced by a child. That makes continuity the next highest-value missing behavior because it:

- gives children a second consequence beyond raw size
- closes the loop from reproduction to continuity
- creates risk in fights without inventing a new scoring system
- stays aligned with the extracted model more directly than broader balancing work

## Use-Case Contract

### Use Case

`ReplaceDefeatedCircleWithChild`

### Primary Actor

Player controlling one circle while autonomous circles participate under the same continuity rule.

### Pre-conditions

- a server process can host one bounded world
- two same-shape circles can resolve into a deterministic fight
- the losing circle may have zero or more accumulated children

### Trigger

The server advances a tick in which a same-shape fight is resolved and one participant loses.

### Success Outcome

- if the loser has zero children, it disappears as before
- if the loser has one or more children, it remains represented in the world through one replacement circle
- replacement consumes exactly one child count from the defeated lineage
- later snapshots expose the surviving replacement state clearly enough for the client to inspect

### Failure Or Rejection Cases

- if the loser has no children, this slice must not invent a replacement
- if replacement requires creating more than one active circle, this slice is too broad
- if replacement needs unspecified inheritance logic to function, the build should choose the smallest deterministic defaults and record them explicitly

## Main Business Rules

1. Same-shape fights still resolve deterministically first.
2. Defeat does not always mean disappearance.
3. A defeated circle with at least one child transfers continuity through exactly one replacement circle.
4. Replacement consumes one child from the defeated lineage.
5. The replacement occupies the defeated circle's position in the active world.
6. This slice must apply the same rule to player and autonomous circles.
7. This slice does not yet model ancestry history, genetic inheritance, or continuity across multiple generations.

## Minimal Domain Concepts In Scope

- `Defeat`
- `Child Replacement`
- `Continuity`
- `Active Circle`
- `Consumed Child`
- `World Snapshot`

## Bounded Continuity Interpretation

The source material says a child may replace a dead circle, but it does not define the exact inheritance model.

This slice therefore chooses the smallest inspectable interpretation:

- replacement preserves the defeated circle's identity as one continuing active circle
- one child is consumed during replacement
- shape remains the defeated circle's shape
- replacement starts with the defeated circle's current position and a deterministic baseline energy state chosen during build
- any remaining child count stays on the replacement circle

This interpretation keeps continuity explicit without requiring a separate lineage subsystem yet.

## Required Runtime Contract Changes

The existing snapshot contract may already be sufficient if replacement is represented as the continued presence of an active circle after defeat.

This slice should ensure:

- snapshots make it clear when the loser did not disappear
- the current interaction outcome still identifies the fight winner and loser
- the replacement circle state is visible through normal snapshot fields

An additional continuity flag may be added only if the build cannot keep the behavior inspectable otherwise.

## Required Ports Or Boundaries

- server-side defeat handling that can either remove or replace the loser
- server-side rule for consuming one child during replacement
- client-side rendering that can show the continued active circle after defeat
- deterministic tests for both disappearance and replacement cases

## Build Guidance

- keep replacement inside the existing fight-resolution flow instead of adding a generalized respawn system
- prefer one deterministic replacement rule over configurable continuity options
- avoid adding persistent history tracking yet
- if a baseline energy value is needed for the replacement, choose one explicit rule and record it in the implementation artifact
- keep the result inspectable in snapshots and tests rather than in logs only

## Initial Test Plan

### Server tests

- a loser with zero children still disappears
- a loser with at least one child remains active through replacement
- replacement consumes exactly one child count
- replacement preserves active participation in later snapshots
- the same continuity rule works for player and autonomous losers

### Contract tests

- the existing snapshot schema remains sufficient unless one explicit continuity field is introduced
- fight outcomes still identify winner and loser even when replacement occurs

### Integration tests

- the client receives a fight resolution where the loser remains present due to child replacement
- later snapshots show the consumed child count and continued active circle

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and receives a snapshot where one circle already carries at least one child
2. a same-shape fight occurs
3. the childless loser disappears, or the losing lineage with children remains through replacement
4. the client receives a later snapshot showing the continuity outcome clearly

## Done Criteria

- defeat can now lead to either disappearance or child replacement
- replacement consumes exactly one child
- player and autonomous circles follow the same continuity rule
- snapshots and tests make the continuity outcome inspectable
- the slice does not implement full lineage trees, mutations, or separate child entities

## Out Of Scope Follow-Ups

- persistent lineage identifiers
- inheritance of mixed or mutated traits
- replacement after non-fight deaths
- multiple children splitting into multiple active circles
- historical ancestry views in the client
