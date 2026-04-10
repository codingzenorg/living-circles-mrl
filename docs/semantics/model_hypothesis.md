# Model Hypothesis

## Purpose

This document captures the current model hypothesis for the target domain.

Use it during `extract`, `refine`, and `build` to define vocabulary, boundaries, use cases, and slice candidates.

---

## Current Hypothesis

Living Circles is a system-driven 2D ecosystem game in which autonomous circles coexist in a shared world and follow a minimal life cycle driven by movement, energy use, feeding, collision-based interaction, growth, and death. The model is not character-centric and should not depend on explicit AI behavior trees or scripted roles. Its intended identity is an emergent simulation with a player participating inside the same rule set as every other circle.

The comparison to agar.io is useful as a validation baseline for familiarity, responsiveness, and playability, but not as a target for feature parity. Living Circles should preserve its own identity by centering energy, shape-based interaction, lineage, and emergence over direct control and pure size dominance.

The source material also carries an important architectural signal: Living Circles is conceived as a JavaScript browser game rendered in 2D canvas, communicating over WebSocket with an authoritative Go server. This should be treated as part of the initial system hypothesis, not as an incidental implementation detail to be ignored during refinement.

## Core Concepts

- `World`: the bounded shared environment where circles and food exist and interact.
- `Circle`: the main living entity that moves, spends energy, eats, collides, grows, reproduces, fights, and dies.
- `Food`: world resource consumed by circles to recover energy.
- `Energy`: the central variable connecting survival, movement, combat readiness, and reproduction capacity.
- `Shape`: the circle classification that determines interaction semantics during collisions.
- `Child`: accumulated or spawned continuity unit that can increase a circle's leverage and may replace a dead circle.
- `Lineage`: the continuity relationship between a circle and its children across death and replacement.
- `Player Circle`: a human-participating circle that still follows the same domain rules as non-player circles.

## Likely Entities And Boundaries

### World

The world owns the simulation space, boundaries, and co-presence of circles and food. It is the container in which real-time interactions occur.

### Circle

The circle is the primary domain entity. It carries position, energy, shape, growth state, and possibly lineage-related continuity. It is the actor through which nearly all important game outcomes emerge.

### Food

Food is a consumable resource rather than a strategic actor. Its role is to sustain the energy loop and thereby the population dynamics.

## Key Value Objects Or Domain Attributes

- `Position`: spatial location inside the world.
- `Energy Level`: current survivability and action capacity.
- `Shape`: interaction category used to determine whether a collision leads toward fight or reproduction.
- `Size` or `Area`: current leverage for collection and conflict outcomes.
- `Tick`: discrete simulation step for authoritative updates.

## Major State Transitions

- `Move -> Spend Energy`: movement is never free and consumes the central resource.
- `Eat -> Recover Energy`: feeding replenishes the resource required for continued life.
- `Collide -> Resolve Interaction`: collisions are the trigger for domain events.
- `Same Shape -> Fight`: same-shape collisions should trend toward conflict behavior.
- `Different Shape -> Reproduce`: different-shape collisions should trend toward reproduction behavior.
- `Accumulate Children -> Grow`: child accumulation increases leverage, likely affecting area and conflict capability.
- `Energy Collapse Or Defeat -> Die`: a circle can exit the active system through death.
- `Death -> Disappear Or Be Replaced`: death can either terminate the line or continue it through child replacement.

## Candidate Use Cases

- start and run a shared world simulation
- place circles and food into the world
- process player input for a player circle without granting hidden advantages
- advance one simulation tick
- move circles and apply energy cost
- resolve feeding when circles reach food
- resolve collision outcomes between circles
- determine whether an interaction becomes fight or reproduction
- update growth and lineage after successful interactions
- remove dead circles or replace them through continuity rules
- broadcast authoritative state to connected clients

## Governing Domain Principles

- emergence over scripted behavior
- simple rules with dense interaction
- energy as the unifying model variable
- fairness through shared rules for player and non-player circles
- low-cost rule resolution to support many entities
- conditions for outcomes rather than direct control of outcomes

## Architectural Pressure From The Source Material

The current source evidence implies a multi-runtime system rather than a single-process application:

- a JavaScript browser client is responsible for 2D canvas rendering and player input
- a Go server is responsible for authoritative state, simulation ticks, collision handling, and rule resolution
- the runtimes communicate over WebSocket with frequent state updates

This matters because it affects where simulation authority lives, how responsiveness is evaluated, and how future slices should describe boundaries. The domain model remains shared, but the repository should not pretend the system is naturally a single-runtime monolith if the intended product shape is client/server from the beginning.

## Relationship To Agar.io

Living Circles intentionally overlaps with agar.io-like expectations in these areas:

- real-time shared world
- browser client plus authoritative server
- movement, collisions, growth, and bounded map
- simple mechanics that produce complex outcomes

Living Circles intentionally diverges in these areas:

- interaction is shape-based rather than pure size dominance
- reproduction and lineage matter, not only consumption
- emergence and ecosystem dynamics matter more than player skill expression alone
- the player is one perturbing participant, not a privileged commander

The agar.io documents should be treated as reference material for evaluation and language, not as direct implementation requirements.

## Unresolved Tensions And Ambiguities

- The exact semantics of `fight` are not yet explicit: whether resolution is based on energy, size, children, randomness, or another factor remains open.
- The exact semantics of `reproduce` are not yet explicit: whether a child is spawned immediately, accumulated abstractly, or represented as a separate circle remains open.
- The relation between `children`, `growth`, and `power` needs clearer modeling. The source says children increase area and fight power, but not the exact mechanism.
- The death continuity rule is intentionally open. A child may replace the dead circle, but the trigger and replacement semantics are not yet defined.
- The degree of autonomy for non-player circles is implied by emergence, but the movement policy for autonomous circles is still unspecified.
- Determinism is presented as optional but valuable. The target model should decide whether replayability and reproducibility are part of the intended core.
- The extracted architecture is JavaScript canvas client plus authoritative Go server over WebSocket, but the repository still needs an explicit architectural decision on how strongly that should shape the initial pack and directory structure.
- Real concurrent responsiveness is now an observed runtime risk: opening a second browser session against the local demo caused visible slowdown in both the simulation and player movement, so multi-client responsiveness should be treated as an explicit evaluation and refinement concern rather than only an inferred transport concern.

## Source Evidence

- `work/sources/living-circles-initial/living-circles.md`
- `work/sources/living-circles-initial/agar-io-reference.md`
- `work/sources/living-circles-initial/agar-io-glossary.md`
