# Agar.io Reference Model

## 1. Purpose

This document defines the expected mechanics of a generic agar.io-like game.

It will be used as a **validation baseline** to compare Living Circles behavior.

---

## 2. Core Concept

Agar.io is a real-time multiplayer game where players control a growing cell in a shared environment.

Goal:
- grow by consuming resources and other players
- avoid being consumed

---

## 3. Core Mechanics

### 3.1 Movement

- player controls direction
- movement speed inversely proportional to size

---

### 3.2 Growth

Growth happens by:

- eating food particles
- consuming smaller players

---

### 3.3 Size Dominance

Rule:

- larger entities consume smaller ones
- smaller entities must evade

---

### 3.4 Splitting

Players can:

- split into multiple parts
- gain speed and attack capability

Trade-off:
- increased vulnerability

---

### 3.5 Merging

Split parts:

- recombine after time
- restore original mass

---

### 3.6 Mass Loss

Over time:

- large entities lose mass
- prevents infinite growth

---

### 3.7 Boundaries

- world is finite
- entities cannot leave the map

---

### 3.8 Multiplayer Interaction

- all players exist in the same world
- interactions are continuous and real-time

---

## 4. System Characteristics

### 4.1 Real-Time

- continuous updates
- low latency required

---

### 4.2 Competitive

- zero-sum interactions
- dominance hierarchy

---

### 4.3 Skill-Based

- movement control
- timing (split/merge)
- positioning

---

### 4.4 Simple Rules, Complex Outcomes

- minimal mechanics
- emergent gameplay

---

## 5. Technical Expectations

### 5.1 Client

- browser-based canvas
- smooth rendering

---

### 5.2 Server

- authoritative simulation
- handles collisions and rules

---

### 5.3 Networking

- WebSocket or similar
- frequent updates

---

## 6. Comparison Axes (for Validation)

Living Circles can be compared against:

### 6.1 Growth Model

- agar.io: consumption-based
- living circles: energy + reproduction

---

### 6.2 Interaction Model

- agar.io: dominance (bigger eats smaller)
- living circles: shape-based (fight vs reproduce)

---

### 6.3 Control Model

- agar.io: direct player control
- living circles: shared rules, partial emergence

---

### 6.4 Complexity Source

- agar.io: player skill
- living circles: system dynamics

---

### 6.5 Entity Structure

- agar.io: single entity with splits
- living circles: parent + children structure

---

## 7. Validation Goal

The goal is NOT to replicate agar.io.

The goal is to evaluate:

- familiarity
- playability
- responsiveness
- engagement

While preserving the distinct identity of Living Circles.