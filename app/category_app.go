package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
)

// Initialize DB and routes.
func (a *App) CategoryInitialize() {
	a.initializeCategoryRoutes()
}

// Defines routes.
func (a *App) initializeCategoryRoutes() {
	// Authorized routes.
	//a.Router.Handle("/category", a.isAuthorized(a.CreateCategory)).Methods("POST")
	//a.Router.Handle("/categories", a.isAuthorized(a.GetCategories)).Methods("GET")

	a.Router.HandleFunc("/category", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/categories", a.getCategories).Methods("GET")

}

// Route handlers


// Gets list of category with count and start variables from URL.
func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	sort := r.URL.Query().Get("sort")
	field := r.URL.Query().Get("field")

	if sort == ""{
		sort = "ASC"
	}
	if field == ""{
		field = "name"
	}
	if limit < 1 {
		limit = 20
	}
	// Min start is 0;
	if page < 1 {
		page = 1
	}

	category, err := model.GetCategories(d.Database, field, sort, limit, page)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, category)
}

// Inserts new category into db.
func (a *App) createCategory(w http.ResponseWriter, r *http.Request) {
	var dt model.Categories
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateCategory(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}
