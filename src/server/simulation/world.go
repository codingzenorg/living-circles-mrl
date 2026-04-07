package simulation

import (
	"encoding/json"
	"math"
	"strconv"
)

const (
	DefaultWorldWidth                           = 800.0
	DefaultWorldHeight                          = 600.0
	DefaultExpandedWorldWidth                   = 1600.0
	DefaultExpandedWorldHeight                  = 1200.0
	DefaultPlayerRadius                         = 12.0
	DefaultAttachedChildRadius                  = 4.0
	DefaultAttachedChildOrbitGap                = 8.0
	DefaultChildOrbitSpeed                      = 0.12
	DefaultAutonomousID                         = "circle-2"
	DefaultSecondaryID                          = "circle-3"
	DefaultTertiaryID                           = "circle-4"
	DefaultQuaternaryID                         = "circle-5"
	DefaultQuinaryID                            = "circle-6"
	DefaultSenaryID                             = "circle-7"
	DefaultSeptenaryID                          = "circle-8"
	DefaultOctonaryID                           = "circle-9"
	DefaultPlayerEnergy                         = 100.0
	DefaultReplacementEnergy                    = 100.0
	DefaultMaxEnergy                            = 100.0
	DefaultMoveSpeed                            = 8.0
	DefaultMoveCost                             = 1.0
	DefaultFoodRadius                           = 6.0
	DefaultFoodEnergy                           = 10.0
	DefaultFoodRegenDelay                       = int64(12)
	DefaultRegionalFoodPressureDistance         = 180.0
	DefaultFoodPriorityDistance                 = 140.0
	DefaultLowEnergyFoodThreshold               = 40.0
	DefaultThreatAvoidanceDistance              = 120.0
	DefaultBlockedReproductionAvoidanceDistance = 120.0
	DefaultCrowdingDistance                     = 120.0
	DefaultCrowdingThreshold                    = 2
	DefaultCrowdingMoveCost                     = 1.0
	DefaultReproductionMinEnergy                = 15.0
	DefaultReproductionCost                     = 10.0
	DefaultPlayerShape                          = "triangle"
	DefaultAutoShape                            = "square"
	DefaultExpandedAutonomousCount              = 8
	DefaultExpandedFoodCount                    = DefaultExpandedAutonomousCount + 3
	DefaultExpandedFoodSeed               int64 = 73
	DefaultExpandedAutonomousSeed         int64 = 131
	DefaultExpandedAutonomousStateSeed    int64 = 211
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
	AttachedChildren []AttachedChild `json:"attached_children"`
}

func (circle *PlayerCircle) UnmarshalJSON(data []byte) error {
	type playerCircleJSON struct {
		ID               string          `json:"id"`
		LineageID        string          `json:"lineage_id"`
		Generation       int             `json:"generation"`
		Shape            string          `json:"shape"`
		X                float64         `json:"x"`
		Y                float64         `json:"y"`
		Radius           float64         `json:"radius"`
		Energy           float64         `json:"energy"`
		AttachedChildren []AttachedChild `json:"attached_children"`
	}

	var decoded playerCircleJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	circle.ID = decoded.ID
	circle.LineageID = decoded.LineageID
	circle.Generation = decoded.Generation
	circle.Shape = decoded.Shape
	circle.X = decoded.X
	circle.Y = decoded.Y
	circle.Radius = decoded.Radius
	circle.Energy = decoded.Energy
	circle.AttachedChildren = decoded.AttachedChildren
	return nil
}

func (circle *AutonomousCircle) UnmarshalJSON(data []byte) error {
	type autonomousCircleJSON struct {
		ID               string          `json:"id"`
		LineageID        string          `json:"lineage_id"`
		Generation       int             `json:"generation"`
		Shape            string          `json:"shape"`
		X                float64         `json:"x"`
		Y                float64         `json:"y"`
		Radius           float64         `json:"radius"`
		Energy           float64         `json:"energy"`
		AttachedChildren []AttachedChild `json:"attached_children"`
	}

	var decoded autonomousCircleJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	circle.ID = decoded.ID
	circle.LineageID = decoded.LineageID
	circle.Generation = decoded.Generation
	circle.Shape = decoded.Shape
	circle.X = decoded.X
	circle.Y = decoded.Y
	circle.Radius = decoded.Radius
	circle.Energy = decoded.Energy
	circle.AttachedChildren = decoded.AttachedChildren
	return nil
}

