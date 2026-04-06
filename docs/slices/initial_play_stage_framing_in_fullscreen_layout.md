# Slice: Initial Play Stage Framing In Fullscreen Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current fullscreen demo refinement track. The header, legend, and HUD are now tighter, and the support column is visually lighter. The remaining presentation pressure is that the play surface still reads mostly as a plain rectangular canvas dropped into the page, rather than as a clearly framed main stage distinct from surrounding support chrome.

## Discovery Scope

Establish the smallest client-facing framing improvement for the fullscreen play surface:

- the main play area should feel more clearly staged
- the canvas should remain the dominant visual surface
- the framing should not reintroduce heavy chrome
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- camera systems
- zoom systems
- animated transitions
- broad visual redesign beyond play-surface framing

## Why This Slice Next

Recent fullscreen-oriented slices already:

- expanded the demo to a fullscreen layout
- reduced side-column visual weight
- collapsed the legend
- reduced header footprint
- reduced HUD footprint

That means the next pressure is not more UI subtraction, but improving how the canvas is visually staged as the main surface within the now-lighter surrounding layout.

## Use-Case Contract

### Use Case

`RenderFramedPlayStage`

### Primary Actor

The player viewing the fullscreen demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders the fullscreen layout
- the canvas is already the dominant world surface
- the runtime contract remains sufficient and unchanged

### Trigger

The page is rendered in the fullscreen demo layout.

### Success Outcome

- the play surface feels more clearly framed
- the canvas remains dominant without becoming heavier
- surrounding support UI remains secondary

### Failure Or Rejection Cases

- if framing adds heavy decorative chrome, the slice fails
- if the support area becomes more dominant again, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may strengthen play-stage framing without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. Surrounding support UI should remain secondary.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Play Stage`
- `Stage Framing`
- `Canvas Primacy`
- `Secondary Support`

## Bounded Interpretation

This slice chooses the smallest useful framing improvement:

- strengthen the separation between play surface and surrounding page
- preserve the current information layout
- avoid adding new systems or decorative overload

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes play-surface framing

## Build Guidance

- prefer subtle but intentional framing over ornament
- keep the canvas as the visual anchor
- avoid increasing support-column visual weight while improving the stage

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the play surface feels more clearly staged
- the canvas remains visually dominant
- surrounding UI still reads as secondary

## Scenario Definition

Start a local server and open one browser client on a desktop-sized viewport.

Scenario steps:

1. the page renders in the fullscreen layout
2. the player looks at the relationship between the canvas and surrounding UI
3. the canvas reads more clearly as the main stage without adding heavy chrome

## Done Criteria

- play-stage framing is stronger
- the canvas remains dominant
- surrounding UI remains secondary
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- camera systems
- animations
- zoom controls
- broader page redesign
