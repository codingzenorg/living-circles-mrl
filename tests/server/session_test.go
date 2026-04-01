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
