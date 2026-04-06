package server_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

func promotedChildPosition(ownerID, childID string, x, y float64, tick int64) (float64, float64) {
	angle := promotedChildAngle(ownerID, childID, tick)
	orbitRadius := simulation.DefaultPlayerRadius + simulation.DefaultAttachedChildOrbitGap + simulation.DefaultAttachedChildRadius
	return x + math.Cos(angle)*orbitRadius, y + math.Sin(angle)*orbitRadius
}

func promotedChildAngle(ownerID, childID string, tick int64) float64 {
	baseAngle := float64(testHashString(ownerID+":"+childID)%360) * math.Pi / 180
	return baseAngle + float64(tick)*simulation.DefaultChildOrbitSpeed
}

func testHashString(value string) int {
	hash := 17
	for _, char := range value {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		return -hash
	}
	return hash
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func TestAdvanceMovesPlayerWhenIntentExists(t *testing.T) {
	session := simulation.NewSession()
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	after := session.Advance()

	if after.Player.X <= before.Player.X {
		t.Fatalf("expected player to move right, before=%v after=%v", before.Player.X, after.Player.X)
	}
}

func TestAdvanceConsumesEnergyWhenMovementOccurs(t *testing.T) {
	session := simulation.NewSession()
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 0, Y: -1})
	after := session.Advance()

	if after.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected energy to decrease, before=%v after=%v", before.Player.Energy, after.Player.Energy)
	}
}

func TestPlayerAvoidsCrowdingCostWhenNeighborhoodIsSparse(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:               simulation.DefaultWorldWidth,
		WorldHeight:              simulation.DefaultWorldHeight,
		UseExpandedPopulation:    false,
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		SecondaryAutonomousShape: "",
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		AutonomousEnergy:         simulation.DefaultPlayerEnergy,
		AutonomousX:              80,
		AutonomousY:              80,
	})
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	after := session.Advance()

	expectedEnergy := before.Player.Energy - simulation.DefaultMoveCost
	if after.Player.Energy != expectedEnergy {
		t.Fatalf("expected sparse player movement energy %v, got %v", expectedEnergy, after.Player.Energy)
	}
}

func TestPlayerPaysCrowdingCostWhenTwoNeighborsAreNearby(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:                1000,
		WorldHeight:               800,
		UseExpandedPopulation:     false,
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "triangle",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          simulation.DefaultPlayerEnergy,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
		PlayerX:                   120,
		PlayerY:                   120,
		AutonomousX:               180,
		AutonomousY:               120,
		SecondaryAutonomousX:      120,
		SecondaryAutonomousY:      180,
	})
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	after := session.Advance()

	expectedEnergy := before.Player.Energy - simulation.DefaultMoveCost - simulation.DefaultCrowdingMoveCost
	if after.Player.Energy != expectedEnergy {
		t.Fatalf("expected crowded player movement energy %v, got %v", expectedEnergy, after.Player.Energy)
	}
}

func TestAutonomousPaysCrowdingCostWhenTwoNeighborsAreNearby(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:                1000,
		WorldHeight:               800,
		UseExpandedPopulation:     false,
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "triangle",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
		PlayerX:                   120,
		PlayerY:                   120,
		AutonomousX:               180,
		AutonomousY:               120,
		SecondaryAutonomousX:      180,
		SecondaryAutonomousY:      180,
	})
	before := session.Snapshot()

	after := session.Advance()

	expectedEnergy := before.AutonomousCircles[0].Energy - simulation.DefaultMoveCost - simulation.DefaultCrowdingMoveCost
	if after.AutonomousCircles[0].Energy != expectedEnergy {
		t.Fatalf("expected crowded autonomous movement energy %v, got %v", expectedEnergy, after.AutonomousCircles[0].Energy)
	}
}

func TestPlayerWithZeroEnergyDisappearsWithoutChildren(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "square",
		PlayerEnergy:     1,
		AutonomousEnergy: 100,
	})

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	snapshot := session.Advance()

	if snapshot.Player != nil {
		t.Fatalf("expected player to disappear on zero energy, got %+v", snapshot.Player)
	}
}

func TestPlayerWithZeroEnergyReplacesThroughChildContinuity(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:         "triangle",
		AutonomousShape:     "square",
		PlayerEnergy:        0,
		AutonomousEnergy:    100,
		PlayerChildrenCount: 1,
	})
	before := session.Snapshot()
	expectedChildID := before.Player.AttachedChildren[0].ID
	expectedX, expectedY := promotedChildPosition(before.Player.ID, before.Player.AttachedChildren[0].ID, before.Player.X, before.Player.Y, 1)

	snapshot := session.Advance()

	if snapshot.Player == nil {
		t.Fatal("expected player replacement after zero-energy collapse")
	}
	if len(snapshot.Player.AttachedChildren) != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", len(snapshot.Player.AttachedChildren))
	}
	if snapshot.Player.Energy != simulation.DefaultReplacementEnergy {
		t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.Player.Energy)
	}
	if snapshot.Player.LineageID != before.Player.LineageID {
		t.Fatalf("expected replacement lineage %q, got %q", before.Player.LineageID, snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != before.Player.Generation+1 {
		t.Fatalf("expected replacement generation %d, got %d", before.Player.Generation+1, snapshot.Player.Generation)
	}
	if !almostEqual(snapshot.Player.X, expectedX) || !almostEqual(snapshot.Player.Y, expectedY) {
		t.Fatalf("expected replacement at promoted child position (%v, %v), got (%v, %v)", expectedX, expectedY, snapshot.Player.X, snapshot.Player.Y)
	}
	if snapshot.Interaction == nil {
		t.Fatal("expected continuity interaction after zero-energy promotion")
	}
	if snapshot.Interaction.Kind != "death_promoted_child" {
		t.Fatalf("expected death_promoted_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != before.Player.ID || snapshot.Interaction.TargetID != before.Player.ID {
		t.Fatalf("expected continuity interaction to identify player %q, got source=%q target=%q", before.Player.ID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if snapshot.Interaction.PromotedChildID != expectedChildID {
		t.Fatalf("expected promoted child id %q, got %q", expectedChildID, snapshot.Interaction.PromotedChildID)
	}
	if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) {
		t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren))
	}
}

func TestAutonomousCircleWithZeroEnergyDisappearsWithoutChildren(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "triangle",
		PlayerEnergy:     100,
		AutonomousEnergy: 0,
	})
	snapshot := session.Advance()

	if len(snapshot.AutonomousCircles) != 0 {
		t.Fatalf("expected zero-energy autonomous circle to disappear, got %d circles", len(snapshot.AutonomousCircles))
	}
}

func TestAutonomousCircleWithZeroEnergyReplacesThroughChildContinuity(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        0,
		AutonomousChildrenCount: 1,
		DisableFoodSeeking:      true,
	})
	before := session.Snapshot()
	expectedChildID := before.AutonomousCircles[0].AttachedChildren[0].ID
	expectedX, expectedY := promotedChildPosition(before.AutonomousCircles[0].ID, before.AutonomousCircles[0].AttachedChildren[0].ID, before.AutonomousCircles[0].X, before.AutonomousCircles[0].Y, 1)
	snapshot := session.Advance()

	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous replacement after zero-energy collapse, got %d circles", len(snapshot.AutonomousCircles))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if snapshot.AutonomousCircles[0].Energy != simulation.DefaultReplacementEnergy {
		t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.AutonomousCircles[0].LineageID != before.AutonomousCircles[0].LineageID {
		t.Fatalf("expected replacement lineage %q, got %q", before.AutonomousCircles[0].LineageID, snapshot.AutonomousCircles[0].LineageID)
	}
	if snapshot.AutonomousCircles[0].Generation != before.AutonomousCircles[0].Generation+1 {
		t.Fatalf("expected replacement generation %d, got %d", before.AutonomousCircles[0].Generation+1, snapshot.AutonomousCircles[0].Generation)
	}
	if !almostEqual(snapshot.AutonomousCircles[0].X, expectedX) || !almostEqual(snapshot.AutonomousCircles[0].Y, expectedY) {
		t.Fatalf("expected replacement at promoted child position (%v, %v), got (%v, %v)", expectedX, expectedY, snapshot.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].Y)
	}
	if snapshot.Interaction == nil {
		t.Fatal("expected continuity interaction after zero-energy promotion")
	}
	if snapshot.Interaction.Kind != "death_promoted_child" {
		t.Fatalf("expected death_promoted_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != before.AutonomousCircles[0].ID || snapshot.Interaction.TargetID != before.AutonomousCircles[0].ID {
		t.Fatalf("expected continuity interaction to identify autonomous circle %q, got source=%q target=%q", before.AutonomousCircles[0].ID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if snapshot.Interaction.PromotedChildID != expectedChildID {
		t.Fatalf("expected promoted child id %q, got %q", expectedChildID, snapshot.Interaction.PromotedChildID)
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
		t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestPositiveEnergyDoesNotTriggerCollapseDeath(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "square",
		PlayerEnergy:     2,
		AutonomousEnergy: 100,
	})

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	snapshot := session.Advance()

	if snapshot.Player == nil {
		t.Fatal("expected positive-energy player to remain active")
	}
	if snapshot.Player.Energy != 1 {
		t.Fatalf("expected energy to decrease without death, got %v", snapshot.Player.Energy)
	}
}

func TestAdvanceKeepsPlayerInsideBounds(t *testing.T) {
	session := simulation.NewSession()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	var snapshot simulation.Snapshot
	for range 60 {
		snapshot = session.Advance()
	}
	if snapshot.Player == nil {
		t.Fatal("expected player to remain active while checking bounds")
	}

	if snapshot.Player.X > snapshot.World.Width-simulation.DefaultPlayerRadius {
		t.Fatalf("expected x to stay inside bounds, got %v", snapshot.Player.X)
	}
}

func TestPlayerChildGrowthDoesNotChangeWallClamp(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              simulation.DefaultPlayerShape,
		AutonomousShape:          simulation.DefaultAutoShape,
		SecondaryAutonomousShape: "",
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:      3,
		AutonomousEnergy:         0,
	})

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	var snapshot simulation.Snapshot
	for range 60 {
		snapshot = session.Advance()
	}

	if snapshot.Player == nil {
		t.Fatal("expected player to remain active while checking wall clamp")
	}
	if snapshot.Player.X != snapshot.World.Width-simulation.DefaultPlayerRadius {
		t.Fatalf("expected player wall clamp at visible body boundary %v, got %v", snapshot.World.Width-simulation.DefaultPlayerRadius, snapshot.Player.X)
	}
}

func TestAutonomousChildGrowthDoesNotChangeWallClamp(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         simulation.DefaultPlayerShape,
		AutonomousShape:                     simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:            "",
		PlayerEnergy:                        0,
		AutonomousEnergy:                    simulation.DefaultPlayerEnergy,
		AutonomousX:                         simulation.DefaultWorldWidth - 20,
		AutonomousY:                         simulation.DefaultWorldHeight / 2,
		AutonomousChildrenCount:             3,
		DisableFoodSeeking:                  true,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	snapshot := session.Advance()

	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected one autonomous circle, got %d", len(snapshot.AutonomousCircles))
	}
	if snapshot.AutonomousCircles[0].X != snapshot.World.Width-simulation.DefaultPlayerRadius {
		t.Fatalf("expected autonomous wall clamp at visible body boundary %v, got %v", snapshot.World.Width-simulation.DefaultPlayerRadius, snapshot.AutonomousCircles[0].X)
	}
}

