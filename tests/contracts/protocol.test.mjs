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
  assert.deepEqual(schema.required, ["type", "tick", "world", "player", "autonomous_circles", "foods"]);
  assert.deepEqual(schema.properties.autonomous_circles.items.required, ["id", "x", "y", "radius", "energy"]);
  assert.deepEqual(schema.properties.foods.items.required, ["id", "x", "y", "radius"]);
});
