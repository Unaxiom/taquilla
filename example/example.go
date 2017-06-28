package main

import (
	"time"

	"github.com/Unaxiom/taquilla"
	"github.com/Unaxiom/ulogger"
)

var log *ulogger.Logger

func init() {
	log = ulogger.New()
	taquilla.Setup(float64(7))
}

func main() {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func(count int) {
			tokenChan := make(chan string)
			// log.Infoln("Requesting token for process --> ", count)
			// log.Infoln("Requesting token for process --> ", i)
			taquilla.Req(tokenChan)
			token := <-tokenChan
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			<-time.After(time.Second * time.Duration(3))
			taquilla.Rel(token)
		}(i)
	}
	<-ch
}
