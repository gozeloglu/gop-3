# gop-3

Post Office Protocol - Version 3 (POP3) Go Client

## Example

```go
func main() {
    pop, err := pop3.Connect("mail.pop3.com:110")
    if err != nil {
        log.Fatalf(err.Error())
    }
    
    fmt.Println(pop.GreetingMsg)    // Message starts with "+OK"
    fmt.Println(pop.IsAuthorized)   // true
    
    // QUIT state
    q, err := pop.Quit()
    if err != nil {
        log.Fatalf(err.Error())
    }
    fmt.Println(q)  // Prints: "QUIT"
}
```

#### References
* [RFC 1939 POP3](https://www.ietf.org/rfc/rfc1939.txt)