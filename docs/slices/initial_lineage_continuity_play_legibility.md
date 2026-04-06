# Slice: Initial Lineage Continuity Play Legibility

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is unless one minimal readability field is clearly necessary

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current play-legibility direction, but shifts from motion and local pressure readability to continuity readability. The simulation already expresses lineage, generations, attached children, promoted-child continuity, and child-based survival, yet much of that meaning still reads more like instrumentation than live world meaning.

## Discovery Scope

Establish the smallest client-facing readability improvement that helps a player perceive continuity and lineage directly from the world view:

- attached children should read more clearly as continuity reserves rather than only as extra decoration
- lineage continuity through promotion or survival should become easier to notice during ordinary play
- the cue should stay grounded in already-authoritative lineage and child state
- the change should remain presentation-focused rather than adding new server rules

This slice does **not** attempt to implement:

- new lineage mechanics
- inherited variation
- family trees
- persistent lineage scoring
- new reproduction or death rules
- new server-side continuity logic

## Why This Slice Next

Recent slices established:

- stronger shape and interaction legibility
- stronger crowding-pressure legibility
- stronger food-pressure legibility
- stronger autonomy-intent legibility
- stronger player-motion legibility

That means the running world is becoming easier to read in the short term, but one of the model’s most distinctive semantic areas still remains weakly legible during ordinary play: continuity. Children, lineage IDs, and promotion outcomes exist, but their meaning still relies too heavily on labels and event logs rather than on the living scene itself.

The next pressure is therefore to make lineage continuity easier to perceive without expanding lineage mechanics or turning the client into a genealogy tool.

## Use-Case Contract

### Use Case

`RenderImmediateLineageContinuityMeaning`

### Primary Actor

The player observing and steering inside the running world.

### Pre-conditions

- the authoritative server already exposes lineage ID, generation, and attached-child state
- the browser client already renders the player, autonomous circles, and attached children
- continuity outcomes already appear through the current event log and snapshot state

### Trigger

A world snapshot is rendered on the client.

### Success Outcome

- attached children read more clearly as continuity-bearing reserves
- promoted or ongoing lineage becomes easier to notice during ordinary play
- the player can infer more of the system’s continuity meaning from the world view itself

### Failure Or Rejection Cases

- if the slice only adds more text, the legibility goal is missed
- if the client invents lineage semantics beyond bounded use of authoritative state, scope is exceeded
- if the visual treatment becomes noisy or overwhelms the existing danger/food/crowding/motion cues, readability worsens

## Main Business Rules

1. The server remains authoritative for all lineage and continuity semantics.
2. The client may use existing authoritative snapshot data to improve visual legibility.
3. The slice should prefer direct visual cues over added textual explanation.
4. Continuity-bearing children and lineage persistence should become easier to notice.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Lineage`
- `Continuity`
- `Attached Child`
- `Promoted Survival`
- `Readable Continuity Cue`

## Bounded Legibility Interpretation

This slice chooses the smallest useful interpretation of continuity readability:

- use local visual treatments that connect a parent and its attached children more clearly
- use bounded continuity emphasis when a lineage continues through ordinary play
- keep the cues tied to current authoritative lineage and child state rather than wider historical tooling

This preserves the current immediate-play presentation while making one of the project’s distinguishing ideas more visible.

## Required Runtime Contract Changes

The current contract is likely sufficient because the client already receives:

- lineage IDs
- generation values
- attached children
- continuity-related interaction outcomes

Build should avoid extending the contract unless one tiny authoritative readability field is clearly necessary.

## Required Ports Or Boundaries

- client rendering logic
- deterministic client-side tests only if the repo already has a practical path for them; otherwise rely on existing integration coverage and manual verification notes
- implementation notes that record the chosen continuity-legibility cues

## Build Guidance

- prefer one or two strong continuity cues over a broad annotation system
- keep the treatment local to parent-child groupings and recent continuity events
- stay visually compatible with the current dark-mode canvas language
- avoid turning lineage into a dense debug overlay

## Initial Test Plan

### Integration tests

- existing snapshots remain sufficient for the new rendering
- no contract change is required unless build documents why it became necessary

### Manual verification

- attached children should feel more visibly tied to continuity rather than only to clutter
- continuity through survival or promotion should be easier to notice when it happens
- the cue should not drown out danger, food, crowding, motion, or nearby intent cues

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player enters the larger world with lineages, attached children, and live interaction
2. circles move, collide, consume food, and sometimes preserve continuity through attached children
3. the client renders stronger local continuity cues tied to existing authoritative state
4. the player can more easily read lineage persistence from the world itself

## Done Criteria

- the client presents clearer lineage-continuity meaning
- the world remains authoritative
- the slice improves ordinary play readability without a broad UI rewrite
- existing server semantics remain unchanged

## Out Of Scope Follow-Ups

- family trees
- explicit ancestry views
- persistent lineage scoring
- inherited variation
- new server-side continuity mechanics
