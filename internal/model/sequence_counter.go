package model

// SequenceCounter stores a durable counter for ID generation.
// Key allows for future expansion (per-tenant sequences, etc).
type SequenceCounter struct {
	Key   string `gorm:"primaryKey"`
	Value uint64 `gorm:"not null"`
}

