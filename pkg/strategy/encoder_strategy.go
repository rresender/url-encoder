package strategy

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"sync"

	"github.com/rresender/url-enconder/pkg/encoder"
)

type EncodingStrategy interface {
	GenerateID(tenantID, originalURL string) (uint64, error)
	Encode(id uint64, length int) string
}

type RandomBase36Strategy struct{}

func (r *RandomBase36Strategy) GenerateID(_, _ string) (uint64, error) {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b[:]), nil
}

func (r *RandomBase36Strategy) Encode(id uint64, _ int) string {
	return encoder.Base36Encode(id)
}

type SequentialBase36Strategy struct {
	counter uint64
	mu      sync.Mutex
}

func NewSequentialBase36Strategy() *SequentialBase36Strategy {
	return &SequentialBase36Strategy{
		counter: 238328,
	}
}

func (s *SequentialBase36Strategy) GenerateID(_, _ string) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	return s.counter, nil
}

func (s *SequentialBase36Strategy) Encode(id uint64, _ int) string {
	return encoder.Base36Encode(id)
}

type TenantIDBase36Strategy struct {
}

func (t *TenantIDBase36Strategy) GenerateID(tenantID, originalURL string) (uint64, error) {
	hash := sha256.Sum256([]byte(tenantID + "|" + originalURL))
	return binary.BigEndian.Uint64(hash[:8]), nil
}

func (t *TenantIDBase36Strategy) Encode(id uint64, length int) string {
	return encoder.DynamicLengthEncode(id, int(length))
}
