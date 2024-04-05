package data

type FieldsGetter interface {
	GetShardID() uint32
	GetTransactionPool() *TransactionPool
	GetHeaderGasConsumption() *HeaderGasConsumption
	GetAlteredAccounts() map[string]*AlteredAccount
	GetNotarizedHeadersHashes() []string
	GetNumberOfShards() uint32
	GetSignersIndexes() []uint64
	GetHighestFinalBlockNonce() uint64
	GetHighestFinalBlockHash() []byte
}

type BlockDataGetter interface {
	GetShardID() uint32
	GetHeaderType() string
	GetHeaderHash() []byte
	GetBody() *Body
	GetIntraShardMiniBlocks() []*MiniBlock
	GetScheduledRootHash() []byte
	GetScheduledAccumulatedFees() []byte
	GetScheduledDeveloperFees() []byte
	GetScheduledGasProvided() uint64
	GetScheduledGasPenalized() uint64
	GetScheduledGasRefunded() uint64
}
