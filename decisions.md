# Decisions Log

## Purpose

This document records relevant architectural and implementation decisions for the adopting repository.

Use it to preserve reasoning, avoid re-discussing settled trade-offs without context, and document deviations from `architecture.md` and `groundrules.md`.

---

## Entry Template

```md
## DEC-XXXX - Title

- Date: YYYY-MM-DD
- Status: proposed | accepted | superseded | rejected
- Owners: human | codex | both

### Context
What problem or tension led to this decision?

### Decision
What was decided?

### Consequences
What becomes easier, harder, or different because of this?

### Alternatives considered
What other options were considered and why were they not chosen?

### Notes
Any additional implementation guidance, migration note, or follow-up.
```

---

## Index

Add entries as the repository evolves.

## DEC-0001 - Separate MRL Core From Implementation Packs

- Date: 2026-03-29
- Status: accepted
- Owners: both

### Context
The starter was presenting Python plus a DDD-inspired modular monolith as if that were the default shape of MRL itself. That creates confusion when a repository needs another language, another architecture such as event sourcing, or more than one runtime.

### Decision
The repository now distinguishes between:

- MRL core, which stays artifact-driven and architecture-agnostic
- implementation packs, which define language, architecture, structure, and testing defaults

The current repository keeps `python_ddd_monolith` as the example selected pack.

### Consequences
It becomes easier to reuse the same refinement workflow across Python, JavaScript, Go, event-sourced, and polyglot client/server repositories. It also becomes necessary to make the selected pack explicit in architecture docs and slice docs.

### Alternatives considered
Keep one universal Python starter and treat every other shape as an undocumented deviation. This was rejected because it would keep conflating MRL with one implementation style.

### Notes
Future pack additions should live under `docs/packs/` and should be referenced by slice documents when the runtime topology matters.

## DEC-0002 - Adopt Polyglot Client Server For Living Circles

- Date: 2026-04-01
- Status: accepted
- Owners: both

### Context
The initial extraction and first slice definition established that Living Circles is intended as:

- a JavaScript browser game
- rendered with 2D canvas
- connected over WebSocket
- to an authoritative Go simulation server

Keeping `python_ddd_monolith` as the selected pack would leave the repository architecture misaligned with the extracted product shape before the first executable slice.

### Decision
This repository now adopts `polyglot_client_server` as its selected pack.

The repository should therefore:

- treat the browser client and Go server as first-class runtimes
- keep semantics shared under `docs/semantics/`
- make runtime contracts explicit under a shared contract boundary
- design early slices around server authority and client/server interaction rather than a fake single-runtime abstraction

### Consequences
It becomes easier to refine and build slices that match the intended game shape from the start. It also means:

- `architecture.md` must describe a client/server layout rather than a Python monolith
- `src/` and `tests/` should evolve toward client, server, integration, and contract structure
- deterministic behavior must be preserved across a runtime boundary, not only inside one process

### Alternatives considered
Keep `python_ddd_monolith` temporarily and treat the extracted client/server topology as only an implementation detail to revisit later. This was rejected because it would encourage the first code to be scaffolded against the wrong authority model and wrong layout.

### Notes
This decision does not force full multiplayer scope or complex network realism in the first slice. Early slices should still remain small, deterministic, and behavior-focused.
