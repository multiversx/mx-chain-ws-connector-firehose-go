package testscommon

// LoggerStub -
type LoggerStub struct {
	InfoCalled func(message string, args ...interface{})
}

// Info -
func (l *LoggerStub) Info(message string, args ...interface{}) {
	if l.InfoCalled != nil {
		l.InfoCalled(message, args)
	}
}

// IsInterfaceNil -
func (l *LoggerStub) IsInterfaceNil() bool {
	return l == nil
}
