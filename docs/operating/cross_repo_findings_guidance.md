# Cross-Repository Findings Handling (MRL-Starter)

## Purpose

Define how a repository should handle **findings discovered locally** that impact other repositories.

This avoids:
- intrusive changes across repositories
- ownership confusion
- hidden coupling
- incorrect “implementation from the outside”

---

## Core Principle

> A repository can observe and report a problem, but only the owning repository should define and implement the fix.

---

## What is a Finding?

A finding is:

- a mismatch in contract (API, event, payload, env var)
- an unexpected behavior from another repository
- an integration failure
- a runtime inconsistency
- a deploy/runtime side-effect

---

## Decision: Skill or Guidance?

This should be a **guidance**, not a skill.

Why:

- It is **behavioral and architectural**, not procedural
- It applies broadly, not as a repeatable execution script
- It influences how agents think, not just what they do

However, you MAY add a lightweight skill later for:
- extracting payload examples
- generating comparison (expected vs actual)

---

## Standard Flow

### Step 1 — Create a Finding (Local Repository)

Create:

```
work/findings/<date>-<short-name>/finding.md
```

Example:

```
work/findings/2026-04-29-candidate-image-dispatch/finding.md
```

---

## Finding Document Template

```md
# Finding: <short name>

## Context

Where this was discovered:
- feature
- deploy
- runtime
- test

## Observed Behavior

What actually happened.

Include:
- payloads
- logs
- traces
- screenshots (if relevant)

## Expected Behavior

What should have happened.

Reference:
- contracts
- docs
- assumptions

## Impact

What breaks because of this:
- deploy
- runtime
- user flow
- data integrity

## Suspected Source

Which repository likely owns the issue.

## Evidence

Concrete data:
- logs
- request/response samples
- event payloads

## Suggested Direction (Optional)

High-level direction only.

Do NOT:
- define internal implementation of another repository
- prescribe exact code changes

## Owning Repository

The repository responsible for fixing the issue.

## Local Impact

What this repository must do (if anything):
- adapt temporarily
- block deploy
- add validation
```

---

## Step 2 — Escalate to Management (If Needed)

If the finding affects more than one repository:

- create or update a campaign in `management`
- attach this finding as evidence

Do NOT coordinate multi-repo changes locally.

---

## Step 3 — Create Repository Handoff

For the owning repository:

- create a Change Packet in `management`
- or communicate via direct handoff

The owning repository must:
- interpret the finding
- decide implementation
- update its own knowledge

---

## Step 4 — Local Repository Behavior

The discovering repository may:

- add validation to prevent silent failures
- log better diagnostics
- document the expected contract

The discovering repository must NOT:

- implement fixes for another repository
- assume ownership of foreign behavior
- embed permanent workarounds without explicit decision

---

## Anti-Patterns

### 1. Writing implementation.md for another repository

Bad:

```
infra-platform/work/changes/.../implementation.md
```

Why:

- breaks ownership boundaries
- creates coupling
- makes future maintenance unclear

---

### 2. Silent workaround

Bad:

```
if payload.invalid:
    fix_it_here()
```

Why:

- hides real issue
- spreads incorrect logic

---

### 3. Expanding scope locally

Bad:

```
fix infra + fix app + redefine contract
```

Why:

- creates implicit cross-repo changes
- breaks coordination

---

## Recommended Practices

- keep findings small and precise
- include real evidence
- identify ownership clearly
- escalate when contracts are unclear
- avoid guessing other repository internals
- prefer explicit contracts over assumptions

---

## Mental Model

```
repository = observer + executor of its own scope
management = coordinator
owning repository = problem solver
```

---

## Final Rule

A finding is complete when:

1. it clearly describes the issue
2. ownership is identified
3. evidence is provided
4. it does not prescribe implementation details of another repository
5. it can be used to create a proper cross-repo change

The fix must always live in the owning repository.
