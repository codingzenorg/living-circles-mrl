# Slice: Initial Active Orientation Usability Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using the reduced-cost active orientation-support path
- player-follow viewport plus minimap demo client
- Go authoritative server and browser client in their current integrated form

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport behavior unchanged, but checks whether the cheaper active orientation path still gives enough usable orientation during real play-speed movement.

## Discovery Scope

Measure the usability consequences of the recent active orientation reduction.

This slice should:

- preserve the current transport behavior unchanged
- preserve the current client behavior unchanged
- make active orientation freshness and staleness explicit during movement
- record whether the optimized path is still plausibly usable before any further transport cuts

This slice does **not** attempt to implement:

- another optimization
- another protocol change
- gameplay changes
- browser performance tuning

## Why This Slice Next

The current state now shows:

- active transport cost dropped materially
- the local-detail path still stays high-cadence
- orientation support is now refreshed more sparsely for active clients

That is a good optimization result, but it creates a new risk:

- transport got cheaper
- yet orientation usefulness may now degrade faster than the transport measurements reveal
- the next useful step is to validate the active experience before cutting more

## Use-Case Contract

### Use Case

`MeasureActiveOrientationUsabilityAfterTransportReduction`

### Primary Actor

The repository’s MRL loop evaluating whether the optimized active orientation path is still fit for use.

### Pre-conditions

- the reduced active orientation-support cadence is already built
- active-client transport measurements already show lower bytes/sec
- the browser client still uses cached minimap orientation between fresh updates

### Trigger

A bounded local evaluation run measures or records how often active orientation is stale during ordinary movement.

### Success Outcome

- the repo has explicit evidence about active orientation freshness under movement
- the next step can be chosen from evidence: keep, tune, or reverse
- no behavior changed during the measurement slice

### Failure Or Rejection Cases

- if the slice changes runtime behavior, it failed
- if the evidence is too vague to judge usability, it failed
- if the slice drifts into another optimization, scope drifted

## Main Business Rules

1. This is a measurement/evaluation slice, not an optimization slice.
2. Active transport behavior remains unchanged.
3. Client behavior remains unchanged.
4. The result must say whether active orientation freshness is still plausibly usable.

## Minimal Domain Concepts In Scope

- `Active Orientation Freshness`
- `Active Orientation Staleness`
- `Viewport Movement`
- `Play-Speed Usability`

## Bounded Measurement Interpretation

This slice chooses the smallest useful validation:

- measure or record how often active snapshots arrive with stale orientation during deterministic movement
- compare that with the current active cadence and minimap cache behavior
- record the result in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- existing transport measurement or integration helpers
- implementation artifact updated with the usability finding

## Build Guidance

- do not change runtime behavior
- prefer explicit deterministic evidence
- make the result actionable for the next loop: keep, tune, or revisit

## Initial Test Plan

### Server or measurement tests

- active orientation freshness/staleness ratio can be measured deterministically over a bounded movement window
- the measurement remains deterministic across repeated runs

### Contract tests

- none beyond the current contract validation, because runtime shape remains unchanged

### Integration tests

- active movement under the current optimized path can be observed without changing behavior

## Scenario Definition

Run one deterministic active client over a bounded movement window and record:

1. total active snapshots
2. snapshots with `orientation_fresh = true`
3. snapshots with `orientation_fresh = false`
4. whether local active detail remained continuous while orientation was stale

## Done Criteria

- active orientation usability is explicitly recorded
- the next step can be chosen from evidence rather than from transport bytes alone
- no gameplay or transport behavior changed

## Out Of Scope Follow-Ups

- further active transport reduction
- reverting the optimization
- passive transport changes
- browser rendering optimization
- gameplay changes
