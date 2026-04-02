# Implementation

## Slice

`docs/slices/initial_energy_collapse_death.md`

## Implemented Shape

- Go server under `src/server/`
- browser client under `src/client/`
- explicit shared contract files under `src/shared_contracts/`
- deterministic server and integration tests under `tests/`
- authoritative food initialization and consumption inside the Go world model
- two deterministic autonomous circles participating under the same movement and energy rules in the default demo world
- explicit shape identity and current interaction classification in snapshots
- deterministic same-shape fight resolution with loser removal
- default demo visibility for both same-shape and different-shape interaction paths
- deterministic different-shape reproduction resolution with child accumulation counts
- deterministic radius growth derived from child accumulation
- deterministic child replacement on defeat when the loser has available children
- zero-energy collapse now causes death or replacement continuity
- browser demo reset through an authoritative server restart endpoint

## Runtime Contract

### Client to server

`movement_intent`

```json
{
  "type": "movement_intent",
  "direction": {
    "x": 1,
    "y": 0
  }
}
```

### Server to client

`world_snapshot`

```json
{
  "type": "world_snapshot",
  "tick": 1,
  "world": {
    "width": 800,
    "height": 600
  },
  "player": {
    "id": "player-1",
    "shape": "triangle",
    "x": 408,
    "y": 300,
    "radius": 16,
    "energy": 99,
    "children_count": 1
  },
  "autonomous_circles": [
    {
      "id": "circle-2",
      "shape": "triangle",
      "x": 268,
      "y": 300,
      "radius": 12,
      "energy": 100,
      "children_count": 0
    },
    {
      "id": "circle-3",
      "shape": "square",
      "x": 532,
      "y": 300,
      "radius": 16,
      "energy": 99,
      "children_count": 1
    }
  ],
  "interaction": {
    "active": false,
    "resolved": true,
    "kind": "fight_resolved",
    "source_id": "player-1",
    "target_id": "circle-2",
    "winner_id": "player-1",
    "loser_id": "circle-2"
  },
  "foods": [
    {
      "id": "food-1",
      "x": 432,
      "y": 300,
      "radius": 6
    }
  ]
}
```

## Deliberate Simplifications

- one player-controlled circle only
- exactly two autonomous circles in the default demo world
- deterministic autonomous movement policy: first autonomous circle moves right, second moves left
- deterministic shape assignment in the default demo world: player `triangle`, same-shape autonomous `triangle`, different-shape autonomous `square`
- deterministic fixed food placement
- no food respawn
- child accumulation is represented as an integer count on each circle
- different-shape reproduction resolves without spawning separate child entities
- radius is derived from child accumulation with a fixed per-child increment
- continuity is limited to one-child replacement after fight defeat
- zero energy is now a death threshold rather than only a movement stop condition
- demo reset recreates the initial world state without restarting the server process
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- player energy clamps to a maximum of `100`
- autonomous circles follow deterministic index-based directions so the default demo exposes both interaction paths
- same-shape overlap resolves as `fight_resolved`; different-shape overlap resolves as `reproduce_resolved`
- resolved reproduction awards exactly one child accumulation unit to each participating circle
- each accumulated child increases radius by a fixed deterministic amount of `4`
- a circle pair may reproduce at most once while continuously overlapping and must separate before reproducing again
- same-shape fights resolve in one tick using: higher energy wins, then larger radius, then player wins exact ties
- a fight loser with at least one child remains active through immediate replacement, consuming exactly one child
- replacement stays at the defeated circle position and resets to deterministic baseline energy `100`
- a zero-energy circle follows the same removal-or-replacement rule used after fight defeat
- `POST /reset` rebuilds the session from its initial config, resets tick to `0`, clears intent, and broadcasts the fresh snapshot

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for child accumulation, derived radius growth, fight winner selection, loser removal, child replacement continuity, and zero-energy collapse death
- contract test for explicit snapshot shape including child counts and resolved reproduction outcomes
- integration tests for WebSocket snapshots with visible growth, resolved reproduction, fight continuity through child replacement, zero-energy collapse continuity, and authoritative reset
