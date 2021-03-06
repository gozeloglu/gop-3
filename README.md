# gop-3 [![GoDoc](https://godoc.org/github.com/gozeloglu/gop-3?status.svg)](https://godoc.org/github.com/gozeloglu/gop-3) [![Go Report Card](https://goreportcard.com/badge/github.com/gozeloglu/gop-3)](https://goreportcard.com/report/github.com/gozeloglu/gop-3)  [![Release](https://img.shields.io/badge/Release-v0.2.0-blue)](https://github.com/gozeloglu/gop-3/releases/tag/v0.2.0) ![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/gozeloglu/gop-3?filename=go.mod) [![RFC 1939](https://img.shields.io/badge/Official%20Doc-RFC%201939-yellowgreen)](https://www.ietf.org/rfc/rfc1939.txt) ![LICENSE](https://img.shields.io/badge/license-MIT-green)

### Post Office Protocol - Version 3 (POP3) Go Client

GOP-3 (Go + POP-3) is a POP-3 client for Go. It has experimental purpose and it is still under
development. [RFC 1939](https://www.ietf.org/rfc/rfc1939.txt) document has been followed while developing package.

#### Features - Commands

* USER
* PASS
* STAT
* LIST
* DELE
* RETR
* NOOP
* RSET
* QUIT
* TOP

### Installation

You can download with the following command.

```shell
go get github.com/gozeloglu/gop-3
```

## Example

```go
package main

import (
	"fmt"
	"github.com/gozeloglu/gop-3"
	"log"
	"os"
)

func main() {
	pop, err := pop3.Connect("mail.pop3.com:110", nil, false)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(pop.GreetingMsg())  // Message starts with "+OK"
	fmt.Println(pop.IsAuthorized()) // true
	fmt.Println(pop.IsEncrypted())  // false

	// USER command
	username := os.Getenv("POP3_USER") // Read from env
	u, err := pop.User(username)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(u)

	// PASS command
	password := os.Getenv("POP3_PASSWORD") // Read from env
	pass, err := pop.Pass(password)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(pass)

	// STAT command
	stat, err := pop.Stat()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(stat)

	// LIST command
	l, _ := pop.List()
	if len(l) == 0 {
		fmt.Println(l)
	}

	// LIST <mail-num> command
	ll, _ := pop.List(1)
	fmt.Println(ll[0])
	
	// TOP msgNum n 
	top, _ := pop.Top(1, 10)
	fmt.Println(top)

	// DELE command
	dele, err := pop.Dele("1")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(dele)

	// RETR command 
	retr, err := pop.Retr("1")
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, m := range retr {
		fmt.Println(m)
	}

	// NOOP command
	noop, err := pop.Noop()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(noop)

	// RSET command
	rset, err := pop.Rset()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(rset)

	// QUIT state
	q, err := pop.Quit()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(q) // Prints: "QUIT"
}
```

### Run & Test

**Note:** If you run the tests, you firstly need to have a GMail account that enables POP3 connections. Also, you have to
save mail address and password in your local environment.

If you make changes, make sure that all tests are passed. You can run the tests with the following command.

```shell
go test pop3/* -v 
```

If you want to run only one test, you can type the following command.

```shell
go test pop3/* -v -run <test_function_name>
```

Example:

```shell
go test pop3/* -v -run TestStat
```

### References

* [RFC 1939 POP3](https://www.ietf.org/rfc/rfc1939.txt)

### LICENSE

[MIT](https://github.com/gozeloglu/gop-3/blob/main/LICENSE)