type InteractionClassification struct {
	Active                 bool     `json:"active"`
	Resolved               bool     `json:"resolved"`
	Kind                   string   `json:"kind"`
	ContactOrigin          string   `json:"contact_origin,omitempty"`
	ContactPathKind        string   `json:"contact_path_kind,omitempty"`
	DistributionKind       string   `json:"distribution_kind,omitempty"`
	SourceID               string   `json:"source_id"`
	TargetID               string   `json:"target_id"`
	SourceChildID          string   `json:"source_child_id,omitempty"`
	TargetChildID          string   `json:"target_child_id,omitempty"`
	WinnerID               string   `json:"winner_id"`
	LoserID                string   `json:"loser_id"`
	PromotedChildID        string   `json:"promoted_child_id,omitempty"`
	AbsorbedChildID        string   `json:"absorbed_child_id,omitempty"`
	SourcePaidChild        bool     `json:"source_paid_child,omitempty"`
	TargetPaidChild        bool     `json:"target_paid_child,omitempty"`
	SourceBlockedCapacity  bool     `json:"source_blocked_capacity,omitempty"`
	TargetBlockedCapacity  bool     `json:"target_blocked_capacity,omitempty"`
	SourceCapacityValue    float64  `json:"source_capacity_value,omitempty"`
	TargetCapacityValue    float64  `json:"target_capacity_value,omitempty"`
	SourceEnergyComponent  float64  `json:"source_energy_component,omitempty"`
	TargetEnergyComponent  float64  `json:"target_energy_component,omitempty"`
	SourceReserveComponent float64  `json:"source_reserve_component,omitempty"`
	TargetReserveComponent float64  `json:"target_reserve_component,omitempty"`
	ReproductionThreshold  float64  `json:"reproduction_threshold,omitempty"`
	ReproductionCost       float64  `json:"reproduction_cost,omitempty"`
	SourcePaidChildID      string   `json:"source_paid_child_id,omitempty"`
	TargetPaidChildID      string   `json:"target_paid_child_id,omitempty"`
	CreatedChildIDs        []string `json:"created_child_ids,omitempty"`
	SourceCreatedChildIDs  []string `json:"source_created_child_ids,omitempty"`
	TargetCreatedChildIDs  []string `json:"target_created_child_ids,omitempty"`
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
	WorldWidth                          float64
	WorldHeight                         float64
	UseExpandedPopulation               bool
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

type ContactDetails struct {
	Origin        string
	PathKind      string
	SourceChildID string
	TargetChildID string
}

func NewWorld() *World {
	return NewWorldWithConfig(Config{
		WorldWidth:                DefaultExpandedWorldWidth,
		WorldHeight:               DefaultExpandedWorldHeight,
		UseExpandedPopulation:     true,
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
	worldWidth := configuredOrDefault(config.WorldWidth, DefaultWorldWidth)
	worldHeight := configuredOrDefault(config.WorldHeight, DefaultWorldHeight)
	playerX := configuredOrDefault(config.PlayerX, worldWidth/2)
	playerY := configuredOrDefault(config.PlayerY, worldHeight/2)
	autonomousX := configuredOrDefault(config.AutonomousX, worldWidth/2-140)
	autonomousY := configuredOrDefault(config.AutonomousY, worldHeight/2)
	autonomousCircles := []AutonomousCircle{
		{
			ID:               DefaultAutonomousID,
			LineageID:        lineageIDFor(DefaultAutonomousID),
			Generation:       0,
			Shape:            config.AutonomousShape,
			X:                autonomousX,
			Y:                autonomousY,
			Radius:           DefaultPlayerRadius,
			Energy:           config.AutonomousEnergy,
			AttachedChildren: initialAttachedChildren(DefaultAutonomousID, config.AutonomousChildrenCount),
		},
	}
	if config.SecondaryAutonomousShape != "" {
		secondaryX := configuredOrDefault(config.SecondaryAutonomousX, worldWidth/2+140)
		secondaryY := configuredOrDefault(config.SecondaryAutonomousY, worldHeight/2)
		autonomousCircles = append(autonomousCircles, AutonomousCircle{
			ID:               DefaultSecondaryID,
			LineageID:        lineageIDFor(DefaultSecondaryID),
			Generation:       0,
			Shape:            config.SecondaryAutonomousShape,
			X:                secondaryX,
			Y:                secondaryY,
			Radius:           DefaultPlayerRadius,
			Energy:           config.SecondaryAutonomousEnergy,
			AttachedChildren: initialAttachedChildren(DefaultSecondaryID, config.SecondaryChildrenCount),
		})
	}
	if config.UseExpandedPopulation {
		reserved := []Vector{
			{X: playerX, Y: playerY},
			{X: autonomousX, Y: autonomousY},
		}
		if config.SecondaryAutonomousShape != "" {
			reserved = append(reserved, Vector{
				X: configuredOrDefault(config.SecondaryAutonomousX, worldWidth/2+140),
				Y: configuredOrDefault(config.SecondaryAutonomousY, worldHeight/2),
			})
		}
		autonomousCircles = append(autonomousCircles, defaultExpandedAutonomousCircles(worldWidth, worldHeight, reserved)...)
	}

	activeCircleCount := len(autonomousCircles)
	if config.PlayerEnergy > 0 {
		activeCircleCount++
	}
	foodSlots := defaultFoodSlots(worldWidth, worldHeight, config.UseExpandedPopulation, activeCircleCount)
	player := &PlayerCircle{
		ID:               playerID,
		LineageID:        lineageIDFor(playerID),
		Generation:       0,
		Shape:            config.PlayerShape,
		X:                playerX,
		Y:                playerY,
		Radius:           DefaultPlayerRadius,
		Energy:           config.PlayerEnergy,
		AttachedChildren: initialAttachedChildren(playerID, config.PlayerChildrenCount),
	}

	return &World{
		bounds: Bounds{
			Width:  worldWidth,
			Height: worldHeight,
		},
		player:                              player,
		autonomousCircles:                   autonomousCircles,
		autonomousDirections:                initialAutonomousDirections(len(autonomousCircles)),
		nextChildID:                         totalInitialChildren(player, autonomousCircles) + 1,
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

func defaultExpandedAutonomousCircles(worldWidth, worldHeight float64, reserved []Vector) []AutonomousCircle {
	specs := []struct {
		id     string
		shape  string
		energy float64
	}{
		{id: DefaultTertiaryID},
		{id: DefaultQuaternaryID},
		{id: DefaultQuinaryID},
		{id: DefaultSenaryID},
		{id: DefaultSeptenaryID},
		{id: DefaultOctonaryID},
	}
	stateMix := seededExpandedAutonomousStateMix(len(specs))
	for index := range specs {
		specs[index].shape = stateMix[index].shape
		specs[index].energy = stateMix[index].energy
	}

	positions := seededExpandedAutonomousPositions(worldWidth, worldHeight, len(specs), reserved)
	circles := make([]AutonomousCircle, 0, len(specs))
	for index, spec := range specs {
		circles = append(circles, AutonomousCircle{
			ID:               spec.id,
			LineageID:        lineageIDFor(spec.id),
			Generation:       0,
			Shape:            spec.shape,
			X:                positions[index].X,
			Y:                positions[index].Y,
			Radius:           DefaultPlayerRadius,
			Energy:           spec.energy,
			AttachedChildren: initialAttachedChildren(spec.id, 0),
		})
	}

	return circles
}

func seededExpandedAutonomousStateMix(count int) []struct {
	shape  string
	energy float64
} {
	seed := DefaultExpandedAutonomousStateSeed
	mix := make([]struct {
		shape  string
		energy float64
	}, 0, count)

	advanceSeed := func() int64 {
		seed = (seed*214013 + 2531011) & 0x7fffffff
		return seed
	}

	for range count {
		shape := DefaultPlayerShape
		if advanceSeed()%2 == 1 {
			shape = DefaultAutoShape
		}
		energy := 82 + float64(advanceSeed()%15)
		mix = append(mix, struct {
			shape  string
			energy float64
		}{
			shape:  shape,
			energy: energy,
		})
	}

	return mix
}

func seededExpandedAutonomousPositions(worldWidth, worldHeight float64, count int, reserved []Vector) []Vector {
	seed := DefaultExpandedAutonomousSeed
	margin := 120.0
	minDistance := 180.0
	positions := make([]Vector, 0, count)
	occupied := append([]Vector(nil), reserved...)

	advanceSeed := func() int64 {
		seed = (seed*1103515245 + 12345) & 0x7fffffff
		return seed
	}

	nextCoordinate := func(size float64) float64 {
		usable := math.Max(1, size-margin*2)
		return margin + math.Mod(float64(advanceSeed()), usable)
	}

	for index := 0; index < count; index++ {
		var candidate Vector
		found := false
		for attempt := 0; attempt < 400; attempt++ {
			candidate = Vector{
				X: math.Round(nextCoordinate(worldWidth)),
				Y: math.Round(nextCoordinate(worldHeight)),
			}

			tooClose := false
			for _, other := range occupied {
				if math.Hypot(candidate.X-other.X, candidate.Y-other.Y) < minDistance {
					tooClose = true
					break
				}
			}
			if tooClose {
				continue
			}

			found = true
			break
		}

		if !found {
			candidate = Vector{
				X: margin + float64((index*173)%int(math.Max(1, worldWidth-margin*2))),
				Y: margin + float64((index*257)%int(math.Max(1, worldHeight-margin*2))),
			}
		}

		positions = append(positions, candidate)
		occupied = append(occupied, candidate)
	}

	return positions
}

func defaultFoodSlots(worldWidth, worldHeight float64, expanded bool, activeCircleCount int) []Food {
	centerX := worldWidth / 2
	centerY := worldHeight / 2

	narrowSlots := []Food{
		{ID: "food-1", X: centerX + 32, Y: centerY, Radius: DefaultFoodRadius},
		{ID: "food-2", X: centerX - 108, Y: centerY, Radius: DefaultFoodRadius},
		{ID: "food-3", X: centerX + 120, Y: centerY + 84, Radius: DefaultFoodRadius},
	}
	if !expanded {
		return append([]Food(nil), narrowSlots...)
	}

	slotCount := activeCircleCount + 2
	if slotCount < DefaultExpandedFoodCount {
		slotCount = DefaultExpandedFoodCount
	}
	return seededExpandedFoodSlots(worldWidth, worldHeight, slotCount)
}

func seededExpandedFoodSlots(worldWidth, worldHeight float64, slotCount int) []Food {
	seed := DefaultExpandedFoodSeed
	margin := 96.0
	centerX := worldWidth / 2
	centerY := worldHeight / 2
	minDistanceBetweenFoods := 72.0
	minDistanceFromCenter := 92.0

	advanceSeed := func() int64 {
		seed = (seed*1664525 + 1013904223) & 0x7fffffff
		return seed
	}

	nextCoordinate := func(size float64) float64 {
		usable := math.Max(1, size-margin*2)
		return margin + math.Mod(float64(advanceSeed()), usable)
	}

	slots := make([]Food, 0, slotCount)
	for index := 0; index < slotCount; index++ {
		var x float64
		var y float64
		found := false
		for attempt := 0; attempt < 400; attempt++ {
			x = math.Round(nextCoordinate(worldWidth))
			y = math.Round(nextCoordinate(worldHeight))

			if math.Hypot(x-centerX, y-centerY) < minDistanceFromCenter {
				continue
			}

			tooClose := false
			for _, existing := range slots {
				if math.Hypot(existing.X-x, existing.Y-y) < minDistanceBetweenFoods {
					tooClose = true
					break
				}
			}
			if tooClose {
				continue
			}

			found = true
			break
		}

		if !found {
			x = margin + float64((index*137)%int(math.Max(1, worldWidth-margin*2)))
			y = margin + float64((index*211)%int(math.Max(1, worldHeight-margin*2)))
		}

		slots = append(slots, Food{
			ID:     "food-" + strconv.Itoa(index+1),
			X:      x,
			Y:      y,
			Radius: DefaultFoodRadius,
		})
	}

	return slots
}

func (w *World) Advance(tick int64, intent Vector) Snapshot {
	w.lastInteraction = nil
	w.player = w.advanceCircle(w.player, intent)
	w.advanceAutonomousCircles(tick)

	w.consumeOverlappingFood(tick)
	w.regenerateFood(tick)
	w.resolveEnergyCollapse(tick)
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
		if !tracked || tick-missingSince < w.foodRegenDelay(slot) {
			continue
		}

		w.foods = append(w.foods, slot)
		delete(w.missingFoodSince, slot.ID)
	}
}

func (w *World) foodRegenDelay(slot Food) int64 {
	missingCount := len(w.missingFoodSince)
	delay := DefaultFoodRegenDelay
	if missingCount <= 1 {
		return delay
	}

	delay += int64(missingCount-1) * 2
	delay += int64(w.localMissingFoodCount(slot.ID)) * 2
	return delay
}

func (w *World) localMissingFoodCount(slotID string) int {
	slot, found := w.foodSlotByID(slotID)
	if !found {
		return 0
	}

	nearbyMissing := 0
	for missingID := range w.missingFoodSince {
		if missingID == slotID {
			continue
		}
		other, ok := w.foodSlotByID(missingID)
		if !ok {
			continue
		}
		if distanceBetween(slot.X, slot.Y, other.X, other.Y) <= DefaultRegionalFoodPressureDistance {
			nearbyMissing++
		}
	}

	return nearbyMissing
}

func (w *World) foodSlotByID(id string) (Food, bool) {
	for _, slot := range w.foodSlots {
		if slot.ID == id {
			return slot, true
		}
	}
	return Food{}, false
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
	circle.Energy = math.Max(0, circle.Energy-w.playerCrowdingCost(*circle))

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
		circle.Energy = math.Max(0, circle.Energy-w.autonomousCrowdingCost(circle))
		if circle.X == DefaultPlayerRadius || circle.X == w.bounds.Width-DefaultPlayerRadius {
			w.autonomousDirections[index] = Vector{X: -intent.X, Y: intent.Y}
		}
		w.autonomousCircles[index] = circle
	}
}

func (w *World) playerCrowdingCost(circle PlayerCircle) float64 {
	neighbors := 0
	for _, other := range w.autonomousCircles {
		if other.Energy <= 0 {
			continue
		}
		if distanceBetween(circle.X, circle.Y, other.X, other.Y) <= DefaultCrowdingDistance {
			neighbors++
		}
	}
	if neighbors < DefaultCrowdingThreshold {
		return 0
	}

	return DefaultCrowdingMoveCost
}

func (w *World) autonomousCrowdingCost(circle AutonomousCircle) float64 {
	neighbors := 0
	if w.player != nil && w.player.Energy > 0 && distanceBetween(circle.X, circle.Y, w.player.X, w.player.Y) <= DefaultCrowdingDistance {
		neighbors++
	}
	for _, other := range w.autonomousCircles {
		if other.ID == circle.ID || other.Energy <= 0 {
			continue
		}
		if distanceBetween(circle.X, circle.Y, other.X, other.Y) <= DefaultCrowdingDistance {
			neighbors++
		}
	}
	if neighbors < DefaultCrowdingThreshold {
		return 0
	}

	return DefaultCrowdingMoveCost
}

func (w *World) autonomousIntent(circle AutonomousCircle, index int, tick int64) Vector {
	if w.disableFoodSeeking {
		return w.autonomousDirections[index]
	}

	foodTarget, foodDistance, foodFound := nearestFoodTarget(circle, w.foods, tick)
	if foodFound && circle.Energy < DefaultLowEnergyFoodThreshold {
		return w.adjustForCrowding(circle, Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		})
	}

	threatTarget, threatFound := nearestThreatTarget(circle, w.player, w.autonomousCircles, tick)
	if !w.disableThreatAvoidance && threatFound {
		return w.adjustForCrowding(circle, Vector{
			X: circle.X - threatTarget.X,
			Y: circle.Y - threatTarget.Y,
		})
	}

	blockedTarget, blockedFound := nearestBlockedReproductionTarget(circle, w.player, w.autonomousCircles, tick)
	if !w.disableBlockedReproductionAvoidance && blockedFound {
		return w.adjustForCrowding(circle, Vector{
			X: circle.X - blockedTarget.X,
			Y: circle.Y - blockedTarget.Y,
		})
	}

	if foodFound && foodDistance <= DefaultFoodPriorityDistance {
		return w.adjustForCrowding(circle, Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		})
	}

	interactionTarget, _, interactionFound := nearestInteractionTarget(circle, w.player, w.autonomousCircles, tick)
	if interactionFound {
		return w.adjustForCrowding(circle, Vector{
			X: interactionTarget.X - circle.X,
			Y: interactionTarget.Y - circle.Y,
		})
	}

	if foodFound {
		return w.adjustForCrowding(circle, Vector{
			X: foodTarget.X - circle.X,
			Y: foodTarget.Y - circle.Y,
		})
	}

	return w.adjustForCrowding(circle, w.autonomousDirections[index])
}

func (w *World) adjustForCrowding(circle AutonomousCircle, intent Vector) Vector {
	normalized := normalize(intent)
	if normalized.X == 0 && normalized.Y == 0 {
		return intent
	}

	currentNeighbors := w.localCrowdingNeighbors(circle.X, circle.Y, circle.ID)
	nextX := clamp(circle.X+normalized.X*w.speed, DefaultPlayerRadius, w.bounds.Width-DefaultPlayerRadius)
	nextY := clamp(circle.Y+normalized.Y*w.speed, DefaultPlayerRadius, w.bounds.Height-DefaultPlayerRadius)
	nextNeighbors := w.localCrowdingNeighbors(nextX, nextY, circle.ID)

	if nextNeighbors > currentNeighbors+1 {
		return Vector{
			X: -normalized.X,
			Y: -normalized.Y,
		}
	}

	return intent
}

func (w *World) localCrowdingNeighbors(x, y float64, selfID string) int {
	neighbors := 0
	if w.player != nil && w.player.Energy > 0 && w.player.ID != selfID && distanceBetween(x, y, w.player.X, w.player.Y) <= DefaultCrowdingDistance {
		neighbors++
	}
	for _, other := range w.autonomousCircles {
		if other.ID == selfID || other.Energy <= 0 {
			continue
		}
		if distanceBetween(x, y, other.X, other.Y) <= DefaultCrowdingDistance {
			neighbors++
		}
	}

	return neighbors
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
			contactDetails, interacting := circlesInteract(*w.player, circle, tick)
			if !interacting {
				continue
			}

			pairKey := overlapPairKey(w.player.ID, circle.ID)
			currentOverlapPairs[pairKey] = struct{}{}

			if w.player.Shape != circle.Shape {
				if _, exists := w.activeOverlapPairs[pairKey]; exists {
					continue
				}

				w.resolveReproduction(circle.ID, tick, contactDetails)
				w.activeOverlapPairs = currentOverlapPairs
				return
			}

			w.resolveFight(circle.ID, contactDetails, tick)
			w.activeOverlapPairs = currentOverlapPairs
			return
		}
	}

	for left := 0; left < len(w.autonomousCircles); left++ {
		for right := left + 1; right < len(w.autonomousCircles); right++ {
			leftCircle := w.autonomousCircles[left]
			rightCircle := w.autonomousCircles[right]
			contactDetails, interacting := autonomousCirclesInteract(leftCircle, rightCircle, tick)
			if !interacting {
				continue
			}

			pairKey := overlapPairKey(leftCircle.ID, rightCircle.ID)
			currentOverlapPairs[pairKey] = struct{}{}

			if leftCircle.Shape != rightCircle.Shape {
				if _, exists := w.activeOverlapPairs[pairKey]; exists {
					continue
				}

				w.resolveAutonomousReproduction(left, right, tick, contactDetails)
				w.activeOverlapPairs = currentOverlapPairs
				return
			}

			w.resolveAutonomousFight(left, right, contactDetails, tick)
			w.activeOverlapPairs = currentOverlapPairs
			return
		}
	}

	w.activeOverlapPairs = currentOverlapPairs
}

