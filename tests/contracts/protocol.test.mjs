import test from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";

import { CONTRACT_VERSION, MESSAGE_TYPES, createMovementIntent } from "../../src/shared_contracts/messages/protocol.js";

test("movement intent helper produces the expected message shape", () => {
  const message = createMovementIntent(1, -1);

  assert.equal(CONTRACT_VERSION, 1);
  assert.deepEqual(message, {
    type: MESSAGE_TYPES.movementIntent,
    direction: {
      x: 1,
      y: -1,
    },
  });
});

test("world snapshot schema remains explicit and parseable", async () => {
  const file = await readFile(new URL("../../src/shared_contracts/schemas/world_snapshot.schema.json", import.meta.url), "utf8");
  const schema = JSON.parse(file);

  assert.equal(schema.properties.type.const, MESSAGE_TYPES.worldSnapshot);
  assert.deepEqual(schema.required, ["type", "tick", "world", "player", "autonomous_circles", "interaction", "foods"]);
  assert.deepEqual(schema.properties.player.anyOf[1].required, ["id", "lineage_id", "generation", "shape", "x", "y", "radius", "energy", "children_count", "attached_children"]);
  assert.deepEqual(schema.properties.autonomous_circles.items.required, ["id", "lineage_id", "generation", "shape", "x", "y", "radius", "energy", "children_count", "attached_children"]);
  assert.deepEqual(schema.properties.player.anyOf[1].properties.attached_children.items.required, ["id", "owner_id", "orbit_slot", "x", "y", "radius"]);
  assert.deepEqual(schema.properties.autonomous_circles.items.properties.attached_children.items.required, ["id", "owner_id", "orbit_slot", "x", "y", "radius"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].required, ["active", "resolved", "kind", "source_id", "target_id"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.kind.enum, ["reproduce_resolved", "reproduce_blocked_energy", "fight_resolved", "fight_absorbed_child", "death_promoted_child"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.contact_origin.enum, ["parent_body", "attached_child"]);
  assert.deepEqual(schema.properties.foods.items.required, ["id", "x", "y", "radius"]);
});
