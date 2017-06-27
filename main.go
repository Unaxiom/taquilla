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

var log *ulogger.Logger

// memoryRequiredPerProcess stores the memory required for running each process, in megabytes
var memoryRequiredPerProcess float64

// pipeline consists of all the list of semaphores that are currently available
var pipeline []semaphore

// online consists of all the semaphores that are under execution
var online []semaphore

// pipelineMutex is the mutex that needs to be acquired before performing any operations on the pipeline array
var pipelineMutex = &sync.Mutex{}

// Setup accepts the average memory requirement for a process
func Setup(memoryRequiredInMB float64) {
	memoryRequiredPerProcess = memoryRequiredInMB
	log = ulogger.New()
	log.SetLogLevel(ulogger.DebugLevel)
	log.Infoln("Set memoryRequiredPerProcess to ", memoryRequiredPerProcess)
}

// Req accepts a string channel via which an access token is returned. The caller function can then begin its execution.
func Req(returnChan chan string) {
	log.Debugln("Requested new semaphore")
	var ticket semaphore
	ticket.Token = strings.Join(strings.Split(uuid.NewV4().String(), "-"), "")
	log.Infoln("Generated new token --> ", ticket.Token)
	ticket.CallerChan = returnChan
	ticket.Type = "" // This can be implemented in a later version
	ticket.Seat = make(chan int, 1)
	// Insert a dummy value into the seat
	ticket.Seat <- 0
	go func() {
		// Insert the actual operating value into the seat
		ticket.Seat <- 1
	}()
	// Attach this ticket to the global list of tickets
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()
	pipeline = append(pipeline, ticket)
	go processNextTicket()
}

// processNextTicket processes the next available ticket in the pipeline
func processNextTicket() {
	log.Debugln("Started processNextTicket()")
	pipelineMutex.Lock()
	defer pipelineMutex.Unlock()
	var ticket semaphore
	if len(pipeline) == 0 {
		return
	}
	ticket = pipeline[0]
	// Check here if memory is available
	if memoryRequiredPerProcess < checkAvailableMemory() {
		log.Warningln("Ticket --> ", ticket.Token, " does not have sufficient memory to process...")
		return
	}
	// Process can be run here
	ticket.CallerChan <- ticket.Token
	pipeline = pipeline[1:]
	online = append(online, ticket)
}

// Rel accepts the allotted token to the process and removes it from memory; consequently, it processes the next available process
func Rel(token string) {
	pipelineMutex.Lock()
	log.Debugln("Trying to release token --> ", token)
	defer pipelineMutex.Unlock()
	for i, ticket := range online {
		if ticket.Token == token {
			online = append(online[0:i], online[i+1:]...)
			log.Warningln("Released token --> ", ticket.Token)
			return
		}
	}
}

// checkAvailableMemory returns the available memory as a float64
func checkAvailableMemory() float64 {
	return float64(0)
}

type semaphore struct {
	Token      string      // Stores the UUID token
	CallerChan chan string // Stores the channel that will be used to return the Token
	Seat       chan int    // This is the actual semaphore
	Type       string      // This can be used later on, to store the type of this semaphore
}