func circlesInteract(player PlayerCircle, autonomous AutonomousCircle, tick int64) (ContactDetails, bool) {
	parentBodyOverlap := overlaps(player.X, player.Y, DefaultPlayerRadius, autonomous.X, autonomous.Y, DefaultPlayerRadius)

	playerChildren := layoutAttachedChildren(player.ID, player.X, player.Y, player.Radius, player.AttachedChildren, tick)
	autonomousChildren := layoutAttachedChildren(autonomous.ID, autonomous.X, autonomous.Y, autonomous.Radius, autonomous.AttachedChildren, tick)

	for _, child := range playerChildren {
		if overlaps(child.X, child.Y, child.Radius, autonomous.X, autonomous.Y, DefaultPlayerRadius) {
			return ContactDetails{Origin: "attached_child", PathKind: "source_child_to_target_parent", SourceChildID: child.ID}, true
		}
	}

	for _, child := range autonomousChildren {
		if overlaps(player.X, player.Y, DefaultPlayerRadius, child.X, child.Y, child.Radius) {
			return ContactDetails{Origin: "attached_child", PathKind: "source_parent_to_target_child", TargetChildID: child.ID}, true
		}
	}

	for _, playerChild := range playerChildren {
		for _, autonomousChild := range autonomousChildren {
			if overlaps(playerChild.X, playerChild.Y, playerChild.Radius, autonomousChild.X, autonomousChild.Y, autonomousChild.Radius) {
				return ContactDetails{Origin: "attached_child", PathKind: "child_to_child", SourceChildID: playerChild.ID, TargetChildID: autonomousChild.ID}, true
			}
		}
	}

	if parentBodyOverlap {
		return ContactDetails{Origin: "parent_body"}, true
	}

	return ContactDetails{}, false
}