func TestAdvanceDoesNotMoveOnIdleTick(t *testing.T) {
	session := simulation.NewSession()
	before := session.Snapshot()

	after := session.Advance()

	if after.Player.X != before.Player.X || after.Player.Y != before.Player.Y {
		t.Fatalf("expected idle tick to keep player still, before=%+v after=%+v", before.Player, after.Player)
	}
}

func TestNewWorldContainsDeterministicFoodItems(t *testing.T) {
	world := simulation.NewWorld()
	snapshot := world.Snapshot(0)
	recreated := simulation.NewWorld().Snapshot(0)

	if snapshot.World.Width != simulation.DefaultExpandedWorldWidth || snapshot.World.Height != simulation.DefaultExpandedWorldHeight {
		t.Fatalf("expected expanded default world %vx%v, got %vx%v", simulation.DefaultExpandedWorldWidth, simulation.DefaultExpandedWorldHeight, snapshot.World.Width, snapshot.World.Height)
	}
	if len(snapshot.AutonomousCircles) != simulation.DefaultExpandedAutonomousCount {
		t.Fatalf("expected %d autonomous circles, got %d", simulation.DefaultExpandedAutonomousCount, len(snapshot.AutonomousCircles))
	}

	if snapshot.Player.Shape != simulation.DefaultPlayerShape {
		t.Fatalf("expected player shape %q, got %q", simulation.DefaultPlayerShape, snapshot.Player.Shape)
	}
	if len(snapshot.Player.AttachedChildren) != 1 {
		t.Fatalf("expected player children count to start at one for demo continuity, got %d", len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren))
	}
	if snapshot.Player.LineageID != "lineage-player-1" {
		t.Fatalf("expected player lineage %q, got %q", "lineage-player-1", snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != 0 {
		t.Fatalf("expected player generation 0, got %d", snapshot.Player.Generation)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed player radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.Player.Radius)
	}

	if snapshot.AutonomousCircles[0].Shape != snapshot.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", snapshot.Player.Shape, snapshot.AutonomousCircles[0].Shape)
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 1 {
		t.Fatalf("expected first autonomous children count to start at one for demo continuity, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
		t.Fatalf("expected first autonomous attached children to match count, count=%d attached=%d", len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if snapshot.AutonomousCircles[0].LineageID != "lineage-circle-2" {
		t.Fatalf("expected first autonomous lineage %q, got %q", "lineage-circle-2", snapshot.AutonomousCircles[0].LineageID)
	}
	if snapshot.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected first autonomous generation 0, got %d", snapshot.AutonomousCircles[0].Generation)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed first autonomous radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.AutonomousCircles[0].Radius)
	}

	if snapshot.AutonomousCircles[1].Shape != simulation.DefaultAutoShape {
		t.Fatalf("expected second autonomous shape %q, got %q", simulation.DefaultAutoShape, snapshot.AutonomousCircles[1].Shape)
	}
	if len(snapshot.AutonomousCircles[1].AttachedChildren) != 0 {
		t.Fatalf("expected second autonomous children count to start at zero, got %d", len(snapshot.AutonomousCircles[1].AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[1].AttachedChildren) != len(snapshot.AutonomousCircles[1].AttachedChildren) {
		t.Fatalf("expected second autonomous attached children to match count, count=%d attached=%d", len(snapshot.AutonomousCircles[1].AttachedChildren), len(snapshot.AutonomousCircles[1].AttachedChildren))
	}
	if snapshot.AutonomousCircles[1].LineageID != "lineage-circle-3" {
		t.Fatalf("expected second autonomous lineage %q, got %q", "lineage-circle-3", snapshot.AutonomousCircles[1].LineageID)
	}
	if snapshot.AutonomousCircles[1].Generation != 0 {
		t.Fatalf("expected second autonomous generation 0, got %d", snapshot.AutonomousCircles[1].Generation)
	}
	if snapshot.AutonomousCircles[1].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected base autonomous radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.AutonomousCircles[1].Radius)
	}
	if snapshot.AutonomousCircles[2].ID != simulation.DefaultTertiaryID ||
		snapshot.AutonomousCircles[3].ID != simulation.DefaultQuaternaryID ||
		snapshot.AutonomousCircles[4].ID != simulation.DefaultQuinaryID ||
		snapshot.AutonomousCircles[5].ID != simulation.DefaultSenaryID ||
		snapshot.AutonomousCircles[6].ID != simulation.DefaultSeptenaryID ||
		snapshot.AutonomousCircles[7].ID != simulation.DefaultOctonaryID {
		t.Fatalf("expected deterministic expanded autonomous ids, got %+v", snapshot.AutonomousCircles)
	}
	for index, circle := range snapshot.AutonomousCircles {
		if recreated.AutonomousCircles[index].X != circle.X || recreated.AutonomousCircles[index].Y != circle.Y {
			t.Fatalf(
				"expected recreated world to keep deterministic autonomous placement at %d, before=(%v,%v) after=(%v,%v)",
				index,
				circle.X,
				circle.Y,
				recreated.AutonomousCircles[index].X,
				recreated.AutonomousCircles[index].Y,
			)
		}
	}
	expectedFoodCount := len(snapshot.AutonomousCircles) + 3
	if len(snapshot.Foods) != expectedFoodCount {
		t.Fatalf("expected expanded food baseline to derive from active population, expected %d got %d", expectedFoodCount, len(snapshot.Foods))
	}
	if len(snapshot.Foods) != simulation.DefaultExpandedFoodCount {
		t.Fatalf("expected expanded default food count %d, got %d", simulation.DefaultExpandedFoodCount, len(snapshot.Foods))
	}
	if len(recreated.Foods) != len(snapshot.Foods) {
		t.Fatalf("expected recreated world to keep food count %d, got %d", len(snapshot.Foods), len(recreated.Foods))
	}
	for index, food := range snapshot.Foods {
		if food.ID != "food-"+strconv.Itoa(index+1) {
			t.Fatalf("expected deterministic expanded food id food-%d, got %q", index+1, food.ID)
		}
		if food.X == simulation.DefaultExpandedWorldWidth/2 && food.Y == simulation.DefaultExpandedWorldHeight/2 {
			t.Fatalf("expected seeded food layout to avoid exact world center, got %+v", food)
		}
		if recreated.Foods[index] != food {
			t.Fatalf("expected recreated world to keep deterministic food slot at %d, before=%+v after=%+v", index, food, recreated.Foods[index])
		}
	}
}

func TestDefaultSessionResetRestoresExpandedPopulationDeterministically(t *testing.T) {
	session := simulation.NewSession()
	before := session.Snapshot()

	session.Advance()
	reset := session.Reset()

	if reset.World != before.World {
		t.Fatalf("expected reset world bounds %+v, got %+v", before.World, reset.World)
	}
	if len(reset.AutonomousCircles) != len(before.AutonomousCircles) {
		t.Fatalf("expected reset autonomous count %d, got %d", len(before.AutonomousCircles), len(reset.AutonomousCircles))
	}
	if len(reset.Foods) != len(before.Foods) {
		t.Fatalf("expected reset food count %d, got %d", len(before.Foods), len(reset.Foods))
	}
	if reset.Player == nil || before.Player == nil {
		t.Fatal("expected player before and after reset")
	}
	if reset.Player.X != before.Player.X || reset.Player.Y != before.Player.Y {
		t.Fatalf("expected reset player position (%v,%v), got (%v,%v)", before.Player.X, before.Player.Y, reset.Player.X, reset.Player.Y)
	}
	for index, circle := range before.AutonomousCircles {
		if reset.AutonomousCircles[index].X != circle.X || reset.AutonomousCircles[index].Y != circle.Y {
			t.Fatalf(
				"expected reset autonomous %d to match initial position, before=(%v,%v) after=(%v,%v)",
				index,
				circle.X,
				circle.Y,
				reset.AutonomousCircles[index].X,
				reset.AutonomousCircles[index].Y,
			)
		}
	}
	for index, food := range before.Foods {
		if reset.Foods[index] != food {
			t.Fatalf("expected reset food slot %d to match initial slot, before=%+v after=%+v", index, food, reset.Foods[index])
		}
	}
}

func TestCustomNarrowWorldKeepsSmallerDeterministicFoodBaseline(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              simulation.DefaultPlayerShape,
		AutonomousShape:          simulation.DefaultAutoShape,
		SecondaryAutonomousShape: "",
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		AutonomousEnergy:         simulation.DefaultPlayerEnergy,
		UseExpandedPopulation:    false,
		WorldWidth:               simulation.DefaultWorldWidth,
		WorldHeight:              simulation.DefaultWorldHeight,
	})
	snapshot := session.Snapshot()

	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected one autonomous circle in narrow custom world, got %d", len(snapshot.AutonomousCircles))
	}
	if len(snapshot.Foods) != 3 {
		t.Fatalf("expected narrow custom world to keep 3 deterministic food slots, got %d", len(snapshot.Foods))
	}
}

func TestExpandedDefaultWorldKeepsInitialEntitiesInsideBounds(t *testing.T) {
	snapshot := simulation.NewWorld().Snapshot(0)

	assertInsideBounds := func(x, y, radius float64, label string) {
		t.Helper()
		if x < radius || x > snapshot.World.Width-radius || y < radius || y > snapshot.World.Height-radius {
			t.Fatalf("expected %s inside bounds, got x=%v y=%v radius=%v world=%+v", label, x, y, radius, snapshot.World)
		}
	}

	assertInsideBounds(snapshot.Player.X, snapshot.Player.Y, snapshot.Player.Radius, "player")
	for _, circle := range snapshot.AutonomousCircles {
		assertInsideBounds(circle.X, circle.Y, circle.Radius, circle.ID)
	}
	for _, food := range snapshot.Foods {
		assertInsideBounds(food.X, food.Y, food.Radius, food.ID)
	}
}

func TestOverlappingFoodRemovesItAndRestoresEnergy(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	})
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	firstTick := session.Advance()
	secondTick := session.Advance()

	if len(secondTick.Foods) >= len(before.Foods) {
		t.Fatalf("expected one food item to be consumed, before=%d after=%d", len(before.Foods), len(secondTick.Foods))
	}

	if secondTick.Player.Energy <= firstTick.Player.Energy {
		t.Fatalf("expected food consumption to restore energy, before=%v after=%v", firstTick.Player.Energy, secondTick.Player.Energy)
	}
}

func TestConsumedFoodRegeneratesAfterDeterministicDelay(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	})
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	_ = session.Advance()
	consumed := session.Advance()

	if len(consumed.Foods) != len(before.Foods)-1 {
		t.Fatalf("expected one consumed food item, before=%d after=%d", len(before.Foods), len(consumed.Foods))
	}

	var snapshot simulation.Snapshot
	for range simulation.DefaultFoodRegenDelay - 1 {
		session.ApplyIntent(simulation.Vector{})
		snapshot = session.Advance()
	}

	if len(snapshot.Foods) != len(before.Foods)-1 {
		t.Fatalf("expected food to remain missing before regen delay, got %d foods", len(snapshot.Foods))
	}

	session.ApplyIntent(simulation.Vector{})
	regenerated := session.Advance()

	if len(regenerated.Foods) != len(before.Foods) {
		t.Fatalf("expected consumed food to regenerate, got %d foods", len(regenerated.Foods))
	}

	found := false
	for _, food := range regenerated.Foods {
		if food.ID != "food-1" {
			continue
		}
		found = true
		if food.X != 432 || food.Y != 300 {
			t.Fatalf("expected regenerated food to return to original slot, got %+v", food)
		}
	}
	if !found {
		t.Fatal("expected regenerated food slot food-1 to be present")
	}
}

