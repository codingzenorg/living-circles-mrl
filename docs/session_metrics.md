# Session Metric Guidance

## Purpose

This document gives this repository a simple way to track session productivity during MRL work.

Use it as a lightweight operational hint, not as a scoring system. The goal is to notice whether sessions remain controlled, grounded, and productive over time.

---

## Primary Metric

* `output_efficiency = output_tokens / total_tokens`

### Interpretation

* ideal range: `5% - 10%`
* below `5%`: investigate context size, loop friction, or low progress per session
* above `15%`: investigate drift risk, weak grounding, or slice size

### Tracking Intent

* do not maximize output
* prefer input-heavy, output-controlled sessions
* optimize for small, controlled changes, consistent progress, and stable reasoning

### Operational Rule

* if `output_efficiency` stays between `5%` and `10%`, the session pattern is healthy
* if it goes outside that range, inspect context size or slice size before continuing in the same pattern

---

## Supporting Metrics

These metrics are optional but useful when the tooling exposes them.

* `reasoning_efficiency = reasoning_tokens / output_tokens`
* `cache_ratio = cached_tokens / input_tokens`

Interpret them conservatively:

* reasoning efficiency helps indicate whether the session is spending a reasonable amount of explicit reasoning relative to produced output
* cache ratio helps show whether the workflow is benefiting from stable repeated context

These numbers are hints, not targets.

---

## How To Use This In MRL

Track metrics per working session or per short delivery interval.

Review them when:

* a session feels expensive without much progress
* slice size may be too large
* context keeps growing across loops
* the repository wants a lightweight historical view of delivery efficiency

Do not use these metrics to replace artifact review, commit review, or slice evaluation.
They are operational hints for the shape of the work, not truth about the value of the work.

---

## Session - 20260405

### Raw

* total_tokens: 12280777
* input_tokens: 11540299
* cached_tokens: 233647872
* output_tokens: 740478
* reasoning_tokens: 112942

### Derived

* output_efficiency: 6.03%
* reasoning_efficiency: 15.25%
* cache_ratio: 20.25x

### Productivity

* slices_applied: 61
* tokens_per_slice: 201324
* docs_commits: 67

### Notes

* historical measurement recorded for `living-circles-mrl`
* productivity is interpreted cumulatively from repository start through `2026-04-05`
* `git log --until 2026-04-05` shows sixty-one `feat:` commits and sixty-seven `docs:` commits in that cumulative window
* this was a high-volume shaping period, so the session note is best read as a coarse operational snapshot rather than as a per-slice benchmark