func autonomousCirclesInteract(left AutonomousCircle, right AutonomousCircle, tick int64) (ContactDetails, bool) {
	parentBodyOverlap := overlaps(left.X, left.Y, DefaultPlayerRadius, right.X, right.Y, DefaultPlayerRadius)

	leftChildren := layoutAttachedChildren(left.ID, left.X, left.Y, left.Radius, left.AttachedChildren, tick)
	rightChildren := layoutAttachedChildren(right.ID, right.X, right.Y, right.Radius, right.AttachedChildren, tick)

	for _, child := range leftChildren {
		if overlaps(child.X, child.Y, child.Radius, right.X, right.Y, DefaultPlayerRadius) {
			return ContactDetails{Origin: "attached_child", PathKind: "source_child_to_target_parent", SourceChildID: child.ID}, true
		}
	}

	for _, child := range rightChildren {
		if overlaps(left.X, left.Y, DefaultPlayerRadius, child.X, child.Y, child.Radius) {
			return ContactDetails{Origin: "attached_child", PathKind: "source_parent_to_target_child", TargetChildID: child.ID}, true
		}
	}

	for _, leftChild := range leftChildren {
		for _, rightChild := range rightChildren {
			if overlaps(leftChild.X, leftChild.Y, leftChild.Radius, rightChild.X, rightChild.Y, rightChild.Radius) {
				return ContactDetails{Origin: "attached_child", PathKind: "child_to_child", SourceChildID: leftChild.ID, TargetChildID: rightChild.ID}, true
			}
		}
	}

	if parentBodyOverlap {
		return ContactDetails{Origin: "parent_body"}, true
	}

	return ContactDetails{}, false
}

