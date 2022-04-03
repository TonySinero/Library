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

func (a *App) IssueInitialize() {
	a.initializeIssueRoutes()
}

// Defines routes.
func (a *App) initializeIssueRoutes() {
	// Authorized routes.
	//a.Router.Handle("/issue", a.isAuthorized(a.createIssue)).Methods("POST")
	//a.Router.Handle("/issuing", a.isAuthorized(a.getIssuing)).Methods("GET")
	//a.Router.Handle("/issue/{id}", a.isAuthorized(a.getIssue)).Methods("GET")
	//a.Router.Handle("/issue/{id}", a.isAuthorized(a.deleteIssue)).Methods("DELETE")

	a.Router.HandleFunc("/issue", a.createIssue).Methods("POST")
	a.Router.HandleFunc("/issuing", a.getIssuing).Methods("GET")
	a.Router.HandleFunc("/issue/{id}", a.getIssue).Methods("GET")
	a.Router.HandleFunc("/issue/{id}", a.updateIssue).Methods("PUT")
	a.Router.HandleFunc("/issue/{id}", a.deleteIssue).Methods("DELETE")
}

// Route handlers

// Retrieves issue from db using id from URL.
func (a *App) getIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Issue{ID: id}
	if err := dt.GetIssue(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if issue not found in db.
			app.RespondWithError(w, http.StatusNotFound, "issue not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If data found respond with issue object.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Gets list of issue with count and start variables from URL.
func (a *App) getIssuing(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	sort := r.URL.Query().Get("sort")
	field := r.URL.Query().Get("field")

	if sort == ""{
		sort = "ASC"
	}
	if field == ""{
		field = "id"
	}
	if limit < 1 {
		limit = 20
	}
	// Min start is 0;
	if page < 1 {
		page = 1
	}

	issue, err := model.GetIssues(d.Database, field, sort, limit, page)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, issue)
}

// Inserts new issue into db.
func (a *App) createIssue(w http.ResponseWriter, r *http.Request) {
	var dt model.Issue
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateIssue(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created issue.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates issue in db using id from URL.
func (a *App) updateIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Issue
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateIssue(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated issue.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes issue in db using id from URL.
func (a *App) deleteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Issue{ID: id}
	if err := dt.DeleteIssue(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}