package rac

type Logger interface {
	Errorf(msg string, args ...interface{})
}

type nullLogger struct{}

func (n nullLogger) Errorf(msg string, args ...interface{}) {}

var _ Logger = (*nullLogger)(nil)
