package taquilla

import (
	"math"
	"sync"
	"testing"
	"time"

	"fmt"

	"github.com/twinj/uuid"
)

// go test -coverprofile=coverage.out
// go tool cover -html=coverage.out

// var log *ulogger.Logger

func init() {
	// log = ulogger.New()
	// Setup(float64(7))
}

// reset needs to be called when setting each test
func reset() {
	semGraph = semGraphStruct{}
	semGraph.container = make(map[string]*semCounterStruct)
	semGraph.availableTypes = []string{}

	online.list = []semaphore{}
	pipeline.list = []semaphore{}
	memUsage = memGraph{}
	memUsage.value = make(map[int64]sysMemAndOnlineSem)
}

func TestAllOnline(t *testing.T) {
	reset()
	currentAvailableMemory.set(float64(12))
	Setup(float64(2))

	done1 := make(chan int)

	for i := 0; i < 4; i++ {
		go func(i int) {
			token := Req(uuid.NewV4().String())
			<-time.After(time.Second * time.Duration(3))
			fmt.Println("Token received in example is ", token, " for count --> ", i)
			// fmt.Println("Length of online is ", online.length(), " and length of pipeline is ", pipeline.length())
			if i == 0 {
				if online.length() == 4 && pipeline.length() == 0 {
					done1 <- 1
				} else {
					done1 <- 0
				}
			} else {
				return
			}

		}(i)
	}

	returnedVal := <-done1
	if returnedVal == 0 {
		t.Error("Failed")
	}
}

func TestOneOnline(t *testing.T) {
	reset()
	currentAvailableMemory.set(float64(12))
	Setup(float64(11))

	done1 := make(chan int)

	for i := 0; i < 4; i++ {
		go func(i int) {
			token := Req(uuid.NewV4().String())
			<-time.After(time.Second * time.Duration(3))
			fmt.Println("Token received in example is ", token, " for count --> ", i, "\n\nOnline --> ", online.length(), "\nPipeline --> ", pipeline.length())
			if online.length() == 1 && pipeline.length() == 3 {
				done1 <- 1
			} else {
				done1 <- 0
			}
		}(i)
	}

	returnedVal := <-done1
	if returnedVal == 0 {
		t.Error("Failed")
	}
}

func TestDifferentTypes(t *testing.T) {
	// Pass 2 different types of ticketTypes
	reset()
	currentAvailableMemory.set(float64(12))
	Setup(float64(5))

	var wg sync.WaitGroup
	ticketType := uuid.NewV4().String()
	// Fire up all the processes
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(count int) {
			token := Req(ticketType)
			<-time.After(time.Second * time.Duration(1))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			Rel(token)
			wg.Done()
		}(i)
	}
	ticketType2 := uuid.NewV4().String()
	for i := 3; i < 5; i++ {
		wg.Add(1)
		go func(count int) {
			token := Req(ticketType2)
			<-time.After(time.Second * time.Duration(1))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			Rel(token)
			wg.Done()
		}(i)
	}
	// <-ch
	wg.Wait()

	// Check for count of semGraph of these two types and assert them --> test will pass
	if semGraph.container[ticketType].counter != 3 {
		t.Error("TicketType 1 should've had a counter as 3. Found --> ", semGraph.container[ticketType].counter)
	}
	if semGraph.container[ticketType2].counter != 2 {
		t.Error("TicketType 2 should've had a counter as 2. Found --> ", semGraph.container[ticketType2].counter)
	}

}

// TestMemoryCalculation checks if the memory calculation using matrices works as expected
func TestMemoryCalculation(t *testing.T) {
	// Fire 2 types of processes with a setup cost of 2 each
	reset()
	semaphoreMemoryUpdateTime = 1
	currentAvailableMemory.set(float64(12))
	Setup(float64(2))

	var wg sync.WaitGroup
	ticketType := uuid.NewV4().String()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(count int) {
			token := Req(ticketType)
			<-time.After(time.Second * time.Duration(2))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			Rel(token)
			wg.Done()
		}(i)
	}
	ticketType2 := uuid.NewV4().String()
	for i := 30; i < 40; i++ {
		wg.Add(1)
		go func(count int) {
			token := Req(ticketType2)
			<-time.After(time.Second * time.Duration(2))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			Rel(token)
			wg.Done()
		}(i)
	}
	ticketType3 := uuid.NewV4().String()
	for i := 60; i < 70; i++ {
		wg.Add(1)
		go func(count int) {
			token := Req(ticketType3)
			<-time.After(time.Second * time.Duration(2))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			Rel(token)
			wg.Done()
		}(i)
	}
	// <-ch
	wg.Wait()
	// <-time.After(time.Second * time.Duration(10))

	// Check for calculated memory
	// if math.Abs(semGraph.container[ticketType].avgMem-2.0) > 0.1 {
	// 	t.Error("Memory for ticket --> ", ticketType, " couldn't be calculated appropriately. Received --> ", semGraph.container[ticketType].avgMem)
	// }
	// if math.Abs(semGraph.container[ticketType2].avgMem-2.0) > 0.1 {
	// 	t.Error("Memory for ticket --> ", ticketType2, " couldn't be calculated appropriately. Received --> ", semGraph.container[ticketType2].avgMem)
	// }

	// Check for the case where there's a sudden spike in memory consumption of a process
	semaphoreMemoryUpdateTime = 20
	<-time.After(time.Duration(3) * time.Second)
	// Once the check has occurred, use the ticketType, fetch the counter value, and call the updateAvgMem method manually to match the result
	// avgMem := semGraph.container[ticketType].avgMem
	// counter := semGraph.container[ticketType].counter

	updatedAvgMem := (semGraph.container[ticketType].avgMem*float64(semGraph.container[ticketType].counter) + float64(5)) / float64(semGraph.container[ticketType].counter)

	semGraph.container[ticketType].updateAvgMem(float64(5))

	if math.Abs(updatedAvgMem-semGraph.container[ticketType].avgMem) > 0.01 {
		t.Error("Update Avg Memory for ticket --> ", ticketType, " couldn't be tested because expected value was ", updatedAvgMem, " and actual value turned out to be ", semGraph.container[ticketType].avgMem)
	}
	semGraph.container[ticketType].counter = 0
	semGraph.container[ticketType].updateAvgMem(float64(5))
	if math.Abs(updatedAvgMem-semGraph.container[ticketType].avgMem) > 0.01 {
		t.Error("Update Avg Memory for ticket --> ", ticketType, " couldn't be tested after setting counter to 0 because expected value was ", updatedAvgMem, " and actual value turned out to be ", semGraph.container[ticketType].avgMem)
	}

}

