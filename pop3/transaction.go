package pop3

import (
	"fmt"
	"strconv"
	"strings"
)

// sendCmd is the function that send command
// without any argument. It ends with CRLF
// (\r\n). It returns if something goes wrong
// while sending cmd.
func (c *Client) sendCmd(cmd string) error {
	buf := []byte(cmd + "\r\n")
	_, err := c.Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// sendCmdWithArg function sends the POP3 command with
// argument. It returns error if sending command
// will be unsuccessful.
//
// cmd string - command that send will send
// arg string - argument which command takes
func (c Client) sendCmdWithArg(cmd string, arg string) error {
	buf := []byte(cmd + " " + arg + "\r\n")
	_, err := c.Conn.Write(buf[:])
	if err != nil {
		return err
	}
	return nil
}

// readResp reads the command's response.
// It allocates a byte array with size of 512 byte.
// Read and store the response into buf array.
// Finally, the byte array converts to string and
// return.
func (c *Client) readResp() (string, error) {
	var buf [512]byte
	r, err := c.Conn.Read(buf[:])
	if err != nil {
		return "", err
	}

	resp := string(buf[:r])
	return resp, nil
}

// readRespMultiLines reads the response that has multiple
// lines until reaching ".\r\n" character set. Each
// line is added to listResp array.
func (c Client) readRespMultiLines() ([]string, error) {
	var buf [512]byte
	var listResp []string

	r, err := c.Conn.Read(buf[:])
	if err != nil {
		return nil, err
	}
	resp := string(buf[:r])
	listResp = append(listResp, strings.Split(resp, "\r\n")...)

	return listResp, nil
}

// Stat is a TRANSACTION state command. It
// shows that how many mails are in the inbox
// and size of the maildrop in octets. Stat
// takes no parameters. The response string
// starts with "+OK" and continues with the
// number of messages and size of the maildrop.
// separated by space.
// Response type:
// 		+OK xx yy
// Example:
// 		+OK 2 320
func (c *Client) Stat() (string, error) {
	return c.stat()
}

// stat is implementation of the Stat function.
// Sends command and receives response. Returns
// string and error. string is the response of
// the command and error is the unexpected
// situations.
func (c *Client) stat() (string, error) {
	err := c.sendCmd("STAT")
	if err != nil {
		return "", err
	}

	resp, err := c.readResp()
	if err != nil {
		return "", err
	}

	return resp, nil
}

// List returns the mail information. It can take argument
// optionally. There might be 2 different usage.
// Example-1:
// 		C: LIST
// 		S: +OK 2 messages (360 octets)
// 		S: 1 160
// 		S: 2 200
// 		S: .
//
// Example-2:
// 		C: LIST 1
// 		S: +OK 2 160
// 		C: LIST 2
// 		S: +OK 2 200
// 		C: LIST 3
// 		S: -ERR no such message, only 2 messages in maildrop
//
// You do not need to pass any argument. The function
// takes variadic parameter.
//
// msgNum ...int - variadic parameter. It indicates mail
// number that we get.
func (c *Client) List(mainNum ...int) ([]string, error) {
	// TODO Check the client is whether in TRANSACTION state
	return c.list(mainNum)
}

// list is the implementation of the List function.
// It sends the LIST command and reads the response
// that coming from the server. Returns the response
// list and error.
//
// mailNum []int - mail numbers.
func (c *Client) list(mailNum []int) ([]string, error) {
	var err error
	var msg string
	var msgList []string

	if len(mailNum) > 0 {
		err = c.sendCmdWithArg("LIST", strconv.Itoa(mailNum[0]))
	} else {
		err = c.sendCmd("LIST")
	}
	if err != nil {
		return msgList, err
	}

	if len(mailNum) == 0 {
		msgList, err = c.readRespMultiLines()
	} else {
		msg, err = c.readResp()
		msgList = append(msgList, msg)
	}
	return msgList, err
}

// Retr retrieves the mails from the inbox. It indicates
// RETR command in POP-3 protocol. It takes mailNum which
// stands for mail number. In return phase, if the mail
// is retrieved successfully, there are multiple lines.
// The first line starts with "+OK" and follows with total
// message size. On the following line, the POP3 server
// sends the entire message. Finally, command response
// ends with ".\r\n". If retrieving mail fails, the server
// returns "-ERR". The line starts with "-ERR". The function
// returns string array and error. Error is returned for
// unexpected situations like sending command or reading
// response fails. The string array contains the multiple
// line responses.
//
// mailNum string - mail-number.
func (c *Client) Retr(mailNum string) ([]string, error) {
	return c.retr(mailNum)
}

// retr function is implementation of the Retr function.
// It takes the mailNum. Firstly, sends RETR command and
// reads the response which comes from the server.
//
// mailNum string - mail-number.
func (c *Client) retr(mailNum string) ([]string, error) {
	// Send the RETR command
	err := c.sendCmdWithArg("RETR", mailNum)
	if err != nil {
		return nil, err
	}

	// Read the response
	retrResp, err := c.readRespMultiLines()
	if err != nil {
		return nil, err
	}
	return retrResp, nil
}

// Dele function deletes mail that is given as parameter.
// DELE command takes mail number and returns 2 possible
// message which are starts with "+OK" or "-ERR". POP3
// server does not actually delete the mail until POP3
// session enters the UPDATE state.
//
// mailNum string - mail number that will be deleted.
func (c *Client) Dele(mailNum string) (string, error) {
	return c.dele(mailNum)
}

// dele function is the implementation of the Dele function.
// It sends the command and reads the server response.
// Response is a string. Error is returned if unexpected
// situations happen like unsuccessful command send or
// response read.
//
// mailNum string - mail number that will be deleted.
func (c *Client) dele(mailNum string) (string, error) {
	// Send the DELE command.
	cmd := "DELE"
	err := c.sendCmdWithArg(cmd, mailNum)
	if err != nil {
		return "", err
	}

	// Read the DELE command's response.
	deleResp, err := c.readResp()
	if err != nil {
		return "", err
	}
	return deleResp, nil
}

// Noop is a command which does nothing. The POP3
// server replies with a positive response.
// Example:
// 		C: NOOP
// 		S: +OK
// It takes no argument.
func (c *Client) Noop() (string, error) {
	return c.noop()
}

// noop is implementation of the Noop function.
func (c *Client) noop() (string, error) {
	err := c.sendCmd("NOOP")
	if err != nil {
		return "", err
	}
	noop, err := c.readResp()
	if err != nil {
		return "", err
	}
	return noop, nil
}

// Rset is a command which unmark if any message
// is marked as deleted by the POP3 server. It takes
// no argument. The POP3 server replies with positive
// message as follows:
//		C: RSET
// 		S: +OK maildrop has 2 messages.
func (c *Client) Rset() (string, error) {
	return c.rset()
}

// rset is the implementation of the Rset function.
// It sends command and reads response from server.
// Returns string and error. String contains the
// message comes from the server and error is
// returned if something goes wrong while sending
// command or reading response.
func (c Client) rset() (string, error) {
	err := c.sendCmd("RSET")
	if err != nil {
		return "", err
	}

	resp, err := c.readResp()
	if err != nil {
		return "", err
	}

	return resp, nil
}

// User is the function that  authenticates the user.
// It takes username as a parameter and returns server
// response and error if something goes wrong. Firstly,
// the server checks the username and returns response
// message to client. If the message starts with "+OK",
// it means that the username is valid. If authentication
// fails, the server may respond with negative status
// indicator ("-ERR"). In such case, you may either
// issue a new authentication command or may issue the
// QUIT command. The server may return a positive response
// even though no such mailbox exits. Also, the server
// may return a negative status indicator even if the
// username is exists because the mail server does not
// permit plaintext password authentication.
// Example:
//		C: USER testUser
//		S: -ERR no such mailbox
//		...
//		C: USER validUser
//		S: +OK send PASS
//
// name string - username of the mailbox
func (c *Client) User(name string) (string, error) {
	return c.user(name)
}

// user is the implementation of the User function.
// It firstly sends the USER command and read response
// comes from the server. It returns string and error.
// String represents server's response message and
// error returns when unexpected situations if something
// goes wrong like reading response is failed.
//
// name string - username
func (c *Client) user(name string) (string, error) {
	// Send USER command
	cmd := "USER"
	err := c.sendCmdWithArg(cmd, name)
	if err != nil {
		return "", err
	}

	// Read server response
	userResp, err := c.readResp()
	if err != nil {
		return "", nil
	}

	return userResp, nil
}

// Pass is the function that sends password to POP3
// server. This function should be called after User
// function (USER command). The function takes password
// as a string and returns (string, error) pair.
// Password is a plaintext string, and it should be
// read from environment variable. Of course, you
// can declare in code and pass it to function, but
// it is not secure for privacy and security. In
// return, string contains the server response regarding
// result of the PASS command. Server response can
// start with "+OK" and "-ERR". If the password is
// true and authentication step is passed successfully,
// the server returns a message starts with "+OK". If
// the password is not true or the mail server does
// not authenticate with some reasons such as security
// level, the server returns a message contains
// negative status indicator ("-ERR"). Remaining
// message parts can be customized by the mail server.
// It depends on the mail server. For example, GMail
// returns "+OK Welcome" message.
// Example:
// 		C: USER username
// 		S: +OK send PASS
// 		C: PASS wrongPassword
//		S: -ERR Username and password not accepted.
// 		...
// 		C: USER username
// 		S: +OK send PASS
// 		C: PASS rightPassword
//		S: +OK Welcome
//
// Note: Be sure that your mail server accepts less
// secure apps. If not, give permission for less secure
// apps.
func (c *Client) Pass(password string) (string, error) {
	return c.pass(password)
}

// pass is implementation of the Pass function. It takes
// password as an argument and sends the PASS command.
// Then, reads the server's response. Password is a
// plaintext string. The function returns string and
// error. String contains the server message which either
// starts with "+OK" or "-ERR". If something goes wrong
// while sending command or reading response steps, the
// error returns.
func (c *Client) pass(password string) (string, error) {
	// Send PASS command
	cmd := "PASS"
	err := c.sendCmdWithArg(cmd, password)
	if err != nil {
		return "", err
	}

	// Read response message
	passResp, err := c.readResp()
	if err != nil {
		return "", err
	}

	return passResp, nil
}

// Top is a command which fetches message (msgNum) with n lines. To get messages
// from mail server, you need to authenticate. The response starts with "+OK"
// status indicator, and it follows the multiline mail. The server sends header
// of the message, the blank line separating the headers from the body, and then
// the number of lines of the message's body. If you request the message line
// number greater than the number of lines in the body, the mail server sends
// the entire message.
// Possible responses:
// 		+OK message follows
//		-ERR no such message
// Example:
//		C: TOP 1 10
//		S: +OK message follows
//		S: <headers of the message, blank line, and the first 10 lines of the body>
//		S: .
// 		...
//		C: TOP 100 10
//		S: -ERR no such message
//
// msgNum indicates message id starts from 1 and n is line of the message's body.
func (c *Client) Top(msgNum, n int) ([]string, error) {
	return c.top(msgNum, n)
}

// top is the implementation function of the Top function.
func (c *Client) top(msgNum, n int) ([]string, error) {
	if msgNum < 1 {
		return nil, fmt.Errorf("%s message number should be greater than 0", e)
	}
	if n < 0 {
		return nil, fmt.Errorf("%s line count cannot be negative", e)
	}
	cmd := "TOP"
	arg := fmt.Sprintf("%d %d", msgNum, n)
	err := c.sendCmdWithArg(cmd, arg)
	if err != nil {
		return nil, err
	}

	return c.readRespMultiLines()
}
