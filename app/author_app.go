package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
	"net/http"
	"strconv"
)

// Initialize DB and routes.
func (a *App) AuthorInitialize() {
	a.initializeAuthorRoutes()
}

// Defines routes.
func (a *App) initializeAuthorRoutes() {
	// Authorized routes.
	//a.Router.Handle("/author", a.isAuthorized(a.createUser)).Methods("POST")
	//a.Router.Handle("/authors", a.isAuthorized(a.GetUsers)).Methods("GET")
	//a.Router.Handle("/author/{id}", a.isAuthorized(a.updateUser)).Methods("PUT")

	a.Router.HandleFunc("/author", a.createAuthor).Methods("POST")
	a.Router.HandleFunc("/authors", a.getAuthors).Methods("GET")
	a.Router.HandleFunc("/author/{id}", a.updateAuthor).Methods("PUT")
	a.Router.HandleFunc("/post/image", a.PostImage).Methods("POST")
	a.Router.HandleFunc("/load/image", a.LoadImage).Methods("GET")

}

// Route handlers


// Gets list of authors with count and start variables from URL.
func (a *App) getAuthors(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	sort := r.URL.Query().Get("sort")
	field := r.URL.Query().Get("field")

	if sort == ""{
		sort = "ASC"
	}
	if field == ""{
		field = "firstname"
	}
	if limit < 1 {
		limit = 20
	}
	// Min start is 0;
	if page < 1 {
		page = 1
	}

	author, err := model.GetAuthors(d.Database, field, sort, limit, page)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, author)
}

// Inserts new author into db.
func (a *App) createAuthor(w http.ResponseWriter, r *http.Request) {
	var dt model.Author
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateAuthor(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created author.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates author in db using id from URL.
func (a *App) updateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Author
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateAuthor(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated author.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

