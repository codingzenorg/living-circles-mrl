# MRL Starter

## Purpose

This document defines the minimal portable shape for adopting the Model Refinement Lab in a new repository.

Use it when creating a fresh project that should follow the same refinement discipline without copying domain-specific content from this repository. The practical instantiation flow is to clone the starter into a new repository name and then let `extract` create the first semantic baseline.

---

## Split the Repository Into Two Parts

### MRL core

These parts should be portable across domains:

- `readme.md`
- `groundrules.md`
- `decisions.md`
- `AGENTS.md`
- `.agents/skills/`
- `docs/operating/`
- evaluator pattern
- semantic docs
- slice docs

### Pack-specific

These parts define language and architecture defaults and may be replaced:

- `architecture.md`
- `docs/building/project_structure.md`
- `docs/packs/`
- scenario runner pattern
- test shape
- code under `src/`
- tests under `tests/`
- language toolchain files such as `pyproject.toml`, `package.json`, or `go.mod`

### Project-specific

These parts should be rewritten for each new domain:

- `docs/semantics/model_hypothesis.md`
- `docs/semantics/domain_background_knowledge.md`
- `docs/slices/`
- local evaluation artifacts under `runs/`

---

## Minimal Starter Layout

```text
project_root/
  readme.md
  architecture.md
  groundrules.md
  decisions.md
  AGENTS.md
  .gitignore
  pyproject.toml

  .agents/
    skills/
      extract/
      refine/
      build/
      egd/
      release/
      expose/
      living/

  docs/
    operating/
      mrl_reference.md
      mrl_starter.md
      packs.md
      skills_workflow.md
    building/
      project_structure.md
    packs/
      python_ddd_monolith.md
    evaluation/
      scenario_evaluation.md
    semantics/
      model_hypothesis.md
      domain_background_knowledge.md
    slices/

  src/                        # shaped by the selected pack
  tests/                      # shaped by the selected pack
```

---

## First Bootstrap Sequence

A new repository should start with this order:

1. define root strategic docs and `docs/operating/`
2. install repo-local skills under `.agents/skills/`
3. select or define a pack and record it in `decisions.md`
4. keep the semantic docs as placeholders under `docs/semantics/`
5. run `extract` to build the first model baseline
6. run `refine` to create the first slice
7. run `build` to implement the first executable vertical slice
8. run `egd` before deciding `release`

---

## Portability Test

MRL is portable only if a new domain can reuse the same loop and skills by changing semantic artifacts and slices rather than rewriting the workflow itself.

A practical test is:

- keep the operating docs and skills
- select or replace the implementation pack intentionally
- replace semantic docs and slice docs
- implement a new first slice
- run the same loop end to end

If that works, the starter is sufficiently portable.
