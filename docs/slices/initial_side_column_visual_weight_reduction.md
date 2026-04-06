# Slice: Initial Side Column Visual Weight Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo track. The canvas now uses more of the viewport and the support area has moved into a persistent desktop side column. That improved spatial usage, but it also made the support column feel heavier because three bordered panels now sit in a visually continuous stack beside the play surface.

## Discovery Scope

Establish the smallest client-facing reduction in side-column visual weight:

- the side column should feel lighter and more secondary
- the canvas should remain the dominant visual surface
- support readability should be preserved
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- new support panels
- collapsing systems
- major restyling of the entire page beyond the side column's visual weight

## Why This Slice Next

Recent slices already:

- simplified canvas labels
- moved player details outside the canvas
- reduced support text density
- bounded lower support growth
- expanded the demo into a fullscreen layout with a desktop side column

That means the next pressure is no longer page usage. It is the visual mass of the side column now that it occupies a stable role beside the canvas.

## Use-Case Contract

### Use Case

`RenderLightweightSupportColumn`

### Primary Actor

The player observing the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders a desktop side column beside the canvas
- the server remains authoritative
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the side column feels lighter and more secondary
- the canvas remains dominant
- support information stays readable

### Failure Or Rejection Cases

- if support readability is harmed, the slice fails
- if the canvas loses visual priority, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reduce support-column visual weight without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. Support information should remain readable and secondary.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Side Column Weight`
- `Canvas Primacy`
- `Secondary Support`
- `Readable Restraint`

## Bounded Interpretation

This slice chooses the smallest useful visual-weight reduction:

- soften support-panel chrome and density where possible
- preserve current information hierarchy
- avoid introducing new systems or content

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes side-column presentation

## Build Guidance

- prefer restraint over decorative additions
- keep the player card clearly primary within the side column
- avoid reintroducing dense legend or support explanation in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the side column feels visually lighter than before
- the canvas remains clearly dominant
- support information remains readable at a glance

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player observes the side column beside the canvas
3. the side column reads as lighter-weight while preserving support readability

## Done Criteria

- side-column visual weight is reduced
- the canvas remains dominant
- support information remains readable
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- new support systems
- collapsible panels
- server-side legibility fields
- whole-page visual redesign
