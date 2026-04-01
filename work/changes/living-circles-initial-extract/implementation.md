# Implementation

## Slice

`docs/slices/initial_food_collection_and_energy_recovery.md`

## Implemented Shape

- Go server under `src/server/`
- browser client under `src/client/`
- explicit shared contract files under `src/shared_contracts/`
- deterministic server and integration tests under `tests/`
- authoritative food initialization and consumption inside the Go world model

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
    "x": 408,
    "y": 300,
    "radius": 12,
    "energy": 99
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
- deterministic fixed food placement
- no food respawn
- no collisions between circles, fight, reproduction, death, or continuity logic
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- player energy clamps to a maximum of `100`

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for food initialization, movement, energy drain, food consumption, and energy caps
- contract test for explicit snapshot shape including foods
- integration test for WebSocket movement, food removal, and energy recovery flow
