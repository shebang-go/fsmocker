package testdouble

import (
	"fmt"
	"testing"
)

type logger struct {
	t *testing.T
}

func CreateLogger(t *testing.T) *logger {
	l := &logger{
		t: t,
	}
	return l
}

type StructuredLogging interface {
	Module(string) *structuredLogging
	Path(string) *structuredLogging
	Operation(string) *structuredLogging
	Error(error) *structuredLogging
	Msg(format string, args ...interface{}) *structuredLogging
}

type structuredLogging struct {
	module string
	op     string
	path   string
	format string
	err    error
	msg    string
	t      *testing.T
}

func (sl *structuredLogging) Module(v string) *structuredLogging {
	sl.module = v
	return sl
}

func (sl *structuredLogging) Error(err error) *structuredLogging {
	sl.err = err
	return sl
}
func (sl *structuredLogging) Path(v string) *structuredLogging {
	sl.path = v
	return sl
}

func (sl *structuredLogging) Operation(v string) *structuredLogging {
	sl.op = v
	return sl
}

func (sl *structuredLogging) Msg(format string, args ...interface{}) *structuredLogging {
	sl.msg = fmt.Sprintf(format, args...)
	return sl
}

func (sl *structuredLogging) Done() {
	m := fmt.Sprintf("|%-10s|%-30s|%-20v|%-s", sl.op, sl.msg, sl.err, sl.path)
	if sl.t != nil {
		sl.t.Log(m)
	}
}
