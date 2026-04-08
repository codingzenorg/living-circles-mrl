# Impact Analysis

## Change

Define the first build slice for Living Circles around the intended client/server runtime shape rather than the starter repository's example Python monolith shape.

## Why This Matters

The extracted evidence establishes a concrete product topology:

- JavaScript browser client
- 2D canvas rendering
- WebSocket communication
- authoritative Go simulation server

The current repository architecture still describes `python_ddd_monolith` as the selected pack. If build proceeds without addressing that mismatch, the repository will drift away from the extracted model before the first executable slice even exists.

## Impacted Areas

### Architecture documents

- `architecture.md` should stop describing `python_ddd_monolith` as the current selected pack for this repository
- `decisions.md` should record adoption of `polyglot_client_server` if that becomes the accepted direction

### Repository layout

- `src/` should be shaped around client, server, and shared contracts rather than the Python monolith example layout
- `tests/` should include client, server, integration, and contract coverage

### Slice design

- early slices should make runtime targets explicit
- early slices should avoid pretending that browser behavior can be collapsed into a fake single-process abstraction

### Evaluation expectations

- responsiveness, authority, and contract stability become first-class review dimensions
- deterministic testing still matters, but it must be applied across a runtime boundary

## Recommended Decision Pressure

The next implementation-facing update should explicitly decide whether this repository is now adopting `polyglot_client_server` as its selected pack.

If yes:

- update `architecture.md`
- add a decision entry in `decisions.md`
- let build scaffold the repository accordingly

If no:

- document why the repository is intentionally using a different implementation shape than the extracted product topology
- explain how that mismatch will still support valid refinement

## Risks If Ignored

- the first code may be scaffolded in the wrong layout
- client/server authority may be blurred or hidden
- later slices may need avoidable structural rework
- semantic artifacts and implementation structure will diverge early

---

## Change

Reinterpret children as attached orbiting dependents of a parent circle instead of immediately free autonomous circles.

## Why This Matters

The current implementation made reproduction visible by spawning a free active child circle. That was a useful executable step, but it does not match the stated product idea that children orbit their parents.

This is not a cosmetic tweak. It changes:

- what a child entity is
- how reproduction output is represented
- how snapshots should expose child state
- how future fight, food, and continuity rules should relate to visible children

## Impacted Areas

### Simulation model

- reproduction output should become parent-owned attached children rather than immediate free autonomous participants
- world update logic needs an orbit model derived from parent-child state
- the current child-count model becomes transitional rather than the whole embodiment

### Runtime contract

- snapshots need explicit attached-child state
- the current representation of a child as just another autonomous circle is no longer sufficient
- child ownership and orbit inspectability become first-class contract concerns

### Browser rendering

- the client should render children orbiting a parent rather than appearing as free peers in the world
- labels and debugging affordances should clarify which parent owns which child

### Existing semantics

- current radius growth, fight leverage, replacement continuity, and reproduction payment all still depend on `children_count`
- the repository needs an explicit transitional stance on whether visible orbiting children and count-based shortcuts coexist temporarily

### Determinism discipline

- the intended game feel says post-reproduction children are randomly distributed between parents
- deterministic testing means build should use a reproducible authoritative distribution rule rather than unrestricted randomness

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether positive interaction targeting should measure nearness from parent cores only or from current parent-or-child geometry
- how deterministic tie-breaking behaves when multiple eligible targets are equally near through different embodied points
- whether current avoidance rules should stay unchanged while only positive targeting becomes more embodied

## Risks If Ignored

- attached children will remain more meaningful in authoritative contact than in pursuit steering
- the simulation will keep mixing embodied collision geometry with more abstract target selection
- later autonomy slices will keep carrying a gap between visible child placement and positive pursuit behavior

---

## Change

Remove grown derived radius from parent movement boundary clamping.

## Why This Matters

Recent slices already made food collection, circle contact, and same-shape fight resolution less dependent on hidden radius growth. But world-edge movement still clamps parent cores using grown derived radius, which means circles with more children are still kept farther away from walls even though those same children no longer silently extend feeding or contact reach through parent-body size.

The next model pressure is to make movement boundaries follow the same visible parent-body geometry already used in the more embodied feeding and contact slices.

## Impacted Areas

### Simulation model

- parent movement clamping should use the visible parent body rather than grown derived radius
- player and autonomous movement should continue sharing one deterministic boundary rule
- attached-child orbit layout can remain unchanged in this slice

### Runtime contract

- the current snapshot shape is likely sufficient because parent positions and radius are already visible
- build may need no contract change if the new clamp is observable from ordinary snapshots

### Browser rendering

- current rendering may already be enough because boundary position is visible on the canvas
- no visual-size change should be required in this slice

### Existing semantics

- food, contact, fight, reproduction, continuity, and steering should remain unchanged
- derived radius may remain for rendering and other transitional semantics outside boundary clamping

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether boundary clamping should use the fixed visible parent-body size already implied by embodied food and contact
- whether attached children stay free to visually extend beyond the parent-body clamp for now
- how tests should prove that child-derived radius no longer changes wall proximity

## Risks If Ignored

- movement will keep using a hidden growth shortcut after feeding, contact, and fight have already moved away from it
- circles with more children will still feel larger in world navigation even when other embodied slices say parent-body geometry should dominate
- later efforts to remove or narrow radius semantics will have to keep working around inconsistent wall behavior

---

## Change

Add a small directional camera lookahead to viewport mode.

## Why This Matters

Viewport mode now solves several orientation problems:

- the player no longer sees the whole world shrunk down
- the camera follows the player with a deadzone
- the minimap restores whole-world orientation
- local heading and offscreen cues make nearby pressure easier to read

But camera feel is still slightly neutral. Once the player is already moving in a clear direction, the viewport still gives equal space in every direction. That keeps the stage readable, but it leaves some of the new fullscreen eye-candy potential unused.

The next pressure is to make the viewport feel a bit more alive and responsive without crossing into cinematic camera behavior.

## Impacted Areas

### Browser rendering

- the client camera transform should allow a small directional lookahead
- that offset should be derived from recent authoritative player motion
- the stage should still feel anchored and controllable

### Camera behavior

- the current deadzone model should remain intact
- world-edge clamping should remain intact
- lookahead should be small enough to avoid jitter or floatiness

### Existing semantics

- server/world behavior should remain unchanged
- minimap, heading cue, and offscreen awareness should keep working under the shifted viewport
- support panels and other UI should remain secondary to the stage

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether lookahead is derived from the player's recent authoritative heading or only from active input direction
- how large the offset can become before it starts harming control clarity
- how lookahead recenters when the player slows or changes direction

## Risks If Ignored

- viewport mode will remain more readable than the old whole-world view but still flatter than it could be
- the player will gain awareness cues without gaining a stronger sense of forward motion in the stage
- later polish work may need to revisit the camera as a separate concern instead of building on the existing viewport track

---

## Change

Increase the default world scale and startup population while replacing hand-authored expanded food placement with deterministic seeded layout generation.

## Why This Matters

The current expanded world baseline improved the old tiny demo, but it still starts from a visibly curated setup:

- `1200x900` space
- `5` autonomous circles
- food drawn from a fixed ordered slot list around the center

That remains test-friendly, but it limits the sense of scale now that the client uses a player-following viewport. The world can appear larger on screen, yet startup conditions still feel arranged and sparse.

The next pressure is to increase ecological breadth without giving up determinism:

- larger default world
- more default autonomous circles
- more random-looking food layout at startup
- identical reset behavior for the same seed/rule

## Impacted Areas

### Simulation model

- expanded default dimensions may increase again
- expanded default autonomous baseline may increase
- initial food slot generation should move from a hand-authored list to deterministic seeded generation

### Reset behavior

- reset must still rebuild the exact same authoritative startup world
- deterministic tests must remain able to assert startup shape and content

### Browser rendering

- viewport mode and minimap should naturally benefit from the larger space and denser startup world
- no contract change is required if coordinates remain the same kind of state

### Existing semantics

- food regeneration, crowding pressure, autonomy, fights, reproduction, and continuity should remain unchanged
- narrow custom worlds should remain small and explicit for focused tests

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the new expanded default world dimensions
- the new expanded autonomous baseline
- the deterministic seeded rule for initial food slot generation
- how generated food avoids unreadable overlap with edges or initial entities

## Risks If Ignored

- the viewport presentation will continue to outgrow the ecological scale of the startup world
- default startup conditions will keep feeling hand-authored rather than world-like
- future evaluation of population-scale behavior will stay constrained by sparse, arranged initial conditions

---

## Change

Replace the expanded default autonomous startup pattern with deterministic seeded placement.

## Why This Matters

The last slice removed the most obvious authored feeling from food placement, but autonomous startup is still arranged by hand:

- a fixed set of offsets around the center
- a clearly human-authored cluster pattern
- startup space that still reads as staged rather than ecological

That means the world is currently mixed:

- seeded food
- authored autonomous placement
- explicit player spawn

The next pressure is to make the expanded startup state more uniformly world-like without losing deterministic reset behavior.

## Impacted Areas

### Simulation model

- expanded autonomous startup placement should move from fixed offsets to deterministic seeded generation
- player startup can remain explicit and stable
- expanded shape, energy, and child-count semantics can stay as they are

### Reset behavior

- reset must still rebuild the exact same autonomous startup arrangement
- deterministic tests must remain able to assert startup and reset equality

### Browser rendering

- viewport and minimap should naturally benefit from a less staged startup arrangement
- no contract change is required if the same circle coordinates are still exposed

### Existing semantics

- food initialization, regeneration, crowding, autonomy, fights, reproduction, and continuity should remain unchanged after startup
- narrow custom worlds should continue to support explicit hand-set placements for focused tests

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the deterministic seeded rule for expanded autonomous startup placement
- the minimum safe spacing from the player, world edges, and other initial circles
- whether the current deterministic shape ordering remains fixed while only placement changes

## Risks If Ignored

- the larger seeded-food world will still feel visibly staged because autonomous startup remains authored
- startup spatial variety will lag behind the newer viewport presentation
- later ecosystem evaluation will continue to begin from a partially artificial arrangement

---

## Change

Replace the authored startup state mix of the additional expanded autonomous circles with deterministic seeded startup state.

## Why This Matters

The startup world is now less authored in geometry:

- expanded food uses deterministic seeded layout
- expanded autonomous placement uses deterministic seeded positions

But the extra expanded autonomous circles still carry authored per-ID startup state:

- fixed shape assignment
- fixed starting energy
- fixed "who looks stronger or weaker" at tick zero

That means the world is spatially less staged but semantically still arranged. The next pressure is to reduce that remaining authored startup bias without giving up determinism.

## Impacted Areas

### Simulation model

- additional expanded autonomous circles should derive startup shape and energy from deterministic seeded rules
- player state and the first explicitly-configurable circles can remain fixed anchors
- runtime rules should remain unchanged after initialization

### Reset behavior

- reset must still rebuild the exact same expanded startup state mix
- deterministic tests must remain able to assert startup/reset equality

### Browser rendering

- the viewport and minimap should naturally benefit from a startup world that feels less pre-scripted
- no contract change is required if the same snapshot fields remain in use

### Existing semantics

- fights, reproduction, continuity, food, regeneration, and autonomy remain unchanged after startup
- narrow custom worlds should continue to support explicit targeted startup states

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the deterministic seeded rule for extra-circle startup shape selection
- the deterministic seeded rule for startup energy values
- the acceptable state range so startup variety does not silently redefine runtime behavior

## Risks If Ignored

- the startup world will still feel partially authored even after seeded placement work
- early interaction pressure will remain more scripted than the larger world presentation suggests
- later evaluation of ecosystem startup diversity will stay constrained by a fixed authored extra-circle mix

---

## Change

Make food regeneration pressure regional rather than only global.

## Why This Matters

The current world now starts from a much less authored state:

- larger map
- more autonomous circles
- seeded food layout
- seeded expanded autonomous placement
- seeded expanded autonomous startup state mix

That improves startup plausibility, but EGD now says the bigger gap is medium-term ecological consequence. The current food-recovery rule still applies pressure mostly at the whole-world level: missing-slot count changes the delay, but one depleted local area is not strongly distinct from another if the same total number of slots is missing.

The next pressure is to make different parts of the larger world recover differently over time.

## Impacted Areas

### Simulation model

- food regeneration timing should depend partly on local missing-slot pressure
- slot identity and return-to-origin behavior should remain unchanged
- determinism should remain intact

### Population dynamics

- heavily stripped neighborhoods should recover more slowly
- lightly depleted neighborhoods should rebound sooner
- larger world scale should begin to produce region-level divergence

### Browser rendering

- no contract change is required if the rule remains observable through ordinary food presence/absence over time
- current viewport and minimap should naturally expose the resulting divergence without new UI in this slice

### Existing semantics

- feeding, autonomy, fight, reproduction, continuity, and startup seeding should remain unchanged
- narrow custom worlds should remain useful for deterministic focused tests

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the local radius used to measure regional food pressure
- how local missing slots add to regeneration delay
- how to keep the rule understandable and deterministic in tests

## Risks If Ignored

- the larger world may still feel more varied at startup than during ongoing simulation
- region-level collapse and recovery will remain weak
- later ecology slices may need to add stronger systems to compensate for missing local resource consequence

---

## Change

Add regional crowding-driven energy pressure beyond the current local crowding surcharge.

## Why This Matters

The current build now has one real regional ecological difference:

- local food recovery can diverge by neighborhood

But inhabitation cost is still mostly immediate and local:

- movement has a base energy cost
- a local crowding surcharge applies when nearby neighbors are present
- yet dense regions do not become substantially more expensive over time than sparse ones

The next pressure is to make the larger world matter not only through resource return timing, but also through ongoing occupancy cost.

## Impacted Areas

### Simulation model

- energy cost should now consider broader regional density, not only immediate local crowding
- the current local crowding surcharge should remain intact
- the rule should remain deterministic and bounded

### Population dynamics

- dense areas should become more energetically expensive to remain in
- sparse areas should become comparatively safer to occupy
- the larger world should support stronger displacement or settlement patterns

### Browser rendering

- no contract change is required if the new pressure remains visible through ordinary energy changes and movement outcomes
- current viewport, minimap, and support panels should naturally expose its effects without new UI in this slice

### Existing semantics

- food, autonomy, fight, reproduction, continuity, and startup seeding remain unchanged
- narrow custom worlds should remain useful for deterministic focused tests

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the broader radius used to define regional density pressure
- the extra energy penalty beyond the existing local surcharge
- how to keep the rule small and explainable rather than turning it into a generalized heatmap system

## Risks If Ignored

- larger-world regions may still differ more in food timing than in actual inhabitation consequence
- local crowding will remain tactically meaningful without becoming strategically meaningful
- the system may continue to feel more like layered local heuristics than like a world with medium-term ecological structure

---

## Change

Remove grown derived radius from attached-child orbit distance.

## Why This Matters

