export const CONTRACT_VERSION = 1;

export const MESSAGE_TYPES = {
  movementIntent: "movement_intent",
  worldSnapshot: "world_snapshot",
};

export function createMovementIntent(x, y) {
  return {
    type: MESSAGE_TYPES.movementIntent,
    direction: {
      x,
      y,
    },
  };
}
