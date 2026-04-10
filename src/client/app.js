import { CONTRACT_VERSION, createMovementIntent, MESSAGE_TYPES } from "/shared_contracts/messages/protocol.js";

const canvas = document.getElementById("world");
const context = canvas.getContext("2d");
const minimapCanvas = document.getElementById("minimap");
const minimapContext = minimapCanvas.getContext("2d");
const playStageNode = document.querySelector(".play-stage");
const statusNode = document.getElementById("status");
const energyNode = document.getElementById("energy");
const tickNode = document.getElementById("tick");
const renderPressureNode = document.getElementById("render-pressure");
const detailsNode = document.getElementById("details");
const playerCardNode = document.getElementById("player-card");
const npcCardNode = document.getElementById("npc-card");
const resetButton = document.getElementById("reset");

let latestSnapshot = null;
let eventLog = ["Awaiting first world snapshot..."];
let previousAutonomousById = new Map();
let previousPlayerPosition = null;
let recentEffects = [];
let lastInteractionSignature = null;
let previousCamera = null;
let cachedMinimapAutonomousCircles = [];
let cachedMinimapFoods = [];
let cachedLocalFoods = [];
const pressedKeys = new Set();
let activeSocket = null;
let senderIntervalId = null;
let renderPressure = {
  samples: 0,
  avgMs: 0,
  maxMs: 0,
};

const movementKeys = new Set(["arrowleft", "arrowright", "arrowup", "arrowdown", "w", "a", "s", "d"]);
const CROWDING_RADIUS = 120;
const CROWDING_THRESHOLD = 2;
const CROWDING_CUE_DISTANCE = 260;
const FOOD_OPPORTUNITY_RADIUS = 170;
const FOOD_CUE_DISTANCE = 260;
const SCARCITY_THRESHOLD = 1;
const INTENT_CUE_DISTANCE = 260;
const MIN_MOVEMENT_FOR_INTENT = 1.5;
const AFTERGLOW_TTL = 10;
const VIEWPORT_MIN_WIDTH = 640;
const VIEWPORT_MIN_HEIGHT = 420;
const VIEWPORT_BOTTOM_MARGIN = 24;
const CAMERA_DEADZONE_X_RATIO = 0.22;
const CAMERA_DEADZONE_Y_RATIO = 0.2;
const CAMERA_LOOKAHEAD_X_RATIO = 0.1;
const CAMERA_LOOKAHEAD_Y_RATIO = 0.08;
const OFFSCREEN_AWARENESS_DISTANCE = 260;
const OFFSCREEN_EDGE_INSET = 18;
const NAME_ADJECTIVES = ["brave", "calm", "eager", "gentle", "keen", "lucky", "mellow", "nimble", "quiet", "solar", "swift", "vivid"];
const NAME_NOUNS = ["badger", "comet", "falcon", "harbor", "lantern", "meadow", "otter", "panda", "reef", "sable", "thunder", "willow"];

function normalizeKey(key) {
  return key.toLowerCase();
}

function hashString(value) {
  let hash = 0;
  for (let index = 0; index < value.length; index += 1) {
    hash = (hash * 31 + value.charCodeAt(index)) >>> 0;
  }

  return hash;
}

function displayName(id) {
  const hash = hashString(id);
  const adjective = NAME_ADJECTIVES[hash % NAME_ADJECTIVES.length];
  const noun = NAME_NOUNS[Math.floor(hash / NAME_ADJECTIVES.length) % NAME_NOUNS.length];
  return `${adjective}_${noun}`;
}

function childCount(circle) {
  return circle.attached_children.length;
}

function isObserverTransport(snapshot) {
  return snapshot.transport_mode === "observer_orientation_only";
}

function minimapAutonomousCircles(snapshot) {
  return snapshot.minimap_autonomous_circles ?? cachedMinimapAutonomousCircles ?? snapshot.autonomous_circles;
}