The visible child model is now central to feeding, contact, avoidance, continuity, and demo readability. But attached-child layout still uses grown parent radius to decide how far children sit from the parent core. That means one hidden growth abstraction still shapes visible child geometry even after other embodied slices removed derived radius from feeding, contact, fights, and movement boundaries.

The next model pressure is to make orbiting children visually grounded in the visible parent body rather than in a broader transitional growth radius.

## Impacted Areas

### Simulation model

- attached-child layout should use the visible parent body rather than grown derived radius
- player and autonomous parents should continue sharing one deterministic orbit-distance rule
- current child ownership, slot assignment, and orbit motion should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because parent radius and attached-child positions are already visible
- build may need no contract change if the tighter orbit distance is observable from ordinary snapshots

### Browser rendering

- current rendering may already be enough because attached-child positions are directly visible
- no visual-style change should be required in this slice

### Existing semantics

- food, contact, fight, reproduction, continuity, movement, and steering should remain unchanged
- derived radius may remain for parent rendering and other transitional semantics outside child orbit distance

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether orbit distance should use the same fixed visible parent-body size already implied by embodied food, contact, and movement
- whether the current orbit gap remains unchanged once derived radius is removed from layout
- how tests should prove that child-derived parent radius no longer changes orbit distance

## Risks If Ignored

- visible child geometry will keep depending on a hidden growth shortcut even after most other embodied slices have removed it
- children will remain partly more symbolic than physical in the simulation display
- later efforts to narrow or remove transitional radius semantics will still have to work around orbit layout depending on them

---

## Change

Make local crowding pressure more legible during ordinary play.

## Why This Matters

Recent slices made the world larger, increased population, introduced crowding-based energy cost, and added bounded autonomy that reacts to crowding. But that pressure still remains hard to perceive directly from the running world. A player can often see that several circles are nearby without being able to quickly tell whether the local neighborhood is one of the costly dense zones implied by the authoritative rule.

The next pressure is therefore not another server mechanic. It is to make existing crowding pressure more readable on the client without inventing a second simulation.

## Impacted Areas

### Browser rendering

- the client should make dense local neighborhoods easier to read directly from the canvas
- cues should stay local and restrained rather than becoming a full-map overlay
- the current shape-risk cues and encounter log should remain compatible with any new crowding-pressure cues

### Runtime contract

- the current snapshot shape is likely sufficient because circle positions are already available
- build should avoid adding crowding-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- crowding energy cost remains authoritative on the server
- fight, reproduction, food, continuity, and autonomy rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether crowding readability is best shown through nearby spatial halos, emphasized cluster rings, or another small local cue
- how to keep the cue grounded in current visible positions rather than hidden future state
- how to avoid visual noise now that shape-risk cues are already present

## Risks If Ignored

- the newest ecosystem pressure will remain mostly discoverable only through documentation or repeated trial and error
- the player may understand shape-based risk while still missing why some spaces are energetically worse than others
- future evaluation of emergence and playability will stay partly blocked by weak ecosystem-legibility

---

## Change

Make nearby food opportunity and local food scarcity more legible during ordinary play.

## Why This Matters

The simulation now has meaningful resource pressure: food is finite, startup abundance scales with population, and regeneration slows under deeper depletion. But the player still mostly reads this system indirectly by watching energy drop and by gradually noticing whether food happens to be nearby. That makes one of the main ecosystem loops less immediately legible than shape danger or crowding pressure.

The next pressure is therefore to make current food opportunity and food scarcity easier to perceive on the client without inventing hidden timer logic or predictive overlays.

## Impacted Areas

### Browser rendering

- the client should make nearby food-rich space easier to recognize as recovery opportunity
- the client should make locally sparse space easier to recognize as resource pressure
- any new food-pressure cue should remain compatible with the current shape-risk and crowding cues

### Runtime contract

- the current snapshot shape is likely sufficient because visible food and circle positions are already available
- build should avoid adding food-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- food placement, consumption, and regeneration remain authoritative on the server
- crowding, fight, reproduction, continuity, and autonomy rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether food pressure is best shown through nearby area glow, richer food emphasis, player-centered scarcity framing, or another small local cue
- how to keep the cue grounded in current visible food state instead of hidden regeneration timing
- how to avoid visual competition with the new crowding-pressure cues

## Risks If Ignored

- one of the main ecosystem loops will remain less readable than the newer interaction and crowding cues
- energy management will continue to feel more reactive than readable during live play
- future evaluation of ecosystem quality will remain partly blocked by weak resource-legibility

---

## Change

Make nearby autonomous intent more legible during ordinary play.

## Why This Matters

The recent legibility slices made shape danger, crowding pressure, and food pressure easier to read from the world. But nearby autonomous circles still often look behaviorally opaque. A circle may be moving toward food, closing on another circle, or retreating from local danger, yet the player still has to infer too much of that only by watching several moments in sequence.

The next pressure is therefore to make current autonomous behavior easier to interpret on the client without publishing hidden server intent or adding new simulation rules.

## Impacted Areas

### Browser rendering

- the client should make nearby autonomous motion modes more interpretable
- cues should stay local, lightweight, and compatible with the current shape, crowding, and food cues
- the player should be able to distinguish at least food pursuit, social pursuit, and retreat more quickly

### Runtime contract

- the current snapshot shape is likely sufficient because visible world state and positions are already available
- build should avoid adding autonomy-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- autonomy rules remain authoritative on the server
- fight, reproduction, food, crowding, continuity, and movement rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether autonomy readability is best shown through small directional markers, target-side emphasis, retreat cues, or another lightweight local treatment
- how to keep the cue grounded in visible movement and current geometry rather than unverifiable hidden intent
- how to avoid overloading the existing play-legibility layers

## Risks If Ignored

- the world may become more readable than the actors inside it
- evaluation of emergence and playability will stay partly blocked by weak motion-meaning legibility
- future autonomy changes will remain harder to assess quickly during manual play

---

## Change

Make the player's own authoritative movement more legible during ordinary play.

## Why This Matters

The recent play-legibility slices have made the world, nearby pressures, and nearby autonomous circles easier to interpret. But the player’s own immediate motion still lacks strong visual reinforcement. That leaves one of the most basic experiential questions underdeveloped: does authoritative movement feel readable and responsive enough during play?

The next pressure is therefore to improve motion legibility on the client without adding prediction, interpolation, or new movement rules.

## Impacted Areas

### Browser rendering

- the client should make current player movement direction easier to perceive
- cues should remain tightly local to the player and visually compatible with the existing layers
- the cue should improve readability without pretending the client controls authoritative future motion

### Runtime contract

- the current snapshot shape is likely sufficient because authoritative player positions already arrive each tick
- build should avoid adding movement-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- movement remains authoritative on the server
- fight, reproduction, food, crowding, autonomy, continuity, and shape rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether movement legibility is best shown through a small directional indicator, motion trail, short-range echo, or another tightly bounded cue
- how to keep the cue grounded in recent authoritative positions rather than prediction
- how to avoid visual competition with the current world-state legibility layers

## Risks If Ignored

- responsiveness will remain less evaluated than other aspects of the simulation
- the player may understand the world better than their own moment-to-moment presence inside it
- future playability review will remain partly blocked by weak motion readability

---

## Change

Make lineage continuity more legible during ordinary play.

## Why This Matters

The recent play-legibility slices have made the world, nearby actors, and the player’s own movement easier to read. But one of the project’s differentiating ideas still remains comparatively weak in live play: continuity. Attached children, lineage IDs, generations, and promoted-child survival all exist, yet they still read more like explicit data than like a visible continuity structure in the running world.

The next pressure is therefore to make lineage persistence and continuity-bearing children easier to perceive on the client without adding new lineage mechanics or history tooling.

## Impacted Areas

### Browser rendering

- the client should make parent-child continuity relationships more visible
- continuity persistence should become easier to notice when it happens
- cues should remain compatible with the current danger, crowding, food, autonomy, and motion layers

### Runtime contract

- the current snapshot shape is likely sufficient because lineage IDs, generations, attached children, and recent continuity outcomes are already available
- build should avoid adding lineage-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- lineage and continuity remain authoritative on the server
- fight, reproduction, food, crowding, autonomy, and movement rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether continuity is best shown through stronger parent-child connection cues, lineage-group treatment, recent-promotion emphasis, or another lightweight local approach
- how to keep the cue grounded in current authoritative state rather than broad history reconstruction
- how to avoid visual overload now that several other play-legibility layers already exist

## Risks If Ignored

- one of Living Circles’ most distinctive semantics will remain under-read during ordinary play
- lineage will continue to feel more inspectable than experientially meaningful
- future evaluation of the game’s identity will remain partly blocked by weak continuity legibility

---

## Change

Make recent important world outcomes linger briefly in the scene.

## Why This Matters

The recent play-legibility slices made the present state of the world easier to read: danger, crowding, food, nearby autonomy, motion, and lineage continuity are all more visible. But important interaction outcomes still vanish from the world too quickly. A player may notice a log message, yet still miss where the fight, reproduction, or promotion actually occurred.

The next pressure is therefore to keep recent authoritative outcomes briefly visible in-world without turning the client into an event history system.

## Impacted Areas

### Browser rendering

- the client should render brief world-tied aftermath cues for recent important outcomes
- the cues should remain visually restrained and compatible with the existing layers
- the scene should become easier to interpret during and just after fast interactions

### Runtime contract

- the current snapshot shape is likely sufficient because interaction kinds and entity positions are already available
- build should avoid adding event-history fields unless one minimal readability field is clearly justified

### Existing semantics

- interaction outcomes remain authoritative on the server
- fight, reproduction, continuity, food, crowding, autonomy, and movement rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether recent outcomes are best shown through short-lived rings, pulses, local glow, or another restrained cue
- how long the client-local memory should persist without becoming a replay feature
- how to avoid clutter when several outcomes happen close together

## Risks If Ignored

- fast interaction outcomes will remain harder to connect to places in the world
- the encounter log will stay more informative than the scene itself in busy moments
- future playability review will remain partly blocked by weak recent-event readability

---

## Change

Make nearby actionable world space more visually primary than distant background activity.

## Why This Matters

The recent play-legibility slices have made many individual semantics easier to read, and the latest UI cleanup moved a good amount of textual detail out of the canvas. But the expanded world still carries a lot of simultaneous visible activity. Without a stronger sense of nearby focus, the player can still spend unnecessary effort separating immediately relevant space from distant background motion.

The next pressure is therefore to improve player-centered spatial focus without hiding world state or introducing camera/gameplay changes.

## Impacted Areas

### Browser rendering

- the client should subtly emphasize nearby space around the player
- distant entities and food should remain visible but feel less dominant
- the treatment should support the current danger, food, crowding, motion, lineage, and afterglow cues rather than competing with them

### Runtime contract

- the current snapshot shape is likely sufficient because the client already has the player position and the world geometry
- build should avoid adding focus-specific contract fields unless one minimal readability field is clearly justified

### Existing semantics

- all world rules remain authoritative on the server
- fight, reproduction, food, crowding, autonomy, continuity, and movement rules should remain unchanged in this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether nearby focus is best shown through a player-centered falloff, nearby contrast lift, restrained background dimming, or another subtle treatment
- how to preserve full-world visibility while still improving local scanability
- how to prevent the focus treatment from weakening the existing cue system

## Risks If Ignored

- the expanded world will remain more informative but not proportionally easier to read
- adding more local cues later may increase clutter instead of clarity
- future playability evaluation will remain partly blocked by weak spatial focus

---

## Change

Reduce legend density now that multiple supporting panels and cue systems are already present.

## Why This Matters

The interface has steadily become more capable: the canvas carries richer cues, the player card externalizes player state, the NPC panel externalizes non-player state, and the encounter log externalizes recent outcomes. The legend still reflects many incremental additions, so it is starting to act more like a glossary than a quick play aid.

The next pressure is therefore to reduce explanatory density and let the stronger layout carry more of the understanding.

## Impacted Areas

### Browser rendering

- the legend row should become shorter and easier to scan
- the simplified layout should still preserve the main cue families
- the canvas and side panels should remain the primary source of meaning

### Runtime contract

- no contract change should be necessary

### Existing semantics

- all world rules remain authoritative and unchanged
- this is a presentation simplification only

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- which cue families are essential enough to remain in the legend
- which cues are already sufficiently explained by the world view or supporting panels
- how to reduce density without making the interface cryptic

## Risks If Ignored

- the UI may continue to gain meaning while losing scanability
- supporting explanation may start competing with the actual play surface
- future legibility work may add clutter faster than clarity

---

## Change

Reduce density in the support area below the canvas.

## Why This Matters

Recent UI work successfully moved meaning out of the canvas: player state, NPC summaries, and recent encounters now live in dedicated external panels. That improved the play surface, but it also concentrated interface density into the support stack beneath it. If that area keeps growing without hierarchy improvements, the UI will simply move clutter rather than remove it.

The next pressure is therefore to simplify and prioritize the support area so the canvas remains the clear primary surface and the external panels feel genuinely supporting.

## Impacted Areas

### Browser rendering

- the player, NPC, and encounter panels should become easier to scan as a group
- the support area should communicate clearer priority and lower visual load
- the canvas should remain visually primary

### Runtime contract

- no contract change should be necessary

### Existing semantics

- all world rules remain authoritative and unchanged
- this is a presentation simplification only

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether support density is best reduced through tighter grouping, stronger hierarchy, or partial layout compression
- which support block should remain most prominent
- how to reduce scan cost without hiding useful state

## Risks If Ignored

- the UI may simply relocate complexity rather than reducing it
- the support area may start competing with the canvas for attention
- future simplification work may become harder because more data keeps accumulating in the same stacked region

---

## Change

Reduce text density inside the support area.

## Why This Matters

The support layout is now structurally cleaner, but it still carries a fair amount of status text. That means the interface has improved where information is placed, but not yet fully improved how quickly that information can be scanned. If the support area stays text-heavy, it can still act like a status dump even with better grouping.

The next pressure is therefore to reduce support-text density without losing the small set of state that actually helps during ordinary play.

## Impacted Areas

### Browser rendering

- the player and NPC summaries should become more compact
- supporting text should feel lighter and easier to scan
- the canvas should remain the dominant source of attention

### Runtime contract

- no contract change should be necessary

### Existing semantics

- all world rules remain authoritative and unchanged
- this is a presentation simplification only

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- which support text can be shortened without losing meaning
- how to keep player state easiest to parse
- how to keep NPC state useful while lighter-weight than the current line format

## Risks If Ignored

- the support area may remain more verbose than necessary even after structural cleanup
- simplification work may stall at the layout level and leave textual clutter untouched
- future UI refinement may keep moving density around instead of removing it

---

## Change

Increase the default world size, autonomous population, and food capacity so the simulation can operate as a small ecosystem baseline rather than only a tightly curated mechanics demo.

## Why This Matters

The latest EGD result identified the strongest remaining expectation gap: the rules now exist, but the world is still too small and too curated to validate emergence. The current baseline of one player, two autonomous circles, and a tiny deterministic food set makes the game readable as a rule demonstrator, but weak as a `system-driven ecosystem`.

