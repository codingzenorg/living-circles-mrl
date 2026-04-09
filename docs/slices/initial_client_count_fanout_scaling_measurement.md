# Slice: Initial Client-Count Fanout Scaling Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- local websocket clients acting as passive observers
- Go authoritative server with the current optimized transport path
- shared runtime contract unchanged

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and transport unchanged. It extends the existing measurement track so the repository can compare how aggregate server output scales as client count increases under the same passive observer profile.

## Discovery Scope

Measure whether broad snapshot fanout is now the dominant remaining transport pressure.

This slice should:

- reuse the current deterministic multi-client measurement harness
- compare several bounded passive-observer client counts
- record aggregate bytes/sec, per-client bytes/sec, and cadence pressure at each level
- make fanout scaling explicit enough to justify or reject future fanout-oriented optimization

This slice does **not** attempt to implement:

- protocol redesign
- selective subscription logic
- observer-specific throttling
- compression
- gameplay changes

## Why This Slice Next

The repository now has:

- a one-client measurement baseline
- a bounded multi-client baseline
- an idle-versus-moving comparison

That last comparison produced the key signal:

- under the current bounded case, movement changes transport only slightly
- the stronger remaining pressure appears to be broad snapshot fanout itself

The new pressure is:

- find out how aggregate cost scales as passive observer count increases
- determine whether the current server still keeps cadence under higher fanout
- give future transport work a clearer justification than “multiple clients might be expensive”

This is the narrowest next step because it:

- reuses the current measurement harness
- stays in evidence-gathering rather than jumping into redesign
- directly tests the pressure the last slice surfaced

## Use-Case Contract

### Use Case

`MeasureClientCountFanoutScaling`

### Primary Actor

The repository maintainer evaluating whether passive client count is now the main transport pressure.

### Pre-conditions

- deterministic multi-client measurement already exists
- idle-versus-moving comparison already exists
- the current transport path is already optimized enough that further work should be evidence-led

### Trigger

A bounded measurement run executes the same passive-observer scenario at several explicit client counts.

### Success Outcome

- the repository has a deterministic fanout-scaling baseline
- aggregate transport growth by client count is explicit
- per-client stability and cadence pressure by client count are explicit
- later optimization can target fanout specifically, or stop if scaling remains acceptable

### Failure Or Rejection Cases

- if the run is not deterministic enough to compare counts meaningfully, the slice failed
- if the slice changes protocol behavior instead of measuring it, scope drifted
- if the count set becomes too broad or environment-dependent to compare sanely, scope exceeded

## Main Business Rules

1. This is a measurement slice, not a gameplay slice.
2. The transport path remains unchanged.
3. The scenario should stay passive and deterministic across client counts.
4. Aggregate bytes/sec, per-client bytes/sec, and cadence pressure must all be reported.
5. The result should make scaling trends explicit rather than only logging raw numbers.

## Minimal Domain Concepts In Scope

- `Client Count`
- `Fanout Scaling`
- `Aggregate Server Output`
- `Cadence Pressure Under Fanout`

## Bounded Measurement Interpretation

This slice chooses the smallest useful scaling study:

- one deterministic world
- one passive observer scenario
- a small explicit count ladder such as `1`, `4`, and `8`
- the same bounded observation window at each count

It avoids large load-testing infrastructure while still surfacing whether scaling remains roughly linear and acceptable.

## Required Runtime Contract Changes

None.

## Required Ports Or Boundaries

- websocket measurement harness
- implementation or measurement artifact recording the count ladder

## Build Guidance

- keep the count ladder small and explicit
- use the same deterministic world and same bounded window for each count
- report deltas and scaling trend, not only isolated measurements
- preserve easy comparison with the existing one-client and four-client results

## Initial Test Plan

### Server or measurement tests

- fanout-scaling measurement remains deterministic enough across repeated runs
- aggregate bytes/sec increases with client count in the expected direction
- per-client bytes/sec stays interpretable
- cadence pressure remains explicit at each count

### Contract tests

- none beyond the current contract validation

### Integration tests

- not required beyond the existing websocket measurement harness unless build needs a direct helper-level integration proof

## Scenario Definition

Start one local server and run the same passive observer scenario at several client counts.

Scenario steps:

1. measure `1` passive client
2. measure `4` passive clients
3. measure `8` passive clients
4. compare aggregate bytes/sec, per-client bytes/sec, and max inter-snapshot gap
5. record whether fanout scaling is now the clearest remaining transport pressure

## Done Criteria

- the repository has a deterministic fanout-scaling baseline
- client-count scaling is explicit enough to guide the next transport decision
- no gameplay semantics or protocol shape changed

## Out Of Scope Follow-Ups

- observer throttling
- transport sharding
- compression
- browser render profiling
