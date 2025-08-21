package repository

import (
	"context"
	"demoL0/internal/domain"
	"demoL0/internal/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

type OrderRepository interface {
	Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error
	GetByID(ctx context.Context, id string) (*domain.MainOrder, error)
	GetAll(ctx context.Context) ([]*domain.MainOrder, error)
	Delete(tx *gorm.DB, ctx context.Context, id string) error
	Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error
}

func CreateOrderRepo(orm *gorm.DB) OrderRepo {
	return OrderRepo{orm}
}

func (o *OrderRepo) Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error {
	result := tx.WithContext(ctx).Create(order)
	if err := result.Error; err != nil {
		err = utils.WrapAndLog(err)
		return err
	}
	return nil
}

func (o *OrderRepo) GetByID(ctx context.Context, id string) (*domain.MainOrder, error) {
	var order domain.MainOrder
	result := o.db.WithContext(ctx).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, result.Error
	}
	return &order, nil
}

func (o *OrderRepo) GetAll(ctx context.Context) ([]*domain.MainOrder, error) {
	var orders []*domain.MainOrder
	result := o.db.WithContext(ctx).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (o *OrderRepo) Delete(tx *gorm.DB, ctx context.Context, id string) error {
	result := tx.WithContext(ctx).Delete(&domain.MainOrder{}, id)
	if result.Error != nil {
		return result.Error
	}
	fmt.Printf("Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (o *OrderRepo) Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error {
	result := tx.WithContext(ctx).Model(&domain.MainOrder{}).Where("order_uid = ?", id).Updates(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
