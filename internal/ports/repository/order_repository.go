package repository

import (
	"context"
	"demoL0/internal/domain"

	"gorm.io/gorm"
)

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error
	GetByID(ctx context.Context, id string) (*domain.MainOrder, error)
	GetAll(ctx context.Context) ([]*domain.MainOrder, error)
	Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error
	Delete(tx *gorm.DB, ctx context.Context, id string) error
}
