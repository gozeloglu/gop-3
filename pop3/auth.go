package pop3

import (
	"crypto/tls"
	"fmt"
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

	// greetingMsg keeps server response in AUTHORIZATION state.
	greetingMsg string

	// isAuthorized keeps status of AUTHORIZATION state.
	isAuthorized bool

	// isEncrypted stands for whether mail server encrypted with TLS.
	isEncrypted bool
}

const (
	// ok is successful server response's prefix
	ok = "+OK"

	// e is unsuccessful server response's prefix
	e = "-ERR"
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
		isEncrypted: true,
	}

	tlsConn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return *c, err
	}
	c.Conn = tlsConn
	c.Addr = addr

	err = c.readGreetingMsg()
	if err != nil {
		return Client{}, err
	}
	return *c, nil
}

// connectPOP3 connects to given address and returns
// a POP3 Client. This function is implementation of
// Connect() function. Reads server's response sending
// after connecting POP3 server.
func connectPOP3(addr string) (Client, error) {
	c := &Client{}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return *c, err
	}
	c.Conn = conn
	c.Addr = addr

	err = c.readGreetingMsg()
	if err != nil {
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
		return err
	}
	resp := string(buf[:r])

	// If AUTHORIZATION state fails wrt greeting
	// message, returns an error.
	if !c.isAuth(resp) {
		e := "not authorized to POP3 server"
		return fmt.Errorf(e)
	}
	c.greetingMsg = resp
	c.isAuthorized = true

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
	return strings.HasPrefix(greeting, ok)
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
	err := c.sendQuitCmd()
	if err != nil {
		return "", err
	}

	qResp, err := c.readQuitResp()
	if err != nil {
		return "", err
	}

	if isQuit(qResp) {
		c.Conn.Close()
		c.changeClientState()
	}

	return qResp, nil
}

// isQuit checks the server response after
// QUIT command. If it starts with "+OK"
// string, it is closed successfully. It
// returns boolean.
//
// resp string - response message retrieved
// after QUIT command.
func isQuit(resp string) bool {
	return strings.HasPrefix(resp, ok)
}

// sendQuitCmd sends the QUIT command to server
// as a client. The command ends with CRLF(\r\n).
// It indicates that command is terminated.
// The function returns error if occurs while
// sending command.
func (c Client) sendQuitCmd() error {
	buf := []byte("QUIT\r\n")
	_, err := c.Conn.Write(buf)
	return err
}

// readQuitResp reads the response message that comes
// from the server after sending QUIT command. If
// QUIT is done successfully, the server sends a response
// which starts with "+OK". Returns the response msg and
// error if occurs.
func (c *Client) readQuitResp() (string, error) {
	var buf [512]byte
	r, err := c.Conn.Read(buf[:])
	if err != nil {
		return "", err
	}
	resp := string(buf[:r])
	return resp, nil
}

// changeClientState changes the client's state
// after Quit command. If the Quit command is
// successful, Conn, Addr, GreetingMsg, IsAuthorized
// variables changed to nil/empty strings.
func (c *Client) changeClientState() {
	c.Conn = nil
	c.Addr = ""
	c.greetingMsg = ""
	c.isAuthorized = false
}

// GreetingMsg returns the greeting message which
// server response when connected to mail server.
// The message is returned in AUTHORIZATION state.
func (c *Client) GreetingMsg() string {
	return c.greetingMsg
}

// IsAuthorized returns the information that
// keeps the status of AUTHORIZATION state.
func (c *Client) IsAuthorized() bool {
	return c.isAuthorized
}

// IsEncrypted returns the information whether
// the server is encrypted with TLS.
func (c *Client) IsEncrypted() bool {
	return c.isEncrypted
}
