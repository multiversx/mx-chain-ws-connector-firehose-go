package testscommon

import (
	"encoding/json"
	"errors"
)

// MarshallerMock -
type MarshallerMock struct {
}

// Marshal -
func (mm *MarshallerMock) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, errors.New("nil object to serialize from")
	}

	return json.Marshal(obj)
}

// Unmarshal -
func (mm *MarshallerMock) Unmarshal(obj interface{}, buff []byte) error {
	if obj == nil {
		return errors.New("nil object to serialize to")
	}

	if buff == nil {
		return errors.New("nil byte buffer to deserialize from")
	}

	if len(buff) == 0 {
		return errors.New("empty byte buffer to deserialize from")
	}

	return json.Unmarshal(buff, obj)
}

// IsInterfaceNil -
func (mm *MarshallerMock) IsInterfaceNil() bool {
	return mm == nil
}
