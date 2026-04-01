# Domain Background Knowledge

## Purpose

This document captures broad background knowledge about the target domain.

It is not the repository's extracted glossary. It is reference material distilled from books, articles, standards, external systems, or other domain sources.

Use it mainly during expectation-gap detection.

---

## Domain Framing

Living Circles sits at the intersection of three recognizable domains:

- agar.io-like real-time arena games
- emergent simulations with simple-rule complexity
- lightweight ecosystem or population-dynamics sandboxes

This background knowledge is useful during review because players and reviewers will likely import expectations from those domains even when Living Circles intentionally diverges from them.

## Common Concepts In The Domain

- `authoritative server`: the server defines the true simulation state
- `client rendering`: the browser presents a visual approximation of that state
- `tick`: discrete simulation step used to update movement, collisions, and outcomes
- `collision`: overlap between entities that triggers interaction rules
- `bounded world`: a finite map where entities cannot drift forever outside play
- `latency`: delay between input, simulation, and visible response
- `emergence`: complex macro behavior arising from small local rules
- `population dynamics`: rise, decline, dominance, collapse, and recovery patterns within a shared system

## External Expectations Users May Bring

Players familiar with agar.io-like games may expect:

- immediate directional control
- growth through consuming food
- larger entities dominating smaller ones
- a readable risk-reward loop
- strong responsiveness and low perceived latency
- a clear sense of danger, evasion, and opportunity in crowded spaces

Players familiar with ecosystem simulations may expect:

- non-trivial autonomous behavior without explicit scripting
- feedback loops between resources and population size
- visible stability, collapse, and recovery cycles
- fairness and consistency in how rules apply

## Standard Artifacts Or Outputs The Domain Usually Implies

- a real-time world state with circles, food, and boundaries
- server tick updates or state snapshots
- visual motion and collision feedback
- recognizable growth and death events
- sessions in which multiple entities coexist and compete
- comparative evaluation of responsiveness, fairness, and emergent behavior

For this project specifically, the source material implies these runtime expectations:

- a browser-based JavaScript client using 2D canvas rendering
- an authoritative Go simulation server
- WebSocket-based real-time communication

Those expectations matter during evaluation because latency, interpolation, authority, and update cadence are not optional UI details in this genre. They shape whether the experience feels coherent and fair.

## Industry Language Worth Preserving

From agar.io-like references:

- `cell`
- `food`
- `movement`
- `growth`
- `consumption`
- `collision`
- `world`
- `tick`
- `latency`
- `server authority`
- `risk/reward`
- `scaling`

From the Living Circles material:

- `circle`
- `energy`
- `fight`
- `reproduce`
- `children`
- `lineage`
- `continuity`
- `player as perturbation`
- `system-driven ecosystem`

During future refinement, the project should preserve whichever terms best express the intended identity rather than forcing total vocabulary alignment with agar.io.

## Comparison Baselines That Matter During Evaluation

The agar.io references provide a comparison baseline for:

- familiarity of the play loop
- responsiveness of real-time control and updates
- readability of collisions and outcomes
- bounded-world behavior
- scalability of many concurrent entities
- coherence between browser client perception and authoritative server state

They should not be treated as proof that Living Circles needs:

- split mechanics
- merge mechanics
- pure size-dominance combat
- skill-first rather than system-first identity

## Likely Omissions To Watch For During Evaluation

- absence of visible energy pressure, which would weaken the game's core identity
- interactions that feel arbitrary because shape-based rules are not legible to the player
- a world that lacks enough food or collisions to produce emergence
- excessive hidden advantages for the player, which would undermine fairness
- lack of continuity or lineage effects, which would collapse the distinctiveness of the model
- overfitting to agar.io expectations and thereby erasing the project's intended differences
- poor responsiveness or high latency, which would still damage playability even if the domain model is sound
- lack of bounded-world behavior, causing weak encounter density

## Review Heuristics

When reviewing future slices, ask:

- does the behavior reinforce energy as the central game currency?
- can a reviewer see the difference between same-shape conflict and different-shape reproduction?
- do simple rules produce interpretable emergent outcomes?
- does the player experience the same rule system as other circles?
- is the game borrowing from agar.io only where that improves accessibility and not where it would erase the distinct domain identity?

## Source Evidence

- `work/sources/living-circles-initial/living-circles.md`
- `work/sources/living-circles-initial/agar-io-reference.md`
- `work/sources/living-circles-initial/agar-io-glossary.md`
