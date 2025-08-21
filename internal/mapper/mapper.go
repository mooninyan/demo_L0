package mapper

import (
	"demoL0/internal/domain"
	"demoL0/internal/dto"
	"github.com/jinzhu/copier"
)

func MapDtoToModel(orderDto *dto.MainOrder) (order *domain.MainOrder) {
	order = &domain.MainOrder{}

	_ = copier.Copy(&order, &orderDto)
	for _, it := range order.Items {
		it.OrderUID = order.OrderUID
	}
	return order
}

func MapModelToDto(order *domain.MainOrder) (orderDto *dto.MainOrder) {
	orderDto = &dto.MainOrder{}

	_ = copier.Copy(&orderDto, &order)
	return orderDto
}

func MapListModelToDto(all []*domain.MainOrder) []*dto.MainOrder {
	var res = make([]*dto.MainOrder, len(all))
	for i := 0; i < len(all); i++ {
		res[i] = MapModelToDto(all[i])
	}
	return res
}
