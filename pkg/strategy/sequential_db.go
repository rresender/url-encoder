package strategy

import "github.com/rresender/url-enconder/pkg/encoder"

// SequentialDBBase36Strategy uses an external, durable sequence provider
// to generate monotonic IDs that survive restarts and support multiple instances.
type SequentialDBBase36Strategy struct {
	Next func() (uint64, error)
}

func (s *SequentialDBBase36Strategy) GenerateID(_, _ string) uint64 {
	id, err := s.Next()
	if err != nil {
		// Service layer will treat the subsequent create as a failure and return an error.
		// Returning 0 keeps this method signature unchanged.
		return 0
	}
	return id
}

func (s *SequentialDBBase36Strategy) Encode(id uint64, _ int) string {
	return encoder.Base36Encode(id)
}

