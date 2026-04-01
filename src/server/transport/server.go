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

type movementIntentMessage struct {
	Type      string            `json:"type"`
	Direction simulation.Vector `json:"direction"`
}

type Server struct {
	mux       *http.ServeMux
	session   *simulation.Session
	tickEvery time.Duration
	upgrader  websocket.Upgrader

	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
}

func NewServer() *Server {
	server := &Server{
		mux:       http.NewServeMux(),
		session:   simulation.NewSession(),
		tickEvery: 100 * time.Millisecond,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
		conns: make(map[*websocket.Conn]struct{}),
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
	s.mux.Handle(websocketPath, http.HandlerFunc(s.handleWebSocket))
	s.mux.Handle("/", http.FileServer(http.Dir("src/client")))
	s.mux.Handle("/shared_contracts/", http.StripPrefix("/shared_contracts/", http.FileServer(http.Dir("src/shared_contracts"))))
}

func (s *Server) handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	connection, err := s.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	s.addConnection(connection)
	if err := connection.WriteJSON(s.session.Snapshot()); err != nil {
		s.removeConnection(connection)
		_ = connection.Close()
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

func (s *Server) broadcastSnapshot(snapshot simulation.Snapshot) {
	s.mu.Lock()
	connections := make([]*websocket.Conn, 0, len(s.conns))
	for connection := range s.conns {
		connections = append(connections, connection)
	}
	s.mu.Unlock()

	for _, connection := range connections {
		if err := connection.WriteJSON(snapshot); err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) {
				s.removeConnection(connection)
			}
			_ = connection.Close()
		}
	}
}

func (s *Server) addConnection(connection *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[connection] = struct{}{}
}

func (s *Server) removeConnection(connection *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.conns, connection)
}

func (s *Server) closeConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for connection := range s.conns {
		_ = connection.Close()
		delete(s.conns, connection)
	}
}
