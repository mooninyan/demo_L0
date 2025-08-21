package service

import (
	"context"
	"demoL0/internal/dto"
	"demoL0/internal/mapper"
	"demoL0/internal/repository"
	"demoL0/internal/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type OrderService struct {
	db   *gorm.DB
	repo repository.OrderRepository
}

func CreateOrderService(db *gorm.DB, repo repository.OrderRepository) *OrderService {
	return &OrderService{db, repo}
}

func (s *OrderService) Create(ctx context.Context, order *dto.MainOrder) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)
	err = tx.Error
	if err != nil {
		return err
	}

	model := mapper.MapDtoToModel(order)
	err = s.repo.Create(tx, ctx, model)
	return err
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*dto.MainOrder, error) {
	byID, err := s.repo.GetByID(ctx, id)
	return mapper.MapModelToDto(byID), err
}

func (s *OrderService) GetAll(ctx context.Context) ([]*dto.MainOrder, error) {
	all, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.MapListModelToDto(all), nil
}

func (s *OrderService) Delete(ctx context.Context, id string) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)

	err = s.repo.Delete(tx, ctx, id)
	return err
}

func (s *OrderService) Update(ctx context.Context, id string, order *dto.MainOrder) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)

	err = s.repo.Update(tx, ctx, id, mapper.MapDtoToModel(order))
	return err
}

func (s *OrderService) handleTransaction(tx *gorm.DB, err *error) {
	if r := recover(); r != nil {
		_ = utils.WrapAndLog(errors.New("panic"))
		tx.Rollback()
		panic(r)
	}
	if *err != nil {
		_ = utils.WrapAndLog(*err)
		tx.Rollback()
		return
	}
	if *err = tx.Commit().Error; *err != nil {
		_ = utils.WrapAndLog(*err)
		tx.Rollback()
		return
	}
	fmt.Println("transaction handled")
	return
}