The next pressure is therefore not another local mechanic. It is to let the existing mechanics coexist at a scale where:

- more than one interaction path can be active at once
- recovery, depletion, pursuit, and avoidance can happen in parallel
- the player enters a living world rather than a staged pairwise example

## Impacted Areas

### Simulation model

- world initialization should produce a larger bounded map
- initial autonomous-circle count should increase beyond the current two-circle default
- initial food capacity should scale alongside the expanded population
- reset behavior should recreate the same expanded world deterministically

### Runtime contract

- the existing snapshot shape may already be sufficient because world size, autonomous circles, and food arrays are already explicit
- build should avoid adding summary metadata unless inspection becomes meaningfully harder without it

### Browser rendering

- the client must remain readable with a larger world and more simultaneous circles
- this slice should avoid introducing camera or zoom systems unless the larger baseline is unusable without them

### Evaluation expectations

- this slice should make ecosystem-level review more meaningful than the current pairwise demo
- EGD after this slice should be able to assess encounter density, depletion/recovery interplay, and whether the world feels less staged

### Existing semantics

- movement, energy, fight, reproduction, continuity, child ownership, and steering rules should remain unchanged
- this slice should scale the baseline, not redesign the rule set

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- a modest deterministic world-size increase rather than an unbounded scale jump
- a modest deterministic autonomous-population increase that keeps tests and live review practical
- a food-capacity increase that supports the larger baseline without introducing a separate balancing subsystem

## Risks If Ignored

- the simulation may keep drifting toward metadata completeness without proving ecosystem validity
- EGD will keep reporting the world as a small curated demo rather than a living system
- future semantic discussion about emergence will stay underconstrained because the world baseline is too thin to generate meaningful patterns

---

## Change

Make food regeneration respond deterministically to the larger world and population baseline instead of using one static global delay for every situation.

## Why This Matters

The expanded-world slice made the system feel less staged, but it also made the current food-recovery rule more visibly provisional. Right now every consumed food slot returns after one fixed delay, regardless of how many circles are alive or how depleted the world currently is.

That static timer was acceptable in the tiny demo phase. In the larger world, it weakens ecosystem validity because resource recovery still behaves like a lab constant rather than part of the living system.

## Impacted Areas

### Simulation model

- food-slot regeneration timing should become sensitive to one small deterministic pressure rule
- consumed slots should still return to the same identity and position
- depletion and recovery should become more meaningful at world level without redesigning the whole economy

### Runtime contract

- the current snapshot shape is likely still sufficient
- build should avoid adding new fields unless the pressure rule becomes too opaque without a minimal inspectability aid

### Browser rendering

- the client may not need any rendering changes if the timing shift is visible from ordinary food disappearance and return
- this should remain a simulation slice, not a HUD expansion by default

### Evaluation expectations

- EGD after this slice should be able to assess whether the larger world shows more believable depletion and recovery cycles
- this is a direct follow-up to the current ecosystem-validity pressure

### Existing semantics

- movement, energy, fight, reproduction, continuity, child ownership, and steering should remain unchanged
- food slot identity and fixed position should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- one deterministic pressure signal for regeneration timing
- a bounded rule that is simple enough to explain and test
- whether the rule should respond to missing-slot count, active-population pressure, or another single world-level signal

## Risks If Ignored

- the world may be larger, but food recovery will still feel like a static demo constant
- depletion and recovery cycles will remain weaker than the model hypothesis suggests
- future ecosystem evaluation will keep conflating scale-up with actual ecological pressure

---

## Change

Make the initial food-capacity baseline derive from starting population scale instead of remaining a fixed authored count.

## Why This Matters

The world is now larger and food regeneration already responds to depletion pressure, but the initial abundance is still effectively a remembered hand-tuned number. That weakens the ecosystem baseline because world startup still depends on curated slot count rather than on an explicit relation between population scale and resource abundance.

The next pressure is to make initial world composition more model-shaped:

- larger starting populations should imply richer initial food support
- smaller focused worlds should remain narrow and deterministic
- reset should recreate a rule-based baseline instead of a magic count

## Impacted Areas

### Simulation model

- initial food slot construction should derive from initial population scale through one simple deterministic rule
- fixed food identity and position should remain intact
- runtime regeneration behavior should remain unchanged in this slice

### Runtime contract

- the current snapshot shape is likely still sufficient because food arrays and world size are already explicit
- build should avoid adding new summary metadata unless the derivation becomes too opaque without it

### Browser rendering

- the client may not need any change because the effect is visible in the starting snapshot itself
- this should remain a simulation-baseline slice, not a UI slice

### Evaluation expectations

- EGD after this slice should be able to assess whether startup abundance now feels more coherent with the larger world and population scale
- the larger world should read less like a staged demonstrator and more like a repeatable living baseline

### Existing semantics

- regeneration timing, movement, energy, fight, reproduction, continuity, child ownership, and steering should remain unchanged
- custom narrow worlds should remain viable for focused deterministic tests

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- one deterministic derivation rule from starting population to starting food capacity
- whether the rule is based on active-circle count alone or another single initial-world signal
- how custom narrow worlds stay intentionally smaller without forking the whole initialization path

## Risks If Ignored

- initial abundance will remain partly arbitrary even as the rest of the ecosystem becomes more rule-shaped
- larger-world behavior will still start from a curated baseline rather than a model-expressing one
- future balancing work will keep carrying unnecessary startup constants

---

## Change

Introduce a local crowding-based energy pressure so denser areas impose an immediate survival cost before explicit collision outcomes occur.

## Why This Matters

The world is now larger, population is higher, startup food support is more rule-shaped, and regeneration pressure already reacts to depletion. But density itself still has little direct cost until food, fight, reproduction, or death events happen.

That keeps the expanded world at risk of feeling like “more actors on a larger board” rather than a genuinely pressured ecosystem. If local clustering has no immediate energetic consequence, then increased population can remain mechanically thinner than the model hypothesis suggests.

## Impacted Areas

### Simulation model

- tick advancement should apply one bounded additional energy pressure based on local nearby-circle density
- the rule should be local rather than global
- player and autonomous circles should remain under the same pressure rule

### Runtime contract

- the current snapshot shape is likely sufficient because energy already appears in ordinary snapshots
- build should avoid new protocol fields unless the chosen rule becomes too opaque without a small inspectability aid

### Browser rendering

- the client may not need any rendering change because the effect should be visible through ordinary energy shifts and resulting outcomes
- this should remain a simulation-behavior slice, not a visualization slice

### Evaluation expectations

- EGD after this slice should be able to assess whether dense regions now create stronger dispersal, collapse, or survival pressure
- this directly strengthens the ecosystem-validity path opened by recent scale and food-pressure slices

### Existing semantics

- movement cost, food placement, food regeneration, fight, reproduction, continuity, child ownership, and steering should remain unchanged
- the new rule should add local pressure, not redesign the rest of the model

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- one neighborhood radius for crowding
- one simple deterministic threshold or formula for additional energy loss
- a pressure magnitude small enough to avoid overwhelming the existing energy loop

## Risks If Ignored

- increased population may continue to read mostly as higher token count rather than as a denser ecosystem
- local spatial concentration will remain less meaningful than the model hypothesis implies
- future ecosystem evaluation will still be missing one important route to collapse and recovery dynamics

---

## Change

Make autonomous steering respond in a bounded way to the local crowding pressure that now exists in the world.

## Why This Matters

The simulation now applies extra energy pressure in dense local neighborhoods. That improves ecosystem validity, but it also creates a new coherence tension: autonomous circles still choose movement as if dense local concentration had no direct downside until after the cost is already paid.

Without a steering adjustment, the model risks splitting into two inconsistent layers:

- the energy system says crowding matters
- the steering system still behaves as if crowding is mostly irrelevant

The next pressure is therefore to let autonomy react minimally to the local pressure already present in the simulation.

## Impacted Areas

### Simulation model

- autonomous steering should consider local crowding as one bounded influence
- the adjustment should remain deterministic and local rather than strategic or global
- the player should remain under unchanged manual steering for now

### Runtime contract

- the current snapshot shape is likely sufficient because movement and energy consequences are already visible
- build should avoid new protocol fields unless the steering adjustment becomes too opaque without a small inspectability aid

### Browser rendering

- no rendering change should be required if the effect is visible through ordinary movement
- this should remain a simulation-behavior slice rather than a UI slice

### Evaluation expectations

- EGD after this slice should be able to assess whether dense clusters now create not only extra cost but also more plausible dispersal behavior
- this directly follows the new local crowding-pressure baseline

### Existing semantics

- food placement, regeneration, fight, reproduction, continuity, child ownership, and the crowding energy rule itself should remain unchanged
- the slice should refine steering coherence, not redesign the world

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- one bounded way for local crowding to influence autonomous direction choice
- how that influence coexists with current food, threat, and interaction priorities
- a deterministic tie-break rule when crowding differences are small or equal

## Risks If Ignored

- autonomy may keep paying for crowding without ever learning to avoid it
- the new crowding-pressure rule may feel bolted on rather than integrated into world behavior
- future ecosystem evaluation will still be missing a stronger route toward dispersal dynamics

---

## Change

Improve player-facing legibility of shape meaning, danger, and opportunity using the already authoritative snapshot state.

## Why This Matters

The simulation now has substantially more world behavior than the earlier demo baseline, but EGD still highlights an important gap: too much of the meaning still depends on reading dense labels and HUD text instead of perceiving the world directly.

That creates a real evaluation problem:

- the server truth may be coherent
- but the player may still struggle to read same-shape danger, different-shape opportunity, or blocked opportunity at ordinary play speed

The next pressure is therefore not another hidden rule. It is to make existing authoritative semantics legible in the client.

## Impacted Areas

### Browser rendering

- client rendering should make shape-driven risk and opportunity easier to read at a glance
- the current dark-mode presentation should remain coherent
- the slice should prefer a few strong visual cues over more text

### Runtime contract

- the current snapshot shape is likely sufficient because shape, energy, children, and interaction outcomes are already present
- build should avoid new contract fields unless one tiny authoritative readability field becomes clearly necessary

### Evaluation expectations

- this slice directly addresses the EGD finding that play meaning is still too dependent on HUD/debug reading
- future EGD should be able to assess readability, not only mechanical completeness

### Existing semantics

- no server-side rule changes should be needed
- fight, reproduction, food, continuity, crowding pressure, and autonomy remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- which nearby semantic distinctions deserve direct visual treatment first
- how to convey danger and opportunity without turning the client into a second rules engine
- how to keep the rendering readable rather than noisy

## Risks If Ignored

- the world can become mechanically richer while still feeling opaque to the player
- evaluation may keep conflating model weakness with presentation weakness
- the simulation may remain harder to judge and enjoy than it needs to be

---

## Change

Expose which side failed the capacity check when reproduction is blocked.

## Why This Matters

The reproduction path is now much more inspectable than before. The runtime can already say which child paid for reproduction, which child was promoted in continuity, which child was absorbed in a fight, and which child triggered attached-child contact. But blocked reproduction still stops at the pair level: `reproduce_blocked_energy` says the interaction failed, while hiding whether the source side, the target side, or both sides lacked enough current capacity.

That leaves a mismatch between authoritative knowledge and inspectable output. The server already knows exactly how the current feasibility rule evaluated each participant, yet the client and tests still need to infer failure identity indirectly from energy and attached-child state.

## Impacted Areas

### Runtime contract

- blocked reproduction outcomes should expose whether the source side, target side, or both sides failed the current capacity rule
- successful reproduction outcomes should remain unchanged and should not emit blocked-capacity identity

### Simulation model

- the current feasibility rule should remain unchanged
- build should surface the result of the existing capacity check rather than inventing a second evaluation path

### Browser inspectability

- the HUD should remain sufficient to read which side failed blocked reproduction
- no larger visualization system is needed if the identity is readable through existing debug output

### Deterministic testing

- tests should prove source-only, target-only, and both-sides blocked-capacity paths
- tests should also prove that successful reproduction emits no blocked-capacity identity

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether blocked-capacity identity is exposed as one field per side rather than a more complex failure object
- whether the same fields remain absent on successful reproduction
- how deterministic tests pin source-only, target-only, and both-sides blocked cases without changing current reproduction semantics

## Risks If Ignored

- blocked reproduction will remain less inspectable than the rest of the current child and reproduction model
- debugging current reproduction feasibility will continue to rely on indirect inference instead of explicit authoritative output
- later reproduction refinements will still have to work around avoidable ambiguity in the blocked path

---

## Change

Expose the concrete child identity consumed during child-paid reproduction.

## Why This Matters

The reproduction path is now close to fully inspectable. The runtime can already say which side paid through a child, which side lacked enough current capacity when reproduction was blocked, and which child or children triggered attached-child contact. But `reproduce_paid_child` still stops short of the concrete payment-child identity: it says the source side, target side, or both used a child, while hiding which attached child was actually consumed.

That leaves a mismatch between authoritative knowledge and inspectable output. The server already knows the exact child ID removed from each paying side, yet the client and tests still need to infer it from the changed attached-child set after resolution.

## Impacted Areas

### Runtime contract

- `reproduce_paid_child` outcomes should expose the concrete consumed child ID for the source side and/or target side
- energy-only reproduction outcomes should remain unchanged and should not emit payment-child identity fields

### Simulation model

- the existing deterministic child-payment path should remain unchanged
- build should surface the child identity already chosen by the current payment rule rather than inventing a new selection rule

### Browser inspectability

- the HUD should remain sufficient to read which child was consumed as payment
- no larger visualization system is needed if the identity is readable through existing debug output

### Deterministic testing

- tests should prove source-only, target-only, and both-sides child-payment identity paths
- tests should also prove that energy-only reproduction emits no payment-child identity

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether consumed payment-child identity is exposed as one field per side rather than a more complex payment object
- whether the same fields remain absent on energy-only reproduction
- how deterministic tests pin source-only, target-only, and both-sides payment cases without changing current reproduction semantics

## Risks If Ignored

- child-paid reproduction will remain less inspectable than continuity promotion, fight absorption, and child-triggered contact
- debugging current payment behavior will continue to rely on indirect inference instead of explicit authoritative output
- later reproduction refinements will still have to work around avoidable ambiguity in the payment path

---

## Change

Expose the concrete child identities created during successful reproduction.

## Why This Matters

The reproduction path is now highly inspectable on the failure and payment sides. The runtime can already say which side paid through a child, which concrete child was consumed as payment, and which side lacked enough current capacity when reproduction was blocked. But successful reproduction still stops short of the created-child identity itself: the server knows exactly which new child IDs it allocated, while the runtime only exposes the resulting attached-child state after redistribution.

That leaves a mismatch between authoritative knowledge and inspectable output. The client and tests still need to infer newly created children by diffing before and after attached-child sets, even though the server already knows the exact created IDs.

## Impacted Areas

### Runtime contract

