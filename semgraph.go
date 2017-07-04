package taquilla

import "sync"

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

// semCounterStruct stores the average memory used by each process, as well as a counter to denote the number of such processes that have already run
type semCounterStruct struct {
	// sync.RWMutex
	avgMem  uint64
	counter uint64
	name    string
}

// incrementCounter atomically increases the counter value by 1
func (c *semCounterStruct) incrementCounter() {
	// c.Lock()
	// defer c.Unlock()
	c.counter++
	log.Errorln("Counter for ", c.name, " has become ", c.counter)
}
