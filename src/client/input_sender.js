function isNeutral(direction) {
  return direction.x === 0 && direction.y === 0;
}

export function shouldSendMovementIntent(lastSentDirection, nextDirection) {
  if (!lastSentDirection) {
    return !isNeutral(nextDirection);
  }

  return lastSentDirection.x !== nextDirection.x || lastSentDirection.y !== nextDirection.y;
}
