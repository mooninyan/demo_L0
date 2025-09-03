package http

import (
	"bytes"
	"context"
	"demoL0/internal/usecase/order"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) Create(ctx context.Context, order *order.MainOrder) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderService) GetByID(ctx context.Context, id string) (*order.MainOrder, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*order.MainOrder), args.Error(1)
}

func (m *MockOrderService) GetAll(ctx context.Context) ([]*order.MainOrder, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*order.MainOrder), args.Error(1)
}

func (m *MockOrderService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderService) Update(ctx context.Context, id string, order *order.MainOrder) error {
	args := m.Called(ctx, id, order)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/order", handler.createOrder).Methods(http.MethodPost)

	ord := order.MainOrder{}
	bodyBytes, _ := json.Marshal(ord)

	mockService.On("Create", mock.Anything, &ord).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetOrderByID(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/order/{id}", handler.getOrderByID).Methods(http.MethodGet)

	expectedOrder := &order.MainOrder{}
	id := "123"

	mockService.On("GetByID", mock.Anything, id).Return(expectedOrder, nil)

	req := httptest.NewRequest(http.MethodGet, "/order/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var respOrder order.MainOrder
	err := json.NewDecoder(w.Body).Decode(&respOrder)
	require.NoError(t, err)

	require.Equal(t, *expectedOrder, respOrder)

	mockService.AssertExpectations(t)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/order/{id}", handler.getOrderByID).Methods(http.MethodGet)

	id := "123"

	mockService.On("GetByID", mock.Anything, id).Return((*order.MainOrder)(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/order/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	mockService.AssertExpectations(t)
}

func TestDeleteOrder(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/order/{id}", handler.deleteOrder).Methods(http.MethodDelete)

	id := "456"

	mockService.On("Delete", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/order/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetAllOrders(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/orders", handler.getAllOrders).Methods(http.MethodGet)

	expectedOrders := []*order.MainOrder{
		{OrderUID: "someUID"},
		{},
	}

	mockService.On("GetAll", mock.Anything).Return(expectedOrders, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	var respOrders []*order.MainOrder
	err := json.NewDecoder(w.Body).Decode(&respOrders)
	require.NoError(t, err)

	require.Equal(t, expectedOrders, respOrders)

	mockService.AssertExpectations(t)
}

func TestUpdateOrder(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewHandler(mockService)
	router := mux.NewRouter()

	router.HandleFunc("/order/{id}", handler.updateOrder).Methods(http.MethodPut)

	id := "789"

	updatedOrder := &order.MainOrder{}

	bodyBytes, _ := json.Marshal(updatedOrder)

	mockService.On("Update", mock.Anything, id, &order.MainOrder{}).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/order/"+id, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": id})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockService.AssertExpectations(t)
}
