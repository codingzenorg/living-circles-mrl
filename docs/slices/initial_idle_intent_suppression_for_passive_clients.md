# Slice: Initial Idle Intent Suppression For Passive Clients

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client movement-intent sender loop
- Go authoritative server active/passive connection handling
- current transport differentiation between active and passive clients

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay rules unchanged, but stops idle browsers from continuously reasserting themselves as active clients through neutral movement intents.

## Discovery Scope

Restore the intended passive-client transport path for browsers that are not actually moving.

This slice should:

- preserve authoritative movement rules unchanged
- preserve active-client responsiveness unchanged
- stop neutral `0,0` intent spam from keeping idle browsers marked as active
- make the existing passive transport differentiation meaningful during ordinary multi-client use

This slice does **not** attempt to implement:

- transport redesign
- client prediction or interpolation
- gameplay changes
- render redesign

## Why This Slice Next

The latest measurement and runtime evidence now point to a concrete likely cause of the observed two-browser slowdown:

- a second idle browser still feels expensive
- the client currently sends a movement intent every `100ms`
- that includes neutral `0,0` intents when the user is not moving
- the server records recent movement intent activity per connection to decide whether a client should receive the active or passive path

That means idle browsers likely keep themselves classified as active even when they are doing nothing.

The next useful step is not a broader transport redesign. It is to restore the intended active/passive distinction by suppressing neutral idle intent churn.

## Use-Case Contract

### Use Case

`SuppressIdleMovementIntent`

### Primary Actor

The browser client and authoritative server coordinating movement intent and transport mode.

### Pre-conditions

- active/passive transport differentiation already exists
- idle browsers currently still send neutral movement intents
- real two-browser slowdown has been observed

### Trigger

A connected browser has no effective movement direction.

### Success Outcome

- idle browsers stop refreshing their active status through neutral intents
- active players still receive full responsiveness while moving
- passive observers can actually use the cheaper transport path in ordinary multi-client use

### Failure Or Rejection Cases

- if active movement responsiveness drops, the slice failed
- if passive clients stop receiving needed initial or fallback state, the slice failed
- if the slice drifts into transport redesign, scope drifted

## Main Business Rules

1. Gameplay and transport contract shape remain unchanged.
2. Active movement must still feel immediate.
3. Neutral idle input must not keep a browser classified as active.
4. The existing passive-path transport policy should become effective during ordinary idle observation.

## Minimal Domain Concepts In Scope

- `Idle Client`
- `Neutral Movement Intent`
- `Active Client Classification`
- `Passive Transport Eligibility`

## Bounded Implementation Interpretation

This slice chooses the smallest useful mitigation:

- suppress redundant neutral movement intents from the browser when there is no effective direction change
- or otherwise ensure the server does not treat repeated neutral intents as recent activity
- leave transport message shape and gameplay rules unchanged

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser sender loop
- server movement-intent activity tracking
- implementation artifact updated with the resulting responsiveness implications

## Build Guidance

- prefer one clear active/passive classification rule over multiple heuristics
- keep movement behavior unchanged for active players
- validate with existing transport measurement helpers and targeted tests where practical

## Initial Test Plan

### Server or measurement tests

- prove idle clients do not keep full active cadence just by sending neutral intent
- prove active movers still keep near-full cadence

### Contract tests

- unchanged beyond current validation

### Integration tests

- add or update a websocket case only if needed to prove the active/passive distinction under neutral idle intent

## Scenario Definition

Run the current demo or deterministic harness with:

1. one actively moving client
2. one idle connected client
3. passive-path transport expected for the idle client
4. active-path responsiveness preserved for the mover

## Done Criteria

- neutral idle intent no longer defeats passive transport
- active movement responsiveness remains intact
- the two-client slowdown mitigation is evidence-backed

## Out Of Scope Follow-Ups

- broader transport redesign
- render redesign
- gameplay changes
- client prediction
