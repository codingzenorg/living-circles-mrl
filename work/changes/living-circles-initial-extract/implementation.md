# Implementation

## Slice

`docs/slices/demo_dual_interaction_visibility.md`

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
      "shape": "triangle",
      "x": 268,
      "y": 300,
      "radius": 12,
      "energy": 99
    },
    {
      "id": "circle-3",
      "shape": "square",
      "x": 532,
      "y": 300,
      "radius": 12,
      "energy": 99
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
- different-shape reproduction remains classification-only
- no continuity, replacement, or child logic after defeat
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- player energy clamps to a maximum of `100`
- autonomous circles follow deterministic index-based directions so the default demo exposes both interaction paths
- same-shape overlap resolves as `fight_resolved`; different-shape overlap maps to `reproduce_candidate`
- same-shape fights resolve in one tick using: higher energy wins, then larger radius, then player wins exact ties

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for default dual-circle composition, fight winner selection, loser removal, and unresolved reproduction classification
- contract test for explicit snapshot shape including resolved fight outcomes
- integration tests for WebSocket snapshots with default dual-circle visibility and same-shape fight resolution
