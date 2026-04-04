package integration_test

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func TestClientReceivesInitialSnapshotAndFightResolution(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:            "triangle",
		AutonomousShape:        "triangle",
		PlayerEnergy:           100,
		AutonomousEnergy:       80,
		DisableThreatAvoidance: true,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	if len(initial.Foods) == 0 {
		t.Fatal("expected initial snapshot to include food items")
	}
	if len(initial.AutonomousCircles) != 1 {
		t.Fatalf("expected one autonomous circle, got %d", len(initial.AutonomousCircles))
	}
	if initial.Player.Shape == "" || initial.AutonomousCircles[0].Shape == "" {
		t.Fatal("expected circle shapes in initial snapshot")
	}
	if initial.Player.Generation != 0 || initial.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected initial generation 0, player=%d autonomous=%d", initial.Player.Generation, initial.AutonomousCircles[0].Generation)
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": -1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read updated snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if snapshot.Interaction != nil {
			if snapshot.Player != nil && snapshot.Player.X >= initial.Player.X {
				t.Fatalf("expected player to move left toward interaction, before=%v after=%v", initial.Player.X, snapshot.Player.X)
			}
			if snapshot.Interaction.Kind != "fight_resolved" {
				t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
			}
			if snapshot.Interaction.WinnerID != "player-1" {
				t.Fatalf("expected player to win resolved fight, got %q", snapshot.Interaction.WinnerID)
			}
			if len(snapshot.AutonomousCircles) != 0 {
				t.Fatalf("expected autonomous circle to be removed after fight, got %d", len(snapshot.AutonomousCircles))
			}
			return
		}
	}
}

func TestClientReceivesFightAbsorptionThroughChildLoss(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        80,
		AutonomousChildrenCount: 1,
		DisableThreatAvoidance:  true,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "fight_absorbed_child" {
			t.Fatalf("expected fight_absorbed_child, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.LoserID != simulation.DefaultAutonomousID {
			t.Fatalf("expected autonomous loser, got %q", snapshot.Interaction.LoserID)
		}
		if len(snapshot.AutonomousCircles) != 1 {
			t.Fatalf("expected autonomous loser to remain active, got %d", len(snapshot.AutonomousCircles))
		}
		if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
			t.Fatalf("expected absorbed conflict to consume one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
		}
		if len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
			t.Fatalf("expected visible attached child removal, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
		}
		if snapshot.AutonomousCircles[0].LineageID != initial.AutonomousCircles[0].LineageID {
			t.Fatalf("expected lineage %q, got %q", initial.AutonomousCircles[0].LineageID, snapshot.AutonomousCircles[0].LineageID)
		}
		if snapshot.AutonomousCircles[0].Generation != initial.AutonomousCircles[0].Generation {
			t.Fatalf("expected generation to remain %d, got %d", initial.AutonomousCircles[0].Generation, snapshot.AutonomousCircles[0].Generation)
		}
		return
	}

	t.Fatal("expected fight absorption through child loss")
}

func TestClientReceivesFightOutcomeDrivenByChildPower(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        100,
		PlayerChildrenCount:     0,
		AutonomousChildrenCount: 1,
		DisableFoodSeeking:      true,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": -1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "fight_resolved" {
			t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.WinnerID != simulation.DefaultAutonomousID {
			t.Fatalf("expected higher-child autonomous circle to win, got %q", snapshot.Interaction.WinnerID)
		}
		return
	}

	t.Fatal("expected child-driven fight resolution snapshot")
}

func TestClientReceivesChildTriggeredReproductionBeforeParentBodiesOverlap(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Player == nil {
			t.Fatal("expected player to remain active after child-triggered reproduction")
		}
		if len(snapshot.AutonomousCircles) != 1 {
			t.Fatalf("expected one autonomous parent, got %d", len(snapshot.AutonomousCircles))
		}
		if snapshot.Interaction.ContactOrigin != "attached_child" {
			t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
		}
		return
	}

	t.Fatal("expected child-triggered reproduction snapshot")
}

func TestClientReceivesChildToChildTriggeredReproduction(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 10 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.ContactOrigin != "attached_child" {
			t.Fatalf("expected attached_child contact origin, got %q", snapshot.Interaction.ContactOrigin)
		}
		if snapshot.Player == nil || len(snapshot.AutonomousCircles) != 1 {
			t.Fatalf("expected both parents to remain, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
		}
		return
	}

	t.Fatal("expected child-to-child reproduction snapshot")
}

