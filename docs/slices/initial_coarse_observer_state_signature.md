# Slice: Initial Coarse Observer State Signature

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- passive websocket observers already using `observer_orientation_only`
- event-driven passive observer refresh already skipping calm repeated snapshots
- Go authoritative server with the current measured fanout path

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and active transport unchanged, but makes passive observer refresh react to coarse whole-world state changes that are meaningful for an observer without reintroducing the old minimap-motion churn.

## Discovery Scope

Restore useful passive observer updates after the aggressive event-driven reduction.

This slice should:

- preserve the current active-client path unchanged
- preserve the current observer-oriented snapshot representation
- keep calm passive fanout much lower than the prior observer-orientation baseline
- refresh passive observers when coarse observer-relevant world state changes, even if there is no explicit interaction event

This slice does **not** attempt to implement:

- compression
- delta encoding
- generalized subscriptions
- gameplay changes
- active-client cadence reduction
- exact minimap-motion-driven passive refresh

## Why This Slice Next

The new passive observer policy now proves that calm repeated observer snapshots can be removed aggressively. That is useful, but it also exposes the next coherence gap:

- passive observers now mainly refresh on interaction or fallback
- broader state changes like total food loss, food return, or population loss can remain stale until fallback
- reusing the full minimap summary as the observer signature would restore unnecessary refreshes caused by ordinary motion

That makes the next pressure explicit:

- observer relevance needs to include coarse whole-world state, not just explicit interaction outcomes
- passive refresh must remain materially cheaper than the prior observer-orientation baseline
- the next bounded step is a coarse observer state signature, not a return to exact orientation churn

## Use-Case Contract

### Use Case

`RefreshIdleObserverTransportOnCoarseObserverStateChange`

### Primary Actor

The Go server broadcasting authoritative snapshots to active and passive websocket clients.

### Pre-conditions

- passive observers already use lower-cadence event-driven `observer_orientation_only` transport
- active clients already keep the local interactive path
- passive fanout is already much lower than the original chatty baseline

### Trigger

A passive observer reaches a potential send tick and the server decides whether a coarse observer-relevant world state changed since the last observer refresh.

### Success Outcome

- active clients remain unchanged
- passive observers stay quiet during ordinary calm motion
- passive observers refresh when coarse whole-world state changes in a meaningful way
- passive fanout remains below the earlier orientation-only passive baseline

### Failure Or Rejection Cases

- if active responsiveness changes, the slice failed
- if passive observers still refresh repeatedly on ordinary minimap drift, the slice failed
- if broad world-state changes remain invisible until fallback, the slice failed

## Main Business Rules

1. This is an observer-signature slice, not a gameplay slice.
2. Active clients keep the current local-detail cadence and payload.
3. Passive observers keep the current orientation-only representation.
4. Passive observer refresh should react to coarse observer-relevant world state changes.
5. Ordinary motion inside the same coarse observer state should not force passive refresh.
6. A deterministic fallback refresh must still exist.

## Minimal Domain Concepts In Scope

- `Passive Observer`
- `Observer Refresh Signature`
- `Coarse Observer State`
- `Observer Fallback Refresh`

## Bounded Optimization Interpretation

This slice chooses the smallest useful observer-signature refinement:

- derive a deterministic coarse observer signature from low-churn observer-relevant state
- include interaction changes and coarse whole-world counts or grouped state changes
- avoid exact minimap-position churn as a refresh trigger

## Required Runtime Contract Changes

None expected.

The current explicit `transport_mode` plus the existing observer-oriented snapshot fields should remain sufficient if only the server’s passive refresh decision changes.

## Required Ports Or Boundaries

- passive observer refresh-signature derivation in the Go transport layer
- deterministic transport measurement harness updated to record the new passive baseline
- implementation artifact updated with the new measured ladder

## Build Guidance

- keep active-client behavior unchanged
- keep the observer-oriented snapshot shape unchanged
- include at least interaction changes and coarse food or population state changes in the observer signature
- avoid exact minimap coordinate drift as a passive refresh trigger
- remeasure passive fanout after the change

## Initial Test Plan

### Server or measurement tests

- passive observer refresh stays quiet when only ordinary motion changes inside the same coarse observer state
- passive observer refresh occurs when coarse food or population state changes
- fallback refresh still occurs after the deterministic interval
- passive fanout remains below the earlier orientation-only passive baseline

### Contract tests

- none beyond the current contract validation, because shape should remain unchanged

### Integration tests

- one active client plus one passive observer can coexist with unchanged active behavior
- a passive observer receives an update when a coarse observer-relevant world change occurs
- a passive observer still receives a slower fallback refresh under calm conditions

## Scenario Definition

Start one local server with one active client and several passive observers.

Scenario steps:

1. connect one active client and several passive observers
2. confirm the active client still receives the current local-detail path
3. confirm ordinary motion alone does not force repeated passive refreshes
4. confirm coarse world-state changes still refresh passive observers
5. compare passive fanout bytes/sec to the prior observer-orientation baseline

## Done Criteria

- passive observer transport stays quiet under calm motion
- passive observers still learn about meaningful coarse world-state changes before fallback
- active-client behavior remains unchanged
- no gameplay semantics or transport shape changed

## Out Of Scope Follow-Ups

- compression
- delta encoding
- generalized subscriptions
- active-client transport changes
- gameplay changes
