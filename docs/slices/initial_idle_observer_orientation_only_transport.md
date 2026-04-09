# Slice: Initial Idle Observer Orientation-Only Transport

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- passive websocket observers connected to the Go server
- active steering client path already preserved by the previous cadence slice
- browser client consuming the authoritative `world_snapshot`

## Architecture Mode

Explicit client/server boundary with the Go server remaining authoritative.

This slice keeps gameplay unchanged, but stops sending passive observers the same local viewport detail that is primarily needed for active play.

## Discovery Scope

Reduce the remaining passive fanout cost after cadence reduction.

This slice should:

- preserve the current active-client transport path
- preserve the current passive observer cadence policy
- keep minimap and world-orientation usefulness for passive observers
- omit local viewport detail for passive observers on their reduced-cadence snapshots

This slice does **not** attempt to implement:

- compression
- delta encoding
- prediction
- gameplay changes
- active-client protocol thinning

## Why This Slice Next

The new measurement state now shows:

- single-client cadence stays unchanged
- passive fanout dropped materially under the new reduced-cadence policy
- aggregate fanout still scales linearly with client count
- mixed active/passive runs still pay for local viewport detail on passive observer ticks

That makes the next pressure explicit:

- passive cadence alone is no longer the only remaining lever
- the remaining passive cost is now mostly about snapshot content relevance
- passive observers primarily need orientation, not the same local viewport stream as an active participant

This is the narrowest next step because it:

- builds directly on the current active/passive transport split
- leaves active responsiveness untouched
- targets observer-only payload rather than redesigning the whole protocol

## Use-Case Contract

### Use Case

`ReduceIdleObserverSnapshotContent`

### Primary Actor

The Go server broadcasting authoritative snapshots to active and passive websocket clients.

### Pre-conditions

- passive observers already use a lower deterministic cadence when client fanout exists
- active clients already keep the current local interactive cadence
- passive fanout is still the clearest measured remaining transport pressure

### Trigger

A connected client is being treated as a passive observer for transport cadence purposes.

### Success Outcome

- active clients keep the current local-detail transport path
- passive observers keep enough world orientation to remain useful
- passive observer snapshots no longer carry unnecessary local viewport detail
- aggregate passive fanout pressure drops below the new cadence-only baseline

### Failure Or Rejection Cases

- if active-client responsiveness changes, the slice failed
- if passive observers can no longer orient in the world, the slice failed
- if the slice requires gameplay changes, scope drifted

## Main Business Rules

1. This is a transport-content slice, not a gameplay slice.
2. Active clients keep the current local-detail snapshot shape.
3. Passive observers may receive an orientation-only transport view.
4. Orientation usefulness for passive observers must remain explicit and sufficient.
5. A client that becomes active again immediately returns to the current active local-detail path.

## Minimal Domain Concepts In Scope

- `Active Client`
- `Passive Observer`
- `Observer Snapshot`
- `Local Viewport Detail`
- `World Orientation Summary`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful content reduction after cadence differentiation:

- keep passive observer cadence as already built
- remove local viewport entity detail from passive observer snapshots
- keep whole-world orientation and counts so passive observers can still orient
- restore full local detail immediately when a client becomes active again

This avoids deltas and compression while still targeting the remaining passive fanout cost.

## Required Runtime Contract Changes

Likely bounded contract evolution for passive snapshots only.

The build should decide one explicit representation, for example:

- passive snapshots omit `player`, `autonomous_circles`, and `foods`
- passive snapshots keep `world`, totals, and minimap/orientation summaries
- an explicit flag indicates observer-orientation-only transport so the browser can render it intentionally

## Required Ports Or Boundaries

- per-connection transport assembly on the Go server
- shared runtime contract for observer-oriented snapshots
- browser transport handling for active versus passive snapshot content
- transport measurement harness comparing passive fanout before and after content reduction

## Build Guidance

- keep the active path unchanged
- keep the passive cadence policy unchanged
- do not invent a generalized subscription system
- make the observer-oriented representation explicit rather than relying on accidental missing fields
- remeasure the passive fanout ladder after the content reduction

## Initial Test Plan

### Server or measurement tests

- passive observer snapshots are smaller than the current cadence-only passive baseline
- active clients still keep the current local-detail shape
- passive fanout aggregate bytes/sec drops below the current post-cadence baseline
- reactivated clients immediately receive the active local-detail path again

### Contract tests

- the shared snapshot contract explicitly covers the observer-oriented representation

### Integration tests

- one active client plus one passive observer can coexist with different snapshot content
- passive observers still receive enough orientation data to remain useful
- a passive observer that sends input resumes the current active local-detail view

## Scenario Definition

Start one local server with one active steering client and several passive observers.

Scenario steps:

1. connect one active client and several passive observers
2. confirm the active client still receives local viewport detail
3. confirm passive observers receive the orientation-only representation
4. compare passive fanout bytes/sec to the current post-cadence baseline
5. confirm a passive observer becomes active again and regains local detail after sending movement intent

## Done Criteria

- passive observer aggregate transport pressure drops again from the new post-cadence baseline
- active-client responsiveness remains unchanged
- passive observers still have useful world orientation
- the observer-oriented representation is explicit in the contract

## Out Of Scope Follow-Ups

- compression
- deltas
- prediction
- generalized field subscription
- gameplay changes
