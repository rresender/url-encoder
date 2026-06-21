package strategy

import "github.com/rresender/url-enconder/pkg/encoder"

// SequentialDBBase36Strategy uses an external, durable sequence provider
// to generate monotonic IDs that survive restarts and support multiple instances.
type SequentialDBBase36Strategy struct {
	Next func() (uint64, error)
}

func (s *SequentialDBBase36Strategy) GenerateID(_, _ string) (uint64, error) {
	return s.Next()
}

func (s *SequentialDBBase36Strategy) Encode(id uint64, _ int) string {
	return encoder.Base36Encode(id)
}

