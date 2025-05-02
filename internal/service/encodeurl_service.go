package service

import (
	"errors"

	"github.com/rresender/url-enconder/internal/cache"
	"github.com/rresender/url-enconder/internal/model"
	"github.com/rresender/url-enconder/internal/repository"
	"github.com/rresender/url-enconder/pkg/strategy"
)

const defaultLength = 4

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

	// Determine encoding strategy
	strategy, exists := s.strategies[request.Strategy]
	if !exists {
		return nil, errors.New("invalid encoding strategy")
	}

	id := strategy.GenerateID(request.TenantID, request.OriginalURL)

	length := defaultLength
	if request.Length != nil {
		length = *request.Length
	}
	encodeURL := strategy.Encode(id, length)

	entity := &model.EncodeURL{
		ID:       encodeURL,
		Original: request.OriginalURL,
		Strategy: request.Strategy,
		TenantID: request.TenantID,
	}

	// Create new entry in repository
	if err := s.repo.Create(entity); err != nil {
		return nil, err
	}

	return s.CacheAndRespond(entity, cacheKey)
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

	model, err := s.repo.FindByID(encodeURL)

	if err != nil {
		return "", err
	}

	s.cache.Set(encodeURL, model)

	return model.Original, nil
}
