package postgres

import (
	"context"
	"demoL0/internal/domain"
	"demoL0/internal/ports/repository"
	"demoL0/internal/utils"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// orderRepo реализует интерфейс OrderRepository для PostgreSQL
type orderRepo struct {
	db *gorm.DB
}

// NewOrderRepo создает новый экземпляр PostgreSQL репозитория заказов
func NewOrderRepo(db *gorm.DB) repository.OrderRepository {
	return &orderRepo{db: db}
}

func (o *orderRepo) Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error {
	result := tx.WithContext(ctx).Create(order)
	if err := result.Error; err != nil {
		err = utils.WrapAndLog(err)
		return err
	}
	return nil
}

func (o *orderRepo) GetByID(ctx context.Context, id string) (*domain.MainOrder, error) {
	var order domain.MainOrder
	result := o.db.WithContext(ctx).Where("order_uid = ?", id).First(&order)
	log.Printf("found order with id %v: %v", id, result)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %s not found", id)
		}
		return nil, result.Error
	}
	return &order, nil
}

func (o *orderRepo) GetAll(ctx context.Context) ([]*domain.MainOrder, error) {
	var orders []*domain.MainOrder
	result := o.db.WithContext(ctx).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (o *orderRepo) Delete(tx *gorm.DB, ctx context.Context, id string) error {
	result := tx.WithContext(ctx).Delete(&domain.MainOrder{}, id)
	if result.Error != nil {
		return result.Error
	}
	fmt.Printf("Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (o *orderRepo) Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error {
	result := tx.WithContext(ctx).Model(&domain.MainOrder{}).Where("order_uid = ?", id).Updates(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
