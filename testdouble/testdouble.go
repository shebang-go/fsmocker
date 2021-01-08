package testdouble

import (
	"fmt"
	"testing"
)

type Option func(*TestDouble)

type OptionData struct {
	t *testing.T
	// enableLogging bool
	err error
}
type TestDoubler interface {
	Log(format string, args ...interface{}) *structuredLogging
	EnableLogging(t *testing.T)
}

type TestDouble struct {
	OptionData
	logger *logger
}

func NewTestDouble(opts ...Option) TestDoubler {
	t := &TestDouble{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}
func WithT(t *testing.T) Option {
	return func(td *TestDouble) {
		td.t = t
	}
}
func WithError(err error) Option {
	return func(td *TestDouble) {
		td.err = err
	}
}
func WithLogging(t *testing.T) Option {
	return func(td *TestDouble) {
		td.t = t
		// td.enableLogging = true
		td.logger = CreateLogger(t)
	}
}

func (td *TestDouble) EnableLogging(t *testing.T) {
	td.t = t
}

func (td *TestDouble) Log(format string, args ...interface{}) *structuredLogging {
	sl := &structuredLogging{
		t:      td.t,
		format: format,
		msg:    fmt.Sprintf(format, args...),
	}
	return sl
}
