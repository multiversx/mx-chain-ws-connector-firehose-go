package process

import "errors"

// ErrNilMarshaller signals that a nil marshaller was provided
var ErrNilMarshaller = errors.New("nil marshaller provided")

// ErrNilOutportBlockData signals that a nil outport block data was provided
var ErrNilOutportBlockData = errors.New("nil outport block data")

// ErrNilWriter signals that a nil write was provided
var ErrNilWriter = errors.New("nil writer provided")

// ErrNilBlockCreator signals that a nil block creator was provided
var ErrNilBlockCreator = errors.New("nil block creator")

// ErrNilPublisher signals that a nil publisher was provided
var ErrNilPublisher = errors.New("nil publisher")

// ErrNilHyperBlocksPool signals that a nil hyper blocks pool was provided
var ErrNilHyperBlocksPool = errors.New("nil hyper blocks pool")

// ErrNilOutportBlocksConverter signals that a nil blocks pool was provided
var ErrNilOutportBlocksConverter = errors.New("nil outport blocks converter")

// ErrNilDataAggregator signals that a nil data aggregator was provided
var ErrNilDataAggregator = errors.New("nil data aggregator provided")

// ErrWrongTypeAssertion signals that a type assertion faled
var ErrWrongTypeAssertion = errors.New("type assertion failed")

// ErrNotSupportedDBType signals that a not supported db type has been provided
var ErrNotSupportedDBType = errors.New("not supported db type")

// ErrInvalidNumberOfPersisters signals that an invalid number of persisters has been provided
var ErrInvalidNumberOfPersisters = errors.New("invalid number of persisters")

// ErrInvalidFilePath signals that an invalid file path has been provided
var ErrInvalidFilePath = errors.New("invalid file path provided")

// ErrNilPruningStorer signals that a nil pruning storer was provide
var ErrNilPruningStorer = errors.New("nil pruning storer")

// ErrNilCacher signals that a nil cacher was provided
var ErrNilCacher = errors.New("nil cacher")

// ErrNilGRPCBlocksHandler singals that nil blocks handler was provided
var ErrNilGRPCBlocksHandler = errors.New("nil gRPC blocks handler")

// ErrInvalidOutportBlock signals that an invalid outport block was provided
var ErrInvalidOutportBlock = errors.New("invalid outport block provided")

// ErrFailedToPutBlockDataToPool signals that it cannot put block data into pool
var ErrFailedToPutBlockDataToPool = errors.New("failed to put block data")

// ErrNilDataPool signals that a nil data pool was provided
var ErrNilDataPool = errors.New("nil data pool")

// ErrNilHyperOutportBlock signals that a nil hyper outport block was provided
var ErrNilHyperOutportBlock = errors.New("nil hyper outport block")

// ErrNilBlockServiceContext signal that the context provided for the block service is nil.
var ErrNilBlockServiceContext = errors.New("nil block service context")

// ErrInvalidValue signals that an invalid value was provded
var ErrInvalidValue = errors.New("invalid value provided")
