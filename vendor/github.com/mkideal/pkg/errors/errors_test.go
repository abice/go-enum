package errors

import (
	"testing"
)

func TestError(t *testing.T) {
	const e = Error("error")
	if e.Error() != "error" {
		t.Errorf("want `error', got `%s'", e.Error())
	}
}

func TestTraceError(t *testing.T) {
	text := "throw error"
	e := Throw(text)
	t.Logf("e.Error=`%v'", e)
}

func TestWrappedError(t *testing.T) {
	e := Error("error")
	en := Wrap(Wrap(Wrap(e)))
	core := Core(en)
	if core != e {
		t.Errorf("core(%v) of en not equal to %v", core, e)
	}
}
