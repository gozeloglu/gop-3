package main

import (
	"fmt"
	"github.com/gozeloglu/gop-3/pop3"
	"log"
)

func main() {
	pop, err := pop3.Connect("pop3.mail.com:110", nil, false)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Print(pop.GreetingMsg)
	fmt.Println(pop.IsAuthorized)
}
