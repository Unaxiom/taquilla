package taquilla

import "sync"

// =============================================================================
// 		semaphoreList struct
type semaphoreList struct {
	sync.RWMutex
	list []semaphore
}

func (l *semaphoreList) append(s semaphore) {
	l.Lock()
	defer l.Unlock()
	l.list = append(l.list, s)
}

func (l *semaphoreList) remove(s semaphore) {
	l.Lock()
	defer l.Unlock()
	for i, sem := range l.list {
		if sem.Token == s.Token {
			l.list = append(l.list[:i], l.list[i+1:]...)
			break
		}
	}
}

// removeByToken removes the semaphore from the list where the token matches
func (l *semaphoreList) removeByToken(token string) semaphore {
	l.Lock()
	defer l.Unlock()
	var semToReturn semaphore
	for i, sem := range l.list {
		if sem.Token == token {
			l.list = append(l.list[:i], l.list[i+1:]...)
			semToReturn = sem
			break
		}
	}
	return semToReturn
}

// Returns the length of the semaphore list
func (l *semaphoreList) length() int {
	l.RLock()
	defer l.RUnlock()
	length := len(l.list)
	return length
}

// getOne returns the first semaphore of the list
func (l *semaphoreList) getOne() (semaphore, bool) {
	l.RLock()
	defer l.RUnlock()
	if len(l.list) == 0 {
		return semaphore{}, false
	}
	return l.list[0], true
}

// getAll returns all the available semaphores in the list at that time
func (l *semaphoreList) getAll() []semaphore {
	l.RLock()
	defer l.RUnlock()
	return l.list[:]
}
