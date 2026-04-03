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
		if snapshot.AutonomousCircles[0].LineageID != initial.AutonomousCircles[0].LineageID {
			t.Fatalf("expected replacement lineage %q, got %q", initial.AutonomousCircles[0].LineageID, snapshot.AutonomousCircles[0].LineageID)
		}
		if snapshot.AutonomousCircles[0].Generation != initial.AutonomousCircles[0].Generation+1 {
			t.Fatalf("expected replacement generation %d, got %d", initial.AutonomousCircles[0].Generation+1, snapshot.AutonomousCircles[0].Generation)
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
	if initial.Player.LineageID != "lineage-player-1" || initial.Player.Generation != 0 {
		t.Fatalf("expected initial player lineage state, got lineage=%q generation=%d", initial.Player.LineageID, initial.Player.Generation)
	}
	if initial.AutonomousCircles[0].ChildrenCount != 1 {
		t.Fatalf("expected first autonomous child count to start at one for demo continuity, got %d", initial.AutonomousCircles[0].ChildrenCount)
	}
	if initial.AutonomousCircles[0].LineageID != "lineage-circle-2" || initial.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected first autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[0].LineageID, initial.AutonomousCircles[0].Generation)
	}
	if initial.AutonomousCircles[1].ChildrenCount != 0 {
		t.Fatalf("expected second autonomous child count to start at zero, got %d", initial.AutonomousCircles[1].ChildrenCount)
	}
	if initial.AutonomousCircles[1].LineageID != "lineage-circle-3" || initial.AutonomousCircles[1].Generation != 0 {
		t.Fatalf("expected second autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[1].LineageID, initial.AutonomousCircles[1].Generation)
	}
	if initial.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected base player radius %v, got %v", simulation.DefaultPlayerRadius, initial.Player.Radius)
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
			if snapshot.Player.Energy >= initial.Player.Energy {
				t.Fatalf("expected player energy to decrease after reproduction, before=%v after=%v", initial.Player.Energy, snapshot.Player.Energy)
			}
			if snapshot.Player.ChildrenCount != 1 || snapshot.AutonomousCircles[0].ChildrenCount != 1 {
				t.Fatalf("expected one child count for both circles after reproduction, player=%d autonomous=%d", snapshot.Player.ChildrenCount, snapshot.AutonomousCircles[0].ChildrenCount)
			}
			if snapshot.AutonomousCircles[0].Energy >= initial.AutonomousCircles[0].Energy {
				t.Fatalf("expected autonomous energy to decrease after reproduction, before=%v after=%v", initial.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
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

func TestClientReceivesBlockedReproductionWhenEnergyIsInsufficient(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:      "triangle",
		AutonomousShape:  "square",
		PlayerEnergy:     simulation.DefaultPlayerEnergy,
		AutonomousEnergy: simulation.DefaultReproductionCost - 1,
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
		return
	}

	t.Fatal("expected blocked reproduction snapshot")
}

func TestClientReceivesReproductionPaidByChildWhenEnergyIsLow(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:             "triangle",
		AutonomousShape:         "square",
		PlayerEnergy:            simulation.DefaultPlayerEnergy,
		AutonomousEnergy:        simulation.DefaultReproductionCost - 1,
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
		if snapshot.AutonomousCircles[0].ChildrenCount != 1 {
			t.Fatalf("expected autonomous circle child count to stay at 1 after child payment and new reproduction, got %d", snapshot.AutonomousCircles[0].ChildrenCount)
		}
		expectedEnergy := previous.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
		if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
			t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
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
		if snapshot.Player.Energy != simulation.DefaultReplacementEnergy {
			t.Fatalf("expected replacement energy %v, got %v", simulation.DefaultReplacementEnergy, snapshot.Player.Energy)
		}
		if snapshot.Player.LineageID != initial.Player.LineageID {
			t.Fatalf("expected replacement lineage %q, got %q", initial.Player.LineageID, snapshot.Player.LineageID)
		}
		if snapshot.Player.Generation != initial.Player.Generation+1 {
			t.Fatalf("expected replacement generation %d, got %d", initial.Player.Generation+1, snapshot.Player.Generation)
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

func foodByID(foods []simulation.Food, id string) (simulation.Food, bool) {
	for _, food := range foods {
		if food.ID == id {
			return food, true
		}
	}

	return simulation.Food{}, false
}
