package dataPool

import "errors"

// ErrNilDataPool signals that a nil data pool was provided
var ErrNilDataPool = errors.New("nil data pool provided")

// ErrNilMetaOutportBlock signals that a nil meta outport block was provided
var ErrNilMetaOutportBlock = errors.New("nil meta outport block was provided")

// ErrFailedToPutBlockDataToPool signals that it cannot put block data into pool
var ErrFailedToPutBlockDataToPool = errors.New("failed to put block data")
