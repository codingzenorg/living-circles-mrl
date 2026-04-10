# Slice: Initial Active Transport Component Measurement

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active websocket clients using `active_local_detail`
- passive observer path already optimized and measured separately
- Go authoritative server with the current local-detail transport

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay and protocol shape unchanged, but makes the active transport payload legible by breaking down which active-detail components dominate serialized cost.

## Discovery Scope

Measure the active payload composition before optimizing it.

This slice should:

- preserve the current active-client transport path unchanged
- preserve the current passive observer path unchanged
- measure how much of the active payload comes from player detail, local autonomous detail, local food detail, interaction detail, and orientation support
- record a deterministic active payload breakdown that can justify the next optimization slice

This slice does **not** attempt to implement:

- compression
- delta encoding
- cadence changes
- gameplay changes
- transport field removal

## Why This Slice Next

The current measurements now show:

- passive observer fanout is much cheaper than before
- active fanout scales almost linearly with client count
- the active `active_local_detail` path is now the clearest remaining transport pressure

That is enough to justify focusing on the active path, but not enough to justify a specific optimization yet.

The new pressure is:

- identify which active payload components actually dominate bytes
- avoid optimizing a minor field while leaving the main active cost untouched
- preserve active responsiveness until the next slice is evidence-backed

That points to a bounded measurement slice:

- measure the current active snapshot with selected components omitted or isolated in a deterministic way
- compare the resulting serialized sizes
- leave gameplay and runtime contract behavior unchanged

## Use-Case Contract

### Use Case

`MeasureActiveTransportPayloadComponents`

### Primary Actor

The Go transport layer inspecting the serialized cost of its current active-client payload.

### Pre-conditions

- active fanout scaling has already been measured
- active clients still use the current full local-detail path
- the next optimization target should be chosen from evidence rather than intuition

### Trigger

A deterministic measurement run is started to compare the serialized size contribution of active payload components.

### Success Outcome

- the dominant active payload components are explicit
- future optimization can target the right active-detail area
- no gameplay or protocol behavior changed

### Failure Or Rejection Cases

- if protocol behavior changes, the slice failed
- if the breakdown is too coarse to guide a next optimization, the slice failed
- if the slice turns into an optimization rather than a measurement, scope drifted

## Main Business Rules

1. This is a measurement slice, not a gameplay slice.
2. Active transport behavior remains unchanged.
3. Passive transport behavior remains unchanged.
4. The breakdown should isolate the major active payload families, not every individual field.
5. The output should be easy to compare with the existing active fanout baseline.

## Minimal Domain Concepts In Scope

- `Active Transport Payload`
- `Player Detail Cost`
- `Local Autonomous Detail Cost`
- `Local Food Detail Cost`
- `Orientation Support Cost`
- `Interaction Cost`

## Bounded Measurement Interpretation

This slice chooses the smallest useful active breakdown:

- take the current deterministic active snapshot
- measure serialized size with selected major components removed or isolated
- record the comparative results in the implementation artifact

## Required Runtime Contract Changes

None expected.

## Required Ports Or Boundaries

- deterministic transport measurement helpers
- implementation artifact updated with the active payload breakdown

## Build Guidance

- keep runtime behavior unchanged
- avoid inventing a new transport shape
- focus on major payload families, not tiny field-level accounting
- make the results directly actionable for the next optimization slice

## Initial Test Plan

### Server or measurement tests

- active payload component measurement is deterministic
- the full active payload exceeds any major isolated subset
- at least one dominant active payload component is clearly identified

### Contract tests

- none beyond the current contract validation, because runtime shape remains unchanged

### Integration tests

- none required unless a helper crosses a runtime boundary

## Scenario Definition

Take a deterministic active snapshot from the current default world and measure:

1. full active snapshot size
2. active snapshot without local autonomous detail
3. active snapshot without local food detail
4. active snapshot without orientation support
5. active snapshot without interaction detail when absent

Compare the resulting sizes and record the dominant cost families.

## Done Criteria

- the dominant active payload families are explicit
- the next active transport optimization target is evidence-backed
- no gameplay or protocol behavior changed

## Out Of Scope Follow-Ups

- compressing payloads
- changing active cadence
- removing fields from the runtime contract
- gameplay changes
