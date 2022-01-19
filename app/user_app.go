package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
)

// Initialize DB and routes.
func (a *App) UserInitialize() {
	a.initializeUserRoutes()
}

// Defines routes.
func (a *App) initializeUserRoutes() {
	// Authorized routes.
	//a.Router.Handle("/user", a.isAuthorized(a.createUser)).Methods("POST")
	//a.Router.Handle("/users", a.isAuthorized(a.GetUsers)).Methods("GET")
	//a.Router.Handle("/user/{id}", a.isAuthorized(a.getUser)).Methods("GET")
	//a.Router.Handle("/user/{id}", a.isAuthorized(a.updateUser)).Methods("PUT")
	//a.Router.Handle("/user/{id}", a.isAuthorized(a.deleteUser)).Methods("DELETE")

	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user/{id}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id}", a.deleteUser).Methods("DELETE")
}

// Route handlers

// Retrieves user from db using id from URL.
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.User{ID: id}
	if err := dt.GetUser(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if user not found in db.
			app.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If data found respond with user object.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Gets list of user with count and start variables from URL.
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit > 1{
		limit = limit
	}
	if limit < 1 {
		limit = 20
	}
	// Min start is 0;
	if page < 1 {
		page = 1
	}

	user, err := model.GetUsers(d.Database, limit, page)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, user)
}

// Inserts new user into db.
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var dt model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created user.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates user in db using id from URL.
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated user.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes user in db using id from URL.
func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.User{ID: id}
	if err := dt.DeleteUser(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}