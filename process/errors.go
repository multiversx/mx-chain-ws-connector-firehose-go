package process

import "errors"

var errNilMarshaller = errors.New("nil marshaller provided")

var errNilOutportBlockData = errors.New("nil outport block data")

var errNilWriter = errors.New("nil writer provided")

var errNilBlockCreator = errors.New("nil block creator provided")

var errNilPublisher = errors.New("nil publisher provided")

var errNilBlocksPool = errors.New("nil blocks pool provided")

var errNilDataAggregator = errors.New("nil data aggregator provided")

// ErrWrongTypeAssertion signals that a type assertion faled
var ErrWrongTypeAssertion = errors.New("type assertion failed")
