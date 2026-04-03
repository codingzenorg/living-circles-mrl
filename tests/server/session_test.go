package server_test

import (
	"testing"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

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
		PlayerEnergy:        1,
		AutonomousEnergy:    100,
		PlayerChildrenCount: 1,
	})
	before := session.Snapshot()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	snapshot := session.Advance()

	if snapshot.Player == nil {
		t.Fatal("expected player replacement after zero-energy collapse")
	}
	if snapshot.Player.ChildrenCount != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", snapshot.Player.ChildrenCount)
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
	})
	before := session.Snapshot()
	snapshot := session.Advance()

	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous replacement after zero-energy collapse, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
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

	if snapshot.Player.X > snapshot.World.Width-snapshot.Player.Radius {
		t.Fatalf("expected x to stay inside bounds, got %v", snapshot.Player.X)
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

	if len(snapshot.Foods) != 3 {
		t.Fatalf("expected three deterministic food items, got %d", len(snapshot.Foods))
	}

	if snapshot.Foods[0].ID != "food-1" || snapshot.Foods[0].X != 432 || snapshot.Foods[0].Y != 300 {
		t.Fatalf("unexpected first food placement: %+v", snapshot.Foods[0])
	}

	if snapshot.Foods[1].ID != "food-2" || snapshot.Foods[1].X != 292 || snapshot.Foods[1].Y != 300 {
		t.Fatalf("unexpected second food placement: %+v", snapshot.Foods[1])
	}

	if len(snapshot.AutonomousCircles) != 2 {
		t.Fatalf("expected two autonomous circles, got %d", len(snapshot.AutonomousCircles))
	}

	if snapshot.Player.Shape != simulation.DefaultPlayerShape {
		t.Fatalf("expected player shape %q, got %q", simulation.DefaultPlayerShape, snapshot.Player.Shape)
	}
	if snapshot.Player.ChildrenCount != 0 {
		t.Fatalf("expected player children count to start at zero, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.Player.LineageID != "lineage-player-1" {
		t.Fatalf("expected player lineage %q, got %q", "lineage-player-1", snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != 0 {
		t.Fatalf("expected player generation 0, got %d", snapshot.Player.Generation)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected base player radius %v, got %v", simulation.DefaultPlayerRadius, snapshot.Player.Radius)
	}

	if snapshot.AutonomousCircles[0].Shape != snapshot.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", snapshot.Player.Shape, snapshot.AutonomousCircles[0].Shape)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected first autonomous children count to start at one for demo continuity, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if snapshot.AutonomousCircles[0].LineageID != "lineage-circle-2" {
		t.Fatalf("expected first autonomous lineage %q, got %q", "lineage-circle-2", snapshot.AutonomousCircles[0].LineageID)
	}
	if snapshot.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected first autonomous generation 0, got %d", snapshot.AutonomousCircles[0].Generation)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected first autonomous radius %v, got %v", simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain, snapshot.AutonomousCircles[0].Radius)
	}

	if snapshot.AutonomousCircles[1].Shape != simulation.DefaultAutoShape {
		t.Fatalf("expected second autonomous shape %q, got %q", simulation.DefaultAutoShape, snapshot.AutonomousCircles[1].Shape)
	}
	if snapshot.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatalf("expected second autonomous children count to start at zero, got %d", snapshot.AutonomousCircles[1].ChildrenCount)
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
}

func TestOverlappingFoodRemovesItAndRestoresEnergy(t *testing.T) {
	session := simulation.NewSession()
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

func TestInitialChildrenIncreaseRadiusDeterministically(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "square",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     2,
		AutonomousChildrenCount: 1,
	})
	snapshot := session.Snapshot()

	expectedPlayerRadius := simulation.DefaultPlayerRadius + 2*simulation.DefaultChildRadiusGain
	expectedAutonomousRadius := simulation.DefaultPlayerRadius + simulation.DefaultChildRadiusGain

	if snapshot.Player.Radius != expectedPlayerRadius {
		t.Fatalf("expected player radius %v, got %v", expectedPlayerRadius, snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[0].Radius != expectedAutonomousRadius {
		t.Fatalf("expected autonomous radius %v, got %v", expectedAutonomousRadius, snapshot.AutonomousCircles[0].Radius)
	}
}

func TestFoodRecoveryRespectsEnergyCap(t *testing.T) {
	session := simulation.NewSession()
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

func TestChildGrowthImprovesFoodCollectionReach(t *testing.T) {
	baseSession := simulation.NewSession()
	baseSession.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	baseAfter := baseSession.Advance()

	grownSession := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "square",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     2,
		AutonomousChildrenCount: 0,
	})
	grownSession.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	after := grownSession.Advance()

	if len(baseAfter.Foods) != 3 {
		t.Fatalf("expected base-radius player to miss food on first move, got %d foods", len(baseAfter.Foods))
	}
	if len(after.Foods) != 2 {
		t.Fatalf("expected larger-radius player to consume food on first move, got %d foods", len(after.Foods))
	}
	if after.Player.Energy != simulation.DefaultMaxEnergy {
		t.Fatalf("expected immediate food recovery from grown reach, got %v", after.Player.Energy)
	}
}

