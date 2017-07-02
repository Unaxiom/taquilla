package taquilla

import "sync"

// =============================================================================
// 		memory struct
type memory struct {
	sync.RWMutex
	available float64
}

// set sets the available memory to value passed
func (m *memory) set(amount float64) {
	m.Lock()
	defer m.Unlock()
	m.available = amount
}

// get gets the available memory
func (m *memory) get() float64 {
	m.Lock()
	defer m.Unlock()
	log.Debugln("Current memory is ", m.available)
	return m.available
}

// add adds the passed value to available memory
func (m *memory) add(amount float64) {
	m.Lock()
	defer m.Unlock()
	m.available += amount
}

// reduce removes the passed value from the available memory
func (m *memory) reduce(amount float64) {
	m.Lock()
	defer m.Unlock()
	m.available -= amount
}