func TestConsumedFoodDoesNotRegenerateBeforeDelayBoundary(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	})

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	_ = session.Advance()
	snapshot := session.Advance()

	if foodIDs(snapshot.Foods)["food-1"] {
		t.Fatal("expected food-1 to be absent immediately after consumption")
	}

	for range simulation.DefaultFoodRegenDelay - 2 {
		session.ApplyIntent(simulation.Vector{})
		snapshot = session.Advance()
	}

	if foodIDs(snapshot.Foods)["food-1"] {
		t.Fatal("expected food-1 to remain absent before the full delay elapses")
	}
}

func TestMultipleMissingFoodSlotsIncreaseRegenDelayDeterministically(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   424,
		PlayerY:                   300,
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousX:               292,
		AutonomousY:               300,
		AutonomousEnergy:          simulation.DefaultPlayerEnergy,
		SecondaryAutonomousEnergy: 0,
	})
	before := session.Snapshot()
	snapshot := session.Advance()

	consumedExtra := false
	if len(snapshot.Foods) <= len(before.Foods)-2 {
		consumedExtra = true
	}

	if !consumedExtra {
		t.Fatal("expected expanded world to consume at least two food slots")
	}

	for range simulation.DefaultFoodRegenDelay {
		snapshot = session.Advance()
	}

	if len(snapshot.Foods) >= len(before.Foods) {
		t.Fatalf("expected at least one slot to remain missing under pressure-scaled regeneration, before=%d after=%d", len(before.Foods), len(snapshot.Foods))
	}
}

func TestFoodRegenerationDoesNotDuplicateActiveSlots(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	})

	var snapshot simulation.Snapshot
	for range simulation.DefaultFoodRegenDelay + 3 {
		snapshot = session.Advance()
	}

	if len(snapshot.Foods) != 3 {
		t.Fatalf("expected active food slots to remain non-duplicated, got %d foods", len(snapshot.Foods))
	}

	ids := foodIDs(snapshot.Foods)
	if len(ids) != len(snapshot.Foods) {
		t.Fatalf("expected unique food ids after regeneration cycles, got %+v", snapshot.Foods)
	}
}

func TestInitialChildrenDoNotChangeVisibleParentRadius(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "square",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     2,
		AutonomousChildrenCount: 1,
	})
	snapshot := session.Snapshot()

	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed player radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed autonomous radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.AutonomousCircles[0].Radius)
	}
	if len(snapshot.Player.AttachedChildren) != 2 {
		t.Fatalf("expected player children to remain visible, got %d", len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 1 {
		t.Fatalf("expected autonomous children to remain visible, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestFoodRecoveryRespectsEnergyCap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	})
	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	_ = session.Advance()
	snapshot := session.Advance()

	if snapshot.Player.Energy > simulation.DefaultMaxEnergy {
		t.Fatalf("expected energy cap to hold, got %v", snapshot.Player.Energy)
	}

	if snapshot.Player.Energy != simulation.DefaultMaxEnergy {
		t.Fatalf("expected food recovery to clamp to max energy, got %v", snapshot.Player.Energy)
	}
}

func TestDerivedParentRadiusAloneDoesNotCollectFoodWithoutVisibleOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		SecondaryAutonomousShape: "",
		PlayerX:                  380,
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:      2,
		AutonomousEnergy:         0,
	})

	before := session.Snapshot()
	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	snapshot := session.Advance()

	if len(snapshot.Foods) != len(before.Foods) {
		t.Fatalf("expected no food consumption without visible overlap, before=%d after=%d", len(before.Foods), len(snapshot.Foods))
	}
	expectedEnergy := simulation.DefaultPlayerEnergy - simulation.DefaultMoveCost
	if snapshot.Player.Energy != expectedEnergy {
		t.Fatalf("expected no immediate food recovery from derived radius alone, got %v", snapshot.Player.Energy)
	}
}

func TestAttachedChildCanCollectFoodForPlayer(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		PlayerX:                  403,
		PlayerEnergy:             80,
		PlayerChildrenCount:      2,
		AutonomousEnergy:         0,
		SecondaryAutonomousShape: "",
	})

	before := session.Snapshot()
	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	snapshot := session.Advance()

	if len(snapshot.Foods) != len(before.Foods)-1 {
		t.Fatalf("expected attached child collection to consume one food, before=%d after=%d", len(before.Foods), len(snapshot.Foods))
	}
	expectedEnergy := before.Player.Energy - simulation.DefaultMoveCost + simulation.DefaultFoodEnergy
	if snapshot.Player.Energy != expectedEnergy {
		t.Fatalf("expected attached child collection to restore player energy to %v, got %v", expectedEnergy, snapshot.Player.Energy)
	}
	if len(snapshot.Player.AttachedChildren) != 2 {
		t.Fatalf("expected attached children to remain attached after collection, got %d", len(snapshot.Player.AttachedChildren))
	}
}

func TestAutonomousCircleMovesWithoutPlayerInput(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         simulation.DefaultPlayerShape,
		AutonomousShape:                     simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:            simulation.DefaultAutoShape,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:                 1,
		AutonomousEnergy:                    80,
		SecondaryAutonomousEnergy:           simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})
	before := session.Snapshot()
	after := session.Advance()

	if after.AutonomousCircles[0].X == before.AutonomousCircles[0].X && after.AutonomousCircles[0].Y == before.AutonomousCircles[0].Y {
		t.Fatalf("expected first autonomous circle to move, before=%+v after=%+v", before.AutonomousCircles[0], after.AutonomousCircles[0])
	}
	if after.AutonomousCircles[1].X == before.AutonomousCircles[1].X && after.AutonomousCircles[1].Y == before.AutonomousCircles[1].Y {
		t.Fatalf("expected second autonomous circle to move, before=%+v after=%+v", before.AutonomousCircles[1], after.AutonomousCircles[1])
	}
}

func TestAutonomousCircleSteersTowardNearestFood(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         simulation.DefaultPlayerShape,
		AutonomousShape:                     simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:            simulation.DefaultAutoShape,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:                 1,
		AutonomousEnergy:                    80,
		SecondaryAutonomousEnergy:           simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})
	before := session.Snapshot()
	after := session.Advance()

	if after.AutonomousCircles[0].X <= before.AutonomousCircles[0].X {
		t.Fatalf("expected first autonomous circle to move right toward nearest food, before=%v after=%v", before.AutonomousCircles[0].X, after.AutonomousCircles[0].X)
	}
	if after.AutonomousCircles[1].Y <= before.AutonomousCircles[1].Y {
		t.Fatalf("expected second autonomous circle to move downward toward nearest food, before=%v after=%v", before.AutonomousCircles[1].Y, after.AutonomousCircles[1].Y)
	}
}

func TestAutonomousCircleConsumesEnergyWhenMoving(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         simulation.DefaultPlayerShape,
		AutonomousShape:                     simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:            simulation.DefaultAutoShape,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:                 1,
		AutonomousEnergy:                    80,
		SecondaryAutonomousEnergy:           simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})
	before := session.Snapshot()
	after := session.Advance()

	if after.AutonomousCircles[0].Energy >= before.AutonomousCircles[0].Energy {
		t.Fatalf("expected first autonomous circle energy to decrease, before=%v after=%v", before.AutonomousCircles[0].Energy, after.AutonomousCircles[0].Energy)
	}
	if after.AutonomousCircles[1].Energy >= before.AutonomousCircles[1].Energy {
		t.Fatalf("expected second autonomous circle energy to decrease, before=%v after=%v", before.AutonomousCircles[1].Energy, after.AutonomousCircles[1].Energy)
	}
}

func TestAutonomousCircleCanConsumeFoodAndRecoverEnergy(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         simulation.DefaultPlayerShape,
		AutonomousShape:                     simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:            simulation.DefaultAutoShape,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:                 1,
		AutonomousEnergy:                    80,
		SecondaryAutonomousEnergy:           simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})
	firstTick := session.Advance()
	secondTick := session.Advance()

	if len(secondTick.Foods) >= len(firstTick.Foods) {
		t.Fatalf("expected autonomous food consumption, before=%d after=%d", len(firstTick.Foods), len(secondTick.Foods))
	}

	if secondTick.AutonomousCircles[0].Energy <= firstTick.AutonomousCircles[0].Energy {
		t.Fatalf("expected autonomous circle energy recovery, before=%v after=%v", firstTick.AutonomousCircles[0].Energy, secondTick.AutonomousCircles[0].Energy)
	}
}

func TestAutonomousAttachedChildCanCollectFoodForParent(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "triangle",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          80,
		AutonomousChildrenCount:   2,
		SecondaryAutonomousShape:  "",
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})

	before := session.Snapshot()
	snapshot := session.Advance()

	if len(snapshot.Foods) != len(before.Foods)-1 {
		t.Fatalf("expected autonomous attached child collection to consume one food, before=%d after=%d", len(before.Foods), len(snapshot.Foods))
	}
	expectedEnergy := before.AutonomousCircles[0].Energy - simulation.DefaultMoveCost + simulation.DefaultFoodEnergy
	if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
		t.Fatalf("expected autonomous attached child collection to restore energy to %v, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 2 {
		t.Fatalf("expected attached children to remain attached after autonomous collection, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestAutonomousFoodSeekingCanCollectOffLaneFood(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousShape:           "triangle",
		AutonomousX:               100,
		AutonomousY:               100,
		AutonomousEnergy:          simulation.DefaultPlayerEnergy,
		SecondaryAutonomousShape:  "square",
		SecondaryAutonomousX:      520,
		SecondaryAutonomousY:      300,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
	})
	before := session.Snapshot()

	var snapshot simulation.Snapshot
	collected := false
	for range 20 {
		snapshot = session.Advance()
		if len(snapshot.Foods) < len(before.Foods) {
			collected = true
			break
		}
	}

	if !collected {
		t.Fatal("expected autonomous food seeking to collect food within the first steering window")
	}
	if snapshot.AutonomousCircles[1].Y <= before.AutonomousCircles[1].Y {
		t.Fatalf("expected second autonomous circle to steer off its original horizontal lane, before=%v after=%v", before.AutonomousCircles[1].Y, snapshot.AutonomousCircles[1].Y)
	}
}

