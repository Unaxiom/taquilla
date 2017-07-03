package main

import (
	"time"

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
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func(count int) {
			token := taquilla.Req(uuid.NewV4().String())
			<-time.After(time.Second * time.Duration(3))
			log.Infoln("Token received in example is ", token, " for count --> ", count)
			taquilla.Rel(token)
		}(i)
	}
	<-ch
}
