package main

import (
	"time"

	"github.com/Unaxiom/taquilla"
	"github.com/Unaxiom/ulogger"
)

var log *ulogger.Logger

func init() {
	log = ulogger.New()
	taquilla.Setup(float64(4))
}

func main() {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func(count int) {
			token := taquilla.Req()
			<-time.After(time.Second * time.Duration(3))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			taquilla.Rel(token)
		}(i)
	}
	<-ch
}
