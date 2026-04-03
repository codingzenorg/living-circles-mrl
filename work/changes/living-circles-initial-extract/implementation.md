# Implementation

## Slice

`docs/slices/initial_attached_child_contact_triggers_parent_interaction.md`

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
- deterministic food-slot regeneration after consumption
- lineage identity and generation are now explicit in active circles and preserved through replacement continuity
- reproduction is now gated by participant energy and consumes energy on success
- autonomous circles now steer deterministically toward nearby food before falling back to baseline drift
- child accumulation now contributes directly to same-shape fight winner selection
- successful reproduction now creates visible attached orbiting children owned by the participating parents
- same-shape conflict can now remove one attached child from the loser before removing the parent
- attached orbiting children can now collect food on behalf of their parent
- parent continuity after death now emits an explicit `death_promoted_child` outcome when one attached child is consumed to preserve the lineage
- attached orbiting children can now trigger parent-level fight or reproduction before parent core bodies overlap
- interaction snapshots now expose `contact_origin` so child-triggered contact is inspectable
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
    "lineage_id": "lineage-player-1",
    "generation": 1,
    "shape": "triangle",
    "x": 408,
    "y": 300,
    "radius": 16,
    "energy": 99,
    "children_count": 1,
    "attached_children": [
      {
        "id": "child-2",
        "owner_id": "player-1",
        "orbit_slot": 0,
        "x": 429.2,
        "y": 304.8,
        "radius": 4
      }
    ]
  },
  "autonomous_circles": [
    {
      "id": "circle-2",
      "lineage_id": "lineage-circle-2",
      "generation": 1,
      "shape": "triangle",
      "x": 268,
      "y": 300,
      "radius": 12,
      "energy": 100,
      "children_count": 0,
      "attached_children": []
    },
    {
      "id": "circle-3",
      "lineage_id": "lineage-circle-3",
      "generation": 0,
      "shape": "square",
      "x": 532,
      "y": 300,
      "radius": 16,
      "energy": 99,
      "children_count": 1,
      "attached_children": [
        {
          "id": "child-1",
          "owner_id": "circle-3",
          "orbit_slot": 0,
          "x": 553.2,
          "y": 304.8,
          "radius": 4
        }
      ]
    }
  ],
  "interaction": {
    "active": false,
    "resolved": true,
    "kind": "death_promoted_child",
    "source_id": "player-1",
    "target_id": "player-1"
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
- autonomous circles prefer the nearest active food target and fall back to deterministic baseline drift when no food target is available
- deterministic shape assignment in the default demo world: player `triangle`, same-shape autonomous `triangle`, different-shape autonomous `square`
- deterministic fixed food placement
- deterministic food regeneration returns consumed slots to their original positions after a fixed delay
- child accumulation remains the source of truth for current radius, fight power, replacement continuity, and child-payment rules, but is now also embodied as attached orbiting children
- attached-child positions now extend a parent's effective food collection reach
- different-shape reproduction resolves through deterministic child redistribution across the participating parents when both participants satisfy the energy rule
- radius is derived from child accumulation with a fixed per-child increment
- continuity is limited to one-child replacement after fight defeat
- lineage is represented only by a stable `lineage_id` plus monotonic `generation`
- zero energy is now a death threshold rather than only a movement stop condition
- demo reset recreates the initial world state without restarting the server process
- no local prediction or interpolation
- one shared movement intent for the connected client
- static world size and player radius

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- each consumed food slot regenerates exactly `12` ticks after consumption
- attached-child food collection gives the same parent energy gain as parent-body food collection
- attached-child body contact against another parent body now counts as authoritative pair contact for interaction triggering
- interaction provenance now records whether contact came from `parent_body` or `attached_child`
- player energy clamps to a maximum of `100`
- autonomous circles deterministically select the nearest active food by distance, breaking ties by food ID
- when no food target is available, autonomous circles fall back to deterministic index-based directions
- same-shape overlap resolves as `fight_resolved`; different-shape overlap resolves as `reproduce_resolved`
- same-shape overlap may also resolve as `fight_absorbed_child` when the loser has an attached child available to absorb the loss
- zero-energy collapse that preserves continuity now resolves as `death_promoted_child`
- pair contact may be initiated by one attached child touching the other parent's body even when the two parent cores do not yet overlap
- same-shape fight ordering is: higher energy, then higher child count, then larger radius, then player exact tie-break
- a same-shape loser with at least one attached child loses exactly one attached child and remains active on that tick
- different-shape overlap without enough energy or an available child payment resolves as `reproduce_blocked_energy`
- resolved reproduction creates exactly two new child units and assigns them deterministically across the participating parents
- attached children remain orbiting dependents of their current parent rather than immediately becoming free autonomous circles
- attached child orbit positions are derived deterministically from owner, child identity, slot, and tick
- reproduction requires a total reproduction capacity of at least `15`
- reproduction costs `10` from each participant on success
- one child may contribute one `10`-point reserve unit toward the threshold check and payment path
- a circle below the reproduction cost may consume exactly one child, convert it into a `10`-point temporary reserve, pay the reproduction cost, and then receive the reproduction result
- each accumulated child increases radius by a fixed deterministic amount of `4`
- a circle pair may reproduce at most once while continuously overlapping and must separate before reproducing again
- same-shape fights resolve in one tick using: higher energy wins, then larger radius, then player wins exact ties
- a fight loser with at least one child remains active through immediate replacement, consuming exactly one child
- replacement stays at the defeated circle position and resets to deterministic baseline energy `100`
- initial circles derive deterministic lineage IDs from their active IDs and start at generation `0`
- replacement continuity preserves `lineage_id` and increments `generation` by exactly `1`
- a zero-energy circle follows the same removal-or-replacement rule used after fight defeat
- when continuity occurs on parent death, one attached child is explicitly consumed and the snapshot exposes `death_promoted_child` for inspectability unless a higher-priority same-tick interaction already exists
- parent-body and attached-child contact still resolve at most one interaction per pair in a tick
- `POST /reset` rebuilds the session from its initial config, resets tick to `0`, clears intent, and broadcasts the fresh snapshot

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for child accumulation, attached orbiting child ownership, attached-child food collection, child absorption during hostile conflict, derived radius growth, fight winner selection including child power, loser removal, child replacement continuity on zero-energy collapse, lineage preservation, food-slot regeneration timing, energy-gated reproduction, and food-seeking autonomy
- contract test for explicit snapshot shape including child counts, attached child state, lineage fields, and resolved, blocked, absorbed, or promoted interaction outcomes
- integration tests for WebSocket snapshots with visible orbiting children after reproduction, attached-child food collection, attached-child loss during hostile conflict, child-triggered contact before parent core overlap, blocked reproduction by low energy, food-seeking autonomous motion, child-driven fight outcomes, zero-energy collapse continuity with explicit promotion signaling, deterministic food regeneration, and authoritative reset
