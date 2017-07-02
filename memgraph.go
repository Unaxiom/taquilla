package taquilla

import "sync"

// ===========================================================================
// memGraph describes a struct that holds the sys allocated memory
type memGraph struct {
	sync.RWMutex
	value map[int64]uint64
}

// set acquires an exclusive lock and updates the map with the current timestamp and the sys memory used
func (m *memGraph) set(timestamp int64, sysMem uint64) {
	m.Lock()
	defer m.Unlock()
	m.value[timestamp] = sysMem
}

// get acquires a read only lock and returns the sys memory at a particular timestamp
func (m *memGraph) get(timestamp int64) uint64 {
	m.RLock()
	defer m.RUnlock()
	return m.value[timestamp]
}
