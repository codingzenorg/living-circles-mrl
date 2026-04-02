package simulation

import "math"

const (
	DefaultWorldWidth        = 800.0
	DefaultWorldHeight       = 600.0
	DefaultPlayerRadius      = 12.0
	DefaultChildRadiusGain   = 4.0
	DefaultAutonomousID      = "circle-2"
	DefaultSecondaryID       = "circle-3"
	DefaultPlayerEnergy      = 100.0
	DefaultReplacementEnergy = 100.0
	DefaultMaxEnergy         = 100.0
	DefaultMoveSpeed         = 8.0
	DefaultMoveCost          = 1.0
	DefaultFoodRadius        = 6.0
	DefaultFoodEnergy        = 10.0
	DefaultPlayerShape       = "triangle"
	DefaultAutoShape         = "square"
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
	ID            string  `json:"id"`
	Shape         string  `json:"shape"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Radius        float64 `json:"radius"`
	Energy        float64 `json:"energy"`
	ChildrenCount int     `json:"children_count"`
}

type AutonomousCircle struct {
	ID            string  `json:"id"`
	Shape         string  `json:"shape"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Radius        float64 `json:"radius"`
	Energy        float64 `json:"energy"`
	ChildrenCount int     `json:"children_count"`
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
	bounds               Bounds
	player               *PlayerCircle
	autonomousCircles    []AutonomousCircle
	autonomousDirections []Vector
	foods                []Food
	moveCost             float64
	speed                float64
	maxEnergy            float64
	foodGain             float64
	lastInteraction      *InteractionClassification
	activeOverlapPairs   map[string]struct{}
}

type Config struct {
	PlayerShape               string
	AutonomousShape           string
	SecondaryAutonomousShape  string
	PlayerEnergy              float64
	AutonomousEnergy          float64
	SecondaryAutonomousEnergy float64
	PlayerChildrenCount       int
	AutonomousChildrenCount   int
	SecondaryChildrenCount    int
}

