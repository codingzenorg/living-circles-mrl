import { CONTRACT_VERSION, createMovementIntent, MESSAGE_TYPES } from "/shared_contracts/messages/protocol.js";

const canvas = document.getElementById("world");
const context = canvas.getContext("2d");
const statusNode = document.getElementById("status");
const energyNode = document.getElementById("energy");
const tickNode = document.getElementById("tick");
const detailsNode = document.getElementById("details");
const resetButton = document.getElementById("reset");

let latestSnapshot = null;
let eventLog = ["Awaiting first world snapshot..."];
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

function distanceBetween(a, b) {
  return Math.hypot(a.x - b.x, a.y - b.y);
}

function playerRiskState(circle, player, interaction) {
  if (!player || circle.id === player.id) {
    return "neutral";
  }

  const nearPlayer = distanceBetween(circle, player) <= 220;
  if (!nearPlayer) {
    return "neutral";
  }

  if (circle.shape === player.shape) {
    return "danger";
  }

  if (!interaction) {
    return "opportunity";
  }

  const involvesCircle = interaction.source_id === circle.id || interaction.target_id === circle.id;
  const involvesPlayer = interaction.source_id === player.id || interaction.target_id === player.id;
  if (!involvesCircle || !involvesPlayer) {
    return "opportunity";
  }

  if (interaction.kind === "reproduce_blocked_energy") {
    return "blocked";
  }

  if (interaction.kind.startsWith("reproduce")) {
    return "opportunity";
  }

  return "danger";
}

function interactionSummary(interaction) {
  if (!interaction) {
    return null;
  }

  const summaries = {
    fight_candidate: "Same-shape collision: contest is active.",
    fight_absorbed_child: "Same-shape collision: one attached child absorbed the loss.",
    fight_resolved: "Same-shape collision: the fight resolved.",
    reproduce_candidate: "Different-shape collision: reproduction is possible.",
    reproduce_resolved: "Different-shape collision: reproduction resolved.",
    reproduce_paid_child: "Different-shape collision: reproduction used child reserve.",
    reproduce_blocked_energy: "Different-shape collision: reproduction is blocked by current capacity.",
    death_promoted_child: "Collapse continuity: a promoted child preserved the lineage.",
  };

  return summaries[interaction.kind] ?? `Interaction: ${interaction.kind}`;
}

function renderEventLog() {
  detailsNode.innerHTML = eventLog.map((entry) => `<li>${entry}</li>`).join("");
}

function pushEventLog(message) {
  if (!message) {
    return;
  }

  if (eventLog[0] === message) {
    return;
  }

  eventLog = [message, ...eventLog].slice(0, 5);
  renderEventLog();
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
    energyNode.textContent = `E ${snapshot.player.energy.toFixed(0)} · C ${childCount(snapshot.player)} · G ${snapshot.player.generation}`;
  } else {
    energyNode.textContent = "E defeated";
  }

  tickNode.textContent = `T ${snapshot.tick} · F ${snapshot.foods.length} · O ${snapshot.autonomous_circles.length}`;
  pushEventLog(interactionSummary(snapshot.interaction));
}

function drawCircle(circle, isPlayer, player) {
  const matchesPlayerShape = !isPlayer && player && circle.shape === player.shape;
  const relationToPlayer = playerRiskState(circle, player, latestSnapshot?.interaction);
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
    context.strokeStyle = "#ff6f7f";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 5, 0, Math.PI * 2);
    context.stroke();
  }

  if (!isPlayer && relationToPlayer === "danger") {
    context.strokeStyle = "rgba(255, 95, 109, 0.9)";
    context.lineWidth = 6;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 10, 0, Math.PI * 2);
    context.stroke();
  }

  if (!isPlayer && relationToPlayer === "opportunity") {
    context.strokeStyle = "rgba(117, 229, 149, 0.9)";
    context.lineWidth = 5;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 9, 0, Math.PI * 2);
    context.stroke();
  }

  if (!isPlayer && relationToPlayer === "blocked") {
    context.setLineDash([8, 6]);
    context.strokeStyle = "rgba(255, 201, 92, 0.95)";
    context.lineWidth = 5;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 9, 0, Math.PI * 2);
    context.stroke();
    context.setLineDash([]);
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
    : `${circle.id} ${relationToPlayer === "danger" ? "danger" : relationToPlayer === "opportunity" ? "open" : relationToPlayer === "blocked" ? "blocked" : matchesPlayerShape ? "match" : "other"} (${circle.shape}) c:${children} o:${circle.attached_children.length} g:${circle.generation}`;
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
    eventLog = ["World restarted."];
    renderEventLog();
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
renderEventLog();

if (!latestSnapshot) {
  context.fillStyle = "#e4f3f8";
  context.font = "24px Georgia";
  context.fillText("Waiting for server snapshot...", 24, 48);
}
