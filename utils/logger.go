package utils

// Logger is a minimal logger interface which satisfies oxy logging needs.
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type DefaultLogger struct{}

func (*DefaultLogger) Debugf(string, ...interface{}) {}
func (*DefaultLogger) Infof(string, ...interface{})  {}
func (*DefaultLogger) Warnf(string, ...interface{})  {}
func (*DefaultLogger) Errorf(string, ...interface{}) {}
func (*DefaultLogger) Fatalf(string, ...interface{}) {}
