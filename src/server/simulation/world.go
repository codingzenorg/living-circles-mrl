package simulation

import "math"

const (
	DefaultWorldWidth   = 800.0
	DefaultWorldHeight  = 600.0
	DefaultPlayerRadius = 12.0
	DefaultAutonomousID = "circle-2"
	DefaultPlayerEnergy = 100.0
	DefaultMaxEnergy    = 100.0
	DefaultMoveSpeed    = 8.0
	DefaultMoveCost     = 1.0
	DefaultFoodRadius   = 6.0
	DefaultFoodEnergy   = 10.0
	DefaultPlayerShape  = "triangle"
	DefaultAutoShape    = "square"
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
	Shape  string  `json:"shape"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
	Energy float64 `json:"energy"`
}

type AutonomousCircle struct {
	ID     string  `json:"id"`
	Shape  string  `json:"shape"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
	Energy float64 `json:"energy"`
}

type InteractionClassification struct {
	Active   bool   `json:"active"`
	Resolved bool   `json:"resolved"`
	Kind     string `json:"kind"`
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	WinnerID string `json:"winner_id"`
	LoserID  string `json:"loser_id"`
}

type Food struct {
	ID     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

type Snapshot struct {
	Type              string                     `json:"type"`
	Tick              int64                      `json:"tick"`
	World             Bounds                     `json:"world"`
	Player            *PlayerCircle              `json:"player"`
	AutonomousCircles []AutonomousCircle         `json:"autonomous_circles"`
	Interaction       *InteractionClassification `json:"interaction"`
	Foods             []Food                     `json:"foods"`
}

type World struct {
	bounds            Bounds
	player            *PlayerCircle
	autonomousCircles []AutonomousCircle
	foods             []Food
	moveCost          float64
	speed             float64
	maxEnergy         float64
	foodGain          float64
	lastInteraction   *InteractionClassification
}

type Config struct {
	PlayerShape      string
	AutonomousShape  string
	PlayerEnergy     float64
	AutonomousEnergy float64
}

func NewWorld() *World {
	return NewWorldWithConfig(Config{
		PlayerShape:      DefaultPlayerShape,
		AutonomousShape:  DefaultAutoShape,
		PlayerEnergy:     DefaultPlayerEnergy,
		AutonomousEnergy: DefaultPlayerEnergy,
	})
}

func NewWorldWithShapes(playerShape, autonomousShape string) *World {
	return NewWorldWithConfig(Config{
		PlayerShape:      playerShape,
		AutonomousShape:  autonomousShape,
		PlayerEnergy:     DefaultPlayerEnergy,
		AutonomousEnergy: DefaultPlayerEnergy,
	})
}

func NewWorldWithConfig(config Config) *World {
	return &World{
		bounds: Bounds{
			Width:  DefaultWorldWidth,
			Height: DefaultWorldHeight,
		},
		player: &PlayerCircle{
			ID:     "player-1",
			Shape:  config.PlayerShape,
			X:      DefaultWorldWidth / 2,
			Y:      DefaultWorldHeight / 2,
			Radius: DefaultPlayerRadius,
			Energy: config.PlayerEnergy,
		},
		autonomousCircles: []AutonomousCircle{
			{
				ID:     DefaultAutonomousID,
				Shape:  config.AutonomousShape,
				X:      DefaultWorldWidth/2 - 140,
				Y:      DefaultWorldHeight / 2,
				Radius: DefaultPlayerRadius,
				Energy: config.AutonomousEnergy,
			},
		},
		foods: []Food{
			{ID: "food-1", X: DefaultWorldWidth/2 + 32, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
			{ID: "food-2", X: DefaultWorldWidth/2 - 108, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
			{ID: "food-3", X: DefaultWorldWidth/2 + 120, Y: DefaultWorldHeight/2 + 84, Radius: DefaultFoodRadius},
		},
		moveCost:  DefaultMoveCost,
		speed:     DefaultMoveSpeed,
		maxEnergy: DefaultMaxEnergy,
		foodGain:  DefaultFoodEnergy,
	}
}

func (w *World) Advance(tick int64, intent Vector) Snapshot {
	w.lastInteraction = nil
	w.player = w.advanceCircle(w.player, intent)
	w.advanceAutonomousCircles(tick)

	w.consumeOverlappingFood()
	w.resolveCircleInteractions()

	return w.Snapshot(tick)
}

func (w *World) Snapshot(tick int64) Snapshot {
	var player *PlayerCircle
	if w.player != nil {
		copy := *w.player
		player = &copy
	}

	var interaction *InteractionClassification
	if w.lastInteraction != nil {
		copy := *w.lastInteraction
		interaction = &copy
	}

	return Snapshot{
		Type:              "world_snapshot",
		Tick:              tick,
		World:             w.bounds,
		Player:            player,
		AutonomousCircles: append([]AutonomousCircle(nil), w.autonomousCircles...),
		Interaction:       interaction,
		Foods:             append([]Food(nil), w.foods...),
	}
}

func (w *World) consumeOverlappingFood() {
	remaining := make([]Food, 0, len(w.foods))
	for _, food := range w.foods {
		if w.player != nil && overlaps(w.player.X, w.player.Y, w.player.Radius, food.X, food.Y, food.Radius) {
			w.player.Energy = math.Min(w.maxEnergy, w.player.Energy+w.foodGain)
			continue
		}

		consumed := false
		for index, circle := range w.autonomousCircles {
			if overlaps(circle.X, circle.Y, circle.Radius, food.X, food.Y, food.Radius) {
				circle.Energy = math.Min(w.maxEnergy, circle.Energy+w.foodGain)
				w.autonomousCircles[index] = circle
				consumed = true
				break
			}
		}

		if consumed {
			continue
		}

		remaining = append(remaining, food)
	}

	w.foods = remaining
}

func (w *World) advanceCircle(circle *PlayerCircle, intent Vector) *PlayerCircle {
	if circle == nil {
		return nil
	}
	if circle.Energy <= 0 {
		return circle
	}

	normalized := normalize(intent)
	if normalized.X == 0 && normalized.Y == 0 {
		return circle
	}

	circle.X = clamp(circle.X+normalized.X*w.speed, circle.Radius, w.bounds.Width-circle.Radius)
	circle.Y = clamp(circle.Y+normalized.Y*w.speed, circle.Radius, w.bounds.Height-circle.Radius)
	circle.Energy = math.Max(0, circle.Energy-w.moveCost)

	return circle
}

func (w *World) advanceAutonomousCircles(tick int64) {
	for index, circle := range w.autonomousCircles {
		if circle.Energy <= 0 {
			continue
		}

		intent := autonomousIntent(tick, index)
		normalized := normalize(intent)
		if normalized.X == 0 && normalized.Y == 0 {
			continue
		}

		circle.X = clamp(circle.X+normalized.X*w.speed, circle.Radius, w.bounds.Width-circle.Radius)
		circle.Y = clamp(circle.Y+normalized.Y*w.speed, circle.Radius, w.bounds.Height-circle.Radius)
		circle.Energy = math.Max(0, circle.Energy-w.moveCost)
		w.autonomousCircles[index] = circle
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

func overlaps(ax, ay, ar, bx, by, br float64) bool {
	return math.Hypot(ax-bx, ay-by) <= ar+br
}

func autonomousIntent(tick int64, index int) Vector {
	return Vector{X: 1, Y: 0}
}

func (w *World) resolveCircleInteractions() {
	if w.player == nil {
		return
	}

	for _, circle := range w.autonomousCircles {
		if !overlaps(w.player.X, w.player.Y, w.player.Radius, circle.X, circle.Y, circle.Radius) {
			continue
		}

		if w.player.Shape != circle.Shape {
			w.lastInteraction = &InteractionClassification{
				Active:   true,
				Resolved: false,
				Kind:     "reproduce_candidate",
				SourceID: w.player.ID,
				TargetID: circle.ID,
			}
			return
		}

		w.resolveFight(circle.ID)
		return
	}
}

func (w *World) resolveFight(opponentID string) {
	opponentIndex := -1
	for index, circle := range w.autonomousCircles {
		if circle.ID == opponentID {
			opponentIndex = index
			break
		}
	}

	if opponentIndex == -1 || w.player == nil {
		return
	}

	opponent := w.autonomousCircles[opponentIndex]
	winnerID, loserID := determineFightOutcome(*w.player, opponent)

	w.lastInteraction = &InteractionClassification{
		Active:   false,
		Resolved: true,
		Kind:     "fight_resolved",
		SourceID: w.player.ID,
		TargetID: opponent.ID,
		WinnerID: winnerID,
		LoserID:  loserID,
	}

	if loserID == w.player.ID {
		w.player = nil
		return
	}

	w.autonomousCircles = append(w.autonomousCircles[:opponentIndex], w.autonomousCircles[opponentIndex+1:]...)
}

func determineFightOutcome(player PlayerCircle, opponent AutonomousCircle) (string, string) {
	if player.Energy > opponent.Energy {
		return player.ID, opponent.ID
	}
	if opponent.Energy > player.Energy {
		return opponent.ID, player.ID
	}
	if player.Radius > opponent.Radius {
		return player.ID, opponent.ID
	}
	if opponent.Radius > player.Radius {
		return opponent.ID, player.ID
	}

	return player.ID, opponent.ID
}