func TestClientReceivesAutonomousOnlyReproductionWithoutPlayerInput(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 10 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
			t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
		}
		if len(snapshot.AutonomousCircles) != 2 {
			t.Fatalf("expected both autonomous circles to remain, got %d", len(snapshot.AutonomousCircles))
		}
		if snapshot.Player == nil || snapshot.Player.X != initial.Player.X || snapshot.Player.Y != initial.Player.Y {
			t.Fatalf("expected player to remain uninvolved at %+v, got %+v", initial.Player, snapshot.Player)
		}
		return
	}

	t.Fatal("expected autonomous-only reproduction snapshot")
}

func TestClientReceivesAutonomousSeekingReproductionWithoutPlayerInput(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
			t.Fatalf("expected autonomous pair ids %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
		}
		if snapshot.Player == nil || snapshot.Player.X != initial.Player.X || snapshot.Player.Y != initial.Player.Y {
			t.Fatalf("expected player to remain uninvolved at %+v, got %+v", initial.Player, snapshot.Player)
		}
		return
	}

	t.Fatal("expected autonomous seeking reproduction snapshot")
}

func TestClientReceivesLowEnergyAutonomousFoodPriority(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction while low-energy autonomous seeks food, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].Y < initial.AutonomousCircles[0].Y {
			return
		}
	}

	t.Fatal("expected low-energy autonomous circle to steer toward food")
}

func TestClientReceivesShapeAwareAutonomousTargetChoice(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved from preferred different-shape target, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.SourceID != "player-1" || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
			t.Fatalf("expected different-shape target %q to be chosen ahead of the player, got %q -> %q", simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
		}
		return
	}

	t.Fatal("expected shape-aware autonomous target selection snapshot")
}

func TestClientReceivesFeasibilityAwareAutonomousFallback(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "fight_resolved" {
			t.Fatalf("expected infeasible reproduction target to be skipped in favor of fight_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.SourceID != simulation.DefaultAutonomousID || snapshot.Interaction.TargetID != simulation.DefaultSecondaryID {
			t.Fatalf("expected fallback to feasible target %q -> %q, got %q -> %q", simulation.DefaultAutonomousID, simulation.DefaultSecondaryID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
		}
		return
	}

	t.Fatal("expected feasibility-aware autonomous fallback snapshot")
}

func TestClientReceivesFightAwareAutonomousFoodFallback(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction while nearby stronger same-shape threat triggers retreat, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to retreat from nearby stronger same-shape threat")
}

func TestClientReceivesThreatAvoidanceAgainstPlayer(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction during threat avoidance, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to retreat from nearby stronger same-shape player threat")
}

func TestClientReceivesBlockedReproductionAvoidanceAgainstPlayer(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction during blocked-reproduction avoidance, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to retreat from nearby blocked different-shape player")
}

func TestClientReceivesChildTriggeredThreatAvoidanceBeforeParentOverlap(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	parentDistance := math.Hypot(initial.Player.X-initial.AutonomousCircles[0].X, initial.Player.Y-initial.AutonomousCircles[0].Y)
	if parentDistance < simulation.DefaultThreatAvoidanceDistance {
		t.Fatalf("expected parent body to stay outside threat window, got distance=%v", parentDistance)
	}
	childInsideWindow := false
	for _, child := range initial.Player.AttachedChildren {
		distance := math.Hypot(child.X-initial.AutonomousCircles[0].X, child.Y-initial.AutonomousCircles[0].Y)
		if distance < simulation.DefaultThreatAvoidanceDistance {
			childInsideWindow = true
			break
		}
	}
	if !childInsideWindow {
		t.Fatal("expected attached child to enter threat window before parent body")
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction during child-triggered threat avoidance, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to retreat from child-triggered same-shape threat")
}

func TestClientReceivesChildTriggeredBlockedReproductionAvoidanceBeforeParentOverlap(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	parentDistance := math.Hypot(initial.Player.X-initial.AutonomousCircles[0].X, initial.Player.Y-initial.AutonomousCircles[0].Y)
	if parentDistance < simulation.DefaultBlockedReproductionAvoidanceDistance {
		t.Fatalf("expected parent body to stay outside blocked-reproduction window, got distance=%v", parentDistance)
	}
	childInsideWindow := false
	for _, child := range initial.Player.AttachedChildren {
		distance := math.Hypot(child.X-initial.AutonomousCircles[0].X, child.Y-initial.AutonomousCircles[0].Y)
		if distance < simulation.DefaultBlockedReproductionAvoidanceDistance {
			childInsideWindow = true
			break
		}
	}
	if !childInsideWindow {
		t.Fatal("expected attached child to enter blocked-reproduction window before parent body")
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no immediate interaction during child-triggered blocked-reproduction avoidance, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to retreat from child-triggered blocked-reproduction target")
}

func TestClientReceivesDefaultDualInteractionDemoSnapshot(t *testing.T) {
	server := transport.NewServer()
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	if initial.Player == nil {
		t.Fatal("expected player in initial snapshot")
	}
	if len(initial.AutonomousCircles) != 2 {
		t.Fatalf("expected two autonomous circles, got %d", len(initial.AutonomousCircles))
	}
	if initial.AutonomousCircles[0].Shape != initial.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", initial.Player.Shape, initial.AutonomousCircles[0].Shape)
	}
	if initial.AutonomousCircles[1].Shape == initial.Player.Shape {
		t.Fatalf("expected second autonomous circle to differ from player shape %q", initial.Player.Shape)
	}
	if initial.Player.ChildrenCount != 1 {
		t.Fatalf("expected player child count to start at one for demo continuity, got %d", initial.Player.ChildrenCount)
	}
	if len(initial.Player.AttachedChildren) != initial.Player.ChildrenCount {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", initial.Player.ChildrenCount, len(initial.Player.AttachedChildren))
	}
	if initial.Player.LineageID != "lineage-player-1" || initial.Player.Generation != 0 {
		t.Fatalf("expected initial player lineage state, got lineage=%q generation=%d", initial.Player.LineageID, initial.Player.Generation)
	}
	if initial.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected first autonomous child count to start at one for demo continuity, got %d", initial.AutonomousCircles[0].ChildrenCount)
	}
	if len(initial.AutonomousCircles[0].AttachedChildren) != initial.AutonomousCircles[0].ChildrenCount {
		t.Fatalf("expected first autonomous attached children to match count, count=%d attached=%d", initial.AutonomousCircles[0].ChildrenCount, len(initial.AutonomousCircles[0].AttachedChildren))
	}
	if initial.AutonomousCircles[0].LineageID != "lineage-circle-2" || initial.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected first autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[0].LineageID, initial.AutonomousCircles[0].Generation)
	}
	if initial.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatalf("expected second autonomous child count to start at zero, got %d", initial.AutonomousCircles[1].ChildrenCount)
	}
	if len(initial.AutonomousCircles[1].AttachedChildren) != initial.AutonomousCircles[1].ChildrenCount {
		t.Fatalf("expected second autonomous attached children to match count, count=%d attached=%d", initial.AutonomousCircles[1].ChildrenCount, len(initial.AutonomousCircles[1].AttachedChildren))
	}
	if initial.AutonomousCircles[1].LineageID != "lineage-circle-3" || initial.AutonomousCircles[1].Generation != 0 {
		t.Fatalf("expected second autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[1].LineageID, initial.AutonomousCircles[1].Generation)
	}
	if initial.Player.Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected player demo radius %v, got %v", simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain, initial.Player.Radius)
	}
	if initial.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected first autonomous radius %v, got %v", simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain, initial.AutonomousCircles[0].Radius)
	}
}

