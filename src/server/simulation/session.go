package simulation

import "sync"

type Session struct {
	mu     sync.Mutex
	tick   int64
	intent Vector
	world  *World
	config Config
}

func NewSession() *Session {
	return NewSessionWithConfig(Config{
		PlayerShape:               DefaultPlayerShape,
		AutonomousShape:           DefaultPlayerShape,
		SecondaryAutonomousShape:  DefaultAutoShape,
		PlayerEnergy:              DefaultPlayerEnergy,
		PlayerChildrenCount:       1,
		AutonomousEnergy:          80,
		SecondaryAutonomousEnergy: DefaultPlayerEnergy,
		AutonomousChildrenCount:   1,
	})
}

func NewSessionWithShapes(playerShape, autonomousShape string) *Session {
	return NewSessionWithConfig(Config{
		PlayerShape:      playerShape,
		AutonomousShape:  autonomousShape,
		PlayerEnergy:     DefaultPlayerEnergy,
		AutonomousEnergy: DefaultPlayerEnergy,
	})
}

func NewSessionWithConfig(config Config) *Session {
	return &Session{
		world:  NewWorldWithConfig(config),
		config: config,
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

func (s *Session) Reset() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tick = 0
	s.intent = Vector{}
	s.world = NewWorldWithConfig(s.config)

	return s.world.Snapshot(s.tick)
}
