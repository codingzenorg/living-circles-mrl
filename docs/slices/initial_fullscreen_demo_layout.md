# Slice: Initial Fullscreen Demo Layout

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract used as-is

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world truth while the client improves how that truth is legibly presented.

This slice continues the current UI refinement track, but shifts from bounded support growth to overall page usage. The current demo keeps a constrained central layout that helped during early iteration, yet the world is now larger and the simulation is visually denser. That makes the limited page width feel artificially narrow for ordinary play and observation.

## Discovery Scope

Establish the smallest client-facing full-screen layout that expands the demo surface without changing semantics:

- the canvas should make fuller use of the available viewport
- support panels should remain readable without overtaking the canvas
- the layout should still work on narrower screens
- the change should stay in presentation and layout, not world semantics

This slice does **not** attempt to implement:

- new game mechanics
- new server fields
- camera systems
- zoom systems
- fullscreen browser API integration
- broader design-system rework beyond layout usage

## Why This Slice Next

Recent slices already:

- reduced legend density
- reduced support text density
- stabilized the player card position
- bounded lower support-panel growth

Those changes mean the support area is now better behaved. The next pressure is no longer support sprawl, but the opposite: the demo no longer uses enough of the available viewport for the larger world it now renders.

## Use-Case Contract

### Use Case

`RenderFullscreenDemoLayout`

### Primary Actor

The player observing and controlling the browser demo on a desktop-sized screen.

### Pre-conditions

- the browser client already renders a canvas plus support area
- the server remains authoritative
- the runtime contract remains sufficient and unchanged

### Trigger

The page is opened during ordinary play.

### Success Outcome

- the demo uses substantially more of the available screen
- the canvas becomes more spatially generous
- support information remains readable and secondary

### Failure Or Rejection Cases

- if support panels become dominant again, the slice fails
- if the layout becomes unusable on narrower screens, the slice fails
- if world semantics or contract shape changes, scope is exceeded

## Main Business Rules

1. The server remains authoritative for all world semantics.
2. The client may reorganize screen usage without altering those semantics.
3. The canvas should remain the dominant visual surface.
4. Support panels should remain available but secondary.
5. Existing world mechanics remain unchanged.

## Minimal Domain Concepts In Scope

- `Fullscreen Layout`
- `Canvas Primacy`
- `Secondary Support Area`
- `Responsive Presentation`

## Bounded Interpretation

This slice chooses the smallest useful full-screen shift:

- expand the main container and canvas usage across the viewport
- reorganize support presentation only as needed to remain readable
- preserve the current information set and hierarchy

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- client presentation only
- implementation notes only if the slice materially changes layout usage

## Build Guidance

- prefer layout and sizing changes over adding new UI systems
- keep desktop play more spacious while preserving a reasonable mobile fallback
- avoid reintroducing dense explanatory chrome in the same slice

## Initial Test Plan

### Validation

- existing tests should remain sufficient because no contract or server behavior changes are expected

### Manual verification

- the page uses substantially more of the viewport on desktop
- the canvas feels more spacious
- the support area remains readable and secondary
- narrow screens still collapse to a usable layout

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the player opens the demo on a desktop-sized viewport
2. the layout renders using much more of the available screen width and height
3. the player can still read support information without it overtaking the play surface

## Done Criteria

- the demo uses more of the available screen
- the canvas is more spacious
- support information remains secondary and readable
- existing world semantics remain unchanged

## Out Of Scope Follow-Ups

- camera zoom
- fullscreen browser API
- custom panel docking systems
- mobile-specific redesign
