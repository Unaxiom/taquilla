package taquilla

import (
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
	Setup(float64(7))
}

func TestAllOnline(t *testing.T) {
	currentAvailableMemory.set(float64(12))
	Setup(float64(2))
	online.list = []semaphore{}
	pipeline.list = []semaphore{}
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
		t.Fatal("Failed")
	}
}

func TestOneOnline(t *testing.T) {
	currentAvailableMemory.set(float64(12))
	Setup(float64(11))
	online.list = []semaphore{}
	pipeline.list = []semaphore{}
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
		t.Fatal("Failed")
	}
}

func TestDifferentTypes(t *testing.T) {
	// Pass 2 different types of ticketTypes
	currentAvailableMemory.set(float64(12))
	Setup(float64(5))
	online.list = []semaphore{}
	pipeline.list = []semaphore{}
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
		t.Fatal("TicketType 1 should've had a counter as 3. Found --> ", semGraph.container[ticketType].counter)
	}
	if semGraph.container[ticketType2].counter != 2 {
		t.Fatal("TicketType 2 should've had a counter as 2. Found --> ", semGraph.container[ticketType2].counter)
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
// 	// 		t.Fatal("Invalid token received.")
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
		t.Fatal("Invalid token received")
	} else {
		fmt.Println("Received token --> ", token)
	}
}
