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
func (a *App) BooksInitialize() {
	a.initializeBooksRoutes()
}

// Defines routes.
func (a *App) initializeBooksRoutes() {
	// Authorized routes.
	//a.Router.Handle("/book/number", a.isAuthorized(a.createNumberBook)).Methods("POST")
	//a.Router.Handle("/books/number", a.isAuthorized(a.getNumberBooks)).Methods("GET")
	//a.Router.Handle("/book/number/{id}", a.isAuthorized(a.getNumberBook)).Methods("GET")
	//a.Router.Handle("/book/number/{id}", a.isAuthorized(a.updateNumberBook)).Methods("PUT")
	//a.Router.Handle("/book/number/{id}", a.isAuthorized(a.deleteBook)).Methods("DELETE")

	a.Router.HandleFunc("/book/number", a.createNumberBook).Methods("POST")
	a.Router.HandleFunc("/books/number", a.getNumberBooks).Methods("GET")
	a.Router.HandleFunc("/book/number/{id}", a.getNumberBook).Methods("GET")
	a.Router.HandleFunc("/book/number/{id}", a.updateNumberBook).Methods("PUT")
	a.Router.HandleFunc("/book/number/{id}", a.deleteNumberBook).Methods("DELETE")

}

// Route handlers

// Retrieves book from db using id from URL.
func (a *App) getNumberBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Books{ID: id}
	if err := dt.GetNumberBook(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if book not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Book not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If data found respond with book object.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Gets list of book with count and start variables from URL.
func (a *App) getNumberBooks(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	sort := r.URL.Query().Get("sort")
	field := r.URL.Query().Get("field")

	if sort == ""{
		sort = "ASC"
	}
	if field == ""{
		field = "number_of_book"
	}
	if limit < 1 {
		limit = 20
	}
	// Min start is 0;
	if page < 1 {
		page = 1
	}

	book, err := model.GetNumberBooks(d.Database, field, sort, limit, page)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, book)
}

// Inserts new book into db.
func (a *App) createNumberBook(w http.ResponseWriter, r *http.Request) {
	var dt model.Books
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateNumberBook(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created book.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates book in db using id from URL.
func (a *App) updateNumberBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Books
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateNumberBook(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated book.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes book in db using id from URL.
func (a *App) deleteNumberBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Books{ID: id}
	if err := dt.DeleteAllBooks(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
