import test from "node:test";
import assert from "node:assert/strict";

import { shouldSendMovementIntent } from "../../src/client/input_sender.js";

test("idle sender suppresses initial neutral intent", () => {
  assert.equal(shouldSendMovementIntent(null, { x: 0, y: 0 }), false);
});

test("sender emits when movement begins", () => {
  assert.equal(shouldSendMovementIntent(null, { x: 1, y: 0 }), true);
});

test("sender suppresses repeated unchanged movement", () => {
  assert.equal(shouldSendMovementIntent({ x: 1, y: 0 }, { x: 1, y: 0 }), false);
});

test("sender emits one stop intent when movement returns to neutral", () => {
  assert.equal(shouldSendMovementIntent({ x: 1, y: 0 }, { x: 0, y: 0 }), true);
  assert.equal(shouldSendMovementIntent({ x: 0, y: 0 }, { x: 0, y: 0 }), false);
});
