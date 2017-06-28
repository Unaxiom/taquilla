package taquilla

import (
	"strings"

	"sync"

	"github.com/Unaxiom/ulogger"
	"github.com/twinj/uuid"
)

// Implementation:
// 1. Setup --> accepts the average memory of each long standing process
// 2. Req --> accepts a string chan; requests for a new semaphore --> fills up a semaphore and
// 				pushes it to the pipeline array, and calls next
// 3. processNextTicket --> processes a semaphore if sufficient memory is available,
//				and returns the associated UUID/token via the associated
//				string chan; the actual process can then begin
//
// 4. Rel --> accepts the token and removes the semaphore from the
//				pipeline and calls next

// log is the ulogger object
var log *ulogger.Logger

// currentAvailableMemory stores the available memory
var currentAvailableMemory memory

// memoryRequiredPerProcess stores the memory required for running each process, in megabytes
var memoryRequiredPerProcess float64

// pipeline consists of all the list of semaphores that are currently available
var pipeline semaphoreList

// online consists of all the semaphores that are under execution
var online semaphoreList

// Setup accepts the average memory requirement for a process
func Setup(memoryRequiredInMB float64) {
	memoryRequiredPerProcess = memoryRequiredInMB
	log = ulogger.New()
	// log.SetLogLevel(ulogger.DebugLevel)
	log.SetLogLevel(ulogger.FatalLevel)
	log.Debugln("Set memoryRequiredPerProcess to ", memoryRequiredPerProcess)
	currentAvailableMemory.set(float64(12))
}

// Req accepts a string channel via which an access token is returned. The caller function can then begin its execution.
func Req(returnChan chan string) {
	log.Debugln("Requested new semaphore")
	var ticket semaphore
	ticket.Token = strings.Join(strings.Split(uuid.NewV4().String(), "-"), "")
	log.Debugln("Generated new token --> ", ticket.Token)
	ticket.CallerChan = returnChan
	ticket.Type = "" // This can be implemented in a later version
	pipeline.append(ticket)
	go processNextTicket()
}

// processNextTicket processes the next available ticket in the pipeline
func processNextTicket() {
	log.Errorln("Started processNextTicket()")
	var ticket semaphore
	if pipeline.length() == 0 {
		log.Errorln("Length of pipeline is 0!")
		return
	}
	ticket = pipeline.getOne()
	// Check here if memory is available
	if memoryRequiredPerProcess > currentAvailableMemory.get() {
		log.Errorln("Ticket --> ", ticket.Token, " does not have sufficient memory to process...")
		return
	}
	// Process can be run here
	currentAvailableMemory.reduce(memoryRequiredPerProcess)
	ticket.CallerChan <- ticket.Token
	pipeline.remove(ticket)
	online.append(ticket)
	log.Warningln("Pipeline is ", pipeline.list, "\nAnd Online is ", online.list)
}

// Rel accepts the allotted token to the process and removes it from memory; consequently, it processes the next available process
func Rel(token string) {
	log.Debugln("Trying to release token --> ", token)
	online.removeByToken(token)
	currentAvailableMemory.add(memoryRequiredPerProcess)
	go processNextTicket()
}

type semaphore struct {
	Token      string      // Stores the UUID token
	CallerChan chan string // Stores the channel that will be used to return the Token
	Seat       chan int    // This is the actual semaphore
	Type       string      // This can be used later on, to store the type of this semaphore
}

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

// =============================================================================
// 		semaphoreList struct
type semaphoreList struct {
	sync.Mutex
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
		if sem == s {
			l.list = append(l.list[:i], l.list[i+1:]...)
			break
		}
	}
}

// removeByToken removes the semaphore from the list where the token matches
func (l *semaphoreList) removeByToken(token string) {
	l.Lock()
	defer l.Unlock()
	for i, sem := range l.list {
		if sem.Token == token {
			l.list = append(l.list[:i], l.list[i+1:]...)
			break
		}
	}
}

// Returns the length of the semaphore list
func (l *semaphoreList) length() int {
	l.Lock()
	defer l.Unlock()
	length := len(l.list)
	return length
}

// getOne returns the first semaphore of the list
func (l *semaphoreList) getOne() semaphore {
	l.Lock()
	defer l.Unlock()
	return l.list[0]
}
