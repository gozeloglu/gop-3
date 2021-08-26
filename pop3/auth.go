package pop3

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

// Client is POP3 client. Keeps the net.Conn, Addr of the POP3
// server's address, GreetingMsg of the POP3 server when connection
// established, and IsAuthorized boolean value when AUTHORIZATION
// state completed.
type Client struct {
	// Conn is connection for POP3 clients.
	Conn net.Conn

	// Addr is POP3 address.
	Addr string

	// GreetingMsg keeps server response in AUTHORIZATION state.
	GreetingMsg string

	// IsAuthorized keeps status of AUTHORIZATION state.
	IsAuthorized bool

	// IsEncrypted stands for whether mail server encrypted with TLS.
	IsEncrypted bool
}

const (
	// OK is successful server response's prefix
	OK = "+OK"

	// ERR is unsuccessful server response's prefix
	ERR = "-ERR"
)

// Connect create and make a connection with POP3
// server. Takes only address of the POP3 server and
// returns Client and error.
//
// addr string - POP3 mail server address. It contains
// host and port number.
// POP3 default port: 110
// POP3 (with TLS) default port: 995
// tlsConf *tls.Config - TLS configuration for POP3 servers.
// You can pass <nil> if there is no configuration.
// isEncryptedTLS bool - Indicates that POP3 server whether
// is encrypted. You can pass false the server is not
// encrypted.
func Connect(addr string, tlsConf *tls.Config, isEncryptedTLS bool) (Client, error) {
	if isEncryptedTLS {
		return connectPOP3TLS(addr, tlsConf)
	}
	return connectPOP3(addr)
}

// connectPOP3TLS connects to given address and returns
// a POP3 (encrypted with TLS) Client. This function is specialized
// for TLS encrypted POP3 servers (995 port).
//
// addr string - POP3 server address.
// config *tls.Config -  TLS configuration for POP3 server.
func connectPOP3TLS(addr string, config *tls.Config) (Client, error) {
	c := &Client{
		Conn:         nil,
		Addr:         "",
		GreetingMsg:  "",
		IsAuthorized: false,
		IsEncrypted:  true,
	}

	tlsConn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Println(err)
		return *c, err
	}
	c.Conn = tlsConn
	c.Addr = addr

	err = c.readGreetingMsg()
	if err != nil {
		log.Println(err)
		return Client{}, err
	}
	return *c, nil
}

// connectPOP3 connects to given address and returns
// a POP3 Client. This function is implementation of
// Connect() function. Reads server's response sending
// after connecting POP3 server.
func connectPOP3(addr string) (Client, error) {
	c := &Client{
		Conn:         nil,
		Addr:         "",
		GreetingMsg:  "",
		IsAuthorized: false,
		IsEncrypted:  false,
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return *c, err
	}
	c.Conn = conn
	c.Addr = addr

	err = c.readGreetingMsg()
	if err != nil {
		log.Println(err)
		return Client{}, err
	}

	return *c, nil
}

// readGreetingMsg reads the server response
// in AUTHORIZATION step. It starts with "+OK"
// string if it is successful. Response message
// keeps in Client.
// Returns error if reading or response message
// fails.
func (c *Client) readGreetingMsg() error {
	// buffer for reading server's response
	// message in AUTHORIZATION state.
	// Response message length may be up to
	// 512 characters.
	var buf [512]byte

	r, err := c.Conn.Read(buf[:])
	if err != nil {
		log.Println(err)
		return err
	}
	resp := string(buf[:r])

	// If AUTHORIZATION state fails wrt greeting
	// message, returns an error.
	if !c.isAuth(resp) {
		e := "not authorized to POP3 server"
		log.Println(e)
		return fmt.Errorf(e)
	}
	c.GreetingMsg = resp
	c.IsAuthorized = true

	log.Println(resp)
	return nil
}

// isAuth checks the greeting messages which comes
// from the server after connecting to given address.
// If the greeting message starts with "+OK" string,
// we can make sure that, the server is POP3 server.
// It returns bool value.
//
// greeting string - greeting message comes from the
// server.
func (c *Client) isAuth(greeting string) bool {
	return strings.HasPrefix(greeting, OK)
}

// Quit closes the POP3 connection with POP3
// server. It just sends "QUIT" command and get
// response from the server. The Quit function
// returns server response and error. Server
// response may start with "+OK" or "-ERR".
func (c *Client) Quit() (string, error) {
	return c.quit()
}

// quit is implementation of the Quit()
// function. Sends "QUIT\r\n" command to POP3
// server. Closes Conn if server response
// contains "+OK".
func (c *Client) quit() (string, error) {
	buf := []byte("QUIT\r\n")
	w, err := c.Conn.Write(buf)
	if err != nil {
		log.Println(err)
		return "", err
	}

	resp := string(buf[:w])

	// Close the connection and change the
	// state of the client.
	if isQuit(resp) {
		defer c.Conn.Close()
		defer c.changeClientState()
	}

	log.Println(resp)
	return resp, nil
}

// isQuit checks the server response after
// QUIT command. If it starts with "+OK"
// string, it is closed successfully. It
// returns boolean.
//
// resp string - response message retrieved
// after QUIT command.
func isQuit(resp string) bool {
	return strings.HasPrefix(resp, OK)
}

// changeClientState changes the client's state
// after Quit command. If the Quit command is
// successful, Conn, Addr, GreetingMsg, IsAuthorized
// variables changed to nil/empty strings.
func (c *Client) changeClientState() {
	c.Conn = nil
	c.Addr = ""
	c.GreetingMsg = ""
	c.IsAuthorized = false

}
