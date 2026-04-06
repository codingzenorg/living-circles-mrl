# Slice: Initial Header Footprint Reduction In Fullscreen Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The legend has been collapsed and the side column is lighter, but the page still opens with a relatively tall title and introductory paragraph above the play area. In the fullscreen layout, that header now claims more vertical space than it needs for ordinary use.

## Discovery Scope

Establish the smallest client-facing reduction in header footprint:

- the top area should take less vertical space
- the canvas should begin higher on the page
- the project identity should remain clear
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- branding redesign
- onboarding or tutorial systems
- broader page restructuring beyond the top header footprint

## Why This Slice Next

Recent fullscreen-oriented slices already:

- moved the support area into a side column
- reduced side-column visual weight
- collapsed the legend into a compact line

That means the next obvious interruption of the play surface is the tall header block itself. The canvas can start higher if the title and intro copy become more compact.

## Use-Case Contract

### Use Case

`RenderLowFootprintHeader`

### Primary Actor

The player opening the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders the fullscreen layout
- the runtime contract remains sufficient and unchanged
- project identity is already present in the page header

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the top header occupies less space
- the canvas begins higher and feels more immediate
- the page still clearly identifies the demo

### Failure Or Rejection Cases

- if the page loses clear identity, the slice fails
- if the change adds new onboarding or tutorial systems, scope is exceeded
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reduce header footprint without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. The page should still clearly identify the demo.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Header Footprint`
- `Canvas Immediacy`
- `Project Identity`
- `UI Restraint`

## Bounded Interpretation

This slice chooses the smallest useful header reduction:

- tighten title and intro presentation
- preserve clear identity with less vertical space
- avoid adding new explanatory systems

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes header presentation

## Build Guidance

- prefer a shorter, sharper header over decorative additions
- keep the page understandable at a glance
- avoid shifting explanatory burden back into the legend or support column

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the top header occupies less space
- the canvas starts higher in the viewport
- the page still clearly identifies the demo

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player scans the top of the page
3. the header reads clearly but takes less vertical space before the play surface begins

## Done Criteria

- header footprint is reduced
- the canvas begins higher
- project identity remains clear
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- branding redesign
- tutorials
- help overlays
- server-side cue metadata
