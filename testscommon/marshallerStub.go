package testscommon

// MarshallerStub -
type MarshallerStub struct {
	MarshalCalled   func(obj interface{}) ([]byte, error)
	UnmarshalCalled func(obj interface{}, buff []byte) error
}

// Marshal -
func (ms *MarshallerStub) Marshal(obj interface{}) ([]byte, error) {
	if ms.MarshalCalled != nil {
		return ms.MarshalCalled(obj)
	}
	return make([]byte, 0), nil
}

// Unmarshal -
func (ms *MarshallerStub) Unmarshal(obj interface{}, buff []byte) error {
	if ms.UnmarshalCalled != nil {
		return ms.UnmarshalCalled(obj, buff)
	}
	return nil
}

// IsInterfaceNil -
func (ms *MarshallerStub) IsInterfaceNil() bool {
	return ms == nil
}
