package repository

import (
	"github.com/rresender/url-enconder/internal/model"
	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

type EncodeURLRepository interface {
	Create(encodeURL *model.EncodeURL) error
	FindByID(id string) (*model.EncodeURL, error)
	FindByOriginalURL(tenantID, originalURL string) (*model.EncodeURL, error)
	NextSequence(key string) (uint64, error)
}

type encodeURLRepository struct {
	db *gorm.DB
}

func NewEncodeURLRepository(db *gorm.DB) EncodeURLRepository {
	return &encodeURLRepository{db: db}
}

func (r *encodeURLRepository) Create(encodeURL *model.EncodeURL) error {
	return r.db.Create(encodeURL).Error
}

func (r *encodeURLRepository) FindByID(id string) (*model.EncodeURL, error) {
	var encodeURL model.EncodeURL
	err := r.db.Where("id = ?", id).First(&encodeURL).Error
	return &encodeURL, err
}

func (r *encodeURLRepository) FindByOriginalURL(tenantID, originalURL string) (*model.EncodeURL, error) {
	var shortURL model.EncodeURL
	err := r.db.Where("tenant_id = ? AND original = ?", tenantID, originalURL).First(&shortURL).Error
	return &shortURL, err
}

func (r *encodeURLRepository) NextSequence(key string) (uint64, error) {
	var next uint64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var counter model.SequenceCounter

		// Best-effort row locking where supported (e.g. Postgres).
		// SQLite will effectively serialize writes at the DB level anyway.
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("key = ?", key).
			First(&counter).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				counter = model.SequenceCounter{
					Key:   key,
					Value: 238328,
				}
				if err := tx.Create(&counter).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		counter.Value++
		if err := tx.Model(&model.SequenceCounter{}).
			Where("key = ?", key).
			Update("value", counter.Value).Error; err != nil {
			return err
		}

		next = counter.Value
		return nil
	})

	return next, err
}
