package resource

import (
	"errors"
	"runtime"
	"strings"
	"sync"
)

var ErrCapacityFull = errors.New("capacity is full for this game type")

type Class string

const (
	Light   Class = "light"
	Physics Class = "physics"
	HeavyAI Class = "heavy_ai"
)

// GetClass categorizes the game based on its ID.
func GetClass(gameID string) Class {
	// We can expand this list as new games are added by ALs.
	switch strings.ToLower(gameID) {
	case "billiards", "football", "pong", "car-race":
		return Physics
	case "chess", "checkers":
		return HeavyAI
	default:
		// "ludo", "tic-tac-toe", "backgammon" fall here
		return Light
	}
}

type Manager struct {
	mu          sync.RWMutex
	maxRooms    map[Class]int
	activeRooms map[Class]int
}

func NewManager() *Manager {
	cpus := runtime.NumCPU()
	
	// Conservative limits for a 3GB RAM / 1 CPU server
	maxLight := cpus * 100
	maxPhysics := cpus * 20
	maxHeavyAI := cpus * 10

	return &Manager{
		maxRooms: map[Class]int{
			Light:   maxLight,
			Physics: maxPhysics,
			HeavyAI: maxHeavyAI,
		},
		activeRooms: make(map[Class]int),
	}
}

func (m *Manager) RequestRoom(class Class) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeRooms[class] >= m.maxRooms[class] {
		return ErrCapacityFull
	}
	
	m.activeRooms[class]++
	return nil
}

func (m *Manager) ReleaseRoom(class Class) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeRooms[class] > 0 {
		m.activeRooms[class]--
	}
}
