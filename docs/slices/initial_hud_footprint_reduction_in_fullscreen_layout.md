# Slice: Initial HUD Footprint Reduction In Fullscreen Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The header is now tighter, the legend has collapsed into a compact line, and the side column is lighter. The remaining top-of-page pressure now sits mainly in the HUD row, which still occupies a full-width control band above the play surface.

## Discovery Scope

Establish the smallest client-facing reduction in HUD footprint:

- the top HUD should take less visual and vertical space
- the canvas should remain the dominant visual surface
- the essential controls and status should remain readable
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- hidden controls
- settings or menu systems
- broader redesign beyond the HUD footprint

## Why This Slice Next

Recent fullscreen-oriented slices already:

- reduced side-column visual weight
- collapsed the legend
- reduced the header footprint

That means the next obvious interruption of the play surface is the HUD itself. The demo can feel more immediate if the status and reset controls become visually lighter while staying usable.

## Use-Case Contract

### Use Case

`RenderLowFootprintHUD`

### Primary Actor

The player using the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders the fullscreen layout
- the runtime contract remains sufficient and unchanged
- status and reset controls are already present in the HUD

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the HUD occupies less space and attention
- the canvas begins to feel more immediate
- core status and reset remain clear and usable

### Failure Or Rejection Cases

- if status or reset becomes hard to use, the slice fails
- if the change adds new systems instead of reducing HUD footprint, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reduce HUD footprint without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. Core status and reset affordances should remain available.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `HUD Footprint`
- `Canvas Immediacy`
- `Readable Status`
- `Usable Reset`

## Bounded Interpretation

This slice chooses the smallest useful HUD reduction:

- tighten or simplify the status row
- preserve the essential control and signal set
- avoid adding replacement systems

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes HUD presentation

## Build Guidance

- prefer concise status presentation over decorative additions
- keep reset easy to find
- avoid pushing explanatory burden back into the legend or side column

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the HUD occupies less space in the fullscreen layout
- the canvas feels more immediate
- reset and core status remain readable and usable

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player scans the top controls and status row
3. the HUD reads more lightly while remaining usable

## Done Criteria

- HUD footprint is reduced
- the canvas remains dominant
- status and reset remain clear
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- settings menus
- hidden controls
- onboarding systems
- server-side status metadata
