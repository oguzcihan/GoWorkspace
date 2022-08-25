package handler

import (
	"TestApp/internal/dtos"
	. "TestApp/internal/service"
	"TestApp/internal/utils"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	emptyBody = "empty_body"
	jsonKey   = "Content-Type"
	jsonValue = "application/json"
)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(_service UserService) UserHandler {
	return &userHandler{service: _service}
}

type userHandler struct {
	service UserService
}

var (
	user        dtos.UserDto
	customError utils.CustomError

	EmptyId = utils.NewError("cannot_empty_id", 404)
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
	if errors.As(errs, &customError) {
		//if an error occurs while recording
		_ = JSON(w, customError.Status, errs)
		return
	}

	//Registration successful code send
	_ = JSON(w, http.StatusCreated, resCreateUser)

}

func (uHandler *userHandler) Save(w http.ResponseWriter, r *http.Request) {

	bodyError := json.NewDecoder(r.Body).Decode(&user)
	if bodyError != nil {
		http.Error(w, emptyBody, http.StatusBadRequest)
		return
	}

	requestId := mux.Vars(r)["id"]
	if len(requestId) == 0 {
		http.Error(w, "Id boş olamaz", http.StatusBadRequest)
		return
	}

	requestErr := utils.RequestValidate(user)
	if requestErr != nil {
		//will work if validate
		//validate yaparken status code ile gösterilecek
		_ = JSON(w, http.StatusBadRequest, requestErr)
		return
	}

	userId, errId := strconv.Atoi(requestId) //hata olursa?
	if errId != nil {
		return
	}

	resUpdateUser, err := uHandler.service.Save(r.Context(), user, userId)
	if errors.As(err, &customError) {
		_ = JSON(w, customError.Status, err)
		return
	}
	_ = JSON(w, http.StatusOK, resUpdateUser)

}

func (uHandler *userHandler) Delete(w http.ResponseWriter, r *http.Request) {

	requestId := mux.Vars(r)["id"]
	if len(requestId) == 0 {
		//http.Error(w, "Id boş olamaz", http.StatusBadRequest)
		_ = JSON(w, http.StatusBadRequest, EmptyId)
		return
	}

	userId, convertError := strconv.Atoi(requestId) //hata olursa?
	if convertError != nil {
		return
	}

	resDeleteUser := uHandler.service.Delete(r.Context(), userId)
	if resDeleteUser != nil {
		//http.Error(w, resDeleteUser.Error(), http.StatusInternalServerError)
		_ = JSON(w, http.StatusBadRequest, resDeleteUser)
		return
	}
	_ = JSON(w, http.StatusOK, SuccessUserDelete)
}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set(jsonKey, jsonValue)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}

/*
Marshall ve ioreader bellekten gelen verilerde kullanmak performansı daha iyi yapar
Body den alınan verilerde NewDecoder ve Decode daha performanslıdır.
*/
