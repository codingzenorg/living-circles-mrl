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

type movementIntentMessage struct {
	Type      string            `json:"type"`
	Direction simulation.Vector `json:"direction"`
}

type clientConnection struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
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

type Server struct {
	mux       *http.ServeMux
	session   *simulation.Session
	tickEvery time.Duration
	upgrader  websocket.Upgrader

	mu    sync.Mutex
	conns map[*websocket.Conn]*clientConnection
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
			s.broadcastSnapshot(s.session.Advance())
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
	if err := client.WriteJSON(s.session.Snapshot()); err != nil {
		s.removeConnection(connection)
		_ = client.Close()
		return
	}

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
	}
}

func (s *Server) handleReset(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snapshot := s.session.Reset()
	s.broadcastSnapshot(snapshot)

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(snapshot)
}

func (s *Server) broadcastSnapshot(snapshot simulation.Snapshot) {
	s.mu.Lock()
	connections := make([]*clientConnection, 0, len(s.conns))
	for _, connection := range s.conns {
		connections = append(connections, connection)
	}
	s.mu.Unlock()

	for _, connection := range connections {
		if err := connection.WriteJSON(snapshot); err != nil {
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
