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
	err := c.sendStatCmd()
	if err != nil {
		log.Println(err)
		return "", err
	}

	resp, err := c.readStatResp()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return resp, nil
}

// sendStatCmd sends STAT command to server.
// It ends with CRLF (\r\n). It returns if
// something goes wrong while sending cmd.
func (c *Client) sendStatCmd() error {
	buf := []byte("STAT\r\n")
	_, err := c.Conn.Write(buf)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// readStatResp reads the STAT command's response.
// It allocates a byte array with size of 512 byte.
// Read and store the response into buf array.
// Finally, the byte array converts to string and
// return.
func (c *Client) readStatResp() (string, error) {
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
// C: LIST
// S: +OK 2 messages (360 octets)
// S: 1 160
// S: 2 200
// S: .
//
// C: LIST 1
// S: +OK 2 160
// C: LIST 2
// S: +OK 2 200
// C: LIST 3
// S: -ERR no such message, only 2 messages in maildrop
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
			return listResp, err
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
