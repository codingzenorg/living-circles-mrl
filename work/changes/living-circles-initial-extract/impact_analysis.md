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
