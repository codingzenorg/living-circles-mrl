# Slice: Initial Child Power In Fight Resolution

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract carrying the resulting fight outcome through ordinary snapshots

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for fight outcome calculation.

This slice does not add animation, damage-over-time, or multi-tick combat. It only makes child accumulation a direct input to same-shape fight resolution.

## Discovery Scope

Establish the smallest deterministic rule that makes children explicitly increase fight power:

- same-shape fight resolution should consider child-based power directly
- child-based power should be readable as an explicit rule rather than only an indirect radius side effect
- the fight outcome should remain one-tick and deterministic

This slice does **not** attempt to implement:

- combat rounds
- hit points
- damage exchange
- probabilistic critical outcomes
- separate attack and defense stats
- fight aftermath rewards beyond the current continuity and removal rules

## Why This Slice Next

The current implementation already supports:

- same-shape fight resolution
- radius growth from children
- lineage continuity after defeat
- energy-gated reproduction
- food-seeking autonomy

But the source material says children:

- increase area
- increase power in fights

Right now only the first part is explicit. Children affect fights only indirectly through radius, and only as a tie-break. That leaves the model under-specified where the source is more direct.

Making child power explicit is the narrowest next step because it:

- resolves a major remaining semantic ambiguity
- reinforces children as a meaningful strategic state
- keeps combat deterministic and cheap
- avoids opening a full combat subsystem

## Use-Case Contract

### Use Case

`ResolveSameShapeFightWithChildPower`

### Primary Actor

Player and autonomous circles colliding under the same-shape fight rule.

### Pre-conditions

- two same-shape circles overlap
- each circle exposes energy, radius, and child count
- the server resolves fights authoritatively

### Trigger

The server resolves a same-shape collision during a simulation tick.

### Success Outcome

- fight winner selection includes child-based power directly
- the resulting winner and loser are still exposed through the existing interaction snapshot
- replacement continuity and loser removal continue to work as before

### Failure Or Rejection Cases

- if children still affect fights only indirectly, the slice fails its purpose
- if fight power becomes opaque or hard to inspect, the slice weakens the model
- if combat expands into a larger subsystem, the slice exceeds scope

## Main Business Rules

1. Same-shape fights remain authoritative server-side behavior.
2. Child accumulation must contribute directly to fight power.
3. Fight resolution remains deterministic and one-tick.
4. Energy still matters in fight resolution.
5. This slice does not change different-shape reproduction behavior.

## Minimal Domain Concepts In Scope

- `Fight Power`
- `Child-Based Power`
- `Same-Shape Fight`
- `Winner Selection`
- `World Snapshot`

## Bounded Fight Interpretation

This slice chooses the smallest inspectable interpretation:

- same-shape fights compare a deterministic power ordering
- child count is an explicit term in that ordering
- radius and energy may remain part of the ordering
- ties are still resolved deterministically

This avoids hit points, battle phases, and opaque formulas while making child power real.

## Required Runtime Contract Changes

No new message types are required.

The existing interaction snapshot should remain sufficient unless build determines that one small explanatory field is necessary for inspectability.

## Required Ports Or Boundaries

- server-side fight outcome calculation with child-based power
- deterministic tests covering child-driven winner changes
- client-side rendering through ordinary snapshots

## Build Guidance

- prefer one explicit ordering or one explicit fight-power formula
- keep the rule inspectable in tests
- do not add animation or temporal combat phases
- preserve the existing loser removal and replacement-continuity paths
- avoid hidden randomness

## Initial Test Plan

### Server tests

- a higher-child circle can beat a lower-child circle even when radius or energy alone would not have decided it before
- child power is deterministic and inspectable in fight outcomes
- existing replacement and loser-removal behavior still works with the new winner rule

### Contract tests

- the existing snapshot schema remains sufficient unless a new explanatory field is added intentionally

### Integration tests

- the client receives a same-shape fight outcome where child power is the decisive factor

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two same-shape circles overlap
2. one circle has higher child-based fight power
3. the authoritative snapshot shows that circle winning
4. the result remains compatible with existing replacement or removal behavior

## Done Criteria

- children now influence same-shape fights directly
- the winner rule remains deterministic and inspectable
- existing fight aftermath rules still hold
- tests cover child-driven winner selection
- the slice avoids expanding into a full combat subsystem

## Out Of Scope Follow-Ups

- multi-tick combat
- damage systems
- explicit attack or defense stats
- combat animations
- fight-generated resources
