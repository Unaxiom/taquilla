package taquilla

import (
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
	done1 := make(chan int)

	for i := 0; i < 4; i++ {
		go func(i int) {
			token := Req(uuid.NewV4().String())
			<-time.After(time.Second * time.Duration(3))
			fmt.Println("Token received in example is ", token, " for count --> ", i)
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

func TestSemReq(t *testing.T) {
	// totalSems := 5
	// fmt.Println("Total sems is ", totalSems)
	// chanList := make([]chan string, totalSems)
	// for i := 0; i < totalSems; i++ {
	// 	ch := make(chan string)
	// 	chanList = append(chanList, ch)
	// 	fmt.Println("Created channel")
	// }
	// ch := make(chan int)
	// for i := 0; i < totalSems; i++ {
	// 	go func(i int) {
	// 		token := Req(uuid.NewV4().String())
	// 		<-time.After(time.Second * time.Duration(3))
	// 		fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 		chanList[i] <- token
	// 		Rel(token)

	// 	}(i)
	// }
	// t.Parallel()
	done1 := make(chan int)
	done2 := make(chan int)
	done3 := make(chan int)
	done4 := make(chan int)
	chanList := make([]chan int, 0)
	chanList = []chan int{done1, done2, done3, done4}

	for i := 0; i < 4; i++ {
		go func(i int) {
			// t.Parallel()
			token := Req(uuid.NewV4().String())
			<-time.After(time.Second * time.Duration(3))
			fmt.Println("Token received in example is ", token, " for count --> ", i)
			fmt.Println("Length of online is ", online.length(), " and length of pipeline is ", pipeline.length())
			if online.length() == 1 && pipeline.length() == 3 {
				ch := chanList[i]
				ch <- 1
				return
			}
			// chanList[i] <- token
			// Rel(token)

			// done1 <- 1
			ch := chanList[i]
			ch <- 1

		}(i)
	}

	// go func(i int) {
	// 	// t.Parallel()
	// 	token := Req(uuid.NewV4().String())
	// 	<-time.After(time.Second * time.Duration(3))
	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 	fmt.Println("Length of online is ", online.length())
	// 	// chanList[i] <- token
	// 	Rel(token)

	// 	done1 <- 1

	// }(0)
	// go func(i int) {
	// 	// t.Parallel()
	// 	token := Req(uuid.NewV4().String())
	// 	<-time.After(time.Second * time.Duration(3))
	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 	// chanList[i] <- token
	// 	Rel(token)
	// 	done2 <- 2
	// }(1)

	// go func(i int) {
	// 	// t.Parallel()
	// 	token := Req(uuid.NewV4().String())
	// 	<-time.After(time.Second * time.Duration(3))
	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 	// chanList[i] <- token
	// 	Rel(token)

	// }(2)

	// go func(i int) {
	// 	// t.Parallel()
	// 	token := Req(uuid.NewV4().String())
	// 	<-time.After(time.Second * time.Duration(3))
	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 	// chanList[i] <- token
	// 	Rel(token)

	// }(3)

	// go func(i int) {
	// 	// t.Parallel()
	// 	token := Req(uuid.NewV4().String())
	// 	<-time.After(time.Second * time.Duration(3))
	// 	fmt.Println("Token received in example is ", token, " for count --> ", i)
	// 	// chanList[i] <- token
	// 	Rel(token)

	// }(4)

	// <-ch
	// for i := 0; i < totalSems; i++ {
	// 	token := <-chanList[i]
	// 	if len(token) <= 0 {
	// 		t.Fatal("Invalid token received.")
	// 	} else {
	// 		log.Infoln("Test passed!")
	// 	}
	// }

	// var token string
	// fmt.Println("Waiting for channels")
	// token = <-chanList[0]
	// fmt.Println("Received token 1 --> ", token)
	// processToken(token, t)
	// token = <-chanList[1]
	// processToken(token, t)
	// token = <-chanList[2]
	// processToken(token, t)
	// token = <-chanList[3]
	// processToken(token, t)
	// token = <-chanList[4]
	// processToken(token, t)
	fmt.Println(<-done1)
	fmt.Println(<-done2)
	fmt.Println(<-done3)
	fmt.Println(<-done4)
}

func processToken(token string, t *testing.T) {
	if len(token) <= 0 {
		t.Fatal("Invalid token received")
	} else {
		fmt.Println("Received token --> ", token)
	}
}
