package sweetest

import (
	"fmt"
	"testing"
)

type testTest struct {
	called bool
	fatal  bool
	msg    string
	t      T
}

func (tt *testTest) reset() {
	tt.called = false
	tt.fatal = false
	tt.msg = ""
}

func (tt *testTest) Fatal(v ...interface{}) {
	tt.called = true
	tt.fatal = true
	tt.msg = v[0].(string)
}

func (tt *testTest) Error(v ...interface{}) {
	tt.called = true
	tt.fatal = false
	tt.msg = v[0].(string)
}

func (tt *testTest) Helper() {
}

func (tt *testTest) expectNothing() {
	tt.t.Helper()
	if tt.called {
		tt.t.Error("Expected no errors")
	}
	tt.reset()
}

func (tt *testTest) expectError(msg string) {
	tt.t.Helper()
	defer tt.reset()
	if !tt.called {
		tt.t.Error("Expected an error")
		return
	}
	if tt.msg != msg {
		tt.t.Error(fmt.Sprintf("Bad error: %s", tt.msg))
	}
	if tt.fatal {
		tt.t.Error("Expected a non fatal error")
	}
}

func TestAssert(t *testing.T) {
	capture := &testTest{}
	e := NewExpect(capture)

	e.Expect("a").Equal("a")
	if capture.called {
		t.Error("Should not have been called")
	}

	capture.reset()

	e.Expect("a").Equal("b")
	if !capture.called {
		t.Error("Should have been called")
	} else if capture.msg != `expect "a" to equal "b"` {
		t.Error(capture.msg)
	}

	capture.reset()

	e.Expect("a").Not().Equal("b")
	if capture.called {
		t.Error("Should not have been called: ", capture.msg)
		return
	}

	capture.reset()

	e.Expect("a").Not().Equal("a")
	if !capture.called {
		t.Error("Should have been called")
	} else if capture.msg != `expect "a" not to equal "a"` {
		t.Error(capture.msg)
	}

	capture.reset()

	e.Expect("a").Fatal().Equal("b")
	if !capture.called {
		t.Error("Should have been called")
	} else if capture.msg != `expect "a" to equal "b"` {
		t.Error(capture.msg)
	} else if !capture.fatal {
		t.Error("Should be fatal")

	}

	capture.reset()

	e.Expect("a").Name("Custom").Equal("b")
	if !capture.called {
		t.Error("Should have been called")
	} else if capture.msg != `expect Custom "a" to equal "b"` {
		t.Error(capture.msg)
	}

	capture.reset()

	e.Expect("a").Message("Custom").Equal("b")
	if !capture.called {
		t.Error("Should have been called")
	} else if capture.msg != `Custom` {
		t.Error(capture.msg)
	}
}
