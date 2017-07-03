package taquilla

import "sync"

// ===========================================================================
// semGraphStruct describes a struct that holds the memory characteristics of each process
type semGraphStruct struct {
	sync.RWMutex
	container map[string]*semCounterStruct
}

// set sets the average memory of the specified semaphore
func (s *semGraphStruct) set(semName string) {
	s.Lock()
	defer s.Unlock()
	s.container[semName].updateAvg()
}

// semCounterStruct stores the average memory used by each process, as well as a counter to denote the number of such processes that have already run
type semCounterStruct struct {
	sync.RWMutex
	avgMem  uint64
	counter uint64
}

// updateAvg increases the counter by 1, and calculates the average memory for this counter variable
func (c *semCounterStruct) updateAvg() {
	c.Lock()
	defer c.Unlock()
}
