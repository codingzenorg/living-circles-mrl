# Slice: Initial Event-Driven Idle Observer Transport

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- passive websocket observers already using `observer_orientation_only` transport
- active steering clients already using `active_local_detail`
- Go authoritative server with the current measured fanout path

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and active transport unchanged, but stops sending lower-cadence passive observer snapshots purely because the passive cadence timer elapsed when the observer-oriented content itself has not meaningfully changed.

## Discovery Scope

Reduce the remaining passive observer transport cost after the orientation-only split.

This slice should:

- preserve the current active-client path unchanged
- preserve the current observer-oriented snapshot representation
- send passive observer snapshots when observer-relevant state actually changes
- retain a slower deterministic fallback refresh so passive observers do not drift indefinitely

This slice does **not** attempt to implement:

- compression
- delta encoding
- generalized subscriptions
- gameplay changes
- active-client cadence reduction

## Why This Slice Next

The new measurement state now shows:

- passive orientation-only transport already reduced fanout below the cadence-only observer baseline
- single-client active transport remains unchanged
- passive fanout still scales linearly with client count
- passive observer ticks can still resend the same orientation-only payload even when nothing observer-relevant changed

That makes the next pressure explicit:

- the remaining passive cost is no longer mostly local detail
- it is now mostly repeated observer-oriented snapshots under calm world conditions
- the next bounded win is event-driven observer refresh, not another content redesign

This is the narrowest next step because it:

- builds directly on the explicit observer transport mode
- keeps active play untouched
- targets redundant passive refreshes rather than already-optimized snapshot content

## Use-Case Contract

### Use Case

`RefreshIdleObserverTransportOnlyWhenObserverStateChanges`

### Primary Actor

The Go server broadcasting authoritative snapshots to active and passive websocket clients.

### Pre-conditions

- passive observers already use lower-cadence orientation-only transport when fanout exists
- active clients already keep the local interactive path
- passive fanout remains the clearest measured transport pressure

### Trigger

A passive observer connection reaches a potential send tick and the server decides whether observer-relevant state has materially changed since the last observer snapshot for that connection class.

### Success Outcome

- active clients remain unchanged
- passive observers still receive useful orientation updates
- redundant observer snapshots are skipped when observer-relevant state is unchanged
- aggregate passive fanout pressure drops below the current orientation-only observer baseline

### Failure Or Rejection Cases

- if active responsiveness changes, the slice failed
- if passive observers can drift too long without useful updates, the slice failed
- if the slice requires a generalized multi-protocol transport redesign, scope drifted

## Main Business Rules

1. This is an observer-refresh policy slice, not a gameplay slice.
2. Active clients keep the current local-detail cadence and payload.
3. Passive observers keep the current orientation-only representation.
4. Passive observer snapshots should refresh when observer-relevant state materially changes.
5. A deterministic fallback refresh must still exist for passive observers.

## Minimal Domain Concepts In Scope

- `Passive Observer`
- `Observer-Relevant State`
- `Observer Refresh Signature`
- `Observer Fallback Refresh`

## Bounded Optimization Interpretation

This slice chooses the smallest meaningful event-driven observer policy:

- derive a deterministic signature from observer-oriented snapshot content
- send passive observer snapshots when that signature changes
- otherwise skip passive sends until a slower fallback interval is reached

This avoids deltas and compression while still targeting the remaining passive fanout cost.

## Required Runtime Contract Changes

Likely none.

The current explicit `transport_mode` plus existing `orientation_fresh` and observer-oriented snapshot fields should already be enough if only the passive send policy changes.

## Required Ports Or Boundaries

- passive observer send decision in the Go transport server
- deterministic observer-signature derivation from the current observer-oriented snapshot
- transport measurement harness updated to compare the new passive baseline

## Build Guidance

- keep active-client behavior unchanged
- keep the observer-oriented snapshot shape unchanged
- define clearly what counts as observer-relevant change
- include interaction changes in the observer-relevant path if they materially affect passive usefulness
- measure the passive fanout ladder again after the policy change

## Initial Test Plan

### Server or measurement tests

- passive observer snapshots are skipped when the observer-oriented signature is unchanged
- passive observer snapshots still refresh on observer-relevant change
- fallback refresh still occurs after a deterministic interval
- passive fanout aggregate bytes/sec drops below the current observer-orientation baseline

### Contract tests

- none beyond the current contract validation, because shape should remain unchanged

### Integration tests

- one active client plus one passive observer can coexist with unchanged active behavior
- a passive observer still receives updates when meaningful observer state changes
- a passive observer still receives a slower fallback refresh under calm conditions

## Scenario Definition

Start one local server with one active client and several passive observers.

Scenario steps:

1. connect one active client and several passive observers
2. confirm the active client still receives the current local-detail path
3. confirm calm passive observers stop receiving repeated unchanged observer snapshots
4. confirm observer snapshots still arrive on observer-relevant change or fallback refresh
5. compare passive fanout bytes/sec to the current orientation-only observer baseline

## Done Criteria

- passive observer aggregate transport pressure drops again from the new observer-orientation baseline
- active-client behavior remains unchanged
- passive observers still remain useful and orienting
- no gameplay semantics or active transport shape changed

## Out Of Scope Follow-Ups

- compression
- delta encoding
- generalized subscriptions
- active-client transport changes
- gameplay changes
