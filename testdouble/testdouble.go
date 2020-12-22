package testdouble

import (
	"testing"
)

type Option func(*OptionData)

type OptionData struct {
	t             *testing.T
	enableLogging bool
	err           error
}

type TestDouble struct {
	OptionData
}

func NewTestDouble(opts ...Option) *TestDouble {
	t := &TestDouble{}
	for _, opt := range opts {
		opt(&t.OptionData)
	}
	return t
}
func WithT(t *testing.T) Option {
	return func(s *OptionData) {
		s.t = t
	}
}
func WithError(err error) Option {
	return func(s *OptionData) {
		s.err = err
	}
}
func WithLogging(t *testing.T) Option {
	return func(s *OptionData) {
		s.t = t
		s.enableLogging = true
	}
}

func (td *TestDouble) Log(format string, args ...interface{}) {
	if td.t != nil {
		td.t.Logf(format, args...)
	}
}