func TestAttachedChildCanMakeFoodEffectivelyNearestForTargeting(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "square",
		PlayerX:                             700,
		PlayerY:                             500,
		PlayerEnergy:                        100,
		PlayerChildrenCount:                 0,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "",
		AutonomousX:                         369,
		AutonomousY:                         300,
		AutonomousEnergy:                    100,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	parentToFoodOne := math.Hypot(before.Foods[0].X-before.AutonomousCircles[0].X, before.Foods[0].Y-before.AutonomousCircles[0].Y)
	parentToFoodTwo := math.Hypot(before.Foods[1].X-before.AutonomousCircles[0].X, before.Foods[1].Y-before.AutonomousCircles[0].Y)
	if parentToFoodOne >= parentToFoodTwo {
		t.Fatalf("expected parent body to be nearer to food-1 than food-2, food-1=%v food-2=%v", parentToFoodOne, parentToFoodTwo)
	}

	childToFoodTwo := parentToFoodTwo
	for _, child := range before.AutonomousCircles[0].AttachedChildren {
		distance := math.Hypot(before.Foods[1].X-child.X, before.Foods[1].Y-child.Y)
		if distance < childToFoodTwo {
			childToFoodTwo = distance
		}
	}
	if childToFoodTwo >= parentToFoodOne {
		t.Fatalf("expected attached child to make food-2 effectively nearer than food-1, food-1=%v food-2=%v", parentToFoodOne, childToFoodTwo)
	}

	snapshot := session.Advance()
	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during child-aware food targeting, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to move left toward child-aware food target, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestLowEnergyChildAwareFoodTargetingCanOverrideParentOnlyDistance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "square",
		PlayerX:                             700,
		PlayerY:                             500,
		PlayerEnergy:                        100,
		PlayerChildrenCount:                 0,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "",
		AutonomousX:                         369,
		AutonomousY:                         300,
		AutonomousEnergy:                    39,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction while low-energy autonomous seeks child-aware food target, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= before.AutonomousCircles[0].X {
		t.Fatalf("expected low-energy autonomous circle to move left toward child-aware food target, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAutonomousKeepsOrdinarySteeringWhenCrowdingDifferenceIsSmall(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:                          1000,
		WorldHeight:                         800,
		UseExpandedPopulation:               false,
		PlayerShape:                         "square",
		PlayerX:                             700,
		PlayerY:                             500,
		PlayerEnergy:                        100,
		PlayerChildrenCount:                 0,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "",
		AutonomousX:                         369,
		AutonomousY:                         300,
		AutonomousEnergy:                    100,
		AutonomousChildrenCount:             1,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	snapshot := session.Advance()

	if snapshot.AutonomousCircles[0].X <= before.AutonomousCircles[0].X {
		t.Fatalf("expected ordinary food steering to remain in control when crowding difference is small, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAutonomousAvoidsMovingIntoClearlyMoreCrowdedDirection(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		WorldWidth:                          1000,
		WorldHeight:                         800,
		UseExpandedPopulation:               false,
		PlayerShape:                         "square",
		PlayerX:                             486,
		PlayerY:                             260,
		PlayerEnergy:                        100,
		PlayerChildrenCount:                 0,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "square",
		AutonomousX:                         360,
		AutonomousY:                         260,
		SecondaryAutonomousX:                487,
		SecondaryAutonomousY:                260,
		AutonomousEnergy:                    100,
		SecondaryAutonomousEnergy:           100,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	snapshot := session.Advance()

	if snapshot.AutonomousCircles[0].X >= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to reverse away from clearly more crowded direction, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAutonomousInteractionSeekingCanCreatePreferredDifferentShapeInteractionWithoutPlayerMovement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 80,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected autonomous interaction-seeking interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected shape-aware different-shape preference to produce reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestAutonomousInteractionSeekingCanCreateReproductionWithoutPlayerMovement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected autonomous interaction-seeking reproduction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestLowEnergyAutonomousCirclePrefersFoodOverInteractionTarget(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          39,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction while low-energy circle seeks food, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].Y >= 500 {
		t.Fatalf("expected low-energy autonomous circle to steer toward food lane, got y=%v", snapshot.AutonomousCircles[0].Y)
	}
}

func TestStableEnergyAutonomousCircleStillSeeksInteractionTarget(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          40,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected stable-energy autonomous circle to still seek interaction")
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected deterministic autonomous target %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestInteractionSeekingPrefersDifferentShapeTargetWhenAvailable(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   220,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      260,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected shape-aware interaction seeking to produce contact")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected preferred different-shape target to produce reproduction, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != "player-1" || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected different-shape target %q to be chosen ahead of the player, got %q -> %q", simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestInteractionSeekingFallsBackWhenNoDifferentShapeTargetExists(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   220,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "triangle",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      420,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 80,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
		DisableThreatAvoidance:    true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fallback interaction target selection")
	}
	if snapshot.Interaction.Kind != "fight_resolved" {
		t.Fatalf("expected same-shape fallback to produce fight_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != "player-1" || snapshot.Interaction.TargetID != simulation.DefaultAutonomousID {
		t.Fatalf("expected fallback to nearest same-shape target %q -> %q, got %q -> %q", "player-1", simulation.DefaultAutonomousID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestInteractionSeekingSkipsInfeasibleDifferentShapeTarget(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   220,
		PlayerY:                   500,
		PlayerEnergy:              14,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "triangle",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      124,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 80,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fallback target interaction after skipping infeasible reproduction target")
	}
	if snapshot.Interaction.Kind != "fight_resolved" {
		t.Fatalf("expected infeasible different-shape target to be skipped in favor of same-shape fight, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected fallback to nearest feasible target %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestInteractionSeekingSkipsLosingSameShapeTarget(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  180,
		PlayerY:                  500,
		PlayerEnergy:             100,
		PlayerChildrenCount:      2,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         40,
		AutonomousChildrenCount:  0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction while losing same-shape and blocked reproduction targets are skipped, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= 100 {
		t.Fatalf("expected losing autonomous circle to retreat from nearby stronger same-shape threat, got x=%v", snapshot.AutonomousCircles[0].X)
	}
}

func TestNearbyBlockedDifferentShapePlayerTriggersAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "square",
		PlayerX:                  180,
		PlayerY:                  500,
		PlayerEnergy:             simulation.DefaultReproductionCost - 1,
		PlayerChildrenCount:      0,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         100,
		AutonomousChildrenCount:  0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during blocked-reproduction avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= 100 {
		t.Fatalf("expected autonomous circle to retreat from nearby blocked different-shape player, got x=%v", snapshot.AutonomousCircles[0].X)
	}
}

func TestNearbyBlockedDifferentShapeAutonomousTriggersAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      180,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: simulation.DefaultReproductionCost - 1,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during blocked-reproduction avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= 100 {
		t.Fatalf("expected autonomous circle to retreat from nearby blocked different-shape autonomous target, got x=%v", snapshot.AutonomousCircles[0].X)
	}
}

func TestFeasibleDifferentShapeTargetDoesNotTriggerBlockedReproductionAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      180,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected feasible different-shape pursuit to remain active")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected feasible different-shape target to remain reproducible, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected pursuit of feasible different-shape target %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestDistantBlockedDifferentShapeTargetDoesNotOverrideOrdinarySteering(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "square",
		PlayerX:                  260,
		PlayerY:                  500,
		PlayerEnergy:             simulation.DefaultReproductionCost - 1,
		PlayerChildrenCount:      0,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         100,
		AutonomousChildrenCount:  0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction while blocked target stays outside avoidance radius, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].Y >= 500 {
		t.Fatalf("expected distant blocked target to leave ordinary food steering in control, got y=%v", snapshot.AutonomousCircles[0].Y)
	}
}

func TestAttachedChildCanMakeDifferentShapeTargetEffectivelyNearestForPursuit(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "square",
		PlayerX:                             225,
		PlayerY:                             500,
		PlayerEnergy:                        100,
		PlayerChildrenCount:                 1,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "square",
		AutonomousX:                         100,
		AutonomousY:                         500,
		SecondaryAutonomousX:                100,
		SecondaryAutonomousY:                620,
		AutonomousEnergy:                    100,
		SecondaryAutonomousEnergy:           100,
		AutonomousChildrenCount:             0,
		SecondaryChildrenCount:              0,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	if before.Player == nil {
		t.Fatal("expected player in initial snapshot")
	}
	playerParentDistance := math.Hypot(before.Player.X-before.AutonomousCircles[0].X, before.Player.Y-before.AutonomousCircles[0].Y)
	secondaryDistance := math.Hypot(before.AutonomousCircles[1].X-before.AutonomousCircles[0].X, before.AutonomousCircles[1].Y-before.AutonomousCircles[0].Y)
	if playerParentDistance <= secondaryDistance {
		t.Fatalf("expected player parent body to be farther than secondary target, player=%v secondary=%v", playerParentDistance, secondaryDistance)
	}

	playerChildDistance := playerParentDistance
	for _, child := range before.Player.AttachedChildren {
		distance := math.Hypot(child.X-before.AutonomousCircles[0].X, child.Y-before.AutonomousCircles[0].Y)
		if distance < playerChildDistance {
			playerChildDistance = distance
		}
	}
	if playerChildDistance >= secondaryDistance {
		t.Fatalf("expected player attached child to be effectively nearer than secondary target, child=%v secondary=%v", playerChildDistance, secondaryDistance)
	}

	snapshot := session.Advance()
	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during child-aware pursuit, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X <= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to move right toward player effective target, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAttachedChildCanMakeWinningSameShapeFallbackEffectivelyNearest(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		PlayerX:                             225,
		PlayerY:                             500,
		PlayerEnergy:                        20,
		PlayerChildrenCount:                 1,
		AutonomousShape:                     "triangle",
		SecondaryAutonomousShape:            "triangle",
		AutonomousX:                         100,
		AutonomousY:                         500,
		SecondaryAutonomousX:                100,
		SecondaryAutonomousY:                620,
		AutonomousEnergy:                    100,
		SecondaryAutonomousEnergy:           20,
		AutonomousChildrenCount:             1,
		SecondaryChildrenCount:              0,
		DisableThreatAvoidance:              true,
		DisableBlockedReproductionAvoidance: true,
	})

	before := session.Snapshot()
	if before.Player == nil {
		t.Fatal("expected player in initial snapshot")
	}
	playerParentDistance := math.Hypot(before.Player.X-before.AutonomousCircles[0].X, before.Player.Y-before.AutonomousCircles[0].Y)
	secondaryDistance := math.Hypot(before.AutonomousCircles[1].X-before.AutonomousCircles[0].X, before.AutonomousCircles[1].Y-before.AutonomousCircles[0].Y)
	if playerParentDistance <= secondaryDistance {
		t.Fatalf("expected player parent body to be farther than secondary fallback target, player=%v secondary=%v", playerParentDistance, secondaryDistance)
	}

	playerChildDistance := playerParentDistance
	for _, child := range before.Player.AttachedChildren {
		distance := math.Hypot(child.X-before.AutonomousCircles[0].X, child.Y-before.AutonomousCircles[0].Y)
		if distance < playerChildDistance {
			playerChildDistance = distance
		}
	}
	if playerChildDistance >= secondaryDistance {
		t.Fatalf("expected player attached child to be effectively nearer than secondary fallback target, child=%v secondary=%v", playerChildDistance, secondaryDistance)
	}

	snapshot := session.Advance()
	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during child-aware same-shape fallback pursuit, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X <= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to move right toward player fallback target, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAttachedChildCanTriggerSameShapeThreatAvoidanceBeforeParentOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  230,
		PlayerY:                  500,
		PlayerEnergy:             100,
		PlayerChildrenCount:      4,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         40,
		AutonomousChildrenCount:  0,
	})

	before := session.Snapshot()
	if before.Player == nil {
		t.Fatal("expected player in initial snapshot")
	}
	parentDistance := math.Hypot(before.Player.X-before.AutonomousCircles[0].X, before.Player.Y-before.AutonomousCircles[0].Y)
	if parentDistance < simulation.DefaultThreatAvoidanceDistance {
		t.Fatalf("expected parent body to stay outside threat window, got distance=%v", parentDistance)
	}
	childInsideWindow := false
	for _, child := range before.Player.AttachedChildren {
		distance := math.Hypot(child.X-before.AutonomousCircles[0].X, child.Y-before.AutonomousCircles[0].Y)
		if distance < simulation.DefaultThreatAvoidanceDistance {
			childInsideWindow = true
			break
		}
	}
	if !childInsideWindow {
		t.Fatal("expected attached child to enter threat window before parent body")
	}

	snapshot := session.Advance()
	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during child-triggered threat avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to retreat from child-triggered same-shape threat, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestAttachedChildCanTriggerBlockedReproductionAvoidanceBeforeParentOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "square",
		PlayerX:                  225,
		PlayerY:                  500,
		PlayerEnergy:             1,
		PlayerChildrenCount:      1,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         100,
		AutonomousChildrenCount:  0,
	})

	before := session.Snapshot()
	if before.Player == nil {
		t.Fatal("expected player in initial snapshot")
	}
	parentDistance := math.Hypot(before.Player.X-before.AutonomousCircles[0].X, before.Player.Y-before.AutonomousCircles[0].Y)
	if parentDistance < simulation.DefaultBlockedReproductionAvoidanceDistance {
		t.Fatalf("expected parent body to stay outside blocked-reproduction window, got distance=%v", parentDistance)
	}
	childInsideWindow := false
	for _, child := range before.Player.AttachedChildren {
		distance := math.Hypot(child.X-before.AutonomousCircles[0].X, child.Y-before.AutonomousCircles[0].Y)
		if distance < simulation.DefaultBlockedReproductionAvoidanceDistance {
			childInsideWindow = true
			break
		}
	}
	if !childInsideWindow {
		t.Fatal("expected attached child to enter blocked-reproduction window before parent body")
	}

	snapshot := session.Advance()
	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during child-triggered blocked-reproduction avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= before.AutonomousCircles[0].X {
		t.Fatalf("expected autonomous circle to retreat from child-triggered blocked-reproduction target, before=%v after=%v", before.AutonomousCircles[0].X, snapshot.AutonomousCircles[0].X)
	}
}

func TestInteractionSeekingPrefersWinningSameShapeFallback(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  220,
		PlayerY:                  500,
		PlayerEnergy:             60,
		PlayerChildrenCount:      0,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         100,
		AutonomousChildrenCount:  1,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected winning same-shape fallback interaction")
	}
	if snapshot.Interaction.Kind != "fight_resolved" && snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected winning same-shape fallback to produce fight outcome, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != "player-1" || snapshot.Interaction.TargetID != simulation.DefaultAutonomousID {
		t.Fatalf("expected winning fallback against %q, got %q -> %q", simulation.DefaultAutonomousID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
}

func TestNearbyStrongerSameShapePlayerTriggersThreatAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  180,
		PlayerY:                  500,
		PlayerEnergy:             100,
		PlayerChildrenCount:      2,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         40,
		AutonomousChildrenCount:  0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during threat avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= 100 {
		t.Fatalf("expected autonomous circle to move away from stronger nearby player threat, got x=%v", snapshot.AutonomousCircles[0].X)
	}
}

func TestNearbyStrongerSameShapeAutonomousTriggersThreatAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "triangle",
		AutonomousX:               100,
		AutonomousY:               500,
		SecondaryAutonomousX:      180,
		SecondaryAutonomousY:      500,
		AutonomousEnergy:          40,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    2,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction during autonomous threat avoidance, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].X >= 100 {
		t.Fatalf("expected autonomous circle to move away from stronger nearby autonomous threat, got x=%v", snapshot.AutonomousCircles[0].X)
	}
}

