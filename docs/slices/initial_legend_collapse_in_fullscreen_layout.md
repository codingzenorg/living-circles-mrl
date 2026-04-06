# Slice: Initial Legend Collapse In Fullscreen Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The legend has already been reduced substantially, and the side column is now lighter, but the legend still occupies a full-width band above the play surface. In the fullscreen layout, that band increasingly feels redundant because the canvas, player panel, NPC panel, and encounter log already carry most of the meaning.

## Discovery Scope

Establish the smallest client-facing legend reduction that fits the fullscreen layout:

- the legend should claim less visual space
- the canvas should remain the dominant surface
- core cue meaning should remain recoverable
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- tutorials or onboarding systems
- hiding all cue explanation entirely
- broader redesign beyond the legend's footprint in the fullscreen layout

## Why This Slice Next

Recent slices already:

- reduced legend density
- simplified support text
- bounded support growth
- moved the support area into a fullscreen side column
- reduced the side column's visual weight

That means the next presentation pressure is no longer the support column itself. It is the wide legend strip that still sits above the play surface and competes with the canvas in the fullscreen layout.

## Use-Case Contract

### Use Case

`RenderLowFootprintLegend`

### Primary Actor

The player observing the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders the fullscreen layout
- support panels remain readable
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the legend takes less space and visual attention
- the canvas remains more dominant
- cue meaning remains recoverable without adding new systems

### Failure Or Rejection Cases

- if cue meaning becomes too hard to recover, the slice fails
- if a new tutorial or onboarding system is added instead of shrinking the legend, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reduce legend footprint without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. Core cue meaning should remain recoverable.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Legend Footprint`
- `Canvas Primacy`
- `Recoverable Meaning`
- `UI Restraint`

## Bounded Interpretation

This slice chooses the smallest useful legend reduction:

- reduce the legend's footprint or prominence in the fullscreen layout
- preserve enough explanation for the strongest cue families
- avoid adding replacement systems

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes legend behavior

## Build Guidance

- prefer restraint and compactness over adding new controls
- keep the legend available, but less dominant
- avoid reintroducing support-density growth elsewhere in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the legend occupies less space in the fullscreen layout
- the canvas feels less interrupted by explanatory UI
- the strongest cue meanings are still recoverable

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player scans the top area before focusing on the canvas
3. the legend now takes less space and attention while remaining useful

## Done Criteria

- legend footprint is reduced
- the canvas remains more dominant
- cue meaning remains recoverable
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- onboarding systems
- help overlays
- server-side cue metadata
- broader page redesign
