# Slice: Initial Post-Reduction Active Transport Reassessment

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using the current reduced-cost active transport path
- Go authoritative server and browser client in their current integrated form
- current passive observer path unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport behavior, and rendering behavior unchanged, but remeasures the active transport payload after the recent cadence and freshness reductions so the next optimization target is chosen from current evidence instead of stale pre-reduction numbers.

## Discovery Scope

Reassess the current active payload composition after the recent active-path reductions.

This slice should:

- preserve current runtime behavior unchanged
- preserve current transport shape unchanged
- preserve current browser behavior unchanged
- make the dominant remaining active payload family explicit in the current repo state

This slice does **not** attempt to implement:

- another transport optimization
- another render optimization
- gameplay changes
- protocol redesign

## Why This Slice Next

The transport path has changed materially since the earlier active component measurement:

- active orientation support is no longer sent every tick
- local autonomous detail is now refreshed less often under real fanout
- local foods already use change-plus-fallback refresh
- observer-mode fixes clarified responsiveness behavior without re-expanding the active path

That means the older “orientation is dominant” read may no longer be the right next decision basis. The repo needs one fresh active payload breakdown in the current state before choosing another optimization slice.

## Use-Case Contract

### Use Case

`MeasureCurrentActiveTransportComposition`

### Primary Actor

The repository’s MRL loop choosing the next active-path responsiveness target.

### Pre-conditions

- the reduced active transport path is already built
- recent measurements showed payload fanout still matters under two active clients
- prior active component numbers were captured before the latest active-path reductions

### Trigger

A bounded local measurement run compares the current full active snapshot against current major active payload subsets.

### Success Outcome

- the dominant remaining active payload family is explicit in the current repo state
- the next active-path optimization can target the current dominant family rather than an outdated one
- no runtime behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the result does not clarify the current active payload composition, it failed
- if the slice drifts into implementing another optimization, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Active transport behavior remains unchanged.
3. The result must reflect the current post-reduction active path, not an earlier baseline.
4. The result should be specific enough to justify the next optimization target.

## Minimal Domain Concepts In Scope

- `Current Active Payload`
- `Post-Reduction Transport Composition`
- `Dominant Remaining Active Family`
- `Two-Active Responsiveness Pressure`

## Bounded Measurement Interpretation

This slice chooses the smallest useful reassessment:

- reuse the existing active transport component measurement shape when possible
- run it against the current transport implementation
- record which current family dominates after the recent reductions

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- existing server-side measurement helpers
- implementation artifact updated with the reassessment result

## Build Guidance

- do not change transport behavior
- prefer updating existing measurement helpers over creating a parallel system
- make the output directly comparable to the earlier active component measurement

## Initial Test Plan

### Server or measurement tests

- active component reassessment remains deterministic
- the result reflects the current transport implementation

### Contract tests

- unchanged beyond current validation

### Integration tests

- none required unless the measurement helper crosses a new boundary

## Scenario Definition

Measure the current active transport payload and compare:

1. full current active snapshot
2. current snapshot without player detail
3. current snapshot without local autonomous detail
4. current snapshot without local food detail
5. current snapshot without current orientation support
6. current snapshot without interaction detail

## Done Criteria

- the current dominant active payload family is explicit
- the next active transport slice can be chosen from current evidence
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- transport optimization
- render optimization
- gameplay changes
- protocol redesign
