package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pejovski/wish-list/controller"
	myerr "github.com/pejovski/wish-list/error"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler interface {
	GetList() http.HandlerFunc
	AddItem() http.HandlerFunc
	RemoveItem() http.HandlerFunc
}

type handler struct {
	controller controller.Controller
}

func newHandler(c controller.Controller) Handler {
	s := handler{
		controller: c,
	}

	return s
}

func (h handler) GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userId := params["user_id"]
		if userId == "" {
			logrus.Warnln("User id not found")
			http.Error(w, "User id not found", http.StatusBadRequest)
			return
		}

		list, err := h.controller.GetList(userId)
		if err != nil {
			logrus.Errorf("Failed to get wish list for user %s. Error: %s", userId, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, list, http.StatusOK)
	}
}

func (h handler) AddItem() http.HandlerFunc {

	type request struct {
		ProductId string `json:"product_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userId := params["user_id"]

		if userId == "" {
			logrus.Warnln("User id not found")
			http.Error(w, "User id not found", http.StatusBadRequest)
			return
		}

		var req request
		err := h.decode(w, r, &req)
		if err != nil {
			logrus.Errorf("Failed to decode request. Error: %s", err)
			http.Error(w, "Request body incorrect", http.StatusBadRequest)
			return
		}

		if req.ProductId == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		err = h.controller.AddItem(userId, req.ProductId)
		if err != nil {
			if err == myerr.ErrItemAlreadyExist {
				http.Error(w, "Item already added", http.StatusMethodNotAllowed)
				return
			}
			logrus.Errorf("Failed to add product %s for user %s. Error: %s", req.ProductId, userId, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusAccepted)
	}
}

func (h handler) RemoveItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		userId := params["user_id"]
		productId := params["product_id"]

		if userId == "" || productId == "" {
			logrus.Warnln("User or product id not found")
			http.Error(w, "User or product id not found", http.StatusBadRequest)
			return
		}

		err := h.controller.RemoveItem(userId, productId)
		if err != nil {
			logrus.Errorf("Failed to remove product %s for user %s. Error: %s", productId, userId, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusNoContent)
	}
}

func (h handler) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.Errorf("Failed to encode data. Error: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h handler) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
