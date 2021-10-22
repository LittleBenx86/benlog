package security_test

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"testing"
)

type User struct {
	Name      string
	Authority string
}

func Test_ABACMixedWithRBAC(t *testing.T) {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act, eft")
	m.AddDef("e", "e", "some(where (p.eft == allow)) && !some(where (p.eft == deny))")
	m.AddDef("m", "m", "r.sub.Authority == p.sub && regexMatch(r.obj, p.obj) && regexMatch(r.act, p.act)")

	a := fileadapter.NewAdapter("./policy_test.csv")

	enforcer, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		t.Fatal(err)
	}

	u := User{Name: "ben", Authority: "anonymous"}
	if ok, err := enforcer.Enforce(u, "/ue/v1/app", "GET"); err != nil {
		t.Logf("error %s\n", err)
	} else {
		if ok {
			t.Logf("%+v security validation ok\n", u)
		} else {
			t.Logf("%+v security validation failed\n", u)
		}
	}

	u = User{Name: "ben2", Authority: "anonymous"}
	if ok, err := enforcer.Enforce(u, "/ue/v1/metrics", "GET"); err != nil {
		t.Logf("error %s\n", err)
	} else {
		if ok {
			t.Logf("%+v security validation ok\n", u)
		} else {
			t.Logf("%+v security validation failed\n", u)
		}
	}

	u = User{Name: "ben2", Authority: "anonymous"}
	if ok, err := enforcer.Enforce(u, "/ue/v1/metrics/mem", "GET"); err != nil {
		t.Logf("error %s\n", err)
	} else {
		if ok {
			t.Logf("%+v security validation ok\n", u)
		} else {
			t.Logf("%+v security validation failed\n", u)
		}
	}

	u = User{Name: "ben2", Authority: "anonymous"}
	if ok, err := enforcer.Enforce(u, "/ue/v1/metrics", "POST"); err != nil {
		t.Logf("error %s\n", err)
	} else {
		if ok {
			t.Logf("%+v security validation ok\n", u)
		} else {
			t.Logf("%+v security validation failed\n", u)
		}
	}

	u = User{Name: "ben2", Authority: "anonymous"}
	if ok, err := enforcer.Enforce(u, "/ue/v1/metrics", "PUT"); err != nil {
		t.Logf("error %s\n", err)
	} else {
		if ok {
			t.Logf("%+v security validation ok\n", u)
		} else {
			t.Logf("%+v security validation failed\n", u)
		}
	}
}
