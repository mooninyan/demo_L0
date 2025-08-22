package rest

import (
	"context"
	"demoL0/internal/dto"
	"demoL0/internal/service"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	orderService *service.OrderService
}

func NewHandler(service *service.OrderService) *Handler {
	return &Handler{
		orderService: service,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	orders := r.PathPrefix("/order").Subrouter()
	{
		orders.HandleFunc("", h.createOrder).Methods(http.MethodPost)
		orders.HandleFunc("", h.getAllOrders).Methods(http.MethodGet)
		orders.HandleFunc("/{id}", h.getOrderByID).Methods(http.MethodGet)
		orders.HandleFunc("/{id}", h.deleteOrder).Methods(http.MethodDelete)
		orders.HandleFunc("/{id}", h.updateOrder).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	log.Println("get by id call")
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getOrderByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	order, err := h.orderService.GetByID(timeout, id)
	if err != nil {
		log.Println("getOrderByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if order == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(order)
	if err != nil {
		log.Println("getOrderByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order dto.MainOrder
	if err = json.Unmarshal(reqBytes, &order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.orderService.Create(timeout, &order)
	if err != nil {
		log.Println("createOrder() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteOrder() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.orderService.Delete(timeout, id)
	if err != nil {
		log.Println("deleteOrder() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("get all orders call")
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	orders, err := h.orderService.GetAll(timeout)
	if err != nil {
		log.Println("getAllOrders() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(orders)
	if err != nil {
		log.Println("getAllOrders() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp dto.MainOrder
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.orderService.Update(timeout, id, &inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return "", errors.New("id can't be null")
	}

	return id, nil
}
