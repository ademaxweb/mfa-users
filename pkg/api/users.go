package api

import (
	"encoding/json"
	"errors"
	"github.com/ademaxweb/mfa-go-core/pkg/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"users/pkg/db"
)

func (a *Api) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.db.GetAllUsers()
	if err != nil {
		if errors.Is(err, db.NotFound) {
			SendResponse(w, http.StatusNotFound, nil)
			return
		}
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	SendResponse(w, http.StatusOK, users)
}

func (a *Api) getUser(w http.ResponseWriter, r *http.Request) {
	a.write("GetUserByID")
	mv := mux.Vars(r)
	id, err := strconv.Atoi(mv["id"])
	if err != nil {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}
	user, err := a.db.GetUser(id)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			SendResponse(w, http.StatusNotFound, nil)
			return
		}
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	SendResponse(w, http.StatusOK, user)
}

func (a *Api) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	a.write("GetUserByEmail")
	mv := mux.Vars(r)
	email, ok := mv["email"]
	if !ok {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}
	user, err := a.db.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			SendResponse(w, http.StatusNotFound, nil)
			return
		}
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	SendResponse(w, http.StatusOK, user)
}

func (a *Api) createUser(w http.ResponseWriter, r *http.Request) {
	var newUser data.User
	type response struct {
		Id int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}
	id, err := a.db.CreateUser(newUser)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	SendResponse(w, http.StatusOK, response{Id: id})
}

func (a *Api) updateUser(w http.ResponseWriter, r *http.Request) {
	mv := mux.Vars(r)
	id, err := strconv.Atoi(mv["id"])
	if err != nil {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}

	var newUserData data.User
	if err := json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}

	err = a.db.UpdateUser(id, newUserData)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			SendResponse(w, http.StatusNotFound, nil)
			return
		}
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}

	SendResponse(w, http.StatusOK, nil)
}

func (a *Api) deleteUser(w http.ResponseWriter, r *http.Request) {
	mv := mux.Vars(r)
	id, err := strconv.Atoi(mv["id"])
	if err != nil {
		SendResponse(w, http.StatusBadRequest, nil)
		return
	}

	err = a.db.DeleteUser(id)
	if err != nil {
		if errors.Is(err, db.NotFound) {
			SendResponse(w, http.StatusNotFound, nil)
			return
		}
		SendResponse(w, http.StatusInternalServerError, err)
		return
	}

	SendResponse(w, http.StatusOK, nil)
}
