# Expectation Gap Detection

## Slice Under Review

`docs/slices/initial_active_transport_stop_point_reassessment.md`

## Evaluation Mode

Lightweight artifact-led EGD.

This pass reviews:

- the current slice definition
- the current implementation artifact
- recent transport measurement history
- fresh validation status

It does **not** claim an Ollama-backed scenario packet run.

## Evidence Reviewed

- [initial_active_transport_stop_point_reassessment.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/slices/initial_active_transport_stop_point_reassessment.md)
- [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- fresh validation:
  - `go test ./... -timeout 60s`
  - `npm run test:contracts`

## Current Built State

The active transport path is materially better than the earlier baseline:

- one-active aggregate bytes dropped from `8942` to `7524`
- two-active aggregate bytes dropped from `17884` to `13980`
- active orientation support is still the dominant active payload family, but much smaller than before
- active orientation usability still shows a bounded `2 fresh / 2 stale` split over the `300ms` movement window

The recent optimization history is also now explicit:

- active orientation-summary compaction was a real win
- active local-food compaction produced no measurable win
- active local-autonomous compaction produced no measurable win
- player-only precision had already produced no measurable win
- event-driven active-orientation refresh regressed the active path and was correctly rejected

## Main EGD Read

The active transport optimization track looks close to a sensible stopping point for now.

The current expectation gap is no longer:

- "the repo is still obviously wasting bytes in one easy active payload family"

The stronger gap is now:

- "the repo does not yet know whether remaining responsiveness pain is meaningfully caused by transport at all, or whether the next bottleneck is elsewhere"

That is a different kind of pressure.

The recent history shows that:

- first-order active transport wins were real
- later micro-optimizations are now mostly low-yield or non-winning
- continuing to iterate on tiny active payload families risks churn without meaningful responsiveness improvement

## Plausible Gaps

### 1. Stop-point judgment is now more important than another micro-optimization

The implementation evidence already suggests that recent active transport work is entering diminishing returns. A reasonable expectation is that the repo should now make an explicit stop-or-pivot decision instead of silently continuing the same optimization loop.

### 2. Responsiveness pressure may now live in another layer

The model hypothesis records real concurrent responsiveness risk. The current active transport baselines are better, but there is still no equally explicit evidence that these remaining user-visible slowdowns are dominated by transport rather than:

- server scheduling or fanout behavior at real play scale
- browser-side play-feel issues under multi-window use
- broader architectural tradeoffs between authoritative updates and perceived responsiveness

### 3. The repo may be near the limit of bounded compaction without a larger redesign

The current optimization track has stayed admirably bounded. The flip side is that the remaining wins may no longer be available through small compaction slices. If responsiveness still feels weak in real play, the next meaningful move may need to be:

- a larger transport redesign
- a scheduling/fanout slice
- or a deliberate pivot away from transport into measured play responsiveness

## What Looks Strong

- The repo did not force weak optimization ideas into history.
- Non-winning attempts were treated as useful evidence, not as failures to hide.
- The active path is much healthier than the earlier baseline.
- The optimization track remained disciplined and measurable.

## Recommendation

Return to `build`, but not for another tiny active payload compaction by default.

Best next directions:

1. build a stop-point reassessment artifact for the current active transport path
2. if responsiveness is still a problem after that, pivot to a new bottleneck-measurement slice instead of another micro-compaction
3. only return to active transport optimization if a clearer larger target emerges

## Conclusion

The expectation gap is no longer "active transport is still obviously too chatty."

The stronger gap is "the repo needs to decide whether active transport work should pause and responsiveness work should pivot elsewhere."
