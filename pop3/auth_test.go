package pop3

import (
	"strings"
	"testing"
)

var (
	c    = Client{}
	quit = "QUIT"
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

	if c.IsAuthorized != false {
		t.Errorf("Expected c.IsAuthorized %v, got: %v", false, c.IsAuthorized)
	}

	if !strings.Contains(got, quit) {
		t.Errorf("expected %s, got %s", quit, got)
	}
}