function minimapFoods(snapshot) {
  return snapshot.minimap_foods ?? cachedMinimapFoods ?? snapshot.foods;
}

function localFoods(snapshot) {
  return snapshot.foods ?? cachedLocalFoods;
}

function distanceBetween(a, b) {
  return Math.hypot(a.x - b.x, a.y - b.y);
}

function normalizedVector(x, y) {
  const magnitude = Math.hypot(x, y);
  if (magnitude === 0) {
    return null;
  }

  return { x: x / magnitude, y: y / magnitude, magnitude };
}

function dot(a, b) {
  return a.x * b.x + a.y * b.y;
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

function nearestByDistance(origin, items) {
  let nearest = null;
  let nearestDistance = Number.POSITIVE_INFINITY;

  for (const item of items) {
    const currentDistance = distanceBetween(origin, item);
    if (currentDistance < nearestDistance) {
      nearest = item;
      nearestDistance = currentDistance;
    }
  }

  return nearest;
}

function inferAutonomyIntent(circle, circles, foods, player) {
  if (!player || distanceBetween(circle, player) > INTENT_CUE_DISTANCE) {
    return null;
  }

  const previous = previousAutonomousById.get(circle.id);
  if (!previous) {
    return null;
  }

  const movement = normalizedVector(circle.x - previous.x, circle.y - previous.y);
  if (!movement || movement.magnitude < MIN_MOVEMENT_FOR_INTENT) {
    return null;
  }

  const visibleFoods = foods.filter((food) => distanceBetween(circle, food) <= FOOD_OPPORTUNITY_RADIUS);
  const nearestFood = nearestByDistance(circle, visibleFoods);
  if (nearestFood) {
    const toFood = normalizedVector(nearestFood.x - circle.x, nearestFood.y - circle.y);
    if (toFood && dot(movement, toFood) >= 0.7) {
      return "food";
    }
  }

  const candidateCircles = circles.filter((other) => other.id !== circle.id && distanceBetween(circle, other) <= FOOD_CUE_DISTANCE);
  const nearestCircle = nearestByDistance(circle, candidateCircles);
  if (!nearestCircle) {
    return null;
  }

  const toCircle = normalizedVector(nearestCircle.x - circle.x, nearestCircle.y - circle.y);
  if (!toCircle) {
    return null;
  }

  const alignment = dot(movement, toCircle);
  if (nearestCircle.shape === circle.shape && alignment <= -0.55) {
    return "retreat";
  }
  if (alignment >= 0.7) {
    return "social";
  }

  return null;
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

function effectStyle(kind) {
  if (kind === "death_promoted_child") {
    return {
      stroke: "rgba(255, 238, 163, 0.95)",
      glow: "rgba(255, 238, 163, 0.22)",
    };
  }

  if (kind.startsWith("reproduce")) {
    return {
      stroke: "rgba(117, 229, 149, 0.9)",
      glow: "rgba(117, 229, 149, 0.2)",
    };
  }

  return {
    stroke: "rgba(255, 132, 132, 0.92)",
    glow: "rgba(255, 132, 132, 0.18)",
  };
}

function entityPositionById(snapshot, id) {
  if (!id) {
    return null;
  }

  if (snapshot.player?.id === id) {
    return { x: snapshot.player.x, y: snapshot.player.y };
  }

  const autonomous = snapshot.autonomous_circles.find((circle) => circle.id === id);
  if (autonomous) {
    return { x: autonomous.x, y: autonomous.y };
  }

  return null;
}

function interactionAnchor(snapshot, interaction) {
  const source = entityPositionById(snapshot, interaction.source_id);
  const target = entityPositionById(snapshot, interaction.target_id);

  if (source && target) {
    return {
      x: (source.x + target.x) / 2,
      y: (source.y + target.y) / 2,
    };
  }

  return source || target;
}

function interactionSignature(interaction) {
  if (!interaction) {
    return null;
  }

  return [
    interaction.kind,
    interaction.source_id ?? "",
    interaction.target_id ?? "",
    interaction.promoted_child_id ?? "",
    interaction.absorbed_child_id ?? "",
  ].join("|");
}

function trackRecentInteraction(snapshot) {
  const interaction = snapshot.interaction;
  const signature = interactionSignature(interaction);

  if (!interaction || !interaction.kind || signature === lastInteractionSignature) {
    lastInteractionSignature = signature;
    return;
  }

  const anchoredKinds = new Set([
    "fight_absorbed_child",
    "fight_resolved",
    "reproduce_resolved",
    "reproduce_paid_child",
    "death_promoted_child",
  ]);

  if (!anchoredKinds.has(interaction.kind)) {
    lastInteractionSignature = signature;
    return;
  }

  const anchor = interactionAnchor(snapshot, interaction);
  if (!anchor) {
    lastInteractionSignature = signature;
    return;
  }

  recentEffects = [
    {
      ...effectStyle(interaction.kind),
      x: anchor.x,
      y: anchor.y,
      ttl: AFTERGLOW_TTL,
      kind: interaction.kind,
    },
    ...recentEffects,
  ].slice(0, 6);

  lastInteractionSignature = signature;
}

function hasContinuityReserve(circle) {
  return childCount(circle) > 0;
}

function isPromotedContinuity(circle, interaction) {
  if (!interaction || interaction.kind !== "death_promoted_child") {
    return false;
  }

  return interaction.source_id === circle.id || interaction.target_id === circle.id;
}

function renderEventLog() {
  detailsNode.innerHTML = eventLog.map((entry) => `<li>${entry}</li>`).join("");
}

function renderPlayerCard(player, pressure, foodState) {
  if (player === "observer") {
    playerCardNode.innerHTML = `
      <div class="player-stat">
        <div class="player-identity">
          <span class="player-name">Observer</span>
          <span class="player-state">orientation-only</span>
        </div>
      </div>
    `;
    return;
  }

  if (!player) {
    playerCardNode.innerHTML = `
      <div class="player-stat">
        <div class="player-identity">
          <span class="player-name">Defeated</span>
          <span class="player-state">offline</span>
        </div>
      </div>
    `;
    return;
  }

  const state = pressure || foodState ? `${pressure}${pressure && foodState ? " · " : ""}${foodState}` : "stable";

  playerCardNode.innerHTML = `
    <div class="player-stat">
      <div class="player-identity">
        <span class="player-name">${displayName(player.id)}</span>
        <span class="player-state">${state}</span>
      </div>
      <div class="player-meta">
        <span class="player-meta-item"><strong>${player.shape}</strong> shape</span>
        <span class="player-meta-item">E <strong>${player.energy.toFixed(0)}</strong></span>
        <span class="player-meta-item">C <strong>${childCount(player)}</strong></span>
        <span class="player-meta-item">G <strong>${player.generation}</strong></span>
      </div>
    </div>
  `;
}

function renderNpcCard(circles) {
  if (!circles.length) {
    npcCardNode.innerHTML = "<li>No active NPCs.</li>";
    return;
  }

  npcCardNode.innerHTML = circles.map((circle) => `
    <li>${displayName(circle.id)} · ${circle.shape} · e:${circle.energy.toFixed(0)} c:${childCount(circle)} g:${circle.generation}</li>
  `).join("");
}

function recordRenderPressure(durationMs) {
  const clampedDuration = Math.max(0, durationMs);
  const windowSize = 30;
  const previousSamples = renderPressure.samples;
  const nextSamples = Math.min(windowSize, previousSamples + 1);
  const carriedWeight = nextSamples === windowSize ? windowSize - 1 : previousSamples;
  const nextAverage = ((renderPressure.avgMs * carriedWeight) + clampedDuration) / nextSamples;
  const nextMax = previousSamples >= windowSize
    ? Math.max(clampedDuration, renderPressure.maxMs * 0.92)
    : Math.max(renderPressure.maxMs, clampedDuration);

  renderPressure = {
    samples: nextSamples,
    avgMs: nextAverage,
    maxMs: nextMax,
  };
  renderPressureNode.textContent = `Render ${nextAverage.toFixed(1)}ms · max ${nextMax.toFixed(1)}ms`;
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

function contentBoxSize(node) {
  const rect = node.getBoundingClientRect();
  const style = window.getComputedStyle(node);
  const horizontalPadding = parseFloat(style.paddingLeft) + parseFloat(style.paddingRight);
  const verticalPadding = parseFloat(style.paddingTop) + parseFloat(style.paddingBottom);

  return {
    width: Math.max(0, Math.floor(rect.width - horizontalPadding)),
    height: Math.max(0, Math.floor(rect.height - verticalPadding)),
  };
}

function viewportSize(world) {
  const stageSize = contentBoxSize(playStageNode);
  const availableWidth = Math.max(VIEWPORT_MIN_WIDTH, stageSize.width);
  const stageRect = playStageNode.getBoundingClientRect();
  const availableHeight = Math.max(
    VIEWPORT_MIN_HEIGHT,
    Math.floor(window.innerHeight - stageRect.top - VIEWPORT_BOTTOM_MARGIN),
  );

  return {
    width: Math.min(world.width, availableWidth),
    height: Math.min(world.height, availableHeight),
  };
}

function cameraFocus(snapshot) {
  if (snapshot.player && previousPlayerPosition) {
    const movement = normalizedVector(
      snapshot.player.x - previousPlayerPosition.x,
      snapshot.player.y - previousPlayerPosition.y,
    );
    if (movement && movement.magnitude >= MIN_MOVEMENT_FOR_INTENT) {
      const viewport = viewportSize(snapshot.world);
      return {
        x: snapshot.player.x + movement.x * Math.floor(viewport.width * CAMERA_LOOKAHEAD_X_RATIO),
        y: snapshot.player.y + movement.y * Math.floor(viewport.height * CAMERA_LOOKAHEAD_Y_RATIO),
      };
    }
  }

  if (snapshot.player) {
    return { x: snapshot.player.x, y: snapshot.player.y };
  }

  if (snapshot.autonomous_circles.length > 0) {
    return {
      x: snapshot.autonomous_circles[0].x,
      y: snapshot.autonomous_circles[0].y,
    };
  }

  return {
    x: snapshot.world.width / 2,
    y: snapshot.world.height / 2,
  };
}

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

function cameraRect(snapshot, viewport) {
  const focus = cameraFocus(snapshot);
  const maxX = Math.max(0, snapshot.world.width - viewport.width);
  const maxY = Math.max(0, snapshot.world.height - viewport.height);
  const deadzoneHalfWidth = Math.floor(viewport.width * CAMERA_DEADZONE_X_RATIO);
  const deadzoneHalfHeight = Math.floor(viewport.height * CAMERA_DEADZONE_Y_RATIO);
  const baseCamera = previousCamera && previousCamera.width === viewport.width && previousCamera.height === viewport.height
    ? previousCamera
    : {
        x: clamp(Math.round(focus.x - viewport.width / 2), 0, maxX),
        y: clamp(Math.round(focus.y - viewport.height / 2), 0, maxY),
        width: viewport.width,
        height: viewport.height,
      };

  let cameraX = baseCamera.x;
  let cameraY = baseCamera.y;

  const deadzoneLeft = cameraX + viewport.width / 2 - deadzoneHalfWidth;
  const deadzoneRight = cameraX + viewport.width / 2 + deadzoneHalfWidth;
  const deadzoneTop = cameraY + viewport.height / 2 - deadzoneHalfHeight;
  const deadzoneBottom = cameraY + viewport.height / 2 + deadzoneHalfHeight;

  if (focus.x < deadzoneLeft) {
    cameraX -= Math.round(deadzoneLeft - focus.x);
  } else if (focus.x > deadzoneRight) {
    cameraX += Math.round(focus.x - deadzoneRight);
  }

  if (focus.y < deadzoneTop) {
    cameraY -= Math.round(deadzoneTop - focus.y);
  } else if (focus.y > deadzoneBottom) {
    cameraY += Math.round(focus.y - deadzoneBottom);
  }

  return {
    x: clamp(cameraX, 0, maxX),
    y: clamp(cameraY, 0, maxY),
    width: viewport.width,
    height: viewport.height,
  };
}

function draw(snapshot) {
  const drawStartedAt = performance.now();
  const observerTransport = isObserverTransport(snapshot);
  const foods = localFoods(snapshot);
  const circles = snapshot.player ? [snapshot.player, ...snapshot.autonomous_circles] : [...snapshot.autonomous_circles];
  const playerFoodPressure = snapshot.player ? foodPressureAt(snapshot.player, foods) : null;
  trackRecentInteraction(snapshot);

  const viewport = viewportSize(snapshot.world);
  const camera = cameraRect(snapshot, viewport);
  previousCamera = camera;

  canvas.width = viewport.width;
  canvas.height = viewport.height;
  canvas.style.width = `${viewport.width}px`;
  canvas.style.height = `${viewport.height}px`;

  context.clearRect(0, 0, canvas.width, canvas.height);
  context.save();
  context.translate(-camera.x, -camera.y);

  context.fillStyle = "rgba(6, 22, 30, 0.9)";
  context.fillRect(camera.x, camera.y, viewport.width, viewport.height);

  context.strokeStyle = "rgba(127, 174, 188, 0.28)";
  context.lineWidth = 2;
  context.strokeRect(1, 1, snapshot.world.width - 2, snapshot.world.height - 2);

  drawRecentEffects();
  drawCrowdingZones(circles, snapshot.player);
  drawFoodZones(foods, snapshot.player);
  drawLineageLinks(circles, snapshot.interaction);

  context.fillStyle = "#ff8a5b";
  for (const food of foods) {
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
    drawCircle(circle, false, snapshot.player, circles, foods);
  }

  if (snapshot.player) {
    drawCircle(snapshot.player, true, snapshot.player, circles, foods);
    drawPlayerHeadingCue(snapshot.player);
    const pressure = isCrowded(snapshot.player, circles) ? "pressure" : "";
    const foodState = playerFoodPressure?.scarcity ? "scarce" : playerFoodPressure?.opportunity ? "food-rich" : "";
    energyNode.textContent = `${displayName(snapshot.player.id)}`;
    renderPlayerCard(snapshot.player, pressure, foodState);
  } else {
    if (observerTransport) {
      energyNode.textContent = "observer";
      renderPlayerCard("observer", "", "");
    } else {
      energyNode.textContent = "defeated";
      renderPlayerCard(null, "", "");
    }
  }

  renderNpcCard(snapshot.autonomous_circles);

  const totalFoods = snapshot.total_foods ?? minimapFoods(snapshot).length;
  const totalAutonomousCircles = snapshot.total_autonomous_circles ?? minimapAutonomousCircles(snapshot).length;
  tickNode.textContent = `${snapshot.tick} · ${totalFoods}f · ${totalAutonomousCircles}o`;
  pushEventLog(interactionSummary(snapshot.interaction));
  previousAutonomousById = new Map(snapshot.autonomous_circles.map((circle) => [circle.id, { x: circle.x, y: circle.y }]));
  previousPlayerPosition = snapshot.player ? { x: snapshot.player.x, y: snapshot.player.y } : null;
  recentEffects = recentEffects
    .map((effect) => ({ ...effect, ttl: effect.ttl - 1 }))
    .filter((effect) => effect.ttl > 0);
  drawOffscreenFoodAwareness(snapshot, camera);
  drawOffscreenAwareness(snapshot, camera);
  context.restore();
  drawMinimap(snapshot, camera);
  recordRenderPressure(performance.now() - drawStartedAt);
}

function drawPlayerHeadingCue(player) {
  if (!previousPlayerPosition) {
    return;
  }

  const movement = normalizedVector(player.x - previousPlayerPosition.x, player.y - previousPlayerPosition.y);
  if (!movement || movement.magnitude < MIN_MOVEMENT_FOR_INTENT) {
    return;
  }

  const tipX = player.x + movement.x * (player.radius + 18);
  const tipY = player.y + movement.y * (player.radius + 18);
  const wingX = movement.y * 5;
  const wingY = -movement.x * 5;

  context.strokeStyle = "rgba(255, 224, 160, 0.96)";
  context.lineWidth = 3;
  context.beginPath();
  context.moveTo(player.x, player.y);
  context.lineTo(tipX, tipY);
  context.moveTo(tipX, tipY);
  context.lineTo(tipX - movement.x * 8 + wingX, tipY - movement.y * 8 + wingY);
  context.moveTo(tipX, tipY);
  context.lineTo(tipX - movement.x * 8 - wingX, tipY - movement.y * 8 - wingY);
  context.stroke();
}

function drawOffscreenAwareness(snapshot, camera) {
  if (!snapshot.player) {
    return;
  }

  const viewportLeft = camera.x;
  const viewportRight = camera.x + camera.width;
  const viewportTop = camera.y;
  const viewportBottom = camera.y + camera.height;
  const nearby = snapshot.autonomous_circles.filter((circle) => {
    const distance = distanceBetween(circle, snapshot.player);
    if (distance > OFFSCREEN_AWARENESS_DISTANCE) {
      return false;
    }

    return circle.x < viewportLeft || circle.x > viewportRight || circle.y < viewportTop || circle.y > viewportBottom;
  }).slice(0, 6);

  for (const circle of nearby) {
    const relation = playerRiskState(circle, snapshot.player, snapshot.interaction);
    const relativeX = clamp(circle.x, viewportLeft + OFFSCREEN_EDGE_INSET, viewportRight - OFFSCREEN_EDGE_INSET) - camera.x;
    const relativeY = clamp(circle.y, viewportTop + OFFSCREEN_EDGE_INSET, viewportBottom - OFFSCREEN_EDGE_INSET) - camera.y;

    context.save();
    context.translate(camera.x, camera.y);
    context.strokeStyle = relation === "danger"
      ? "rgba(255, 95, 109, 0.92)"
      : relation === "blocked"
        ? "rgba(255, 201, 92, 0.92)"
        : "rgba(117, 229, 149, 0.9)";
    context.lineWidth = 3;
    context.beginPath();
    context.arc(relativeX, relativeY, 8, 0, Math.PI * 2);
    context.stroke();
    context.restore();
  }
}

function drawOffscreenFoodAwareness(snapshot, camera) {
  if (!snapshot.player) {
    return;
  }

  const viewportLeft = camera.x;
  const viewportRight = camera.x + camera.width;
  const viewportTop = camera.y;
  const viewportBottom = camera.y + camera.height;
  const nearbyFoods = localFoods(snapshot).filter((food) => {
    const distance = distanceBetween(food, snapshot.player);
    if (distance > OFFSCREEN_AWARENESS_DISTANCE) {
      return false;
    }

    return food.x < viewportLeft || food.x > viewportRight || food.y < viewportTop || food.y > viewportBottom;
  }).slice(0, 6);

  for (const food of nearbyFoods) {
    const relativeX = clamp(food.x, viewportLeft + OFFSCREEN_EDGE_INSET, viewportRight - OFFSCREEN_EDGE_INSET) - camera.x;
    const relativeY = clamp(food.y, viewportTop + OFFSCREEN_EDGE_INSET, viewportBottom - OFFSCREEN_EDGE_INSET) - camera.y;

    context.save();
    context.translate(camera.x, camera.y);
    context.fillStyle = "rgba(255, 138, 91, 0.88)";
    context.beginPath();
    context.arc(relativeX, relativeY, 4, 0, Math.PI * 2);
    context.fill();
    context.restore();
  }
}

function drawMinimap(snapshot, camera) {
  const width = minimapCanvas.width;
  const height = minimapCanvas.height;
  const scaleX = width / snapshot.world.width;
  const scaleY = height / snapshot.world.height;
  const circles = minimapAutonomousCircles(snapshot);
  const foods = minimapFoods(snapshot);

  minimapContext.clearRect(0, 0, width, height);
  minimapContext.fillStyle = "rgba(6, 18, 24, 0.96)";
  minimapContext.fillRect(0, 0, width, height);

  minimapContext.strokeStyle = "rgba(126, 166, 178, 0.32)";
  minimapContext.lineWidth = 1;
  minimapContext.strokeRect(0.5, 0.5, width - 1, height - 1);

  minimapContext.fillStyle = "rgba(255, 138, 91, 0.82)";
  for (const food of foods) {
    const size = Math.min(5, 1 + (food.count ?? 1));
    minimapContext.fillRect(
      Math.round(food.x * scaleX) - Math.floor(size / 2),
      Math.round(food.y * scaleY) - Math.floor(size / 2),
      size,
      size,
    );
  }

  for (const circle of circles) {
    minimapContext.fillStyle = circle.shape === "triangle" ? "#6fd5ff" : "#c08cff";
    minimapContext.beginPath();
    minimapContext.arc(circle.x * scaleX, circle.y * scaleY, Math.min(5, 1.5 + (circle.count ?? 1)), 0, Math.PI * 2);
    minimapContext.fill();
  }

  if (snapshot.player) {
    minimapContext.fillStyle = "#ff8a5b";
    minimapContext.beginPath();
    minimapContext.arc(snapshot.player.x * scaleX, snapshot.player.y * scaleY, 3, 0, Math.PI * 2);
    minimapContext.fill();
  }

  minimapContext.strokeStyle = "rgba(255, 240, 170, 0.92)";
  minimapContext.lineWidth = 1.5;
  minimapContext.strokeRect(
    Math.round(camera.x * scaleX) + 0.5,
    Math.round(camera.y * scaleY) + 0.5,
    Math.max(1, Math.round(camera.width * scaleX) - 1),
    Math.max(1, Math.round(camera.height * scaleY) - 1),
  );
}

function drawRecentEffects() {
  for (const effect of recentEffects) {
    const progress = effect.ttl / AFTERGLOW_TTL;
    const radius = 24 + (1 - progress) * 18;
    const gradient = context.createRadialGradient(effect.x, effect.y, 6, effect.x, effect.y, radius);
    gradient.addColorStop(0, effect.glow);
    gradient.addColorStop(1, "rgba(0, 0, 0, 0)");
    context.fillStyle = gradient;
    context.beginPath();
    context.arc(effect.x, effect.y, radius, 0, Math.PI * 2);
    context.fill();

    context.strokeStyle = effect.stroke;
    context.globalAlpha = progress;
    context.lineWidth = 2.5;
    context.beginPath();
    context.arc(effect.x, effect.y, radius - 6, 0, Math.PI * 2);
    context.stroke();
    context.globalAlpha = 1;
  }
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

function drawIntentCue(circle, intent) {
  const previous = previousAutonomousById.get(circle.id);
  if (!previous) {
    return;
  }

  const movement = normalizedVector(circle.x - previous.x, circle.y - previous.y);
  if (!movement || movement.magnitude < MIN_MOVEMENT_FOR_INTENT) {
    return;
  }

  const cueX = circle.x + movement.x * (circle.radius + 18);
  const cueY = circle.y + movement.y * (circle.radius + 18);

  if (intent === "food") {
    context.fillStyle = "rgba(94, 224, 138, 0.95)";
    context.beginPath();
    context.arc(cueX, cueY, 5, 0, Math.PI * 2);
    context.fill();
    return;
  }

  if (intent === "social") {
    context.strokeStyle = "rgba(111, 213, 255, 0.95)";
    context.lineWidth = 3;
    context.beginPath();
    context.moveTo(circle.x, circle.y);
    context.lineTo(cueX, cueY);
    context.stroke();
    return;
  }

  if (intent === "retreat") {
    context.strokeStyle = "rgba(255, 130, 130, 0.95)";
    context.lineWidth = 3;
    context.beginPath();
    context.moveTo(cueX, cueY);
    context.lineTo(cueX - movement.x * 8 + movement.y * 5, cueY - movement.y * 8 - movement.x * 5);
    context.moveTo(cueX, cueY);
    context.lineTo(cueX - movement.x * 8 - movement.y * 5, cueY - movement.y * 8 + movement.x * 5);
    context.stroke();
  }
}

function drawLineageLinks(circles, interaction) {
  for (const circle of circles) {
    if (!hasContinuityReserve(circle)) {
      continue;
    }

    const promoted = isPromotedContinuity(circle, interaction);
    for (const child of circle.attached_children) {
      context.strokeStyle = promoted ? "rgba(255, 236, 158, 0.95)" : "rgba(255, 220, 122, 0.5)";
      context.lineWidth = promoted ? 2.5 : 1.5;
      context.beginPath();
      context.moveTo(circle.x, circle.y);
      context.lineTo(child.x, child.y);
      context.stroke();
    }
  }
}

function drawCircle(circle, isPlayer, player, circles, foods) {
  const matchesPlayerShape = !isPlayer && player && circle.shape === player.shape;
  const relationToPlayer = playerRiskState(circle, player, latestSnapshot?.interaction);
  const crowded = shouldRenderCrowdingCue(circle, circles, player);
  const intent = !isPlayer ? inferAutonomyIntent(circle, circles, foods, player) : null;
  const continuityReserve = hasContinuityReserve(circle);
  const promoted = isPromotedContinuity(circle, latestSnapshot?.interaction);
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

  if (continuityReserve) {
    context.strokeStyle = promoted ? "rgba(255, 240, 170, 0.95)" : "rgba(255, 226, 140, 0.7)";
    context.lineWidth = promoted ? 4 : 2;
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius + 18, 0, Math.PI * 2);
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

  if (intent) {
    drawIntentCue(circle, intent);
  }

  drawAttachedChildren(circle, color);

  context.fillStyle = "#e4f3f8";
  context.font = "16px Georgia";
  const children = childCount(circle);
  const label = isPlayer
    ? displayName(circle.id)
    : displayName(circle.id);
  context.fillText(label, circle.x - 28, circle.y - circle.radius - 10);
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
    if (snapshot.foods_fresh ?? true) {
      cachedLocalFoods = snapshot.foods ?? [];
    } else {
      snapshot.foods = cachedLocalFoods;
    }
    if (snapshot.orientation_fresh) {
      cachedMinimapAutonomousCircles = snapshot.minimap_autonomous_circles ?? [];
      cachedMinimapFoods = snapshot.minimap_foods ?? [];
    } else {
      snapshot.minimap_autonomous_circles = cachedMinimapAutonomousCircles;
      snapshot.minimap_foods = cachedMinimapFoods;
    }
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

    if (snapshot.orientation_fresh) {
      cachedMinimapAutonomousCircles = snapshot.minimap_autonomous_circles ?? [];
      cachedMinimapFoods = snapshot.minimap_foods ?? [];
    } else {
      snapshot.minimap_autonomous_circles = cachedMinimapAutonomousCircles;
      snapshot.minimap_foods = cachedMinimapFoods;
    }
    if (isObserverTransport(snapshot)) {
      snapshot.foods = [];
    } else if (snapshot.foods_fresh ?? true) {
      cachedLocalFoods = snapshot.foods ?? [];
    } else {
      snapshot.foods = cachedLocalFoods;
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

window.addEventListener("resize", () => {
  if (latestSnapshot) {
    draw(latestSnapshot);
  }
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