func TestClientReceivesAutonomousFoodSeekingMotion(t *testing.T) {
	server := transport.NewServer()
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 10 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if snapshot.AutonomousCircles[1].Y > initial.AutonomousCircles[1].Y {
			return
		}
	}

	t.Fatal("expected second autonomous circle to steer downward toward food")
}

func TestClientReceivesResolvedReproductionWithoutRepeatAccumulation(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "square",
		PlayerEnergy:     100,
		AutonomousEnergy: 100,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	resolvedSeen := false
	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if snapshot.Interaction != nil && snapshot.Interaction.Kind == "reproduce_resolved" {
			resolvedSeen = true
			if snapshot.Player == nil {
				t.Fatal("expected player to remain active after reproduction")
			}
			if len(snapshot.AutonomousCircles) != len(initial.AutonomousCircles) {
				t.Fatalf("expected no new autonomous circles after reproduction, before=%d after=%d", len(initial.AutonomousCircles), len(snapshot.AutonomousCircles))
			}
			if snapshot.Player.Energy >= initial.Player.Energy {
				t.Fatalf("expected player energy to decrease after reproduction, before=%v after=%v", initial.Player.Energy, snapshot.Player.Energy)
			}
			if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != initial.Player.ChildrenCount+initial.AutonomousCircles[0].ChildrenCount+2 {
				t.Fatalf("expected reproduction to distribute two children across the pair, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
			}
			if snapshot.AutonomousCircles[0].Energy >= initial.AutonomousCircles[0].Energy {
				t.Fatalf("expected autonomous energy to decrease after reproduction, before=%v after=%v", initial.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
			}
			if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount || len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
				t.Fatalf("expected attached children to match counts after reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), snapshot.Player.ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren), snapshot.AutonomousCircles[0].ChildrenCount)
			}
			continue
		}

		if resolvedSeen {
			if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount > initial.Player.ChildrenCount+initial.AutonomousCircles[0].ChildrenCount+2 {
				t.Fatalf("expected no repeat accumulation during continuous overlap, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
			}
			if snapshot.Player.ChildrenCount > 0 && len(snapshot.Player.AttachedChildren) > 0 {
				if snapshot.Player.AttachedChildren[0].X == initial.Player.X && snapshot.Player.AttachedChildren[0].Y == initial.Player.Y {
					t.Fatal("expected attached child position to be authoritative world state, not parent center")
				}
			}
			return
		}
	}

	t.Fatal("expected resolved reproduction snapshot")
}

func TestClientReceivesBlockedReproductionWhenEnergyIsInsufficient(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		DisableBlockedReproductionAvoidance: true,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_blocked_energy" {
			t.Fatalf("expected reproduce_blocked_energy, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Player == nil {
			t.Fatal("expected player to remain active after blocked reproduction")
		}
		if snapshot.Player.ChildrenCount != 0 || snapshot.AutonomousCircles[0].ChildrenCount != 0 {
			t.Fatalf("expected blocked reproduction to preserve child counts, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
		}
		if len(snapshot.Player.AttachedChildren) != 0 || len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
			t.Fatalf("expected blocked reproduction to preserve attached children, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
		}
		return
	}

	t.Fatal("expected blocked reproduction snapshot")
}

func TestClientReceivesReproductionPaidByChildWhenEnergyIsLow(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		AutonomousChildrenCount:             1,
		DisableBlockedReproductionAvoidance: true,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	var previous simulation.Snapshot
	for range 40 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Interaction == nil {
			previous = snapshot
			continue
		}

		if snapshot.Interaction.Kind != "reproduce_resolved" {
			t.Fatalf("expected reproduce_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Player.ChildrenCount+snapshot.AutonomousCircles[0].ChildrenCount != previous.Player.ChildrenCount+previous.AutonomousCircles[0].ChildrenCount+1 {
			t.Fatalf("expected one child payment plus two redistributed children, before player=%d autonomous=%d after player=%d autonomous=%d", previous.Player.ChildrenCount, previous.AutonomousCircles[0].ChildrenCount, snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
		}
		expectedEnergy := previous.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
		if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
			t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
		}
		if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount || len(snapshot.AutonomousCircles[0].AttachedChildren) != snapshot.AutonomousCircles[0].ChildrenCount {
			t.Fatalf("expected attached children to match counts after child-payment reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), snapshot.Player.ChildrenCount, len(snapshot.AutonomousCircles[0].AttachedChildren), snapshot.AutonomousCircles[0].ChildrenCount)
		}
		return
	}

	t.Fatal("expected reproduction snapshot paid by child")
}

func TestResetEndpointReturnsAndBroadcastsInitialSnapshot(t *testing.T) {
	server := transport.NewServer()
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": 1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	var advanced simulation.Snapshot
	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))
	for range 20 {
		if err := connection.ReadJSON(&advanced); err != nil {
			t.Fatalf("read advanced snapshot: %v", err)
		}
		if advanced.Tick > 0 {
			break
		}
	}

	if advanced.Tick == 0 {
		t.Fatal("expected advanced snapshot before reset")
	}

	response, err := http.Post(httpServer.URL+"/reset", "application/json", nil)
	if err != nil {
		t.Fatalf("post reset: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected reset status 200, got %d", response.StatusCode)
	}

	var resetResponse simulation.Snapshot
	if err := json.NewDecoder(response.Body).Decode(&resetResponse); err != nil {
		t.Fatalf("decode reset response: %v", err)
	}

	if resetResponse.Tick != 0 {
		t.Fatalf("expected reset response tick 0, got %d", resetResponse.Tick)
	}
	if resetResponse.Player == nil || resetResponse.Player.X != initial.Player.X || resetResponse.Player.Energy != initial.Player.Energy {
		t.Fatalf("expected reset response player to match initial snapshot, initial=%+v reset=%+v", initial.Player, resetResponse.Player)
	}

	var broadcast simulation.Snapshot
	if err := connection.ReadJSON(&broadcast); err != nil {
		t.Fatalf("read reset broadcast: %v", err)
	}

	if broadcast.Tick != 0 {
		t.Fatalf("expected reset broadcast tick 0, got %d", broadcast.Tick)
	}
	if broadcast.Player == nil || broadcast.Player.X != initial.Player.X || broadcast.Player.Energy != initial.Player.Energy {
		t.Fatalf("expected reset broadcast player to match initial snapshot, initial=%+v broadcast=%+v", initial.Player, broadcast.Player)
	}
}

func TestClientReceivesEnergyCollapseReplacement(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:         "triangle",
		AutonomousShape:     "square",
		PlayerEnergy:        1,
		AutonomousEnergy:    100,
		PlayerChildrenCount: 1,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": 1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Player == nil {
			t.Fatal("expected replacement player to remain active after zero-energy collapse")
		}
		if snapshot.Player.ChildrenCount != 0 {
			t.Fatalf("expected replacement to consume one child, got %d", snapshot.Player.ChildrenCount)
		}
		if len(snapshot.Player.AttachedChildren) != snapshot.Player.ChildrenCount {
			t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", snapshot.Player.ChildrenCount, len(snapshot.Player.AttachedChildren))
		}
		if snapshot.Player.Energy != simulation.DefaultReplacementEnergy {
			t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.Player.Energy)
		}
		if snapshot.Player.LineageID != initial.Player.LineageID {
			t.Fatalf("expected replacement lineage %q, got %q", initial.Player.LineageID, snapshot.Player.LineageID)
		}
		if snapshot.Player.Generation != initial.Player.Generation+1 {
			t.Fatalf("expected replacement generation %d, got %d", initial.Player.Generation+1, snapshot.Player.Generation)
		}
		if snapshot.Interaction == nil {
			t.Fatal("expected continuity interaction after zero-energy promotion")
		}
		if snapshot.Interaction.Kind != "death_promoted_child" {
			t.Fatalf("expected death_promoted_child, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.SourceID != initial.Player.ID || snapshot.Interaction.TargetID != initial.Player.ID {
			t.Fatalf("expected continuity interaction to identify player %q, got source=%q target=%q", initial.Player.ID, snapshot.Interaction.SourceID, snapshot.Interaction.TargetID)
		}
		return
	}

	t.Fatal("expected energy-collapse replacement snapshot")
}

