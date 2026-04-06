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

func TestClientReceivesCrowdingEnergyPressureInClusteredWorld(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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

	for {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read updated snapshot: %v", err)
		}

		if snapshot.Tick == 0 || snapshot.Player == nil {
			continue
		}

		expectedEnergy := initial.Player.Energy - simulation.DefaultMoveCost - simulation.DefaultCrowdingMoveCost
		if snapshot.Player.Energy != expectedEnergy {
			t.Fatalf("expected crowded player energy %v, got %v", expectedEnergy, snapshot.Player.Energy)
		}
		return
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
	expectedChildID := initial.AutonomousCircles[0].AttachedChildren[0].ID

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
		if snapshot.Interaction.AbsorbedChildID != expectedChildID {
			t.Fatalf("expected absorbed child id %q, got %q", expectedChildID, snapshot.Interaction.AbsorbedChildID)
		}
		if len(snapshot.AutonomousCircles) != 1 {
			t.Fatalf("expected autonomous loser to remain active, got %d", len(snapshot.AutonomousCircles))
		}
		if len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
			t.Fatalf("expected absorbed conflict to consume one child, got %d", len(snapshot.AutonomousCircles[0].AttachedChildren))
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
		if snapshot.Interaction.ContactPathKind != "source_child_to_target_parent" {
			t.Fatalf("expected source_child_to_target_parent path kind, got %q", snapshot.Interaction.ContactPathKind)
		}
		if snapshot.Interaction.SourceChildID != "player-1-child-1" {
			t.Fatalf("expected source child id player-1-child-1, got %q", snapshot.Interaction.SourceChildID)
		}
		if snapshot.Interaction.TargetChildID != "" {
			t.Fatalf("expected empty target child id, got %q", snapshot.Interaction.TargetChildID)
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
		AutonomousX:               238.511,
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
			t.Fatalf("expected both parents to remain, player=%v autonomous=%d", snapshot.Player != nil, len(snapshot.AutonomousCircles))
		}
		return
	}

	t.Fatal("expected child-to-child reproduction snapshot")
}

func TestClientDoesNotReceiveContactFromDerivedRadiusAlone(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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

	for range 6 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}
		if snapshot.Tick == 0 {
			continue
		}
		if snapshot.Interaction != nil {
			t.Fatalf("expected no interaction when only derived parent radius would have overlapped, got %+v", snapshot.Interaction)
		}
		return
	}

	t.Fatal("expected at least one post-initial snapshot without interaction")
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