- successful reproduction outcomes should expose the concrete created child IDs
- blocked reproduction outcomes should remain unchanged and should not emit created-child identity

### Simulation model

- the existing deterministic child-allocation and redistribution path should remain unchanged
- build should surface the IDs already allocated by the current creation rule rather than inventing a new allocation path

### Browser inspectability

- the HUD should remain sufficient to read which child IDs were created
- no larger visualization system is needed if the identity is readable through existing debug output

### Deterministic testing

- tests should prove that energy-paid and child-paid reproduction both expose the created child IDs
- tests should also prove that blocked reproduction emits no created-child identity

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether created-child identity is exposed as one list shared across the interaction rather than split per side
- whether blocked reproduction keeps those fields absent
- how deterministic tests pin the exposed created IDs to the actual newly attached children without changing current reproduction semantics

## Risks If Ignored

- successful reproduction will remain less inspectable than the rest of the current child and reproduction model
- debugging child creation and redistribution will continue to rely on indirect set-diffing instead of explicit authoritative output
- later reproduction refinements will still have to work around avoidable ambiguity in the creation path

---

## Change

Expose which side received each created child during successful reproduction.

## Why This Matters

The reproduction path is now highly inspectable at the child-creation level. The runtime can already say which new child IDs were created, which concrete child was consumed as payment, and which side lacked enough current capacity when reproduction was blocked. But successful reproduction still stops short of the redistribution result itself: the server knows which side received each created child, while the runtime only exposes the raw created IDs and the final attached-child sets.

That leaves a mismatch between authoritative knowledge and inspectable output. The client and tests still need to infer created-child ownership by diffing per-side attached-child sets, even though the server already knows the exact deterministic allocation result.

## Impacted Areas

### Runtime contract

- successful reproduction outcomes should expose which created child IDs were allocated to the source side and which were allocated to the target side
- blocked reproduction outcomes should remain unchanged and should not emit created-child ownership

### Simulation model

- the existing deterministic redistribution path should remain unchanged
- build should surface the ownership implied by the current allocation rule rather than inventing a new redistribution path

### Browser inspectability

- the HUD should remain sufficient to read created-child ownership
- no larger visualization system is needed if the identity is readable through existing debug output

### Deterministic testing

- tests should prove that energy-paid and child-paid reproduction both expose created-child ownership
- tests should also prove that blocked reproduction emits no created-child ownership

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether created-child ownership is exposed as one field per side rather than as richer child objects
- whether blocked reproduction keeps those fields absent
- how deterministic tests pin the exposed ownership to the actual newly attached children on each side without changing current reproduction semantics

## Risks If Ignored

- successful reproduction will remain less inspectable at the redistribution level than the rest of the current child and reproduction model
- debugging allocation outcomes will continue to rely on indirect per-side set-diffing instead of explicit authoritative output
- later reproduction refinements will still have to work around avoidable ambiguity in the ownership path

---

## Change

Expose the triggering attached-child contact path kind during interaction.

## Why This Matters

The interaction path is now highly inspectable at the identity level. The runtime can already say whether contact came from `parent_body` or `attached_child`, and it can name the source-side and target-side child IDs that participated. But the runtime still does not say whether the actual trigger path was child-to-parent or child-to-child. The server knows which geometric path fired first, while the runtime only exposes enough data for the client and tests to infer that path indirectly.

That leaves a mismatch between authoritative knowledge and inspectable output. Child-triggered interaction remains less explicit about geometry than recent reproduction and continuity slices are about child usage.

## Impacted Areas

### Runtime contract

- `attached_child` interactions should expose the triggering contact path kind
- parent-body-only interactions should remain unchanged and should not emit a child contact path kind

### Simulation model

- the existing deterministic contact precedence should remain unchanged
- build should surface the chosen path kind already implied by the current geometry and precedence rule rather than inventing a new selection rule

### Browser inspectability

- the HUD should remain sufficient to read the contact path kind
- no larger visualization system is needed if the path kind is readable through existing debug output

### Deterministic testing

- tests should prove source-child-to-target-parent, source-parent-to-target-child, and child-to-child path kinds
- tests should also prove that parent-body-only contact emits no child path kind

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether contact path kind is exposed as a small enum rather than inferred from child ID presence
- whether parent-body-only contact keeps that field absent
- how deterministic tests pin each path kind without changing current contact semantics

## Risks If Ignored

- child-triggered interaction will remain less inspectable at the geometry-path level than the rest of the current child model
- debugging contact precedence will continue to rely on indirect inference instead of explicit authoritative output
- later interaction refinements will still have to work around avoidable ambiguity in the trigger path

---

## Change

Expose the deterministic redistribution kind used during successful reproduction.

## Why This Matters

Successful reproduction is now highly inspectable at the identity and ownership levels. The runtime can already say which child IDs were created and which side received each created child. But it still does not say which redistribution case was selected by the authoritative rule. The server already knows whether the result was source-only, split, or target-only, while the runtime only exposes enough detail for the client and tests to infer that kind indirectly from the per-side ownership lists.

That leaves a mismatch between authoritative knowledge and inspectable output. Redistribution remains less explicit than other recent child and contact slices even though the underlying rule is already deterministic and bounded.

## Impacted Areas

### Runtime contract

- successful reproduction outcomes should expose the selected redistribution kind
- blocked reproduction outcomes should remain unchanged and should not emit redistribution kind

### Simulation model

- the existing deterministic redistribution path should remain unchanged
- build should surface the already-selected case implied by the current allocation rule rather than inventing a new path

### Browser inspectability

- the HUD should remain sufficient to read redistribution kind
- no larger visualization system is needed if the kind is readable through existing debug output

### Deterministic testing

- tests should prove source-only, split, and target-only successful reproduction kinds
- tests should also prove that blocked reproduction emits no redistribution kind

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether redistribution kind is exposed as a small enum rather than inferred from ownership lists
- whether blocked reproduction keeps that field absent
- how deterministic tests pin the exposed kind to the actual created-child ownership on each side without changing current reproduction semantics

## Risks If Ignored

- successful reproduction will remain less inspectable at the rule-selection level than the rest of the current child and reproduction model
- debugging redistribution outcomes will continue to rely on indirect ownership inference instead of explicit authoritative output
- later reproduction refinements will still have to work around avoidable ambiguity in the redistribution path

---

## Change

Expose the concrete attached-child identity used during child-triggered interaction contact.

## Why This Matters

The simulation already treats attached children as real geometry in food collection, interaction triggering, avoidance, and positive targeting. Recent inspectability slices also made child-dependent outcomes explicit for continuity, fight absorption, and child-paid reproduction. But the contact layer still stops one step short: the runtime can say that contact came from `attached_child`, yet it cannot say which concrete child actually triggered the interaction.

That leaves a mismatch between what the authoritative server knows and what the client or tests can inspect. If a pair begins through child-to-parent or child-to-child contact, the current snapshot still requires indirect positional inference instead of making the embodied trigger explicit.

## Impacted Areas

### Runtime contract

- interaction payloads should expose the participating source-side and/or target-side child ID when `contact_origin` is `attached_child`
- parent-body-only contact should remain unchanged and should not emit child identity fields

### Simulation model

- existing contact detection and pair selection should remain unchanged
- build should surface the child identity already implied by the chosen contact geometry rather than inventing a new contact selection rule

### Browser inspectability

- the client HUD should remain sufficient to read which attached child triggered contact
- no larger visualization system is needed if the identity is readable through existing debug output

### Deterministic testing

- tests should prove source-child-only, target-child-only, and child-to-child contact identity paths
- tests should also prove that parent-body-only contact emits no child identity fields

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether interaction payloads expose one field per side rather than a more complex contact-event structure
- whether child identity remains absent for parent-body-only contact
- how deterministic tests pin the participating child IDs without altering current contact semantics

## Risks If Ignored

- child-triggered contact will remain less inspectable than child-based promotion, absorption, and payment
- the simulation will keep relying on positional inference for one of its most embodied child mechanics
- later debugging of interaction-triggering behavior will stay harder than necessary

---

## Change

Expose promoted child identity explicitly during continuity.

## Why This Matters

The continuity model is now more embodied than before: one attached child is consumed and the continuing active parent reappears at that promoted child's visible position. But the runtime still does not expose which child was promoted. That leaves continuity only partly inspectable, because the server knows the exact child identity while clients and tests must infer it indirectly from position.

The next model pressure is to make continuity identity explicit in the runtime contract without changing the actual continuity rule.

## Impacted Areas

### Runtime contract

- `death_promoted_child` outcomes should expose the promoted child identity explicitly
- the contract extension should stay small and scoped to continuity inspectability

### Simulation model

- continuity resolution already deterministically selects a promoted child
- build should surface that existing choice rather than redesigning promotion behavior

### Browser rendering

- the client or debug HUD may need a small display update so promoted-child identity is readable during demos
- this should remain inspectability-oriented rather than becoming a new effect system

### Tests

- server and integration tests should prove the exposed promoted child ID matches the consumed child and the continuing active position
- contract tests should validate the minimal schema extension

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether promoted child identity belongs inside the existing interaction payload or another equally small runtime field
- how to keep the extension continuity-specific rather than turning the snapshot into a broader event-history system
- how the client should surface the promoted child identity with minimal UI change

---

## Change

Expose child-payment identity explicitly during `reproduce_paid_child`.

## Why This Matters

The reproduction model is now partially inspectable: `reproduce_paid_child` distinguishes child-funded payment from ordinary energy-funded reproduction. But the runtime still does not expose which participant actually paid through a child. That leaves the server with exact payer knowledge that clients and tests can only infer indirectly from before-and-after child sets.

The next model pressure is to make child-payment identity explicit in the runtime contract without changing the actual reproduction rule.

## Impacted Areas

### Runtime contract

- `reproduce_paid_child` outcomes should expose which participant or participants paid through a child
- the contract extension should stay small and scoped to reproduction inspectability

### Simulation model

- child-payment logic already deterministically decides whether the player, the opponent, or neither uses a child
- build should surface that existing choice rather than redesigning payment behavior

### Browser rendering

- the client or debug HUD may need a small display update so child-payment identity is readable during demos
- this should remain inspectability-oriented rather than becoming a new reproduction effect system

### Tests

- server and integration tests should prove the exposed payer identity matches the actual child-payment path
- contract tests should validate the minimal schema extension

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether child-payment identity belongs inside the existing interaction payload or another equally small runtime field
- how to keep the extension reproduction-specific rather than turning the snapshot into a broader event-history system
- how the client should surface child-payment identity with minimal UI change

## Risks If Ignored

- `reproduce_paid_child` will remain less inspectable than `death_promoted_child` and `fight_absorbed_child`
- tests will keep inferring payer identity indirectly instead of reading the authoritative source
- later reproduction refinements will have weaker runtime evidence about which participant actually consumed a child for payment

---

## Change

Expose absorbed child identity explicitly during `fight_absorbed_child`.

## Why This Matters

The fight model is now already partially inspectable: `fight_absorbed_child` distinguishes child absorption from full defeat, and the loser remains active. But the runtime still does not expose which child was absorbed. That leaves the server with exact child-loss knowledge that clients and tests can only infer indirectly from the remaining attached-child set.

The next model pressure is to make absorbed-child identity explicit in the runtime contract without changing the actual fight rule.

## Impacted Areas

### Runtime contract

- `fight_absorbed_child` outcomes should expose the absorbed child identity explicitly
- the contract extension should stay small and scoped to fight inspectability

### Simulation model

- fight absorption already deterministically selects and consumes one child
- build should surface that existing choice rather than redesigning combat behavior

### Browser rendering

- the client or debug HUD may need a small display update so absorbed-child identity is readable during demos
- this should remain inspectability-oriented rather than becoming a new combat effect system

### Tests

- server and integration tests should prove the exposed absorbed child ID matches the consumed child
- contract tests should validate the minimal schema extension

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether absorbed child identity belongs inside the existing interaction payload or another equally small runtime field
- how to keep the extension fight-specific rather than turning the snapshot into a broader event-history system
- how the client should surface absorbed-child identity with minimal UI change

## Risks If Ignored

- `fight_absorbed_child` will remain less inspectable than `death_promoted_child`
- tests will keep inferring child loss indirectly instead of reading the authoritative source
- later fight refinements will have weaker runtime evidence about which visible child was actually consumed

---

## Change

Remove dead derived-radius state from the server.

## Why This Matters

The embodied-radius transition is behaviorally complete: parent radius no longer changes food reach, contact reach, fight tie-breaks, movement boundaries, orbit distance, or rendered body size. But the server still carries dead radius-growth structure such as `derivedRadius(...)` and `DefaultChildRadiusGain`, which now only disguise the fact that parent radius is fixed.

The next model pressure is to make the implementation say what the simulation already does, instead of preserving inert growth scaffolding.

## Impacted Areas

### Simulation model

- parent radius initialization and child-sync paths should express fixed radius directly
- dead growth helpers and constants should be removed or collapsed
- gameplay behavior should remain unchanged

### Tests

- existing fixed-radius tests should continue to prove that child-related state changes do not alter parent radius
- tests or implementation notes that still describe active child-driven radius growth need alignment

### Runtime contract

- no wire-level change is expected
- this slice should preserve the current snapshot shape exactly

### Implementation memory

- implementation notes should stop describing child-derived parent radius growth as active behavior once build removes the dead state

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether to delete dead helpers outright rather than leave no-op abstractions for hypothetical future growth
- how to keep the cleanup bounded to internal state and documentation alignment
- which tests should act as the fixed-radius regression proof after dead growth state is removed

## Risks If Ignored

- the code will keep implying a growth mechanic that no longer exists
- future refinement will have to reason through dead scaffolding when changing parent-body semantics
- implementation memory will stay partially inconsistent about whether radius growth is active or only historical

## Risks If Ignored

- continuity will remain less inspectable than contact origin and child-paid reproduction
- tests will keep inferring promotion identity indirectly from position instead of reading the authoritative source
- later continuity work will have weaker runtime evidence about which visible child actually became the continuing line

---

## Change

Remove the mirrored derived child-count field from Go-side snapshots.

## Why This Matters

Attached children are already the single authoritative child state in the simulation, and the runtime JSON contract already removed `children_count`. But the Go-side snapshot/test path still carries a mirrored `ChildrenCount` convenience field. That means child state is still represented twice inside the server/test boundary even though one representation is already fully derivable from the other.

The next model pressure is to complete the single-source-of-truth move consistently across the in-process boundary, not only on the wire.

## Impacted Areas

### Simulation model

- core gameplay semantics should remain unchanged
- only snapshot-facing representation and local derived reads should change

### Go-side snapshot consumers

- deterministic tests and any local server-side readers must derive child quantity from attached children
- convenience assertions that still use `ChildrenCount` need to move to `len(AttachedChildren)`

### Runtime contract

- no wire-level change is expected
- this slice should preserve the current contract shape exactly

### Implementation memory

