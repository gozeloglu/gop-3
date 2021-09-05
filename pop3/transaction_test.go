package pop3

import (
	"strings"
	"testing"
)

func TestUserCmd(t *testing.T) {
	addr := "mail.btopenworld.com:110"
	pop, err := Connect(addr, nil, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	u, err := pop.User("testUser")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}
}

func TestUserCmdWithTLS(t *testing.T) {
	addr := "mail.btopenworld.com:995"
	pop, err := Connect(addr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	u, err := pop.User("testUser")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}
}

func TestUserGMail(t *testing.T) {
	addr := "pop.gmail.com:995"
	pop, err := Connect(addr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	u, err := pop.User("testUser")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}
}
