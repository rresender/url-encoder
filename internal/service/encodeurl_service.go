package service

import (
	"errors"

	"github.com/rresender/url-enconder/internal/cache"
	"github.com/rresender/url-enconder/internal/model"
	"github.com/rresender/url-enconder/internal/repository"
	"github.com/rresender/url-enconder/pkg/strategy"
	"gorm.io/gorm"
)

const defaultLength = 4

var ErrInvalidEncodingStrategy = errors.New("invalid encoding strategy")
var ErrNotFound = errors.New("not found")

type EncodeURLService interface {
	CreateEncodeURL(request *model.CreateEncodeURLRequest) (*model.EncodeURLResponse, error)
	GetOriginalURL(encodeURL string) (string, error)
}

type encodeURLService struct {
	repo       repository.EncodeURLRepository
	cache      cache.EncodeURLCache
	strategies map[string]strategy.EncodingStrategy
}

func NewEncodeURLService(repo repository.EncodeURLRepository, cache cache.EncodeURLCache) EncodeURLService {
	strategies := make(map[string]strategy.EncodingStrategy)
	strategies["random"] = &strategy.RandomBase36Strategy{}
	strategies["sequential"] = strategy.NewSequentialBase36Strategy()
	strategies["sequential_db"] = &strategy.SequentialDBBase36Strategy{
		Next: func() (uint64, error) {
			return repo.NextSequence("global")
		},
	}
	strategies["tenant"] = &strategy.TenantIDBase36Strategy{}
	return &encodeURLService{
		repo:       repo,
		cache:      cache,
		strategies: strategies,
	}
}

func (s *encodeURLService) CreateEncodeURL(request *model.CreateEncodeURLRequest) (*model.EncodeURLResponse, error) {
	cacheKey := request.TenantID + "|" + request.OriginalURL

	// Check cache first
	if cached, ok := s.cache.Get(cacheKey); ok {
		return &model.EncodeURLResponse{
			EncodeURL:   cached.ID,
			OriginalURL: cached.Original,
			TenantID:    cached.TenantID,
		}, nil
	}

	// Check repository for existing entry
	existing, err := s.repo.FindByOriginalURL(request.TenantID, request.OriginalURL)
	if err == nil {
		return s.CacheAndRespond(existing, cacheKey)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Determine encoding strategy
	strategy, exists := s.strategies[request.Strategy]
	if !exists {
		return nil, ErrInvalidEncodingStrategy
	}

	length := defaultLength
	if request.Length != nil {
		length = *request.Length
	}

	// Create new entry in repository.
	//
	// For "random" and "sequential", multiple concurrent requests can race and attempt
	// to create the same logical mapping. We rely on the DB uniqueness constraint on
	// (tenant_id, original) and retry with a different generated ID if needed.
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		id, err := strategy.GenerateID(request.TenantID, request.OriginalURL)
		if err != nil {
			return nil, err
		}
		encodeURL := strategy.Encode(id, length)

		entity := &model.EncodeURL{
			ID:       encodeURL,
			Original: request.OriginalURL,
			Strategy: request.Strategy,
			TenantID: request.TenantID,
		}

		if err := s.repo.Create(entity); err == nil {
			return s.CacheAndRespond(entity, cacheKey)
		} else {
			lastErr = err

			// If another request created the mapping, reuse it.
			existing, findErr := s.repo.FindByOriginalURL(request.TenantID, request.OriginalURL)
			if findErr == nil {
				return s.CacheAndRespond(existing, cacheKey)
			}
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return nil, findErr
			}
		}
	}

	return nil, lastErr
}

func (s *encodeURLService) CacheAndRespond(entity *model.EncodeURL, cacheKey string) (*model.EncodeURLResponse, error) {
	s.cache.Set(entity.ID, entity)
	s.cache.Set(cacheKey, entity)
	return &model.EncodeURLResponse{
		EncodeURL:   entity.ID,
		OriginalURL: entity.Original,
		TenantID:    entity.TenantID,
	}, nil
}

func (s *encodeURLService) GetOriginalURL(encodeURL string) (string, error) {
	if cached, ok := s.cache.Get(encodeURL); ok {
		return cached.Original, nil
	}

	m, err := s.repo.FindByID(encodeURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	s.cache.Set(encodeURL, m)
	return m.Original, nil
}
