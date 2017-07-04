package taquilla

import "sync"

// ===========================================================================
// memGraph describes a struct that holds the sys allocated memory
type memGraph struct {
	sync.RWMutex
	value map[int64]sysMemAndOnlineSem
}

// set acquires an exclusive lock and updates the map with the current timestamp and the sys memory used
func (m *memGraph) set(timestamp int64, sysMem uint64) {
	m.Lock()
	defer m.Unlock()
	var localvar sysMemAndOnlineSem
	localvar.sysMem = sysMem
	localvar.online = online.getAll()
	m.value[timestamp] = localvar
	// log.Errorln("\n\n================================\nMem Graph is\n", len(localvar.online), "\n\n=================================")
}

// get acquires a read only lock and returns the sys memory at a particular timestamp
func (m *memGraph) get(timestamp int64) sysMemAndOnlineSem {
	m.RLock()
	defer m.RUnlock()
	return m.value[timestamp]
}

// sysMemAndOnlineSem is a struct that stores the system memory as well as the list of online semaphores running at that particular timestamp
type sysMemAndOnlineSem struct {
	sysMem uint64
	online []semaphore
}
