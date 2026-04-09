package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/codingzen/living-circles-mrl/src/server/simulation"
)

const websocketPath = "/ws"

const DefaultTickEvery = 100 * time.Millisecond
const DefaultRecentMovementIntentTicks = int64(40)
const DefaultPassiveObserverCadenceTicks = int64(3)

type movementIntentMessage struct {
	Type      string            `json:"type"`
	Direction simulation.Vector `json:"direction"`
}

type clientConnection struct {
	conn *websocket.Conn

	writeMu sync.Mutex
	stateMu sync.Mutex

	hasRecentMovementIntent bool
	lastMovementIntentTick  int64
}

func (c *clientConnection) WriteJSON(value any) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteJSON(value)
}

func (c *clientConnection) Close() error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.Close()
}

func (c *clientConnection) RecordMovementIntent(currentTick int64) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	c.hasRecentMovementIntent = true
	c.lastMovementIntentTick = currentTick
}

func (c *clientConnection) IsActiveAtTick(currentTick int64) bool {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if !c.hasRecentMovementIntent {
		return false
	}

	return currentTick-c.lastMovementIntentTick <= DefaultRecentMovementIntentTicks
}

func (c *clientConnection) ShouldReceiveSnapshot(currentTick int64, reducePassiveCadence bool, force bool) bool {
	if force {
		return true
	}
	if !reducePassiveCadence {
		return true
	}
	if c.IsActiveAtTick(currentTick) {
		return true
	}

	return currentTick%DefaultPassiveObserverCadenceTicks == 0
}

type Server struct {
	mux       *http.ServeMux
	session   *simulation.Session
	tickEvery time.Duration
	upgrader  websocket.Upgrader

	mu                  sync.Mutex
	conns               map[*websocket.Conn]*clientConnection
	lastOrientationSig  string
	lastOrientationTick int64
	lastFoodSig         string
	lastFoodTick        int64
	lastBroadcastTick   int64
}

func NewServer() *Server {
	return NewServerWithSession(simulation.NewSession())
}

func NewServerWithShapes(playerShape, autonomousShape string) *Server {
	return NewServerWithSession(simulation.NewSessionWithShapes(playerShape, autonomousShape))
}

func NewServerWithSession(session *simulation.Session) *Server {
	server := &Server{
		mux:       http.NewServeMux(),
		session:   session,
		tickEvery: DefaultTickEvery,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
		conns: make(map[*websocket.Conn]*clientConnection),
	}

	server.routes()
	return server
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.tickEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.closeConnections()
			return nil
		case <-ticker.C:
			s.broadcastSnapshot(s.session.Advance(), false)
		}
	}
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	})
	s.mux.HandleFunc("/reset", s.handleReset)
	s.mux.Handle(websocketPath, http.HandlerFunc(s.handleWebSocket))
	s.mux.Handle("/", http.FileServer(http.Dir("src/client")))
	s.mux.Handle("/shared_contracts/", http.StripPrefix("/shared_contracts/", http.FileServer(http.Dir("src/shared_contracts"))))
}

func (s *Server) handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	connection, err := s.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	client := s.addConnection(connection)
	transportSnapshot := BuildViewportSnapshot(s.session.Snapshot(), true)
	if err := client.WriteJSON(transportSnapshot); err != nil {
		s.removeConnection(connection)
		_ = client.Close()
		return
	}
	s.recordOrientationRefresh(transportSnapshot)
	s.recordFoodRefresh(transportSnapshot)

	go s.readMessages(connection)
}

func (s *Server) readMessages(connection *websocket.Conn) {
	defer func() {
		s.removeConnection(connection)
		_ = connection.Close()
	}()

	for {
		_, payload, err := connection.ReadMessage()
		if err != nil {
			return
		}

		var message movementIntentMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			continue
		}

		if message.Type != "movement_intent" {
			continue
		}

		s.session.ApplyIntent(message.Direction)
		s.recordMovementIntent(connection)
	}
}

