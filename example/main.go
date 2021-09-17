package main

import (
	"fmt"
	"github.com/gozeloglu/gop-3/pop3"
	"log"
	"os"
)

const (
	userKey = "POP3_USER"
	passwordKey = "POP3_PASSWORD"
)

func main() {
	pop, err := pop3.Connect("pop.gmail.com:995", nil, true)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Print(pop.GreetingMsg())	// Message starts with "+OK"
	fmt.Println(pop.IsAuthorized())	// true
	fmt.Println(pop.IsEncrypted())	// true

	username := os.Getenv(userKey)
	u, err := pop.User(username)	// USER command
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(u)		// Starts with "+OK"

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)	// PASS command
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(p)		// Starts with "+OK"

	stat, _ := pop.Stat()	// STAT command
	fmt.Println(stat)		// Starts with "+OK". Returns total msg number and size.

	list, _ := pop.List()	// LIST command
	fmt.Println(list)		// Array of message number and size

	list, _ = pop.List(1)	// LIST <arg> command
	fmt.Println(list)				// 1st message size

	n, _ := pop.Noop()	// NOOP command
	fmt.Println(n)		// +OK
}
