package resource

import (
	"errors"
	"runtime"
	"sync"
)

var ErrCapacityFull = errors.New("capacity is full for this game type")

type Class string

const (
	Light   Class = "light"    // Ludo, Tic-Tac-Toe
	Physics Class = "physics"  // Billiards, Football
	HeavyAI Class = "heavy_ai" // Chess
)

type Manager struct {
	mu          sync.RWMutex
	maxRooms    map[Class]int
	activeRooms map[Class]int
}

// NewManager automatically detects hardware and sets safe limits.
func NewManager() *Manager {
	cpus := runtime.NumCPU()
	
	// If the server has 1 CPU, these numbers will be small.
	// If you upgrade the server on Runflare, these numbers will increase automatically.
	// We use a conservative multiplier to prevent the 3GB RAM from OOM (Out Of Memory).
	
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

// GetStats returns current usage so the Hub can decide whether to accept new players.
func (m *Manager) GetStats() (map[Class]int, map[Class]int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	maxCopy := make(map[Class]int)
	activeCopy := make(map[Class]int)
	for k, v := range m.maxRooms {
		maxCopy[k] = v
	}
	for k, v := range m.activeRooms {
		activeCopy[k] = v
	}
	return maxCopy, activeCopy
}
