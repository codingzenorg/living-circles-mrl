# Living Circles

## 1. What is Living Circles

Living Circles is a 2D emergent simulation game where autonomous entities (circles) exist in a shared environment and follow a minimal life cycle:

- Move → spend energy  
- Eat → recover energy  
- Interact → fight or reproduce  
- Grow → gain advantage  
- Die → exit or transfer continuity  

The system is designed so that **complex behaviors emerge from simple interaction rules**, rather than explicit control or scripting.

It is not a character-driven game — it is a **system-driven ecosystem**.

---

## 2. Motivations

### 2.1 Emergence over Control

The primary motivation is to explore:

> How far can simple rules generate complex, lifelike dynamics?

Inspired by:
- cellular automata
- agar.io-like growth systems
- ecosystem simulations

---

### 2.2 System Thinking

The game is a **playground for system behavior**:

- feedback loops
- resource competition
- population dynamics
- dominance and collapse cycles

---

### 2.3 Player as Perturbation

The player is not a direct controller of units.

Instead, the player:
- participates as one circle
- influences the environment indirectly
- competes within the same rules as others

---

### 2.4 Platform Potential

Living Circles can evolve into:
- a multiplayer system
- a simulation sandbox
- a base for multiple game modes

---

## 3. Characteristics

### 3.1 Minimal Rule Set

The system is intentionally small:

- movement consumes energy
- food restores energy
- collisions trigger interactions

---

### 3.2 Shape-Based Interaction

Each circle has a **type/shape**:

- same shape → fight  
- different shape → reproduce  

This creates:
- tribal behavior
- implicit grouping
- emergent alliances/conflicts

---

### 3.3 Energy-Centric Life

Energy is the core variable:

- defines survival
- enables reproduction
- limits movement

---

### 3.4 Growth as Leverage

Circles can accumulate "children":

- increases area (better food collection)
- increases power in fights
- creates risk (more to lose)

---

### 3.5 Death and Continuity

When a circle dies:

- it may disappear
- or a child may replace it

This introduces:
- lineage
- continuity
- inheritance potential

---

### 3.6 No Explicit AI

Behavior emerges from:

- rules
- environment
- interactions

There is no traditional AI system.

---

## 4. Minimal Architecture

### 4.1 Client (Browser)

- Canvas rendering (2D)
- Input handling (player movement)
- WebSocket connection
- Local interpolation (optional)

---

### 4.2 Server (Go Backend)

Responsibilities:

- authoritative game state
- entity simulation loop
- collision detection
- rule resolution (fight/reproduction)
- state broadcasting

---

### 4.3 Communication

- WebSocket (real-time)
- tick-based updates

---

### 4.4 Game Loop

Server-side loop:

1. process inputs
2. update positions
3. resolve collisions
4. apply rules (fight/reproduce)
5. update energy
6. broadcast state

---

### 4.5 State Model (Minimal)

Entities:

- Circle
- Food
- World

---

## 5. Key Design Aspects

### 5.1 Emergence First

Avoid:
- scripted behaviors
- complex AI layers

Prefer:
- simple rules + interaction density

---

### 5.2 Deterministic Core (Optional)

Consider:

- deterministic simulation
- reproducible states

Useful for:
- debugging
- replay systems

---

### 5.3 Fairness via Rules

All entities (including player):

- follow the same rules
- no hidden advantages

---

### 5.4 Scale Through Simplicity

System must support:

- many entities
- frequent interactions

Thus:
- rules must be cheap
- data structures must be simple

---

### 5.5 Energy as Unifying Metric

Energy connects:

- movement
- survival
- combat
- reproduction

This keeps the system coherent.

---

### 5.6 Replace Control with Interaction

Instead of:
- controlling outcomes

Design:
- conditions for outcomes

---

## 6. MRL Extraction Notes

This document should allow extraction of:

- entities (Circle, Food, World)
- attributes (energy, size, shape, children)
- events (collision, reproduction, fight, death)
- rules (interaction matrix)
- invariants (energy > 0, size > 0)

The system is intentionally minimal to allow iterative refinement.