package sweetest

import (
	"fmt"
	"reflect"
)

type T interface {
	Fatal(...interface{})
	Error(...interface{})
	Helper()
}

type E interface {
	Expect(val interface{}) Expectation
	getT() T
}

type expectBase struct {
	t T
}

func (b expectBase) Expect(val interface{}) Expectation {
	return Expect(b.t, val)
}

func (b expectBase) getT() T {
	return b.t
}

func NewExpect(t T) E {
	return expectBase{
		t: t,
	}
}

type Expectation interface {
	Not() Expectation
	Fatal() Expectation
	Message(string, ...interface{}) Expectation
	Name(string, ...interface{}) Expectation

	Equal(b interface{}) bool
}

type expectation struct {
	t       T
	val     interface{}
	fatal   bool
	invert  bool
	message string
	name    string
}

func Expect(t T, val interface{}) Expectation {
	return &expectation{
		t:   t,
		val: val,
	}
}

func (e *expectation) getT() T {
	return e.t
}

func (e *expectation) Fatal() Expectation {
	e.fatal = true
	return e
}

func (e *expectation) Not() Expectation {
	e.invert = true
	fmt.Printf("I: %v\n", e.invert)
	return e
}

func (e *expectation) Name(name string, params ...interface{}) Expectation {
	e.name = fmt.Sprintf(name, params...)
	return e
}

func (e *expectation) Message(msg string, params ...interface{}) Expectation {
	e.message = fmt.Sprintf(msg, params...)
	return e
}

func (e *expectation) met(met bool, msg string, params ...interface{}) bool {
	e.t.Helper()
	joinMsg := ""
	if e.invert && met {
		joinMsg = " not to "
	} else if !e.invert && !met {
		joinMsg = " to "
	} else {
		return true
	}

	expMsg := e.message

	if expMsg == "" {
		expMsg = "expect "
		if e.name != "" {
			expMsg += e.name + " "
		}
		expMsg += fmt.Sprintf("%#v", e.val)
		expMsg += joinMsg
		expMsg += fmt.Sprintf(msg, params...)
	}

	if e.fatal {
		e.t.Fatal(expMsg)
	} else {
		e.t.Error(expMsg)
	}
	return false
}
func (e *expectation) Equal(expect interface{}) bool {
	e.t.Helper()
	return e.met(reflect.DeepEqual(e.val, expect), "equal %#v", expect)
}
