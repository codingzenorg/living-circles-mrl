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
const CROWDING_RADIUS = 120;
const CROWDING_THRESHOLD = 2;
const CROWDING_CUE_DISTANCE = 260;
const FOOD_OPPORTUNITY_RADIUS = 170;
const FOOD_CUE_DISTANCE = 260;
const SCARCITY_THRESHOLD = 1;

function normalizeKey(key) {
  return key.toLowerCase();
}

function childCount(circle) {
  return circle.attached_children.length;
}

function distanceBetween(a, b) {
  return Math.hypot(a.x - b.x, a.y - b.y);
}

function nearbyCircles(circle, circles, radius) {
  return circles.filter((other) => other.id !== circle.id && distanceBetween(circle, other) <= radius);
}

function crowdingCount(circle, circles) {
  return nearbyCircles(circle, circles, CROWDING_RADIUS).length;
}

function isCrowded(circle, circles) {
  return crowdingCount(circle, circles) >= CROWDING_THRESHOLD;
}

function shouldRenderCrowdingCue(circle, circles, player) {
  if (!player) {
    return false;
  }

  return isCrowded(circle, circles) && distanceBetween(circle, player) <= CROWDING_CUE_DISTANCE;
}

function nearbyFoods(anchor, foods, radius) {
  return foods.filter((food) => distanceBetween(anchor, food) <= radius);
}

function foodPressureAt(anchor, foods) {
  const nearby = nearbyFoods(anchor, foods, FOOD_OPPORTUNITY_RADIUS);
  return {
    nearbyCount: nearby.length,
    scarcity: nearby.length <= SCARCITY_THRESHOLD,
    opportunity: nearby.length >= 2,
  };
}

function shouldRenderFoodOpportunityCue(anchor, foods, player) {
  if (!player) {
    return false;
  }

  return distanceBetween(anchor, player) <= FOOD_CUE_DISTANCE && foodPressureAt(anchor, foods).opportunity;
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
  const circles = snapshot.player ? [snapshot.player, ...snapshot.autonomous_circles] : [...snapshot.autonomous_circles];
  const playerFoodPressure = snapshot.player ? foodPressureAt(snapshot.player, snapshot.foods) : null;

  canvas.width = snapshot.world.width;
  canvas.height = snapshot.world.height;

  context.clearRect(0, 0, canvas.width, canvas.height);

  context.fillStyle = "rgba(6, 22, 30, 0.9)";
  context.fillRect(0, 0, canvas.width, canvas.height);

  context.strokeStyle = "rgba(127, 174, 188, 0.28)";
  context.lineWidth = 2;
  context.strokeRect(1, 1, canvas.width - 2, canvas.height - 2);

  drawCrowdingZones(circles, snapshot.player);
  drawFoodZones(snapshot.foods, snapshot.player);

  context.fillStyle = "#ff8a5b";
  for (const food of snapshot.foods) {
    const nearbyOpportunity = snapshot.player && distanceBetween(food, snapshot.player) <= FOOD_CUE_DISTANCE;
    if (nearbyOpportunity) {
      const gradient = context.createRadialGradient(food.x, food.y, food.radius, food.x, food.y, 28);
      gradient.addColorStop(0, "rgba(103, 221, 129, 0.28)");
      gradient.addColorStop(1, "rgba(103, 221, 129, 0)");
      context.fillStyle = gradient;
      context.beginPath();
      context.arc(food.x, food.y, 28, 0, Math.PI * 2);
      context.fill();
    }

    context.fillStyle = "#ff8a5b";
    context.beginPath();
    context.arc(food.x, food.y, food.radius, 0, Math.PI * 2);
    context.fill();
  }

  for (const circle of snapshot.autonomous_circles) {
    drawCircle(circle, false, snapshot.player, circles);
  }

  if (snapshot.player) {
    drawCircle(snapshot.player, true, snapshot.player, circles);
    const pressure = isCrowded(snapshot.player, circles) ? " · pressure" : "";
    const foodState = playerFoodPressure?.scarcity ? " · scarce" : playerFoodPressure?.opportunity ? " · food" : "";
    energyNode.textContent = `E ${snapshot.player.energy.toFixed(0)} · C ${childCount(snapshot.player)} · G ${snapshot.player.generation}${pressure}${foodState}`;
  } else {
    energyNode.textContent = "E defeated";
  }

  tickNode.textContent = `T ${snapshot.tick} · F ${snapshot.foods.length} · O ${snapshot.autonomous_circles.length}`;
  pushEventLog(interactionSummary(snapshot.interaction));
}

function drawCrowdingZones(circles, player) {
  for (const circle of circles) {
    if (!shouldRenderCrowdingCue(circle, circles, player)) {
      continue;
    }

    const gradient = context.createRadialGradient(circle.x, circle.y, circle.radius + 8, circle.x, circle.y, 92);
    gradient.addColorStop(0, "rgba(255, 170, 61, 0.18)");
    gradient.addColorStop(0.65, "rgba(255, 128, 61, 0.08)");
    gradient.addColorStop(1, "rgba(255, 128, 61, 0)");
    context.fillStyle = gradient;
    context.beginPath();
    context.arc(circle.x, circle.y, 92, 0, Math.PI * 2);
    context.fill();
  }
}

function drawFoodZones(foods, player) {
  if (!player) {
    return;
  }

  const playerFoodPressure = foodPressureAt(player, foods);

  if (playerFoodPressure.opportunity) {
    const gradient = context.createRadialGradient(player.x, player.y, 16, player.x, player.y, 120);
    gradient.addColorStop(0, "rgba(103, 221, 129, 0.12)");
    gradient.addColorStop(0.7, "rgba(103, 221, 129, 0.06)");
    gradient.addColorStop(1, "rgba(103, 221, 129, 0)");
    context.fillStyle = gradient;
    context.beginPath();
    context.arc(player.x, player.y, 120, 0, Math.PI * 2);
    context.fill();
  }

  if (playerFoodPressure.scarcity) {
    context.strokeStyle = "rgba(106, 196, 226, 0.85)";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(player.x, player.y, 72, 0, Math.PI * 2);
    context.stroke();
  }

  for (const food of foods) {
    if (!shouldRenderFoodOpportunityCue(food, foods, player)) {
      continue;
    }

    const gradient = context.createRadialGradient(food.x, food.y, food.radius + 4, food.x, food.y, 46);
    gradient.addColorStop(0, "rgba(103, 221, 129, 0.18)");
    gradient.addColorStop(1, "rgba(103, 221, 129, 0)");
    context.fillStyle = gradient;
    context.beginPath();
    context.arc(food.x, food.y, 46, 0, Math.PI * 2);
    context.fill();
  }
}

function drawCircle(circle, isPlayer, player, circles) {
  const matchesPlayerShape = !isPlayer && player && circle.shape === player.shape;
  const relationToPlayer = playerRiskState(circle, player, latestSnapshot?.interaction);
  const crowded = shouldRenderCrowdingCue(circle, circles, player);
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

  if (crowded) {
    context.strokeStyle = "rgba(255, 170, 61, 0.85)";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 15, 0, Math.PI * 2);
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
    ? `YOU ${circle.id} (${circle.shape}) c:${children} o:${circle.attached_children.length} g:${circle.generation}${crowded ? " crowded" : ""}`
    : `${circle.id} ${relationToPlayer === "danger" ? "danger" : relationToPlayer === "opportunity" ? "open" : relationToPlayer === "blocked" ? "blocked" : matchesPlayerShape ? "match" : "other"} (${circle.shape}) c:${children} o:${circle.attached_children.length} g:${circle.generation}${crowded ? " crowded" : ""}`;
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
