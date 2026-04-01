package simulation

import "sync"

type Session struct {
	mu     sync.Mutex
	tick   int64
	intent Vector
	world  *World
}

func NewSession() *Session {
	return &Session{
		world: NewWorld(),
	}
}

func (s *Session) ApplyIntent(intent Vector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.intent = intent
}

func (s *Session) Advance() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tick++
	return s.world.Advance(s.tick, s.intent)
}

func (s *Session) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.world.Snapshot(s.tick)
}
