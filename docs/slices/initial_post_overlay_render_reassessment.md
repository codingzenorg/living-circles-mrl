# Slice: Initial Post Overlay Render Reassessment

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- browser client in viewport mode with current live render-pressure plus world and overlay render breakdowns
- current reduced-cost overlay path after glow, offscreen, cue, and recent-effect narrowing
- Go authoritative server unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay, transport, and visual behavior unchanged, but re-establishes the current dominant browser render pressure after the recent sequence of bounded client-side reductions.

## Discovery Scope

Reassess the client render path after the recent overlay optimizations.

This slice should:

- preserve current runtime behavior unchanged
- preserve transport behavior unchanged
- use the existing live client render instrumentation
- record which top-level and subfamily buckets now appear dominant after the recent reductions

This slice does **not** attempt to implement:

- further render optimization
- gameplay changes
- transport changes
- UI redesign

## Why This Slice Next

The repository has now applied several bounded client render optimizations in sequence:

- label cost reduction
- glow overlay cost reduction
- offscreen overlay cost reduction
- cue overlay cost reduction
- recent-effect overlay cost reduction

That means the old optimization priorities may no longer hold. Continuing to trim client drawing without reassessment would become speculative again.

The next useful step is therefore:

- remeasure the current dominant render buckets
- decide whether the next work still belongs in rendering
- or stop optimizing and return to broader evaluation

## Use-Case Contract

### Use Case

`ReassessClientRenderPressure`

### Primary Actor

The browser client measuring its own current draw path.

### Pre-conditions

- live total render pressure already exists
- world and overlay subfamily measurements already exist
- several bounded render-cost reductions have already landed

### Trigger

A bounded local render measurement run reviews the current client render path after the recent reductions.

### Success Outcome

- the current dominant render families are explicit again
- the next loop can decide between another render slice, EGD, or stopping the optimization track
- no runtime behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the result does not clarify the new dominant render pressure, it failed
- if the slice drifts into another optimization, scope drifted

## Main Business Rules

1. This is a measurement slice, not an optimization slice.
2. Default runtime behavior remains unchanged.
3. The result should clarify the new dominant render families after the recent reductions.
4. The result should be simple enough to justify whether the optimization track should continue.

## Minimal Domain Concepts In Scope

- `Post-Reduction Render Pressure`
- `Dominant Render Family`
- `Optimization Reassessment`
- `Browser Draw Baseline`

## Bounded Measurement Interpretation

This slice chooses the smallest useful reassessment:

- use the existing live render instrumentation
- review top-level, world, and overlay buckets after the recent reductions
- record the resulting dominant families in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- browser client measurement only
- implementation artifact updated with the new dominant-render read

## Build Guidance

- do not change gameplay, transport, or current visuals
- prefer a concise recorded conclusion over more instrumentation
- use the current measurement path instead of adding a new profiling system

## Initial Test Plan

### Server or measurement tests

- existing server and contract validation only, because no behavior changes are expected

### Contract tests

- unchanged beyond current validation

### Integration tests

- none required unless the measurement path changes unexpectedly

## Scenario Definition

Run the current viewport client after the recent overlay reductions and record:

1. the current dominant top-level render family
2. the current dominant world subfamily when relevant
3. the current dominant overlay subfamily when relevant
4. whether another render optimization slice is still justified

## Done Criteria

- the current dominant render pressure is explicit again
- the next loop can justify whether to continue or stop client render optimization
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- another render optimization by default
- transport optimization
- gameplay changes
- visual redesign
