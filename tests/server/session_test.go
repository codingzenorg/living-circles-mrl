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

func TestAdvanceKeepsPlayerInsideBounds(t *testing.T) {
	session := simulation.NewSession()

	session.ApplyIntent(simulation.Vector{X: 1, Y: 0})
	var snapshot simulation.Snapshot
	for range 200 {
		snapshot = session.Advance()
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

	if snapshot.AutonomousCircles[0].Shape != snapshot.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", snapshot.Player.Shape, snapshot.AutonomousCircles[0].Shape)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
		t.Fatalf("expected first autonomous children count to start at zero, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
	}

	if snapshot.AutonomousCircles[1].Shape != simulation.DefaultAutoShape {
		t.Fatalf("expected second autonomous shape %q, got %q", simulation.DefaultAutoShape, snapshot.AutonomousCircles[1].Shape)
	}
	if snapshot.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatalf("expected second autonomous children count to start at zero, got %d", snapshot.AutonomousCircles[1].ChildrenCount)
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
}

func TestDefaultWorldSupportsDifferentShapeReproductionPath(t *testing.T) {
	session := simulation.NewSession()

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
	if snapshot.AutonomousCircles[1].ChildrenCount != 1 {
		t.Fatalf("expected autonomous child accumulation after reproduction, got %d", snapshot.AutonomousCircles[1].ChildrenCount)
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

func TestDifferentShapeOverlapProducesResolvedReproduction(t *testing.T) {
	session := simulation.NewSessionWithShapes("triangle", "square")

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
	if snapshot.Player.ChildrenCount != 1 {
		t.Fatalf("expected player child accumulation, got %d", snapshot.Player.ChildrenCount)
	}
	if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected autonomous child accumulation, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
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
	session := simulation.NewSessionWithShapes("triangle", "square")

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

	for range 120 {
		session.ApplyIntent(simulation.Vector{})
		snapshot = session.Advance()
		if snapshot.Interaction != nil {
			break
		}
	}

	if snapshot.Interaction == nil {
		t.Fatal("expected second reproduction after re-overlap")
	}
	if snapshot.Interaction.Kind != "reproduce_resolved" {
		t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
	}
	if snapshot.Player.ChildrenCount != 2 || snapshot.AutonomousCircles[0].ChildrenCount != 2 {
		t.Fatalf("expected second reproduction to increment child counts, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
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
