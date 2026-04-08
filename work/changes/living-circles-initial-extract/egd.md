# Expectation Gap Detection

## Scope

- reviewed against [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- reviewed against [domain_background_knowledge.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/domain_background_knowledge.md)
- reviewed against the current slice [initial_event_driven_local_food_transport.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/slices/initial_event_driven_local_food_transport.md)
- reviewed against the current implementation summary in [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- validated with `go test ./...` and `npm run test:contracts` on April 8, 2026

## Summary

This slice is a good optimization slice, and it succeeds on its own terms. The transport path is now materially less chatty than the earlier full-snapshot baseline, and the current optimization ladder remains coherent:

- viewport culling
- lower-cadence orientation summaries
- compact minimap summaries
- reduced local precision
- event-driven orientation refresh
- event-driven local food refresh

That is a real systems improvement. The browser client/server boundary now better matches the genre expectation of an authoritative server without wasting as much bandwidth on obviously static payload.

But this EGD also makes a useful distinction: this slice improves transport efficiency, not domain meaning. It reduces overhead and helps responsiveness headroom, yet it does not materially advance the ecological or lineage distinctiveness questions that earlier EGDs identified. So the main gap is no longer “this optimization is missing”; the main gap remains what kind of world behavior should be strengthened next.

## Main Findings

### 1. The transport path is now meaningfully more disciplined

This is the strongest positive finding.

The system no longer treats all snapshot data as equally deserving of the same cadence. It now distinguishes between:

- local circles, which remain high-cadence and interaction-relevant
- orientation summaries, which are now event-driven with fallback refresh
- local food detail, which is now also event-driven with fallback refresh

That better matches actual information volatility. In genre terms, this is a healthier authoritative-update model than a uniformly chatty full-state stream.

### 2. The client contract remains understandable despite the optimization

This matters because transport optimization often creates brittle or opaque state-sync behavior.

The current design avoided that failure mode:

- stale ticks are explicit through freshness flags
- the browser client reuses the last valid data instead of inventing prediction
- the transport still looks like one readable snapshot shape rather than a generalized delta protocol

That keeps inspectability intact, which is aligned with the lab’s architectural goals.

### 3. Responsiveness headroom likely improved, but the project still lacks a stronger runtime-pressure evaluation artifact

The transport cost trend is clearly moving in the right direction. That helps the domain background expectation of low perceived latency and trustworthy authority.

But this EGD can only infer likely responsiveness gains from payload reductions and passing tests. It still does not have:

- a sustained multi-client pressure scenario
- a browser-side frame/jank observation artifact
- a measured end-to-end latency or update-staleness report

So the optimization is plausible and valuable, but still not fully evaluated under stronger real-time pressure.

### 4. The slice improves technical coherence more than ecosystem coherence

This is not a criticism of the slice; it is a scope observation.

Relative to the model hypothesis, this transport work:

- supports the intended client/server architecture
- supports authoritative-server fairness
- supports real-time browser play

But it does not change:

- regional ecological consequence
- long-horizon population dynamics
- lineage significance
- autonomous niche formation

So the major semantic gaps identified by earlier EGDs still remain the main gaps after this slice.

### 5. The protocol is now approaching a point where blind optimization becomes less attractive than targeted measurement

The current optimization path has stayed disciplined because each step had a clear local pressure and measurable effect.

That is still good. But the slice history suggests a new risk:

- continuing to optimize transport mechanically without a stronger pressure model

The next transport slice should probably come from one of two clearer pressures:

- real observed runtime pressure, such as multi-client load or visible browser lag
- a remaining large redundant category in the current transport shape

Without that, additional optimization may start chasing smaller wins while semantic work remains more valuable.

## Secondary Findings

- The current transport strategy still preserves a strong authoritative-server feel without introducing prediction complexity.
- The minimap-oriented worldview remains coherent after these optimizations; the client still has enough orientation data to navigate the larger world.
- The protocol is becoming more asymmetric by volatility, which is a good sign that the boundary is being shaped by meaning rather than by naive completeness.

## Main Expectation Gaps

1. Transport efficiency is improving, but runtime-pressure evaluation is still weaker than the optimization work itself.
2. The largest remaining gaps are still ecological and experiential, not raw transport shape.
3. The current build still lacks stronger evidence of long-horizon population dynamics and lineage consequence.
4. Further transport work should now be justified by either measured pressure or a clearly dominant remaining redundant payload class.

## Recommendation

Return to `refine`, but choose deliberately between two directions instead of continuing transport work by default.

### Direction A: one more bounded transport measurement or optimization slice

Choose this only if the goal is still network/runtime headroom.

Best candidates:

1. measure or simulate multi-client transport pressure
2. identify whether high-cadence local circle detail is now the dominant remaining cost
3. only then decide whether a new transport slice is still worth it

### Direction B: return to ecological or experiential meaning

Choose this if the goal is still to strengthen Living Circles as a model rather than mainly as a transport system.

Best candidates:

1. a longer-run population-dynamics evaluation slice
2. a stronger lineage-consequence slice
3. a play-feel evaluation slice focused on responsiveness and readability in the larger world

## Best Next Slice Directions

### Option A: Multi-client transport pressure measurement

Add a bounded measurement slice that estimates bytes/sec and snapshot behavior across several simultaneous clients under the current optimized transport path.

Why this fits:

- it would tell you whether more transport optimization is truly justified
- it keeps scope smaller than a new protocol redesign
- it converts current intuition about lag into evidence

### Option B: Longer-run ecosystem evaluation

Add a bounded evaluation slice that observes a larger world for longer runs and records whether current ecological pressures produce visible regional dominance, collapse, recovery, or stable niches.

Why this fits:

- the main semantic gaps are still ecological
- the current build now has enough structural pressure to justify stronger observation
- it can reveal whether gameplay/model work is more urgent than more transport work

### Option C: Player responsiveness evaluation

Add a bounded experience-oriented evaluation slice that records whether the larger map, viewport, minimap, and current transport policy still feel responsive and legible at play speed.

Why this fits:

- the domain background explicitly says responsiveness matters
- transport optimization is currently being used as a proxy for responsiveness
- a direct evaluation artifact would reduce that guesswork

## Return-To-Loop Recommendation

Recommended next phase: `refine`

Recommended intent:

- if the immediate concern is lag, refine a measurement slice before more transport redesign
- otherwise, switch back toward ecosystem or play evaluation instead of continuing optimization by inertia
