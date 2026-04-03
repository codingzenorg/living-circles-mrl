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

- whether visible orbiting children are a one-to-one embodiment of `children_count`
- how post-reproduction child ownership is assigned while remaining deterministic
- whether current radius growth stays temporarily active alongside orbiting children

## Risks If Ignored

- reproduction output will keep drifting away from the intended game feel
- later changes may require replacing tests and client assumptions built around free spawned child circles
- child-related semantics will remain split between counters and bodies without an explicit bridge

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
