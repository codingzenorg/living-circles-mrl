package transport

import (
	"slices"
	"strconv"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

const (
	DefaultViewportInterestWidth  = 1200.0
	DefaultViewportInterestHeight = 840.0
	DefaultViewportInterestMargin = 180.0
	DefaultOrientationEveryTicks  = int64(5)
	DefaultMinimapClusterSize     = 480.0
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
	Type                     string                                `json:"type"`
	Tick                     int64                                 `json:"tick"`
	World                    simulation.Bounds                     `json:"world"`
	Player                   *simulation.PlayerCircle              `json:"player"`
	AutonomousCircles        []simulation.AutonomousCircle         `json:"autonomous_circles"`
	Interaction              *simulation.InteractionClassification `json:"interaction"`
	Foods                    []simulation.Food                     `json:"foods"`
	OrientationFresh         bool                                  `json:"orientation_fresh"`
	MinimapAutonomousCircles []MinimapAutonomousCircle             `json:"minimap_autonomous_circles"`
	MinimapFoods             []MinimapFood                         `json:"minimap_foods"`
	TotalAutonomousCircles   int                                   `json:"total_autonomous_circles"`
	TotalFoods               int                                   `json:"total_foods"`
}

func BuildViewportSnapshot(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	return buildViewportSnapshot(snapshot, includeOrientation, true)
}

func BuildViewportSnapshotExactOrientation(snapshot simulation.Snapshot, includeOrientation bool) Snapshot {
	return buildViewportSnapshot(snapshot, includeOrientation, false)
}

func buildViewportSnapshot(snapshot simulation.Snapshot, includeOrientation bool, compactOrientation bool) Snapshot {
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
				key, x, y := minimapClusterKey(circle.X, circle.Y, snapshot.World)
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
				key, x, y := minimapClusterKey(food.X, food.Y, snapshot.World)
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
		Type:                     snapshot.Type,
		Tick:                     snapshot.Tick,
		World:                    snapshot.World,
		Player:                   snapshot.Player,
		AutonomousCircles:        localAutonomous,
		Interaction:              snapshot.Interaction,
		Foods:                    localFoods,
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

func minimapClusterKey(x, y float64, bounds simulation.Bounds) (string, float64, float64) {
	cellSize := DefaultMinimapClusterSize
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
