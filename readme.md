# Living Circles

## Purpose

Living Circles is an experimental multiplayer ecosystem game.

The current repository explores a browser-based JavaScript canvas client connected by WebSocket to an authoritative Go server. The game model being shaped here is centered on circles as living entities constrained by energy, interacting through movement, food collection, shape-based collision meaning, and later reproduction or conflict.

At this stage, the repository is not trying to present a finished game. It is trying to make the game model explicit, testable, and evolvable in small slices.

---

## Why MRL Was Chosen

This project adopted the Model Refinement Lab because Living Circles began with incomplete but meaningful ideas rather than a fixed specification.

The early material described:

- energy as the central regulating variable
- shape as the signal that changes the meaning of collisions
- a real-time browser game connected to an authoritative backend
- agar.io as a reference point, but not as a cloning target

That kind of project benefits from a workflow that separates:

- source evidence and semantic extraction
- slice definition
- implementation
- expectation review

MRL was chosen to keep those concerns explicit in repository artifacts instead of relying on one long conversation or on undocumented implementation drift.

The working loop used here is:

```text
extract -> refine -> build -> egd -> release -> expose -> living -> extract
```

---

## How This Repository Was Initiated

This repository started from an MRL starter template and was then reshaped around the Living Circles domain.

The initial work in this repository followed this progression:

1. keep the MRL operating documents and repository structure
2. extract the original Living Circles notes into semantic artifacts and change requests
3. replace the starter's default architectural assumptions with a polyglot client/server shape
4. implement the first deterministic slices of the game model

So MRL is the method used to evolve the project, but Living Circles is the actual subject of the repository.

---

## Current State

The implemented baseline currently includes:

- a Go authoritative simulation server
- a browser client rendered with JavaScript canvas 2D
- WebSocket movement input and server snapshots
- energy consumption on movement
- deterministic food collection and energy recovery
- deterministic autonomous circles
- shape-based interaction classification
- same-shape fight resolution
- a default demo world that exposes both:
  - same-shape fight behavior
  - different-shape `reproduce_candidate` behavior

The game semantics are still intentionally incomplete. Reproduction outcomes, children, continuity, and broader ecosystem behavior remain future refinement work.

---

## Repository Use

This repository is organized to preserve project memory in artifacts:

```text
.agents/skills/            # repo-local MRL skills
/docs/operating/           # MRL model and workflow docs
/docs/semantics/           # extracted game meaning and reference knowledge
/docs/slices/              # one slice document per increment
/work/sources/             # curated raw evidence and original source material
/work/changes/             # request, impact, and implementation artifacts
/src/                      # client, server, and shared runtime code
/tests/                    # deterministic executable specification
```

Useful files to read first:

- `docs/semantics/model_hypothesis.md`
- `docs/semantics/domain_background_knowledge.md`
- `architecture.md`
- `decisions.md`
- `docs/operating/mrl_reference.md`
- `docs/operating/skills_workflow.md`

---

## Running The Current Demo

Use the Go server as the runtime entrypoint:

```bash
source "$HOME/.nvm/nvm.sh"
nvm use
go run ./src/server/cmd/livingcircles
```

Then open:

```text
http://localhost:8080
```

Contract-side JavaScript tests can be run with:

```bash
source "$HOME/.nvm/nvm.sh"
nvm use
npm run test:contracts
```

Go tests can be run with:

```bash
go test ./...
```

---

## Notes On MRL Usage Here

MRL in this repository is being used as a disciplined shaping process:

- `extract` records source evidence and semantic meaning
- `refine` defines one bounded slice
- `build` implements that slice with deterministic tests
- later phases evaluate, accept, expose, and feed new evidence back into the loop

The important point is that MRL is here to serve the evolution of Living Circles, not to overshadow it.

## Licensing

This repository uses split licensing.

- Code under `src/` and `tests/` is licensed under `MPL-2.0`.
- MRL artifacts, prompts, planning material, and documentation in `docs/` and the repository root are licensed under `MIT`, unless a file states otherwise.

See `LICENSE`, `LICENSES/MPL-2.0.txt`, and `LICENSES/MIT.txt`.
