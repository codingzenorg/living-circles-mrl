# Slice: Initial Idle Observer Transport Cadence Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- passive websocket clients that are connected but not actively steering
- Go authoritative server with the current optimized transport path
- shared runtime contract unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and message shape unchanged, but stops treating passive observers as if they always need the same full high-cadence local stream as an actively steering player.

## Discovery Scope

Reduce the clearest remaining measured transport pressure: broad passive snapshot fanout.

This slice should:

- preserve the current transport path for the active player experience
- identify clients with no recent movement intent as passive observers
- send those passive observers snapshots at a lower deterministic cadence
- leave snapshot content and gameplay semantics unchanged

This slice does **not** attempt to implement:

- selective field subscription
- different protocol shapes for different clients
- compression
- prediction
- gameplay changes

## Why This Slice Next

The current measurement track now shows:

- movement changes transport only slightly in the bounded idle-versus-moving comparison
- passive client-count fanout scales almost linearly
- per-client cost stays roughly flat while aggregate output rises directly with client count

That makes the next pressure explicit:

- the server is paying the full current cadence for clients that may not need it
- the cheapest next mitigation is cadence differentiation, not another content redesign

This is the narrowest next step because it:

- uses existing connection state rather than a new protocol
- targets the measured fanout pressure directly
- preserves the active player path while thinning passive observer output

## Use-Case Contract

### Use Case

`ReducePassiveObserverSnapshotCadence`

### Primary Actor

The Go server broadcasting authoritative state to both active and passive websocket clients.

### Pre-conditions

- the current transport path is already optimized by content and freshness rules
- the server already receives explicit movement intents from active clients
- passive observer fanout is now the clearest measured remaining transport pressure

### Trigger

A connected client has not sent recent movement intent and is therefore treated as a passive observer for cadence purposes.

### Success Outcome

- active clients keep the current interactive cadence
- passive observers receive deterministic lower-cadence snapshots
- aggregate fanout pressure drops under passive multi-client ladders
- protocol shape remains unchanged

### Failure Or Rejection Cases

- if active clients lose responsiveness, the slice failed
- if passive observers can become stale for too long to remain useful, the slice failed
- if the slice requires a protocol split or observer-specific snapshot format, scope drifted

## Main Business Rules

1. This is a transport policy slice, not a gameplay slice.
2. Active clients keep the current cadence.
3. Passive observers may receive a lower deterministic cadence.
4. A client becomes active again as soon as fresh movement intent arrives.
5. Snapshot content stays the same; only cadence differs.

## Minimal Domain Concepts In Scope

- `Active Client`
- `Passive Observer`
- `Recent Movement Intent`
- `Observer Cadence`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful observer-specific mitigation:

- track whether a connection has sent movement intent recently
- apply a lower fixed cadence to passive observers
- immediately restore the current cadence when input resumes

This avoids a more complex subscription model while still targeting the measured fanout pressure.

## Required Runtime Contract Changes

None.

## Required Ports Or Boundaries

- websocket connection state on the server
- broadcast cadence decision per connection
- transport measurement harness updated to compare passive fanout before and after the cadence change

## Build Guidance

- keep the active/passive threshold deterministic and explicit
- do not delay active-player snapshots
- preserve current snapshot structure so the browser client needs little or no change
- measure the passive fanout ladder again after the change

## Initial Test Plan

### Server or measurement tests

- active clients still receive the current cadence
- passive observers receive the reduced cadence
- clients that resume movement return to active cadence
- passive fanout aggregate bytes/sec drops below the current baseline

### Contract tests

- none beyond the current contract validation, because shape stays unchanged

### Integration tests

- one active client plus one passive observer can coexist with different effective cadences
- passive observers still receive enough snapshots to remain orienting and useful

## Scenario Definition

Start one local server with a small mix of active and passive clients.

Scenario steps:

1. connect one active steering client and several passive observers
2. verify that the active client keeps the current cadence
3. verify that passive observers receive fewer snapshots
4. compare passive fanout bytes/sec to the current fanout baseline
5. verify that an observer becomes active again when movement intent resumes

## Done Criteria

- passive observer fanout pressure drops measurably
- active player responsiveness remains unchanged
- no protocol shape or gameplay semantics changed

## Out Of Scope Follow-Ups

- separate observer protocol shapes
- compression
- prediction
- transport sharding
