package strategy

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"sync"

	"github.com/rresender/url-enconder/pkg/encoder"
)

type EncodingStrategy interface {
	GenerateID(tenantID, originalURL string) uint64
	Encode(id uint64, length int) string
}

type RandomBase36Strategy struct{}

func (r *RandomBase36Strategy) GenerateID(_, _ string) uint64 {
	return uint64(rand.Int63n(1_000_000_000))
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

func (s *SequentialBase36Strategy) GenerateID(_, _ string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	return s.counter
}

func (s *SequentialBase36Strategy) Encode(id uint64, _ int) string {
	return encoder.Base36Encode(id)
}

type TenantIDBase36Strategy struct {
}

func (t *TenantIDBase36Strategy) GenerateID(tenantID, originalURL string) uint64 {
	hash := sha256.Sum256([]byte(tenantID + "|" + originalURL))
	return binary.BigEndian.Uint64(hash[:8])
}

func (t *TenantIDBase36Strategy) Encode(id uint64, length int) string {
	return encoder.DynamicLengthEncode(id, int(length))
}