func NewWorld() *World {
	return NewWorldWithConfig(Config{
		PlayerShape:               DefaultPlayerShape,
		AutonomousShape:           DefaultPlayerShape,
		SecondaryAutonomousShape:  DefaultAutoShape,
		PlayerEnergy:              DefaultPlayerEnergy,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: DefaultPlayerEnergy,
		AutonomousChildrenCount:   1,
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
	autonomousCircles := []AutonomousCircle{
		{
			ID:            DefaultAutonomousID,
			Shape:         config.AutonomousShape,
			X:             DefaultWorldWidth/2 - 140,
			Y:             DefaultWorldHeight / 2,
			Radius:        derivedRadius(config.AutonomousChildrenCount),
			Energy:        config.AutonomousEnergy,
			ChildrenCount: config.AutonomousChildrenCount,
		},
	}
	if config.SecondaryAutonomousShape != "" {
		autonomousCircles = append(autonomousCircles, AutonomousCircle{
			ID:            DefaultSecondaryID,
			Shape:         config.SecondaryAutonomousShape,
			X:             DefaultWorldWidth/2 + 140,
			Y:             DefaultWorldHeight / 2,
			Radius:        derivedRadius(config.SecondaryChildrenCount),
			Energy:        config.SecondaryAutonomousEnergy,
			ChildrenCount: config.SecondaryChildrenCount,
		})
	}

	return &World{
		bounds: Bounds{
			Width:  DefaultWorldWidth,
			Height: DefaultWorldHeight,
		},
		player: &PlayerCircle{
			ID:            "player-1",
			Shape:         config.PlayerShape,
			X:             DefaultWorldWidth / 2,
			Y:             DefaultWorldHeight / 2,
			Radius:        derivedRadius(config.PlayerChildrenCount),
			Energy:        config.PlayerEnergy,
			ChildrenCount: config.PlayerChildrenCount,
		},
		autonomousCircles:    autonomousCircles,
		autonomousDirections: initialAutonomousDirections(len(autonomousCircles)),
		foods: []Food{
			{ID: "food-1", X: DefaultWorldWidth/2 + 32, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
			{ID: "food-2", X: DefaultWorldWidth/2 - 108, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
			{ID: "food-3", X: DefaultWorldWidth/2 + 120, Y: DefaultWorldHeight/2 + 84, Radius: DefaultFoodRadius},
		},
		moveCost:           DefaultMoveCost,
		speed:              DefaultMoveSpeed,
		maxEnergy:          DefaultMaxEnergy,
		foodGain:           DefaultFoodEnergy,
		activeOverlapPairs: make(map[string]struct{}),
	}
}

func (w *World) Advance(tick int64, intent Vector) Snapshot {
	w.lastInteraction = nil
	w.player = w.advanceCircle(w.player, intent)
	w.advanceAutonomousCircles()

	w.consumeOverlappingFood()
	w.resolveEnergyCollapse()
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

func (w *World) advanceAutonomousCircles() {
	for index, circle := range w.autonomousCircles {
		if circle.Energy <= 0 {
			continue
		}

		intent := w.autonomousDirections[index]
		normalized := normalize(intent)
		if normalized.X == 0 && normalized.Y == 0 {
			continue
		}

		circle.X = clamp(circle.X+normalized.X*w.speed, circle.Radius, w.bounds.Width-circle.Radius)
		circle.Y = clamp(circle.Y+normalized.Y*w.speed, circle.Radius, w.bounds.Height-circle.Radius)
		circle.Energy = math.Max(0, circle.Energy-w.moveCost)
		if circle.X == circle.Radius || circle.X == w.bounds.Width-circle.Radius {
			w.autonomousDirections[index] = Vector{X: -intent.X, Y: intent.Y}
		}
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

func (w *World) resolveCircleInteractions() {
	if w.player == nil {
		w.activeOverlapPairs = make(map[string]struct{})
		return
	}

	currentOverlapPairs := make(map[string]struct{})
	for _, circle := range w.autonomousCircles {
		if !overlaps(w.player.X, w.player.Y, w.player.Radius, circle.X, circle.Y, circle.Radius) {
			continue
		}

		pairKey := overlapPairKey(w.player.ID, circle.ID)
		currentOverlapPairs[pairKey] = struct{}{}

		if w.player.Shape != circle.Shape {
			if _, exists := w.activeOverlapPairs[pairKey]; exists {
				continue
			}

			w.resolveReproduction(circle.ID)
			w.activeOverlapPairs = currentOverlapPairs
			return
		}

		w.resolveFight(circle.ID)
		w.activeOverlapPairs = currentOverlapPairs
		return
	}

	w.activeOverlapPairs = currentOverlapPairs
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
		w.player = replaceOrRemovePlayer(w.player)
		return
	}

	replacedOpponent, active := replaceOrRemoveAutonomous(opponent)
	if !active {
		w.autonomousCircles = append(w.autonomousCircles[:opponentIndex], w.autonomousCircles[opponentIndex+1:]...)
		w.autonomousDirections = append(w.autonomousDirections[:opponentIndex], w.autonomousDirections[opponentIndex+1:]...)
		return
	}

	w.autonomousCircles[opponentIndex] = replacedOpponent
}

func (w *World) resolveReproduction(opponentID string) {
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

	w.player.ChildrenCount++
	w.player.Radius = derivedRadius(w.player.ChildrenCount)

	opponent := w.autonomousCircles[opponentIndex]
	opponent.ChildrenCount++
	opponent.Radius = derivedRadius(opponent.ChildrenCount)
	w.autonomousCircles[opponentIndex] = opponent

	w.lastInteraction = &InteractionClassification{
		Active:   false,
		Resolved: true,
		Kind:     "reproduce_resolved",
		SourceID: w.player.ID,
		TargetID: opponent.ID,
	}
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

func overlapPairKey(a, b string) string {
	if a < b {
		return a + ":" + b
	}

	return b + ":" + a
}

func initialAutonomousDirections(count int) []Vector {
	directions := make([]Vector, 0, count)
	for index := range count {
		if index%2 == 1 {
			directions = append(directions, Vector{X: -1, Y: 0})
			continue
		}

		directions = append(directions, Vector{X: 1, Y: 0})
	}

	return directions
}

func derivedRadius(childrenCount int) float64 {
	return DefaultPlayerRadius + float64(childrenCount)*DefaultChildRadiusGain
}

func (w *World) resolveEnergyCollapse() {
	if w.player != nil && w.player.Energy == 0 {
		w.player = replaceOrRemovePlayer(w.player)
	}

	for index, circle := range w.autonomousCircles {
		if circle.Energy != 0 {
			continue
		}

		replaced, active := replaceOrRemoveAutonomous(circle)
		if !active {
			w.autonomousCircles = append(w.autonomousCircles[:index], w.autonomousCircles[index+1:]...)
			w.autonomousDirections = append(w.autonomousDirections[:index], w.autonomousDirections[index+1:]...)
			index--
			continue
		}

		w.autonomousCircles[index] = replaced
	}
}

func replaceOrRemovePlayer(circle *PlayerCircle) *PlayerCircle {
	if circle == nil {
		return nil
	}
	if circle.ChildrenCount == 0 {
		return nil
	}

	circle.ChildrenCount--
	circle.Radius = derivedRadius(circle.ChildrenCount)
	circle.Energy = DefaultReplacementEnergy

	return circle
}

func replaceOrRemoveAutonomous(circle AutonomousCircle) (AutonomousCircle, bool) {
	if circle.ChildrenCount == 0 {
		return AutonomousCircle{}, false
	}

	circle.ChildrenCount--
	circle.Radius = derivedRadius(circle.ChildrenCount)
	circle.Energy = DefaultReplacementEnergy

	return circle, true
}
