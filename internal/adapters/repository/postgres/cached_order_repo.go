package postgres

import (
	"context"
	"demoL0/internal/adapters/cache/inmemory"
	"demoL0/internal/domain"
	"demoL0/internal/ports/repository"
	"errors"
	"log"

	"gorm.io/gorm"
)

// cachedOrderRepository реализует кэшированный репозиторий заказов
type cachedOrderRepository struct {
	repo  repository.OrderRepository
	cache *inmemory.MapCache[string, *domain.MainOrder]
}

// NewCachedOrderRepository создает новый экземпляр кэшированного репозитория
func NewCachedOrderRepository(repo repository.OrderRepository, cache *inmemory.MapCache[string, *domain.MainOrder]) repository.OrderRepository {
	cachedRepository := &cachedOrderRepository{
		repo:  repo,
		cache: cache,
	}
	_, err := cachedRepository.GetAll(context.TODO())
	if err != nil {
		log.Println("failed to initial fill cache")
	}
	return cachedRepository
}

func (cr *cachedOrderRepository) Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error {
	if _, ok := cr.cache.Get(order.OrderUID); ok {
		return errors.New("order already exists")
	}

	err := cr.repo.Create(tx, ctx, order)
	if err != nil {
		return err
	}

	cr.cache.Set(order.OrderUID, order)
	return nil
}

func (cr *cachedOrderRepository) GetByID(ctx context.Context, id string) (*domain.MainOrder, error) {
	if order, ok := cr.cache.Get(id); ok {
		log.Printf("[cache] id: %v, order %v", id, order)
		return order, nil
	}
	byID, err := cr.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cr.cache.Set(id, byID)
	return byID, err
}

func (cr *cachedOrderRepository) GetAll(ctx context.Context) ([]*domain.MainOrder, error) {
	if res, ok := cr.cache.GetAll(); ok {
		return res, nil
	}
	all, err := cr.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, it := range all {
		cr.cache.Set(it.OrderUID, it)
	}
	return all, nil
}

func (cr *cachedOrderRepository) Delete(tx *gorm.DB, ctx context.Context, id string) error {
	err := cr.repo.Delete(tx, ctx, id)
	if err != nil {
		return err
	}
	cr.cache.Delete(id)
	return nil
}

func (cr *cachedOrderRepository) Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error {
	err := cr.repo.Update(tx, ctx, id, order)
	if err != nil {
		return err
	}
	cr.cache.Set(id, order)
	return err
}
