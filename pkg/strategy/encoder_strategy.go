package strategy

import (
	"math/rand"

	"github.com/rresender/url-enconder/pkg/encoder"
)

type EncodingStrategy interface {
	GenerateID() int64
	Encode(id int64) string
}

type RandomBase62Strategy struct{}

func (r *RandomBase62Strategy) GenerateID() int64 {
	return rand.Int63n(1_000_000_000)
}

func (r *RandomBase62Strategy) Encode(id int64) string {
	return encoder.Base62Encode(id)
}

type SequentialBase62Strategy struct {
	counter int64
}

func (s *SequentialBase62Strategy) GenerateID() int64 {
	s.counter++
	return s.counter
}

func (s *SequentialBase62Strategy) Encode(id int64) string {
	return encoder.Base62Encode(id)
}
