package order

import (
	"context"
	"demoL0/internal/domain"
	"demoL0/internal/ports/repository"
	"demoL0/internal/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

// Service определяет интерфейс для бизнес-логики заказов
type Service interface {
	Create(ctx context.Context, order *MainOrder) error
	GetByID(ctx context.Context, id string) (*MainOrder, error)
	GetAll(ctx context.Context) ([]*MainOrder, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, order *MainOrder) error
}

// Service реализует бизнес-логику для заказов
type service struct {
	db   *gorm.DB
	repo repository.OrderRepository
}

// NewService создает новый экземпляр сервиса заказов
func NewService(db *gorm.DB, repo repository.OrderRepository) Service {
	return &service{
		db:   db,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, order *MainOrder) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)
	err = tx.Error
	if err != nil {
		return err
	}

	model := mapDtoToModel(order)
	err = s.repo.Create(tx, ctx, model)
	return err
}

func (s *service) GetByID(ctx context.Context, id string) (*MainOrder, error) {
	byID, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapModelToDto(byID), nil
}

func (s *service) GetAll(ctx context.Context) ([]*MainOrder, error) {
	all, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapListModelToDto(all), nil
}

func (s *service) Delete(ctx context.Context, id string) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)

	err = s.repo.Delete(tx, ctx, id)
	return err
}

func (s *service) Update(ctx context.Context, id string, order *MainOrder) (err error) {
	tx := s.db.Begin()
	defer s.handleTransaction(tx, &err)

	err = s.repo.Update(tx, ctx, id, mapDtoToModel(order))
	return err
}

func (s *service) handleTransaction(tx *gorm.DB, err *error) {
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
	log.Println("transaction handled")
	return
}

// mapDtoToModel конвертирует DTO в domain модель
func mapDtoToModel(dto *MainOrder) *domain.MainOrder {
	if dto == nil {
		return nil
	}

	items := make([]domain.Item, len(dto.Items))
	for i, item := range dto.Items {
		items[i] = domain.Item{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			RID:         item.RID,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}

	return &domain.MainOrder{
		OrderUID:    dto.OrderUID,
		TrackNumber: dto.TrackNumber,
		Entry:       dto.Entry,
		Delivery: domain.Delivery{
			Name:    dto.Delivery.Name,
			Phone:   dto.Delivery.Phone,
			Zip:     dto.Delivery.Zip,
			City:    dto.Delivery.City,
			Address: dto.Delivery.Address,
			Region:  dto.Delivery.Region,
			Email:   dto.Delivery.Email,
		},
		Payment: domain.Payment{
			Transaction:  dto.Payment.Transaction,
			RequestID:    dto.Payment.RequestID,
			Currency:     dto.Payment.Currency,
			Provider:     dto.Payment.Provider,
			Amount:       dto.Payment.Amount,
			PaymentDt:    dto.Payment.PaymentDt,
			Bank:         dto.Payment.Bank,
			DeliveryCost: dto.Payment.DeliveryCost,
			GoodsTotal:   dto.Payment.GoodsTotal,
			CustomFee:    dto.Payment.CustomFee,
		},
		Items:             items,
		Locale:            dto.Locale,
		InternalSignature: dto.InternalSignature,
		CustomerID:        dto.CustomerID,
		DeliveryService:   dto.DeliveryService,
		Shardkey:          dto.Shardkey,
		SmID:              dto.SmID,
		DateCreated:       dto.DateCreated,
		OofShard:          dto.OofShard,
	}
}

// mapModelToDto конвертирует domain модель в DTO
func mapModelToDto(model *domain.MainOrder) *MainOrder {
	if model == nil {
		return nil
	}

	items := make([]Item, len(model.Items))
	for i, item := range model.Items {
		items[i] = Item{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			RID:         item.RID,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}

	return &MainOrder{
		OrderUID:    model.OrderUID,
		TrackNumber: model.TrackNumber,
		Entry:       model.Entry,
		Delivery: Delivery{
			Name:    model.Delivery.Name,
			Phone:   model.Delivery.Phone,
			Zip:     model.Delivery.Zip,
			City:    model.Delivery.City,
			Address: model.Delivery.Address,
			Region:  model.Delivery.Region,
			Email:   model.Delivery.Email,
		},
		Payment: Payment{
			Transaction:  model.Payment.Transaction,
			RequestID:    model.Payment.RequestID,
			Currency:     model.Payment.Currency,
			Provider:     model.Payment.Provider,
			Amount:       model.Payment.Amount,
			PaymentDt:    model.Payment.PaymentDt,
			Bank:         model.Payment.Bank,
			DeliveryCost: model.Payment.DeliveryCost,
			GoodsTotal:   model.Payment.GoodsTotal,
			CustomFee:    model.Payment.CustomFee,
		},
		Items:             items,
		Locale:            model.Locale,
		InternalSignature: model.InternalSignature,
		CustomerID:        model.CustomerID,
		DeliveryService:   model.DeliveryService,
		Shardkey:          model.Shardkey,
		SmID:              model.SmID,
		DateCreated:       model.DateCreated,
		OofShard:          model.OofShard,
	}
}

// mapListModelToDto конвертирует список domain моделей в список DTO
func mapListModelToDto(models []*domain.MainOrder) []*MainOrder {
	if models == nil {
		return nil
	}

	result := make([]*MainOrder, len(models))
	for i, model := range models {
		result[i] = mapModelToDto(model)
	}
	return result
}
