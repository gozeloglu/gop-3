# gop-3

#### Post Office Protocol - Version 3 (POP3) Go Client

GOP-3 (Go + POP-3) is a POP-3 client for Go. It has experimental purpose and it is still under development.

### Download

You can download with the following command.

```shell
$ go get github.com/gozeloglu/gop-3
```

## Example

```go
package main

import (
	"fmt"
	"github.com/gozeloglu/gop-3"
	"log"
)

func main() {
	pop, err := pop3.Connect("mail.pop3.com:110")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(pop.GreetingMsg)  // Message starts with "+OK"
	fmt.Println(pop.IsAuthorized) // true

	// STAT command
	stat, err := pop.Stat()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(stat)

	// LIST command
	s, _ := pop.List()
	if len(s) == 0 {
		fmt.Println(s)
	}

	// LIST <mail-num> command
	l, _ := pop.List(1)
	fmt.Println(l[0])

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

	// QUIT state
	q, err := pop.Quit()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(q) // Prints: "QUIT"
}
```

#### References

* [RFC 1939 POP3](https://www.ietf.org/rfc/rfc1939.txt)

:warning: This package is just for testing purposes. It is not development-ready package. If you use in production, be
careful with the package. All issues are welcome and you can open an issue if you face any problem.
