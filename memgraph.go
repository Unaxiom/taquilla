package taquilla

import (
	"fmt"
	"sync"

	"crypto/md5"
	"io"

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
	if len(onlineTypes) != 0 {
		m.value[timestamp] = localvar
	}
}

// get acquires a read only lock and returns the sys memory at a particular timestamp
// func (m *memGraph) get(timestamp int64) sysMemAndOnlineSem {
// 	m.RLock()
// 	defer m.RUnlock()
// 	return m.value[timestamp]
// }

// updateMem calculates the memory required for each semaphore by solving the N-dimensional equation formed from constant monitoring of the memory used by the app
func (m *memGraph) updateMem() {
	// log.Errorln("------------------------------------------ UpdateMem() has been called here ------------------------------------------")
	m.RLock()
	defer m.RUnlock()
	// Generate equation matrix here
	// The different types of semaphores are given by the number of keys of the semGraphStruct --> this will also be the number of variables in an equation --> number of columns in the matrix
	numberOfColumnsOfEqn := semGraph.len()
	// numberOfRowsOfEqn := len(m.value)
	numberOfRowsOfEqn := numberOfColumnsOfEqn
	// numberOfRowsOfEqn := len(listOfTimestampsToMonitor)
	// eqnMatrix := mat64.NewDense(numberOfRowsOfEqn, numberOfColumnsOfEqn, nil)
	// resultMatrix := mat64.NewDense(numberOfRowsOfEqn, 1, nil)
	// eqnMatrix.Set()
	semGraph.RLock()
	defer semGraph.RUnlock()

	// Outer loop is memGraph.value keys
	var rowCount = 0

	// ===========================================================================
	// Generate all the available equations and corresponding results
	equationMap := make(map[int][]float64)
	resultMap := make(map[int][]float64)

	// Store all the available ticket types in a list
	var ticketTypeList []string
	for _, ticketType := range semGraph.availableTypes {
		ticketTypeList = append(ticketTypeList, ticketType)
	}

	for _, sysMem := range m.value {
		// fmt.Println("*******")
		// fmt.Println("Timestamp: ", timestamp, " Row: ", rowCount, " ")
		// resultMatrix.SetRow(rowCount, []float64{float64(sysMem.sysMem)})
		resultMap[rowCount] = []float64{float64(sysMem.sysMem)}
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
		// eqnMatrix.SetRow(rowCount, equationRow)
		equationMap[rowCount] = equationRow
		// fmt.Println("Timestamp: ", timestamp, " Row: ", rowCount, ", ", equationRow)
		// fmt.Println("Equation Row is ", equationRow)
		// fmt.Println("*******")
		rowCount++
	}

	// Remove all the duplicates here
	localChan := make(chan int)
	go func(ch chan int) {
		equationMap = removeDuplicateValuesFromMap(equationMap)
		newResultMap := make(map[int][]float64)
		for key := range equationMap {
			newResultMap[key] = resultMap[key]
		}
		// resultMap = removeDuplicateValuesFromMap(resultMap)
		resultMap = newResultMap
		ch <- 1
	}(localChan)
	<-localChan

	// Pick up the last numberOfColumnsOfEqn elements from both the equationMap and resultMap --> this would give the exact matrices
	if len(equationMap) < numberOfColumnsOfEqn {
		return
	}
	// log.ErrorDump("ResultMap after deduplication is ", resultMap)
	// log.ErrorDump("EquationMap after deduplication is ", equationMap)
	eqnMatrix := mat64.NewDense(numberOfRowsOfEqn, numberOfColumnsOfEqn, nil)
	resultMatrix := mat64.NewDense(numberOfRowsOfEqn, 1, nil)
	var count = 0
	for key := range equationMap {
		if count >= numberOfColumnsOfEqn {
			break
		}
		resultMatrix.SetRow(count, resultMap[key])
		eqnMatrix.SetRow(count, equationMap[key])
		count++
	}

	// log.ErrorDump("=======================================================================================")
	// log.ErrorDump("Eqn Matrix is: ", eqnMatrix, "\n\n")
	// log.ErrorDump("Result Matrix is: ", resultMatrix, "\n\n")
	// log.ErrorDump("=======================================================================================")
	go computeEquationSolution(eqnMatrix, resultMatrix, ticketTypeList)

}

// sysMemAndOnlineSem is a struct that stores the system memory as well as the list of online semaphores running at that particular timestamp
type sysMemAndOnlineSem struct {
	sysMem uint64
	online []string
}

// removeDuplicateValuesFromMap accepts the address of a map and removes all the duplicate values from the map
func removeDuplicateValuesFromMap(localMap map[int][]float64) map[int][]float64 {
	stringMap := make(map[string]int)
	// Generate string representation of the inner array, and reverse store it (key becomes the value)
	for key, val := range localMap {
		hashHolder := md5.New()
		for _, fl64 := range val {
			io.WriteString(hashHolder, fmt.Sprintf("%f, ", fl64))
		}

		stringMap[fmt.Sprintf("%x", hashHolder.Sum(nil))] = key
	}

	deduplicatedMap := make(map[int][]float64)
	for _, val := range stringMap {
		deduplicatedMap[val] = localMap[val]
	}
	return deduplicatedMap
}