func TestClientReceivesChildAwareDifferentShapePursuitBeforeParentCoreNearest(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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

	playerParentDistance := math.Hypot(initial.Player.X-initial.AutonomousCircles[0].X, initial.Player.Y-initial.AutonomousCircles[0].Y)
	secondaryDistance := math.Hypot(initial.AutonomousCircles[1].X-initial.AutonomousCircles[0].X, initial.AutonomousCircles[1].Y-initial.AutonomousCircles[0].Y)
	if playerParentDistance <= secondaryDistance {
		t.Fatalf("expected player parent body to stay farther than secondary target, player=%v secondary=%v", playerParentDistance, secondaryDistance)
	}
	playerChildDistance := playerParentDistance
	for _, child := range initial.Player.AttachedChildren {
		distance := math.Hypot(child.X-initial.AutonomousCircles[0].X, child.Y-initial.AutonomousCircles[0].Y)
		if distance < playerChildDistance {
			playerChildDistance = distance
		}
	}
	if playerChildDistance >= secondaryDistance {
		t.Fatalf("expected player attached child to be effectively nearer than secondary target, child=%v secondary=%v", playerChildDistance, secondaryDistance)
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
			t.Fatalf("expected no immediate interaction during child-aware pursuit, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X > initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to move toward child-aware different-shape target")
}

func TestClientReceivesChildAwareFoodTargetingBeforeParentCoreNearest(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
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

	parentToFoodOne := math.Hypot(initial.Foods[0].X-initial.AutonomousCircles[0].X, initial.Foods[0].Y-initial.AutonomousCircles[0].Y)
	parentToFoodTwo := math.Hypot(initial.Foods[1].X-initial.AutonomousCircles[0].X, initial.Foods[1].Y-initial.AutonomousCircles[0].Y)
	if parentToFoodOne >= parentToFoodTwo {
		t.Fatalf("expected parent body to be nearer to food-1 than food-2, food-1=%v food-2=%v", parentToFoodOne, parentToFoodTwo)
	}

	childToFoodTwo := parentToFoodTwo
	for _, child := range initial.AutonomousCircles[0].AttachedChildren {
		distance := math.Hypot(initial.Foods[1].X-child.X, initial.Foods[1].Y-child.Y)
		if distance < childToFoodTwo {
			childToFoodTwo = distance
		}
	}
	if childToFoodTwo >= parentToFoodOne {
		t.Fatalf("expected attached child to make food-2 effectively nearer than food-1, food-1=%v food-2=%v", parentToFoodOne, childToFoodTwo)
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
			t.Fatalf("expected no immediate interaction during child-aware food targeting, got %+v", snapshot.Interaction)
		}
		if snapshot.AutonomousCircles[0].X < initial.AutonomousCircles[0].X {
			return
		}
	}

	t.Fatal("expected autonomous circle to move left toward child-aware food target")
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
	if initial.World.Width != simulation.DefaultExpandedWorldWidth || initial.World.Height != simulation.DefaultExpandedWorldHeight {
		t.Fatalf("expected expanded world %vx%v, got %vx%v", simulation.DefaultExpandedWorldWidth, simulation.DefaultExpandedWorldHeight, initial.World.Width, initial.World.Height)
	}
	if len(initial.AutonomousCircles) != simulation.DefaultExpandedAutonomousCount {
		t.Fatalf("expected %d autonomous circles, got %d", simulation.DefaultExpandedAutonomousCount, len(initial.AutonomousCircles))
	}
	if len(initial.Foods) != simulation.DefaultExpandedFoodCount {
		t.Fatalf("expected %d food slots, got %d", simulation.DefaultExpandedFoodCount, len(initial.Foods))
	}
	if initial.AutonomousCircles[0].Shape != initial.Player.Shape {
		t.Fatalf("expected first autonomous circle to match player shape %q, got %q", initial.Player.Shape, initial.AutonomousCircles[0].Shape)
	}
	if initial.AutonomousCircles[1].Shape == initial.Player.Shape {
		t.Fatalf("expected second autonomous circle to differ from player shape %q", initial.Player.Shape)
	}
	if len(initial.Player.AttachedChildren) != 1 {
		t.Fatalf("expected player child count to start at one for demo continuity, got %d", len(initial.Player.AttachedChildren))
	}
	if len(initial.Player.AttachedChildren) != len(initial.Player.AttachedChildren) {
		t.Fatalf("expected player attached children to match count, count=%d attached=%d", len(initial.Player.AttachedChildren), len(initial.Player.AttachedChildren))
	}
	if initial.Player.LineageID != "lineage-player-1" || initial.Player.Generation != 0 {
		t.Fatalf("expected initial player lineage state, got lineage=%q generation=%d", initial.Player.LineageID, initial.Player.Generation)
	}
	if len(initial.AutonomousCircles[0].AttachedChildren) != 1 {
		t.Fatalf("expected first autonomous child count to start at one for demo continuity, got %d", len(initial.AutonomousCircles[0].AttachedChildren))
	}
	if len(initial.AutonomousCircles[0].AttachedChildren) != len(initial.AutonomousCircles[0].AttachedChildren) {
		t.Fatalf("expected first autonomous attached children to match count, count=%d attached=%d", len(initial.AutonomousCircles[0].AttachedChildren), len(initial.AutonomousCircles[0].AttachedChildren))
	}
	if initial.AutonomousCircles[0].LineageID != "lineage-circle-2" || initial.AutonomousCircles[0].Generation != 0 {
		t.Fatalf("expected first autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[0].LineageID, initial.AutonomousCircles[0].Generation)
	}
	if len(initial.AutonomousCircles[1].AttachedChildren) != 0 {
		t.Fatalf("expected second autonomous child count to start at zero, got %d", len(initial.AutonomousCircles[1].AttachedChildren))
	}
	if len(initial.AutonomousCircles[1].AttachedChildren) != len(initial.AutonomousCircles[1].AttachedChildren) {
		t.Fatalf("expected second autonomous attached children to match count, count=%d attached=%d", len(initial.AutonomousCircles[1].AttachedChildren), len(initial.AutonomousCircles[1].AttachedChildren))
	}
	if initial.AutonomousCircles[1].LineageID != "lineage-circle-3" || initial.AutonomousCircles[1].Generation != 0 {
		t.Fatalf("expected second autonomous lineage state, got lineage=%q generation=%d", initial.AutonomousCircles[1].LineageID, initial.AutonomousCircles[1].Generation)
	}
	if initial.AutonomousCircles[2].ID != simulation.DefaultTertiaryID || initial.AutonomousCircles[3].ID != simulation.DefaultQuaternaryID || initial.AutonomousCircles[4].ID != simulation.DefaultQuinaryID {
		t.Fatalf("expected deterministic expanded autonomous ids, got %+v", initial.AutonomousCircles)
	}
	if initial.Player.Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed player radius %v, got %v", simulation.DefaultPlayerRadius, initial.Player.Radius)
	}
	if initial.AutonomousCircles[0].Radius != simulation.DefaultPlayerRadius {
		t.Fatalf("expected fixed first autonomous radius %v, got %v", simulation.DefaultPlayerRadius, initial.AutonomousCircles[0].Radius)
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
			if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != len(initial.Player.AttachedChildren)+len(initial.AutonomousCircles[0].AttachedChildren)+2 {
				t.Fatalf("expected reproduction to distribute two children across the pair, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
			}
			if snapshot.AutonomousCircles[0].Energy >= initial.AutonomousCircles[0].Energy {
				t.Fatalf("expected autonomous energy to decrease after reproduction, before=%v after=%v", initial.AutonomousCircles[0].Energy, snapshot.AutonomousCircles[0].Energy)
			}
			if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) || len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
				t.Fatalf("expected attached children to match counts after reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
			}
			continue
		}

		if resolvedSeen {
			if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) > len(initial.Player.AttachedChildren)+len(initial.AutonomousCircles[0].AttachedChildren)+2 {
				t.Fatalf("expected no repeat accumulation during continuous overlap, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
			}
			if len(snapshot.Player.AttachedChildren) > 0 && len(snapshot.Player.AttachedChildren) > 0 {
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
		if len(snapshot.Interaction.CreatedChildIDs) != 0 {
			t.Fatalf("expected no created child ids on blocked reproduction, got %+v", snapshot.Interaction.CreatedChildIDs)
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
		if snapshot.Player == nil {
			t.Fatal("expected player to remain active after blocked reproduction")
		}
		if len(snapshot.Player.AttachedChildren) != 0 || len(snapshot.AutonomousCircles[0].AttachedChildren) != 0 {
			t.Fatalf("expected blocked reproduction to preserve child counts, player=%d autonomous=%d", len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
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
		PlayerX:                             200,
		PlayerY:                             300,
		AutonomousX:                         228,
		AutonomousY:                         300,
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultReproductionCost - 1,
		AutonomousChildrenCount:             1,
		DisableBlockedReproductionAvoidance: true,
		DisableFoodSeeking:                  true,
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

		if snapshot.Interaction.Kind != "reproduce_paid_child" {
			t.Fatalf("expected reproduce_paid_child, got %q", snapshot.Interaction.Kind)
		}
		if len(snapshot.Interaction.CreatedChildIDs) != 2 {
			t.Fatalf("expected two created child ids, got %d", len(snapshot.Interaction.CreatedChildIDs))
		}
		sourceCreated := newlyOwnedChildIDs(previous.Player.AttachedChildren, snapshot.Player.AttachedChildren)
		targetCreated := newlyOwnedChildIDs(previous.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren)
		if snapshot.Interaction.SourcePaidChild {
			t.Fatal("expected player not to pay with child in this reproduction path")
		}
		if !snapshot.Interaction.TargetPaidChild {
			t.Fatal("expected autonomous target to pay with child in this reproduction path")
		}
		assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
		assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
		assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
		assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
		assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
		assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
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
		for _, childID := range newlyCreatedChildIDs(previous.Player.AttachedChildren, snapshot.Player.AttachedChildren, previous.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren) {
			if _, exists := createdChildIDs[childID]; !exists {
				t.Fatalf("expected created child id %q in interaction payload, got %+v", childID, snapshot.Interaction.CreatedChildIDs)
			}
		}
		assertChildIDSetEqual(t, snapshot.Interaction.SourceCreatedChildIDs, sourceCreated, "source created child ids")
		assertChildIDSetEqual(t, snapshot.Interaction.TargetCreatedChildIDs, targetCreated, "target created child ids")
		assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
		if len(snapshot.Player.AttachedChildren)+len(snapshot.AutonomousCircles[0].AttachedChildren) != len(previous.Player.AttachedChildren)+len(previous.AutonomousCircles[0].AttachedChildren)+1 {
			t.Fatalf("expected one child payment plus two redistributed children, before player=%d autonomous=%d after player=%d autonomous=%d", len(previous.Player.AttachedChildren), len(previous.AutonomousCircles[0].AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
		}
		expectedEnergy := previous.AutonomousCircles[0].Energy - simulation.DefaultMoveCost
		if snapshot.AutonomousCircles[0].Energy != expectedEnergy {
			t.Fatalf("expected autonomous energy to be %v after movement plus child payment, got %v", expectedEnergy, snapshot.AutonomousCircles[0].Energy)
		}
		if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) || len(snapshot.AutonomousCircles[0].AttachedChildren) != len(snapshot.AutonomousCircles[0].AttachedChildren) {
			t.Fatalf("expected attached children to match counts after child-payment reproduction, player attached=%d count=%d autonomous attached=%d count=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren), len(snapshot.AutonomousCircles[0].AttachedChildren))
		}
		return
	}

	t.Fatal("expected reproduction snapshot paid by child")
}

func TestClientReceivesOrdinaryResolvedReproductionWhenNoChildPaymentIsUsed(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:                         "triangle",
		AutonomousShape:                     "square",
		PlayerEnergy:                        simulation.DefaultPlayerEnergy,
		AutonomousEnergy:                    simulation.DefaultPlayerEnergy,
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
			t.Fatalf("expected reproduce_resolved for energy-paid reproduction, got %q", snapshot.Interaction.Kind)
		}
		assertFloatEqual(t, snapshot.Interaction.SourceCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.Player.Energy, len(snapshot.Player.AttachedChildren)))
		assertFloatEqual(t, snapshot.Interaction.TargetCapacityValue, expectedReportedCapacity(snapshot.Interaction.Kind, snapshot.AutonomousCircles[0].Energy, len(snapshot.AutonomousCircles[0].AttachedChildren)))
		assertCapacityComponentsMatch(t, snapshot.Interaction.SourceCapacityValue, snapshot.Interaction.SourceEnergyComponent, snapshot.Interaction.SourceReserveComponent)
		assertCapacityComponentsMatch(t, snapshot.Interaction.TargetCapacityValue, snapshot.Interaction.TargetEnergyComponent, snapshot.Interaction.TargetReserveComponent)
		assertFloatEqual(t, snapshot.Interaction.ReproductionThreshold, simulation.DefaultReproductionMinEnergy)
		assertFloatEqual(t, snapshot.Interaction.ReproductionCost, simulation.DefaultReproductionCost)
		sourceCreated := newlyOwnedChildIDs(previous.Player.AttachedChildren, snapshot.Player.AttachedChildren)
		targetCreated := newlyOwnedChildIDs(previous.AutonomousCircles[0].AttachedChildren, snapshot.AutonomousCircles[0].AttachedChildren)
		if len(snapshot.Interaction.CreatedChildIDs) != 2 {
			t.Fatalf("expected two created child ids for energy-paid reproduction, got %d", len(snapshot.Interaction.CreatedChildIDs))
		}
		assertChildIDSetEqual(t, snapshot.Interaction.SourceCreatedChildIDs, sourceCreated, "source created child ids")
		assertChildIDSetEqual(t, snapshot.Interaction.TargetCreatedChildIDs, targetCreated, "target created child ids")
		assertDistributionKindMatchesOwnership(t, snapshot.Interaction.DistributionKind, sourceCreated, targetCreated)
		return
	}

	t.Fatal("expected ordinary resolved reproduction snapshot")
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
		PlayerEnergy:        0,
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
	expectedChildID := initial.Player.AttachedChildren[0].ID
	expectedX, expectedY := promotedChildPosition(initial.Player.ID, initial.Player.AttachedChildren[0].ID, initial.Player.X, initial.Player.Y, 1)

	_ = connection.SetReadDeadline(time.Now().Add(3 * time.Second))

	for range 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Player == nil {
			t.Fatal("expected replacement player to remain active after zero-energy collapse")
		}
		if len(snapshot.Player.AttachedChildren) != 0 {
			t.Fatalf("expected replacement to consume one child, got %d", len(snapshot.Player.AttachedChildren))
		}
		if len(snapshot.Player.AttachedChildren) != len(snapshot.Player.AttachedChildren) {
			t.Fatalf("expected attached children to remain synchronized, count=%d attached=%d", len(snapshot.Player.AttachedChildren), len(snapshot.Player.AttachedChildren))
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
		if math.Abs(snapshot.Player.X-expectedX) > 1e-9 || math.Abs(snapshot.Player.Y-expectedY) > 1e-9 {
			t.Fatalf("expected replacement at promoted child position (%v, %v), got (%v, %v)", expectedX, expectedY, snapshot.Player.X, snapshot.Player.Y)
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
		if snapshot.Interaction.PromotedChildID != expectedChildID {
			t.Fatalf("expected promoted child id %q, got %q", expectedChildID, snapshot.Interaction.PromotedChildID)
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
		initialFood, _ := foodByID(initial.Foods, "food-1")
		if food.X != initialFood.X || food.Y != initialFood.Y {
			t.Fatalf("expected regenerated food-1 to return to its original slot, got %+v", food)
		}
		return
	}

	t.Fatal("expected regenerated food after deterministic delay")
}

func TestClientSeesPressureScaledFoodRegenDelayInExpandedWorld(t *testing.T) {
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

	_ = connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	sawTwoMissing := false
	for range simulation.DefaultFoodRegenDelay + 20 {
		var snapshot simulation.Snapshot
		if err := connection.ReadJSON(&snapshot); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}

		if snapshot.Tick == 0 {
			continue
		}

		if len(snapshot.Foods) <= len(initial.Foods)-2 {
			sawTwoMissing = true
			continue
		}

		if sawTwoMissing && snapshot.Tick <= simulation.DefaultFoodRegenDelay+4 && len(snapshot.Foods) == len(initial.Foods) {
			t.Fatalf("expected expanded world to stay partially depleted longer than the base regen delay, tick=%d foods=%d", snapshot.Tick, len(snapshot.Foods))
		}
	}

	if !sawTwoMissing {
		t.Fatal("expected expanded world to consume at least two food slots during the observation window")
	}
}

func TestClientReceivesFoodCollectionThroughAttachedChild(t *testing.T) {
	server := transport.NewServerWithSession(simulation.NewSessionWithConfig(simulation.Config{
		PlayerShape:               "triangle",
		AutonomousShape:           "square",
		SecondaryAutonomousShape:  "",
		PlayerX:                   403,
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
