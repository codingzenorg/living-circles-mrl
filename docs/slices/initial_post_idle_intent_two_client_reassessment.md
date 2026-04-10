# Slice: Initial Post Idle Intent Two Client Reassessment

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client movement-intent sender loop after idle suppression
- Go authoritative server active/passive connection handling under two browsers
- current transport measurement and responsiveness evidence path

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged. It only reassesses the observed two-browser slowdown now that idle browsers no longer churn neutral movement intent.

## Discovery Scope

Measure whether the recent idle-intent suppression actually restores the intended passive-path behavior strongly enough to improve ordinary two-browser responsiveness.

This slice should:

- preserve gameplay rules unchanged
- preserve transport shape unchanged
- reuse the current active/passive transport policy as built
- compare the new two-browser baseline against the prior slowdown evidence
- record whether the remaining pressure is still mainly server tick/broadcast fanout

This slice does **not** attempt to implement:

- another transport redesign
- client render redesign
- gameplay changes
- prediction or interpolation

## Why This Slice Next

The repo now has a concrete mitigation for a likely cause of the observed slowdown:

- idle browsers no longer send repeated neutral `0,0` intents
- unchanged repeated movement directions are suppressed
- one explicit stop intent still preserves active control semantics

That should make the passive transport path reachable for truly idle browsers. But the real question is no longer theoretical. It is whether the original “second browser makes movement slow” signal still appears after this mitigation.

The next useful step is therefore not another optimization guess. It is to remeasure and reassess the two-browser case using the new idle-intent baseline.

## Use-Case Contract

### Use Case

`ReassessTwoClientResponsiveness`

### Primary Actor

The browser clients and authoritative server under ordinary concurrent use.

### Pre-conditions

- idle-intent suppression is already built
- active/passive transport differentiation already exists
- earlier runtime evidence recorded perceptible slowdown with a second browser

### Trigger

Two browsers connect after the idle-intent suppression slice has landed.

### Success Outcome

- the repository has a clear post-fix two-browser responsiveness reading
- the next mitigation, if any, is based on that new evidence
- if the slowdown is materially reduced, the repo avoids unnecessary further optimization guesses

### Failure Or Rejection Cases

- if the slice drifts into implementation changes, scope drifted
- if the measurement cannot distinguish idle-second-browser behavior from active-second-browser behavior, the slice is incomplete

## Main Business Rules

1. Runtime behavior remains unchanged in this slice.
2. The recent idle-intent suppression is treated as the new baseline.
3. The two-browser case should distinguish at least:
   - one active browser plus one idle browser
   - two simultaneously active browsers, if practical
4. The result should clarify whether the remaining pressure is:
   - mostly resolved
   - still mainly passive/observer related
   - still mainly active fanout or tick/broadcast related

## Minimal Domain Concepts In Scope

- `Idle Browser`
- `Active Browser`
- `Passive Transport Reachability`
- `Two Client Responsiveness`
- `Post Mitigation Baseline`

## Bounded Implementation Interpretation

This slice chooses the smallest useful next step:

- reuse current measurement helpers where possible
- add only the minimal measurement or artifact updates needed to reassess the two-browser case
- record the result in implementation evidence

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser/runtime evidence path
- server transport measurement path
- implementation artifact updated with the reassessment result

## Build Guidance

- prefer deterministic measurement over intuition
- treat single-client and two-client comparisons as the baseline frame
- if helpful, distinguish idle-second-client and active-second-client cases explicitly

## Initial Test Plan

### Server or measurement tests

- prove the measurement path can distinguish the post-fix two-client cases deterministically
- record whether passive-path reachability improved the idle-second-client case

### Contract tests

- unchanged beyond current validation

### Integration tests

- reuse existing transport responsiveness proofs unless one focused new case is required

## Scenario Definition

Run the current demo or deterministic harness with:

1. one active moving browser
2. one second idle browser using the new idle-intent suppression behavior
3. optionally, one comparison case with both browsers active
4. record whether responsiveness and snapshot pressure improved relative to the prior evidence

## Done Criteria

- post-fix two-browser responsiveness is explicitly recorded
- the repo knows whether the idle-intent mitigation was sufficient
- the next mitigation, if any, is grounded in the new result

## Out Of Scope Follow-Ups

- new transport redesign
- render redesign
- gameplay changes
- client prediction