func TestDistantThreatDoesNotOverrideOrdinarySteering(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  260,
		PlayerY:                  500,
		PlayerEnergy:             100,
		PlayerChildrenCount:      2,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         40,
		AutonomousChildrenCount:  0,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no immediate interaction while distant threat stays outside avoidance radius, got %+v", snapshot.Interaction)
	}
	if snapshot.AutonomousCircles[0].Y >= 500 {
		t.Fatalf("expected distant threat to leave ordinary food fallback in control, got y=%v", snapshot.AutonomousCircles[0].Y)
	}
}

func TestWeakerSameShapeTargetDoesNotTriggerThreatAvoidance(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		PlayerX:                  180,
		PlayerY:                  500,
		PlayerEnergy:             20,
		PlayerChildrenCount:      0,
		AutonomousShape:          "triangle",
		SecondaryAutonomousShape: "",
		AutonomousX:              100,
		AutonomousY:              500,
		AutonomousEnergy:         100,
		AutonomousChildrenCount:  1,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected ordinary same-shape pursuit when target is not threatening")
	}
	if snapshot.Interaction.Kind != "fight_resolved" && snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected ordinary fight outcome, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.TargetID != simulation.DefaultAutonomousID {
		t.Fatalf("expected pursuit of non-threatening same-shape target %q, got %q", simulation.DefaultAutonomousID, snapshot.Interaction.TargetID)
	}
}

func TestDefaultWorldSupportsSameShapeFightPath(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               simulation.DefaultPlayerShape,
		AutonomousShape:           simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:  simulation.DefaultAutoShape,
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:       1,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:   1,
		DisableThreatAvoidance:    true,
	})
	before := session.Snapshot()
	expectedChildID := before.AutonomousCircles[0].AttachedChildren[0].ID

	var snapshot simulation.Snapshot
	for range 20 {
		session.ApplyIntent(simulation.Vector{X: -1, Y: 0})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected same-shape interaction in default world")
	}
	if snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.TargetID != simulation.DefaultAutonomousID {
		t.Fatalf("expected fight against %q, got %q", simulation.DefaultAutonomousID, snapshot.Interaction.TargetID)
	}
	if snapshot.Interaction.AbsorbedChildID != expectedChildID {
		t.Fatalf("expected absorbed child id %q, got %q", expectedChildID, snapshot.Interaction.AbsorbedChildID)
	}
	if len(snapshot.AutonomousCircles) != 2 {
		t.Fatalf("expected both autonomous circles to remain, got %d circles", len(snapshot.AutonomousCircles))
	}
	foundAbsorbedLoser := false
	for _, circle := range snapshot.AutonomousCircles {
		if circle.ID != simulation.DefaultAutonomousID {
			continue
		}
		foundAbsorbedLoser = true
		if len(circle.AttachedChildren) != 0 {
			t.Fatalf("expected absorbed conflict to remove one child, got %d", len(circle.AttachedChildren))
		}
		if len(circle.AttachedChildren) != 0 {
			t.Fatalf("expected absorbed conflict to remove visible attached child, got %d", len(circle.AttachedChildren))
		}
		if circle.Generation != before.AutonomousCircles[0].Generation {
			t.Fatalf("expected absorbed conflict to keep generation %d, got %d", before.AutonomousCircles[0].Generation, circle.Generation)
		}
		break
	}
	if !foundAbsorbedLoser {
		t.Fatalf("expected autonomous loser %q to remain active after child absorption", simulation.DefaultAutonomousID)
	}
}

func TestDefaultWorldSupportsDifferentShapeReproductionPath(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               simulation.DefaultPlayerShape,
		AutonomousShape:           simulation.DefaultPlayerShape,
		SecondaryAutonomousShape:  simulation.DefaultAutoShape,
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: simulation.DefaultPlayerEnergy,
		AutonomousChildrenCount:   1,
		DisableFoodSeeking:        true,
	})
	before := session.Snapshot()

	var snapshot simulation.Snapshot
	for range 20 {
		session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected different-shape interaction in default world")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.AutonomousCircles) != len(before.AutonomousCircles) {
		t.Fatalf("expected no new autonomous circles from reproduction, before=%d after=%d", len(before.AutonomousCircles), len(snapshot.AutonomousCircles))
	}
	if snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected reproduction against %q, got %q", simulation.DefaultSecondaryID, snapshot.Interaction.TargetID)
	}
	if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[1].AttachedChildren) != len(before.Player.AttachedChildren)+len(before.AutonomousCircles[1].AttachedChildren)+2 {
		t.Fatalf("expected reproduction to distribute two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", len(before.Player.AttachedChildren), len(before.AutonomousCircles[1].AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[1].AttachedChildren))
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through movement and reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) {
		t.Fatalf("expected player attached children to match count after reproduction, count=%d attached=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[1].AttachedChildren) != len(snapshot.AutonomousCircles[1].AttachedChildren) {
		t.Fatalf("expected autonomous attached children to match count after reproduction, count=%d attached=%d", len(snapshot.AutonomousCircles[1].AttachedChildren), len(snapshot.AutonomousCircles[1].AttachedChildren))
	}
	if snapshot.AutonomousCircles[1].Energy >= before.AutonomousCircles[1].Energy {
		t.Fatalf("expected autonomous energy to decrease through movement and reproduction, before=%v after=%v", before.AutonomousCircles[1].Energy, snapshot.AutonomousCircles[1].Energy)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed player radius after reproduction, got %v", snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[1].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed autonomous radius after reproduction, got %v", snapshot.AutonomousCircles[1].Radius)
	}
}

func TestSameShapeOverlapProducesFightCandidate(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:            "triangle",
		AutonomousShape:        "triangle",
		PlayerEnergy:           100,
		AutonomousEnergy:       100,
		DisableThreatAvoidance: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected active interaction classification")
	}
	if snapshot.Interaction.Kind != "fight_resolved" {
		t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.WinnerID != "player-1" {
		t.Fatalf("expected player to win tie-break fight, got winner %q", snapshot.Interaction.WinnerID)
	}
	if len(snapshot.AutonomousCircles) != 0 {
		t.Fatalf("expected losing autonomous circle to be removed, got %d autonomous circles", len(snapshot.AutonomousCircles))
	}
}

func TestAutonomousLoserWithChildAbsorbsFightLossAndRemainsActive(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        80,
		AutonomousChildrenCount: 1,
		DisableThreatAvoidance:  true,
	})
	before := session.Snapshot()
	expectedChildID := before.AutonomousCircles[0].AttachedChildren[0].ID

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.LoserID != simulation.DefaultAutonomousID {
		t.Fatalf("expected autonomous loser, got %q", snapshot.Interaction.LoserID)
	}
	if snapshot.Interaction.AbsorbedChildID != expectedChildID {
		t.Fatalf("expected absorbed child id %q, got %q", expectedChildID, snapshot.Interaction.AbsorbedChildID)
	}
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous loser to remain active, got %d circles", len(snapshot.AutonomousCircles))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected absorbed conflict to consume one child, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected visible attached child removal, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if snapshot.AutonomousCircles[0].Energy >= before.AutonomousCircles[0].Energy {
		t.Fatalf("expected loser energy to reflect ordinary movement, before=%v after=%v", before.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.AutonomousCircles[0].LineageID != before.AutonomousCircles[0].LineageID {
		t.Fatalf("expected lineage to stay %q, got %q", before.AutonomousCircles[0].LineageID, snapshot.AutonomousCircles[0].LineageID)
	}
	if snapshot.AutonomousCircles[0].Generation != before.AutonomousCircles[0].Generation {
		t.Fatalf("expected absorbed conflict to keep generation %d, got %d", before.AutonomousCircles[0].Generation, snapshot.AutonomousCircles[0].Generation)
	}
}

func TestPlayerAttachedChildCanTriggerFightBeforeParentBodiesOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               141,
		AutonomousY:               308,
		PlayerEnergy:              100,
		PlayerChildrenCount:       1,
		AutonomousEnergy:          80,
		AutonomousChildrenCount:   1,
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})

	var snapshot simulation.Snapshot
	for range 4 {
		session.ApplyIntent(simulation.Vector{X: -1, Y: 0})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected child-triggered fight interaction")
	}
	if snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player == nil {
		t.Fatal("expected player to remain active after winning child-triggered fight")
	}
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous loser to remain active after child absorption, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.Interaction.ContactOrigin != "attached_child" {
		t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
	}
	if snapshot.Interaction.ContactPathKind != "source_child_to_target_parent" {
		t.Fatalf("expected source_child_to_target_parent path kind, got %q", snapshot.Interaction.ContactPathKind)
	}
	if snapshot.Interaction.SourceChildID != "player-1-child-1" {
		t.Fatalf("expected source child id player-1-child-1, got %q", snapshot.Interaction.SourceChildID)
	}
	if snapshot.Interaction.TargetChildID != "" {
		t.Fatalf("expected empty target child id, got %q", snapshot.Interaction.TargetChildID)
	}
}

