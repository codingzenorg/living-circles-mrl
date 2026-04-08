# Expectation Gap Detection

## Scope

- reviewed against [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- reviewed against [domain_background_knowledge.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/domain_background_knowledge.md)
- reviewed against the current slice [initial_regional_food_regeneration_pressure.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/slices/initial_regional_food_regeneration_pressure.md)
- reviewed against the current implementation summary in [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- validated with `go test ./...` and `npm run test:contracts` on April 8, 2026

## Summary

The current build continues moving in the right direction. The default world is no longer only a plausible startup snapshot; it now has a first region-sensitive recovery rule. That matters because the previous EGD identified the main gap as medium-term ecological consequence rather than startup plausibility.

This slice improves that gap, but only partially. The system can now let different neighborhoods recover at different rates, which is more ecosystem-like than a purely global regeneration delay. But the world still does not yet show stronger longer-horizon outcomes such as visible regional dominance, durable niches, or lineage-level population consequences.

## Main Findings

### 1. The build now has its first meaningful region-level ecological difference

This is the strongest positive finding in the current EGD.

Before this slice:

- startup state varied more than before
- the world was larger
- food and expanded autonomous startup were less authored
- but recovery pressure still remained too global

Now:

- food still regenerates deterministically into the same slot
- but locally depleted neighborhoods recover more slowly than less depleted ones

That is a real step toward the model hypothesis, because world regions can now diverge over time even when the underlying food model stays simple.

### 2. The simulation is becoming more ecosystem-like, but still mostly through additive heuristic rules

The system now combines:

- energy costs
- food recovery
- global and local regeneration pressure
- crowding pressure
- food seeking
- avoidance heuristics
- shape-based interaction
- child continuity

This is coherent and increasingly rich. But it still reads as a layered bundle of local rules rather than a smaller set of highly generative ecological principles.

That remains an open design question rather than a bug:

- is the intended identity satisfied by cumulative bounded rules?
- or should future slices try to produce more from fewer mechanisms?

### 3. Medium-term consequence improved, but long-horizon dynamics are still not yet evidenced

The regional regeneration slice helps the world support:

- uneven recovery
- local depletion pressure
- more meaningful large-world geography

But EGD still cannot point to stronger long-horizon patterns such as:

- stable population clusters
- repeated collapse/recovery cycles visible at world scale
- strong area-based competition pressure
- clearly differentiated autonomous niches

So the gap has narrowed from “no medium-term ecological differentiation” to “some medium-term differentiation exists, but it is still not enough to prove the larger ecosystem claim.”

### 4. The larger viewport and minimap now have a stronger systemic reason to exist

This is an important secondary consequence.

Earlier, the viewport improvements were mostly presentation gains:

- player-follow camera
- deadzone
- minimap
- offscreen awareness
- heading cue

Now that recovery pressure is region-sensitive, the larger world is not just more eye candy. It has a better systemic justification:

- where the player and autonomous circles move matters more over time
- different parts of the map can meaningfully differ in recovery pace

That strengthens the coherence between the presentation layer and the world model.

### 5. Lineage is still present but still not one of the strongest ecosystem drivers

Nothing in this slice changes the earlier lineage conclusion.

The build still has:

- attached children as the authoritative child model
- continuity through promoted children
- visible lineage and generation
- explicit inspectability of child-driven events

But lineage still contributes more to continuity and local outcomes than to larger population patterns.

That means the distinctiveness of Living Circles still leans more on:

- energy
- shape semantics
- continuity embodiment

than on truly strong lineage-driven ecology.

### 6. Process quality remains strong: the repo is still catching real system issues during validation

The previous loop surfaced and fixed a websocket transport race during validation. This loop added a direct simulation-level test for regional regeneration timing and still passed the broader integration suite.

That is a positive signal:

- the process is not only adding features
- it is preserving authority, determinism, and regression discipline while the world model grows

## Secondary Findings

- The world now has a more credible relationship between large map scale and resource dynamics.
- Deterministic seeding continues to be a strong design choice; it allows richer initial conditions without sacrificing reset/test repeatability.
- The implementation summary still contains some stale older simplification bullets, so repo-memory compression remains a cleanup candidate.

## Main Expectation Gaps

1. The world now supports local recovery divergence, but still lacks clearer long-horizon population dynamics.
2. Autonomy still appears mostly heuristic-layered rather than strongly niche-forming.
3. Lineage remains visible and inspectable, but still not strongly consequential at ecosystem scale.
4. Play quality is better supported by the model now, but still lightly evaluated as experience rather than just structure.

## Recommendation

Return to `refine`, but stay in the medium-term consequence direction rather than returning to startup-world polish.

The strongest next directions are:

1. a regional ecological pressure slice beyond food timing alone
2. a longer-run evaluation slice for population dynamics
3. a lineage-consequence slice only if you want distinctiveness over pure ecosystem pressure

## Best Next Slice Directions

### Option A: Regional crowding or energy pressure

Extend the current local differentiation logic so crowded or overused areas become more energetically expensive beyond food recovery timing alone.

Why this fits:

- it builds on the new regional resource-pressure foundation
- it would deepen area-level divergence over time
- it could create stronger movement and settlement patterns without adding terrain systems

### Option B: Longer-run population-dynamics evaluation

Add a bounded evaluation artifact or scenario that observes the world over longer runs and records whether regional depletion, recovery, dominance, and continuity patterns actually emerge.

Why this fits:

- the system now has enough medium-term pressure to justify a more serious observation pass
- the next uncertainty is increasingly empirical rather than purely structural
- it can reveal whether more behavior slices are needed or whether the current rule set already produces enough variety

### Option C: Lineage consequence slice

Make lineage matter more than continuity bookkeeping, for example through a bounded rule that gives longer-lived or higher-generation lines a meaningful ecological tradeoff or advantage.

Why this fits:

- lineage is already deeply embodied and inspectable
- it still feels weaker than the model hypothesis suggests
- it would increase distinctiveness relative to a more generic resource-competition sandbox

## Return-To-Loop Recommendation

Recommended next phase: `refine`

Recommended intent:

- keep building toward medium-term consequence, not startup variety
- prefer either regional ecological pressure or a longer-run evaluation slice
- only switch to lineage consequence next if distinctiveness is now more important than broader ecosystem dynamics
