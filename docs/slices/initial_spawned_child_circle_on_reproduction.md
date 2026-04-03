# Slice: Initial Spawned Child Circle On Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract carrying a newly active child circle through ordinary snapshots

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether reproduction creates a new active circle.

This slice does not add full lineage trees, mutation, or separate reproductive castes. It turns one part of the current abstract child model into an observable active entity.

## Discovery Scope

Establish the smallest deterministic rule that makes reproduction create a visible child circle:

- a successful different-shape reproduction may spawn one new active circle
- the spawned child should be positioned deterministically relative to the parents
- the child should enter the same world rules as every other circle

This slice does **not** attempt to implement:

- multiple children per reproduction
- random mutation
- inherited trait bundles
- delayed incubation
- separate juvenile/adult phases
- branching ancestry history views

## Why This Slice Next

The current implementation already supports:

- reproduction eligibility and cost
- child accumulation as an abstract counter
- radius growth from children
- lineage continuity after defeat
- explicit fight and food loops

But reproduction is still only partially embodied. The source material leaves open whether children are abstract or spawned, and the current model still lacks any direct transition from reproduction into a new active participant.

Spawning one child circle is the narrowest next step because it:

- makes reproduction visible as population change rather than only hidden accumulation
- strengthens the ecosystem aspect of the model
- creates a direct bridge between reproduction and emergent encounter density
- stays smaller than full inheritance or lifecycle systems

## Use-Case Contract

### Use Case

`SpawnChildCircleFromReproduction`

### Primary Actor

Two different-shape circles reproducing under the current authoritative reproduction rules.

### Pre-conditions

- reproduction is successful under the current energy and child-reserve rules
- the server owns authoritative world state
- the world can host more than the initial set of circles

### Trigger

The server resolves a successful different-shape reproduction during a simulation tick.

### Success Outcome

- one new active child circle is added to the world
- the child has deterministic initial position, shape, size, and energy
- later snapshots include the child as a normal active participant

### Failure Or Rejection Cases

- if successful reproduction still only changes counters, the slice fails its purpose
- if child spawning is random or opaque, the slice loses inspectability
- if the child enters the world under different movement or collision rules, fairness is weakened

## Main Business Rules

1. Successful reproduction may create one new active child circle.
2. Child spawning is authoritative server-side behavior.
3. Child spawn position and initial state must be deterministic.
4. A spawned child follows the same world rules as every other active circle.
5. This slice does not remove the existing notion of lineage.

## Minimal Domain Concepts In Scope

- `Spawned Child Circle`
- `Reproduction Outcome`
- `Initial Child State`
- `World Participation`
- `World Snapshot`

## Bounded Child Interpretation

This slice chooses the smallest inspectable embodiment rule:

- exactly one child circle is spawned per successful reproduction
- the child starts with a deterministic baseline shape and energy policy
- the child is immediately active in the same simulation

This avoids incubation, trait inheritance, and probabilistic offspring generation while making reproduction materially visible.

## Required Runtime Contract Changes

The existing snapshot contract may remain structurally sufficient if spawned children are represented as ordinary active circles.

Build should only extend the contract if a small explicit field is necessary to preserve inspectability of parent/child origin.

## Required Ports Or Boundaries

- server-side child creation during reproduction resolution
- deterministic ID and lineage handling for spawned children
- client-side rendering of the child as an ordinary active circle
- deterministic tests covering spawn count and initial child state

## Build Guidance

- prefer one explicit spawn rule over a generic factory system
- keep initial child state small and deterministic
- avoid adding generalized inheritance mechanics in this slice
- preserve the current reproduction eligibility and energy-cost rules
- keep the child inside the same movement, collision, and death model

## Initial Test Plan

### Server tests

- a successful reproduction adds exactly one active child circle
- the child has deterministic initial state
- blocked reproduction does not create a child
- the child participates in later ticks under ordinary rules

### Contract tests

- the existing snapshot schema remains sufficient unless one explicit origin field is added intentionally

### Integration tests

- the client receives a later snapshot with an additional active circle created by reproduction

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two different-shape circles reproduce successfully
2. the next authoritative snapshot includes one new active child circle
3. later ticks show that child participating under the same world rules

## Done Criteria

- successful reproduction now creates one active child circle
- spawn state is deterministic and inspectable
- blocked reproduction does not create children
- the child follows the same world rules as other circles
- tests cover spawn creation and non-creation paths

## Out Of Scope Follow-Ups

- multiple offspring
- mutation
- trait inheritance
- incubation delays
- ancestry tree visualization
