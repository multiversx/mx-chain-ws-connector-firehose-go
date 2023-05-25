package process

import "errors"

var errInvalidOperationType = errors.New("invalid/unknown operation type")

var errNilMarshaller = errors.New("nil marshaller provided")

var errNilLogger = errors.New("nil logger provided")

var errNilOutportBlockData = errors.New("nil outport block data")