func TestAutonomousAttachedChildCanTriggerReproductionBeforeParentBodiesOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               142,
		AutonomousY:               308,
		PlayerEnergy:              100,
		PlayerChildrenCount:       1,
		AutonomousEnergy:          100,
		AutonomousChildrenCount:   0,
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})

	var snapshot simulation.Snapshot
	for range 4 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected child-triggered reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player == nil {
		t.Fatal("expected player to remain active after reproduction")
	}
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous parent to remain active, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.Interaction.ContactOrigin != "attached_child" {
		t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
	}
	if snapshot.Interaction.ContactPathKind != "source_child_to_target_parent" {
		t.Fatalf("expected source_child_to_target_parent path kind, got %q", snapshot.Interaction.ContactPathKind)
	}
	if snapshot.Interaction.SourceChildID != "player-1-child-1" {
		t.Fatalf("expected source child id player-1-child-1, got %q", snapshot.Interaction.SourceChildID)
	}
	if snapshot.Interaction.TargetChildID != "" {
		t.Fatalf("expected empty target child id, got %q", snapshot.Interaction.TargetChildID)
	}
}

func TestAttachedChildrenCanTriggerFightThroughChildToChildContact(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               238.511,
		AutonomousY:               294.944,
		PlayerEnergy:              100,
		PlayerChildrenCount:       2,
		AutonomousEnergy:          80,
		AutonomousChildrenCount:   2,
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})

	snapshot := session.Advance()

	if snapshot.Interaction == nil {
		t.Fatal("expected child-to-child fight interaction")
	}
	if snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.ContactOrigin != "attached_child" {
		t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
	}
	if snapshot.Interaction.ContactPathKind != "child_to_child" {
		t.Fatalf("expected child_to_child path kind, got %q", snapshot.Interaction.ContactPathKind)
	}
	if snapshot.Interaction.SourceChildID != "player-1-child-2" {
		t.Fatalf("expected source child id player-1-child-2, got %q", snapshot.Interaction.SourceChildID)
	}
	if snapshot.Interaction.TargetChildID != "circle-2-child-1" {
		t.Fatalf("expected target child id circle-2-child-1, got %q", snapshot.Interaction.TargetChildID)
	}
	if snapshot.Player == nil || len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected both parents to remain after absorbed fight, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != 1 {
		t.Fatalf("expected autonomous loser to lose one child, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestAttachedChildrenCanTriggerReproductionThroughChildToChildContact(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               238.511,
		AutonomousY:               294.944,
		PlayerEnergy:              100,
		PlayerChildrenCount:       2,
		AutonomousEnergy:          100,
		AutonomousChildrenCount:   2,
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})
	before := session.Snapshot()

	snapshot := session.Advance()

	if snapshot.Interaction == nil {
		t.Fatal("expected child-to-child reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.ContactOrigin != "attached_child" {
		t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
	}
	if snapshot.Interaction.ContactPathKind != "child_to_child" {
		t.Fatalf("expected child_to_child path kind, got %q", snapshot.Interaction.ContactPathKind)
	}
	if snapshot.Interaction.SourceChildID != "player-1-child-2" {
		t.Fatalf("expected source child id player-1-child-2, got %q", snapshot.Interaction.SourceChildID)
	}
	if snapshot.Interaction.TargetChildID != "circle-2-child-1" {
		t.Fatalf("expected target child id circle-2-child-1, got %q", snapshot.Interaction.TargetChildID)
	}
	if snapshot.Player == nil || len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected both parents to remain after reproduction, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
	}
	if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != len(before.Player.AttachedChildren)+len(before.AutonomousCircles[0].AttachedChildren)+2 {
		t.Fatalf("expected reproduction to add two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", len(before.Player.AttachedChildren), len(before.AutonomousCircles[0].AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestDerivedRadiusAloneDoesNotStartContactWithoutVisibleOverlap(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               227,
		AutonomousY:               308,
		PlayerEnergy:              100,
		PlayerChildrenCount:       2,
		AutonomousEnergy:          0,
		AutonomousChildrenCount:   0,
		SecondaryAutonomousEnergy: 0,
		DisableFoodSeeking:        true,
	})

	snapshot := session.Advance()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no interaction when only derived parent radius would have overlapped, got %+v", snapshot.Interaction)
	}
}

func TestParentBodyContactDoesNotExposeChildIdentity(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:            "triangle",
		AutonomousShape:        "triangle",
		PlayerEnergy:           100,
		AutonomousEnergy:       100,
		DisableThreatAvoidance: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected parent-body fight interaction")
	}
	if snapshot.Interaction.ContactOrigin != "parent_body" {
		t.Fatalf("expected parent_body contact origin, got %q", snapshot.Interaction.ContactOrigin)
	}
	if snapshot.Interaction.ContactPathKind != "" {
		t.Fatalf("expected empty contact path kind, got %q", snapshot.Interaction.ContactPathKind)
	}
	if snapshot.Interaction.SourceChildID != "" {
		t.Fatalf("expected empty source child id, got %q", snapshot.Interaction.SourceChildID)
	}
	if snapshot.Interaction.TargetChildID != "" {
		t.Fatalf("expected empty target child id, got %q", snapshot.Interaction.TargetChildID)
	}
}

func TestAutonomousCirclesCanResolveFightWithoutPlayerInvolvement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "triangle",
		AutonomousX:               200,
		AutonomousY:               300,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      300,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 80,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
		DisableFoodSeeking:        true,
	})
	before := session.Snapshot()

	snapshot := session.Advance()

	if snapshot.Interaction == nil {
		t.Fatal("expected autonomous-only fight interaction")
	}
	if snapshot.Interaction.Kind != "fight_resolved" {
		t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if snapshot.Interaction.WinnerID != simulation.DefaultAutonomousID {
		t.Fatalf("expected higher-energy autonomous winner %q, got %q", simulation.DefaultAutonomousID, snapshot.Interaction.WinnerID)
	}
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected losing autonomous circle to be removed, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.Player == nil || snapshot.Player.X != before.Player.X || snapshot.Player.Y != before.Player.Y {
		t.Fatalf("expected player to remain uninvolved at %+v, got %+v", before.Player, snapshot.Player)
	}
}

func TestAutonomousCirclesCanResolveReproductionWithoutPlayerInvolvement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "square",
		AutonomousX:               200,
		AutonomousY:               300,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      300,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
		DisableFoodSeeking:        true,
	})
	before := session.Snapshot()

	snapshot := session.Advance()

	if snapshot.Interaction == nil {
		t.Fatal("expected autonomous-only reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if len(snapshot.AutonomousCircles) != 2 {
		t.Fatalf("expected both autonomous parents to remain, got %d circles", len(snapshot.AutonomousCircles))
	}
	left, found := autonomousByID(snapshot.AutonomousCircles, simulation.DefaultAutonomousID)
	if !found {
		t.Fatalf("expected autonomous circle %q after reproduction", simulation.DefaultAutonomousID)
	}
	right, found := autonomousByID(snapshot.AutonomousCircles, simulation.DefaultSecondaryID)
	if !found {
		t.Fatalf("expected autonomous circle %q after reproduction", simulation.DefaultSecondaryID)
	}
	beforeLeft, _ := autonomousByID(before.AutonomousCircles, simulation.DefaultAutonomousID)
	beforeRight, _ := autonomousByID(before.AutonomousCircles, simulation.DefaultSecondaryID)
	sourceCreated := newlyOwnedChildIDs(beforeLeft.AttachedChildren, left.AttachedChildren)
	targetCreated := newlyOwnedChildIDs(beforeRight.AttachedChildren, right.AttachedChildren)
	assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, left.Energy, len(left.AttachedChildren)))
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, right.Energy, len(right.AttachedChildren)))
	assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
	assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
	assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
	assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
	if len(left.AttachedChildren)+len(right.AttachedChildren) != len(beforeLeft.AttachedChildren)+len(beforeRight.AttachedChildren)+2 {
		t.Fatalf("expected reproduction to add two children across the autonomous pair, before left=%d right=%d after left=%d right=%d", len(beforeLeft.AttachedChildren), len(beforeRight.AttachedChildren), len(left.AttachedChildren), len(right.AttachedChildren))
	}
	assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
	if snapshot.Player == nil || snapshot.Player.X != before.Player.X || snapshot.Player.Y != before.Player.Y {
		t.Fatalf("expected player to remain uninvolved at %+v, got %+v", before.Player, snapshot.Player)
	}
}

func TestDifferentShapeOverlapProducesResolvedReproduction(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "square")
	before := session.Snapshot()
	var snapshot simulation.Snapshot
	for range 20 {
		before = session.Snapshot()
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected active interaction classification")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.Interaction.CreatedChildIDs) != 2 {
		t.Fatalf("expected two created child ids, got %d", len(snapshot.Interaction.CreatedChildIDs))
	}
	sourceCreated := newlyOwnedChildIDs(before.Player.AttachedChildren, snapshot.Player.AttachedChildren)
	targetCreated := newlyOwnedChildIDs(before.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren)
	if len(snapshot.Interaction.SourceCreatedChildIDs) != len(sourceCreated) {
		t.Fatalf("expected %d source created child ids, got %d", len(sourceCreated), len(snapshot.Interaction.SourceCreatedChildIDs))
	}
	if len(snapshot.Interaction.TargetCreatedChildIDs) != len(targetCreated) {
		t.Fatalf("expected %d target created child ids, got %d", len(targetCreated), len(snapshot.Interaction.TargetCreatedChildIDs))
	}
	createdChildIDs := map[string]struct{}{}
	for _, childID := range snapshot.Interaction.CreatedChildIDs {
		createdChildIDs[childID] = struct{}{}
	}
	for _, childID := range newlyCreatedChildIDs(before.Player.AttachedChildren, snapshot.Player.AttachedChildren, before.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren) {
		if _, exists := createdChildIDs[childID]; !exists {
			t.Fatalf("expected created child id %q in interaction payload, got %+v", childID, snapshot.Interaction.CreatedChildIDs)
		}
	}
	assertChildIDSetEqual(t, snapshot.Interaction.SourceCreatedChildIDs, sourceCreated, "source created child ids")
	assertChildIDSetEqual(t, snapshot.Interaction.TargetCreatedChildIDs, targetCreated, "target created child ids")
	assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
	if len(snapshot.AutonomousCircles) != len(before.AutonomousCircles) {
		t.Fatalf("expected no new autonomous circles, before=%d after=%d", len(before.AutonomousCircles), len(snapshot.AutonomousCircles))
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != len(before.Player.AttachedChildren)+len(before.AutonomousCircles[0].AttachedChildren)+2 {
		t.Fatalf("expected reproduction to distribute two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", len(before.Player.AttachedChildren), len(before.AutonomousCircles[0].AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
		t.Fatalf("expected autonomous attached children to match count, count=%d attached=%d", len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if snapshot.AutonomousCircles[0].Energy >= before.AutonomousCircles[0].Energy {
		t.Fatalf("expected autonomous energy to decrease through reproduction, before=%v after=%v", before.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed player radius after reproduction, got %v", snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed autonomous radius after reproduction, got %v", snapshot.AutonomousCircles[0].Radius)
	}
}

func TestDifferentShapeOverlapBlocksReproductionWhenEnergyIsInsufficient(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		DisableBlockedReproductionAvoidance: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_blocked_energy" {
		t.Fatalf("expected reproduce_blocked_energy, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.Interaction.CreatedChildIDs) != 0 {
		t.Fatalf("expected no created child ids on blocked reproduction, got %+v", snapshot.Interaction.CreatedChildIDs)
	}
	if len(snapshot.Interaction.SourceCreatedChildIDs) != 0 || len(snapshot.Interaction.TargetCreatedChildIDs) != 0 {
		t.Fatalf("expected no created child ownership on blocked reproduction, source=%+v target=%+v", snapshot.Interaction.SourceCreatedChildIDs, snapshot.Interaction.TargetCreatedChildIDs)
	}
	if snapshot.Interaction.DistributionKind != "" {
		t.Fatalf("expected no distribution kind on blocked reproduction, got %q", snapshot.Interaction.DistributionKind)
	}
	if snapshot.Interaction.SourceBlockedCapacity {
		t.Fatal("expected source side to have enough reproduction capacity")
	}
	if !snapshot.Interaction.TargetBlockedCapacity {
		t.Fatal("expected target side to be marked as blocked")
	}
	if snapshot.Interaction.SourceCapacityValue < simulation.DefaultReproductionMinEnergy {
		t.Fatalf("expected source capacity to meet reproduction threshold, got %v", snapshot.Interaction.SourceCapacityValue)
	}
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
	assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
	assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
	assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
	assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected blocked reproduction to avoid spawning autonomous circles, got %d autonomous circles", len(snapshot.AutonomousCircles))
	}
	if len(snapshot.Player.AttachedChildren) != 0 || len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected blocked reproduction to preserve child counts, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if len(snapshot.Player.AttachedChildren) != 0 || len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected blocked reproduction to preserve attached children, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestDifferentShapeOverlapBlocksReproductionWhenPlayerCapacityIsInsufficient(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultReproductionCost - 1,
		AutonomousEnergy:                    simulation.DefaultPlayerEnergy,
		DisableBlockedReproductionAvoidance: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_blocked_energy" {
		t.Fatalf("expected reproduce_blocked_energy, got %q", snapshot.Interaction.Kind)
	}
	if !snapshot.Interaction.SourceBlockedCapacity {
		t.Fatal("expected source side to be marked as blocked")
	}
	if snapshot.Interaction.TargetBlockedCapacity {
		t.Fatal("expected target side to have enough reproduction capacity")
	}
	assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
	assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
	assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
	assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
	assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
}

func TestDifferentShapeOverlapBlocksReproductionWhenBothSidesLackCapacity(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultReproductionCost - 1,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		DisableBlockedReproductionAvoidance: true,
	})
	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_blocked_energy" {
		t.Fatalf("expected reproduce_blocked_energy, got %q", snapshot.Interaction.Kind)
	}
	if !snapshot.Interaction.SourceBlockedCapacity || !snapshot.Interaction.TargetBlockedCapacity {
		t.Fatalf("expected both sides to be marked as blocked, source=%v target=%v", snapshot.Interaction.SourceBlockedCapacity, snapshot.Interaction.TargetBlockedCapacity)
	}
	assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
	assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
	assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
	assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
	assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
}

