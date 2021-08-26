package pop3

import (
	"strings"
	"testing"
)

var (
	c    = Client{}
	ok = "+OK"
)

func TestIsAuth(t *testing.T) {
	resp := "+OK Hello POP3 Server"
	auth := c.isAuth(resp)

	if !auth {
		t.Errorf("Expected: %v, got: %v.", true, auth)
	}
}

func TestIsAuthFalse(t *testing.T) {
	resp := "-ERR Some problem"
	auth := c.isAuth(resp)

	if auth {
		t.Errorf("Expected: %v, got: %v.", false, auth)
	}
}

func TestConnect(t *testing.T) {
	addr := "mail.btopenworld.com:110"
	pop, err := Connect(addr, nil, false)

	if pop.Conn == nil {
		t.Errorf("c.Conn is nil.")
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if !pop.IsAuthorized {
		t.Errorf("Expected: %v, got: %v", true, pop.IsAuthorized)
	}

	if pop.Addr != addr {
		t.Errorf("Expected: %s, got: %s", addr, pop.Addr)
	}
}

func TestConnectTLS(t *testing.T) {
	addr := "mail.btopenworld.com:995"
	pop, err := Connect(addr, nil, true)

	if pop.Conn == nil {
		t.Errorf("c.Conn is nil.")
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if !pop.IsAuthorized {
		t.Errorf("Expected: %v, got: %v", true, pop.IsAuthorized)
	}

	if pop.Addr != addr {
		t.Errorf("Expected: %s, got: %s", addr, pop.Addr)
	}
}

func TestClient_Quit(t *testing.T) {
	addr := "mail.btopenworld.com:110"
	pop, err := Connect(addr, nil, false)

	if pop.Conn == nil {
		t.Errorf("c.Conn is nil.")
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	got, err := pop.Quit()
	if err != nil {
		t.Errorf(err.Error())
	}

	if pop.IsAuthorized != false {
		t.Errorf("Expected c.IsAuthorized %v, got: %v", false, pop.IsAuthorized)
	}

	if !strings.Contains(got, ok) {
		t.Errorf("expected %s, got %s", ok, got)
	}
}

func TestClientTLS_Quit(t *testing.T) {
	addr := "mail.btopenworld.com:995"
	popTLS, err := Connect(addr, nil, true)

	if popTLS.Conn == nil {
		t.Errorf(err.Error())
	}

	got, err := popTLS.Quit()
	if err != nil {
		t.Errorf(err.Error())
	}

	if popTLS.IsAuthorized != false {
		t.Errorf("expected popTLS.IsAuthorized: %v, got: %v", false, popTLS.IsAuthorized)
	}

	if !strings.Contains(got, ok) {
		t.Errorf("expected: %s, got: %s", ok, got)
	}
}
