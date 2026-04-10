# Expectation Gap Detection

## Scope

- reviewed against [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- reviewed against [domain_background_knowledge.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/domain_background_knowledge.md)
- reviewed against the current slice [initial_post_overlay_render_reassessment.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/slices/initial_post_overlay_render_reassessment.md)
- reviewed against the current implementation summary in [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- reviewed as a lightweight artifact-led EGD rather than an Ollama-backed scenario packet
- validated with `go test ./... -timeout 60s` and `npm run test:contracts` on April 10, 2026

## Summary

This optimization track has reached a sensible stopping point.

The recent sequence of client-side render slices succeeded on its own terms:

- render pressure became visible
- the browser path was broken into major families
- world and overlay buckets were further decomposed
- several bounded overlay and label costs were reduced
- the post-reduction read now shows a materially calmer render path

That is a real implementation improvement. The browser-side presentation is more disciplined and less speculative than before.

The important EGD conclusion is that the main remaining gap is no longer “keep trimming overlays.” The current live read now shows:

- `world` as the top-level dominant family at about `1.5ms`
- `overlay` materially below it at about `0.6ms`
- `support` and `minimap` already small
- world-side labels and circle drawing as the largest remaining client costs

So this EGD does not point to another obvious overlay optimization. It points to a change in priority: either stop the render-optimization track for now, or only continue if the goal is a specific world-side presentation change such as label/circle drawing. Relative to the broader model hypothesis, the more meaningful remaining gaps are still ecological and experiential rather than raw browser overlay cost.

## Main Findings

### 1. The render optimization sequence was disciplined and successful

This is the strongest positive finding.

The recent work did not just “make things lighter” in an ad hoc way. It followed an evidence path:

- total render measurement
- family measurement
- world subfamily measurement
- overlay subfamily measurement
- bounded reductions
- reassessment

That is exactly the right shape for client-side optimization work in this repository. The current browser path is now better understood and more intentionally shaped.

### 2. Overlay trimming no longer looks like the highest-value next step

This is the main decision-level finding.

The reassessment result shows that the overlay path is now materially smaller than the world path. That matters because the recent slice history could otherwise have drifted into optimization inertia.

Instead, the current evidence says:

- overlay work is now present but not obviously dominant
- the remaining top client cost is world-side drawing
- further overlay trimming is likely to deliver smaller gains while risking legibility erosion

That is a healthy place to stop.

### 3. The browser path is calmer, but the project still lacks a stronger play-feel evaluation artifact

The current render instrumentation makes browser draw cost visible, which is useful. But this EGD still cannot say much about felt responsiveness beyond inference.

What is still missing:

- a direct play-feel or responsiveness evaluation artifact
- a browser-side frame stability artifact over a longer movement window
- a stronger comparison between measured draw cost and actual perceived roughness

So the render work is technically coherent, but still not the same thing as a user-facing responsiveness evaluation.

### 4. The main semantic gaps remain outside this optimization track

This is not a criticism of the slices; it is a scope boundary.

Relative to the model hypothesis and background knowledge, the current render work improves:

- visual discipline
- measurement quality
- confidence about client cost families
- the likelihood that browser rendering is not the immediate weak point

But it does not materially advance:

- long-horizon ecological consequence
- population dynamics
- lineage significance beyond inspectability
- stronger user-facing evaluation of fairness and emergence

So the largest remaining gaps are still domain and experience questions, not client overlay cost questions.

### 5. If client render work continues, it should now target world-side drawing or stop entirely

The post-reduction read gives a clear boundary.

If another client render slice is chosen, the most defensible next targets are:

- world-side label drawing
- world-side circle drawing density or style

What no longer looks justified by default:

- another overlay-specific trim
- more browser micro-optimizations without fresh evidence of actual roughness

That means the current track has reached a point where a stop-or-switch decision is healthier than automatic continuation.

## Secondary Findings

- The render HUD and tooltip now provide enough live evidence to support local review without introducing a heavier profiling system.
- The support panels successfully carried more identity/detail load while canvas work became lighter.
- The minimap remained intact through the optimization sequence, which helps preserve whole-world orientation while local overlays became cheaper.

## Main Expectation Gaps

1. The repository still lacks a stronger play-feel evaluation artifact that relates measured render cost to perceived responsiveness.
2. The main remaining client cost appears to have moved from overlays to world-side drawing, but that may no longer be the most valuable project problem.
3. The larger remaining gaps are still ecological and experiential rather than purely graphical.
4. Continuing client render optimization by default would now risk diminishing returns.

## Recommendation

Do not continue the current overlay-optimization track by inertia.

Recommended next phase: `refine`

Recommended intent:

- either switch to a broader evaluation or domain-pressure direction
- or, if client performance is still the explicit concern, refine one world-side render slice rather than another overlay slice

## Best Next Slice Directions

### Option A: Player responsiveness or play-feel evaluation

Add a bounded evaluation slice that relates the current transport and render measurements to actual play responsiveness and readability during ordinary movement.

Why this fits:

- the domain background explicitly says responsiveness matters
- transport and render work are currently being used as proxies for play feel
- the project now needs a more direct experience-oriented read

### Option B: Return to ecological or lineage meaning

Switch back toward the model’s stronger differentiators instead of continuing optimization work.

Best candidates:

- longer-run ecosystem evaluation
- regional or population-dynamics consequence
- stronger lineage consequence

Why this fits:

- the main semantic gaps are still there
- the current optimization track no longer shows one glaring technical hotspot
- this is more aligned with Living Circles’ distinct identity than more browser trimming

### Option C: One world-side render slice

Choose this only if client rendering is still the explicit concern.

Best candidates:

- world label/circle drawing reduction or restyling
- a world-side draw simplification that keeps play-surface readability intact

Why this fits:

- the current live read points to world-side work, not overlay work, as the largest remaining client cost
- it is the only clearly evidence-backed continuation of the current render track

## Return-To-Loop Recommendation

Recommended next phase: `refine`

Recommended direction:

- prefer a play-feel or ecological direction next
- only continue client render optimization if there is a specific reason to target world-side drawing
