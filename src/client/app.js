import { CONTRACT_VERSION, createMovementIntent, MESSAGE_TYPES } from "/shared_contracts/messages/protocol.js";

const canvas = document.getElementById("world");
const context = canvas.getContext("2d");
const statusNode = document.getElementById("status");
const energyNode = document.getElementById("energy");
const tickNode = document.getElementById("tick");
const resetButton = document.getElementById("reset");

let latestSnapshot = null;
const pressedKeys = new Set();
let activeSocket = null;
let senderIntervalId = null;

const movementKeys = new Set(["arrowleft", "arrowright", "arrowup", "arrowdown", "w", "a", "s", "d"]);

function normalizeKey(key) {
  return key.toLowerCase();
}

function currentDirection() {
  const vector = { x: 0, y: 0 };

  if (pressedKeys.has("arrowleft") || pressedKeys.has("a")) {
    vector.x -= 1;
  }
  if (pressedKeys.has("arrowright") || pressedKeys.has("d")) {
    vector.x += 1;
  }
  if (pressedKeys.has("arrowup") || pressedKeys.has("w")) {
    vector.y -= 1;
  }
  if (pressedKeys.has("arrowdown") || pressedKeys.has("s")) {
    vector.y += 1;
  }

  return vector;
}

function draw(snapshot) {
  canvas.width = snapshot.world.width;
  canvas.height = snapshot.world.height;

  context.clearRect(0, 0, canvas.width, canvas.height);

  context.fillStyle = "rgba(23, 49, 58, 0.06)";
  context.fillRect(0, 0, canvas.width, canvas.height);

  context.strokeStyle = "rgba(23, 49, 58, 0.15)";
  context.lineWidth = 2;
  context.strokeRect(1, 1, canvas.width - 2, canvas.height - 2);

  context.fillStyle = "#d85f3d";
  for (const food of snapshot.foods) {
    context.beginPath();
    context.arc(food.x, food.y, food.radius, 0, Math.PI * 2);
    context.fill();
  }

  for (const circle of snapshot.autonomous_circles) {
    drawCircle(circle, false, snapshot.player);
  }

  if (snapshot.player) {
    drawCircle(snapshot.player, true, snapshot.player);
    energyNode.textContent = `Energy: ${snapshot.player.energy.toFixed(0)} · Children: ${snapshot.player.children_count} · Generation: ${snapshot.player.generation}`;
  } else {
    energyNode.textContent = "Energy: defeated";
  }

  const interaction = snapshot.interaction ? snapshot.interaction.kind : "none";
  tickNode.textContent = `Tick: ${snapshot.tick} · Food: ${snapshot.foods.length} · Others: ${snapshot.autonomous_circles.length} · Interaction: ${interaction}`;
}

function drawCircle(circle, isPlayer, player) {
  const matchesPlayerShape = !isPlayer && player && circle.shape === player.shape;
  const color = isPlayer ? "#d85f3d" : circle.shape === "triangle" ? "#3b8ea5" : "#8c6bb1";
  context.fillStyle = color;

  if (circle.shape === "triangle") {
    context.beginPath();
    context.moveTo(circle.x, circle.y - circle.radius);
    context.lineTo(circle.x - circle.radius, circle.y + circle.radius);
    context.lineTo(circle.x + circle.radius, circle.y + circle.radius);
    context.closePath();
    context.fill();
  } else if (circle.shape === "square") {
    context.fillRect(circle.x - circle.radius, circle.y - circle.radius, circle.radius * 2, circle.radius * 2);
  } else {
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius, 0, Math.PI * 2);
    context.fill();
  }

  if (matchesPlayerShape) {
    context.strokeStyle = "#b63b29";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 5, 0, Math.PI * 2);
    context.stroke();
  }

  if (isPlayer) {
    context.strokeStyle = "#f4d35e";
    context.lineWidth = 4;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 8, 0, Math.PI * 2);
    context.stroke();

    context.fillStyle = "#f7f3e8";
    context.beginPath();
    context.arc(circle.x, circle.y, 3, 0, Math.PI * 2);
    context.fill();
  }

  drawAttachedChildren(circle, color);

  context.fillStyle = "#17313a";
  context.font = "16px Georgia";
  const label = isPlayer
    ? `YOU ${circle.id} (${circle.shape}) c:${circle.children_count} o:${circle.attached_children.length} g:${circle.generation}`
    : `${circle.id} ${matchesPlayerShape ? "match" : "other"} (${circle.shape}) c:${circle.children_count} o:${circle.attached_children.length} g:${circle.generation}`;
  context.fillText(label, circle.x - 40, circle.y - circle.radius - 10);

  context.font = "12px Georgia";
  context.fillText(circle.lineage_id, circle.x - 40, circle.y + circle.radius + 18);
}

