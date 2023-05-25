package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

type logDataProcessor struct {
	marshaller        marshal.Marshalizer
	logger            Logger
	operationHandlers map[string]func(marshalledData []byte) error
}

// NewLogDataProcessor creates a data processor able to receive data from a ws outport driver and log events
func NewLogDataProcessor(marshaller marshal.Marshalizer, logger Logger) (DataProcessor, error) {
	if check.IfNil(marshaller) {
		return nil, errNilMarshaller
	}
	if check.IfNil(logger) {
		return nil, errNilLogger
	}

	dp := &logDataProcessor{
		marshaller: marshaller,
		logger:     logger,
	}

	dp.operationHandlers = map[string]func(marshalledData []byte) error{
		outport.TopicSaveBlock:             dp.saveBlock,
		outport.TopicRevertIndexedBlock:    dp.revertIndexedBlock,
		outport.TopicSaveRoundsInfo:        dp.saveRounds,
		outport.TopicSaveValidatorsRating:  dp.saveValidatorsRating,
		outport.TopicSaveValidatorsPubKeys: dp.saveValidatorsPubKeys,
		outport.TopicSaveAccounts:          dp.saveAccounts,
		outport.TopicFinalizedBlock:        dp.finalizedBlock,
	}

	return dp, nil
}

// ProcessPayload will process the received payload, if the topic is known.
func (dp *logDataProcessor) ProcessPayload(payload []byte, topic string) error {
	operationHandler, found := dp.operationHandlers[topic]
	if !found {
		return fmt.Errorf("%w, operation type for topic = %s, received data = %s",
			errInvalidOperationType, topic, string(payload))
	}

	return operationHandler(payload)
}

func (dp *logDataProcessor) saveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := dp.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	if outportBlock.BlockData == nil {
		return errNilOutportBlockData
	}

	dp.logger.Info("received payload", "topic", outport.TopicSaveBlock, "hash", outportBlock.BlockData.HeaderHash)

	return nil
}

func (dp *logDataProcessor) revertIndexedBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := dp.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicRevertIndexedBlock)

	return nil
}

func (dp *logDataProcessor) saveRounds(marshalledData []byte) error {
	roundsInfo := &outport.RoundsInfo{}
	err := dp.marshaller.Unmarshal(roundsInfo, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicSaveRoundsInfo)

	return nil
}

func (dp *logDataProcessor) saveValidatorsRating(marshalledData []byte) error {
	ratingData := &outport.ValidatorsRating{}
	err := dp.marshaller.Unmarshal(ratingData, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicSaveValidatorsRating)

	return nil
}

func (dp *logDataProcessor) saveValidatorsPubKeys(marshalledData []byte) error {
	validatorsPubKeys := &outport.ValidatorsPubKeys{}
	err := dp.marshaller.Unmarshal(validatorsPubKeys, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicSaveValidatorsPubKeys)

	return nil
}

func (dp *logDataProcessor) saveAccounts(marshalledData []byte) error {
	accounts := &outport.Accounts{}
	err := dp.marshaller.Unmarshal(accounts, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicSaveAccounts)

	return nil
}

func (dp *logDataProcessor) finalizedBlock(marshalledData []byte) error {
	finalizedBlock := &outport.FinalizedBlock{}
	err := dp.marshaller.Unmarshal(finalizedBlock, marshalledData)
	if err != nil {
		return err
	}

	dp.logger.Info("received payload", "topic", outport.TopicFinalizedBlock)

	return nil
}

// Close will signal via a log that the data processor is closed
func (dp *logDataProcessor) Close() error {
	dp.logger.Info("data processor closed")
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *logDataProcessor) IsInterfaceNil() bool {
	return dp == nil
}
