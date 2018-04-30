package sweetest

import "fmt"

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

func (tt *testTest) Fatalf(s string, v ...interface{}) {
	tt.called = true
	tt.fatal = true
	tt.msg = fmt.Sprintf(s, v...)
}

func (tt *testTest) Errorf(s string, v ...interface{}) {
	tt.called = true
	tt.fatal = false
	tt.msg = fmt.Sprintf(s, v...)
}

func (tt *testTest) Logf(s string, v ...interface{}) {

}

func (tt *testTest) Helper() {
}

func (tt *testTest) expectNothing() {
	tt.t.Helper()
	if tt.called {
		tt.t.Errorf("Expected no errors")
	}
	tt.reset()
}

func (tt *testTest) expectError(msg string) {
	tt.t.Helper()
	defer tt.reset()
	if !tt.called {
		tt.t.Errorf("Expected an error")
		return
	}
	if tt.msg != msg {
		tt.t.Errorf("Bad error: %s", tt.msg)
	}
	if tt.fatal {
		tt.t.Errorf("Expected a non fatal error")
	}
}

func (tt *testTest) expectFatal(msg string) {
	tt.t.Helper()
	defer tt.reset()
	if !tt.called {
		tt.t.Errorf("Expected an error")
		return
	}
	if tt.msg != msg {
		tt.t.Errorf("Bad error: %s", tt.msg)
	}
	if !tt.fatal {
		tt.t.Errorf("Expected a fatal error")
	}
}
