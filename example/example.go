package main

import (
	"time"

	"sync"

	"github.com/Unaxiom/taquilla"
	"github.com/Unaxiom/ulogger"
	"github.com/twinj/uuid"
)

var log *ulogger.Logger

func init() {
	log = ulogger.New()
	taquilla.Setup(float64(2))
}

func main() {
	// ch := make(chan int)
	var wg sync.WaitGroup
	ticketType := uuid.NewV4().String()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(count int) {
			token := taquilla.Req(ticketType)
			<-time.After(time.Second * time.Duration(4))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			taquilla.Rel(token)
			wg.Done()
		}(i)
	}
	ticketType2 := uuid.NewV4().String()
	for i := 10; i < 20; i++ {
		wg.Add(1)
		go func(count int) {
			token := taquilla.Req(ticketType2)
			<-time.After(time.Second * time.Duration(4))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			taquilla.Rel(token)
			wg.Done()
		}(i)
	}
	ticketType3 := uuid.NewV4().String()
	for i := 10; i < 20; i++ {
		wg.Add(1)
		go func(count int) {
			token := taquilla.Req(ticketType3)
			<-time.After(time.Second * time.Duration(4))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			taquilla.Rel(token)
			wg.Done()
		}(i)
	}
	// <-ch
	wg.Wait()
	<-time.After(time.Second * time.Duration(10))
}
