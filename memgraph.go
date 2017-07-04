package taquilla

import (
	"fmt"
	"sync"

	"github.com/gonum/matrix/mat64"
)

// ===========================================================================
// memGraph describes a struct that holds the sys allocated memory
type memGraph struct {
	sync.RWMutex
	value map[int64]sysMemAndOnlineSem
}

// set acquires an exclusive lock and updates the map with the current timestamp and the sys memory used
func (m *memGraph) set(timestamp int64, sysMem uint64, onlineTypes []string) {
	m.Lock()
	defer m.Unlock()
	var localvar sysMemAndOnlineSem
	localvar.sysMem = sysMem
	localvar.online = onlineTypes
	m.value[timestamp] = localvar
	// log.Errorln("\n\n================================\nMem Graph is\n", len(localvar.online), "\n\n=================================")
}

// get acquires a read only lock and returns the sys memory at a particular timestamp
func (m *memGraph) get(timestamp int64) sysMemAndOnlineSem {
	m.RLock()
	defer m.RUnlock()
	return m.value[timestamp]
}

// updateMem calculates the memory required for each semaphore by solving the N-dimensional equation formed from constant monitoring of the memory used by the app
func (m *memGraph) updateMem() {
	m.RLock()
	defer m.RUnlock()
	// Generate equation matrix here
	// The different types of semaphores are given by the number of keys of the semGraphStruct --> this will also be the number of variables in an equation --> number of columns in the matrix
	numberOfColumnsOfEqn := semGraph.len()
	numberOfRowsOfEqn := len(m.value)
	eqnMatrix := mat64.NewDense(numberOfRowsOfEqn, numberOfColumnsOfEqn, nil)
	resultMatrix := mat64.NewDense(numberOfRowsOfEqn, 1, nil)
	// eqnMatrix.Set()
	semGraph.RLock()
	defer semGraph.RUnlock()

	// Outer loop is memGraph.value keys
	var rowCount = 0
	for timestamp, sysMem := range m.value {
		fmt.Println("*******")
		// fmt.Println("Timestamp: ", timestamp, " Row: ", rowCount, " ")
		resultMatrix.SetRow(rowCount, []float64{float64(sysMem.sysMem)})
		var equationRow []float64
		for _, ticketType := range semGraph.availableTypes {
			// fmt.Println("Ticket Type --> ", ticketType)
			var semFoundCounter = 0
			for _, semType := range sysMem.online {
				if ticketType == semType {
					// fmt.Println("Found ticketType...")
					semFoundCounter++
				}
			}
			equationRow = append(equationRow, float64(semFoundCounter))

		}
		eqnMatrix.SetRow(rowCount, equationRow)
		fmt.Println("Timestamp: ", timestamp, " Row: ", rowCount, ", ", equationRow)
		// fmt.Println("Equation Row is ", equationRow)
		fmt.Println("*******")
		rowCount++
	}

	log.ErrorDump("=======================================================================================")
	log.ErrorDump("Eqn Matrix is: ", eqnMatrix, "\n\n")
	log.ErrorDump("Result Matrix is: ", resultMatrix, "\n\n")
	log.ErrorDump("=======================================================================================")
	go computeEquationSolution(eqnMatrix, resultMatrix)

}

// sysMemAndOnlineSem is a struct that stores the system memory as well as the list of online semaphores running at that particular timestamp
type sysMemAndOnlineSem struct {
	sysMem uint64
	online []string
}
