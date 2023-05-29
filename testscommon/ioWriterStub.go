package testscommon

// IoWriterStub -
type IoWriterStub struct {
	WriteCalled func(p []byte) (n int, err error)
	CloseCalled func() error
}

// Write -
func (iws *IoWriterStub) Write(p []byte) (n int, err error) {
	if iws.WriteCalled != nil {
		return iws.WriteCalled(p)
	}

	return 0, err
}

// Close -
func (iws *IoWriterStub) Close() error {
	if iws.CloseCalled != nil {
		return iws.CloseCalled()
	}

	return nil
}
