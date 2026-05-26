package repository

import (
	"testing"

	"github.com/rresender/url-enconder/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNextSequence_IsMonotonic(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.AutoMigrate(&model.SequenceCounter{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := NewEncodeURLRepository(db)

	a, err := repo.NextSequence("global")
	if err != nil {
		t.Fatalf("NextSequence 1: %v", err)
	}
	b, err := repo.NextSequence("global")
	if err != nil {
		t.Fatalf("NextSequence 2: %v", err)
	}
	if b != a+1 {
		t.Fatalf("expected monotonic sequence, got a=%d b=%d", a, b)
	}
}