func TestDifferentShapeOverlapConsumesChildAsReproductionPayment(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerX:                             200,
		PlayerY:                             300,
		AutonomousX:                         228,
		AutonomousY:                         300,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		AutonomousChildrenCount:             1,
		DisableBlockedReproductionAvoidance: true,
		DisableFoodSeeking:                  true,
	})

	var snapshot simulation.Snapshot
	before := session.Snapshot()
	for range 20 {
		before = session.Snapshot()
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_paid_child" {
		t.Fatalf("expected reproduce_paid_child, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.Interaction.CreatedChildIDs) != 2 {
		t.Fatalf("expected two created child ids, got %d", len(snapshot.Interaction.CreatedChildIDs))
	}
	sourceCreated := newlyOwnedChildIDs(before.Player.AttachedChildren, snapshot.Player.AttachedChildren)
	targetCreated := newlyOwnedChildIDs(before.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren)
	if snapshot.Interaction.SourceBlockedCapacity || snapshot.Interaction.TargetBlockedCapacity {
		t.Fatalf("expected no blocked-capacity flags on successful reproduction, source=%v target=%v", snapshot.Interaction.SourceBlockedCapacity, snapshot.Interaction.TargetBlockedCapacity)
	}
	assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
	if snapshot.Interaction.SourcePaidChild {
		t.Fatal("expected player not to pay with child in this reproduction path")
	}
	if !snapshot.Interaction.TargetPaidChild {
		t.Fatal("expected autonomous target to pay with child in this reproduction path")
	}
	if snapshot.Interaction.SourcePaidChildID != "" {
		t.Fatalf("expected no source paid child id, got %q", snapshot.Interaction.SourcePaidChildID)
	}
	if snapshot.Interaction.TargetPaidChildID != "circle-2-child-1" {
		t.Fatalf("expected target paid child id circle-2-child-1, got %q", snapshot.Interaction.TargetPaidChildID)
	}
	createdChildIDs := map[string]struct{}{}
	for _, childID := range snapshot.Interaction.CreatedChildIDs {
		createdChildIDs[childID] = struct{}{}
	}
	for _, childID := range newlyCreatedChildIDs(before.Player.AttachedChildren, snapshot.Player.AttachedChildren, before.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren) {
		if _, exists := createdChildIDs[childID]; !exists {
			t.Fatalf("expected created child id %q in interaction payload, got %+v", childID, snapshot.Interaction.CreatedChildIDs)
		}
	}
	assertChildIDSetEqual(t, snapshot.Interaction.SourceCreatedChildIDs, sourceCreated, "source created child ids")
	assertChildIDSetEqual(t, snapshot.Interaction.TargetCreatedChildIDs, targetCreated, "target created child ids")
	assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
	if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != len(before.Player.AttachedChildren)+len(before.AutonomousCircles[0].AttachedChildren)+1 {
		t.Fatalf("expected one child to be spent and two children to be redistributed, before player=%d autonomous=%d after player=%d autonomous=%d", len(before.Player.AttachedChildren), len(before.AutonomousCircles[0].AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	expectedEnergy := before.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
	if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
		t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
	}
	if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) || len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
		t.Fatalf("expected attached children to match counts after child-payment reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestDifferentShapeOverlapPaidByEnergyRemainsOrdinaryResolvedReproduction(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultPlayerEnergy,
		DisableBlockedReproductionAvoidance: true,
	})
	before := session.Snapshot()

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction interaction")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved for energy-paid reproduction, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.Interaction.CreatedChildIDs) != 2 {
		t.Fatalf("expected two created child ids for energy-paid reproduction, got %d", len(snapshot.Interaction.CreatedChildIDs))
	}
	sourceCreated := newlyOwnedChildIDs(before.Player.AttachedChildren, snapshot.Player.AttachedChildren)
	targetCreated := newlyOwnedChildIDs(before.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren)
	if snapshot.Interaction.SourcePaidChild || snapshot.Interaction.TargetPaidChild {
		t.Fatalf("expected ordinary resolved reproduction to omit child-payment identity, got source=%v target=%v", snapshot.Interaction.SourcePaidChild, snapshot.Interaction.TargetPaidChild)
	}
	assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
	assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
	if snapshot.Interaction.SourcePaidChildID != "" || snapshot.Interaction.TargetPaidChildID != "" {
		t.Fatalf("expected ordinary resolved reproduction to omit paid child ids, got source=%q target=%q", snapshot.Interaction.SourcePaidChildID, snapshot.Interaction.TargetPaidChildID)
	}
	assertChildIDSetEqual(t, snapshot.Interaction.SourceCreatedChildIDs, sourceCreated, "source created child ids")
	assertChildIDSetEqual(t, snapshot.Interaction.TargetCreatedChildIDs, targetCreated, "target created child ids")
	assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
}

func assertDistributionKindMatchesOwnership(t *testing.T, got string, sourceCreated []string, targetCreated []string) {
	t.Helper()

	want := ""
	switch {
	case len(sourceCreated) == 2 && len(targetCreated) == 0:
		want = "source_only"
	case len(sourceCreated) == 1 && len(targetCreated) == 1:
		want = "split"
	case len(sourceCreated) == 0 && len(targetCreated) == 2:
		want = "target_only"
	default:
		t.Fatalf("unexpected reproduction ownership split: source=%d target=%d", len(sourceCreated), len(targetCreated))
	}

	if got != want {
		t.Fatalf("expected distribution kind %q, got %q", want, got)
	}
}

func assertFloatEqual(t *testing.T, got float64, want float64) {
	t.Helper()

	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func expectedReportedCapacity(kind string, energy float64, childrenCount int) float64 {
	if kind == "reproduce_resolved" || kind == "reproduce_paid_child" {
		return energy + simulation.DefaultReproductionCost
	}

	if childrenCount == 0 {
		return energy
	}

	return energy + simulation.DefaultReproductionCost
}

func assertCapacityComponentsMatch(t *testing.T, total float64, energyComponent float64, reserveComponent float64) {
	t.Helper()

	if total != energyComponent+reserveComponent {
		t.Fatalf("expected capacity components to sum to %v, got energy=%v reserve=%v", total, energyComponent, reserveComponent)
	}
}

func newlyCreatedChildIDs(beforeSource []simulation.AttachedChild, afterSource []simulation.AttachedChild, beforeTarget []simulation.AttachedChild, afterTarget []simulation.AttachedChild) []string {
	beforeIDs := make(map[string]struct{}, len(beforeSource)+len(beforeTarget))
	for _, child := range beforeSource {
		beforeIDs[child.ID] = struct{}{}
	}
	for _, child := range beforeTarget {
		beforeIDs[child.ID] = struct{}{}
	}

	created := make([]string, 0, 2)
	for _, child := range afterSource {
		if _, exists := beforeIDs[child.ID]; !exists {
			created = append(created, child.ID)
		}
	}
	for _, child := range afterTarget {
		if _, exists := beforeIDs[child.ID]; !exists {
			created = append(created, child.ID)
		}
	}

	return created
}

func newlyOwnedChildIDs(before []simulation.AttachedChild, after []simulation.AttachedChild) []string {
	beforeIDs := make(map[string]struct{}, len(before))
	for _, child := range before {
		beforeIDs[child.ID] = struct{}{}
	}

	created := make([]string, 0, 2)
	for _, child := range after {
		if _, exists := beforeIDs[child.ID]; !exists {
			created = append(created, child.ID)
		}
	}

	return created
}

func assertChildIDSetEqual(t *testing.T, got []string, expected []string, label string) {
	t.Helper()

	if len(got) != len(expected) {
		t.Fatalf("expected %d %s, got %d", len(expected), label, len(got))
	}

	gotSet := make(map[string]struct{}, len(got))
	for _, childID := range got {
		gotSet[childID] = struct{}{}
	}

	for _, childID := range expected {
		if _, exists := gotSet[childID]; !exists {
			t.Fatalf("expected %s to include %q, got %+v", label, childID, got)
		}
	}
}

func TestContinuousOverlapDoesNotRepeatChildAccumulation(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:        "triangle",
		AutonomousShape:    "square",
		PlayerEnergy:       100,
		AutonomousEnergy:   100,
		DisableFoodSeeking: true,
	})
	originalOpponentID := simulation.DefaultAutonomousID

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected reproduction resolution")
	}
	originalOpponent, found := autonomousByID(snapshot.AutonomousCircles, originalOpponentID)
	if !found {
		t.Fatalf("expected original autonomous circle %q after reproduction", originalOpponentID)
	}
	if len(snapshot.Player.AttachedChildren)+len(originalOpponent.AttachedChildren) != 2 {
		t.Fatalf("expected first reproduction to create two visible child units across the pair, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(originalOpponent.AttachedChildren))
	}

	for range 3 {
		snapshot = session.Advance()
		originalOpponent, found = autonomousByID(snapshot.AutonomousCircles, originalOpponentID)
		if !found {
			t.Fatalf("expected original autonomous circle %q during continuous overlap", originalOpponentID)
		}
		if len(snapshot.Player.AttachedChildren)+len(originalOpponent.AttachedChildren) != 2 {
			t.Fatalf("expected continuous overlap to avoid repeat accumulation, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(originalOpponent.AttachedChildren))
		}
	}
}

func TestPairCanReproduceAgainAfterSeparating(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:        "triangle",
		AutonomousShape:    "square",
		PlayerEnergy:       100,
		AutonomousEnergy:   100,
		DisableFoodSeeking: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected first reproduction resolution")
	}

	separated := false
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction == nil && snapshot.AutonomousCircles[0].X > snapshot.Player.X+snapshot.Player.Radius+snapshot.AutonomousCircles[0].Radius {
			separated = true
			break
		}
	}

	if !separated {
		t.Fatal("expected pair to separate after first reproduction")
	}

	beforePlayerChildren := len(snapshot.Player.AttachedChildren)
	beforeAutonomousChildren := len(snapshot.AutonomousCircles[0].AttachedChildren)
	for range 120 {
		session.ApplyIntent(simulation.Vector{})
		snapshot = session.Advance()
		if snapshot.Interaction != nil && snapshot.Interaction.Kind == "reproduce_resolved" {
			break
		}
		beforePlayerChildren = len(snapshot.Player.AttachedChildren)
		beforeAutonomousChildren = len(snapshot.AutonomousCircles[0].AttachedChildren)
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected second reproduction after re-overlap")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != beforePlayerChildren+beforeAutonomousChildren+2 {
		t.Fatalf(
			"expected second reproduction to add two child units across the pair, before player=%d autonomous=%d after player=%d autonomous=%d",
			beforePlayerChildren,
			beforeAutonomousChildren,
			len(snapshot.Player.AttachedChildren),
			len(snapshot.AutonomousCircles[0].AttachedChildren),
		)
	}
}

