package pop3

import (
	"os"
	"strings"
	"testing"
)

const (
	btAddr       = "mail.btopenworld.com:110"
	btTLSAddr    = "mail.btopenworld.com:995"
	gmailTLSAddr = "pop.gmail.com:995"
)

func TestUserCmd(t *testing.T) {
	pop, err := Connect(btAddr, nil, false)
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
	pop, err := Connect(btTLSAddr, nil, true)
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
	pop, err := Connect(gmailTLSAddr, nil, true)
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

// You need to save GMail username and password
// as environment variables. Environment variable
// names should be "POP3_USER" and "POP3_PASSWORD".
// NOTE: If you are working with GMail, you need to
// change security level and give permission to
// less secure apps. You can go to the following
// link and give permission to less secure apps.
// https://myaccount.google.com/lesssecureapps
func TestPassCmd(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// read username from env variable
	username := os.Getenv("POP3_USER")
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	// read password from env variable
	password := os.Getenv("POP3_PASSWORD")
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}
}

func TestStat(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv("POP3_USER")
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv("POP3_PASSWORD")
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	s, err := pop.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(s, ok) {
		t.Errorf("expected: %s, got: %s", ok, s)
	}
}

func TestStatUnauthorized(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	s, err := pop.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(s, e) {
		t.Errorf("expected: %s, got: %s", s, e)
	}
}

func TestStatErr(t *testing.T) {
	pop, err := Connect(btAddr, nil, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	s, err := pop.Stat()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(s, e) {
		t.Errorf("expected: %s, got: %s", e, s)
	}
}