function drawAttachedChildren(circle, color) {
  for (const child of circle.attached_children) {
    context.fillStyle = color;
    context.globalAlpha = 0.7;
    context.beginPath();
    context.arc(child.x, child.y, child.radius, 0, Math.PI * 2);
    context.fill();

    context.globalAlpha = 1;
    context.strokeStyle = "#17313a";
    context.lineWidth = 1;
    context.beginPath();
    context.arc(child.x, child.y, child.radius, 0, Math.PI * 2);
    context.stroke();
  }
}

function setStatus(message) {
  statusNode.textContent = `${message} · contract v${CONTRACT_VERSION}`;
}

async function resetWorld() {
  if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) {
    setStatus("Reset unavailable while disconnected");
    return;
  }

  resetButton.disabled = true;
  setStatus("Resetting world");

  try {
    const response = await fetch("/reset", {
      method: "POST",
    });

    if (!response.ok) {
      throw new Error(`reset failed with status ${response.status}`);
    }

    const snapshot = await response.json();
    latestSnapshot = snapshot;
    draw(snapshot);
    setStatus("Connected");
  } catch (_error) {
    setStatus("Reset failed");
  } finally {
    resetButton.disabled = false;
  }
}

function ensureSenderLoop() {
  if (senderIntervalId !== null) {
    return;
  }

  senderIntervalId = window.setInterval(() => {
    if (!activeSocket || activeSocket.readyState !== WebSocket.OPEN) {
      return;
    }

    const direction = currentDirection();
    activeSocket.send(JSON.stringify(createMovementIntent(direction.x, direction.y)));
  }, 100);
}

function connect() {
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const socket = new WebSocket(`${protocol}://${window.location.host}/ws`);
  activeSocket = socket;

  socket.addEventListener("open", () => {
    setStatus("Connected");
    ensureSenderLoop();
  });

  socket.addEventListener("message", (event) => {
    const snapshot = JSON.parse(event.data);
    if (snapshot.type !== MESSAGE_TYPES.worldSnapshot) {
      return;
    }

    latestSnapshot = snapshot;
    draw(snapshot);
  });

  socket.addEventListener("close", () => {
    if (activeSocket === socket) {
      activeSocket = null;
    }
    setStatus("Disconnected");
    setTimeout(connect, 1000);
  });

  socket.addEventListener("error", () => {
    setStatus("Connection error");
  });

}

function handleMovementKey(event, isPressed) {
  const key = normalizeKey(event.key);
  if (!movementKeys.has(key)) {
    return;
  }

  event.preventDefault();
  if (isPressed) {
    pressedKeys.add(key);
  } else {
    pressedKeys.delete(key);
  }
}

window.addEventListener("keydown", (event) => {
  handleMovementKey(event, true);
});

window.addEventListener("keyup", (event) => {
  handleMovementKey(event, false);
});

window.addEventListener("blur", () => {
  pressedKeys.clear();
});

resetButton.addEventListener("click", () => {
  resetWorld();
});

setStatus("Connecting");
connect();

if (!latestSnapshot) {
  context.fillStyle = "#17313a";
  context.font = "24px Georgia";
  context.fillText("Waiting for server snapshot...", 24, 48);
}
