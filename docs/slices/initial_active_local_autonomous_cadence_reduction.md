# Slice: Initial Active Local Autonomous Cadence Reduction

## Selected Pack

`polyglot_client_server`

## Runtime Targets

- active-client local-detail transport under concurrent active browsers
- websocket snapshot fanout for local autonomous-circle detail
- current active player responsiveness path

## Architecture Mode

Explicit client/server boundary with the server remaining authoritative.

This slice keeps gameplay rules unchanged, but reduces one specific family inside the active local-detail path: local autonomous-circle detail does not need to refresh at the same cadence as player movement itself.

## Discovery Scope

Reduce the dominant remaining active payload fanout without touching passive observer behavior or active player control responsiveness.

This slice should:

- preserve gameplay rules unchanged
- preserve active player movement cadence unchanged
- preserve player detail and interaction detail unchanged
- keep passive transport behavior unchanged
- reduce only the cadence or refresh frequency of local autonomous-circle detail inside active snapshots

This slice does **not** attempt to implement:

- a full delta protocol
- compression
- gameplay changes
- client prediction or interpolation
- passive transport redesign

## Why This Slice Next

The last two slices now make the remaining active-path pressure explicit:

- post-idle-intent reassessment showed the idle second browser now truly reaches the passive path
- the new two-active timing measurement showed aggregate payload roughly doubles while inter-snapshot timing stays effectively flat near one tick
- earlier active component measurement already identified orientation support as one large cost, and that cost was reduced

That means the remaining slowdown is less likely to be immediate tick-loop collapse and more likely to be active payload fanout.

The next bounded move is to thin one expensive active payload family without touching the player’s own authoritative responsiveness. Local autonomous-circle detail is the cleanest target:

- it is useful at high cadence, but less critical than player motion and interaction state
- it can be refreshed less often while the client reuses the last valid local autonomous set
- it avoids reopening player control semantics or passive-path logic

## Use-Case Contract

### Use Case

`ReduceActiveLocalAutonomousFanout`

### Primary Actor

The authoritative server sending active local-detail snapshots to active browser clients.

### Pre-conditions

- active local-detail transport is already established
- passive-path transport is already functioning for idle observers
- two-active-client pressure now reads more like payload fanout than tick collapse

### Trigger

An active browser continues receiving local-detail snapshots under ordinary play.

### Success Outcome

- active player control and player-centric responsiveness remain intact
- local autonomous detail no longer refreshes on every active tick when unchanged
- two-active-client aggregate active payload drops from the current baseline

### Failure Or Rejection Cases

- if player responsiveness drops, the slice failed
- if local autonomous state drifts too long or becomes unusable, the slice failed
- if the slice expands into a broader delta/compression redesign, scope drifted

## Main Business Rules

1. Gameplay and authoritative server semantics remain unchanged.
2. The active player path must remain near-full cadence.
3. The player’s own detail and interaction detail remain high priority.
4. Local autonomous-circle detail may refresh less often or on change plus fallback.
5. The client may reuse the last valid local autonomous set between fresh active-autonomous refreshes.

## Minimal Domain Concepts In Scope

- `Active Client`
- `Local Autonomous Detail`
- `Player-Critical Responsiveness`
- `Active Payload Fanout`
- `Fresh Autonomous Detail`

## Bounded Implementation Interpretation

This slice chooses the smallest useful mitigation:

- keep the existing active snapshot shape where practical
- introduce one explicit freshness rule for local autonomous detail
- reuse the last valid local autonomous set on the client between refreshes
- measure the resulting active-path reduction against the current two-active baseline

## Required Runtime Contract Changes

Likely minimal:

- an explicit freshness signal for local autonomous detail if omission is used

## Required Ports Or Boundaries

- server active snapshot assembly
- browser active snapshot reuse path
- measurement helpers and implementation artifact updated with the new baseline

## Build Guidance

- prefer one clear freshness rule over several heuristics
- keep active player responsiveness unchanged
- validate against the current one-active and two-active baselines

## Initial Test Plan

### Server or measurement tests

- prove active local autonomous detail can be omitted or thinned while player detail remains full cadence
- prove two-active aggregate payload drops below the current `17884` byte baseline over `300ms`

### Contract tests

- extend only if a minimal freshness field is introduced

### Integration tests

- add or update a websocket case only if needed to prove active clients still recover fresh autonomous detail correctly

## Scenario Definition

Run the deterministic active-client harness with:

1. one active client baseline
2. two active clients baseline
3. reduced active autonomous-detail cadence or change-driven refresh
4. record the post-change active aggregate bytes and cadence outcome

## Done Criteria

- active player responsiveness remains intact
- active local autonomous detail no longer contributes full cost every tick by default
- the two-active aggregate baseline drops materially below the current measurement

## Out Of Scope Follow-Ups

- full delta protocol
- compression
- gameplay changes
- passive transport redesign