func (w *World) resolveFight(opponentID string, contactDetails ContactDetails, tick int64) {
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
		Active:          false,
		Resolved:        true,
		Kind:            "",
		ContactOrigin:   contactDetails.Origin,
		ContactPathKind: contactDetails.PathKind,
		SourceID:        w.player.ID,
		TargetID:        opponent.ID,
		SourceChildID:   contactDetails.SourceChildID,
		TargetChildID:   contactDetails.TargetChildID,
		WinnerID:        winnerID,
		LoserID:         loserID,
	}

	if loserID == w.player.ID {
		if childCountForPlayer(*w.player) > 0 {
			absorbedChildID, _ := consumedPlayerChildID(w.player)
			consumePlayerChild(w.player)
			w.lastInteraction.Kind = "fight_absorbed_child"
			w.lastInteraction.AbsorbedChildID = absorbedChildID
			return
		}
		w.lastInteraction.Kind = "fight_resolved"
		w.player = replaceOrRemovePlayer(w.player, tick)
		return
	}

	if childCountForAutonomous(opponent) > 0 {
		absorbedChildID, _ := consumedAutonomousChildID(opponent)
		consumeAutonomousChild(&opponent)
		w.autonomousCircles[opponentIndex] = opponent
		w.lastInteraction.Kind = "fight_absorbed_child"
		w.lastInteraction.AbsorbedChildID = absorbedChildID
		return
	}

	w.lastInteraction.Kind = "fight_resolved"
	replacedOpponent, active := replaceOrRemoveAutonomous(opponent, tick)
	if !active {
		w.autonomousCircles = append(w.autonomousCircles[:opponentIndex], w.autonomousCircles[opponentIndex+1:]...)
		w.autonomousDirections = append(w.autonomousDirections[:opponentIndex], w.autonomousDirections[opponentIndex+1:]...)
		return
	}

	w.autonomousCircles[opponentIndex] = replacedOpponent
}

func (w *World) resolveReproduction(opponentID string, tick int64, contactDetails ContactDetails) {
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
	playerCapacity := reproductionCapacityForPlayer(*w.player)
	opponentCapacity := reproductionCapacityForAutonomous(opponent)
	playerEnergyComponent, playerReserveComponent := reproductionCapacityComponents(w.player.Energy, childCountForPlayer(*w.player))
	opponentEnergyComponent, opponentReserveComponent := reproductionCapacityComponents(opponent.Energy, childCountForAutonomous(opponent))
	playerPaid, playerUsedChild, playerPaidChildID := payPlayerReproductionCost(w.player)
	opponentPaid, opponentUsedChild, opponentPaidChildID, opponent := payAutonomousReproductionCost(opponent)
	if !playerPaid || !opponentPaid {
		w.lastInteraction = &InteractionClassification{
			Active:                 false,
			Resolved:               true,
			Kind:                   "reproduce_blocked_energy",
			ContactOrigin:          contactDetails.Origin,
			ContactPathKind:        contactDetails.PathKind,
			SourceID:               w.player.ID,
			TargetID:               opponent.ID,
			SourceChildID:          contactDetails.SourceChildID,
			TargetChildID:          contactDetails.TargetChildID,
			SourceBlockedCapacity:  !playerPaid,
			TargetBlockedCapacity:  !opponentPaid,
			SourceCapacityValue:    playerCapacity,
			TargetCapacityValue:    opponentCapacity,
			SourceEnergyComponent:  playerEnergyComponent,
			TargetEnergyComponent:  opponentEnergyComponent,
			SourceReserveComponent: playerReserveComponent,
			TargetReserveComponent: opponentReserveComponent,
			ReproductionThreshold:  DefaultReproductionMinEnergy,
			ReproductionCost:       DefaultReproductionCost,
		}
		return
	}

	distribution := reproductionDistributionCase(tick, *w.player, opponent)
	createdChildIDs, sourceCreatedChildIDs, targetCreatedChildIDs := w.assignReproductionChildren(&opponent, distribution)
	w.autonomousCircles[opponentIndex] = opponent

	kind := "reproduce_resolved"
	if playerUsedChild || opponentUsedChild {
		kind = "reproduce_paid_child"
	}

	w.lastInteraction = &InteractionClassification{
		Active:                 false,
		Resolved:               true,
		Kind:                   kind,
		ContactOrigin:          contactDetails.Origin,
		ContactPathKind:        contactDetails.PathKind,
		DistributionKind:       reproductionDistributionKind(distribution),
		SourceID:               w.player.ID,
		TargetID:               opponent.ID,
		SourceChildID:          contactDetails.SourceChildID,
		TargetChildID:          contactDetails.TargetChildID,
		SourceCapacityValue:    playerCapacity,
		TargetCapacityValue:    opponentCapacity,
		SourceEnergyComponent:  playerEnergyComponent,
		TargetEnergyComponent:  opponentEnergyComponent,
		SourceReserveComponent: playerReserveComponent,
		TargetReserveComponent: opponentReserveComponent,
		ReproductionThreshold:  DefaultReproductionMinEnergy,
		ReproductionCost:       DefaultReproductionCost,
		SourcePaidChild:        playerUsedChild,
		TargetPaidChild:        opponentUsedChild,
		SourcePaidChildID:      playerPaidChildID,
		TargetPaidChildID:      opponentPaidChildID,
		CreatedChildIDs:        createdChildIDs,
		SourceCreatedChildIDs:  sourceCreatedChildIDs,
		TargetCreatedChildIDs:  targetCreatedChildIDs,
	}
}

