# Implementation

## Slice

`docs/slices/initial_authoritative_movement_loop.md`

## Implemented Shape

- Go server under `src/server/`
- browser client under `src/client/`
- explicit shared contract files under `src/shared_contracts/`
- deterministic server and integration tests under `tests/`

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
  }
}
```

## Deliberate Simplifications

- one player-controlled circle only
- no food, collisions, fight, reproduction, death, or continuity logic
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed one implementation choice not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero

This keeps the loop deterministic and prevents negative energy while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for movement, energy drain, idle ticks, and bounds
- contract test for explicit message shape
- integration test for WebSocket movement and snapshot flow
