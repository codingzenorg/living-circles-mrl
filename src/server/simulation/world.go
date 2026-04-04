package simulation

import "math"

const (
	DefaultWorldWidth                           = 800.0
	DefaultWorldHeight                          = 600.0
	DefaultPlayerRadius                         = 12.0
	DefaultChildRadiusGain                      = 4.0
	DefaultAttachedChildRadius                  = 4.0
	DefaultAttachedChildOrbitGap                = 8.0
	DefaultChildOrbitSpeed                      = 0.12
	DefaultAutonomousID                         = "circle-2"
	DefaultSecondaryID                          = "circle-3"
	DefaultPlayerEnergy                         = 100.0
	DefaultReplacementEnergy                    = 100.0
	DefaultMaxEnergy                            = 100.0
	DefaultMoveSpeed                            = 8.0
	DefaultMoveCost                             = 1.0
	DefaultFoodRadius                           = 6.0
	DefaultFoodEnergy                           = 10.0
	DefaultFoodRegenDelay                       = int64(12)
	DefaultFoodPriorityDistance                 = 140.0
	DefaultLowEnergyFoodThreshold               = 40.0
	DefaultThreatAvoidanceDistance              = 120.0
	DefaultBlockedReproductionAvoidanceDistance = 120.0
	DefaultReproductionMinEnergy                = 15.0
	DefaultReproductionCost                     = 10.0
	DefaultPlayerShape                          = "triangle"
	DefaultAutoShape                            = "square"
)

