package integration_test

import (
	"context"
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
		PlayerShape:      "triangle",
		AutonomousShape:  "triangle",
		PlayerEnergy:     100,
		AutonomousEnergy: 80,
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
	if initial.Player.ChildrenCount != 0 || initial.AutonomousCircles[0].ChildrenCount != 0 || initial.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatal("expected child counts to start at zero")
	}
	if initial.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected base player radius %v, got %v", simulation.DefaultPlayerRadius, initial.Player.Radius)
	}
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
			if snapshot.Player.ChildrenCount != 1 || snapshot.AutonomousCircles[0].ChildrenCount != 1 {
				t.Fatalf("expected one child count for both circles after reproduction, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
			}
			expectedRadius := simulation.DefaultPlayerRadius + simulation.DefaultChildRadiusGain
			if snapshot.Player.Radius != expectedRadius || snapshot.AutonomousCircles[0].Radius != expectedRadius {
				t.Fatalf("expected grown radius %v after reproduction, player=%v autonomous=%v", expectedRadius, snapshot.Player.Radius, snapshot.AutonomousCircles[0].Radius)
			}
			continue
		}

		if resolvedSeen {
			if snapshot.Player.ChildrenCount > 1 || snapshot.AutonomousCircles[0].ChildrenCount > 1 {
				t.Fatalf("expected no repeat accumulation during continuous overlap, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
			}
			return
		}
	}

	t.Fatal("expected resolved reproduction snapshot")
}