func TestAttachedChildrenOrbitTheirParentAcrossTicks(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "square")

	var resolved simulation.Snapshot
	for range 20 {
		resolved = session.Advance()
		if resolved.Interaction != nil {
			break
		}
	}

	if resolved.Interaction == nil || resolved.Interaction.Kind != "reproduce_resolved" {
		t.Fatal("expected resolved reproduction before checking orbit motion")
	}
	if len(resolved.Player.AttachedChildren) == 0 {
		t.Fatal("expected player to own at least one attached child in this deterministic reproduction path")
	}

	firstChild := resolved.Player.AttachedChildren[0]
	later := session.Advance()
	if len(later.Player.AttachedChildren) == 0 {
		t.Fatal("expected attached child to remain with player on later tick")
	}

	laterChild := later.Player.AttachedChildren[0]
	if laterChild.ID != firstChild.ID {
		t.Fatalf("expected same attached child identity across ticks, before=%q after=%q", firstChild.ID, laterChild.ID)
	}
	if laterChild.OwnerID != resolved.Player.ID {
		t.Fatalf("expected child to remain attached to player %q, got owner %q", resolved.Player.ID, laterChild.OwnerID)
	}
	if laterChild.X == firstChild.X && laterChild.Y == firstChild.Y {
		t.Fatalf("expected attached child orbit position to change across ticks, before=(%v,%v) after=(%v,%v)", firstChild.X, firstChild.Y, laterChild.X, laterChild.Y)
	}
}

func TestDerivedParentRadiusDoesNotExpandAttachedChildOrbitDistance(t *testing.T) {
	base := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		SecondaryAutonomousShape: "",
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:      1,
		AutonomousEnergy:         0,
	})
	grown := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		SecondaryAutonomousShape: "",
		PlayerEnergy:             simulation.DefaultPlayerEnergy,
		PlayerChildrenCount:      3,
		AutonomousEnergy:         0,
	})

	baseSnapshot := base.Snapshot()
	grownSnapshot := grown.Snapshot()

	if len(baseSnapshot.Player.AttachedChildren) == 0 || len(grownSnapshot.Player.AttachedChildren) == 0 {
		t.Fatal("expected attached children in both snapshots")
	}

	baseChild := baseSnapshot.Player.AttachedChildren[0]
	grownChild := grownSnapshot.Player.AttachedChildren[0]
	baseDistance := math.Hypot(baseChild.X-baseSnapshot.Player.X, baseChild.Y-baseSnapshot.Player.Y)
	grownDistance := math.Hypot(grownChild.X-grownSnapshot.Player.X, grownChild.Y-grownSnapshot.Player.Y)
	expectedDistance := simulation.DefaultPlayerRadius + simulation.DefaultAttachedChildOrbitGap + simulation.DefaultAttachedChildRadius

	if math.Abs(baseDistance-expectedDistance) > 1e-9 {
		t.Fatalf("expected base orbit distance %v, got %v", expectedDistance, baseDistance)
	}
	if math.Abs(grownDistance-expectedDistance) > 1e-9 {
		t.Fatalf("expected grown orbit distance to stay at visible-body distance %v, got %v", expectedDistance, grownDistance)
	}
}

func TestNonOverlappingCirclesHaveNoInteractionClassification(t *testing.T) {
	session := simulation.NewSession()
	snapshot := session.Snapshot()

	if snapshot.Interaction != nil {
		t.Fatalf("expected no interaction classification, got %+v", snapshot.Interaction)
	}
}

func TestHigherEnergyCircleWinsFight(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:            "triangle",
		AutonomousShape:        "triangle",
		PlayerEnergy:           100,
		AutonomousEnergy:       80,
		DisableThreatAvoidance: true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != "player-1" {
		t.Fatalf("expected higher-energy player to win, got %q", snapshot.Interaction.WinnerID)
	}
}

func TestPlayerWinsExactEqualEnergyAndChildFightTie(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     0,
		AutonomousChildrenCount: 0,
		DisableThreatAvoidance:  true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != "player-1" {
		t.Fatalf("expected player to win exact equal-energy equal-child tie, got %q", snapshot.Interaction.WinnerID)
	}
}

func TestHigherChildCountWinsEqualEnergyFight(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     0,
		AutonomousChildrenCount: 1,
		DisableThreatAvoidance:  true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		session.ApplyIntent(simulation.Vector{X: -1, Y: 0})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != simulation.DefaultAutonomousID {
		t.Fatalf("expected autonomous circle with child presence to win, got %q", snapshot.Interaction.WinnerID)
	}
}

func TestAdditionalChildCountDoesNotStackFightPowerWhenBothSidesHaveChildren(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		PlayerChildrenCount:     3,
		AutonomousEnergy:        100,
		AutonomousChildrenCount: 1,
		DisableFoodSeeking:      true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != "player-1" {
		t.Fatalf("expected exact tie rule to decide once both sides have child presence, got %q", snapshot.Interaction.WinnerID)
	}
}

func TestAutonomousLowerIDWinsExactEqualEnergyAndChildFightTie(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "square",
		PlayerX:                   700,
		PlayerY:                   500,
		PlayerEnergy:              100,
		PlayerChildrenCount:       0,
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "triangle",
		AutonomousX:               200,
		AutonomousY:               300,
		SecondaryAutonomousX:      220,
		SecondaryAutonomousY:      300,
		AutonomousEnergy:          100,
		SecondaryAutonomousEnergy: 100,
		AutonomousChildrenCount:   0,
		SecondaryChildrenCount:    0,
		DisableFoodSeeking:        true,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != simulation.DefaultAutonomousID {
		t.Fatalf("expected lower-id autonomous circle to win exact tie, got %q", snapshot.Interaction.WinnerID)
	}
}

func TestPlayerCanBeRemovedWhenLosingFight(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "triangle",
		PlayerEnergy:     50,
		AutonomousEnergy: 100,
	})

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.WinnerID != simulation.DefaultAutonomousID {
		t.Fatalf("expected autonomous circle to win, got %q", snapshot.Interaction.WinnerID)
	}
	if snapshot.Player != nil {
		t.Fatalf("expected player to be removed after losing fight, got %+v", snapshot.Player)
	}
}

func TestPlayerLoserWithChildAbsorbsFightLossAndRemainsActive(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            50,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     1,
		AutonomousChildrenCount: 0,
	})
	before := session.Snapshot()
	expectedChildID := before.Player.AttachedChildren[0].ID

	var snapshot simulation.Snapshot
	for range 20 {
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.Kind != "fight_absorbed_child" {
		t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.LoserID != "player-1" {
		t.Fatalf("expected player loser, got %q", snapshot.Interaction.LoserID)
	}
	if snapshot.Interaction.AbsorbedChildID != expectedChildID {
		t.Fatalf("expected absorbed child id %q, got %q", expectedChildID, snapshot.Interaction.AbsorbedChildID)
	}
	if snapshot.Player == nil {
		t.Fatal("expected player to remain active after child absorption")
	}
	if len(snapshot.Player.AttachedChildren) != 0 {
		t.Fatalf("expected absorbed conflict to consume one child, got %d", len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.Player.AttachedChildren) != 0 {
		t.Fatalf("expected visible attached child removal, got %d", len(snapshot.Player.AttachedChildren))
	}
	if snapshot.Player.Energy != before.Player.Energy {
		t.Fatalf("expected child absorption not to invent extra energy change, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if snapshot.Player.LineageID != before.Player.LineageID {
		t.Fatalf("expected lineage to stay %q, got %q", before.Player.LineageID, snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != before.Player.Generation {
		t.Fatalf("expected absorbed conflict to keep generation %d, got %d", before.Player.Generation, snapshot.Player.Generation)
	}
}

func TestResetRestoresInitialWorldState(t *testing.T) {
	session := simulation.NewSession()
	initial := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	for range 10 {
		_ = session.Advance()
	}

	reset := session.Reset()

	if reset.Tick != 0 {
		t.Fatalf("expected reset tick 0, got %d", reset.Tick)
	}
	if reset.Player == nil {
		t.Fatal("expected player after reset")
	}
	if reset.Player.X != initial.Player.X || reset.Player.Y != initial.Player.Y {
		t.Fatalf("expected player position reset to %+v, got %+v", initial.Player, reset.Player)
	}
	if reset.Player.Energy != initial.Player.Energy {
		t.Fatalf("expected player energy reset to %v, got %v", initial.Player.Energy, reset.Player.Energy)
	}
	if len(reset.Player.AttachedChildren) != len(initial.Player.AttachedChildren) {
		t.Fatalf("expected player child count reset to %d, got %d", len(initial.Player.AttachedChildren), len(reset.Player.AttachedChildren))
	}
	if len(reset.AutonomousCircles) != len(initial.AutonomousCircles) {
		t.Fatalf("expected %d autonomous circles after reset, got %d", len(initial.AutonomousCircles), len(reset.AutonomousCircles))
	}
	if len(reset.Foods) != len(initial.Foods) {
		t.Fatalf("expected %d foods after reset, got %d", len(initial.Foods), len(reset.Foods))
	}
	if reset.Interaction != nil {
		t.Fatalf("expected no interaction after reset, got %+v", reset.Interaction)
	}
}

func foodIDs(foods []simulation.Food) map[string]bool {
	ids := make(map[string]bool, len(foods))
	for _, food := range foods {
		ids[food.ID] = true
	}

	return ids
}

func autonomousByID(circles []simulation.AutonomousCircle, id string) (simulation.AutonomousCircle, bool) {
	for _, circle := range circles {
		if circle.ID == id {
			return circle, true
		}
	}

	return simulation.AutonomousCircle{}, false
}
