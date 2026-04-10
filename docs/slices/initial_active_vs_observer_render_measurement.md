# Slice: Initial Active Vs Observer Render Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in ordinary viewport mode
- active player transport path and `observer_orientation_only` path
- current server and client behavior unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and rendering behavior unchanged, but measures whether observer-mode render cost is genuinely higher than active-mode render cost or whether the difference is mostly browser scheduling noise.

## Discovery Scope

Measure the current render-pressure difference between an active tab and an observer tab under the current two-browser model.

This slice should:

- preserve current runtime behavior unchanged
- preserve current transport shape unchanged
- preserve current render instrumentation unchanged unless a narrow comparison helper is needed
- make the active-versus-observer render difference explicit enough to guide the next step

This slice does **not** attempt to implement:

- another transport optimization
- another render optimization
- gameplay changes
- browser-priority workarounds

## Why This Slice Next

The current runtime now exposes a new ambiguity:

- observer mode is transport-lighter than active mode
- observer mode now still renders a live world view instead of a blank or frozen scene
- in ordinary manual use, the observer tab’s `render/world` indicator appears roughly twice as high as the active tab’s

That number is not automatically evidence of a real render asymmetry:

- observer mode still draws most of the same scene
- browser tab scheduling can skew local draw timing
- the current instrumentation does not yet say whether the difference is stable, meaningful, or mostly incidental

The next useful step is therefore measurement, not another optimization guess.

## Use-Case Contract

### Use Case

`MeasureActiveVsObserverRenderPressure`

### Primary Actor

The repository’s MRL loop evaluating whether observer-mode render pressure is a real client-side cost difference.

### Pre-conditions

- observer mode now remains live and visually informative
- the current render instrumentation already exposes top-level and subcomponent timings
- a manual two-browser observation suggested higher observer-side render timing

### Trigger

A bounded local comparison run records render-pressure readings for an active tab and an observer tab under the same current runtime.

### Success Outcome

- the repo can say whether observer-mode render pressure is materially different from active-mode render pressure
- the repo can say whether the difference looks more like real draw work or likely browser scheduling noise
- no runtime behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the result still cannot distinguish active versus observer render pressure meaningfully, it failed
- if the slice drifts into optimization, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Active and observer runtime behavior remain unchanged.
3. The result must compare active and observer render pressure directly enough to guide the next decision.
4. The result should explicitly separate local render work from likely scheduling bias when possible.

## Minimal Domain Concepts In Scope

- `Active Render Pressure`
- `Observer Render Pressure`
- `World Draw Cost`
- `Render Timing Bias`
- `Two-Tab Comparison`

## Bounded Measurement Interpretation

This slice chooses the smallest useful comparison:

- use one active tab and one observer tab against the same current server
- record bounded render readings for both roles over a short comparable movement window
- summarize whether the observer difference appears stable enough to justify a render optimization slice

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser-side measurement only
- implementation artifact updated with the comparison result

## Build Guidance

- do not change gameplay, transport, or rendering behavior
- reuse the current render instrumentation when possible
- prefer a concise comparison result over a new profiling system

## Initial Test Plan

### Server or measurement tests

- add a bounded comparison helper only if the current instrumentation is insufficient

### Contract tests

- unchanged beyond current validation

### Integration tests

- only if a new bounded comparison helper crosses the current client/server boundary

## Scenario Definition

Run the current demo with:

1. one active tab
2. one observer tab
3. a short comparable movement window
4. recorded active and observer render readings, especially total render and `world` cost

## Done Criteria

- the active-versus-observer render difference is more explicit than the current anecdote
- the next step can be chosen from evidence rather than from one suspicious live reading
- no gameplay, transport, or render behavior changed

## Out Of Scope Follow-Ups

- render optimization
- transport optimization
- gameplay changes
- browser scheduling workarounds
