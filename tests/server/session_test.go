package server_test

import (
	"math"
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
	if snapshot.Interaction == nil {
		t.Fatal("expected continuity interaction after zero-energy promotion")
	}
	if snapshot.Interaction.Kind != "death_promoted_child" {
		t.Fatalf("expected death_promoted_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != before.Player.ID || snapshot.Interaction.TargetID != before.Player.ID {
		t.Fatalf("expected continuity interaction to identify player %q, got source=%q target=%q", before.Player.ID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount {
		t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", snapshot.Player.ChildrenCount, len(snapshot.Player.AttachedChildren))
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
	if snapshot.Interaction == nil {
		t.Fatal("expected continuity interaction after zero-energy promotion")
	}
	if snapshot.Interaction.Kind != "death_promoted_child" {
		t.Fatalf("expected death_promoted_child, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Interaction.SourceID != before.AutonomousCircles[0].ID || snapshot.Interaction.TargetID != before.AutonomousCircles[0].ID {
		t.Fatalf("expected continuity interaction to identify autonomous circle %q, got source=%q target=%q", before.AutonomousCircles[0].ID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
		t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", snapshot.AutonomousCircles[0].ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren))
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
	if snapshot.Player.ChildrenCount != 1 {
		t.Fatalf("expected player children count to start at one for demo continuity, got %d", snapshot.Player.ChildrenCount)
	}
	if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", snapshot.Player.ChildrenCount, len(snapshot.Player.AttachedChildren))
	}
	if snapshot.Player.LineageID != "lineage-player-1" {
		t.Fatalf("expected player lineage %q, got %q", "lineage-player-1", snapshot.Player.LineageID)
	}
	if snapshot.Player.Generation != 0 {
		t.Fatalf("expected player generation 0, got %d", snapshot.Player.Generation)
	}
	if snapshot.Player.Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected player demo radius %v, got %v", simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain, snapshot.Player.Radius)
	}

	if snapshot.AutonomousCircles[0].Shape != snapshot.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", snapshot.Player.Shape, snapshot.AutonomousCircles[0].Shape)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected first autonomous children count to start at one for demo continuity, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
		t.Fatalf("expected first autonomous attached children to match count, count=%d attached=%d", snapshot.AutonomousCircles[0].ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren))
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
	if len(snapshot.AutonomousCircles[1].AttachedChildren) != snapshot.AutonomousCircles[1].ChildrenCount {
		t.Fatalf("expected second autonomous attached children to match count, count=%d attached=%d", snapshot.AutonomousCircles[1].ChildrenCount, len(snapshot.AutonomousCircles[1].AttachedChildren))
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

func TestChildGrowthAloneDoesNotImproveFoodCollectionReach(t *testing.T) {
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
	if len(after.Foods) != 3 {
		t.Fatalf("expected enlarged derived radius alone to miss food on first move, got %d foods", len(after.Foods))
	}
	expectedEnergy := simulation.DefaultPlayerEnergy - simulation.DefaultMoveCost
	if after.Player.Energy != expectedEnergy {
		t.Fatalf("expected no immediate food recovery from derived radius alone, got %v", after.Player.Energy)
	}
}

func TestAttachedChildCanCollectFoodForPlayer(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:              "triangle",
		AutonomousShape:          "square",
		PlayerX:                  389,
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
	if len(snapshot.AutonomousCircles) != 2 {
		t.Fatalf("expected both autonomous circles to remain, got %d circles", len(snapshot.AutonomousCircles))
	}
	foundAbsorbedLoser := false
	for _, circle := range snapshot.AutonomousCircles {
		if circle.ID != simulation.DefaultAutonomousID {
			continue
		}
		foundAbsorbedLoser = true
		if circle.ChildrenCount != 0 {
			t.Fatalf("expected absorbed conflict to remove one child, got %d", circle.ChildrenCount)
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
	if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[1].ChildrenCount != before.Player.ChildrenCount+before.AutonomousCircles[1].ChildrenCount+2 {
		t.Fatalf("expected reproduction to distribute two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", before.Player.ChildrenCount, before.AutonomousCircles[1].ChildrenCount, snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[1].ChildrenCount)
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through movement and reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount {
		t.Fatalf("expected player attached children to match count after reproduction, count=%d attached=%d", snapshot.Player.ChildrenCount, len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[1].AttachedChildren) != snapshot.AutonomousCircles[1].ChildrenCount {
		t.Fatalf("expected autonomous attached children to match count after reproduction, count=%d attached=%d", snapshot.AutonomousCircles[1].ChildrenCount, len(snapshot.AutonomousCircles[1].AttachedChildren))
	}
	if snapshot.AutonomousCircles[1].Energy >= before.AutonomousCircles[1].Energy {
		t.Fatalf("expected autonomous energy to decrease through movement and reproduction, before=%v after=%v", before.AutonomousCircles[1].Energy, snapshot.AutonomousCircles[1].Energy)
	}
	if snapshot.Player.ChildrenCount > 0 && snapshot.Player.Radius <= before.Player.Radius {
		t.Fatalf("expected player radius growth when player owns children, before=%v after=%v", before.Player.Radius, snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[1].ChildrenCount > before.AutonomousCircles[1].ChildrenCount && snapshot.AutonomousCircles[1].Radius <= before.AutonomousCircles[1].Radius {
		t.Fatalf("expected autonomous radius growth when autonomous owns children, before=%v after=%v", before.AutonomousCircles[1].Radius, snapshot.AutonomousCircles[1].Radius)
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
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected autonomous loser to remain active, got %d circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected absorbed conflict to consume one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
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
}

func TestAttachedChildrenCanTriggerFightThroughChildToChildContact(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "triangle",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               254.511,
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
	if snapshot.Player == nil || len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected both parents to remain after absorbed fight, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected autonomous loser to lose one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}
}

func TestAttachedChildrenCanTriggerReproductionThroughChildToChildContact(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   200,
		PlayerY:                   300,
		AutonomousX:               254.511,
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
	if snapshot.Player == nil || len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected both parents to remain after reproduction, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
	}
	if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != before.Player.ChildrenCount+before.AutonomousCircles[0].ChildrenCount+2 {
		t.Fatalf("expected reproduction to add two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", before.Player.ChildrenCount, before.AutonomousCircles[0].ChildrenCount, snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
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
	if left.ChildrenCount+right.ChildrenCount != beforeLeft.ChildrenCount+beforeRight.ChildrenCount+2 {
		t.Fatalf("expected reproduction to add two children across the autonomous pair, before left=%d right=%d after left=%d right=%d", beforeLeft.ChildrenCount, beforeRight.ChildrenCount, left.ChildrenCount, right.ChildrenCount)
	}
	if snapshot.Player == nil || snapshot.Player.X != before.Player.X || snapshot.Player.Y != before.Player.Y {
		t.Fatalf("expected player to remain uninvolved at %+v, got %+v", before.Player, snapshot.Player)
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
	if len(snapshot.AutonomousCircles) != len(before.AutonomousCircles) {
		t.Fatalf("expected no new autonomous circles, before=%d after=%d", len(before.AutonomousCircles), len(snapshot.AutonomousCircles))
	}
	if snapshot.Player.Energy >= before.Player.Energy {
		t.Fatalf("expected player energy to decrease through reproduction, before=%v after=%v", before.Player.Energy, snapshot.Player.Energy)
	}
	if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != before.Player.ChildrenCount+before.AutonomousCircles[0].ChildrenCount+2 {
		t.Fatalf("expected reproduction to distribute two children across the pair, before player=%d autonomous=%d after player=%d autonomous=%d", before.Player.ChildrenCount, before.AutonomousCircles[0].ChildrenCount, snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", snapshot.Player.ChildrenCount, len(snapshot.Player.AttachedChildren))
	}
	if len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
		t.Fatalf("expected autonomous attached children to match count, count=%d attached=%d", snapshot.AutonomousCircles[0].ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
	if snapshot.AutonomousCircles[0].Energy >= before.AutonomousCircles[0].Energy {
		t.Fatalf("expected autonomous energy to decrease through reproduction, before=%v after=%v", before.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
	}
	if snapshot.Player.ChildrenCount > 0 && snapshot.Player.Radius <= before.Player.Radius {
		t.Fatalf("expected player radius growth when player owns children, before=%v after=%v", before.Player.Radius, snapshot.Player.Radius)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount > before.AutonomousCircles[0].ChildrenCount && snapshot.AutonomousCircles[0].Radius <= before.AutonomousCircles[0].Radius {
		t.Fatalf("expected autonomous radius growth when autonomous owns children, before=%v after=%v", before.AutonomousCircles[0].Radius, snapshot.AutonomousCircles[0].Radius)
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
	if len(snapshot.AutonomousCircles) != 1 {
		t.Fatalf("expected blocked reproduction to avoid spawning autonomous circles, got %d autonomous circles", len(snapshot.AutonomousCircles))
	}
	if snapshot.Player.ChildrenCount != 0 || snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected blocked reproduction to preserve child counts, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
	}
	if len(snapshot.Player.AttachedChildren) != 0 || len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
		t.Fatalf("expected blocked reproduction to preserve attached children, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
	}
}

func TestDifferentShapeOverlapConsumesChildAsReproductionPayment(t *testing.T) {
	session := simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerX:                             200,
		PlayerY:                             300,
		AutonomousX:                         232,
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
	if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != before.Player.ChildrenCount+before.AutonomousCircles[0].ChildrenCount+1 {
		t.Fatalf("expected one child to be spent and two children to be redistributed, before player=%d autonomous=%d after player=%d autonomous=%d", before.Player.ChildrenCount, before.AutonomousCircles[0].ChildrenCount, snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
	}
	expectedEnergy := before.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
	if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
		t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
	}
	if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount || len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
		t.Fatalf("expected attached children to match counts after child-payment reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), snapshot.Player.ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren), snapshot.AutonomousCircles[0].ChildrenCount)
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
	if snapshot.Player.ChildrenCount+originalOpponent.ChildrenCount != 2 {
		t.Fatalf("expected first reproduction to create two visible child units across the pair, player=%d autonomous=%d", snapshot.Player.ChildrenCount, originalOpponent.ChildrenCount)
	}

	for range 3 {
		snapshot = session.Advance()
		originalOpponent, found = autonomousByID(snapshot.AutonomousCircles, originalOpponentID)
		if !found {
			t.Fatalf("expected original autonomous circle %q during continuous overlap", originalOpponentID)
		}
		if snapshot.Player.ChildrenCount+originalOpponent.ChildrenCount != 2 {
			t.Fatalf("expected continuous overlap to avoid repeat accumulation, player=%d autonomous=%d", snapshot.Player.ChildrenCount, originalOpponent.ChildrenCount)
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
		if snapshot.Interaction != nil && snapshot.Interaction.Kind == "reproduce_resolved" {
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
	if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != beforePlayerChildren+beforeAutonomousChildren+2 {
		t.Fatalf(
			"expected second reproduction to add two child units across the pair, before player=%d autonomous=%d after player=%d autonomous=%d",
			beforePlayerChildren,
			beforeAutonomousChildren,
			snapshot.Player.ChildrenCount,
			snapshot.AutonomousCircles[0].ChildrenCount,
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
		t.Fatalf("expected higher-child autonomous circle to win, got %q", snapshot.Interaction.WinnerID)
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
	if snapshot.Player == nil {
		t.Fatal("expected player to remain active after child absorption")
	}
	if snapshot.Player.ChildrenCount != 0 {
		t.Fatalf("expected absorbed conflict to consume one child, got %d", snapshot.Player.ChildrenCount)
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

func autonomousByID(circles []simulation.AutonomousCircle, id string) (simulation.AutonomousCircle, bool) {
	for _, circle := range circles {
		if circle.ID == id {
			return circle, true
		}
	}

	return simulation.AutonomousCircle{}, false
}
