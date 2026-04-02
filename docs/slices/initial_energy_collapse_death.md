# Slice: Initial Energy Collapse Death

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for energy-driven death outcomes

## Architecture Mode

Explicit client/server boundary with the server resolving zero-energy collapse into death or continuity.

This slice extends the current movement and continuity model without introducing hunger subsystems, decay timers, or environmental hazards beyond the existing energy loop.

## Discovery Scope

Establish the smallest deterministic death rule for energy collapse:

- energy is no longer only a movement throttle
- when a circle reaches zero energy, it dies on the server
- childless circles disappear
- circles with children continue through the same one-child replacement rule already used after fight defeat
- the death outcome is visible in later snapshots

This slice does **not** attempt to implement:

- gradual health loss or damage-over-time
- partial starvation states
- nonzero energy penalties beyond movement stop
- replacement inheritance beyond the already chosen continuity rule
- food respawn balancing or ecosystem tuning
- death causes other than zero-energy collapse

## Why This Slice Next

The source material says energy:

- defines survival
- enables reproduction
- limits movement

The current implementation only covers the third point directly. Zero energy currently freezes movement but does not create death, which weakens the model's central variable.

Making energy collapse lethal is the narrowest next step because it:

- turns energy into an actual survival boundary
- strengthens the meaning of feeding and resource pressure
- reuses the existing continuity rule instead of inventing a second death model
- keeps the rule cheap, inspectable, and deterministic

## Use-Case Contract

### Use Case

`ResolveEnergyCollapseDeath`

### Primary Actor

Player controlling one circle while autonomous circles follow the same survival rule.

### Pre-conditions

- a server process can host one bounded world
- a circle can spend energy through movement
- a circle can reach zero energy
- a circle may or may not have accumulated children

### Trigger

The server advances a tick after a circle reaches zero energy.

### Success Outcome

- a zero-energy circle is treated as dead
- if it has no children, it disappears from active participation
- if it has one or more children, one child is consumed and continuity continues through replacement
- later snapshots expose the resulting world state clearly enough for the client to inspect

### Failure Or Rejection Cases

- if energy is above zero, no energy-collapse death occurs
- if a zero-energy circle remains frozen forever without dying, this slice fails its purpose
- if energy-collapse continuity uses a different replacement rule than fight defeat without a documented reason, the slice becomes inconsistent

## Main Business Rules

1. Zero energy means death, not only immobility.
2. Energy-collapse death is resolved authoritatively on the server.
3. Childless circles disappear when they die from energy collapse.
4. Circles with at least one child continue through the same one-child replacement rule already used after fight defeat.
5. The same energy-collapse rule applies to player and autonomous circles.
6. This slice does not add a distinct starvation-specific inheritance model.

## Minimal Domain Concepts In Scope

- `Energy Collapse`
- `Death`
- `Continuity`
- `Replacement`
- `Active Circle`
- `World Snapshot`

## Bounded Survival Interpretation

The source material says energy defines survival, but it does not specify a richer starvation model.

This slice therefore chooses the smallest explicit interpretation:

- energy equal to `0` is the death threshold
- death resolves immediately on the authoritative server
- the already-established replacement rule is reused if children are available

That keeps survival simple and coherent with the rest of the current model.

## Required Runtime Contract Changes

The existing snapshot contract may already be sufficient if energy-collapse death is represented through:

- disappearance of the dead circle
- or continued presence through replacement
- plus the normal snapshot fields and current interaction/death context when applicable

If the build needs an explicit death outcome marker for inspectability, it may add one only if the existing contract proves too ambiguous.

## Required Ports Or Boundaries

- server-side zero-energy detection
- server-side death handling shared between player and autonomous circles
- reuse or adaptation of the current continuity rule for non-fight death
- deterministic tests for disappearance and replacement on energy collapse

## Build Guidance

- keep the death threshold explicit: `energy == 0`
- prefer reusing the current replacement rule over introducing a second continuity mechanism
- avoid a generalized status-effect framework
- make the outcome inspectable through normal snapshots and tests
- do not rebalance food, movement cost, or growth in this slice unless needed for deterministic tests

## Initial Test Plan

### Server tests

- a childless player with zero energy disappears
- a childless autonomous circle with zero energy disappears
- a zero-energy circle with children remains active through replacement
- replacement consumes exactly one child
- zero-energy death does not trigger while energy is still above zero

### Contract tests

- the current snapshot schema remains sufficient unless an explicit death marker is added

### Integration tests

- the client receives a later snapshot where a zero-energy circle disappears or continues through replacement
- the browser-visible state is enough to inspect that the reset world can be replayed after energy-collapse scenarios

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. the client connects and moves until one circle reaches zero energy
2. the server resolves zero-energy death
3. a childless lineage disappears, or a lineage with children remains through replacement
4. the client receives a later snapshot showing the energy-collapse outcome clearly

## Done Criteria

- zero energy now causes death
- death can still lead to disappearance or replacement depending on child availability
- player and autonomous circles follow the same rule
- tests cover both disappearance and replacement paths
- the slice does not introduce hunger timers, mutations, or a second continuity model

## Out Of Scope Follow-Ups

- slower starvation stages before death
- food respawn balancing to sustain populations
- multiple death causes with distinct consequences
- explicit death history or lineage visualization
- decay or recovery from near-zero energy
