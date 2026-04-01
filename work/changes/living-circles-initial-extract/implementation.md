# Implementation

## Slice

`docs/slices/initial_shape_identity_and_collision_classification.md`

## Implemented Shape

- Go server under `src/server/`
- browser client under `src/client/`
- explicit shared contract files under `src/shared_contracts/`
- deterministic server and integration tests under `tests/`
- authoritative food initialization and consumption inside the Go world model
- one deterministic autonomous circle participating under the same movement and energy rules
- explicit shape identity and current interaction classification in snapshots

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
    "radius": 12,
    "energy": 99
  },
  "autonomous_circles": [
    {
      "id": "circle-2",
      "shape": "square",
      "x": 268,
      "y": 300,
      "radius": 12,
      "energy": 99
    }
  ],
  "interaction": {
    "active": true,
    "kind": "reproduce_candidate",
    "source_id": "player-1",
    "target_id": "circle-2"
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
- one autonomous circle only
- deterministic autonomous movement policy moving right on every tick
- deterministic shape assignment: player `triangle`, autonomous circle `square`
- deterministic fixed food placement
- no food respawn
- no fight or reproduction resolution, death, or continuity logic
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- player energy clamps to a maximum of `100`
- the autonomous circle follows a fixed deterministic direction: right on every tick
- same-shape overlap maps to `fight_candidate`; different-shape overlap maps to `reproduce_candidate`

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for shape assignment, overlap classification, shared food consumption, and energy caps
- contract test for explicit snapshot shape including shape identity and interaction classification
- integration test for WebSocket snapshots with active interaction classification
