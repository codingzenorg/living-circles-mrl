package simulation

import "math"

const (
	DefaultWorldWidth   = 800.0
	DefaultWorldHeight  = 600.0
	DefaultPlayerRadius = 12.0
	DefaultPlayerEnergy = 100.0
	DefaultMoveSpeed    = 8.0
	DefaultMoveCost     = 1.0
)

type Bounds struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PlayerCircle struct {
	ID     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
	Energy float64 `json:"energy"`
}

type Snapshot struct {
	Type   string       `json:"type"`
	Tick   int64        `json:"tick"`
	World  Bounds       `json:"world"`
	Player PlayerCircle `json:"player"`
}

type World struct {
	bounds   Bounds
	player   PlayerCircle
	moveCost float64
	speed    float64
}

func NewWorld() *World {
	return &World{
		bounds: Bounds{
			Width:  DefaultWorldWidth,
			Height: DefaultWorldHeight,
		},
		player: PlayerCircle{
			ID:     "player-1",
			X:      DefaultWorldWidth / 2,
			Y:      DefaultWorldHeight / 2,
			Radius: DefaultPlayerRadius,
			Energy: DefaultPlayerEnergy,
		},
		moveCost: DefaultMoveCost,
		speed:    DefaultMoveSpeed,
	}
}

func (w *World) Advance(tick int64, intent Vector) Snapshot {
	if w.player.Energy > 0 {
		normalized := normalize(intent)
		if normalized.X != 0 || normalized.Y != 0 {
			w.player.X = clamp(w.player.X+normalized.X*w.speed, w.player.Radius, w.bounds.Width-w.player.Radius)
			w.player.Y = clamp(w.player.Y+normalized.Y*w.speed, w.player.Radius, w.bounds.Height-w.player.Radius)
			w.player.Energy = math.Max(0, w.player.Energy-w.moveCost)
		}
	}

	return w.Snapshot(tick)
}

func (w *World) Snapshot(tick int64) Snapshot {
	return Snapshot{
		Type:   "world_snapshot",
		Tick:   tick,
		World:  w.bounds,
		Player: w.player,
	}
}

func normalize(vector Vector) Vector {
	length := math.Hypot(vector.X, vector.Y)
	if length == 0 {
		return Vector{}
	}

	return Vector{
		X: vector.X / length,
		Y: vector.Y / length,
	}
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Min(maximum, math.Max(minimum, value))
}
