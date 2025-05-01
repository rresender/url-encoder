package repository

import (
	"github.com/rresender/url-enconder/internal/model"
	"gorm.io/gorm"
)

type EncodeURLRepository interface {
	Create(encodeURL *model.EncodeURL) error
	FindByID(id string) (*model.EncodeURL, error)
	FindByOriginalURL(tenantID, originalURL string) (*model.EncodeURL, error)
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
