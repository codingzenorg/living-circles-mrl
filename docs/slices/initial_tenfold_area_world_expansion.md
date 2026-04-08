# Slice: Initial Tenfold Area World Expansion

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract where ordinary snapshots already carry the larger world bounds without requiring a new protocol surface

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for world size, placement, and reset state.

This slice changes the size of the authoritative expanded world, not the viewport model. The player should continue moving through the same on-screen camera window, using the minimap as the main whole-world orientation aid.

## Discovery Scope

Establish the smallest coherent map-scale change that makes the expanded world roughly `10x` larger in total area:

- keep the current expanded-world aspect ratio
- scale width and height by approximately `sqrt(10)`
- preserve the current viewport, deadzone camera, offscreen cues, and minimap model
- adjust only the default seeded expanded baseline enough to avoid an obviously empty-feeling world

This slice does **not** attempt to implement:

- zoom controls
- a larger on-screen viewport
- minimap redesign
- chunking, streaming, or server sharding
- terrain systems
- new autonomy or reproduction rules

## Why This Slice Next

The current presentation model is now built around:

- a fixed on-screen viewport
- a player-follow camera
- a minimap for whole-world orientation

That makes a much larger authoritative map not only possible, but desirable. The current expanded world is materially larger than the old baseline, yet it is still modest relative to the viewport/navigation model now in place.

The next pressure is:

- make the map feel substantially larger as navigable space
- keep the player’s on-screen window stable
- let the minimap matter more as an orientation tool

This slice is the narrowest next step because it:

- changes map scale rather than introducing a new system
- stays compatible with the existing viewport/minimap model
- can be implemented mostly through world constants and deterministic expanded-world seeding

## Use-Case Contract

### Use Case

`StartExpandedWorldAtTenfoldArea`

### Primary Actor

Any session or reset using the expanded default world baseline.

### Pre-conditions

- the expanded world already has deterministic width, height, seeded food, and seeded autonomous startup
- the browser already renders a bounded viewport into authoritative world coordinates
- the minimap already provides a whole-world orientation aid

### Trigger

A new expanded world is created or reset.

### Success Outcome

- the expanded world starts at approximately `10x` the previous total area
- width and height preserve the current aspect ratio
- viewport behavior remains unchanged on screen
- the default seeded population and food baseline remain plausible for the larger territory

### Failure Or Rejection Cases

- if the viewport size changes with the map, scope is exceeded
- if the map becomes obviously empty because startup seeding was not reconsidered, the slice is incomplete
- if the change introduces nondeterministic startup layouts, reset/test coherence weakens

## Main Business Rules

1. The world-size change is authoritative server-side behavior.
2. The expanded world should preserve its current aspect ratio while increasing total area by about `10x`.
3. The viewport and minimap behavior remain client-side presentation layers over the same authoritative coordinates.
4. The expanded startup population and food baseline may be adjusted only as much as needed to keep the larger map plausible.
5. Reset should recreate the same larger expanded world deterministically.
6. Food, fight, reproduction, continuity, child, and steering rules remain unchanged.

## Minimal Domain Concepts In Scope

- `World Bounds`
- `Expanded Baseline`
- `Seeded Startup Layout`
- `Viewport Orientation`

## Bounded Scale Interpretation

This slice chooses the smallest coherent interpretation of “`10x` greater map”:

- treat `10x` as total area, not linear dimensions
- preserve the current aspect ratio
- choose rounded width and height close to `sqrt(10)` scaling of the current expanded world
- keep the viewport size fixed while the authoritative map grows

This avoids turning the slice into a camera redesign while still delivering a meaningfully larger world.

## Required Runtime Contract Changes

The current contract is likely sufficient because snapshots already expose `world.width` and `world.height`.

Build should avoid new contract fields unless a tiny inspectability aid becomes necessary.

## Required Ports Or Boundaries

- server-side expanded world initialization
- deterministic seeded food and autonomous placement for the new larger bounds
- deterministic tests that prove reset and startup coherence under the new map size
- implementation notes that record the chosen expanded dimensions and baseline adjustments

## Build Guidance

- preserve the current expanded-world aspect ratio
- document the chosen rounded dimensions and their relation to the old expanded area
- keep narrow custom worlds unchanged
- adjust expanded startup food/population only if the larger map would otherwise feel obviously sparse
- avoid mixing new ecology or camera rules into the same slice

## Initial Test Plan

### Server tests

- the default expanded world reports the new larger deterministic bounds
- the expanded startup population and food baseline still initialize deterministically inside those bounds
- reset reproduces the same larger expanded world exactly

### Contract tests

- the current snapshot schema remains sufficient because world bounds are already explicit

### Integration tests

- the client still receives authoritative larger bounds and seeded startup layout through ordinary snapshots
- viewport mode and minimap continue to function without contract changes

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the expanded world starts with the new much larger bounds
2. the browser still shows the same viewport-sized play window
3. the player navigates through the larger map using the minimap for orientation
4. reset reproduces the same larger seeded world

## Done Criteria

- the expanded world area is approximately `10x` the previous expanded map area
- the aspect ratio is preserved
- the on-screen viewport model remains unchanged
- expanded startup seeding remains deterministic and plausible
- tests cover world bounds and reset/startup coherence

## Out Of Scope Follow-Ups

- zoom controls
- minimap redesign
- procedural chunk streaming
- terrain or biome systems
- new ecological rules introduced only to compensate for scale
