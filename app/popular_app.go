package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
	"net/http"
)

// Initialize DB and routes.
func (a *App) PopularInitialize() {
	a.initializePopularRoutes()
}

// Defines routes.
func (a *App) initializePopularRoutes() {
	// Authorized routes.
	//a.Router.Handle("/popular", a.isAuthorized(a.createPopular)).Methods("POST")
	//a.Router.Handle("/populars", a.isAuthorized(a.getPopulars)).Methods("GET")
	//a.Router.Handle("/popular/{id}", a.isAuthorized(a.updatePopular)).Methods("PUT")
	//a.Router.Handle("/popular/{id}", a.isAuthorized(a.deletePopular)).Methods("DELETE")

	a.Router.HandleFunc("/popular", a.createPopular).Methods("POST")
	a.Router.HandleFunc("/populars", a.getPopulars).Methods("GET")
	a.Router.HandleFunc("/popular/{id}", a.updatePopular).Methods("PUT")
	a.Router.HandleFunc("/popular/{id}", a.deletePopular).Methods("DELETE")
}

// Route handlers

// Gets list of popular with count and start variables from URL.
func (a *App) getPopulars(w http.ResponseWriter, r *http.Request) {
	user, err := model.GetPopulars(d.Database)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.RespondWithJSON(w, http.StatusOK, user)
}

// Inserts new popular into db.
func (a *App) createPopular(w http.ResponseWriter, r *http.Request) {
	var dt model.Popular
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreatePopular(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created popular.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates popular in db using id from URL.
func (a *App) updatePopular(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Popular
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdatePopular(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated popular.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes popular in db using id from URL.
func (a *App) deletePopular(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Popular{ID: id}
	if err := dt.DeletePopular(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
