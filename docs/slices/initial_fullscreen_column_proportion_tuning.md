# Slice: Initial Fullscreen Column Proportion Tuning

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The canvas is framed, the header and HUD are tighter, and the side column has clearer internal hierarchy. The remaining presentation pressure is the fixed relationship between the play stage and the support column: the current desktop split can still feel slightly rigid rather than intentionally proportioned around a dominant stage and a genuinely secondary support rail.

## Discovery Scope

Establish the smallest client-facing adjustment to fullscreen column proportions:

- the play stage should feel more dominant
- the support column should remain readable but more clearly secondary
- the layout should stay usable across common desktop widths
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- new support panels
- responsive breakpoint redesign from scratch
- broader redesign beyond the desktop proportion between play stage and support column

## Why This Slice Next

Recent fullscreen-oriented slices already:

- expanded the demo to fullscreen
- reduced header, legend, and HUD footprint
- reduced side-column visual weight
- improved play-stage framing
- improved side-column internal hierarchy

That means the next pressure is not support content, but spatial proportion. The page should communicate a clearer large-stage / small-support relationship through width allocation.

## Use-Case Contract

### Use Case

`RenderBalancedFullscreenColumns`

### Primary Actor

The player viewing the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders the fullscreen layout
- the runtime contract remains sufficient and unchanged
- both play stage and support column remain present

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the play stage feels more dominant through proportion
- the support column remains readable and usable
- the desktop layout feels less rigid

### Failure Or Rejection Cases

- if the support column becomes cramped or unreadable, the slice fails
- if the stage does not feel more dominant, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may tune column proportions without altering those semantics.
3. The play stage should remain the dominant visual surface.
4. The support column should remain readable and usable.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Column Proportion`
- `Stage Dominance`
- `Secondary Support Rail`
- `Desktop Balance`

## Bounded Interpretation

This slice chooses the smallest useful proportion adjustment:

- tune the fullscreen desktop split between play stage and support column
- preserve the current information layout and hierarchy
- avoid adding new systems or redesigning responsive behavior from scratch

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes fullscreen column proportions

## Build Guidance

- prefer modest proportional tuning over large breakpoint changes
- keep the play stage clearly dominant
- avoid making the support column feel cramped

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the play stage feels more dominant on desktop
- the support column remains readable and usable
- the overall fullscreen layout feels less rigid

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player observes the proportional balance between the play stage and support column
3. the play stage feels more dominant while the support rail stays usable

## Done Criteria

- fullscreen column proportions feel better balanced
- the play stage is more clearly dominant
- the support column remains readable
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- new support systems
- responsive redesign from scratch
- camera or zoom work
- server-side layout metadata