func TestAutonomousCircleMovesWithoutPlayerInput(t *testing.T) {
	session := simulation.NewSession()
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
	session := simulation.NewSession()
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
	session := simulation.NewSession()
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
	session := simulation.NewSession()
	firstTick := session.Advance()
	secondTick := session.Advance()

	if len(secondTick.Foods) >= len(firstTick.Foods) {
		t.Fatalf("expected autonomous food consumption, before=%d after=%d", len(firstTick.Foods), len(secondTick.Foods))
	}

	if secondTick.AutonomousCircles[0].Energy <= firstTick.AutonomousCircles[0].Energy {
		t.Fatalf("expected autonomous circle energy recovery, before=%v after=%v", firstTick.AutonomousCircles[0].Energy, secondTick.AutonomousCircles[0].Energy)
	}
}

func TestAutonomousFoodSeekingCanCollectOffLaneFood(t *testing.T) {
	session := simulation.NewSession()

	var snapshot simulation.Snapshot
	collected := false
	for range 20 {
		snapshot = session.Advance()
		if len(snapshot.Foods) < 3 {
			collected = true
			break
		}
	}

	if !collected {
		t.Fatal("expected autonomous food seeking to collect food within the first steering window")
	}
	if snapshot.AutonomousCircles[1].Y <= 300 {
		t.Fatalf("expected second autonomous circle to steer off its original horizontal lane, got y=%v", snapshot.AutonomousCircles[1].Y)
	}
}

