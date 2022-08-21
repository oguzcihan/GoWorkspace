package handler

import (
	"TestApp/internal/dtos"
	. "TestApp/internal/service"
	"TestApp/internal/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	emptyBody = "empty_body"
	jsonKey   = "Content-Type"
	jsonValue = "application/json"
)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(_service UserService) UserHandler {
	return &userHandler{service: _service}
}

type userHandler struct {
	service UserService
}

var (
	user dtos.UserDto
)

func (uHandler *userHandler) Create(w http.ResponseWriter, r *http.Request) {

	bodyError := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if bodyError != nil {
		//is body empty?
		_ = JSON(w, http.StatusBadRequest, bodyError)
		return
	}

	requestErr := utils.RequestValidate(user)
	if requestErr != nil {
		//will work if validate
		_ = JSON(w, http.StatusBadRequest, requestErr)
		return
	}

	//if true,registration start
	resCreateUser, errs := uHandler.service.Create(r.Context(), user)
	if errs != nil {
		//if an error occurs while recording
		_ = JSON(w, 404, errs)
		return
	}

	//Registration successful code send
	_ = JSON(w, http.StatusCreated, resCreateUser)

}

func (uHandler *userHandler) Update(w http.ResponseWriter, r *http.Request) {

	bodyError := json.NewDecoder(r.Body).Decode(&user)
	if bodyError != nil {
		http.Error(w, emptyBody, http.StatusBadRequest)
		return
	}

	requestId := mux.Vars(r)["id"]
	if len(requestId) == 0 {
		http.Error(w, "Id", http.StatusBadRequest)
		return
	}

	resUpdateUser, err := uHandler.service.Update(r.Context(), user, requestId)
	if err != nil {
		errJson := JSON(w, http.StatusBadRequest, resUpdateUser)
		if errJson != nil {
			return
		}
	}
	errJson := JSON(w, http.StatusOK, resUpdateUser)
	//http.Error(w, "Güncellendi", http.StatusOK)
	if errJson != nil {
		return
	}

}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set(jsonKey, jsonValue)
	//w.WriteHeader()
	return json.NewEncoder(w).Encode(res)
}

/*
Marshall ve ioreader bellekten gelen verilerde kullanmak performansı daha iyi yapar
Body den alınan verilerde NewDecoder ve Decode daha performanslıdır.
*/