- implementation notes should stop describing Go-side child quantity as intentionally duplicated once build removes the mirrored field

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether to remove the mirrored field outright instead of keeping it for test convenience
- how to update Go-side tests without reintroducing another child-count helper abstraction
- how to keep the change bounded to representation cleanup rather than semantic rewrite

## Risks If Ignored

- child state will remain duplicated inside the Go-side boundary after the runtime contract has already been simplified
- tests will keep validating one derived convenience field instead of the actual attached-child representation
- future child-model changes will still have to account for unnecessary in-process duplication

---

## Change

Remove child-derived growth from the visible parent body.

## Why This Matters

The current simulation now uses attached children as the main visible embodiment of accumulated children, and most embodied mechanics no longer depend on derived parent radius. But snapshots and rendering still enlarge the parent body itself with child count. That makes accumulated children visible twice: once as orbiting children and again as a larger parent core.

The next model pressure is to make the visible parent body match the fixed core implied by the embodied feeding, contact, movement, and orbit slices, while leaving non-visual child leverage explicit for now.

## Impacted Areas

### Simulation model

- parent `radius` should become a fixed visible-body value rather than a child-derived growth value
- player and autonomous parent circles should continue sharing one deterministic visible-body rule
- attached children should remain the visible embodiment of accumulated children

### Runtime contract

- the current snapshot shape is likely sufficient because `radius`, `children_count`, and attached-child positions are already visible
- build may need no contract change if the fixed visible body is observable from ordinary snapshots

### Browser rendering

- current rendering may already be enough because parent radius and attached children are already drawn
- no visual-style change should be required beyond the fixed body size itself

### Existing semantics

- child-based fight power, reproduction payment, and continuity can remain unchanged in this slice
- food, contact, movement, orbit distance, reproduction, and steering should remain unchanged
- later slices may still choose to remove additional child-count shortcuts from non-visual semantics

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether parent `radius` should now simply equal the fixed visible parent-core size in all snapshots
- whether existing child-based fight power, payment, and continuity stay unchanged while visible growth is removed
- how tests should prove that attached children remain visible while parent-body growth disappears

## Risks If Ignored

- the simulation will keep showing accumulated children twice in the visible geometry
- the embodied child model will remain partially contradicted by a larger parent core
- later efforts to narrow or remove remaining child-count shortcuts will still have to work around a visually doubled growth model
- how post-reproduction child ownership is assigned while remaining deterministic
- whether current radius growth stays temporarily active alongside orbiting children

## Risks If Ignored

- reproduction output will keep drifting away from the intended game feel
- later changes may require replacing tests and client assumptions built around free spawned child circles
- child-related semantics will remain split between counters and bodies without an explicit bridge

---

## Change

Allow shape-aware interaction-seeking autonomy to consider whether preferred different-shape targets are currently feasible for reproduction.

## Why This Matters

The current implementation now makes shape matter in target choice, which better aligns steering with collision meaning. But it still treats every different-shape target as attractive even when the current reproduction rule would be blocked by insufficient energy or missing child reserve. That keeps steering partially disconnected from the executable rules already in the model.

The next model pressure is to make preferred social pursuit more semantically honest by considering whether the currently preferred reproduction path could actually succeed.

## Impacted Areas

### Simulation model

- interaction-target ordering should consider current reproduction feasibility, not just shape
- deterministic feasibility checks become part of the steering contract
- current movement, energy, and downstream interaction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because energy, children, positions, and outcomes are already visible
- build may need steering provenance only if feasibility-driven choice is too hard to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because circle shapes, energy, and the interaction HUD are visible
- minor HUD adjustments may still help if the new target choice is subtle

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current shape-aware preference should become more precise rather than being removed
- player-targetable and autonomous-targetable steering should continue sharing one candidate set

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the exact feasibility basis for preferring different-shape targets
- the fallback rule when no feasible preferred target exists
- how deterministic tie-breaking behaves across equally feasible candidates

## Risks If Ignored

- steering will keep preferring blocked reproduction opportunities
- autonomous behavior will remain more arbitrary than the rest of the reproduction model
- later ecosystem slices will keep carrying a gap between target pursuit and target feasibility

---

## Change

Allow interaction-seeking autonomous circles to consider whether fallback same-shape targets are currently survivable under the existing fight ordering.

## Why This Matters

The current implementation now treats feasible reproduction as a meaningful steering preference, which makes different-shape pursuit more honest. But when no feasible reproduction target exists, same-shape fallback pursuit is still blind to the existing deterministic fight rule. That means circles can still steer directly into clearly losing conflicts even though the model already knows how those conflicts resolve.

The next model pressure is to make pre-contact steering better match the current same-shape conflict semantics without escalating into a full tactical AI system.

## Impacted Areas

### Simulation model

- fallback target ordering should consider whether a same-shape encounter is currently survivable
- the current deterministic fight ordering becomes part of steering eligibility, not only post-contact resolution
- current movement, energy, fight, and reproduction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because energy, children, shapes, and outcomes are already visible
- build may need steering provenance only if fight-feasibility-driven fallback is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because the HUD shows shapes, children, energy, and interaction outcomes
- minor HUD adjustments may still help if skipping a losing same-shape target is hard to observe

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current reproduction-feasibility-aware preference should remain in place
- player-targetable and autonomous-targetable steering should continue sharing one candidate set
- the current fight ordering should remain the single authoritative basis for deciding whether a same-shape fallback is currently acceptable

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the exact fight-feasibility basis for allowing same-shape fallback pursuit
- the fallback rule when no same-shape target is currently survivable
- how deterministic tie-breaking behaves across equally survivable candidates

## Risks If Ignored

- steering will keep walking into deterministic losing fights as a default fallback
- children and energy will remain more meaningful in post-contact fight resolution than in pre-contact behavior

---

## Change

Allow positive interaction-seeking autonomy to treat attached-child positions as part of effective social nearness.

## Why This Matters

The current implementation already lets attached children matter for contact triggering and for negative steering through avoidance. But once an autonomous circle enters ordinary interaction-seeking mode, target choice still mostly depends on parent-core distance. That leaves the embodied model uneven: the simulation already says children are part of encounter geometry, but positive pursuit only partially reflects that fact.

The next model pressure is to align positive steering with the same visible child geometry that already shapes contact and retreat.

## Impacted Areas

### Simulation model

- interaction-target distance checks should consider the parent body and its attached children
- current reproduction-feasibility-aware and fight-feasibility-aware eligibility rules should stay unchanged
- current movement, energy, and downstream interaction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because attached-child positions are already visible
- build may need steering provenance only if child-aware pursuit is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because attached-child positions and contact outcomes are visible
- minor HUD adjustments may still help if child-aware target choice is hard to observe

### Existing semantics

- low-energy food recovery should remain the first steering priority
- current same-shape threat and blocked-reproduction avoidance should remain unchanged
- current social target eligibility should remain governed by reproduction feasibility and fight feasibility rather than by a new scoring system

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether social nearness is measured from the parent center only or from the nearest current parent-or-child contact point
- how deterministic tie-breaking behaves when multiple attached-child paths are equally near
- whether steering should point at the effective contact point or only use it for ranking

## Risks If Ignored

- pursuit will remain less embodied than contact and avoidance
- visible orbiting children will keep mattering more for reactive behavior than for proactive behavior
- later autonomy slices may accumulate more asymmetry between approach and retreat semantics

---

## Change

Allow food-targeting autonomy to treat attached-child positions as part of effective food nearness.

## Why This Matters

The current implementation already lets attached children matter for actual food collection, for social contact, and for both negative and positive social steering. But food targeting still mostly depends on parent-core distance. That leaves the core energy loop uneven: the simulation already says a child can collect a food slot, while the targeting logic still behaves as if only the parent center can effectively reach it.

The next model pressure is to align food approach with the same visible child geometry that already shapes food collection and social behavior.

## Impacted Areas

### Simulation model

- food-target distance checks should consider the parent body and its attached children
- current low-energy and nearby-food priority rules should stay unchanged
- current movement, collection, regeneration, and social interaction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because attached-child positions and food positions are already visible
- build may need steering provenance only if child-aware food pursuit is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because attached-child positions, food positions, and energy changes are visible
- minor HUD adjustments may still help if child-aware food choice is hard to observe

### Existing semantics

- low-energy food recovery should remain the first steering priority
- current same-shape threat and blocked-reproduction avoidance should remain unchanged
- current social target rules should remain unchanged
- food collection, energy gain, and regeneration should remain the existing authoritative path

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether food nearness is measured from the parent center only or from the nearest current parent-or-child reach point
- how deterministic tie-breaking behaves when multiple child-based food paths are equally near
- whether steering should point toward the chosen food slot exactly as before or expose additional provenance

## Risks If Ignored

- the energy loop will remain less embodied than the social interaction loop
- visible orbiting children will continue to matter for actual collection more than for food-seeking intent
- later autonomy slices may accumulate avoidable asymmetry between feeding and social steering

---

## Change

Make child-paid reproduction visibly grounded in attached-child state.

## Why This Matters

The current implementation already allows a circle to use one child as reserve payment when energy is below the reproduction cost. That keeps the gameplay rule moving, but the visible result can still feel too similar to ordinary reproduction because the payment and the reward are folded together in the final child counts. Attached children are now embodied across feeding, contact, avoidance, pursuit, conflict absorption, and continuity, so reproduction payment is the clearest remaining place where child use can still feel hidden.

The next model pressure is to make child-paid reproduction observably different from energy-only reproduction without changing the broader reproduction model.

## Impacted Areas

### Simulation model

- reproduction resolution should visibly reflect when one participant paid through child reserve
- current reproduction feasibility, cost, and deterministic redistribution rules should stay unchanged unless a minimal explicit adjustment is required

### Runtime contract

- the current snapshot shape may already be sufficient because child counts, attached-child arrays, and interaction outcomes are visible
- build may need one explicit reproduction-outcome distinction only if ordinary snapshots remain too ambiguous

### Browser rendering

- current rendering may already be enough because attached-child state is visible
- minor HUD adjustments may still help if child-paid reproduction is hard to distinguish from energy-only reproduction

### Existing semantics

- current orbiting-child representation should remain the authoritative child embodiment
- current food, fight, continuity, and movement rules should remain unchanged
- player and autonomous circles should follow the same visible child-payment rule

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether ordinary child-count and attached-child changes are enough to make child-paid reproduction legible
- whether one explicit reproduction outcome distinction is needed for inspectability
- how visible child payment coexists with the current deterministic redistribution rule

## Risks If Ignored

- reproduction payment will remain less embodied than other child-related mechanics
- players may continue to read child-paid reproduction as a hidden or net-zero bookkeeping trick
- later efforts to remove transitional count-based shortcuts will have weaker visible support

---

## Change

Reduce the feeding shortcut by making food collection depend on visible parent and attached-child bodies rather than enlarged derived radius alone.

## Why This Matters

The current implementation now makes attached children matter across nearly every core loop, including food targeting and actual child-based collection. But food collection still benefits from the abstract radius growth shortcut, which means the energy loop remains only partially embodied. A circle can still collect food because of a hidden reach expansion even when no visible body touches the slot.

The next model pressure is to reduce one shortcut where the embodied child model is already strong, without trying to remove radius-based abstraction everywhere at once.

## Impacted Areas

### Simulation model

- food consumption checks should depend on parent-core and attached-child overlap rather than enlarged parent radius alone
- current energy gain and food regeneration should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because food positions, parent bodies, and attached children are already visible
- build should add provenance only if embodied non-collection is too subtle to infer

### Browser rendering

- current rendering may already be enough because visible bodies and food positions are shown
- minor HUD changes may still help if “large radius but no collection” is hard to read

### Existing semantics

- food targeting should remain unchanged
- fight, reproduction, continuity, and social steering should remain unchanged
- radius may remain a visual and other-domain property for now even if feeding stops using it as silent reach

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether parent-core overlap and attached-child overlap are the only valid food collection paths
- how to preserve exact-once food consumption when multiple visible bodies overlap the same slot
- whether the current snapshot contract is sufficient to explain embodied non-collection

## Risks If Ignored

- the energy loop will remain partly governed by hidden geometric reach
- the visible orbiting-child model will still be weaker in feeding than it appears
- later removal of broader count-based shortcuts will remain harder to stage cleanly

---

## Change

Reduce the encounter shortcut by making circle-to-circle contact depend on visible parent-core and attached-child bodies rather than enlarged derived parent radius alone.

## Why This Matters

The current implementation now makes food collection depend on visible bodies, but encounter initiation still uses enlarged parent radius for parent-body overlap. That means fights and reproduction can still begin because of a hidden grown reach even when no visible body touches. After tightening food collection, this is the clearest remaining embodied-versus-abstract mismatch.

The next model pressure is to reduce one more shortcut where the orbiting-child model already provides a concrete visible geometry.

## Impacted Areas

### Simulation model

- player-autonomous and autonomous-autonomous contact checks should depend on visible parent-core and attached-child overlap rather than enlarged parent radius alone
- current fight and reproduction outcome rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because parent bodies, attached children, and `contact_origin` are already visible
- build should add provenance only if embodied non-contact is too subtle to infer

### Browser rendering

- current rendering may already be enough because visible bodies and interaction outcomes are shown
- minor HUD changes may still help if “grown radius but no contact” is hard to read

### Existing semantics

- child-triggered parent and child-to-child contact should remain unchanged
- food collection should remain unchanged
- fight winner ordering, reproduction payment, and continuity rules should remain unchanged
- radius may remain a visual and other-domain property for now even if it stops silently enlarging contact reach

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether visible parent-core and attached-child overlap are the only valid contact paths
- how to preserve exact-once pair contact when multiple visible bodies overlap
- whether the current snapshot contract is sufficient to explain embodied non-contact

## Risks If Ignored

- fights and reproduction will remain partly governed by hidden geometric reach
- the visible orbiting-child model will still be weaker in encounter initiation than it appears
- later removal of broader count-based shortcuts will remain harder to stage cleanly

---

## Change

Remove derived radius as a same-shape fight tie-break once energy and child power are already considered.

## Why This Matters

The current implementation has now tightened feeding reach and contact reach to visible embodied geometry, and child count already exists as an explicit fight-power input. But same-shape fights still fall back to derived radius after energy and child count, which leaves one remaining hidden growth shortcut deciding combat even after child power became explicit.

The next model pressure is to remove that radius-based combat shortcut before taking on broader questions about visual size or movement boundaries.

## Impacted Areas

### Simulation model

- same-shape fight winner selection should no longer use derived radius after energy and child count
- exact ties still need one deterministic final rule
- contact initiation, child absorption, and continuity should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because energy, child count, winner, loser, and interaction kind are already visible
- build should add provenance only if the removed radius tie-break is too subtle to infer

### Browser rendering

- current rendering may already be enough because radius can remain visible even if it no longer decides same-shape fights
- minor HUD changes may still help if users confuse visual size with current combat tie-breaks

### Existing semantics

