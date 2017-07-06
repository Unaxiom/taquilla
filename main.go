package taquilla

import (
	"strings"
	"time"

	"runtime"

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

// memUsageGraph stores all the time
var memUsage memGraph

// semGraph stores a map between the semaphore types and the corresponding semCounter
var semGraph semGraphStruct

// memorySamplingTime is the time interval after which the memory used by the application is calculated
const memorySamplingTime = 1

// semaphoreMemoryUpdateTime is the time interval after which a semaphore's memory usage is calculated
const semaphoreMemoryUpdateTime = 5

func init() {
	log = ulogger.New()
	// log.SetLogLevel(ulogger.DebugLevel)
	log.SetLogLevel(ulogger.ErrorLevel)
	memUsage.value = make(map[int64]sysMemAndOnlineSem)
	semGraph.container = make(map[string]*semCounterStruct)
	go countMemoryUsage()
	go updateSemaphoreMemUsage()
}

// countMemoryUsage keeps reading the current memory usage
func countMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// log.Errorln("==================\n\nHeapAlloc: ", m.HeapAlloc, "\nHeapSys: ", m.HeapSys, "\nHeapObjects: ", m.HeapObjects, "\nStackSize: ", m.StackSys, "\n\n====================================")
	log.Debugln("Sys Mem: ", m.Sys)
	// go memUsage.set(time.Now().Unix(), m.Sys, online.getAll())
	go memUsage.set(time.Now().Unix(), uint64(currentAvailableMemory.total-currentAvailableMemory.getFree()), online.getAllTypes())
	<-time.After(time.Second * time.Duration(memorySamplingTime))
	go countMemoryUsage()
}

// updateSemaphoreMemUsage keeps calculating the memory used by each semaphore
func updateSemaphoreMemUsage() {
	<-time.After(time.Second * time.Duration(semaphoreMemoryUpdateTime))
	go memUsage.updateMem()
	go updateSemaphoreMemUsage()
}

// Setup accepts the average memory requirement for a process
func Setup(memoryRequiredInMB float64) {
	memoryRequiredPerProcess = memoryRequiredInMB
	log.Debugln("Set memoryRequiredPerProcess to ", memoryRequiredPerProcess)
	currentAvailableMemory.set(float64(12))
}

// Req accepts a title for the semaphore and an access token is returned. The caller function can then begin its execution.
func Req(ticketType string) string {
	log.Debugln("Requested new semaphore")
	var ticket semaphore
	ticket.ReqTime = time.Now().Unix()
	ticket.Token = strings.Join(strings.Split(uuid.NewV4().String(), "-"), "")
	ticket.Type = ticketType
	// ticket.ReqSemList = online.getAll()
	log.Debugln("Generated new token --> ", ticket.Token)
	ticketChan := make(chan string)
	ticket.CallerChan = ticketChan
	go updateSemTypeCounter(ticketType)
	pipeline.append(ticket)
	go processNextTicket(ticket)
	<-ticketChan
	return ticket.Token
}

// processNextTicket processes the next available ticket in the pipeline
func processNextTicket(ticket semaphore) {
	log.Debugln("Started processNextTicket()")
	// Check here if memory is available
	if memoryRequiredPerProcess > currentAvailableMemory.getFree() {
		log.Warningln("Ticket --> ", ticket.Token, " does not have sufficient memory to process...")
		return
	}
	// Process can be run here
	ticket.CallerChan <- ticket.Token
	pipeline.remove(ticket)
	online.appendAndReduceMemory(ticket)
	log.Debugln("Pipeline is ", pipeline.list, "\nAnd Online is ", online.list)
}

// Rel accepts the allotted token to the process and removes it from memory; consequently, it processes the next available process
func Rel(token string) {
	log.Debugln("Trying to release token --> ", token)

	ticket := online.removeByToken(token)
	ticket.RelTime = time.Now().Unix()
	// ticket.RelSemList = online.getAll()

	newTicket, presence := pipeline.getOne()
	if !presence {
		log.Warningln("No elements found in the pipeline!")
		return
	}
	log.Infoln("New Ticket is ", newTicket)
	go processNextTicket(newTicket)
}

// updateSemTypeCounter accepts the ticketType of the semaphore, and increments the count of the total number of semaphores created of this type
func updateSemTypeCounter(ticketType string) {
	// Check if the ticketType exists already
	semGraph.set(ticketType)
}

// // updateSemaphoreCharacteristics accepts a semaphore, calculates the average memory used by this semaphore, and updates
// func updateSemaphoreCharacteristics(ticket semaphore) {

// }