type Bounds struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AttachedChild struct {
	ID        string  `json:"id"`
	OwnerID   string  `json:"owner_id"`
	OrbitSlot int     `json:"orbit_slot"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Radius    float64 `json:"radius"`
}

type PlayerCircle struct {
	ID               string          `json:"id"`
	LineageID        string          `json:"lineage_id"`
	Generation       int             `json:"generation"`
	Shape            string          `json:"shape"`
	X                float64         `json:"x"`
	Y                float64         `json:"y"`
	Radius           float64         `json:"radius"`
	Energy           float64         `json:"energy"`
	ChildrenCount    int             `json:"children_count"`
	AttachedChildren []AttachedChild `json:"attached_children"`
}

type AutonomousCircle struct {
	ID               string          `json:"id"`
	LineageID        string          `json:"lineage_id"`
	Generation       int             `json:"generation"`
	Shape            string          `json:"shape"`
	X                float64         `json:"x"`
	Y                float64         `json:"y"`
	Radius           float64         `json:"radius"`
	Energy           float64         `json:"energy"`
	ChildrenCount    int             `json:"children_count"`
	AttachedChildren []AttachedChild `json:"attached_children"`
}

type InteractionClassification struct {
	Active        bool   `json:"active"`
	Resolved      bool   `json:"resolved"`
	Kind          string `json:"kind"`
	ContactOrigin string `json:"contact_origin,omitempty"`
	SourceID      string `json:"source_id"`
	TargetID      string `json:"target_id"`
	WinnerID      string `json:"winner_id"`
	LoserID       string `json:"loser_id"`
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
	bounds                              Bounds
	player                              *PlayerCircle
	autonomousCircles                   []AutonomousCircle
	autonomousDirections                []Vector
	nextChildID                         int
	disableFoodSeeking                  bool
	disableThreatAvoidance              bool
	disableBlockedReproductionAvoidance bool
	foodSlots                           []Food
	foods                               []Food
	missingFoodSince                    map[string]int64
	moveCost                            float64
	speed                               float64
	maxEnergy                           float64
	foodGain                            float64
	lastInteraction                     *InteractionClassification
	activeOverlapPairs                  map[string]struct{}
}

type Config struct {
	PlayerShape                         string
	AutonomousShape                     string
	SecondaryAutonomousShape            string
	PlayerX                             float64
	PlayerY                             float64
	AutonomousX                         float64
	AutonomousY                         float64
	SecondaryAutonomousX                float64
	SecondaryAutonomousY                float64
	PlayerEnergy                        float64
	AutonomousEnergy                    float64
	SecondaryAutonomousEnergy           float64
	PlayerChildrenCount                 int
	AutonomousChildrenCount             int
	SecondaryChildrenCount              int
	DisableFoodSeeking                  bool
	DisableThreatAvoidance              bool
	DisableBlockedReproductionAvoidance bool
}

func NewWorld() *World {
	return NewWorldWithConfig(Config{
		PlayerShape:               DefaultPlayerShape,
		AutonomousShape:           DefaultPlayerShape,
		SecondaryAutonomousShape:  DefaultAutoShape,
		PlayerEnergy:              DefaultPlayerEnergy,
		PlayerChildrenCount:       1,
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
	playerID := "player-1"
	playerX := configuredOrDefault(config.PlayerX, DefaultWorldWidth/2)
	playerY := configuredOrDefault(config.PlayerY, DefaultWorldHeight/2)
	autonomousX := configuredOrDefault(config.AutonomousX, DefaultWorldWidth/2-140)
	autonomousY := configuredOrDefault(config.AutonomousY, DefaultWorldHeight/2)
	autonomousCircles := []AutonomousCircle{
		{
			ID:               DefaultAutonomousID,
			LineageID:        lineageIDFor(DefaultAutonomousID),
			Generation:       0,
			Shape:            config.AutonomousShape,
			X:                autonomousX,
			Y:                autonomousY,
			Radius:           derivedRadius(0),
			Energy:           config.AutonomousEnergy,
			ChildrenCount:    config.AutonomousChildrenCount,
			AttachedChildren: initialAttachedChildren(DefaultAutonomousID, config.AutonomousChildrenCount),
		},
	}
	if config.SecondaryAutonomousShape != "" {
		secondaryX := configuredOrDefault(config.SecondaryAutonomousX, DefaultWorldWidth/2+140)
		secondaryY := configuredOrDefault(config.SecondaryAutonomousY, DefaultWorldHeight/2)
		autonomousCircles = append(autonomousCircles, AutonomousCircle{
			ID:               DefaultSecondaryID,
			LineageID:        lineageIDFor(DefaultSecondaryID),
			Generation:       0,
			Shape:            config.SecondaryAutonomousShape,
			X:                secondaryX,
			Y:                secondaryY,
			Radius:           derivedRadius(0),
			Energy:           config.SecondaryAutonomousEnergy,
			ChildrenCount:    config.SecondaryChildrenCount,
			AttachedChildren: initialAttachedChildren(DefaultSecondaryID, config.SecondaryChildrenCount),
		})
	}

	foodSlots := []Food{
		{ID: "food-1", X: DefaultWorldWidth/2 + 32, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
		{ID: "food-2", X: DefaultWorldWidth/2 - 108, Y: DefaultWorldHeight / 2, Radius: DefaultFoodRadius},
		{ID: "food-3", X: DefaultWorldWidth/2 + 120, Y: DefaultWorldHeight/2 + 84, Radius: DefaultFoodRadius},
	}

	return &World{
		bounds: Bounds{
			Width:  DefaultWorldWidth,
			Height: DefaultWorldHeight,
		},
		player: &PlayerCircle{
			ID:               playerID,
			LineageID:        lineageIDFor(playerID),
			Generation:       0,
			Shape:            config.PlayerShape,
			X:                playerX,
			Y:                playerY,
			Radius:           derivedRadius(0),
			Energy:           config.PlayerEnergy,
			ChildrenCount:    config.PlayerChildrenCount,
			AttachedChildren: initialAttachedChildren(playerID, config.PlayerChildrenCount),
		},
		autonomousCircles:                   autonomousCircles,
		autonomousDirections:                initialAutonomousDirections(len(autonomousCircles)),
		nextChildID:                         totalInitialChildren(config) + 1,
		disableFoodSeeking:                  config.DisableFoodSeeking,
		disableThreatAvoidance:              config.DisableThreatAvoidance,
		disableBlockedReproductionAvoidance: config.DisableBlockedReproductionAvoidance,
		foodSlots:                           foodSlots,
		foods:                               append([]Food(nil), foodSlots...),
		missingFoodSince:                    make(map[string]int64),
		moveCost:                            DefaultMoveCost,
		speed:                               DefaultMoveSpeed,
		maxEnergy:                           DefaultMaxEnergy,
		foodGain:                            DefaultFoodEnergy,
		activeOverlapPairs:                  make(map[string]struct{}),
	}
}

func configuredOrDefault(value, fallback float64) float64 {
	if value == 0 {
		return fallback
	}

	return value
}

func (w *World) Advance(tick int64, intent Vector) Snapshot {
	w.lastInteraction = nil
	w.player = w.advanceCircle(w.player, intent)
	w.advanceAutonomousCircles(tick)

	w.consumeOverlappingFood(tick)
	w.regenerateFood(tick)
	w.resolveEnergyCollapse()
	w.resolveCircleInteractions(tick)

	return w.Snapshot(tick)
}

func (w *World) Snapshot(tick int64) Snapshot {
	var player *PlayerCircle
	if w.player != nil {
		copy := snapshotPlayerCircle(*w.player, tick)
		player = &copy
	}

	var interaction *InteractionClassification
	if w.lastInteraction != nil {
		copy := *w.lastInteraction
		interaction = &copy
	}

	autonomousCircles := make([]AutonomousCircle, 0, len(w.autonomousCircles))
	for _, circle := range w.autonomousCircles {
		autonomousCircles = append(autonomousCircles, snapshotAutonomousCircle(circle, tick))
	}

	return Snapshot{
		Type:              "world_snapshot",
		Tick:              tick,
		World:             w.bounds,
		Player:            player,
		AutonomousCircles: autonomousCircles,
		Interaction:       interaction,
		Foods:             append([]Food(nil), w.foods...),
	}
}

func (w *World) consumeOverlappingFood(tick int64) {
	remaining := make([]Food, 0, len(w.foods))
	for _, food := range w.foods {
		if w.player != nil && playerCollectsFood(*w.player, food, tick) {
			w.player.Energy = math.Min(w.maxEnergy, w.player.Energy+w.foodGain)
			w.missingFoodSince[food.ID] = tick
			continue
		}

		consumed := false
		for index, circle := range w.autonomousCircles {
			if autonomousCollectsFood(circle, food, tick) {
				circle.Energy = math.Min(w.maxEnergy, circle.Energy+w.foodGain)
				w.autonomousCircles[index] = circle
				w.missingFoodSince[food.ID] = tick
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

func (w *World) regenerateFood(tick int64) {
	active := make(map[string]struct{}, len(w.foods))
	for _, food := range w.foods {
		active[food.ID] = struct{}{}
	}

	for _, slot := range w.foodSlots {
		if _, exists := active[slot.ID]; exists {
			continue
		}

		missingSince, tracked := w.missingFoodSince[slot.ID]
		if !tracked || tick-missingSince < DefaultFoodRegenDelay {
			continue
		}

		w.foods = append(w.foods, slot)
		delete(w.missingFoodSince, slot.ID)
	}
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

	circle.X = clamp(circle.X+normalized.X*w.speed, DefaultPlayerRadius, w.bounds.Width-DefaultPlayerRadius)
	circle.Y = clamp(circle.Y+normalized.Y*w.speed, DefaultPlayerRadius, w.bounds.Height-DefaultPlayerRadius)
	circle.Energy = math.Max(0, circle.Energy-w.moveCost)

	return circle
}

func (w *World) advanceAutonomousCircles(tick int64) {
	for index, circle := range w.autonomousCircles {
		if circle.Energy <= 0 {
			continue
		}

		intent := w.autonomousIntent(circle, index, tick)
		normalized := normalize(intent)
		if normalized.X == 0 && normalized.Y == 0 {
			continue
		}

		circle.X = clamp(circle.X+normalized.X*w.speed, DefaultPlayerRadius, w.bounds.Width-DefaultPlayerRadius)
		circle.Y = clamp(circle.Y+normalized.Y*w.speed, DefaultPlayerRadius, w.bounds.Height-DefaultPlayerRadius)
		circle.Energy = math.Max(0, circle.Energy-w.moveCost)
		if circle.X == DefaultPlayerRadius || circle.X == w.bounds.Width-DefaultPlayerRadius {
			w.autonomousDirections[index] = Vector{X: -intent.X, Y: intent.Y}
		}
		w.autonomousCircles[index] = circle
	}
}

func (w *World) autonomousIntent(circle AutonomousCircle, index int, tick int64) Vector {
	if w.disableFoodSeeking {
		return w.autonomousDirections[index]
	}

	foodTarget, foodDistance, foodFound := nearestFoodTarget(circle, w.foods, tick)
	if foodFound && circle.Energy < DefaultLowEnergyFoodThreshold {
		return Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		}
	}

	threatTarget, threatFound := nearestThreatTarget(circle, w.player, w.autonomousCircles, tick)
	if !w.disableThreatAvoidance && threatFound {
		return Vector{
			X: circle.X - threatTarget.X,
			Y: circle.Y - threatTarget.Y,
		}
	}

	blockedTarget, blockedFound := nearestBlockedReproductionTarget(circle, w.player, w.autonomousCircles, tick)
	if !w.disableBlockedReproductionAvoidance && blockedFound {
		return Vector{
			X: circle.X - blockedTarget.X,
			Y: circle.Y - blockedTarget.Y,
		}
	}

	if foodFound && foodDistance <= DefaultFoodPriorityDistance {
		return Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		}
	}

	interactionTarget, _, interactionFound := nearestInteractionTarget(circle, w.player, w.autonomousCircles, tick)
	if interactionFound {
		return Vector{
			X: interactionTarget.X - circle.X,
			Y: interactionTarget.Y - circle.Y,
		}
	}

	if foodFound {
		return Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		}
	}

	return w.autonomousDirections[index]
}

func nearestFoodTarget(circle AutonomousCircle, foods []Food, tick int64) (Food, float64, bool) {
	var selected Food
	bestDistance := 0.0
	found := false

	for _, food := range foods {
		distance := effectiveFoodDistance(circle, food, tick)
		if !found || distance < bestDistance || (distance == bestDistance && food.ID < selected.ID) {
			selected = food
			bestDistance = distance
			found = true
		}
	}

	return selected, bestDistance, found
}

func effectiveFoodDistance(circle AutonomousCircle, food Food, tick int64) float64 {
	bestDistance := distanceBetween(circle.X, circle.Y, food.X, food.Y)
	for _, child := range layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick) {
		distance := distanceBetween(child.X, child.Y, food.X, food.Y)
		if distance < bestDistance {
			bestDistance = distance
		}
	}

	return bestDistance
}

func nearestInteractionTarget(circle AutonomousCircle, player *PlayerCircle, candidates []AutonomousCircle, tick int64) (Vector, string, bool) {
	var selected Vector
	selectedID := ""
	bestDistance := 0.0
	found := false
	selectedPriority := 0

	if player != nil && player.Energy > 0 {
		playerPriority := interactionTargetPriorityAgainstPlayer(circle, *player)
		if playerPriority > 0 {
			target, distance := nearestInteractionPointToPlayer(circle, *player, tick)
			selected = target
			selectedID = player.ID
			bestDistance = distance
			found = true
			selectedPriority = playerPriority
		}
	}

	for _, candidate := range candidates {
		if candidate.ID == circle.ID || candidate.Energy <= 0 {
			continue
		}

		candidatePriority := interactionTargetPriorityAgainstAutonomous(circle, candidate)
		if candidatePriority == 0 {
			continue
		}

		target, distance := nearestInteractionPointToAutonomous(circle, candidate, tick)
		if !found ||
			candidatePriority > selectedPriority ||
			(candidatePriority == selectedPriority && (distance < bestDistance || (distance == bestDistance && candidate.ID < selectedID))) {
			selected = target
			selectedID = candidate.ID
			bestDistance = distance
			found = true
			selectedPriority = candidatePriority
		}
	}

	return selected, selectedID, found
}

func nearestThreatTarget(circle AutonomousCircle, player *PlayerCircle, candidates []AutonomousCircle, tick int64) (Vector, bool) {
	var selected Vector
	selectedID := ""
	bestDistance := 0.0
	found := false

	if player != nil && player.Energy > 0 && playerThreatensAutonomous(circle, *player) {
		target, distance := nearestInteractionPointToPlayer(circle, *player, tick)
		if distance < DefaultThreatAvoidanceDistance {
			selected = target
			selectedID = player.ID
			bestDistance = distance
			found = true
		}
	}

	for _, candidate := range candidates {
		if candidate.ID == circle.ID || candidate.Energy <= 0 || !autonomousThreatensAutonomous(circle, candidate) {
			continue
		}

		target, distance := nearestInteractionPointToAutonomous(circle, candidate, tick)
		if distance >= DefaultThreatAvoidanceDistance {
			continue
		}

		if !found || distance < bestDistance || (distance == bestDistance && candidate.ID < selectedID) {
			selected = target
			selectedID = candidate.ID
			bestDistance = distance
			found = true
		}
	}

	return selected, found
}

func nearestBlockedReproductionTarget(circle AutonomousCircle, player *PlayerCircle, candidates []AutonomousCircle, tick int64) (Vector, bool) {
	var selected Vector
	selectedID := ""
	bestDistance := 0.0
	found := false

	if player != nil && player.Energy > 0 && blockedReproductionWithPlayer(circle, *player) {
		target, distance := nearestInteractionPointToPlayer(circle, *player, tick)
		if distance < DefaultBlockedReproductionAvoidanceDistance {
			selected = target
			selectedID = player.ID
			bestDistance = distance
			found = true
		}
	}

	for _, candidate := range candidates {
		if candidate.ID == circle.ID || candidate.Energy <= 0 || !blockedReproductionWithAutonomous(circle, candidate) {
			continue
		}

		target, distance := nearestInteractionPointToAutonomous(circle, candidate, tick)
		if distance >= DefaultBlockedReproductionAvoidanceDistance {
			continue
		}

		if !found || distance < bestDistance || (distance == bestDistance && candidate.ID < selectedID) {
			selected = target
			selectedID = candidate.ID
			bestDistance = distance
			found = true
		}
	}

	return selected, found
}

func nearestInteractionPointToPlayer(circle AutonomousCircle, player PlayerCircle, tick int64) (Vector, float64) {
	selected := Vector{X: player.X, Y: player.Y}
	bestDistance := distanceBetween(circle.X, circle.Y, player.X, player.Y)

	for _, child := range layoutAttachedChildren(player.ID, player.X, player.Y, player.Radius, player.AttachedChildren, tick) {
		distance := distanceBetween(circle.X, circle.Y, child.X, child.Y)
		if distance < bestDistance {
			selected = Vector{X: child.X, Y: child.Y}
			bestDistance = distance
		}
	}

	return selected, bestDistance
}

func nearestInteractionPointToAutonomous(circle AutonomousCircle, candidate AutonomousCircle, tick int64) (Vector, float64) {
	selected := Vector{X: candidate.X, Y: candidate.Y}
	bestDistance := distanceBetween(circle.X, circle.Y, candidate.X, candidate.Y)

	for _, child := range layoutAttachedChildren(candidate.ID, candidate.X, candidate.Y, candidate.Radius, candidate.AttachedChildren, tick) {
		distance := distanceBetween(circle.X, circle.Y, child.X, child.Y)
		if distance < bestDistance {
			selected = Vector{X: child.X, Y: child.Y}
			bestDistance = distance
		}
	}

	return selected, bestDistance
}

func interactionTargetPriorityAgainstPlayer(circle AutonomousCircle, player PlayerCircle) int {
	if circle.Shape != player.Shape {
		if reproductionFeasibleWithPlayer(circle, player) {
			return 2
		}
		return 0
	}

	if fightFeasibleWithPlayer(circle, player) {
		return 1
	}

	return 0
}

func interactionTargetPriorityAgainstAutonomous(circle AutonomousCircle, candidate AutonomousCircle) int {
	if circle.Shape != candidate.Shape {
		if reproductionFeasibleWithAutonomous(circle, candidate) {
			return 2
		}
		return 0
	}

	if fightFeasibleWithAutonomous(circle, candidate) {
		return 1
	}

	return 0
}

func reproductionFeasibleWithPlayer(circle AutonomousCircle, player PlayerCircle) bool {
	return reproductionCapacity(circle.Energy, childCountForAutonomous(circle)) >= DefaultReproductionMinEnergy &&
		reproductionCapacity(player.Energy, childCountForPlayer(player)) >= DefaultReproductionMinEnergy
}

func reproductionFeasibleWithAutonomous(left AutonomousCircle, right AutonomousCircle) bool {
	return reproductionCapacity(left.Energy, childCountForAutonomous(left)) >= DefaultReproductionMinEnergy &&
		reproductionCapacity(right.Energy, childCountForAutonomous(right)) >= DefaultReproductionMinEnergy
}

func blockedReproductionWithPlayer(circle AutonomousCircle, player PlayerCircle) bool {
	return circle.Shape != player.Shape && !reproductionFeasibleWithPlayer(circle, player)
}

func blockedReproductionWithAutonomous(left AutonomousCircle, right AutonomousCircle) bool {
	return left.Shape != right.Shape && !reproductionFeasibleWithAutonomous(left, right)
}

func fightFeasibleWithPlayer(circle AutonomousCircle, player PlayerCircle) bool {
	winnerID, _ := determineFightOutcome(player, circle)
	return winnerID == circle.ID
}

func fightFeasibleWithAutonomous(left AutonomousCircle, right AutonomousCircle) bool {
	winnerID, _ := determineAutonomousFightOutcome(left, right)
	return winnerID == left.ID
}

func playerThreatensAutonomous(circle AutonomousCircle, player PlayerCircle) bool {
	if circle.Shape != player.Shape {
		return false
	}

	winnerID, _ := determineFightOutcome(player, circle)
	return winnerID == player.ID
}

func autonomousThreatensAutonomous(circle AutonomousCircle, candidate AutonomousCircle) bool {
	if circle.Shape != candidate.Shape {
		return false
	}

	winnerID, _ := determineAutonomousFightOutcome(circle, candidate)
	return winnerID == candidate.ID
}

func distanceBetween(ax, ay, bx, by float64) float64 {
	return math.Hypot(ax-bx, ay-by)
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

func playerCollectsFood(circle PlayerCircle, food Food, tick int64) bool {
	if overlaps(circle.X, circle.Y, DefaultPlayerRadius, food.X, food.Y, food.Radius) {
		return true
	}

	children := layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	for _, child := range children {
		if overlaps(child.X, child.Y, child.Radius, food.X, food.Y, food.Radius) {
			return true
		}
	}

	return false
}

func autonomousCollectsFood(circle AutonomousCircle, food Food, tick int64) bool {
	if overlaps(circle.X, circle.Y, DefaultPlayerRadius, food.X, food.Y, food.Radius) {
		return true
	}

	children := layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	for _, child := range children {
		if overlaps(child.X, child.Y, child.Radius, food.X, food.Y, food.Radius) {
			return true
		}
	}

	return false
}

func (w *World) resolveCircleInteractions(tick int64) {
	currentOverlapPairs := make(map[string]struct{})
	if w.player != nil {
		for _, circle := range w.autonomousCircles {
			contactOrigin, interacting := circlesInteract(*w.player, circle, tick)
			if !interacting {
				continue
			}

			pairKey := overlapPairKey(w.player.ID, circle.ID)
			currentOverlapPairs[pairKey] = struct{}{}

			if w.player.Shape != circle.Shape {
				if _, exists := w.activeOverlapPairs[pairKey]; exists {
					continue
				}

				w.resolveReproduction(circle.ID, tick, contactOrigin)
				w.activeOverlapPairs = currentOverlapPairs
				return
			}

			w.resolveFight(circle.ID, contactOrigin)
			w.activeOverlapPairs = currentOverlapPairs
			return
		}
	}

	for left := 0; left < len(w.autonomousCircles); left++ {
		for right := left + 1; right < len(w.autonomousCircles); right++ {
			leftCircle := w.autonomousCircles[left]
			rightCircle := w.autonomousCircles[right]
			contactOrigin, interacting := autonomousCirclesInteract(leftCircle, rightCircle, tick)
			if !interacting {
				continue
			}

			pairKey := overlapPairKey(leftCircle.ID, rightCircle.ID)
			currentOverlapPairs[pairKey] = struct{}{}

			if leftCircle.Shape != rightCircle.Shape {
				if _, exists := w.activeOverlapPairs[pairKey]; exists {
					continue
				}

				w.resolveAutonomousReproduction(left, right, tick, contactOrigin)
				w.activeOverlapPairs = currentOverlapPairs
				return
			}

			w.resolveAutonomousFight(left, right, contactOrigin)
			w.activeOverlapPairs = currentOverlapPairs
			return
		}
	}

	w.activeOverlapPairs = currentOverlapPairs
}

func circlesInteract(player PlayerCircle, autonomous AutonomousCircle, tick int64) (string, bool) {
	parentBodyOverlap := overlaps(player.X, player.Y, DefaultPlayerRadius, autonomous.X, autonomous.Y, DefaultPlayerRadius)

	playerChildren := layoutAttachedChildren(player.ID, player.X, player.Y, player.Radius, player.AttachedChildren, tick)
	autonomousChildren := layoutAttachedChildren(autonomous.ID, autonomous.X, autonomous.Y, autonomous.Radius, autonomous.AttachedChildren, tick)

	for _, child := range playerChildren {
		if overlaps(child.X, child.Y, child.Radius, autonomous.X, autonomous.Y, DefaultPlayerRadius) {
			return "attached_child", true
		}
	}

	for _, child := range autonomousChildren {
		if overlaps(player.X, player.Y, DefaultPlayerRadius, child.X, child.Y, child.Radius) {
			return "attached_child", true
		}
	}

	for _, playerChild := range playerChildren {
		for _, autonomousChild := range autonomousChildren {
			if overlaps(playerChild.X, playerChild.Y, playerChild.Radius, autonomousChild.X, autonomousChild.Y, autonomousChild.Radius) {
				return "attached_child", true
			}
		}
	}

	if parentBodyOverlap {
		return "parent_body", true
	}

	return "", false
}

func autonomousCirclesInteract(left AutonomousCircle, right AutonomousCircle, tick int64) (string, bool) {
	parentBodyOverlap := overlaps(left.X, left.Y, DefaultPlayerRadius, right.X, right.Y, DefaultPlayerRadius)

	leftChildren := layoutAttachedChildren(left.ID, left.X, left.Y, left.Radius, left.AttachedChildren, tick)
	rightChildren := layoutAttachedChildren(right.ID, right.X, right.Y, right.Radius, right.AttachedChildren, tick)

	for _, child := range leftChildren {
		if overlaps(child.X, child.Y, child.Radius, right.X, right.Y, DefaultPlayerRadius) {
			return "attached_child", true
		}
	}

	for _, child := range rightChildren {
		if overlaps(left.X, left.Y, DefaultPlayerRadius, child.X, child.Y, child.Radius) {
			return "attached_child", true
		}
	}

	for _, leftChild := range leftChildren {
		for _, rightChild := range rightChildren {
			if overlaps(leftChild.X, leftChild.Y, leftChild.Radius, rightChild.X, rightChild.Y, rightChild.Radius) {
				return "attached_child", true
			}
		}
	}

	if parentBodyOverlap {
		return "parent_body", true
	}

	return "", false
}

func (w *World) resolveFight(opponentID string, contactOrigin string) {
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
		Active:        false,
		Resolved:      true,
		Kind:          "",
		ContactOrigin: contactOrigin,
		SourceID:      w.player.ID,
		TargetID:      opponent.ID,
		WinnerID:      winnerID,
		LoserID:       loserID,
	}

	if loserID == w.player.ID {
		if childCountForPlayer(*w.player) > 0 {
			consumePlayerChild(w.player)
			w.lastInteraction.Kind = "fight_absorbed_child"
			return
		}
		w.lastInteraction.Kind = "fight_resolved"
		w.player = replaceOrRemovePlayer(w.player)
		return
	}

	if childCountForAutonomous(opponent) > 0 {
		consumeAutonomousChild(&opponent)
		w.autonomousCircles[opponentIndex] = opponent
		w.lastInteraction.Kind = "fight_absorbed_child"
		return
	}

	w.lastInteraction.Kind = "fight_resolved"
	replacedOpponent, active := replaceOrRemoveAutonomous(opponent)
	if !active {
		w.autonomousCircles = append(w.autonomousCircles[:opponentIndex], w.autonomousCircles[opponentIndex+1:]...)
		w.autonomousDirections = append(w.autonomousDirections[:opponentIndex], w.autonomousDirections[opponentIndex+1:]...)
		return
	}

	w.autonomousCircles[opponentIndex] = replacedOpponent
}

func (w *World) resolveReproduction(opponentID string, tick int64, contactOrigin string) {
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
	playerPaid, playerUsedChild := payPlayerReproductionCost(w.player)
	opponentPaid, opponentUsedChild, opponent := payAutonomousReproductionCost(opponent)
	if !playerPaid || !opponentPaid {
		w.lastInteraction = &InteractionClassification{
			Active:        false,
			Resolved:      true,
			Kind:          "reproduce_blocked_energy",
			ContactOrigin: contactOrigin,
			SourceID:      w.player.ID,
			TargetID:      opponent.ID,
		}
		return
	}

	w.assignReproductionChildren(&opponent, tick)
	w.autonomousCircles[opponentIndex] = opponent

	kind := "reproduce_resolved"
	if playerUsedChild || opponentUsedChild {
		kind = "reproduce_paid_child"
	}

	w.lastInteraction = &InteractionClassification{
		Active:        false,
		Resolved:      true,
		Kind:          kind,
		ContactOrigin: contactOrigin,
		SourceID:      w.player.ID,
		TargetID:      opponent.ID,
	}
}

func determineAutonomousFightOutcome(left AutonomousCircle, right AutonomousCircle) (string, string) {
	if left.Energy > right.Energy {
		return left.ID, right.ID
	}
	if right.Energy > left.Energy {
		return right.ID, left.ID
	}
	if childCountForAutonomous(left) > childCountForAutonomous(right) {
		return left.ID, right.ID
	}
	if childCountForAutonomous(right) > childCountForAutonomous(left) {
		return right.ID, left.ID
	}
	if left.ID < right.ID {
		return left.ID, right.ID
	}

	return right.ID, left.ID
}

func (w *World) resolveAutonomousFight(leftIndex int, rightIndex int, contactOrigin string) {
	if leftIndex < 0 || rightIndex < 0 || leftIndex >= len(w.autonomousCircles) || rightIndex >= len(w.autonomousCircles) {
		return
	}

	leftCircle := w.autonomousCircles[leftIndex]
	rightCircle := w.autonomousCircles[rightIndex]
	winnerID, loserID := determineAutonomousFightOutcome(leftCircle, rightCircle)

	w.lastInteraction = &InteractionClassification{
		Active:        false,
		Resolved:      true,
		Kind:          "",
		ContactOrigin: contactOrigin,
		SourceID:      leftCircle.ID,
		TargetID:      rightCircle.ID,
		WinnerID:      winnerID,
		LoserID:       loserID,
	}

	loserIndex := leftIndex
	if loserID == rightCircle.ID {
		loserIndex = rightIndex
	}
	loser := w.autonomousCircles[loserIndex]
	if childCountForAutonomous(loser) > 0 {
		consumeAutonomousChild(&loser)
		w.autonomousCircles[loserIndex] = loser
		w.lastInteraction.Kind = "fight_absorbed_child"
		return
	}

	w.lastInteraction.Kind = "fight_resolved"
	replacedLoser, active := replaceOrRemoveAutonomous(loser)
	if active {
		w.autonomousCircles[loserIndex] = replacedLoser
		return
	}

	w.autonomousCircles = append(w.autonomousCircles[:loserIndex], w.autonomousCircles[loserIndex+1:]...)
	w.autonomousDirections = append(w.autonomousDirections[:loserIndex], w.autonomousDirections[loserIndex+1:]...)
}

func (w *World) resolveAutonomousReproduction(leftIndex int, rightIndex int, tick int64, contactOrigin string) {
	if leftIndex < 0 || rightIndex < 0 || leftIndex >= len(w.autonomousCircles) || rightIndex >= len(w.autonomousCircles) {
		return
	}

	leftCircle := w.autonomousCircles[leftIndex]
	rightCircle := w.autonomousCircles[rightIndex]
	leftPaid, leftUsedChild, leftCircle := payAutonomousReproductionCost(leftCircle)
	rightPaid, rightUsedChild, rightCircle := payAutonomousReproductionCost(rightCircle)
	if !leftPaid || !rightPaid {
		w.lastInteraction = &InteractionClassification{
			Active:        false,
			Resolved:      true,
			Kind:          "reproduce_blocked_energy",
			ContactOrigin: contactOrigin,
			SourceID:      leftCircle.ID,
			TargetID:      rightCircle.ID,
		}
		return
	}

	w.assignAutonomousPairReproductionChildren(&leftCircle, &rightCircle, tick)
	w.autonomousCircles[leftIndex] = leftCircle
	w.autonomousCircles[rightIndex] = rightCircle

	kind := "reproduce_resolved"
	if leftUsedChild || rightUsedChild {
		kind = "reproduce_paid_child"
	}

	w.lastInteraction = &InteractionClassification{
		Active:        false,
		Resolved:      true,
		Kind:          kind,
		ContactOrigin: contactOrigin,
		SourceID:      leftCircle.ID,
		TargetID:      rightCircle.ID,
	}
}

func determineFightOutcome(player PlayerCircle, opponent AutonomousCircle) (string, string) {
	if player.Energy > opponent.Energy {
		return player.ID, opponent.ID
	}
	if opponent.Energy > player.Energy {
		return opponent.ID, player.ID
	}
	if childCountForPlayer(player) > childCountForAutonomous(opponent) {
		return player.ID, opponent.ID
	}
	if childCountForAutonomous(opponent) > childCountForPlayer(player) {
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
	_ = childrenCount
	return DefaultPlayerRadius
}

func lineageIDFor(circleID string) string {
	return "lineage-" + circleID
}

func intToString(value int) string {
	if value == 0 {
		return "0"
	}

	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}

	digits := make([]byte, 0, 10)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}

	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}

	return sign + string(digits)
}

func payPlayerReproductionCost(circle *PlayerCircle) (bool, bool) {
	if circle == nil {
		return false, false
	}
	if reproductionCapacity(circle.Energy, childCountForPlayer(*circle)) < DefaultReproductionMinEnergy {
		return false, false
	}
	if circle.Energy >= DefaultReproductionCost {
		circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
		return true, false
	}
	if childCountForPlayer(*circle) == 0 {
		return false, false
	}

	consumePlayerChild(circle)
	circle.Energy += DefaultReproductionCost
	circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
	return true, true
}

func payAutonomousReproductionCost(circle AutonomousCircle) (bool, bool, AutonomousCircle) {
	if reproductionCapacity(circle.Energy, childCountForAutonomous(circle)) < DefaultReproductionMinEnergy {
		return false, false, circle
	}
	if circle.Energy >= DefaultReproductionCost {
		circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
		return true, false, circle
	}
	if childCountForAutonomous(circle) == 0 {
		return false, false, circle
	}

	consumeAutonomousChild(&circle)
	circle.Energy += DefaultReproductionCost
	circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
	return true, true, circle
}

func reproductionCapacity(energy float64, childrenCount int) float64 {
	if childrenCount == 0 {
		return energy
	}

	return energy + DefaultReproductionCost
}

func (w *World) resolveEnergyCollapse() {
	if w.player != nil && w.player.Energy == 0 {
		promoted := childCountForPlayer(*w.player) > 0
		w.player = replaceOrRemovePlayer(w.player)
		if promoted && w.player != nil && w.lastInteraction == nil {
			w.lastInteraction = &InteractionClassification{
				Active:   false,
				Resolved: true,
				Kind:     "death_promoted_child",
				SourceID: w.player.ID,
				TargetID: w.player.ID,
			}
		}
	}

	for index, circle := range w.autonomousCircles {
		if circle.Energy != 0 {
			continue
		}

		promoted := childCountForAutonomous(circle) > 0
		replaced, active := replaceOrRemoveAutonomous(circle)
		if !active {
			w.autonomousCircles = append(w.autonomousCircles[:index], w.autonomousCircles[index+1:]...)
			w.autonomousDirections = append(w.autonomousDirections[:index], w.autonomousDirections[index+1:]...)
			index--
			continue
		}

		w.autonomousCircles[index] = replaced
		if promoted && w.lastInteraction == nil {
			w.lastInteraction = &InteractionClassification{
				Active:   false,
				Resolved: true,
				Kind:     "death_promoted_child",
				SourceID: replaced.ID,
				TargetID: replaced.ID,
			}
		}
	}
}

func replaceOrRemovePlayer(circle *PlayerCircle) *PlayerCircle {
	if circle == nil {
		return nil
	}
	if childCountForPlayer(*circle) == 0 {
		return nil
	}

	consumePlayerChild(circle)
	circle.Generation++
	circle.Energy = DefaultReplacementEnergy

	return circle
}

func replaceOrRemoveAutonomous(circle AutonomousCircle) (AutonomousCircle, bool) {
	if childCountForAutonomous(circle) == 0 {
		return AutonomousCircle{}, false
	}

	consumeAutonomousChild(&circle)
	circle.Generation++
	circle.Energy = DefaultReplacementEnergy

	return circle, true
}

func initialAttachedChildren(ownerID string, count int) []AttachedChild {
	children := make([]AttachedChild, 0, count)
	for index := 0; index < count; index++ {
		children = append(children, AttachedChild{
			ID:      ownerID + "-child-" + intToString(index+1),
			OwnerID: ownerID,
		})
	}

	return children
}

func totalInitialChildren(config Config) int {
	return config.PlayerChildrenCount + config.AutonomousChildrenCount + config.SecondaryChildrenCount
}

func snapshotPlayerCircle(circle PlayerCircle, tick int64) PlayerCircle {
	copy := circle
	copy.ChildrenCount = childCountForPlayer(circle)
	copy.AttachedChildren = layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	return copy
}

func snapshotAutonomousCircle(circle AutonomousCircle, tick int64) AutonomousCircle {
	copy := circle
	copy.ChildrenCount = childCountForAutonomous(circle)
	copy.AttachedChildren = layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	return copy
}

func layoutAttachedChildren(ownerID string, x, y, parentRadius float64, children []AttachedChild, tick int64) []AttachedChild {
	if len(children) == 0 {
		return []AttachedChild{}
	}

	positioned := make([]AttachedChild, 0, len(children))
	_ = parentRadius
	orbitRadius := DefaultPlayerRadius + DefaultAttachedChildOrbitGap + DefaultAttachedChildRadius
	for index, child := range children {
		angle := childOrbitAngle(ownerID, child.ID, tick, index, len(children))
		positioned = append(positioned, AttachedChild{
			ID:        child.ID,
			OwnerID:   ownerID,
			OrbitSlot: index,
			X:         x + math.Cos(angle)*orbitRadius,
			Y:         y + math.Sin(angle)*orbitRadius,
			Radius:    DefaultAttachedChildRadius,
		})
	}

	return positioned
}

func childOrbitAngle(ownerID, childID string, tick int64, index, total int) float64 {
	baseAngle := float64(hashString(ownerID+":"+childID)%360) * math.Pi / 180
	slotOffset := 0.0
	if total > 0 {
		slotOffset = float64(index) * (2 * math.Pi / float64(total))
	}

	return baseAngle + slotOffset + float64(tick)*DefaultChildOrbitSpeed
}

func hashString(value string) int {
	hash := 17
	for _, char := range value {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		return -hash
	}
	return hash
}

func (w *World) assignReproductionChildren(opponent *AutonomousCircle, tick int64) {
	distribution := reproductionDistributionCase(tick, *w.player, *opponent)
	switch distribution {
	case 0:
		w.addPlayerChild()
		w.addPlayerChild()
	case 1:
		w.addPlayerChild()
		addAutonomousChild(opponent, w.allocateChildID())
	case 2:
		addAutonomousChild(opponent, w.allocateChildID())
		addAutonomousChild(opponent, w.allocateChildID())
	}
}

func reproductionDistributionCase(tick int64, player PlayerCircle, opponent AutonomousCircle) int {
	seed := int(tick) + len(player.ID) + len(opponent.ID) + player.Generation + opponent.Generation + childCountForPlayer(player) + childCountForAutonomous(opponent)
	return seed % 3
}

func autonomousReproductionDistributionCase(tick int64, left AutonomousCircle, right AutonomousCircle) int {
	seed := int(tick) + len(left.ID) + len(right.ID) + left.Generation + right.Generation + childCountForAutonomous(left) + childCountForAutonomous(right)
	return seed % 3
}

func (w *World) allocateChildID() string {
	childID := "child-" + intToString(w.nextChildID)
	w.nextChildID++
	return childID
}

func (w *World) addPlayerChild() {
	if w.player == nil {
		return
	}
	w.player.AttachedChildren = append(w.player.AttachedChildren, AttachedChild{
		ID:      w.allocateChildID(),
		OwnerID: w.player.ID,
	})
	syncPlayerChildrenState(w.player)
}

func addAutonomousChild(circle *AutonomousCircle, childID string) {
	circle.AttachedChildren = append(circle.AttachedChildren, AttachedChild{
		ID:      childID,
		OwnerID: circle.ID,
	})
	syncAutonomousChildrenState(circle)
}

func (w *World) assignAutonomousPairReproductionChildren(left *AutonomousCircle, right *AutonomousCircle, tick int64) {
	distribution := autonomousReproductionDistributionCase(tick, *left, *right)
	switch distribution {
	case 0:
		addAutonomousChild(left, w.allocateChildID())
		addAutonomousChild(left, w.allocateChildID())
	case 1:
		addAutonomousChild(left, w.allocateChildID())
		addAutonomousChild(right, w.allocateChildID())
	case 2:
		addAutonomousChild(right, w.allocateChildID())
		addAutonomousChild(right, w.allocateChildID())
	}
}

func consumePlayerChild(circle *PlayerCircle) {
	if len(circle.AttachedChildren) > 0 {
		circle.AttachedChildren = circle.AttachedChildren[:len(circle.AttachedChildren)-1]
	}
	syncPlayerChildrenState(circle)
}

func consumeAutonomousChild(circle *AutonomousCircle) {
	if len(circle.AttachedChildren) > 0 {
		circle.AttachedChildren = circle.AttachedChildren[:len(circle.AttachedChildren)-1]
	}
	syncAutonomousChildrenState(circle)
}

func childCountForPlayer(circle PlayerCircle) int {
	return len(circle.AttachedChildren)
}

func childCountForAutonomous(circle AutonomousCircle) int {
	return len(circle.AttachedChildren)
}

func syncPlayerChildrenState(circle *PlayerCircle) {
	circle.ChildrenCount = childCountForPlayer(*circle)
	circle.Radius = derivedRadius(circle.ChildrenCount)
}

func syncAutonomousChildrenState(circle *AutonomousCircle) {
	circle.ChildrenCount = childCountForAutonomous(*circle)
	circle.Radius = derivedRadius(circle.ChildrenCount)
}
