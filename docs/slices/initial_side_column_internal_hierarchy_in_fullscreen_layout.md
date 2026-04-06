# Slice: Initial Side Column Internal Hierarchy In Fullscreen Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The side column is already lighter and the play stage is better framed, but the three support blocks still read with similar visual weight inside the column. The player block should feel clearly primary, with NPCs secondary and recent encounters tertiary.

## Discovery Scope

Establish the smallest client-facing improvement to side-column internal hierarchy:

- player information should feel clearly primary within the side column
- NPC summaries should feel secondary
- recent encounters should feel tertiary
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- new support panels
- collapsible systems
- broader redesign beyond hierarchy inside the existing side column

## Why This Slice Next

Recent fullscreen-oriented slices already:

- expanded the demo layout
- reduced header, legend, and HUD footprint
- reduced side-column visual weight
- strengthened play-stage framing

That means the next pressure is not the existence of the side column, but how its internal blocks compete with one another. The side column should communicate a stronger priority order without adding more chrome.

## Use-Case Contract

### Use Case

`RenderHierarchicalSupportColumn`

### Primary Actor

The player scanning the side column during fullscreen play.

### Pre-conditions

- the browser client already renders the fullscreen layout
- the runtime contract remains sufficient and unchanged
- the side column already includes player, NPC, and recent encounter panels

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the player panel reads as primary
- NPC summaries read as secondary
- recent encounters read as tertiary
- the side column remains readable and secondary to the canvas

### Failure Or Rejection Cases

- if hierarchy becomes less clear, the slice fails
- if the side column regains too much visual weight, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may strengthen side-column hierarchy without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. The player panel should remain the clearest support priority.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Support Hierarchy`
- `Primary Player Panel`
- `Secondary NPC Panel`
- `Tertiary Encounter Panel`

## Bounded Interpretation

This slice chooses the smallest useful hierarchy improvement:

- adjust spacing, typography, or panel emphasis within the side column
- preserve the current information set and order
- avoid adding new systems or chrome-heavy treatments

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes side-column hierarchy

## Build Guidance

- prefer subtle hierarchy cues over decoration
- keep the player panel clearly strongest inside the support column
- avoid increasing overall side-column weight while improving internal order

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the player panel is visually primary within the column
- NPCs read as clearly below the player in importance
- recent encounters read as the lightest support block

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player scans the support column
3. the internal support hierarchy reads more clearly without growing the column's overall weight

## Done Criteria

- side-column internal hierarchy is clearer
- the player panel is primary
- NPC and encounter panels are clearly secondary and tertiary
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- new support panels
- collapsible systems
- server-side prioritization metadata
- broad page redesign
