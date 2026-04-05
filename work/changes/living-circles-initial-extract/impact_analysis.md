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