func (s *Server) handleReset(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snapshot := s.session.Reset()
	s.broadcastSnapshot(snapshot, true)

	writer.Header().Set("Content-Type", "application/json")
	transportSnapshot := BuildViewportSnapshot(snapshot, true)
	s.recordOrientationRefresh(transportSnapshot)
	s.recordFoodRefresh(transportSnapshot)
	_ = json.NewEncoder(writer).Encode(transportSnapshot)
}

func (s *Server) broadcastSnapshot(snapshot simulation.Snapshot, force bool) {
	activeSnapshot := BuildViewportSnapshot(snapshot, true)
	observerSnapshot := BuildObserverSnapshot(snapshot, true)
	orientationSignature := OrientationSummarySignature(activeSnapshot)
	foodSignature := LocalFoodSignature(activeSnapshot)
	includeOrientation := s.shouldRefreshOrientation(snapshot.Tick, orientationSignature)
	includeFoods := s.shouldRefreshFoods(snapshot.Tick, foodSignature)

	activeTransportSnapshot := activeSnapshot
	observerTransportSnapshot := observerSnapshot
	if !includeOrientation {
		activeTransportSnapshot.OrientationFresh = false
		activeTransportSnapshot.MinimapAutonomousCircles = nil
		activeTransportSnapshot.MinimapFoods = nil
		observerTransportSnapshot.OrientationFresh = false
		observerTransportSnapshot.MinimapAutonomousCircles = nil
		observerTransportSnapshot.MinimapFoods = nil
	} else {
		s.recordOrientationRefresh(activeSnapshot)
	}
	if !includeFoods {
		activeTransportSnapshot.FoodsFresh = false
		activeTransportSnapshot.Foods = nil
	} else {
		s.recordFoodRefresh(activeSnapshot)
	}

	s.recordBroadcastTick(snapshot.Tick)

	s.mu.Lock()
	connections := make([]*clientConnection, 0, len(s.conns))
	for _, connection := range s.conns {
		connections = append(connections, connection)
	}
	s.mu.Unlock()
	reducePassiveCadence := len(connections) > 1

	for _, connection := range connections {
		if !connection.ShouldReceiveSnapshot(snapshot.Tick, reducePassiveCadence, force) {
			continue
		}
		transportSnapshot := activeTransportSnapshot
		if reducePassiveCadence && !connection.IsActiveAtTick(snapshot.Tick) {
			transportSnapshot = observerTransportSnapshot
		}
		if err := connection.WriteJSON(transportSnapshot); err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) {
				s.removeConnection(connection.conn)
			}
			_ = connection.Close()
		}
	}
}

func (s *Server) addConnection(connection *websocket.Conn) *clientConnection {
	s.mu.Lock()
	defer s.mu.Unlock()
	client := &clientConnection{conn: connection}
	s.conns[connection] = client
	return client
}

func (s *Server) recordMovementIntent(connection *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := s.conns[connection]
	if client == nil {
		return
	}

	client.RecordMovementIntent(s.lastBroadcastTick)
}

func (s *Server) removeConnection(connection *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.conns, connection)
}

func (s *Server) closeConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, connection := range s.conns {
		_ = connection.Close()
		delete(s.conns, key)
	}
}

func (s *Server) recordBroadcastTick(tick int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastBroadcastTick = tick
}

func (s *Server) shouldRefreshOrientation(currentTick int64, currentSignature string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return ShouldRefreshOrientation(s.lastOrientationSig, s.lastOrientationTick, currentTick, currentSignature)
}

func (s *Server) recordOrientationRefresh(snapshot Snapshot) {
	if !snapshot.OrientationFresh {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastOrientationSig = OrientationSummarySignature(snapshot)
	s.lastOrientationTick = snapshot.Tick
}

func (s *Server) shouldRefreshFoods(currentTick int64, currentSignature string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return ShouldRefreshLocalFoods(s.lastFoodSig, s.lastFoodTick, currentTick, currentSignature)
}

func (s *Server) recordFoodRefresh(snapshot Snapshot) {
	if !snapshot.FoodsFresh {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastFoodSig = LocalFoodSignature(snapshot)
	s.lastFoodTick = snapshot.Tick
}