func TestDefaultWorldSupportsSameShapeFightPath(t *testing.T) {
	session := simulation.NewSession()

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
	if snapshot.Interaction.Kind != "fight_resolved" {
		t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.TargetID != simulation.DefaultAutonomousID {
		t.Fatalf("expected fight against %q, got %q", simulation.DefaultAutonomousID, snapshot.Interaction.TargetID)
	}
	if len(snapshot.AutonomousCircles) != 2 {
		t.Fatalf("expected replacement plus different-shape autonomous circle to remain, got %d circles", len(snapshot.AutonomousCircles))
	}
	foundReplacement := false
	for _, circle := range snapshot.AutonomousCircles {
		if circle.ID != simulation.DefaultAutonomousID {
			continue
		}
		foundReplacement = true
		if circle.ChildrenCount != 0 {
			t.Fatalf("expected replacement to consume one child, got %d", circle.ChildrenCount)
		}
		break
	}
	if !foundReplacement {
		t.Fatalf("expected replacement circle %q to remain active", simulation.DefaultAutonomousID)
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
	if snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
		t.Fatalf("expected reproduction against %q, got %q", simulation.DefaultSecondaryID, snapshot.Interaction.TargetID)
	}
	if snapshot.Player.ChildrenCount != 1 {
		t.Fatalf("expected player child accumulation after reproduction, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through movement and reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected player radius growth after reproduction, got %v", snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[1].ChildrenCount != 1 {
		t.Fatalf("expected autonomous child accumulation after reproduction, got %d", snapshot.AutonomousCircles[1].ChildrenCount)
	}
	if snapshot.AutonomousCircles[1].Energy >= before.AutonomousCircles[1].Energy {
		t.Fatalf("expected autonomous energy to decrease through movement and reproduction, before=%v after=%v", before.AutonomousCircles[1].Energy, snapshot.AutonomousCircles[1].Energy)
	}
	if snapshot.AutonomousCircles[1].Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected autonomous radius growth after reproduction, got %v", snapshot.AutonomousCircles[1].Radius)
	}
}

func TestSameShapeOverlapProducesFightCandidate(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "triangle")

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

func TestAutonomousLoserWithChildRemainsActiveThroughReplacement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        80,
		AutonomousChildrenCount: 1,
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
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.LoserID != simulation.DefaultAutonomousID {
		t.Fatalf("expected autonomous loser, got %q", snapshot.Interaction.LoserID)
	}
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous replacement to remain active, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if snapshot.AutonomousCircles[0].Energy != simulation.DefaultReplacementEnergy {
		t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected replacement radius to reset to base after child consumption, got %v", snapshot.AutonomousCircles[0].Radius)
	}
	if snapshot.AutonomousCircles[0].LineageID != before.AutonomousCircles[0].LineageID {
		t.Fatalf("expected replacement lineage %q, got %q", before.AutonomousCircles[0].LineageID, snapshot.AutonomousCircles[0].LineageID)
	}
	if snapshot.AutonomousCircles[0].Generation != before.AutonomousCircles[0].Generation+1 {
		t.Fatalf("expected replacement generation %d, got %d", before.AutonomousCircles[0].Generation+1, snapshot.AutonomousCircles[0].Generation)
	}
}

func TestDifferentShapeOverlapProducesResolvedReproduction(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "square")
	before := session.Snapshot()

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
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if snapshot.Player.ChildrenCount != 1 {
		t.Fatalf("expected player child accumulation, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected player radius growth, got %v", snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected autonomous child accumulation, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if snapshot.AutonomousCircles[0].Energy >= before.AutonomousCircles[0].Energy {
		t.Fatalf("expected autonomous energy to decrease through reproduction, before=%v after=%v", before.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected autonomous radius growth, got %v", snapshot.AutonomousCircles[0].Radius)
	}
}

func TestDifferentShapeOverlapBlocksReproductionWhenEnergyIsInsufficient(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "square",
		PlayerEnergy:     simulation.DefaultPlayerEnergy,
		AutonomousEnergy: simulation.DefaultReproductionCost - 1,
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
	if snapshot.Player.ChildrenCount != 0 || snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected blocked reproduction to preserve child counts, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
	}
}

func TestDifferentShapeOverlapConsumesChildAsReproductionPayment(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "square",
		PlayerEnergy:            simulation.DefaultPlayerEnergy,
		AutonomousEnergy:        simulation.DefaultReproductionCost - 1,
		AutonomousChildrenCount: 1,
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
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player.ChildrenCount != 1 {
		t.Fatalf("expected player to gain one child, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected autonomous circle to spend one child and gain one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
	expectedEnergy := before.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
	if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
		t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
	}
}

func TestContinuousOverlapDoesNotRepeatChildAccumulation(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "square")

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
	if snapshot.Player.ChildrenCount != 1 || snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected one child each after first reproduction, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
	}

	for range 3 {
		snapshot = session.Advance()
		if snapshot.Player.ChildrenCount != 1 || snapshot.AutonomousCircles[0].ChildrenCount != 1 {
			t.Fatalf("expected continuous overlap to avoid repeat accumulation, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
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

	beforePlayerChildren := snapshot.Player.ChildrenCount
	beforeAutonomousChildren := snapshot.AutonomousCircles[0].ChildrenCount
	for range 120 {
		session.ApplyIntent(simulation.Vector{})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
		beforePlayerChildren = snapshot.Player.ChildrenCount
		beforeAutonomousChildren = snapshot.AutonomousCircles[0].ChildrenCount
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected second reproduction after re-overlap")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player.ChildrenCount != beforePlayerChildren+1 || snapshot.AutonomousCircles[0].ChildrenCount != beforeAutonomousChildren+1 {
		t.Fatalf(
			"expected second reproduction to increment child counts from current state, before player=%d autonomous=%d after player=%d autonomous=%d",
			beforePlayerChildren,
			beforeAutonomousChildren,
			snapshot.Player.ChildrenCount,
			snapshot.AutonomousCircles[0].ChildrenCount,
		)
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
		PlayerShape:      "triangle",
		AutonomousShape:  "triangle",
		PlayerEnergy:     100,
		AutonomousEnergy: 80,
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

func TestLargerRadiusBreaksEqualEnergyFightTie(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     1,
		AutonomousChildrenCount: 0,
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
		t.Fatalf("expected larger-radius player to win tie-break fight, got %q", snapshot.Interaction.WinnerID)
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
		t.Fatalf("expected higher-child autonomous circle to win, got %q", snapshot.Interaction.WinnerID)
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

func TestPlayerLoserWithChildRemainsActiveThroughReplacement(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            50,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     1,
		AutonomousChildrenCount: 0,
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
		t.Fatal("expected fight resolution")
	}
	if snapshot.Interaction.LoserID != "player-1" {
		t.Fatalf("expected player loser, got %q", snapshot.Interaction.LoserID)
	}
	if snapshot.Player == nil {
		t.Fatal("expected player replacement to remain active")
	}
	if snapshot.Player.ChildrenCount != 0 {
		t.Fatalf("expected replacement to consume one child, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.Player.Energy != simulation.DefaultReplacementEnergy {
		t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.Player.Energy)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected replacement radius to reset to base after child consumption, got %v", snapshot.Player.Radius)
	}
	if snapshot.Player.LineageID != before.Player.LineageID {
		t.Fatalf("expected replacement lineage %q, got %q", before.Player.LineageID, snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != before.Player.Generation+1 {
		t.Fatalf("expected replacement generation %d, got %d", before.Player.Generation+1, snapshot.Player.Generation)
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
	if reset.Player.ChildrenCount != initial.Player.ChildrenCount {
		t.Fatalf("expected player child count reset to %d, got %d", initial.Player.ChildrenCount, reset.Player.ChildrenCount)
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
