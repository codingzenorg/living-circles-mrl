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
  assert.deepEqual(schema.required, ["type", "transport_mode", "tick", "world", "player", "autonomous_circles", "interaction", "foods", "foods_fresh", "orientation_fresh", "total_autonomous_circles", "total_foods"]);
  assert.deepEqual(schema.properties.transport_mode.enum, ["active_local_detail", "observer_orientation_only"]);
  assert.deepEqual(schema.properties.player.anyOf[1].required, ["id", "lineage_id", "generation", "shape", "x", "y", "radius", "energy", "attached_children"]);
  assert.deepEqual(schema.properties.autonomous_circles.items.required, ["id", "lineage_id", "generation", "shape", "x", "y", "radius", "energy", "attached_children"]);
  assert.deepEqual(schema.properties.minimap_autonomous_circles.anyOf[0].items.required, ["shape", "x", "y", "count"]);
  assert.deepEqual(schema.properties.player.anyOf[1].properties.attached_children.items.required, ["id", "owner_id", "orbit_slot", "x", "y", "radius"]);
  assert.deepEqual(schema.properties.autonomous_circles.items.properties.attached_children.items.required, ["id", "owner_id", "orbit_slot", "x", "y", "radius"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].required, ["active", "resolved", "kind", "source_id", "target_id"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.kind.enum, ["reproduce_resolved", "reproduce_paid_child", "reproduce_blocked_energy", "fight_resolved", "fight_absorbed_child", "death_promoted_child"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.contact_origin.enum, ["parent_body", "attached_child"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.contact_path_kind.enum, ["source_child_to_target_parent", "source_parent_to_target_child", "child_to_child"]);
  assert.deepEqual(schema.properties.interaction.anyOf[1].properties.distribution_kind.enum, ["source_only", "split", "target_only"]);
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.promoted_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.absorbed_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_paid_child.type, "boolean");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_paid_child.type, "boolean");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_paid_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_paid_child_id.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.created_child_ids.type, "array");
  assert.equal(schema.properties.interaction.anyOf[1].properties.created_child_ids.items.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_created_child_ids.type, "array");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_created_child_ids.items.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_created_child_ids.type, "array");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_created_child_ids.items.type, "string");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_blocked_capacity.type, "boolean");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_blocked_capacity.type, "boolean");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_capacity_value.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_capacity_value.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_energy_component.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_energy_component.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.source_reserve_component.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.target_reserve_component.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.reproduction_threshold.type, "number");
  assert.equal(schema.properties.interaction.anyOf[1].properties.reproduction_cost.type, "number");
  assert.equal(schema.properties.foods_fresh.type, "boolean");
  assert.equal(schema.properties.orientation_fresh.type, "boolean");
  assert.deepEqual(schema.properties.foods.anyOf[0].items.required, ["id", "x", "y", "radius"]);
  assert.deepEqual(schema.properties.minimap_foods.anyOf[0].items.required, ["x", "y", "count"]);
  assert.equal(schema.properties.total_autonomous_circles.type, "integer");
  assert.equal(schema.properties.total_foods.type, "integer");
});
