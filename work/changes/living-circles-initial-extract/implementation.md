# Implementation

## Slice

`docs/slices/initial_idle_observer_orientation_only_transport.md`

## Implemented Shape

- Go server under `src/server/`
- browser client under `src/client/`
- explicit shared contract files under `src/shared_contracts/`
- deterministic server and integration tests under `tests/`
- the support area now has a clearer hierarchy: the player card is primary in the left column, while NPC summaries and recent encounters are grouped into a tighter secondary stack on the right
- the player support card now uses one compact identity row plus a short metric grid instead of six separate labeled tiles
- NPC summaries now use tighter single-line text rows with lighter formatting so the support area scans faster during play
- the player card now spans the full support width, with NPC summaries and recent encounters sharing a lower two-panel row so growing lists do not push the player summary around
- the lower NPC and recent-encounter panels now have bounded vertical growth with internal scrolling so content accumulation does not keep expanding the support area
- the demo now uses a fullscreen desktop layout with the canvas taking the main viewport area and the support column docked to the side while preserving a narrow-screen collapse
- the side support column now uses lighter panel chrome, tighter spacing, and softer internal surfaces so it reads more clearly as secondary to the canvas
- the fullscreen legend now uses a single compact explanatory line instead of a full chip band so it interrupts the play surface less
- the fullscreen header now uses a tighter title and shorter intro copy so the play surface begins higher without losing the page identity
- the fullscreen HUD now uses smaller chips, a tighter control row, and shorter status text so the play surface feels more immediate while reset stays visible
- the fullscreen play surface now sits inside a lighter framed stage container so the canvas reads more clearly as the main staged area within the page
- the fullscreen side column now gives the player panel stronger emphasis while NPCs and recent encounters use progressively lighter treatment, clarifying the internal support hierarchy
- the fullscreen desktop split now gives the play stage a slightly larger share while keeping the support rail readable and secondary
- the browser now renders a bounded player-following viewport into the authoritative world instead of shrinking the full world to fit the stage, with the camera clamped at world bounds
- the player-follow camera now uses a centered deadzone so small movements do not constantly recenter the viewport while world-edge clamping remains intact
- viewport mode now includes a small passive minimap overlay that shows world position and the current viewport window without competing with the main play surface
- viewport mode now adds lightweight edge cues for nearby offscreen circles so local danger or opportunity just outside the camera window is easier to notice
- viewport mode now adds a small local heading cue near the player derived from recent authoritative motion so movement direction is easier to read inside the camera window
- viewport mode now adds lightweight edge cues for nearby offscreen food so local recovery opportunity just outside the camera window is easier to notice
- viewport mode now applies a small camera lookahead in the player's recent authoritative direction of travel while preserving deadzone follow and world-edge clamping
- transport payload measurement is now explicit through a server-side helper that serializes the current authoritative snapshot as-is and combines it with the live tick cadence
- the deterministic default expanded-world baseline currently measures at `6487` bytes per snapshot, or about `64870` bytes/sec per client at the current `10` snapshots/sec cadence
- a deterministic larger-world scenario with doubled expanded autonomous and food counts currently measures at `12539` bytes per snapshot, or about `125390` bytes/sec per client at the same cadence
- websocket and reset transport now send local full-detail autonomous circles and foods for the current viewport neighborhood instead of the full world at equal detail
- transport snapshots now also carry lightweight whole-world minimap summaries and total world counts so the client keeps orientation while the main play surface consumes less payload
- the viewport-thinned transport path now measures at `3907` bytes per snapshot, or about `39070` bytes/sec per client at `10` snapshots/sec, down from the prior `6487` byte / `64870` bytes-sec full-snapshot baseline
- transport snapshots now keep local viewport detail on every tick while refreshing minimap-orientation summaries only every `5` ticks, with the client reusing the last valid orientation summary between refreshes
- the current single-cadence viewport-culling payload measures at `3932` bytes per snapshot, or about `39320` bytes/sec per client at `10` snapshots/sec
- the current dual-cadence transport averages about `1710` bytes per message, or about `17096` bytes/sec per client over the same `10` snapshots/sec cadence window
- minimap-orientation refreshes now use deterministic coarse clusters instead of exact whole-world food and autonomous point lists
- the exact orientation-refresh payload would measure at `3574` bytes per refresh, while the compact clustered orientation refresh now measures at `3311` bytes
- the compact dual-cadence transport now averages about `1585` bytes per message, or about `15854` bytes/sec per client over the same `10` snapshots/sec cadence window
- local transport detail is now serialized with reduced display-sufficient precision at the transport boundary while the server keeps full internal simulation precision
- the compact full-precision local transport payload would measure at `3311` bytes per snapshot, while the reduced-precision local payload now measures at `3275` bytes
- the reduced-precision dual-cadence transport now averages about `1549` bytes per message, or about `15494` bytes/sec per client over the same `10` snapshots/sec cadence window
- the default expanded world baseline now uses a larger bounded space, more autonomous circles, and deterministic seeded food slots instead of a hand-authored expanded food layout
- the expanded default autonomous startup pattern now uses deterministic seeded placement for the additional expanded circles instead of fixed authored offsets
- the additional expanded autonomous circles now also use a deterministic seeded startup shape and energy mix instead of fixed authored per-ID state
- food regeneration pressure now includes a deterministic local component, so slots in more heavily depleted neighborhoods return later than equally old isolated missing slots
- crowded regions now add a small deterministic energy penalty beyond the existing local crowding surcharge, making dense areas more expensive to inhabit over time
- the legend is now reduced to the strongest cue families so the canvas and external panels carry more of the explanatory load during ordinary play
- browser rendering now keeps a very short-lived client-local afterglow for recent resolved fights, reproductions, and continuity-preserving promotion so important outcomes remain briefly visible in-world
- recent event afterglow remains grounded in current authoritative interaction outcomes and current entity positions without introducing server-side event history
- browser rendering now emphasizes continuity-bearing lineage through visible parent-child connection links, reserve halos around circles with attached children, and stronger promotion emphasis when continuity survives through a child
- lineage legibility remains grounded in existing authoritative lineage, generation, attached-child, and continuity outcome state without adding new server fields
- browser rendering now uses short-range motion cues to distinguish nearby autonomous food pursuit, social pursuit, and retreat without adding any server-side intent fields
- autonomy legibility is inferred only from consecutive authoritative positions plus current visible world geometry, not from hidden client rules or server-published intent
- browser rendering now highlights nearby food-rich space with local green recovery cues and marks local scarcity around the player with a cool blue ring, grounded only in current visible food positions
- the player HUD now surfaces nearby food opportunity or scarcity without introducing any server-side food legibility fields
- browser rendering now highlights costly nearby crowding zones with local amber halos and circle-level crowding rings grounded only in visible circle positions
- the player HUD now surfaces immediate crowding pressure without adding any new server-side crowding fields
- browser rendering now uses stronger on-canvas cues for same-shape danger, different-shape opportunity, and blocked reproduction near the player instead of relying mainly on dense HUD text
- the live HUD now summarizes interaction meaning in more readable language while preserving detailed authoritative identifiers for inspection
- authoritative food initialization and consumption inside the Go world model
- expanded default-world food capacity now derives from the initial active-circle count rather than remaining a purely hand-picked slot count
- food-slot regeneration delay now scales deterministically with current missing-slot count rather than remaining one unconditional global timer
- circles now pay a small additional energy cost when at least two other active circles are within a fixed local crowding radius
- autonomous steering now avoids moving deeper into a clearly more crowded neighborhood by reversing the preferred direction when the next step would add more than one additional nearby circle
- expanded default world baseline with a larger bounded map, five deterministic autonomous circles, and eight deterministic food slots
- explicit shape identity and current interaction classification in snapshots
- deterministic same-shape fight resolution with loser removal
- default demo visibility for both same-shape and different-shape interaction paths
- deterministic different-shape reproduction resolution with child accumulation counts
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
- attached orbiting children can now also trigger parent-level interaction through child-to-child overlap
- autonomous circles can now resolve fight and reproduction outcomes against each other without player involvement
- autonomous circles can now seek other autonomous circles when no nearby food has priority
- autonomous circles can now also target the player through the same bounded interaction-seeking rule
- low-energy autonomous circles now prefer food recovery before interaction seeking
- interaction-seeking autonomous circles now prefer different-shape targets before same-shape targets
- interaction-seeking autonomous circles now prefer different-shape targets only when current reproduction is feasible for the pair
- interaction-seeking autonomous circles now treat same-shape fallback targets as eligible only when the current fight ordering says they would not lose
- autonomous circles can now retreat from nearby stronger same-shape threats before ordinary pursuit resumes
- autonomous circles can now retreat from nearby different-shape targets when reproduction is currently blocked
- attached-child positions can now trigger same-shape threat and blocked-reproduction avoidance before parent-core overlap
- attached-child positions can now also influence which eligible social target is effectively nearest during positive interaction seeking
- attached-child positions can now also influence which food slot is effectively nearest during food targeting
- child-paid reproduction is now explicitly distinguished from ordinary energy-paid reproduction in resolved interaction outcomes
- food collection now depends on visible parent-core and attached-child overlap rather than enlarged derived radius alone
- circle-to-circle contact now depends on visible parent-core and attached-child overlap rather than enlarged derived parent radius alone
- same-shape fights no longer use derived radius as a tie-break after energy and child count
- parent movement boundary clamping now uses the visible parent-core body rather than grown derived radius
- attached-child orbit distance now uses the visible parent-core body rather than grown derived radius
- parent `radius` in snapshots and rendering now stays fixed at the visible parent-core body rather than growing from child count
- attached children are now the single authoritative child state, with `children_count` derived for snapshots and readability
- runtime snapshots now expose child state only through attached children, with client-side child quantity derived from `attached_children`
- same-shape fights now treat child-based fight leverage as attached-child presence versus absence rather than raw child-count magnitude
- continuity promotion now repositions the continuing active parent to the promoted child's last visible position
- continuity outcomes now explicitly expose the promoted child identity
- `fight_absorbed_child` outcomes now explicitly expose the absorbed child identity
- `reproduce_paid_child` outcomes now explicitly expose which side paid through a child
- child-triggered interactions now explicitly expose which source-side and/or target-side attached child participated in contact
- `reproduce_blocked_energy` outcomes now explicitly expose which side or sides lacked enough current capacity
- `reproduce_paid_child` outcomes now explicitly expose which concrete child was consumed on each paying side
- successful reproduction outcomes now explicitly expose which new child IDs were created
- successful reproduction outcomes now explicitly expose which side received each created child
- child-triggered interactions now explicitly expose whether the trigger path was child-to-parent or child-to-child
- successful reproduction outcomes now explicitly expose which deterministic child-distribution case was selected
- reproduction outcomes now explicitly expose the decision-time current capacity value on each side
- reproduction outcomes now explicitly expose the governing threshold and cost constants
- reproduction outcomes now explicitly expose the energy and reserve components that compose each side's current capacity
- Go-side snapshot structs now derive child quantity only from attached children instead of mirroring a separate `ChildrenCount` field
- dead derived-radius scaffolding has been removed and parent radius is now expressed directly as the fixed visible body size
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
    "width": 1200,
    "height": 900
  },
  "player": {
    "id": "player-1",
    "lineage_id": "lineage-player-1",
    "generation": 1,
    "shape": "triangle",
    "x": 608,
    "y": 450,
    "radius": 12,
    "energy": 99,
    "attached_children": [
      {
        "id": "child-2",
        "owner_id": "player-1",
        "orbit_slot": 0,
        "x": 629.2,
        "y": 454.8,
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
      "x": 468,
      "y": 450,
      "radius": 12,
      "energy": 100,
      "attached_children": []
    },
    {
      "id": "circle-3",
      "lineage_id": "lineage-circle-3",
      "generation": 0,
      "shape": "square",
      "x": 740,
      "y": 410,
      "radius": 12,
      "energy": 99,
      "attached_children": [
        {
          "id": "child-1",
          "owner_id": "circle-3",
          "orbit_slot": 0,
          "x": 761.2,
          "y": 414.8,
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
      "x": 632,
      "y": 450,
      "radius": 6
    }
  ]
}
```

## Deliberate Simplifications

- one player-controlled circle only
- five autonomous circles in the default demo world
- autonomous circles prefer the nearest active food target and fall back to deterministic baseline drift when no food target is available
- deterministic shape assignment in the default demo world: player `triangle`, same-shape autonomous `triangle`, different-shape autonomous `square`
- deterministic fixed food placement scaled to the expanded default world
- the expanded default world now derives its starting food count as active circles plus two, capped by the available deterministic slot set, while narrow custom worlds keep the smaller three-slot baseline
- local neighborhoods now become more expensive when at least two nearby active circles are within `120` units
- autonomous circles now keep their ordinary steering unless the next step would move them into a neighborhood with more than one additional nearby circle compared with their current local density
- deterministic food regeneration returns consumed slots to their original positions after a fixed delay
- deterministic food regeneration returns consumed slots to their original positions, with delay now increasing under deeper current depletion
- attached children remain the source of truth for fight power, replacement continuity, and child-payment rules, and are embodied as attached orbiting children
- attached-child positions now extend a parent's effective food collection reach
- different-shape reproduction resolves through deterministic child redistribution across the participating parents when both participants satisfy the energy rule
- continuity is limited to one-child replacement after fight defeat
- lineage is represented only by a stable `lineage_id` plus monotonic `generation`
- zero energy is now a death threshold rather than only a movement stop condition
- demo reset recreates the initial world state without restarting the server process
- no local prediction or interpolation
- one shared movement intent for the connected client
- default world size now uses an expanded bounded map while custom config worlds can remain smaller and deterministic

## Surfaced Provisional Rules

The slice needed these implementation choices not fully specified in the refined artifacts:

- when player energy reaches zero, movement stops and energy clamps at zero
- food grants a fixed energy recovery amount of `10`
- each consumed food slot regenerates after `12` ticks when it is the only missing slot, with the delay increasing by `2` ticks for each additional currently missing slot
- each active circle pays an additional `1` energy on a tick when at least two other active circles are within `120` units
- autonomous crowding-aware steering reverses the chosen direction only when the next step would increase local nearby-circle count by more than `1`
- attached-child food collection gives the same parent energy gain as parent-body food collection
- attached-child body contact against another parent body now counts as authoritative pair contact for interaction triggering
- attached-child-to-attached-child overlap also counts as authoritative pair contact for interaction triggering
- autonomous-autonomous pairs now resolve after player-autonomous checks using deterministic index order
- autonomous steering now prefers nearby food within `140` units, otherwise seeks the nearest active circle among the player and autonomous candidates before falling back to baseline drift
- autonomous steering now also prefers food whenever energy is below `40`, even if an interaction target is otherwise eligible
- once interaction seeking is active, target ordering now prefers different-shape candidates over same-shape candidates before applying distance and ID tie-breaks
- different-shape candidates are only preferred when both participants currently satisfy the existing reproduction-capacity rule
- when no feasible reproduction target exists, same-shape candidates are only eligible when the current fight ordering says the acting autonomous circle would win
- when no social target is currently eligible, autonomous steering falls back to the nearest available food target before baseline drift
- same-shape threats within `120` units now override ordinary pursuit when the current fight ordering says the acting autonomous circle would lose
- threat avoidance steers directly away from the nearest qualifying same-shape threat, breaking equal-distance ties by target ID
- blocked different-shape targets within `120` units now override ordinary pursuit when the current reproduction feasibility rule says the pair cannot currently reproduce
- blocked-reproduction avoidance steers directly away from the nearest qualifying different-shape target, breaking equal-distance ties by target ID
- same-shape threat and blocked-reproduction avoidance now measure nearness against the target parent body and its current attached-child positions
- interaction-seeking target selection now measures social nearness against the target parent body and its current attached-child positions
- food-target selection now measures food nearness against the acting parent body and its current attached-child positions
- resolved reproduction now distinguishes ordinary `reproduce_resolved` from `reproduce_paid_child` when a participant consumed one child as reserve payment
- food collection now uses a fixed parent-core body plus attached-child bodies rather than the derived grown radius as silent collection reach
- circle contact now uses a fixed parent-core body plus attached-child bodies rather than the derived grown radius as silent parent contact reach
- same-shape fight ordering now resolves exact ties deterministically without consulting derived radius
- parent movement clamping now uses the fixed visible parent-core body rather than the derived grown radius
- attached-child orbit distance now uses the fixed visible parent-core body rather than the derived grown radius
- parent `radius` now stays fixed at the visible parent-core body even when child count changes
- successful reproduction now exposes whether deterministic redistribution resolved as `source_only`, `split`, or `target_only`
- reproduction outcomes now expose source-side and target-side current capacity values from the authoritative decision point before payment is applied
- reproduction outcomes now also expose the current threshold and cost constants used by that same authoritative rule
- reproduction outcomes now also expose how much of each side's reported capacity comes from direct energy versus child reserve contribution
- child-dependent rule evaluation now derives child quantity from attached children directly
- runtime snapshots no longer expose `children_count`; contract consumers derive child quantity from attached children
- same-shape fight ordering now uses child presence only after energy, not raw child-count magnitude
- interaction provenance now records whether contact came from `parent_body` or `attached_child`
- player energy clamps to a maximum of `100`
- autonomous circles deterministically select the nearest active food by distance, breaking ties by food ID
- when no food target is available, autonomous circles fall back to deterministic index-based directions
- same-shape overlap resolves as `fight_resolved`; different-shape overlap resolves as `reproduce_resolved`
- same-shape overlap may also resolve as `fight_absorbed_child` when the loser has an attached child available to absorb the loss
- zero-energy collapse that preserves continuity now resolves as `death_promoted_child`
- when continuity promotes one attached child, the continuing active parent now takes that child's last visible position at the current tick
- `death_promoted_child` outcomes now expose the promoted child ID explicitly
- `fight_absorbed_child` outcomes now expose the absorbed child ID explicitly
- `reproduce_paid_child` outcomes now expose whether the source side, target side, or both paid through a child
- `reproduce_paid_child` outcomes now also expose the concrete consumed source-side and/or target-side payment child IDs
- successful reproduction outcomes now also expose the concrete created child IDs allocated by the deterministic child-creation path
- successful reproduction outcomes now also expose which created child IDs were allocated to the source side and which were allocated to the target side
- `attached_child` contact origin now also exposes the participating source-side and/or target-side child IDs explicitly
- `attached_child` contact origin now also exposes the concrete trigger path kind used by the existing deterministic contact geometry
- `reproduce_blocked_energy` outcomes now expose whether the source side, target side, or both failed the current capacity check
- Go-side snapshot readers now derive child quantity directly from attached children with no mirrored `ChildrenCount` field
- pair contact may be initiated by one attached child touching the other parent's body even when the two parent cores do not yet overlap
- pair contact may also be initiated by one attached child touching another attached child across the pair
- autonomous-only encounters use the same fight, reproduction, child, and continuity rules as player-involved encounters
- interaction-seeking autonomy now considers the player and other active autonomous circles together, breaking equal-distance ties by target ID
- the low-energy food-priority threshold is checked before interaction seeking
- shape meaning now influences interaction-target ordering before nearest-target fallback
- infeasible reproduction targets now fall back to the nearest remaining eligible target
- same-shape fight ordering is: higher energy, then higher child count, then larger radius, then player exact tie-break
- same-shape fight ordering is now: higher energy, then higher child count, then deterministic exact tie resolution
- a same-shape loser with at least one attached child loses exactly one attached child and remains active on that tick
- different-shape overlap without enough energy or an available child payment resolves as `reproduce_blocked_energy`
- resolved reproduction creates exactly two new child units and assigns them deterministically across the participating parents
- attached children remain orbiting dependents of their current parent rather than immediately becoming free autonomous circles
- attached child orbit positions are derived deterministically from owner, child identity, slot, and tick
- reproduction requires a total reproduction capacity of at least `15`
- reproduction costs `10` from each participant on success
- one child may contribute one `10`-point reserve unit toward the threshold check and payment path
- a circle below the reproduction cost may consume exactly one child, convert it into a `10`-point temporary reserve, pay the reproduction cost, and then receive the reproduction result
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
- autonomous steering now reacts to regional density in a bounded way by reversing a chosen direction when the next step would enter a clearly denser nearby region and the opposite step is less regionally crowded
- food collection now applies a bounded regional yield penalty based on nearby missing food slots around the collected slot, so depleted regions recover less energy even when food reappears
- successful reproduction now applies a bounded regional surcharge based on nearby missing food slots around the reproduction location, while keeping threshold checks and child-payment semantics on the existing baseline path
- the expanded default world now preserves the existing aspect ratio at approximately `10x` the previous total area, while keeping the viewport/minimap model unchanged and increasing the seeded expanded startup population baseline to avoid obvious emptiness
- orientation refresh transport is now driven by deterministic compact-summary change plus a slower fallback interval, instead of a fixed short cadence alone
- local viewport food transport is now driven by visible-food change plus a slower fallback interval, while local circle detail still arrives every tick
- passive observer transport now refreshes on observer-relevant change plus a slower fallback interval, instead of resending the same observer-oriented snapshot whenever the passive cadence timer alone elapses
- passive observer transport now uses a coarse observer state signature: interaction changes still refresh immediately, food state now refreshes by coarse abundance bucket instead of exact minimap motion, and population loss still refreshes by total autonomous count
- the transport measurement harness now also measures deterministic active-client fanout scaling across an explicit moving-client ladder
- the transport measurement helpers now also measure deterministic active payload component cost across the current major active-detail families without changing runtime behavior
- active clients now keep local detail every tick while whole-world orientation support refreshes at a lower deterministic cadence, relying on the existing client-side minimap cache between fresh orientation ticks
- the measurement helpers now also record active orientation freshness versus staleness over a bounded movement window without changing runtime behavior
- the browser client now records a rolling live draw-duration metric and surfaces it in the HUD as a bounded render-pressure indicator during ordinary play
- the repository now includes a deterministic multi-client websocket measurement harness for aggregate bytes/sec, per-client bytes/sec, aggregate snapshot count, and observed inter-snapshot gap under bounded local load
- the multi-client transport harness now also compares deterministic idle and moving-client runs so active-play pressure can be distinguished from passive observer fanout
- the multi-client transport harness now also measures passive client-count fanout scaling across an explicit count ladder

These keep the loop deterministic and prevent energy drift while staying aligned with energy as the constraining movement resource.

## Validation Targets

- deterministic server tests for child accumulation, attached orbiting child ownership, attached-child food collection, child absorption during hostile conflict, derived radius growth, fight winner selection including child power, loser removal, child replacement continuity on zero-energy collapse, lineage preservation, food-slot regeneration timing, energy-gated reproduction, and food-seeking autonomy
- contract test for explicit snapshot shape including child counts, attached child state, lineage fields, and resolved, blocked, absorbed, or promoted interaction outcomes
- integration tests for WebSocket snapshots with visible orbiting children after reproduction, attached-child food collection, attached-child loss during hostile conflict, child-triggered contact before parent core overlap, child-to-child-triggered contact, child-triggered avoidance before parent-core overlap, child-aware positive pursuit before parent-core nearness, child-aware food pursuit before parent-core nearness, explicit `reproduce_paid_child` versus ordinary `reproduce_resolved` outcomes, autonomous-only interaction outcomes, interaction-seeking autonomous outcomes, player-targeted autonomous outcomes, low-energy food-priority steering, shape-aware target choice, feasibility-aware target fallback, fight-feasibility-aware food fallback, threat avoidance against stronger same-shape circles, blocked-reproduction avoidance against nearby different-shape circles, blocked reproduction by low energy, food-seeking autonomous motion, embodied food collection without derived-radius reach, embodied circle contact without derived-radius reach, child-driven fight outcomes, zero-energy collapse continuity with explicit promotion signaling, deterministic food regeneration, and authoritative reset

## Transport Measurement Notes

- viewport-culling plus dual cadence plus compact minimap summaries plus reduced local precision measured about `1549` bytes/message, about `15494` bytes/sec/client at `10` snapshots/sec
- event-driven orientation refresh with the same compact summary path now measures below that fixed dual-cadence baseline over a deterministic fallback window because unchanged orientation summaries are skipped until they materially change or the slower fallback interval is reached
- event-driven local food transport now measures below the event-driven orientation baseline over a deterministic fallback window because unchanged visible-food detail is omitted until it materially changes or the local fallback interval is reached
- deterministic multi-client transport pressure with `4` local websocket clients over `300ms` measured `51904` aggregate bytes across `16` snapshots, about `173013` aggregate bytes/sec and about `43253` bytes/sec/client, with a measured max inter-snapshot gap of `100ms` against an expected tick of `100ms`
- deterministic idle-versus-moving comparison with `4` local websocket clients over `300ms` measured:
  - idle: `51904` aggregate bytes, about `173013` aggregate bytes/sec, `100ms` max gap
  - one moving client: `51892` aggregate bytes, about `172973` aggregate bytes/sec, `101ms` max gap
- under this bounded case, active movement changes the transport profile only slightly, which suggests the current dominant cost is still broad snapshot fanout rather than a large movement-specific throughput spike
- deterministic passive fanout scaling over `300ms` measured:
  - `1` client: `12976` aggregate bytes, about `43253` bytes/sec, `103ms` max gap, `4` snapshots
  - `4` clients: `51904` aggregate bytes, about `173013` bytes/sec, `100ms` max gap, `16` snapshots
  - `8` clients: `103808` aggregate bytes, about `346027` bytes/sec, `100ms` max gap, `32` snapshots
- under this bounded ladder, aggregate output scales almost linearly with passive client count while per-client throughput stays flat, which confirms fanout itself as the clearest current transport pressure
- idle-observer cadence reduction now applies only when there is real observer fanout, so the single-client path keeps the prior cadence while passive multi-client observers receive deterministic lower-cadence snapshots
- deterministic passive fanout scaling over `300ms` with the new cadence policy now measures:
  - `1` client: `12976` aggregate bytes, about `43253` bytes/sec, `101ms` max gap, `4` snapshots
  - `4` clients: `25752` aggregate bytes, about `85840` bytes/sec, `301ms` max gap, `8` snapshots
  - `8` clients: `51504` aggregate bytes, about `171680` bytes/sec, `300ms` max gap, `16` snapshots
- under this bounded ladder, passive aggregate fanout now drops by roughly half at `4` and `8` clients while keeping single-client cadence unchanged
- deterministic mixed `4`-client pressure with one active steering client over `300ms` now measures `32284` aggregate bytes across `10` snapshots, about `107613` aggregate bytes/sec and about `26903` bytes/sec/client, with the active client receiving `4` snapshots while each passive observer receives `2`
- passive multi-client observers now receive an explicit `observer_orientation_only` transport mode instead of the active local-detail view
- observer-oriented snapshots keep world bounds, totals, interaction state, and minimap orientation summaries, while omitting local player, autonomous, and food detail
- active snapshots keep the existing `active_local_detail` transport mode and unchanged local-detail path
- deterministic passive fanout scaling over `300ms` with observer-oriented snapshots now measures:
  - `1` client: `13132` aggregate bytes, about `43773` bytes/sec, `100ms` max gap, `4` snapshots
  - `4` clients: `22988` aggregate bytes, about `76627` bytes/sec, `301ms` max gap, `8` snapshots
  - `8` clients: `45976` aggregate bytes, about `153253` bytes/sec, `301ms` max gap, `16` snapshots
- under this bounded ladder, orientation-only passive snapshots cut the current post-cadence `4`-client passive baseline from `25752` aggregate bytes to `22988`, and the `8`-client passive baseline from `51504` to `45976`
- deterministic mixed `4`-client pressure with one active steering client over `300ms` now measures `30370` aggregate bytes across `10` snapshots, about `101233` aggregate bytes/sec and about `25308` bytes/sec/client, with the active client receiving `4` snapshots while each passive observer receives `2`
- deterministic passive fanout scaling over `300ms` with event-driven calm observer refresh now measures:
  - `1` client: `13132` aggregate bytes, about `43773` bytes/sec, `101ms` max gap, `4` snapshots
  - `4` clients: `13336` aggregate bytes, about `44453` bytes/sec, `0ms` max gap, `4` snapshots
  - `8` clients: `26672` aggregate bytes, about `88907` bytes/sec, `0ms` max gap, `8` snapshots
- under this bounded calm ladder, event-driven observer refresh reduces the `4`-client passive baseline from `22988` aggregate bytes to `13336`, and the `8`-client passive baseline from `45976` to `26672`
- deterministic mixed `4`-client pressure with one active steering client over `300ms` now measures `23131` aggregate bytes across `7` snapshots, about `77103` aggregate bytes/sec and about `19276` bytes/sec/client, with the active client receiving `4` snapshots while each passive observer receives `1`
- the coarse observer signature keeps the current calm passive baseline unchanged in the bounded expanded-world ladder, while restoring passive refresh on deterministic small-world food-count change and autonomous-count change before fallback
- deterministic active-client fanout scaling over `300ms` now measures:
  - `1` active client: `13129` aggregate bytes, about `43763` bytes/sec, `102ms` max gap, `4` snapshots
  - `2` active clients: `26258` aggregate bytes, about `87527` bytes/sec, `102ms` max gap, `8` snapshots
  - `4` active clients: `52516` aggregate bytes, about `175053` bytes/sec, `104ms` max gap, `16` snapshots
- under this bounded ladder, active fanout now scales almost linearly with client count while per-client throughput stays flat at about `43763` bytes/sec, which makes the active local-detail path the clearest remaining transport pressure
- deterministic active payload component measurement on the current default active snapshot now measures:
  - full active payload: `3333` bytes
  - without player detail: `3097` bytes
  - without local autonomous detail: `2946` bytes
  - without local food detail: `3105` bytes
  - without orientation support: `1176` bytes
  - without interaction detail: `3333` bytes in the current no-interaction baseline
- under this bounded breakdown, orientation support is the dominant serialized active payload family in the current default snapshot
- deterministic optimized active transport over `300ms` now measures:
  - `1` active client: `8942` aggregate bytes, about `29807` bytes/sec, `102ms` max gap, `4` snapshots
  - `2` active clients: `17884` aggregate bytes, about `59613` bytes/sec, `102ms` max gap, `8` snapshots
  - `4` active clients: `35768` aggregate bytes, about `119227` bytes/sec, `102ms` max gap, `16` snapshots
- under this bounded ladder, reducing active orientation support cuts the `1`-client active baseline from `13129` bytes to `8942` and the `4`-client active baseline from `52516` to `35768` while keeping full active snapshot cadence
- deterministic active orientation usability over a `300ms` movement window now measures `4` total active snapshots with `2` fresh and `2` stale orientation snapshots at the current optimized cadence
- under that bounded movement window, active local detail remained continuous while orientation support was stale on half the snapshots
- client render pressure is now measured live in the browser as rolling average and max draw duration over recent draws, which makes viewport rendering pressure explicit without changing transport or gameplay behavior
- client render pressure is now also broken into rolling major render families: `world`, `overlay`, `minimap`, and `support`
- the HUD now shows the current dominant render family and exposes the rolling family breakdown through the existing render-pressure indicator
- this keeps default runtime behavior unchanged while making the next client render optimization target evidence-backed during ordinary play
