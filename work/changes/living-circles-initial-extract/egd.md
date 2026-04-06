# Expectation Gap Detection

## Scope

- reviewed against [model_hypothesis.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/model_hypothesis.md)
- reviewed against [domain_background_knowledge.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/semantics/domain_background_knowledge.md)
- reviewed against the current slice [initial_seeded_expanded_autonomous_state_mix.md](/home/henrique/repos/github/codingzen/living-circles-mrl/docs/slices/initial_seeded_expanded_autonomous_state_mix.md)
- reviewed against the current implementation summary in [implementation.md](/home/henrique/repos/github/codingzen/living-circles-mrl/work/changes/living-circles-initial-extract/implementation.md)
- validated with `go test ./...` and `npm run test:contracts` on April 6, 2026

## Summary

The current build now clears an older expectation gap: the world no longer starts as a tiny, obviously hand-authored toy setup. The default startup state is materially broader:

- larger world
- more autonomous circles
- seeded food layout
- seeded expanded autonomous placement
- seeded expanded autonomous startup state mix
- viewport-oriented presentation instead of whole-world shrinking

That is real progress toward the model hypothesis of a system-driven ecosystem. But the main gap has shifted. The repo now has enough spatial breadth and startup variation to feel less staged, yet the system still lacks stronger medium-term ecological differentiation. It is better at startup plausibility than at long-horizon emergent consequence.

## Main Findings

### 1. Startup world plausibility is now much stronger than before

Compared with the earlier EGD, the default world is no longer dominated by obviously authored startup geometry.

What now works better:

- the world is large enough to justify viewport mode
- food no longer starts from a visibly center-authored slot list
- expanded autonomous circles no longer start from a visibly authored offset pattern
- additional expanded circles no longer have a purely hand-scripted startup shape/energy mix

This closes a meaningful gap from the earlier review: the startup world now reads more like a seeded ecosystem state than a deterministic rules exhibit.

### 2. Medium-term ecosystem behavior still looks thin relative to the model hypothesis

The model hypothesis emphasizes:

- emergence
- population dynamics
- collapse and recovery
- the player as one perturbing participant inside a shared rule set

The current system now supports more plausible startup variety, but EGD still cannot point to stronger medium-term consequences such as:

- visible population booms and busts
- region-level overuse and recovery patterns
- durable competitive niches
- meaningful long-horizon lineage turnover

So the expectation gap is no longer “the world is too curated to start.” It is now “the world starts plausibly, but still may not evolve richly enough after startup.”

### 3. Autonomy remains coherent but still reads as layered heuristics

Autonomous behavior is now fairly rich for a bounded rules system:

- food seeking
- energy-aware priorities
- interaction seeking
- fight-feasibility filtering
- blocked-reproduction avoidance
- threat avoidance
- crowding-aware steering

This is coherent and increasingly gameable. But it still looks like an accumulation of local steering rules rather than a deeper ecological dynamic.

That is not automatically a problem. It becomes a problem only if the intended identity requires stronger emergent differentiation than these heuristics can produce.

Current review question:

- is this project aiming for “good bounded ecosystem-like behavior,” or for “deeper emergent dynamics from fewer more generative rules”?

### 4. The world is better staged for play, but still only lightly evaluated as play

The viewport, deadzone, minimap, offscreen cues, heading cue, camera lookahead, and fullscreen layout all move the demo closer to an actual play surface.

That is a strong improvement over the older whole-world-shrunk presentation. But the current repo evidence still says more about legibility affordances than about actual play quality.

Open experience-level questions remain:

- does the larger viewport world feel satisfying to move through over time?
- do offscreen awareness and minimap cues feel helpful rather than supervisory?
- does the player actually perceive meaningful local decisions, or mostly observe a system?
- does the larger world improve play, or only presentation?

So the gap here is evaluation, not missing UI.

### 5. Lineage is inspectable and embodied, but still not yet strategically central

The build now has:

- attached children as the authoritative child model
- continuity through promoted children
- lineage and generation
- visible continuity cues
- explicit identity for promoted, absorbed, paid, created, and contacted children

This is conceptually much stronger than earlier builds. But lineage still functions more as a continuity and inspectability layer than as a major strategic force in the ecosystem.

What still seems thin:

- longer-term lineage advantage or disadvantage
- lineage-level divergence
- visible family persistence patterns that matter beyond immediate continuity

So lineage is no longer missing, but it still is not yet one of the strongest drivers of the ecosystem.

### 6. The recent validation exposed and fixed a real authority/transport weakness

During the latest build validation, the integration suite surfaced a real concurrent websocket write bug around reset broadcasting. That was fixed by serializing writes per connection in the transport layer.

This matters in EGD because it confirms something important:

- the repo is now large and real enough that authority and transport behavior can fail in genre-relevant ways
- integration validation is catching issues that pure local reasoning would miss

That is a positive sign about process quality, even though it surfaced as a bug first.

## Secondary Findings

- The default startup world now better supports the current fullscreen viewport presentation.
- Deterministic seeding is doing useful work: the repo is gaining variety without giving up reset/test reproducibility.
- The implementation summary still contains some outdated older bullets, so repo-memory compression could be improved even though the build itself is coherent.

## Main Expectation Gaps

1. The system is now stronger at startup plausibility than at medium-term ecological consequence.
2. Autonomy still appears heuristic-layered rather than deeply generative.
3. Play quality is still more assumed than evaluated.
4. Lineage is visible and inspectable, but not yet strongly consequential over longer horizons.

## Recommendation

Return to `refine`, but avoid another startup-world slice immediately unless it directly creates medium-term ecological consequences.

The most promising next directions are:

1. a medium-term ecological pressure slice
2. a population-dynamics evaluation slice
3. a lineage-consequence slice

## Best Next Slice Directions

### Option A: Regional resource pressure

Introduce a localized resource or depletion rule that makes different areas of the larger world diverge meaningfully over time, not just at startup.

Why this fits:

- the world is now large enough for regions to matter
- it would create medium-term consequence instead of only startup variety
- it would test whether the larger world actually supports population dynamics

### Option B: Population-dynamics evaluation slice

Add a bounded evaluation artifact or scenario that observes the system for longer runs and records whether dominance, collapse, recovery, and continuity patterns actually appear.

Why this fits:

- the startup world is now more credible
- the next gap is increasingly about what the system does over time
- this could reveal whether more behavior slices are needed or whether current rules are already sufficient

### Option C: Lineage consequence slice

Make lineage matter more than continuity bookkeeping, for example by giving longer-lived or higher-generation lines some explicit ecological consequence.

Why this fits:

- lineage is already present and visible
- it is still weaker than the model hypothesis suggests
- this would increase distinctiveness relative to a more generic circle ecosystem

## Return-To-Loop Recommendation

Recommended next phase: `refine`

Recommended intent:

- stop polishing startup variability for now
- choose a slice that changes medium-term consequence, not only initial conditions
- if uncertainty remains high, choose an evaluation slice before introducing another behavior layer