- energy should remain the primary fight input
- child count should remain the next explicit fight-power input
- reproduction, food, contact, and continuity should remain unchanged
- radius may remain a visual and other-domain property for now even if it stops deciding same-shape ties

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- what deterministic rule replaces radius once energy and child count tie
- whether player and autonomous exact-tie handling stay the same
- whether the current snapshot contract is sufficient to explain the removed radius tie-break

## Risks If Ignored

- fight resolution will remain more abstract than the now-embodied feeding and contact loops
- visible orbiting children will still be weaker in combat than the model suggests
- later removal of remaining count-based shortcuts will stay harder to stage cleanly
- later autonomy slices will keep carrying a gap between steering and the current fight model

---

## Change

Allow autonomous circles to steer away from nearby stronger same-shape threats under one bounded deterministic rule.

## Why This Matters

The current implementation now avoids pursuing clearly losing same-shape fights, which improves target choice. But it still does not let circles respond explicitly when a stronger same-shape threat is already nearby. That means the model still lacks a negative steering expression for hostile asymmetry: circles can decline pursuit, but they do not yet actively avoid immediate danger.

The next model pressure is to make the existing fight model visible in motion before contact, not only in target eligibility and post-contact resolution.

## Impacted Areas

### Simulation model

- steering needs one explicit threat-avoidance step based on the current deterministic fight ordering
- proximity becomes part of whether fight semantics affect motion before contact
- current movement, energy, fight, reproduction, and food rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because positions, shapes, children, energy, and outcomes are already visible
- build may need steering provenance only if avoidance is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because the demo shows shapes and live motion
- minor HUD cues may still help if avoidance is hard to distinguish from food-seeking in some layouts

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current reproduction-feasibility-aware and fight-feasibility-aware pursuit rules should remain in place
- player-targetable and autonomous-targetable steering should continue sharing one candidate set for threat evaluation
- the existing deterministic fight ordering should remain the sole basis for deciding whether a nearby same-shape target is threatening

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the proximity threshold for threat avoidance
- whether threat avoidance runs before or after ordinary pursuit selection
- how deterministic tie-breaking behaves across multiple nearby threats

## Risks If Ignored

- autonomous circles will still feel too passive in the presence of immediately stronger same-shape threats
- the fight model will stay more visible in outcomes than in motion
- later ecosystem slices will keep carrying a gap between threat presence and steering response

---

## Change

Allow autonomous circles to steer away from nearby different-shape targets when reproduction is currently blocked under the existing feasibility rule.

## Why This Matters

The current implementation now expresses explicit negative steering for one kind of social asymmetry: nearby stronger same-shape threats. But different-shape behavior is still incomplete. A nearby different-shape circle that cannot currently reproduce is only non-preferred, not explicitly avoided. That leaves blocked reproduction less visible in motion than losing fight threat.

The next model pressure is to make the current reproduction model shape motion negatively as well as positively before contact occurs.

## Impacted Areas

### Simulation model

- steering needs one explicit blocked-reproduction avoidance step based on the current reproduction feasibility rule
- proximity becomes part of whether blocked reproduction affects motion before contact
- current movement, energy, fight, reproduction, and food rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because positions, shapes, children, energy, and outcomes are already visible
- build may need steering provenance only if blocked-reproduction avoidance is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because the demo shows shapes, labels, and live motion
- minor HUD cues may still help if blocked different-shape retreat is hard to distinguish from same-shape threat retreat in some layouts

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current same-shape threat-avoidance rule should remain in place
- the current reproduction-feasibility-aware preference should remain in place for viable different-shape targets
- player-targetable and autonomous-targetable steering should continue sharing one candidate set for blocked-target evaluation
- the existing reproduction feasibility rule should remain the sole basis for deciding whether a nearby different-shape target is blocked

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- the proximity threshold for blocked-reproduction avoidance
- whether blocked-reproduction avoidance runs before or after same-shape threat avoidance
- how deterministic tie-breaking behaves across multiple blocked different-shape targets

## Risks If Ignored

- blocked reproduction will remain less legible in motion than losing fight threat
- energy and child reserve will stay more visible in different-shape outcomes than in different-shape steering
- later autonomy slices will keep carrying a gap between blocked reproduction presence and steering response

---

## Change

Allow attached children to contribute to same-shape threat and blocked-reproduction avoidance detection.

## Why This Matters

The current implementation now lets orbiting children matter in several places:

- they are visible in the world
- they collect food
- they can absorb conflict loss
- they can be consumed for continuity and reproduction payment
- they can already trigger authoritative parent-level contact before parent cores overlap

But the newer avoidance rules still mostly reason from parent-core proximity. That creates a clean remaining mismatch: child-triggered contact is already real, while child-triggered avoidance is not.

The next model pressure is to align pre-contact avoidance with the embodied orbiting-child geometry already used by the interaction engine.

## Impacted Areas

### Simulation model

- avoidance detection should be able to treat attached-child positions as part of current proximity
- the current same-shape threat and blocked-reproduction categories should remain unchanged
- current movement, energy, fight, reproduction, and food rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because attached-child positions are already exposed
- build may need steering provenance only if child-triggered avoidance is too subtle to infer from ordinary snapshots

### Browser rendering

- current rendering is likely sufficient because attached children are already visible
- minor HUD cues may still help if distinguishing parent-triggered from child-triggered retreat is too subtle

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current same-shape threat-avoidance and blocked-reproduction avoidance rules should remain in place
- the current child-triggered contact model should become more consistent with pre-contact motion
- the existing deterministic child layout should remain the sole geometry source for this slice

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether parent-body and attached-child proximity use the same avoidance thresholds
- how deterministic tie-breaking behaves across parent-triggered and child-triggered avoidance paths
- whether both same-shape threat and blocked-reproduction avoidance gain child-awareness in the same slice

## Risks If Ignored

- visible orbiting children will keep mattering more for contact than for avoidance
- avoidance will remain semantically behind the current interaction geometry
- later autonomy slices will keep carrying a gap between child-triggered contact and child-triggered retreat

---

## Change

Allow interaction-seeking autonomous circles to prioritize targets by shape meaning instead of using only nearest-target selection.

## Why This Matters

The current implementation now lets autonomous circles decide between food recovery and social interaction based on energy. But once they enter interaction-seeking mode, target choice is still shape-blind even though shape is the variable that changes interaction semantics.

The next model pressure is to make shape influence target ordering before contact so steering and outcome meaning are less disconnected.

## Impacted Areas

### Simulation model

- interaction-target ordering should consider whether a candidate implies fight or reproduction
- deterministic priority and tie-breaking become part of the steering contract
- current movement, energy, and downstream interaction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because positions, shapes, and outcomes are already visible
- build may need steering provenance only if shape-driven choice is too hard to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because circle shapes and the interaction HUD are visible
- minor HUD adjustments may still help if the new target preference is subtle

### Existing semantics

- the current low-energy food-recovery rule should remain unchanged
- the current food-priority distance rule should remain coherent
- player-targetable and autonomous-targetable steering should continue sharing one candidate set

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- which shape outcome is preferred when interaction seeking is active
- how fallback behaves when no preferred-outcome target exists
- how equal-distance and equal-priority cases remain deterministic

## Risks If Ignored

- shape will remain too absent from pre-contact behavior
- interaction-seeking autonomy will feel more arbitrary than the rest of the model
- later ecosystem slices will keep carrying a gap between steering meaning and collision meaning

---

## Change

Allow autonomous steering priority to depend on energy condition so low-energy circles prefer food recovery over social interaction.

## Why This Matters

The current implementation lets autonomous circles seek other circles, including the player, which makes the ecosystem more active. But the steering rule still treats a nearly starving circle much like a healthy one whenever food is not inside the current priority distance. That leaves energy underused as a behavioral variable even though it is already the central survival variable in the model.

The next model pressure is to make energy shape autonomous target choice before contact, not only after collisions and movement costs have already been applied.

## Impacted Areas

### Simulation model

- autonomous steering should consider energy level when choosing between food and interaction targets
- deterministic threshold selection becomes part of the steering contract
- current movement, energy recovery, and downstream interaction rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because energy values, positions, and outcomes are already visible
- build may need steering provenance only if the energy-driven choice is hard to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because energy and motion are visible
- minor HUD adjustments may still help if the change is too subtle in the demo

### Existing semantics

- the current food-priority distance rule should remain coherent with the new threshold
- player-targetable and autonomous-targetable interaction seeking should still apply once a circle is sufficiently energized
- the model should remain simple and avoid sliding into personality systems

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- one deterministic energy threshold for preferring food over interaction
- whether the threshold applies before or after the current food-priority-distance check
- how equal-distance and equal-condition tie-breaking remains deterministic

## Risks If Ignored

- energy will remain too weak as a pre-contact behavioral variable
- autonomous steering will feel more arbitrary than survival-driven
- later ecosystem slices will have a weaker behavioral basis for collapse and recovery dynamics

---

## Change

Allow autonomous circles to include the player in deterministic interaction-seeking target selection.

## Why This Matters

The current implementation now lets autonomous circles actively create encounters with other autonomous circles, which strengthens emergence. But the player is still excluded from that steering layer, leaving one more structural asymmetry in a model that says the player should be one participant inside the same ecosystem.

The next model pressure is to let the world engage the player under the same bounded steering logic instead of requiring the player to initiate most contact.

## Impacted Areas

### Simulation model

- autonomous target selection should consider the player as an eligible target when active
- deterministic ordering becomes more important when the player and autonomous circles compete as potential targets
- current movement, energy, and downstream interaction rules should remain unchanged after target choice

### Runtime contract

- the current snapshot shape is likely sufficient because movement and interaction outcomes are already visible
- build may need steering provenance only if player-targeting is hard to infer from ordinary snapshots

### Browser rendering

- current rendering may already be enough because the player, autonomous circles, and interaction HUD are visible
- minor HUD adjustments may still be useful if player-targeting becomes ambiguous in the demo

### Existing semantics

- the current food-priority rule should remain coherent
- player-autonomous fight and reproduction semantics should remain unchanged
- the model should still avoid slipping into explicit hostility or role systems

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- when the player becomes eligible relative to other autonomous targets
- how deterministic tie-breaking works when the player and another circle are equally valid
- whether the player is always eligible while active or only under a bounded condition

## Risks If Ignored

- the ecosystem will remain only partially shared from the player’s perspective
- the player will stay too exempt from non-player initiative
- later emergence-oriented slices will keep carrying this remaining asymmetry

---

## Change

Allow autonomous circles to steer toward interaction opportunities instead of relying almost entirely on food-seeking and drift.

## Why This Matters

The current implementation can now resolve autonomous-autonomous encounters, but autonomous circles still do not actively help create those encounters. That keeps the ecosystem more passive than the model hypothesis suggests and leaves the new shared-rule interaction engine underexercised in ordinary play.

The next model pressure is to make non-player circles behave more like participants in the ecosystem without introducing explicit AI personalities or strategy trees.

## Impacted Areas

### Simulation model

- autonomous target-selection logic should consider other circles, not only food
- deterministic target ordering becomes more important when several eligible circles exist
- current movement, energy, and interaction outcome rules should remain unchanged after target selection

### Runtime contract

- the current snapshot shape is likely sufficient because motion and interaction outcomes are already visible
- build may need steering provenance only if the new motion is too hard to understand from ordinary snapshots

### Browser rendering

- current rendering may already be enough because circle positions, labels, and the interaction HUD are visible
- minor HUD improvements may still be useful if autonomous target choice becomes hard to infer

### Existing semantics

- current food-seeking autonomy should coexist with the new interaction-seeking rule
- player-autonomous and autonomous-autonomous outcomes should remain governed by the same downstream rules
- the model should avoid slipping into hard-coded roles or tactical AI framing

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- when interaction seeking takes priority over food seeking
- how target choice is ordered deterministically across multiple circles
- whether the player is always an eligible target or only under bounded conditions

## Risks If Ignored

- autonomous-autonomous interaction will remain correct but underutilized
- the ecosystem will continue to depend too much on player movement or lucky geometry
- later emergence-oriented slices will rest on still-passive non-player steering

---

## Change

Allow autonomous circles to resolve interactions with other autonomous circles under the same current rules used for player-involved encounters.

## Why This Matters

The current implementation has improved embodiment and fairness in many local rules, but the world is still structurally player-centered: autonomous circles do not yet resolve encounters with each other. That limits emergence because meaningful state change still depends too heavily on player participation.

The next model pressure is to let the ecosystem evolve through non-player encounters while preserving the current deterministic fight and reproduction semantics.

## Impacted Areas

### Simulation model

- overlap detection should extend beyond player-autonomous pairs to autonomous-autonomous pairs
- deterministic ordering becomes more important if several autonomous pairs overlap in one tick
- current child, energy, reproduction, and continuity rules must remain coherent for autonomous-only encounters

### Runtime contract

- the current snapshot shape is likely sufficient because autonomous circle states and the existing interaction object are already visible
- build may need a slightly clearer source/target convention only if autonomous-autonomous outcomes are hard to read

### Browser rendering

- current rendering may already be sufficient because autonomous circles and interaction labels are visible
- minor label improvements may still be needed if autonomous-autonomous results are ambiguous in the demo

### Existing semantics

- player-autonomous interaction should remain unchanged
- current contact-origin handling should continue to work when autonomous pairs use body or child contact
- the ecosystem becomes less player-centered without requiring new AI policies yet

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- how autonomous-autonomous pairs are ordered deterministically when more than one overlap exists
- whether only one autonomous-autonomous interaction resolves per tick or whether it can coexist with a player-involved interaction in the same tick
- whether the existing interaction object is readable enough for autonomous-only encounters

## Risks If Ignored

- the simulation will remain artificially dependent on player involvement
- fairness-through-shared-rules will stay incomplete
- later ecosystem slices will rest on a still player-centered interaction engine

---

## Change

Treat attached-child-to-attached-child overlap as a valid trigger for parent-level interaction.

## Why This Matters

The current implementation now lets orbiting children trigger parent-level interaction when they touch another parent body. That improved the embodied contact model, but it still leaves a visible inconsistency: two child swarms can meet each other without authoritative meaning.

The next model pressure is to complete the minimum contact interpretation by allowing attached-child-to-attached-child overlap to trigger the same parent-level fight or reproduction paths.

## Impacted Areas

### Simulation model

- overlap detection should consider child-to-child contact across a parent pair
- de-duplication needs to stay deterministic when parent-body, child-to-parent, and child-to-child contact all happen in the same tick
- current repeated-reproduction overlap windows must remain coherent when the triggering path is child-to-child

### Runtime contract

- the current contract may remain sufficient if `contact_origin: attached_child` is still enough for inspectability
- build may need a finer-grained provenance value only if child-to-child contact cannot be understood from current snapshots

### Browser rendering

- current orbit rendering may already be enough to explain child-to-child contact
- minor label or debug adjustments may still be needed if provenance is too ambiguous in the live demo

### Existing semantics