func TestClientReceivesRegeneratedFoodAfterDeterministicDelay(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              simulation.DefaultPlayerEnergy,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": 1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(4 * time.Second))

	sawFoodMissing := false
	for range simulation.DefaultFoodRegenDelay + 8 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		_, found := foodByID(snapshot.Foods, "food-1")
		if !found {
			sawFoodMissing = true
			continue
		}

		if !sawFoodMissing {
			continue
		}

		if len(snapshot.Foods) != len(initial.Foods) {
			t.Fatalf("expected full food set after regeneration, got %d foods", len(snapshot.Foods))
		}
		food, _ := foodByID(snapshot.Foods, "food-1")
		if food.X != 432 || food.Y != 300 {
			t.Fatalf("expected regenerated food-1 to return to its original slot, got %+v", food)
		}
		return
	}

	t.Fatal("expected regenerated food after deterministic delay")
}

func TestClientReceivesFoodCollectionThroughAttachedChild(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerEnergy:              80,
		PlayerChildrenCount:       2,
		AutonomousEnergy:          0,
		SecondaryAutonomousEnergy: 0,
	}))
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Run(ctx)
	}()

	websocketURL := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer connection.Close()

	var initial simulation.Snapshot
	if err := connection.ReadJSON(&initial); err != nil {
		t.Fatalf("read initial snapshot: %v", err)
	}

	message := map[string]any{
		"type": "movement_intent",
		"direction": map[string]float64{
			"x": 1,
			"y": 0,
		},
	}
	if err := connection.WriteJSON(message); err != nil {
		t.Fatalf("write movement intent: %v", err)
	}

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 8 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if len(snapshot.Foods) != len(initial.Foods)-1 {
			continue
		}

		if snapshot.Player == nil {
			t.Fatal("expected player to remain active after attached-child collection")
		}
		expectedEnergy := initial.Player.Energy - simulation.DefaultMoveCost + simulation.DefaultFoodEnergy
		if snapshot.Player.Energy != expectedEnergy {
			t.Fatalf("expected player energy %v after attached-child collection, got %v", expectedEnergy, snapshot.Player.Energy)
		}
		if len(snapshot.Player.AttachedChildren) != 2 {
			t.Fatalf("expected attached children to remain visible after collection, got %d", len(snapshot.Player.AttachedChildren))
		}
		return
	}

	t.Fatal("expected food collection through attached child")
}

func foodByID(foods []simulation.Food, id string) (simulation.Food, bool) {
	for _, food := range foods {
		if food.ID == id {
			return food, true
		}
	}

	return simulation.Food{}, false
}
