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
func (a *App) AcceptanceInitialize() {
	a.initializeAcceptanceRoutes()
}

// Defines routes.
func (a *App) initializeAcceptanceRoutes() {
	// Authorized routes.
	//a.Router.Handle("/acceptance", a.isAuthorized(a.createAcceptance)).Methods("POST")
	//a.Router.Handle("/acceptances", a.isAuthorized(a.getAcceptances)).Methods("GET")
	//a.Router.Handle("/acceptance/{id}", a.isAuthorized(a.getAcceptance)).Methods("GET")
	//a.Router.Handle("/acceptance/{id}", a.isAuthorized(a.updateAcceptance)).Methods("PUT")
	//a.Router.Handle("/acceptance/{id}", a.isAuthorized(a.deleteAcceptance)).Methods("DELETE")

	a.Router.HandleFunc("/acceptance", a.createAcceptance).Methods("POST")
	a.Router.HandleFunc("/acceptances", a.getAcceptances).Methods("GET")
	a.Router.HandleFunc("/acceptance/{id}", a.getAcceptance).Methods("GET")
	a.Router.HandleFunc("/acceptance/{id}", a.updateAcceptance).Methods("PUT")
	a.Router.HandleFunc("/acceptance/{id}", a.deleteAcceptance).Methods("DELETE")
	a.Router.HandleFunc("/profit", a.getProfit).Methods("GET")
}

// Route handlers

// Retrieves acceptance from db using id from URL.
func (a *App) getAcceptance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Acceptance{ID: id}
	if err := dt.GetAcceptance(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if acceptance not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Acceptance not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If data found respond with acceptance object.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Gets list of acceptance with count and start variables from URL.
func (a *App) getAcceptances(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	// Default and limit of count is 20.
	if count > 20 || count < 1 {
		count = 20
	}
	// Min start is 0;
	if start < 0 {
		start = 0
	}

	acceptance, err := model.GetAcceptances(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, acceptance)
}

// Inserts new acceptance into db.
func (a *App) createAcceptance(w http.ResponseWriter, r *http.Request) {
	var dt model.Acceptance
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateAcceptance(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created acceptance.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates acceptance in db using id from URL.
func (a *App) updateAcceptance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Acceptance
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateAcceptance(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated acceptance.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes acceptance in db using id from URL.
func (a *App) deleteAcceptance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Acceptance{ID: id}
	if err := dt.DeleteAcceptance(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getProfit(w http.ResponseWriter, r *http.Request) {
	acceptance, err := model.GetProfit(d.Database)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.RespondWithJSON(w, http.StatusOK, acceptance)
}
