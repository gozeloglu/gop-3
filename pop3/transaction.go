package pop3

import (
	"log"
	"strconv"
)

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
		log.Println(err)
		return "", err
	}

	resp, err := c.readResp()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return resp, nil
}

// sendCmd is the function that send command
// without any argument. It ends with CRLF
// (\r\n). It returns if something goes wrong
// while sending cmd.
func (c *Client) sendCmd(cmd string) error {
	buf := []byte(cmd + "\r\n")
	_, err := c.Conn.Write(buf)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return "", err
	}

	resp := string(buf[:r])
	log.Println(resp)
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
	var msgList []string

	if len(mailNum) > 0 {
		err = c.sendListCmd(strconv.Itoa(mailNum[0]))
	} else {
		err = c.sendListCmd("")
	}
	if err != nil {
		log.Println(err)
		return msgList, err
	}

	if len(mailNum) == 0 {
		msgList, err = c.readListLines()
	} else {
		msgList, err = c.readListCmd()
	}
	return msgList, err
}

// sendListCmd sends the LIST command to POP3 server.
// The command ends with CRLF (\r\n). It can take the
// number of the message.
//
// mailNum string - mail number that listing. It can
// be empty string. If it is empty string, only
// "LIST" command is called.
func (c *Client) sendListCmd(mailNum string) error {
	buf := []byte("LIST " + mailNum + "\r\n")
	_, err := c.Conn.Write(buf)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// readListLines reads the response that has multiple
// lines until reaching ".\r\n" character set. Each
// line is added to listResp array.
func (c Client) readListLines() ([]string, error) {
	var buf [512]byte
	var listResp []string

	for string(buf[:]) != ".\r\n" {
		r, err := c.Conn.Read(buf[:])
		if err != nil {
			return nil, err
		}
		listResp = append(listResp, string(buf[:r]))
	}

	return listResp, nil
}

// readListCmd reads the "LIST" command's response.
// It has only one line. Finally, response is converted
// to string and returned it.
func (c *Client) readListCmd() ([]string, error) {
	var buf [512]byte
	var listResp []string

	r, err := c.Conn.Read(buf[:])
	if err != nil {
		log.Println(err)
		return listResp, err
	}
	resp := string(buf[:r])
	listResp = append(listResp, resp)

	return listResp, nil
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
	err := c.sendRetrCmd(mailNum)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Read the response
	retrResp, err := c.readRetrRespLines()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return retrResp, nil
}

// sendRetrCmd sends the RETR command to POP3 server.
// It returns error if writing fails.
//
// mailNum string - mail-number.
func (c *Client) sendRetrCmd(mailNum string) error {
	buf := []byte("RETR " + mailNum + "\r\n")
	_, err := c.Conn.Write(buf[:])
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// readRetrRespLines reads the response lines that come
// from the POP3 server. It returns the string array which
// contains the response lines.
func (c *Client) readRetrRespLines() ([]string, error) {
	var retrResp []string
	var buf [512]byte

	for string(buf[:]) != ".\r\n" {
		r, err := c.Conn.Read(buf[:])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		retrResp = append(retrResp, string(buf[:r]))
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
	err := c.sendDeleCmd(mailNum)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Read the DELE command's response.
	deleResp, err := c.readResp()
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(deleResp)
	return deleResp, nil
}

// sendDeleCmd function sends the DELE command with
// mail number. It returns error if sending command
// will be unsuccessful.
//
// mailNum string - mail number that will be deleted.
func (c Client) sendDeleCmd(mailNum string) error {
	buf := []byte("DELE " + mailNum + "\r\n")
	_, err := c.Conn.Write(buf[:])
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
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
		log.Println(err)
		return "", err
	}
	noop, err := c.readResp()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return "", err
	}

	resp, err := c.readResp()
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(resp)
	return resp, nil
}
