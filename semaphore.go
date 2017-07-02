package taquilla

// ============================================================================
// semaphore describes the internal representation of each semaphore
type semaphore struct {
	Token      string      // Stores the UUID token
	CallerChan chan string // Stores the channel that will be used to return the Token
	Seat       chan int    // This is the actual semaphore
	Type       string      // Stores the type of this semaphore
	ReqTime    int64       // Stores the timestamp of when the semaphore is requested
	RelTime    int64       // Stores the timestamp of when the semaphore is released
	ReqSemList []semaphore // Stores the semaphoreList when the semaphore is requested
	RelSemList []semaphore // Stores the semaphoreList when the semaphore is released
}