func determineAutonomousFightOutcome(left AutonomousCircle, right AutonomousCircle) (string, string) {
	if left.Energy > right.Energy {
		return left.ID, right.ID
	}
	if right.Energy > left.Energy {
		return right.ID, left.ID
	}
	if hasAttachedChildrenAutonomous(left) && !hasAttachedChildrenAutonomous(right) {
		return left.ID, right.ID
	}
	if hasAttachedChildrenAutonomous(right) && !hasAttachedChildrenAutonomous(left) {
		return right.ID, left.ID
	}
	if left.ID < right.ID {
		return left.ID, right.ID
	}

	return right.ID, left.ID
}

func (w *World) resolveAutonomousFight(leftIndex int, rightIndex int, contactDetails ContactDetails, tick int64) {
	if leftIndex < 0 || rightIndex < 0 || leftIndex >= len(w.autonomousCircles) || rightIndex >= len(w.autonomousCircles) {
		return
	}

	leftCircle := w.autonomousCircles[leftIndex]
	rightCircle := w.autonomousCircles[rightIndex]
	winnerID, loserID := determineAutonomousFightOutcome(leftCircle, rightCircle)

	w.lastInteraction = &InteractionClassification{
		Active:          false,
		Resolved:        true,
		Kind:            "",
		ContactOrigin:   contactDetails.Origin,
		ContactPathKind: contactDetails.PathKind,
		SourceID:        leftCircle.ID,
		TargetID:        rightCircle.ID,
		SourceChildID:   contactDetails.SourceChildID,
		TargetChildID:   contactDetails.TargetChildID,
		WinnerID:        winnerID,
		LoserID:         loserID,
	}

	loserIndex := leftIndex
	if loserID == rightCircle.ID {
		loserIndex = rightIndex
	}
	loser := w.autonomousCircles[loserIndex]
	if childCountForAutonomous(loser) > 0 {
		absorbedChildID, _ := consumedAutonomousChildID(loser)
		consumeAutonomousChild(&loser)
		w.autonomousCircles[loserIndex] = loser
		w.lastInteraction.Kind = "fight_absorbed_child"
		w.lastInteraction.AbsorbedChildID = absorbedChildID
		return
	}

	w.lastInteraction.Kind = "fight_resolved"
	replacedLoser, active := replaceOrRemoveAutonomous(loser, tick)
	if active {
		w.autonomousCircles[loserIndex] = replacedLoser
		return
	}

	w.autonomousCircles = append(w.autonomousCircles[:loserIndex], w.autonomousCircles[loserIndex+1:]...)
	w.autonomousDirections = append(w.autonomousDirections[:loserIndex], w.autonomousDirections[loserIndex+1:]...)
}

func (w *World) resolveAutonomousReproduction(leftIndex int, rightIndex int, tick int64, contactDetails ContactDetails) {
	if leftIndex < 0 || rightIndex < 0 || leftIndex >= len(w.autonomousCircles) || rightIndex >= len(w.autonomousCircles) {
		return
	}

	leftCircle := w.autonomousCircles[leftIndex]
	rightCircle := w.autonomousCircles[rightIndex]
	leftCapacity := reproductionCapacityForAutonomous(leftCircle)
	rightCapacity := reproductionCapacityForAutonomous(rightCircle)
	leftEnergyComponent, leftReserveComponent := reproductionCapacityComponents(leftCircle.Energy, childCountForAutonomous(leftCircle))
	rightEnergyComponent, rightReserveComponent := reproductionCapacityComponents(rightCircle.Energy, childCountForAutonomous(rightCircle))
	leftPaid, leftUsedChild, leftPaidChildID, leftCircle := payAutonomousReproductionCost(leftCircle)
	rightPaid, rightUsedChild, rightPaidChildID, rightCircle := payAutonomousReproductionCost(rightCircle)
	if !leftPaid || !rightPaid {
		w.lastInteraction = &InteractionClassification{
			Active:                 false,
			Resolved:               true,
			Kind:                   "reproduce_blocked_energy",
			ContactOrigin:          contactDetails.Origin,
			ContactPathKind:        contactDetails.PathKind,
			SourceID:               leftCircle.ID,
			TargetID:               rightCircle.ID,
			SourceChildID:          contactDetails.SourceChildID,
			TargetChildID:          contactDetails.TargetChildID,
			SourceBlockedCapacity:  !leftPaid,
			TargetBlockedCapacity:  !rightPaid,
			SourceCapacityValue:    leftCapacity,
			TargetCapacityValue:    rightCapacity,
			SourceEnergyComponent:  leftEnergyComponent,
			TargetEnergyComponent:  rightEnergyComponent,
			SourceReserveComponent: leftReserveComponent,
			TargetReserveComponent: rightReserveComponent,
			ReproductionThreshold:  DefaultReproductionMinEnergy,
			ReproductionCost:       DefaultReproductionCost,
		}
		return
	}

	distribution := autonomousReproductionDistributionCase(tick, leftCircle, rightCircle)
	createdChildIDs, sourceCreatedChildIDs, targetCreatedChildIDs := w.assignAutonomousPairReproductionChildren(&leftCircle, &rightCircle, distribution)
	w.autonomousCircles[leftIndex] = leftCircle
	w.autonomousCircles[rightIndex] = rightCircle

	kind := "reproduce_resolved"
	if leftUsedChild || rightUsedChild {
		kind = "reproduce_paid_child"
	}

	w.lastInteraction = &InteractionClassification{
		Active:                 false,
		Resolved:               true,
		Kind:                   kind,
		ContactOrigin:          contactDetails.Origin,
		ContactPathKind:        contactDetails.PathKind,
		DistributionKind:       reproductionDistributionKind(distribution),
		SourceID:               leftCircle.ID,
		TargetID:               rightCircle.ID,
		SourceChildID:          contactDetails.SourceChildID,
		TargetChildID:          contactDetails.TargetChildID,
		SourceCapacityValue:    leftCapacity,
		TargetCapacityValue:    rightCapacity,
		SourceEnergyComponent:  leftEnergyComponent,
		TargetEnergyComponent:  rightEnergyComponent,
		SourceReserveComponent: leftReserveComponent,
		TargetReserveComponent: rightReserveComponent,
		ReproductionThreshold:  DefaultReproductionMinEnergy,
		ReproductionCost:       DefaultReproductionCost,
		SourcePaidChild:        leftUsedChild,
		TargetPaidChild:        rightUsedChild,
		SourcePaidChildID:      leftPaidChildID,
		TargetPaidChildID:      rightPaidChildID,
		CreatedChildIDs:        createdChildIDs,
		SourceCreatedChildIDs:  sourceCreatedChildIDs,
		TargetCreatedChildIDs:  targetCreatedChildIDs,
	}
}

