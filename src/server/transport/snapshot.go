package transport

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

const (
	DefaultViewportInterestWidth     = 1200.0
	DefaultViewportInterestHeight    = 840.0
	DefaultViewportInterestMargin    = 180.0
	DefaultOrientationEveryTicks     = int64(5)
	DefaultOrientationFallbackTicks  = int64(20)
	DefaultLocalAutonomousEveryTicks = int64(2)
	DefaultLocalFoodFallbackTicks    = int64(20)
	DefaultObserverFallbackTicks     = int64(12)
	DefaultMinimapClusterSize        = 480.0
	DefaultActiveMinimapClusterSize  = 960.0
)

type MinimapAutonomousCircle struct {
	Shape string  `json:"shape"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Count int     `json:"count"`
}

type MinimapFood struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Count int     `json:"count"`
}

type Snapshot struct {
	TransportMode            string                                `json:"transport_mode"`
	Type                     string                                `json:"type"`
	Tick                     int64                                 `json:"tick"`
	World                    simulation.Bounds                     `json:"world"`
	Player                   *simulation.PlayerCircle              `json:"player"`
	AutonomousCircles        []simulation.AutonomousCircle         `json:"autonomous_circles"`
	AutonomousFresh          bool                                  `json:"autonomous_fresh"`
	Interaction              *simulation.InteractionClassification `json:"interaction"`
	Foods                    []simulation.Food                     `json:"foods"`
	FoodsFresh               bool                                  `json:"foods_fresh"`
	OrientationFresh         bool                                  `json:"orientation_fresh"`
	MinimapAutonomousCircles []MinimapAutonomousCircle             `json:"minimap_autonomous_circles"`
	MinimapFoods             []MinimapFood                         `json:"minimap_foods"`
	TotalAutonomousCircles   int                                   `json:"total_autonomous_circles"`
	TotalFoods               int                                   `json:"total_foods"`
}

func BuildViewportSnapshot(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	return buildViewportSnapshot(snapshot, includeOrientation, true, true, DefaultActiveMinimapClusterSize)
}

func BuildObserverSnapshot(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	observerSnapshot := buildViewportSnapshot(snapshot, includeOrientation, true, true, DefaultMinimapClusterSize)
	observerSnapshot.TransportMode = "observer_orientation_only"
	observerSnapshot.Foods = []simulation.Food{}
	observerSnapshot.FoodsFresh = false
	return observerSnapshot
}

func BuildViewportSnapshotExactOrientation(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	return buildViewportSnapshot(snapshot, includeOrientation, false, true, DefaultMinimapClusterSize)
}

func BuildViewportSnapshotCompactFullPrecision(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	return buildViewportSnapshot(snapshot, includeOrientation, true, false, DefaultActiveMinimapClusterSize)
}

func OrientationSummarySignature(snapshot Snapshot) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(snapshot.TotalAutonomousCircles))
	builder.WriteByte('|')
	builder.WriteString(strconv.Itoa(snapshot.TotalFoods))

	for _, circle := range snapshot.MinimapAutonomousCircles {
		builder.WriteByte('|')
		builder.WriteString(circle.Shape)
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(circle.X, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(circle.Y, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(circle.Count))
	}

	for _, food := range snapshot.MinimapFoods {
		builder.WriteByte('|')
		builder.WriteString(strconv.FormatFloat(food.X, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(food.Y, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(food.Count))
	}

	return builder.String()
}

func ObserverTransportSignature(snapshot Snapshot) string {
	var builder strings.Builder
	builder.WriteString(snapshot.TransportMode)
	builder.WriteByte('|')
	builder.WriteString(strconv.Itoa(snapshot.TotalAutonomousCircles))
	builder.WriteByte(':')
	builder.WriteString(observerFoodStateBucket(snapshot.TotalFoods))
	builder.WriteByte('|')
	if snapshot.Interaction == nil {
		builder.WriteString("no-interaction")
		return builder.String()
	}

	builder.WriteString(snapshot.Interaction.Kind)
	builder.WriteByte(':')
	builder.WriteString(snapshot.Interaction.SourceID)
	builder.WriteByte(':')
	builder.WriteString(snapshot.Interaction.TargetID)
	builder.WriteByte(':')
	builder.WriteString(snapshot.Interaction.ContactOrigin)
	builder.WriteByte(':')
	builder.WriteString(snapshot.Interaction.ContactPathKind)
	return builder.String()
}

func ShouldRefreshOrientation(lastSignature string, lastRefreshTick, currentTick int64, currentSignature string) bool {
	if currentTick <= lastRefreshTick {
		return true
	}
	if lastSignature == "" {
		return true
	}
	if currentSignature != lastSignature {
		return true
	}
	return currentTick-lastRefreshTick >= DefaultOrientationFallbackTicks
}

func ShouldRefreshObserverTransport(lastSignature string, lastRefreshTick, currentTick int64, currentSignature string) bool {
	if currentTick <= lastRefreshTick {
		return true
	}
	if lastSignature == "" {
		return true
	}
	if currentSignature != lastSignature {
		return true
	}
	return currentTick-lastRefreshTick >= DefaultObserverFallbackTicks
}

func observerFoodStateBucket(totalFoods int) string {
	switch {
	case totalFoods <= 2:
		return "scarce"
	case totalFoods <= 8:
		return "thin"
	default:
		return "abundant"
	}
}

func LocalAutonomousSignature(snapshot Snapshot) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(snapshot.TotalAutonomousCircles))

	for _, circle := range snapshot.AutonomousCircles {
		builder.WriteByte('|')
		builder.WriteString(circle.ID)
		builder.WriteByte(':')
		builder.WriteString(circle.Shape)
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(len(circle.AttachedChildren)))
	}

	return builder.String()
}

func ShouldRefreshLocalAutonomous(lastSignature string, lastRefreshTick, currentTick int64, currentSignature string) bool {
	if currentTick <= lastRefreshTick {
		return true
	}
	if lastSignature == "" {
		return true
	}
	if currentSignature != lastSignature {
		return true
	}
	return currentTick-lastRefreshTick >= DefaultLocalAutonomousEveryTicks
}

func LocalFoodSignature(snapshot Snapshot) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(snapshot.TotalFoods))

	for _, food := range snapshot.Foods {
		builder.WriteByte('|')
		builder.WriteString(food.ID)
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(food.X, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(food.Y, 'f', -1, 64))
		builder.WriteByte(':')
		builder.WriteString(strconv.FormatFloat(food.Radius, 'f', -1, 64))
	}

	return builder.String()
}

func ShouldRefreshLocalFoods(lastSignature string, lastRefreshTick, currentTick int64, currentSignature string) bool {
	if currentTick <= lastRefreshTick {
		return true
	}
	if lastSignature == "" {
		return true
	}
	if currentSignature != lastSignature {
		return true
	}
	return currentTick-lastRefreshTick >= DefaultLocalFoodFallbackTicks
}

func buildViewportSnapshot(snapshot simulation.Snapshot, includeOrientation bool, compactOrientation bool, reducePrecision bool, orientationClusterSize float64) Snapshot {
	focusX := snapshot.World.Width / 2
	focusY := snapshot.World.Height / 2
	if snapshot.Player != nil {
		focusX = snapshot.Player.X
		focusY = snapshot.Player.Y
	} else if len(snapshot.AutonomousCircles) > 0 {
		focusX = snapshot.AutonomousCircles[0].X
		focusY = snapshot.AutonomousCircles[0].Y
	}

	halfWidth := minFloat(DefaultViewportInterestWidth, snapshot.World.Width)/2 + DefaultViewportInterestMargin
	halfHeight := minFloat(DefaultViewportInterestHeight, snapshot.World.Height)/2 + DefaultViewportInterestMargin

	left := focusX - halfWidth
	right := focusX + halfWidth
	top := focusY - halfHeight
	bottom := focusY + halfHeight

	requiredIDs := map[string]struct{}{}
	if snapshot.Interaction != nil {
		if snapshot.Interaction.SourceID != "" {
			requiredIDs[snapshot.Interaction.SourceID] = struct{}{}
		}
		if snapshot.Interaction.TargetID != "" {
			requiredIDs[snapshot.Interaction.TargetID] = struct{}{}
		}
	}

	localAutonomous := make([]simulation.AutonomousCircle, 0, len(snapshot.AutonomousCircles))
	var minimapAutonomous []MinimapAutonomousCircle
	var autonomousClusters map[string]*MinimapAutonomousCircle
	if includeOrientation && compactOrientation {
		autonomousClusters = make(map[string]*MinimapAutonomousCircle)
	} else if includeOrientation {
		minimapAutonomous = make([]MinimapAutonomousCircle, 0, len(snapshot.AutonomousCircles))
	}
	for _, circle := range snapshot.AutonomousCircles {
		if includeOrientation {
			if compactOrientation {
				key, x, y := minimapClusterKey(circle.X, circle.Y, snapshot.World, orientationClusterSize)
				clusterKey := circle.Shape + ":" + key
				cluster := autonomousClusters[clusterKey]
				if cluster == nil {
					autonomousClusters[clusterKey] = &MinimapAutonomousCircle{
						Shape: circle.Shape,
						X:     x,
						Y:     y,
						Count: 1,
					}
				} else {
					cluster.Count++
				}
			} else {
				minimapAutonomous = append(minimapAutonomous, MinimapAutonomousCircle{
					Shape: circle.Shape,
					X:     circle.X,
					Y:     circle.Y,
					Count: 1,
				})
			}
		}
		_, required := requiredIDs[circle.ID]
		if required || pointInsideRect(circle.X, circle.Y, left, right, top, bottom) {
			localAutonomous = append(localAutonomous, circle)
		}
	}

	localFoods := make([]simulation.Food, 0, len(snapshot.Foods))
	var minimapFoods []MinimapFood
	var foodClusters map[string]*MinimapFood
	if includeOrientation && compactOrientation {
		foodClusters = make(map[string]*MinimapFood)
	} else if includeOrientation {
		minimapFoods = make([]MinimapFood, 0, len(snapshot.Foods))
	}
	for _, food := range snapshot.Foods {
		if includeOrientation {
			if compactOrientation {
				key, x, y := minimapClusterKey(food.X, food.Y, snapshot.World, orientationClusterSize)
				cluster := foodClusters[key]
				if cluster == nil {
					foodClusters[key] = &MinimapFood{
						X:     x,
						Y:     y,
						Count: 1,
					}
				} else {
					cluster.Count++
				}
			} else {
				minimapFoods = append(minimapFoods, MinimapFood{
					X:     food.X,
					Y:     food.Y,
					Count: 1,
				})
			}
		}
		if pointInsideRect(food.X, food.Y, left, right, top, bottom) {
			localFoods = append(localFoods, food)
		}
	}
	if includeOrientation && compactOrientation {
		minimapAutonomous = clustersFromMap(autonomousClusters)
		minimapFoods = foodClustersFromMap(foodClusters)
	}

	return Snapshot{
		TransportMode:            "active_local_detail",
		Type:                     snapshot.Type,
		Tick:                     snapshot.Tick,
		World:                    snapshot.World,
		Player:                   roundedPlayer(snapshot.Player, reducePrecision),
		AutonomousCircles:        roundedAutonomousCircles(localAutonomous, reducePrecision),
		AutonomousFresh:          true,
		Interaction:              roundedInteraction(snapshot.Interaction, reducePrecision),
		Foods:                    roundedFoods(localFoods, reducePrecision),
		FoodsFresh:               true,
		OrientationFresh:         includeOrientation,
		MinimapAutonomousCircles: minimapAutonomous,
		MinimapFoods:             minimapFoods,
		TotalAutonomousCircles:   len(snapshot.AutonomousCircles),
		TotalFoods:               len(snapshot.Foods),
	}
}

func pointInsideRect(x, y, left, right, top, bottom float64) bool {
	return x >= left && x <= right && y >= top && y <= bottom
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minimapClusterKey(x, y float64, bounds simulation.Bounds, cellSize float64) (string, float64, float64) {
	column := int(x / cellSize)
	row := int(y / cellSize)
	centerX := minFloat(bounds.Width, float64(column)*cellSize+cellSize/2)
	centerY := minFloat(bounds.Height, float64(row)*cellSize+cellSize/2)
	return simulationKey(column, row), centerX, centerY
}

func simulationKey(column, row int) string {
	return strconv.Itoa(column) + ":" + strconv.Itoa(row)
}

func clustersFromMap(clusters map[string]*MinimapAutonomousCircle) []MinimapAutonomousCircle {
	keys := make([]string, 0, len(clusters))
	for key := range clusters {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	items := make([]MinimapAutonomousCircle, 0, len(clusters))
	for _, key := range keys {
		items = append(items, *clusters[key])
	}
	return items
}

func foodClustersFromMap(clusters map[string]*MinimapFood) []MinimapFood {
	keys := make([]string, 0, len(clusters))
	for key := range clusters {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	items := make([]MinimapFood, 0, len(clusters))
	for _, key := range keys {
		items = append(items, *clusters[key])
	}
	return items
}

func roundTransportFloat(value float64) float64 {
	return math.Round(value)
}

func roundedPlayer(circle *simulation.PlayerCircle, reducePrecision bool) *simulation.PlayerCircle {
	if circle == nil || !reducePrecision {
		return circle
	}

	copy := *circle
	copy.X = roundTransportFloat(copy.X)
	copy.Y = roundTransportFloat(copy.Y)
	copy.Radius = roundTransportFloat(copy.Radius)
	copy.Energy = roundTransportFloat(copy.Energy)
	copy.AttachedChildren = roundedChildren(copy.AttachedChildren, reducePrecision)
	return &copy
}

func roundedAutonomousCircles(circles []simulation.AutonomousCircle, reducePrecision bool) []simulation.AutonomousCircle {
	if !reducePrecision {
		return circles
	}

	rounded := make([]simulation.AutonomousCircle, 0, len(circles))
	for _, circle := range circles {
		copy := circle
		copy.X = roundTransportFloat(copy.X)
		copy.Y = roundTransportFloat(copy.Y)
		copy.Radius = roundTransportFloat(copy.Radius)
		copy.Energy = roundTransportFloat(copy.Energy)
		copy.AttachedChildren = roundedChildren(copy.AttachedChildren, reducePrecision)
		rounded = append(rounded, copy)
	}
	return rounded
}

func roundedChildren(children []simulation.AttachedChild, reducePrecision bool) []simulation.AttachedChild {
	if !reducePrecision {
		return children
	}

	rounded := make([]simulation.AttachedChild, 0, len(children))
	for _, child := range children {
		copy := child
		copy.X = roundTransportFloat(copy.X)
		copy.Y = roundTransportFloat(copy.Y)
		copy.Radius = roundTransportFloat(copy.Radius)
		rounded = append(rounded, copy)
	}
	return rounded
}

func roundedFoods(foods []simulation.Food, reducePrecision bool) []simulation.Food {
	if !reducePrecision {
		return foods
	}

	rounded := make([]simulation.Food, 0, len(foods))
	for _, food := range foods {
		copy := food
		copy.X = roundTransportFloat(copy.X)
		copy.Y = roundTransportFloat(copy.Y)
		copy.Radius = roundTransportFloat(copy.Radius)
		rounded = append(rounded, copy)
	}
	return rounded
}

func roundedInteraction(interaction *simulation.InteractionClassification, reducePrecision bool) *simulation.InteractionClassification {
	if interaction == nil || !reducePrecision {
		return interaction
	}

	copy := *interaction
	copy.SourceCapacityValue = roundTransportFloat(copy.SourceCapacityValue)
	copy.TargetCapacityValue = roundTransportFloat(copy.TargetCapacityValue)
	copy.SourceEnergyComponent = roundTransportFloat(copy.SourceEnergyComponent)
	copy.TargetEnergyComponent = roundTransportFloat(copy.TargetEnergyComponent)
	copy.SourceReserveComponent = roundTransportFloat(copy.SourceReserveComponent)
	copy.TargetReserveComponent = roundTransportFloat(copy.TargetReserveComponent)
	copy.ReproductionThreshold = roundTransportFloat(copy.ReproductionThreshold)
	copy.ReproductionCost = roundTransportFloat(copy.ReproductionCost)
	return &copy
}
