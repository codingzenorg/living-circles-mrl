package integration_test

import (
	"context"
	"encoding/json"
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

func TestClientReceivesFightResolutionWithChildReplacement(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "triangle",
		PlayerEnergy:            100,
		AutonomousEnergy:        80,
		AutonomousChildrenCount: 1,
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

		if snapshot.Interaction.Kind != "fight_resolved" {
			t.Fatalf("expected fight_resolved, got %q", snapshot.Interaction.Kind)
		}
		if snapshot.Interaction.LoserID != simulation.DefaultAutonomousID {
			t.Fatalf("expected autonomous loser, got %q", snapshot.Interaction.LoserID)
		}
		if len(snapshot.AutonomousCircles) != 1 {
			t.Fatalf("expected replacement autonomous circle to remain active, got %d", len(snapshot.AutonomousCircles))
		}
		if snapshot.AutonomousCircles[0].ChildrenCount != 0 {
			t.Fatalf("expected replacement to consume one child, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
		}
		if snapshot.AutonomousCircles[0].Energy != simulation.DefaultReplacementEnergy {
			t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.AutonomousCircles[0].Energy)
		}
		return
	}

	t.Fatal("expected fight resolution with child replacement")
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
	if initial.Player.ChildrenCount != 0 {
		t.Fatalf("expected player child count to start at zero, got %d", initial.Player.ChildrenCount)
	}
	if initial.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected first autonomous child count to start at one for demo continuity, got %d", initial.AutonomousCircles[0].ChildrenCount)
	}
	if initial.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatalf("expected second autonomous child count to start at zero, got %d", initial.AutonomousCircles[1].ChildrenCount)
	}
	if initial.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected base player radius %v, got %v", simulation.DefaultPlayerRadius, initial.Player.Radius)
	}
	if initial.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain {
		t.Fatalf("expected first autonomous radius %v, got %v", simulation.DefaultPlayerRadius+simulation.DefaultChildRadiusGain, initial.AutonomousCircles[0].Radius)
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

		if snapshot.Tick == 0 {
			continue
		}

		if snapshot.Player == nil {
			t.Fatal("expected replacement player to remain active after zero-energy collapse")
		}
		if snapshot.Player.ChildrenCount != 0 {
			t.Fatalf("expected replacement to consume one child, got %d", snapshot.Player.ChildrenCount)
		}
		if snapshot.Player.Energy != simulation.DefaultReplacementEnergy {
			t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.Player.Energy)
		}
		return
	}

	t.Fatal("expected energy-collapse replacement snapshot")
}
