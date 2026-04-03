# Slice: Initial Energy-Gated Reproduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in JavaScript with 2D canvas rendering
- Go server with authoritative simulation state
- shared runtime contract for visible reproduction eligibility and outcome

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative for whether reproduction is allowed and what energy state results from it.

This slice does not add child entities or ancestry branching. It only makes reproduction obey the already-stated energy-centric life model.

## Discovery Scope

Establish the smallest deterministic rule that makes reproduction require both eligibility capacity and a payable energy cost:

- different-shape overlap may reproduce only when both participating circles satisfy the reproduction-capacity threshold
- successful reproduction consumes energy from both participants
- a low-energy circle with an available child may consume that child as the reproduction payment unit
- failed reproduction due to insufficient payment capacity remains visible and inspectable

This slice does **not** attempt to implement:

- spawned child entities
- variable reproduction costs by shape or size
- cooldown systems
- fertility traits or lineage inheritance
- mating choice beyond current collision overlap

## Why This Slice Next

The current implementation already supports:

- energy spending on movement
- energy recovery through food
- energy collapse death
- different-shape reproduction with child accumulation
- explicit lineage continuity after replacement

But reproduction is still energetically free, which weakens the source claim that energy:

- defines survival
- enables reproduction
- limits movement

Making reproduction energy-gated is the narrowest next step because it:

- directly reinforces energy as the central game currency
- adds visible tradeoff between movement, feeding, and reproduction
- keeps the rule cheap and deterministic
- avoids prematurely modeling child entities or inheritance mechanics

## Use-Case Contract

### Use Case

`ResolveEnergyEligibleReproduction`

### Primary Actor

Player and autonomous circles participating under the same authoritative reproduction rules.

### Pre-conditions

- two circles of different shape overlap
- reproduction is currently resolved through child accumulation
- circles already expose current energy in authoritative snapshots

### Trigger

The server resolves a different-shape collision during a simulation tick.

### Success Outcome

- reproduction succeeds only if both circles satisfy the deterministic reproduction-capacity threshold
- both circles pay the deterministic reproduction energy cost
- both circles still gain one child accumulation unit
- later snapshots expose the lower post-reproduction energy

### Failure Or Rejection Cases

- if reproduction succeeds while one participant lacks sufficient reproduction capacity, the slice fails its purpose
- if reproduction changes energy nondeterministically, the rule loses inspectability
- if the player receives a hidden reproduction exemption, fairness is weakened

## Main Business Rules

1. Reproduction is authoritative server-side behavior.
2. Reproduction requires both participants to satisfy the reproduction-capacity threshold.
3. Successful reproduction consumes energy from both participants.
4. A participant without enough energy may consume one child as the reproduction payment unit.
5. Reproduction remains deterministic and pair-based.
6. Failed reproduction due to insufficient payment capacity must remain inspectable.
7. This slice does not change same-shape fight behavior.

## Minimal Domain Concepts In Scope

- `Reproduction Eligibility`
- `Reproduction Energy Cost`
- `Reproduction Capacity`
- `Child Accumulation`
- `World Snapshot`

## Bounded Reproduction Interpretation

This slice chooses the smallest inspectable interpretation:

- one explicit reproduction-capacity threshold
- one explicit reproduction energy cost
- both values apply equally to player and autonomous circles
- one child may contribute one reproduction-cost unit to the threshold check and payment path
- when either participant cannot satisfy the threshold or pay with energy or one child, no child accumulation occurs for that overlap

This avoids fertility systems, probabilistic outcomes, and shape-specific exceptions.

## Required Runtime Contract Changes

The current snapshot contract already exposes:

- circle energy
- child counts
- resolved interaction kind

This slice should extend the contract only as needed to keep failed eligibility visible, for example by exposing a distinct interaction result when reproduction is blocked by insufficient energy.

## Required Ports Or Boundaries

- server-side eligibility check before reproduction resolution
- server-side deterministic energy deduction on successful reproduction
- client-side rendering of the resulting energy change and blocked or resolved interaction state through ordinary snapshots
- deterministic tests covering both successful and blocked reproduction

## Build Guidance

- prefer one explicit reproduction-capacity threshold, one explicit reproduction cost, and one explicit child-payment fallback
- keep eligibility and cost logic in the authoritative server model
- do not add new client-originated message types
- preserve the current “once per continuous overlap” rule unless the slice explicitly requires otherwise
- keep blocked reproduction inspectable in snapshots and tests

## Initial Test Plan

### Server tests

- different-shape overlap reproduces only when both circles satisfy the reproduction-capacity threshold
- successful reproduction deducts the deterministic energy cost from both participants
- a low-energy circle with one child can still reproduce by consuming that child as the payment unit
- blocked reproduction leaves child counts unchanged
- blocked reproduction does not bypass the energy collapse rule indirectly

### Contract tests

- the snapshot schema remains explicit for the chosen blocked-or-resolved interaction representation

### Integration tests

- the client receives a snapshot showing successful energy-paid reproduction
- the client receives a snapshot showing blocked reproduction when one participant lacks energy

## Scenario Definition

Start a local server and open one browser client.

Scenario steps:

1. two different-shape circles approach overlap
2. in one case both satisfy the threshold and reproduction resolves
3. in another case one circle pays through a child because it lacks enough energy
4. in another case one circle lacks both energy and child payment capacity, so reproduction is blocked
5. the demo and snapshots make the difference visible

## Done Criteria

- reproduction now depends on participant energy
- successful reproduction consumes deterministic energy from both circles
- blocked reproduction is visible and testable
- the player and autonomous circles follow the same rule
- tests cover success and low-energy rejection paths

## Out Of Scope Follow-Ups

- spawned child entities
- fertility cooldowns
- inherited traits
- shape-specific reproduction costs
- ancestry branching during reproduction
