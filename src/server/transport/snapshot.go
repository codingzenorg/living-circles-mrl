package transport

import "github.com/codingzen/living-circles-mrl/src/server/simulation"

const (
	DefaultViewportInterestWidth  = 1200.0
	DefaultViewportInterestHeight = 840.0
	DefaultViewportInterestMargin = 180.0
	DefaultOrientationEveryTicks  = int64(5)
)

type MinimapAutonomousCircle struct {
	ID    string  `json:"id"`
	Shape string  `json:"shape"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type MinimapFood struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
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
	if includeOrientation {
		minimapAutonomous = make([]MinimapAutonomousCircle, 0, len(snapshot.AutonomousCircles))
	}
	for _, circle := range snapshot.AutonomousCircles {
		if includeOrientation {
			minimapAutonomous = append(minimapAutonomous, MinimapAutonomousCircle{
				ID:    circle.ID,
				Shape: circle.Shape,
				X:     circle.X,
				Y:     circle.Y,
			})
		}
		_, required := requiredIDs[circle.ID]
		if required || pointInsideRect(circle.X, circle.Y, left, right, top, bottom) {
			localAutonomous = append(localAutonomous, circle)
		}
	}

	localFoods := make([]simulation.Food, 0, len(snapshot.Foods))
	var minimapFoods []MinimapFood
	if includeOrientation {
		minimapFoods = make([]MinimapFood, 0, len(snapshot.Foods))
	}
	for _, food := range snapshot.Foods {
		if includeOrientation {
			minimapFoods = append(minimapFoods, MinimapFood{
				ID: food.ID,
				X:  food.X,
				Y:  food.Y,
			})
		}
		if pointInsideRect(food.X, food.Y, left, right, top, bottom) {
			localFoods = append(localFoods, food)
		}
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
