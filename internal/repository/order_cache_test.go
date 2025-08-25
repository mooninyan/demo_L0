package repository

import (
	"context"
	"testing"

	"demoL0/internal/domain"
	"demoL0/internal/mapcache"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(tx *gorm.DB, ctx context.Context, order *domain.MainOrder) error {
	args := m.Called(tx, ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) GetByID(ctx context.Context, id string) (*domain.MainOrder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.MainOrder), args.Error(1)
}

func (m *MockOrderRepo) GetAll(ctx context.Context) ([]*domain.MainOrder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.MainOrder), args.Error(1)
}

func (m *MockOrderRepo) Delete(tx *gorm.DB, ctx context.Context, id string) error {
	args := m.Called(tx, ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepo) Update(tx *gorm.DB, ctx context.Context, id string, order *domain.MainOrder) error {
	args := m.Called(tx, ctx, id, order)
	return args.Error(0)
}

func createSampleOrder(id string) *domain.MainOrder {
	return &domain.MainOrder{
		OrderUID: id,
	}
}

func TestCachedOrderRepository_Create(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	order := createSampleOrder("order1")

	mockRepo.On("Create", mock.Anything, ctx, order).Return(nil)

	err := cr.Create(nil, ctx, order)
	assert.NoError(t, err)

	cachedOrder, found := cache.Get(order.OrderUID)
	assert.True(t, found)
	assert.Equal(t, order, cachedOrder)

	err = cr.Create(nil, ctx, order)
	assert.Error(t, err)
	assert.Equal(t, "order already exists", err.Error())
}

func TestCachedOrderRepository_GetByID_CacheMiss(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	orderID := "order2"
	order := createSampleOrder(orderID)

	mockRepo.On("GetByID", ctx, orderID).Return(order, nil)

	result, err := cr.GetByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	cachedOrder, found := cache.Get(orderID)
	assert.True(t, found)
	assert.Equal(t, order, cachedOrder)
}

func TestCachedOrderRepository_GetByID_CacheHit(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	orderID := "order3"
	order := createSampleOrder(orderID)

	cache.Set(orderID, order)

	result, err := cr.GetByID(ctx, orderID)
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	mockRepo.AssertNotCalled(t, "GetByID", ctx, orderID)
}

func TestCachedOrderRepository_GetAll_CacheMiss(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	order1 := createSampleOrder("id1")
	order2 := createSampleOrder("id2")
	allOrders := []*domain.MainOrder{order1, order2}

	mockRepo.On("GetAll", ctx).Return(allOrders, nil)

	result, err := cr.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	cachedOrder, found := cache.Get("id1")
	assert.True(t, found)
	assert.Equal(t, order1, cachedOrder)

	cachedOrder, found = cache.Get("id2")
	assert.True(t, found)
	assert.Equal(t, order2, cachedOrder)
}

func TestCachedOrderRepository_GetAll_CacheHit(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	order1 := createSampleOrder("id1")
	order2 := createSampleOrder("id2")

	cache.Set("id1", order1)
	cache.Set("id2", order2)

	result, err := cr.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, order1)
	assert.Contains(t, result, order2)

	mockRepo.AssertNotCalled(t, "GetAll", ctx)
}

func TestCachedOrderRepository_Delete(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	orderID := "orderDel"
	order := createSampleOrder(orderID)

	cache.Set(orderID, order)

	mockRepo.On("Delete", mock.Anything, ctx, orderID).Return(nil)

	err := cr.Delete(nil, ctx, orderID)
	assert.NoError(t, err)

	_, found := cache.Get(orderID)
	assert.False(t, found)
}

func TestCachedOrderRepository_Update(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	cache := mapcache.NewCache[string, *domain.MainOrder]()
	cr := CachedOrderRepository{
		repo:  mockRepo,
		cache: cache,
	}

	ctx := context.Background()
	orderID := "orderUpdate"
	updatedOrder := createSampleOrder(orderID)
	updatedOrder.TrackNumber = "newTrack"

	mockRepo.On("Update", mock.Anything, ctx, orderID, updatedOrder).Return(nil)

	err := cr.Update(nil, ctx, orderID, updatedOrder)
	assert.NoError(t, err)

	cachedOrder, found := cache.Get(orderID)
	assert.True(t, found)
	assert.Equal(t, updatedOrder, cachedOrder)
}
