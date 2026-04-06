# Expectation Gap Detection

## Scope

- reviewed against [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- reviewed against [domain_background_knowledge.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/domain_background_knowledge.md)
- reviewed against current implementation summary in [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- validated with `go test ./...` and `npm run test:contracts` on April 6, 2026

## Summary

The current build is no longer missing core domain categories. It already expresses:

- energy-driven life and death
- food recovery and regeneration
- same-shape conflict
- different-shape reproduction
- child-based continuity
- autonomous participation under shared rules
- attached children as visible embodied state

The main expectation gaps are no longer about missing nouns. They are about whether the simulation produces convincing ecosystem behavior instead of a small deterministic rules demo.

## Main Findings

### 1. The world is still too small and too curated to validate emergence

The current implementation deliberately keeps:

- one player
- exactly two autonomous circles
- deterministic initial placements
- deterministic food slots
- strongly bounded steering heuristics

This is enough for slice validation, but weak for the model hypothesis of a `system-driven ecosystem`. The built state currently proves the rules exist; it does not yet prove that those rules generate interpretable population dynamics.

### 2. Autonomy is still heuristic-driven rather than convincingly ecosystem-like

Autonomous circles now seek food, pursue feasible reproduction, avoid losing fights, and retreat from blocked reproduction. That is coherent, but it still reads more like layered local heuristics than like a dense emergent population model.

This is not necessarily wrong, but it creates a review question:

- is the intended game satisfied by bounded reactive heuristics, or should autonomy become less hand-shaped and more system-derived?

### 3. Lineage exists, but long-term lineage meaning is still thin

The build now exposes lineage, generation, promoted-child identity, absorbed-child identity, payment-child identity, created-child identity, and ownership identity. That makes continuity inspectable, but not yet deeply meaningful.

The current system still lacks stronger lineage consequences such as:

- persistent family differentiation
- lineage-level dominance or collapse patterns
- inherited variation or divergence

So lineage is present, but still semantically shallow relative to how central it feels in the model hypothesis.

### 4. Reproduction is highly inspectable, but likely over-instrumented relative to gameplay value

Recent slices have made reproduction extremely explicit:

- capacity totals
- capacity components
- threshold and cost constants
- blocked-side identity
- payment-side identity
- payment-child identity
- created-child identity
- ownership identity
- distribution kind

This is excellent for debugging and refinement, but the review risk is that the loop has shifted from improving the ecosystem to annotating it. Further slices in this direction are likely to have diminishing value unless a concrete ambiguity is still blocking design decisions.

### 5. Responsiveness and authority are structurally present, but not yet evaluated as experience qualities

The architecture pressure from the source material is represented correctly:

- JavaScript browser client
- 2D canvas rendering
- WebSocket communication
- authoritative Go server

But the current review evidence does not yet say much about:

- perceived responsiveness
- readability of real-time state changes
- whether the player can understand danger and opportunity quickly
- whether authoritative updates feel coherent under play

That is a meaningful gap because responsiveness is part of the genre baseline, not a cosmetic concern.

### 6. Shape semantics are implemented, but their legibility under play still needs evaluation

Same-shape conflict and different-shape reproduction are now explicit, and the demo visually distinguishes shapes. But EGD should still flag a practical review question:

- does a new player actually understand the shape rule quickly from the running simulation, or only after reading HUD/debug information?

The system may be semantically correct but still weakly legible as a play experience.

## Secondary Findings

- Attached children are now the central embodied child model, which is a strong improvement over earlier count-only semantics.
- The staged removal of hidden derived-radius behavior improved model coherence.
- The default demo remains more of a demonstrator than a convincing living world.

## Recommendation

Return to `refine`, but not for another tiny inspectability slice by default.

The best next loop options are:

1. refine a new behavior cluster aimed at ecosystem validity
2. refine an evaluation/exposure slice that captures real play evidence
3. only continue inspectability if a specific ambiguity is still blocking design decisions

## Best Next Slice Directions

### Option A: Population-scale ecosystem slice

Increase the number of autonomous circles and tune world/resource conditions so the system can exhibit visible dominance, collapse, and recovery patterns rather than only pairwise demonstrations.

### Option B: Behavior-pressure slice

Introduce one new ecological pressure that changes medium-term outcomes, such as stronger scarcity, crowding pressure, or more consequential continuity tradeoffs.

### Option C: Play-legibility slice

Improve what a player can infer from the live simulation without relying on debug-heavy HUD details, especially around shape semantics, danger, and reproduction opportunity.

## Return-To-Loop Recommendation

Recommended next phase: `refine`

Recommended intent:

- stop the current inspectability chain unless a specific debugging need remains
- define the next slice around ecosystem validity or play-legibility
- treat the current build as sufficient for a release-style internal checkpoint once one more meaningful behavior or evaluation step is chosen