func determineFightOutcome(player PlayerCircle, opponent AutonomousCircle) (string, string) {
	if player.Energy > opponent.Energy {
		return player.ID, opponent.ID
	}
	if opponent.Energy > player.Energy {
		return opponent.ID, player.ID
	}
	if hasAttachedChildrenPlayer(player) && !hasAttachedChildrenAutonomous(opponent) {
		return player.ID, opponent.ID
	}
	if hasAttachedChildrenAutonomous(opponent) && !hasAttachedChildrenPlayer(player) {
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

func payPlayerReproductionCost(circle *PlayerCircle) (bool, bool, string) {
	if circle == nil {
		return false, false, ""
	}
	if reproductionCapacityForPlayer(*circle) < DefaultReproductionMinEnergy {
		return false, false, ""
	}
	if circle.Energy >= DefaultReproductionCost {
		circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
		return true, false, ""
	}
	if childCountForPlayer(*circle) == 0 {
		return false, false, ""
	}

	paidChildID, _ := consumedPlayerChildID(circle)
	consumePlayerChild(circle)
	circle.Energy += DefaultReproductionCost
	circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
	return true, true, paidChildID
}

func payAutonomousReproductionCost(circle AutonomousCircle) (bool, bool, string, AutonomousCircle) {
	if reproductionCapacityForAutonomous(circle) < DefaultReproductionMinEnergy {
		return false, false, "", circle
	}
	if circle.Energy >= DefaultReproductionCost {
		circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
		return true, false, "", circle
	}
	if childCountForAutonomous(circle) == 0 {
		return false, false, "", circle
	}

	paidChildID, _ := consumedAutonomousChildID(circle)
	consumeAutonomousChild(&circle)
	circle.Energy += DefaultReproductionCost
	circle.Energy = math.Max(0, circle.Energy-DefaultReproductionCost)
	return true, true, paidChildID, circle
}

func reproductionCapacityForPlayer(circle PlayerCircle) float64 {
	return reproductionCapacity(circle.Energy, childCountForPlayer(circle))
}

func reproductionCapacityForAutonomous(circle AutonomousCircle) float64 {
	return reproductionCapacity(circle.Energy, childCountForAutonomous(circle))
}

func reproductionCapacity(energy float64, childrenCount int) float64 {
	if childrenCount == 0 {
		return energy
	}

	return energy + DefaultReproductionCost
}

func reproductionCapacityComponents(energy float64, childrenCount int) (float64, float64) {
	if childrenCount == 0 {
		return energy, 0
	}

	return energy, DefaultReproductionCost
}

func (w *World) resolveEnergyCollapse(tick int64) {
	if w.player != nil && w.player.Energy == 0 {
		promotedChildID, promoted := continuityPromotedChildIDForPlayer(w.player, tick)
		w.player = replaceOrRemovePlayer(w.player, tick)
		if promoted && w.player != nil && w.lastInteraction == nil {
			w.lastInteraction = &InteractionClassification{
				Active:          false,
				Resolved:        true,
				Kind:            "death_promoted_child",
				SourceID:        w.player.ID,
				TargetID:        w.player.ID,
				PromotedChildID: promotedChildID,
			}
		}
	}

	for index, circle := range w.autonomousCircles {
		if circle.Energy != 0 {
			continue
		}

		promotedChildID, promoted := continuityPromotedChildIDForAutonomous(circle, tick)
		replaced, active := replaceOrRemoveAutonomous(circle, tick)
		if !active {
			w.autonomousCircles = append(w.autonomousCircles[:index], w.autonomousCircles[index+1:]...)
			w.autonomousDirections = append(w.autonomousDirections[:index], w.autonomousDirections[index+1:]...)
			index--
			continue
		}

		w.autonomousCircles[index] = replaced
		if promoted && w.lastInteraction == nil {
			w.lastInteraction = &InteractionClassification{
				Active:          false,
				Resolved:        true,
				Kind:            "death_promoted_child",
				SourceID:        replaced.ID,
				TargetID:        replaced.ID,
				PromotedChildID: promotedChildID,
			}
		}
	}
}

func replaceOrRemovePlayer(circle *PlayerCircle, tick int64) *PlayerCircle {
	if circle == nil {
		return nil
	}
	if childCountForPlayer(*circle) == 0 {
		return nil
	}

	promotedX, promotedY, ok := promotedChildPosition(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	consumePlayerChild(circle)
	if ok {
		circle.X = promotedX
		circle.Y = promotedY
	}
	circle.Generation++
	circle.Energy = DefaultReplacementEnergy

	return circle
}

func replaceOrRemoveAutonomous(circle AutonomousCircle, tick int64) (AutonomousCircle, bool) {
	if childCountForAutonomous(circle) == 0 {
		return AutonomousCircle{}, false
	}

	promotedX, promotedY, ok := promotedChildPosition(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	consumeAutonomousChild(&circle)
	if ok {
		circle.X = promotedX
		circle.Y = promotedY
	}
	circle.Generation++
	circle.Energy = DefaultReplacementEnergy

	return circle, true
}

func promotedChildPosition(ownerID string, x, y, parentRadius float64, children []AttachedChild, tick int64) (float64, float64, bool) {
	_, promotedX, promotedY, ok := promotedChildSelection(ownerID, x, y, parentRadius, children, tick)
	return promotedX, promotedY, ok
}

func continuityPromotedChildIDForPlayer(circle *PlayerCircle, tick int64) (string, bool) {
	if circle == nil {
		return "", false
	}
	return continuityPromotedChildID(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
}

func continuityPromotedChildIDForAutonomous(circle AutonomousCircle, tick int64) (string, bool) {
	return continuityPromotedChildID(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
}

func continuityPromotedChildID(ownerID string, x, y, parentRadius float64, children []AttachedChild, tick int64) (string, bool) {
	promotedID, _, _, ok := promotedChildSelection(ownerID, x, y, parentRadius, children, tick)
	return promotedID, ok
}

func consumedPlayerChildID(circle *PlayerCircle) (string, bool) {
	if circle == nil || len(circle.AttachedChildren) == 0 {
		return "", false
	}
	return consumedChildID(circle.AttachedChildren), true
}

func consumedAutonomousChildID(circle AutonomousCircle) (string, bool) {
	if len(circle.AttachedChildren) == 0 {
		return "", false
	}
	return consumedChildID(circle.AttachedChildren), true
}

func consumedChildID(children []AttachedChild) string {
	return children[len(children)-1].ID
}

func promotedChildSelection(ownerID string, x, y, parentRadius float64, children []AttachedChild, tick int64) (string, float64, float64, bool) {
	if len(children) == 0 {
		return "", 0, 0, false
	}

	layout := layoutAttachedChildren(ownerID, x, y, parentRadius, children, tick)
	promoted := layout[len(layout)-1]
	return promoted.ID, promoted.X, promoted.Y, true
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

func totalInitialChildren(player *PlayerCircle, autonomousCircles []AutonomousCircle) int {
	total := 0
	if player != nil {
		total += len(player.AttachedChildren)
	}
	for _, circle := range autonomousCircles {
		total += len(circle.AttachedChildren)
	}

	return total
}

func snapshotPlayerCircle(circle PlayerCircle, tick int64) PlayerCircle {
	copy := circle
	copy.AttachedChildren = layoutAttachedChildren(circle.ID, circle.X, circle.Y, circle.Radius, circle.AttachedChildren, tick)
	return copy
}

func snapshotAutonomousCircle(circle AutonomousCircle, tick int64) AutonomousCircle {
	copy := circle
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

func (w *World) assignReproductionChildren(opponent *AutonomousCircle, distribution int) ([]string, []string, []string) {
	createdChildIDs := make([]string, 0, 2)
	sourceCreatedChildIDs := make([]string, 0, 2)
	targetCreatedChildIDs := make([]string, 0, 2)
	switch distribution {
	case 0:
		firstChildID := w.addPlayerChild()
		secondChildID := w.addPlayerChild()
		createdChildIDs = append(createdChildIDs, firstChildID, secondChildID)
		sourceCreatedChildIDs = append(sourceCreatedChildIDs, firstChildID, secondChildID)
	case 1:
		playerChildID := w.addPlayerChild()
		createdChildIDs = append(createdChildIDs, playerChildID)
		sourceCreatedChildIDs = append(sourceCreatedChildIDs, playerChildID)
		childID := w.allocateChildID()
		addAutonomousChild(opponent, childID)
		createdChildIDs = append(createdChildIDs, childID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, childID)
	case 2:
		firstChildID := w.allocateChildID()
		addAutonomousChild(opponent, firstChildID)
		createdChildIDs = append(createdChildIDs, firstChildID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, firstChildID)
		secondChildID := w.allocateChildID()
		addAutonomousChild(opponent, secondChildID)
		createdChildIDs = append(createdChildIDs, secondChildID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, secondChildID)
	}
	return createdChildIDs, sourceCreatedChildIDs, targetCreatedChildIDs
}

func reproductionDistributionCase(tick int64, player PlayerCircle, opponent AutonomousCircle) int {
	seed := int(tick) + len(player.ID) + len(opponent.ID) + player.Generation + opponent.Generation + childCountForPlayer(player) + childCountForAutonomous(opponent)
	return seed % 3
}

func reproductionDistributionKind(distribution int) string {
	switch distribution {
	case 0:
		return "source_only"
	case 1:
		return "split"
	case 2:
		return "target_only"
	default:
		return ""
	}
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

func (w *World) addPlayerChild() string {
	if w.player == nil {
		return ""
	}
	childID := w.allocateChildID()
	w.player.AttachedChildren = append(w.player.AttachedChildren, AttachedChild{
		ID:      childID,
		OwnerID: w.player.ID,
	})
	syncPlayerChildrenState(w.player)
	return childID
}

func addAutonomousChild(circle *AutonomousCircle, childID string) {
	circle.AttachedChildren = append(circle.AttachedChildren, AttachedChild{
		ID:      childID,
		OwnerID: circle.ID,
	})
	syncAutonomousChildrenState(circle)
}

func (w *World) assignAutonomousPairReproductionChildren(left *AutonomousCircle, right *AutonomousCircle, distribution int) ([]string, []string, []string) {
	createdChildIDs := make([]string, 0, 2)
	sourceCreatedChildIDs := make([]string, 0, 2)
	targetCreatedChildIDs := make([]string, 0, 2)
	switch distribution {
	case 0:
		firstChildID := w.allocateChildID()
		addAutonomousChild(left, firstChildID)
		createdChildIDs = append(createdChildIDs, firstChildID)
		sourceCreatedChildIDs = append(sourceCreatedChildIDs, firstChildID)
		secondChildID := w.allocateChildID()
		addAutonomousChild(left, secondChildID)
		createdChildIDs = append(createdChildIDs, secondChildID)
		sourceCreatedChildIDs = append(sourceCreatedChildIDs, secondChildID)
	case 1:
		leftChildID := w.allocateChildID()
		addAutonomousChild(left, leftChildID)
		createdChildIDs = append(createdChildIDs, leftChildID)
		sourceCreatedChildIDs = append(sourceCreatedChildIDs, leftChildID)
		rightChildID := w.allocateChildID()
		addAutonomousChild(right, rightChildID)
		createdChildIDs = append(createdChildIDs, rightChildID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, rightChildID)
	case 2:
		firstChildID := w.allocateChildID()
		addAutonomousChild(right, firstChildID)
		createdChildIDs = append(createdChildIDs, firstChildID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, firstChildID)
		secondChildID := w.allocateChildID()
		addAutonomousChild(right, secondChildID)
		createdChildIDs = append(createdChildIDs, secondChildID)
		targetCreatedChildIDs = append(targetCreatedChildIDs, secondChildID)
	}
	return createdChildIDs, sourceCreatedChildIDs, targetCreatedChildIDs
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

func hasAttachedChildrenPlayer(circle PlayerCircle) bool {
	return childCountForPlayer(circle) > 0
}

func hasAttachedChildrenAutonomous(circle AutonomousCircle) bool {
	return childCountForAutonomous(circle) > 0
}

func syncPlayerChildrenState(circle *PlayerCircle) {
	circle.Radius = DefaultPlayerRadius
}

func syncAutonomousChildrenState(circle *AutonomousCircle) {
	circle.Radius = DefaultPlayerRadius
}
