package taquilla

import "sync"
import "math"

// ===========================================================================
// semGraphStruct describes a struct that holds the memory characteristics of each process
type semGraphStruct struct {
	sync.RWMutex
	container      map[string]*semCounterStruct
	availableTypes []string // Stores all the available types in a sequence
}

// set sets the average memory of the specified semaphore
func (s *semGraphStruct) set(ticketType string) {
	s.Lock()
	defer s.Unlock()
	// Check if the ticketType already exists in the map
	if s.container[ticketType] == nil {
		// Otherwise, create a new entry here
		// log.Errorln("Creating new entry in semGraphStruct")
		s.availableTypes = append(s.availableTypes, ticketType)
		var newCounter = new(semCounterStruct)
		newCounter.name = ticketType
		s.container[ticketType] = newCounter
	}

	s.container[ticketType].incrementCounter()
}

// len returns the number of keys in the container
func (s *semGraphStruct) len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.availableTypes)
}

// updateMem accepts the ticketType and the memory consumed and updates the avgMem of the underlying semCounterStruct
func (s *semGraphStruct) updateMem(ticketType string, memoryConsumed float64) {
	s.Lock()
	defer s.Unlock()
	s.container[ticketType].updateAvgMem(memoryConsumed)
}

// semCounterStruct stores the average memory used by each process, as well as a counter to denote the number of such processes that have already run
type semCounterStruct struct {
	// sync.RWMutex
	avgMem  float64
	counter uint64
	name    string
}

// incrementCounter atomically increases the counter value by 1
func (c *semCounterStruct) incrementCounter() {
	// c.Lock()
	// defer c.Unlock()
	c.counter++
	// log.Errorln("Counter for ", c.name, " has become ", c.counter)
}

// updateAvgMem accepts the memoryConsumed and updates the average memory of this semaphore
func (c *semCounterStruct) updateAvgMem(memoryConsumed float64) {
	// In case average memory hasn't been previously set, and a sudden spike appears, then it will be set without averaging out the value
	if c.avgMem == 0 {
		c.avgMem = memoryConsumed
		return
	}
	if math.Abs(c.avgMem-memoryConsumed) < 0.1 {
		// The memory consumed is the same as the current average --> calculation can be avoided
		return
	}
	sum := float64(c.counter)*c.avgMem + math.Abs(memoryConsumed)

	if c.counter == 0 {
		return
	}
	c.avgMem = sum / float64(c.counter)

}
