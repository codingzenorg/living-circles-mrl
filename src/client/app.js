import { CONTRACT_VERSION, createMovementIntent, MESSAGE_TYPES } from "/shared_contracts/messages/protocol.js";

const canvas = document.getElementById("world");
const context = canvas.getContext("2d");
const statusNode = document.getElementById("status");
const energyNode = document.getElementById("energy");
const tickNode = document.getElementById("tick");

let latestSnapshot = null;
const pressedKeys = new Set();

function currentDirection() {
  const vector = { x: 0, y: 0 };

  if (pressedKeys.has("ArrowLeft") || pressedKeys.has("a")) {
    vector.x -= 1;
  }
  if (pressedKeys.has("ArrowRight") || pressedKeys.has("d")) {
    vector.x += 1;
  }
  if (pressedKeys.has("ArrowUp") || pressedKeys.has("w")) {
    vector.y -= 1;
  }
  if (pressedKeys.has("ArrowDown") || pressedKeys.has("s")) {
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

  context.fillStyle = "#3b8ea5";
  context.beginPath();
  context.arc(snapshot.player.x, snapshot.player.y, snapshot.player.radius, 0, Math.PI * 2);
  context.fill();

  context.fillStyle = "#17313a";
  context.font = "16px Georgia";
  context.fillText(snapshot.player.id, snapshot.player.x - 28, snapshot.player.y - snapshot.player.radius - 10);

  energyNode.textContent = `Energy: ${snapshot.player.energy.toFixed(0)}`;
  tickNode.textContent = `Tick: ${snapshot.tick} · Food: ${snapshot.foods.length}`;
}

function setStatus(message) {
  statusNode.textContent = `${message} · contract v${CONTRACT_VERSION}`;
}

function connect() {
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const socket = new WebSocket(`${protocol}://${window.location.host}/ws`);

  socket.addEventListener("open", () => {
    setStatus("Connected");
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
    setStatus("Disconnected");
    setTimeout(connect, 1000);
  });

  socket.addEventListener("error", () => {
    setStatus("Connection error");
  });

  setInterval(() => {
    if (socket.readyState !== WebSocket.OPEN) {
      return;
    }

    const direction = currentDirection();
    socket.send(JSON.stringify(createMovementIntent(direction.x, direction.y)));
  }, 100);
}

window.addEventListener("keydown", (event) => {
  pressedKeys.add(event.key);
});

window.addEventListener("keyup", (event) => {
  pressedKeys.delete(event.key);
});

setStatus("Connecting");
connect();

if (!latestSnapshot) {
  context.fillStyle = "#17313a";
  context.font = "24px Georgia";
  context.fillText("Waiting for server snapshot...", 24, 48);
}
