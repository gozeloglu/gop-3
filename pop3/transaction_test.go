package pop3

import (
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

const (
	btAddr       = "mail.btopenworld.com:110"
	btTLSAddr    = "mail.btopenworld.com:995"
	gmailTLSAddr = "pop.gmail.com:995"
	userKey      = "POP3_USER"
	passwordKey  = "POP3_PASSWORD"
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
	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	// read password from env variable
	password := os.Getenv(passwordKey)
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

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
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

func TestList(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	l, err := pop.List()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(l[0], ok) {
		t.Errorf("expected: %s, got: %s", ok, l[0])
	}
}

func TestListUnauthorized(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	l, err := pop.List()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(l[0], e) {
		t.Errorf("expected:%s, got: %s", l, e)
	}
}

func TestListWithArg(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	l, err := pop.List(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(l[0], ok) {
		t.Errorf("expected: %s, got: %s", ok, l[0])
	}

	if len(l) != 1 {
		t.Errorf("expected length: %v, l's length: %v", 1, len(l))
	}
}

func TestListWithArgUnauthorized(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	l, err := pop.List(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(l[0], e) {
		t.Errorf("expected: %s, got: %s", e, l[0])
	}
}

func TestNoop(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	n, err := pop.Noop()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(n, ok) {
		t.Errorf("expected: %s, got: %s", ok, n)
	}
}

func TestRetr(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	l, err := pop.List()
	if err != nil {
		t.Errorf(err.Error())
	}

	mailNum := strings.Split(l[1], " ")[0]
	r, err := pop.Retr(mailNum)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(r[0], ok) {
		t.Errorf("expected prefix: %s, got %s message", ok, r[0])
	}
}

func TestRetrFail(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	_, err = pop.List()
	if err != nil {
		t.Errorf(err.Error())
	}

	r, err := pop.Retr(strconv.Itoa(math.MaxInt64 - 1))
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(r[0], e) {
		t.Errorf("expected prefix: %s, got %s message", ok, r[0])
	}
}

func TestDele(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	d, err := pop.Dele("1")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(d, ok) {
		t.Errorf("expected prefix: %s, got: %s", ok, d)
	}
}

func TestDeleFail(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	d, err := pop.Dele(strconv.Itoa(math.MaxInt64 - 1))
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(d, e) {
		t.Errorf("expected prefix: %s, got: %s", e, d)
	}
}

func TestRset(t *testing.T) {
	pop, err := Connect(gmailTLSAddr, nil, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	username := os.Getenv(userKey)
	u, err := pop.User(username)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(u, ok) {
		t.Errorf("expected: %s, got: %s", ok, u)
	}

	password := os.Getenv(passwordKey)
	p, err := pop.Pass(password)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(p, ok) {
		t.Errorf("expected: %s, got: %s", ok, p)
	}

	d, err := pop.Dele("1")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(d, ok) {
		t.Errorf("expected prefix: %s, got: %s", ok, d)
	}

	r, err := pop.Rset()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.HasPrefix(r, ok) {
		t.Errorf("expected prefix: %s, got: %s", ok, r)
	}
}
