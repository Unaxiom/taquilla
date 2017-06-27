package main

import (
	"time"

	"github.com/Unaxiom/taquilla"
	"github.com/Unaxiom/ulogger"
)

var log *ulogger.Logger

func init() {
	log = ulogger.New()
	taquilla.Setup(float64(5))
}

func main() {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func(count int) {
			tokenChan := make(chan string)
			log.Infoln("Requesting token for process --> ", count)
			taquilla.Req(tokenChan)
			token := <-tokenChan
			log.Warningln("Token received in example is ", token)
			<-time.After(time.Second * time.Duration(5))
			taquilla.Rel(token)
		}(i)
	}
	<-ch
}