// func TestSemReq(t *testing.T) {
// 	// totalSems := 5
// 	// fmt.Println("Total sems is ", totalSems)
// 	// chanList := make([]chan string, totalSems)
// 	// for i := 0; i < totalSems; i++ {
// 	// 	ch := make(chan string)
// 	// 	chanList = append(chanList, ch)
// 	// 	fmt.Println("Created channel")
// 	// }
// 	// ch := make(chan int)
// 	// for i := 0; i < totalSems; i++ {
// 	// 	go func(i int) {
// 	// 		token := Req(uuid.NewV4().String())
// 	// 		<-time.After(time.Second * time.Duration(3))
// 	// 		fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 		chanList[i] <- token
// 	// 		Rel(token)

// 	// 	}(i)
// 	// }
// 	// t.Parallel()
// 	done1 := make(chan int)
// 	done2 := make(chan int)
// 	done3 := make(chan int)
// 	done4 := make(chan int)
// 	chanList := make([]chan int, 0)
// 	chanList = []chan int{done1, done2, done3, done4}

// 	for i := 0; i < 4; i++ {
// 		go func(i int) {
// 			// t.Parallel()
// 			token := Req(uuid.NewV4().String())
// 			<-time.After(time.Second * time.Duration(3))
// 			fmt.Println("Token received in example is ", token, " for count --> ", i)
// 			fmt.Println("Length of online is ", online.length(), " and length of pipeline is ", pipeline.length())
// 			if online.length() == 1 && pipeline.length() == 3 {
// 				ch := chanList[i]
// 				ch <- 1
// 				return
// 			}
// 			// chanList[i] <- token
// 			// Rel(token)

// 			// done1 <- 1
// 			ch := chanList[i]
// 			ch <- 1

// 		}(i)
// 	}

// 	// go func(i int) {
// 	// 	// t.Parallel()
// 	// 	token := Req(uuid.NewV4().String())
// 	// 	<-time.After(time.Second * time.Duration(3))
// 	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 	fmt.Println("Length of online is ", online.length())
// 	// 	// chanList[i] <- token
// 	// 	Rel(token)

// 	// 	done1 <- 1

// 	// }(0)
// 	// go func(i int) {
// 	// 	// t.Parallel()
// 	// 	token := Req(uuid.NewV4().String())
// 	// 	<-time.After(time.Second * time.Duration(3))
// 	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 	// chanList[i] <- token
// 	// 	Rel(token)
// 	// 	done2 <- 2
// 	// }(1)

// 	// go func(i int) {
// 	// 	// t.Parallel()
// 	// 	token := Req(uuid.NewV4().String())
// 	// 	<-time.After(time.Second * time.Duration(3))
// 	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 	// chanList[i] <- token
// 	// 	Rel(token)

// 	// }(2)

// 	// go func(i int) {
// 	// 	// t.Parallel()
// 	// 	token := Req(uuid.NewV4().String())
// 	// 	<-time.After(time.Second * time.Duration(3))
// 	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 	// chanList[i] <- token
// 	// 	Rel(token)

// 	// }(3)

// 	// go func(i int) {
// 	// 	// t.Parallel()
// 	// 	token := Req(uuid.NewV4().String())
// 	// 	<-time.After(time.Second * time.Duration(3))
// 	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
// 	// 	// chanList[i] <- token
// 	// 	Rel(token)

// 	// }(4)

// 	// <-ch
// 	// for i := 0; i < totalSems; i++ {
// 	// 	token := <-chanList[i]
// 	// 	if len(token) <= 0 {
// 	// 		t.Error("Invalid token received.")
// 	// 	} else {
// 	// 		log.Infoln("Test passed!")
// 	// 	}
// 	// }

// 	// var token string
// 	// fmt.Println("Waiting for channels")
// 	// token = <-chanList[0]
// 	// fmt.Println("Received token 1 --> ", token)
// 	// processToken(token, t)
// 	// token = <-chanList[1]
// 	// processToken(token, t)
// 	// token = <-chanList[2]
// 	// processToken(token, t)
// 	// token = <-chanList[3]
// 	// processToken(token, t)
// 	// token = <-chanList[4]
// 	// processToken(token, t)
// 	fmt.Println(<-done1)
// 	fmt.Println(<-done2)
// 	fmt.Println(<-done3)
// 	fmt.Println(<-done4)
// }

func processToken(token string, t *testing.T) {
	if len(token) <= 0 {
		t.Error("Invalid token received")
	} else {
		fmt.Println("Received token --> ", token)
	}
}
