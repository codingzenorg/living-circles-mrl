import { CONTRACT_VERSION, createMovementIntent, MESSAGE_TYPES } from "/shared_contracts/messages/protocol.js";

const canvas = document.getElementById("world");
const context = canvas.getContext("2d");
const statusNode = document.getElementById("status");
const energyNode = document.getElementById("energy");
const tickNode = document.getElementById("tick");

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

  context.fillStyle = "#8c6bb1";
  for (const circle of snapshot.autonomous_circles) {
    context.beginPath();
    context.arc(circle.x, circle.y, circle.radius, 0, Math.PI * 2);
    context.fill();
  }

  context.fillStyle = "#3b8ea5";
  context.beginPath();
  context.arc(snapshot.player.x, snapshot.player.y, snapshot.player.radius, 0, Math.PI * 2);
  context.fill();

  context.fillStyle = "#17313a";
  context.font = "16px Georgia";
  context.fillText(snapshot.player.id, snapshot.player.x - 28, snapshot.player.y - snapshot.player.radius - 10);

  energyNode.textContent = `Energy: ${snapshot.player.energy.toFixed(0)}`;
  tickNode.textContent = `Tick: ${snapshot.tick} · Food: ${snapshot.foods.length} · Others: ${snapshot.autonomous_circles.length}`;
}

function setStatus(message) {
  statusNode.textContent = `${message} · contract v${CONTRACT_VERSION}`;
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

setStatus("Connecting");
connect();

if (!latestSnapshot) {
  context.fillStyle = "#17313a";
  context.font = "24px Georgia";
  context.fillText("Waiting for server snapshot...", 24, 48);
}
