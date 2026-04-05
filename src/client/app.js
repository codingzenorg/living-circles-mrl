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

function childCount(circle) {
  return circle.attached_children.length;
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

  context.fillStyle = "rgba(6, 22, 30, 0.9)";
  context.fillRect(0, 0, canvas.width, canvas.height);

  context.strokeStyle = "rgba(127, 174, 188, 0.28)";
  context.lineWidth = 2;
  context.strokeRect(1, 1, canvas.width - 2, canvas.height - 2);

  context.fillStyle = "#ff8a5b";
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
    energyNode.textContent = `Energy: ${snapshot.player.energy.toFixed(0)} · Children: ${childCount(snapshot.player)} · Generation: ${snapshot.player.generation}`;
  } else {
    energyNode.textContent = "Energy: defeated";
  }

  const interaction = snapshot.interaction ? snapshot.interaction.kind : "none";
  const promotedChild = snapshot.interaction?.promoted_child_id ? ` · Promoted: ${snapshot.interaction.promoted_child_id}` : "";
  const absorbedChild = snapshot.interaction?.absorbed_child_id ? ` · Absorbed: ${snapshot.interaction.absorbed_child_id}` : "";
  const childPayment = snapshot.interaction?.source_paid_child || snapshot.interaction?.target_paid_child
    ? ` · Paid: ${snapshot.interaction?.source_paid_child ? "source" : ""}${snapshot.interaction?.source_paid_child && snapshot.interaction?.target_paid_child ? "+" : ""}${snapshot.interaction?.target_paid_child ? "target" : ""}`
    : "";
  const blockedCapacity = snapshot.interaction?.source_blocked_capacity || snapshot.interaction?.target_blocked_capacity
    ? ` · Blocked: ${snapshot.interaction?.source_blocked_capacity ? "source" : ""}${snapshot.interaction?.source_blocked_capacity && snapshot.interaction?.target_blocked_capacity ? "+" : ""}${snapshot.interaction?.target_blocked_capacity ? "target" : ""}`
    : "";
  const sourceChild = snapshot.interaction?.source_child_id ? ` · Source child: ${snapshot.interaction.source_child_id}` : "";
  const targetChild = snapshot.interaction?.target_child_id ? ` · Target child: ${snapshot.interaction.target_child_id}` : "";
  tickNode.textContent = `Tick: ${snapshot.tick} · Food: ${snapshot.foods.length} · Others: ${snapshot.autonomous_circles.length} · Interaction: ${interaction}${promotedChild}${absorbedChild}${childPayment}${blockedCapacity}${sourceChild}${targetChild}`;
}

function drawCircle(circle, isPlayer, player) {
  const matchesPlayerShape = !isPlayer && player && circle.shape === player.shape;
  const color = isPlayer ? "#ff8a5b" : circle.shape === "triangle" ? "#6fd5ff" : "#c08cff";
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
    context.strokeStyle = "#ffad8c";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 5, 0, Math.PI * 2);
    context.stroke();
  }

  if (isPlayer) {
    context.strokeStyle = "#ffe082";
    context.lineWidth = 4;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 8, 0, Math.PI * 2);
    context.stroke();

    context.fillStyle = "#fff3d4";
    context.beginPath();
    context.arc(circle.x, circle.y, 3, 0, Math.PI * 2);
    context.fill();
  }

  drawAttachedChildren(circle, color);

  context.fillStyle = "#e4f3f8";
  context.font = "16px Georgia";
  const children = childCount(circle);
  const label = isPlayer
    ? `YOU ${circle.id} (${circle.shape}) c:${children} o:${circle.attached_children.length} g:${circle.generation}`
    : `${circle.id} ${matchesPlayerShape ? "match" : "other"} (${circle.shape}) c:${children} o:${circle.attached_children.length} g:${circle.generation}`;
  context.fillText(label, circle.x - 40, circle.y - circle.radius - 10);

  context.font = "12px Georgia";
  context.fillStyle = "#9cb8c0";
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
    context.strokeStyle = "#d7eef6";
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
  context.fillStyle = "#e4f3f8";
  context.font = "24px Georgia";
  context.fillText("Waiting for server snapshot...", 24, 48);
}
