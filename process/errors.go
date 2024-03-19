package process

import "errors"

// ErrNilMarshaller signals that a nil marshaller was provided
var ErrNilMarshaller = errors.New("nil marshaller provided")

// ErrNilOutportBlockData signals that a nil outport block data was provided
var ErrNilOutportBlockData = errors.New("nil outport block data")

// ErrNilWriter signals that a nil write was provided
var ErrNilWriter = errors.New("nil writer provided")

// ErrNilBlockCreator signals that a nil block creator was provided
var ErrNilBlockCreator = errors.New("nil block creator provided")

// ErrNilPublisher signals that a nil publisher was provided
var ErrNilPublisher = errors.New("nil publisher provided")

// ErrNilBlocksPool signals that a nil blocks pool was provided
var ErrNilBlocksPool = errors.New("nil blocks pool provided")

// ErrNilDataAggregator signals that a nil data aggregator was provided
var ErrNilDataAggregator = errors.New("nil data aggregator provided")

// ErrWrongTypeAssertion signals that a type assertion faled
var ErrWrongTypeAssertion = errors.New("type assertion failed")
