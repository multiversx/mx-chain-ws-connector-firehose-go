package hyperOutportBlock

import "github.com/golang/protobuf/ptypes/duration"

// Poll -
func (bs *Service) Poll(nonce uint64, stream serverStream, pollingInterval *duration.Duration) error {
	return bs.poll(nonce, stream, pollingInterval)
}
