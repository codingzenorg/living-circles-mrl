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