- same-shape and different-shape meaning should remain unchanged after contact is detected
- current parent-body and child-to-parent contact rules should coexist cleanly with child-to-child contact
- player and autonomous circles must follow the same rule

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether the existing `attached_child` provenance is sufficient for both child-to-parent and child-to-child contact
- how to prevent duplicate same-tick outcomes when several child bodies overlap across one pair
- whether any child-to-child-specific outcome is needed or remains out of scope

## Risks If Ignored

- orbiting children will still feel only partially real as contact bodies
- the visible swarm model will remain more expressive than the authoritative interaction model
- later attempts to lean less on radius shortcuts will have a weaker embodied basis

---

## Change

Treat attached orbiting children as valid contact points for triggering parent-level interactions.

## Why This Matters

The current implementation has made attached children visible and useful in feeding, payment, conflict absorption, and continuity. But interaction initiation is still mostly parent-core-driven. That keeps the orbiting-child model partially embodied: children matter after outcomes, but not enough in how encounters begin.

The next model pressure is to let visible orbiting children participate directly in authoritative contact detection while still preserving the current parent-level fight and reproduction rules.

## Impacted Areas

### Simulation model

- overlap detection should consider attached-child positions against other parent bodies
- one parent pair still needs deterministic de-duplication so child contact does not produce duplicate same-tick outcomes
- the current overlap-window rule for repeated reproduction must remain coherent when contact originates from a child

### Runtime contract

- the current snapshot may remain sufficient because attached-child positions and interaction outcomes are already visible
- build may need one extra outcome marker only if ordinary snapshots are not enough to explain child-originated contact

### Browser rendering

- no major rendering feature is required if current orbiting-child visuals already make the contact path inspectable
- labels or debug text may still need minor updates if child-originated contact is too ambiguous in the demo

### Existing semantics

- same-shape and different-shape interaction meaning should remain unchanged after contact is detected
- current radius shortcuts remain transitional and should coexist with child-originated contact in this slice
- player and autonomous circles should follow the same child-contact rule

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether attached-child contact counts only against parent bodies or also against other attached children
- how to prevent one pair from triggering the same interaction twice in one tick when parent and child both overlap
- whether the existing interaction object is sufficient for inspectability

## Risks If Ignored

- orbiting children will remain mechanically important but still not central to encounter initiation
- later removal of radius shortcuts will be harder because contact logic will still be too parent-core-centric
- the visible orbital model will continue to lag behind the authoritative interaction model

---

## Change

Make attached orbiting children absorb hostile loss during same-shape conflict before a whole parent circle disappears.

## Why This Matters

The current implementation made children visible and consumable for reproduction or replacement semantics, but hostile conflict still mostly lands on the parent body. That leaves orbiting children visually important but conflict-light.

The next model pressure is to make visible children directly matter in conflict without replacing the existing deterministic fight system.

## Impacted Areas

### Simulation model

- same-shape conflict may now remove one attached child from the loser before full parent defeat
- child ownership, `children_count`, and visible attached-child state must remain synchronized after hostile loss

### Runtime contract

- the existing contract may remain sufficient if child absorption is expressed through changed child counts and attached-child arrays
- if the interaction result needs to distinguish child absorption from full parent defeat, one new explicit outcome may be justified

### Browser rendering

- the demo should make it legible that one orbiting child disappeared while the parent remained active

### Existing semantics

- current winner selection can remain unchanged
- replacement continuity, reproduction payment, and radius growth continue to rely on `children_count`
- the repository must preserve coherence when the same visible child bodies can now be lost through conflict as well as consumption

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether hostile child absorption happens before every parent defeat when a child exists
- whether the parent remains active immediately after that absorbed loss
- whether the existing interaction object needs a new explicit kind for absorbed child loss

## Risks If Ignored

- orbiting children remain visually expressive but mechanically underused in conflict
- the bridge between visible children and parent-level combat remains weak
- later attempts to remove transitional radius shortcuts will have less support from executable behavior

---

## Change

Allow attached orbiting children to collect food on behalf of their parent.

## Why This Matters

Attached children now affect reproduction, conflict absorption, and continuity-related consumption, but they still do not directly participate in feeding. That leaves visible orbiters absent from one of the core ecosystem loops.

The next model pressure is to make orbiting children matter in resource acquisition, not just in loss and bookkeeping.

## Impacted Areas

### Simulation model

- food-overlap detection should consider attached-child positions as authoritative collection points
- parent energy gain and food-slot removal must stay synchronized when a child collects

### Runtime contract

- the current contract may remain sufficient because attached-child positions, foods, and parent energy are already visible
- no extra event channel is required unless collection provenance becomes necessary for inspectability

### Browser rendering

- the current rendering already exposes attached-child positions, which should make child-based collection visible without new UI systems

### Existing semantics

- current radius-based collection remains active, so attached-child collection must coexist with it cleanly
- food regeneration timing should remain unchanged
- player and autonomous circles must follow the same child-based collection rule

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether attached-child overlap is treated exactly like parent-body overlap for food collection
- how to avoid double-consuming one food slot when both a parent and its child overlap it
- whether child-based collection needs any explicit interaction marker or remains visible only through ordinary snapshots

## Risks If Ignored

- orbiting children remain mechanically partial in the core energy loop
- the feeding model will stay more abstract than the conflict model
- later removal of radius shortcuts will be harder because children still will not materially affect feeding reach

---

## Change

Make parent continuity on death explicitly consume one attached child as the visible source of promotion.

## Why This Matters

The current implementation already lets a parent continue through child-based replacement, but that continuity is still mostly an abstract rule. Now that attached children are visible and mechanically meaningful, continuity should be grounded in them too.

The next model pressure is to make “a child replaces the dead parent” visibly tied to the orbiting child model rather than only to a hidden count decrement.

## Impacted Areas

### Simulation model

- death resolution paths should explicitly consume one attached child when continuity occurs
- parent generation, energy reset, and child count must stay synchronized with attached-child removal

### Runtime contract

- the current contract may remain sufficient because attached-child arrays, child count, generation, and parent presence are already visible
- no extra continuity event is required unless the outcome needs stronger inspectability

### Browser rendering

- the current rendering should already make the lost child and continued parent visible without a new UI system

### Existing semantics

- current zero-energy and defeat continuity paths can likely be reused
- replacement energy and radius reset may stay unchanged for now
- player and autonomous circles must follow the same visible-promotion rule

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether every continuity path consumes one attached child visibly
- whether the parent body remains the active representation or a child body becomes the new visible parent
- whether the current interaction model needs an explicit continuity outcome marker

## Risks If Ignored

- continuity remains more abstract than the now-visible child model
- lineage survival will feel weaker than feeding and conflict roles for children
- later removal of counter-only shortcuts will be harder because death continuity still will not be visibly grounded

---

## Change

Make attached children the single authoritative child state and derive `children_count`.

## Why This Matters

The current model now treats attached children as real visible entities for feeding, contact, payment, absorption, and continuity. But the simulation still stores and reads a separate `children_count` authority in many places. That means child semantics are still split between visible child entities and mirrored bookkeeping.

The next model pressure is to make attached children the real authoritative child state and reduce `children_count` to a derived convenience view. That removes a remaining duplication point without changing the current readable contract or child-dependent behaviors.

## Impacted Areas

### Simulation model

- child-dependent rules should derive quantity from attached children rather than a separate stored count authority
- player and autonomous circles should continue sharing one deterministic child-state rule
- current child-based behaviors should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient if `children_count` remains as a derived field
- attached-child arrays become the clearly authoritative child representation

### Browser rendering

- current rendering should remain unchanged because attached children and `children_count` can still be exposed
- no visual changes should be required in this slice

### Existing semantics

- child-based fight power, reproduction payment, and continuity can remain unchanged
- food, contact, movement, orbit, reproduction, and steering should remain unchanged
- later slices may still remove additional child-count shortcuts from rule semantics, but this slice only removes duplicated authority

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether `children_count` remains in snapshots as a derived convenience field
- which child-dependent rules should stop reading separate stored count state first
- how tests should prove that visible child entities and readable child counts can no longer drift apart

## Risks If Ignored

- visible child state and rule inputs will continue to rely on duplicated authority
- future child-model refinement will keep paying synchronization complexity
- later removal of remaining child-count shortcuts will be harder because the model will still have two competing child representations

---

## Change

Remove `children_count` from the runtime contract.

## Why This Matters

Attached children are now the single authoritative child state inside the simulation, but snapshots still expose both `attached_children` and `children_count`. That means the runtime contract still carries duplicated child representation even after the internal authority was cleaned up.

The next model pressure is to make the contract itself match the current embodied child model: child state is visible through child bodies, and counts are derived by consumers when needed.

## Impacted Areas

### Simulation model

- gameplay logic should remain unchanged
- snapshot building should stop emitting mirrored child-count fields

### Runtime contract

- `children_count` should be removed from player and autonomous circle snapshots
- contract consumers should derive child quantity from `attached_children`

### Browser rendering

- the client should derive any displayed child count from `attached_children`
- no visual behavior change should be required beyond that derivation

### Existing semantics

- fight power, reproduction payment, continuity, feeding, contact, movement, orbit, and steering should remain unchanged
- this slice removes contract duplication, not child mechanics

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether `children_count` is removed entirely from the runtime contract now that it is no longer authoritative
- how the client derives readable child quantity from attached children
- how tests should prove that the reduced contract remains inspectable and behaviorally unchanged

## Risks If Ignored

- the contract will keep carrying a mirrored child representation after internal authority has already been simplified
- client and server work will continue to maintain two ways to say the same thing
- later child-model refinements will still have to work around unnecessary contract surface

---

## Change

Reduce same-shape fight power from raw child-count magnitude to child presence only.

## Why This Matters

Attached children are now the authoritative and visible child state across geometry, contract, feeding, and continuity. But same-shape fight resolution still uses raw child-count magnitude as a direct strength ladder. That leaves combat as one of the clearest remaining places where children still act like an abstract numeric stockpile instead of primarily as visible dependents.

The next model pressure is to keep children relevant in fights while reducing that abstract stacking. A smaller embodied step is to let child presence matter, but stop letting larger raw child counts directly scale fight strength.

## Impacted Areas

### Simulation model

- same-shape fight winner selection should still prioritize energy first
- child-based fight leverage should become presence-based rather than magnitude-based
- current absorption, continuity, and payment rules should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because attached children and fight outcomes are already visible
- no extra contract field should be necessary if the changed winner rule remains inspectable from snapshots

### Browser rendering

- current rendering should remain sufficient because attached children are already visible
- no visual changes should be required in this slice

### Existing semantics

- fight absorption can still consume one attached child before full parent defeat
- reproduction payment and continuity can remain unchanged
- feeding, contact, movement, orbit, and steering should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether child-based fight advantage should now mean only “has attached child” versus “does not”
- how exact ties behave once both sides share the same child-presence state
- how tests should prove that larger child counts no longer stack direct fight power

## Risks If Ignored

- combat will remain one of the most abstract child-driven systems after most other child mechanics have been embodied
- larger raw child counts will keep acting like a hidden strength ladder even though visible child bodies already absorb loss and trigger contact
- later combat refinement will still have to unwind direct child-count stacking from the winner rule

---

## Change

Use the promoted child’s visible position for continuity.

## Why This Matters

Continuity already consumes one attached child visibly, but the continuing active parent still remains at the old parent-body position. That means the current rule is only half embodied: a child is consumed, yet the surviving line does not actually continue from that child’s place in the world.

The next model pressure is to make continuity look like child promotion rather than stationary parent persistence by reusing the promoted child’s last visible position.

## Impacted Areas

### Simulation model

- continuity should still consume one attached child
- the continuing active parent should move to the promoted child’s visible position
- lineage, generation, and replacement energy should remain unchanged

### Runtime contract

- the current snapshot shape is likely sufficient because attached-child positions and continuity outcomes are already visible
- no new contract field should be necessary if the promoted-position rule is readable from ordinary snapshots

### Browser rendering

- current rendering should remain sufficient because continuity movement will now be visible on the canvas
- no visual-style change should be required in this slice

### Existing semantics

- continuity eligibility should remain unchanged
- fight absorption, reproduction payment, feeding, contact, movement, orbit, and steering should remain unchanged
- this slice changes continuity placement, not continuity existence

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether both fight-defeat and zero-energy continuity should use the promoted child’s last visible position
- how the promoted child position is selected deterministically when more than one child exists
- how tests should prove that continuity now emerges from the child’s position instead of the parent’s old center

## Risks If Ignored

- continuity will remain less embodied than feeding, contact, and child loss
- the visible child model will still stop short of actual positional promotion
- later continuity refinement will still have to unwind the old parent-centered replacement placement

---

## Change

Expose current reproduction capacity values in interaction outcomes.

## Why This Matters

The reproduction path is now highly inspectable in terms of identity:

- blocked reproduction exposes which side or sides were blocked
- child-paid reproduction exposes which side paid through a child and which concrete child was consumed
- successful reproduction exposes created-child identity, ownership identity, and redistribution kind

But the runtime still hides the actual current capacity values that make those decisions succeed or fail. That keeps an avoidable gap between what the server already knows and what the client and tests can explicitly read.

The next model pressure is to make the current reproduction capacity of each participant explicit without changing the existing feasibility or payment rules.

## Impacted Areas

### Simulation model

- the current reproduction capacity formula should remain unchanged
- the authoritative decision point should expose the evaluated current capacity values for both source and target participants
- blocked-capacity booleans and successful outcomes should remain behaviorally unchanged

### Runtime contract

- the current snapshot shape is close but still insufficient because it exposes blockage identity without exposing the evaluated capacity values
- one minimal extension should expose source-side and target-side reproduction capacity values

### Browser rendering

- current rendering should remain sufficient if the HUD can display the new numeric values
- no visual behavior change should be required beyond that output

### Existing semantics

- reproduction threshold, cost, child reserve contribution, payment identity, creation, ownership, and redistribution should remain unchanged
- contact, fight, continuity, feeding, movement, orbit, and steering should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether reproduction capacity values are exposed on both blocked and successful reproduction outcomes
- whether the values should reflect the authoritative decision-time state before payment or after payment
- how tests should prove that the exposed values match the existing feasibility rule without redefining it

## Risks If Ignored

- blocked reproduction will remain partially inspectable through booleans only
- successful reproduction will still hide the numeric basis for why it was allowed
- future reproduction refinement will continue to rely on implicit server-side arithmetic instead of explicit outcome data

---

## Change

Expose reproduction threshold and cost constants in interaction outcomes.

## Why This Matters

The reproduction path is now inspectable through:

- blocked-capacity identity
- child-payment identity
- created-child identity
- ownership identity
- redistribution kind
- decision-time capacity values

But the governing rule constants are still implicit. The client and tests can now read the evaluated capacities, yet they still need repository knowledge to know what threshold those values were compared against and what cost the successful rule applies.

The next model pressure is to expose the threshold and cost constants explicitly without changing the reproduction rule itself.

## Impacted Areas

### Simulation model

- the current reproduction threshold and payment cost should remain unchanged
- the authoritative decision path should expose those constants as metadata
- feasibility, payment, creation, and redistribution behavior should remain unchanged

### Runtime contract

- the current snapshot shape is close but still insufficient because it exposes capacity values without the rule constants those values are measured against
- one minimal extension should expose reproduction threshold and cost values

### Browser rendering

- current rendering should remain sufficient if the HUD can display the new constants
- no visual behavior change should be required beyond that output

### Existing semantics

- reproduction capacity, blocked-capacity identity, payment identity, creation, ownership, and redistribution should remain unchanged
- contact, fight, continuity, feeding, movement, orbit, and steering should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether threshold and cost constants are exposed on both blocked and successful reproduction outcomes
- whether those constants are surfaced only when interaction kind is reproduction-related
- how tests should prove that the exposed constants match the existing server rule without redefining it

## Risks If Ignored

- capacity values will remain partially context-free in the runtime payload
- the client and tests will still need out-of-band repository knowledge to interpret reproduction outcomes fully
- future reproduction refinement will keep relying on implicit server constants instead of explicit outcome metadata

---

## Change

Expose reproduction capacity components in interaction outcomes.

## Why This Matters

The reproduction path is now inspectable through:

- blocked-capacity identity
- child-payment identity
- created-child identity
- ownership identity
- redistribution kind
- total capacity values
- threshold and cost constants

But the runtime still hides how those total capacity values are composed. The client and tests can now read the total current capacity and the governing constants, yet they still need to infer whether that total came only from energy or whether child reserve contributed to it.

The next model pressure is to expose the energy-versus-reserve split explicitly without changing the reproduction rule itself.

## Impacted Areas

### Simulation model

- the current reproduction capacity formula should remain unchanged
- the authoritative decision path should expose the energy contribution and child-reserve contribution values that compose each side's total capacity
- feasibility, payment, creation, and redistribution behavior should remain unchanged

### Runtime contract

- the current snapshot shape is close but still insufficient because it exposes total capacity without exposing how that total was formed
- one minimal extension should expose energy and reserve contribution values for both sides

### Browser rendering

- current rendering should remain sufficient if the HUD can display the new component values
- no visual behavior change should be required beyond that output

### Existing semantics

- reproduction threshold, cost, capacity totals, blocked-capacity identity, payment identity, creation, ownership, and redistribution should remain unchanged
- contact, fight, continuity, feeding, movement, orbit, and steering should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether capacity components are exposed on both blocked and successful reproduction outcomes
- how direct energy versus child reserve contribution is represented when reserve is absent
- how tests should prove that the exposed components add up to the existing total capacity values without redefining the rule

## Risks If Ignored

- total capacity values will remain only partially interpretable in the runtime payload
- the client and tests will still need implicit knowledge to know when child reserve actually mattered
- future reproduction refinement will keep relying on inferred capacity composition instead of explicit outcome metadata

---

## Change

Expose decision-time child counts in reproduction outcomes.

## Why This Matters

The reproduction path is now inspectable through:

- blocked-capacity identity
- child-payment identity
- created-child identity
- ownership identity
- redistribution kind
- total capacity values
- threshold and cost constants
- energy and reserve components of those totals

But the runtime still hides one important part of the decision context: how many attached children each side had when the reproduction decision was actually made. After successful reproduction, the post-outcome attached-child arrays already include payment and creation effects, so they no longer cleanly preserve the decision-time counts.

The next model pressure is to expose those decision-time child counts explicitly without changing the reproduction rule itself.

## Impacted Areas

### Simulation model

- the current reproduction decision path should remain unchanged
- the authoritative decision point should expose the source-side and target-side attached-child counts before payment and creation are applied
- feasibility, payment, creation, and redistribution behavior should remain unchanged

### Runtime contract

- the current snapshot shape is close but still insufficient because post-outcome attached-child arrays do not preserve the pre-decision counts directly
- one minimal extension should expose source-side and target-side decision-time child counts

### Browser rendering

- current rendering should remain sufficient if the HUD can display the new counts
- no visual behavior change should be required beyond that output

### Existing semantics

- reproduction threshold, cost, capacity totals, capacity components, blocked-capacity identity, payment identity, creation, ownership, and redistribution should remain unchanged
- contact, fight, continuity, feeding, movement, orbit, and steering should remain unchanged

## Recommended Decision Pressure

The next implementation-facing slice should explicitly choose:

- whether decision-time child counts are exposed on both blocked and successful reproduction outcomes
- how those counts should be represented when a side has no attached children
- how tests should prove that the exposed counts remain stable even when payment and creation alter the post-outcome attached-child arrays

## Risks If Ignored

- the decision context for reproduction will remain partially hidden even after the recent inspectability slices
- the client and tests will still need inference to recover pre-payment child availability
- future reproduction refinement will keep relying on implicit pre-decision child state instead of explicit outcome metadata
## Pressure: Support Panel Growth Bounds

Recent UI work improved support hierarchy, reduced legend density, reduced support text density, and stabilized the player card above the lower support row. That clarified where information belongs, but it did not yet constrain how much vertical space the lower support panels can claim as NPC summaries and recent encounter history grow.

This creates a clear presentation pressure:

- the canvas is supposed to remain the dominant play surface
- the player summary is now stable and should remain visually anchored
- lower support growth should not drag the whole support area into page-dominant status

That pressure points to a bounded client-only slice:

- keep the player card stable above the lower row
- bound NPC and encounter panel growth
- preserve readable access within those bounds
- leave server semantics and contract shape unchanged
## Pressure: Fullscreen Demo Layout

The demo started as a tightly constrained evaluation surface, which was useful while the world was small and the support area was still unstable. Recent slices reduced legend density, reduced support text density, stabilized the player card position, and bounded lower support-panel growth. That means the support area is now sufficiently controlled that the main pressure has shifted.

The new pressure is:

- the simulated world is now larger and denser
- the current centered layout leaves available viewport space unused
- the canvas no longer feels proportionate to the larger default world

That points to a bounded client-only slice:

- expand the demo across more of the viewport
- preserve the canvas as the dominant surface
- keep support information readable and secondary
- leave server semantics and contract shape unchanged
## Pressure: Side Column Visual Weight Reduction

The fullscreen layout solved the previous underuse of the available viewport by giving the canvas a larger main surface and docking the support area into a persistent side column. That improved spatial usage, but it also made the support area feel visually heavier because three bordered support blocks now read as a dense vertical mass beside the play surface.

The new pressure is:

- the canvas should remain clearly dominant
- the side column should stay secondary now that it is always visible on desktop
- support readability should be preserved without the current amount of visual chrome

That points to a bounded client-only slice:

- reduce side-column visual weight
- preserve current information hierarchy and readability
- leave server semantics and contract shape unchanged
## Pressure: Legend Collapse In Fullscreen Layout

The legend has already been reduced once, and the support column is now lighter, but the fullscreen layout changed the balance of the page. The legend still occupies a full-width band above the play surface even though the canvas, side column, and encounter log already carry most of the useful meaning.

The new pressure is:

- the legend still interrupts the fullscreen play surface more than it needs to
- the canvas should remain dominant
- the strongest cue meanings should remain recoverable without a heavy explanatory strip

That points to a bounded client-only slice:

- reduce legend footprint or prominence in the fullscreen layout
- preserve core cue recoverability
- leave server semantics and contract shape unchanged
## Pressure: Header Footprint Reduction In Fullscreen Layout

The fullscreen demo now uses more of the viewport, the side column is lighter, and the legend has been collapsed into a compact line. That leaves the remaining top-of-page pressure concentrated in the title and introductory paragraph, which still take more vertical space than they need during ordinary use.

The new pressure is:

- the canvas should begin higher in the viewport
- the page should still clearly identify the demo
- the current header footprint is larger than necessary for a fullscreen play surface

That points to a bounded client-only slice:

- reduce header footprint
- preserve clear project identity
- leave server semantics and contract shape unchanged
## Pressure: HUD Footprint Reduction In Fullscreen Layout

The fullscreen demo has already reduced the legend, header, and side-column weight. That leaves the top HUD row as the most obvious remaining high-footprint UI band above the play surface. It still occupies a full-width control strip even though the key needs are relatively small: connection/state signal, a compact identity/status summary, tick/world summary, and reset.

The new pressure is:

- the HUD still claims more space and attention than it needs
- the canvas should feel more immediate in the fullscreen layout
- core status and reset must remain easy to use

That points to a bounded client-only slice:

- reduce HUD footprint
- preserve readable status and reset access
- leave server semantics and contract shape unchanged
## Pressure: Play Stage Framing In Fullscreen Layout

The fullscreen demo has already reduced header, legend, HUD, and side-column weight. That means the canvas now has room to dominate, but it still reads mostly as a plain rectangle placed inside the page rather than as a clearly framed main stage.

The new pressure is:

- the play surface should feel more intentionally staged
- the canvas should remain dominant without heavier surrounding chrome
- support UI should stay secondary

That points to a bounded client-only slice:

- strengthen play-stage framing
- preserve the current layout and information hierarchy
- leave server semantics and contract shape unchanged
## Pressure: Side Column Internal Hierarchy In Fullscreen Layout

The fullscreen demo now has a lighter side column and a better-framed play stage, but the player, NPC, and recent encounter blocks still read with fairly similar internal weight once the eye moves into that column.

The new pressure is:

- the player panel should read as the clearest support priority
- NPC summaries should read as secondary
- recent encounters should read as tertiary
- the side column should not regain overall visual weight

That points to a bounded client-only slice:

- strengthen internal hierarchy inside the existing side column
- preserve the current information set and order
- leave server semantics and contract shape unchanged
## Pressure: Fullscreen Column Proportion Tuning

The fullscreen layout now has a tighter header, smaller HUD, collapsed legend, lighter side column, stronger play-stage framing, and clearer internal support hierarchy. The remaining presentation pressure is the fixed proportion between the stage and the support rail, which can still feel slightly rigid rather than intentionally tuned around a dominant play surface.

The new pressure is:

- the play stage should feel more dominant through width allocation
- the support rail should remain readable and usable
- the fullscreen layout should feel less rigid on desktop

That points to a bounded client-only slice:

- tune fullscreen column proportions
- preserve the current information set and hierarchy
- leave server semantics and contract shape unchanged
## Pressure: Player Follow Camera Viewport

The fullscreen demo now uses the viewport more effectively, but the browser still renders the entire authoritative world by scaling it down to fit the play stage. That keeps the whole world visible, but it weakens visual density and undercuts the sense of scale.

The new pressure is:

- the world should feel large rather than shrunk
- rendering should stay visually dense and eye-catching
- the browser needs a bounded camera model instead of whole-world scaling

That points to a bounded client-only slice:

- introduce a player-following viewport
- keep authoritative world coordinates unchanged
- clamp the camera to world bounds
- leave server semantics and contract shape unchanged
## Pressure: Camera Deadzone Follow

The viewport slice corrected the presentation model by switching from scaled whole-world overview to a bounded player-following camera. That improved scale and density, but it also exposed a new presentation pressure: the current follow rule keeps the player centered whenever possible, which can make the camera feel mechanically locked and visually busy during small movements.

The new pressure is:

- the viewport should remain clearly player-following
- the camera should feel more comfortable during small movements
- world-edge clamping should remain intact

That points to a bounded client-only slice:

- add a simple camera deadzone
- keep the viewport deterministic and easy to reason about
- leave server semantics and contract shape unchanged
## Pressure: Minimap Orientation In Viewport Mode

The viewport and deadzone slices improved scale and camera comfort, but they also removed the always-visible whole-world overview. That makes the experience stronger moment to moment, while weakening larger-world orientation.

The new pressure is:

- the viewport should remain the primary way to see the world
- the player should regain a sense of where they are in the larger space
- any orientation aid should stay small and secondary

That points to a bounded client-only slice:

- add a small minimap or equivalent orientation aid
- preserve the main viewport as the dominant rendering mode
- leave server semantics and contract shape unchanged
## Pressure: Offscreen Edge Awareness In Viewport Mode

The viewport, deadzone, and minimap slices improved scale, comfort, and large-world orientation. The remaining pressure is local awareness: meaningful entities just outside the visible camera window can disappear too abruptly, while the minimap is better at whole-world position than immediate nearby pressure.

The new pressure is:

- the viewport should remain the primary surface
- the player should gain lightweight awareness of nearby offscreen pressure
- cues should stay local and not become a general tracking overlay

That points to a bounded client-only slice:

- add lightweight offscreen edge awareness
- keep the minimap secondary
- leave server semantics and contract shape unchanged
## Pressure: Player Heading Cue In Viewport Mode

The viewport, deadzone, minimap, and offscreen-awareness slices improved scale, comfort, orientation, and nearby awareness. The remaining pressure is more local: the player still lacks a simple cue for current heading inside the viewport itself.

The new pressure is:

- the main viewport should remain primary
- the player should gain more immediate local direction readability
- the solution should stay lightweight and not become path prediction or navigation UI

That points to a bounded client-only slice:

- add a local player heading cue
- keep it tied to recent authoritative motion
- leave server semantics and contract shape unchanged
## Pressure: Offscreen Food Awareness In Viewport Mode

The viewport, minimap, offscreen circle awareness, and player heading cue now make movement and nearby pressure much easier to read. The remaining immediate gap is recovery awareness: food just outside the viewport can still remain locally invisible until it crosses into view.

The new pressure is:

- the viewport should remain the primary surface
- the player should gain lightweight awareness of nearby offscreen food
- food cues should complement, not compete with, existing offscreen circle awareness

That points to a bounded client-only slice:

- add lightweight offscreen food awareness
- keep the cues local and secondary
- leave server semantics and contract shape unchanged
## Pressure: Regional Crowding Aware Autonomy

The world now differentiates between immediate local crowding pressure and broader regional density pressure. That improves ecological consequence, but autonomous steering still only reacts directly to the more local cost model.

The new pressure is:

- autonomous circles should not keep steering as if dense regions were free to inhabit over time
- regional density should influence movement in a bounded and deterministic way
- the current food, threat, interaction, and local crowding steering layers should remain recognizable

That points to a bounded server-side slice:

- extend autonomous steering to react to regional crowding pressure
- preserve the existing steering stack rather than replacing it
- leave food, regeneration, fight, reproduction, continuity, and contract shape unchanged by default
## Pressure: Regional Food Yield Pressure

The world now has regional recovery timing, local and regional crowding cost, and regional-crowding-aware autonomous steering. Regional scarcity is therefore starting to matter structurally, but food payoff itself is still uniform once a slot is collected.

The new pressure is:

- stripped regions should not feel identical to healthier regions once food finally appears
- regional scarcity should affect both recovery timing and energy payoff
- the change should remain bounded, deterministic, and shared across player and autonomous circles

That points to a bounded server-side slice:

- make food yield region-sensitive based on nearby depletion
- preserve slot identity, slot placement, and regeneration timing
- leave fight, reproduction, continuity, autonomy rules, and contract shape unchanged by default
