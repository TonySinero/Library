package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
)

// Initialize DB and routes.
func (a *App) BookInitialize() {
	a.initializeBookRoutes()
}

// Defines routes.
func (a *App) initializeBookRoutes() {
	// Authorized routes.
	//a.Router.Handle("/book", a.isAuthorized(a.createBook)).Methods("POST")
	//a.Router.Handle("/books", a.isAuthorized(a.getBooks)).Methods("GET")
	//a.Router.Handle("/book/{id}", a.isAuthorized(a.getBook)).Methods("GET")
	//a.Router.Handle("/book/{id}", a.isAuthorized(a.updateBook)).Methods("PUT")
	//a.Router.Handle("/book/{id}", a.isAuthorized(a.deleteBook)).Methods("DELETE")

	a.Router.HandleFunc("/book", a.createBook).Methods("POST")
	a.Router.HandleFunc("/books", a.getBooks).Methods("GET")
	a.Router.HandleFunc("/book/{name}", a.getBook).Methods("GET")
	a.Router.HandleFunc("/book/{id}", a.updateBook).Methods("PUT")
	a.Router.HandleFunc("/book/{id}", a.deleteBook).Methods("DELETE")
	a.Router.HandleFunc("/book/image", a.PostFile).Methods("POST")

}

// Route handlers

// Retrieves book from db using id from URL.
func (a *App) getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	name:= vars["name"]
	dt := model.Book{Name: name}
	if err := dt.GetBook(d.Database); err != nil {
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
func (a *App) getBooks(w http.ResponseWriter, r *http.Request) {
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

	book, err := model.GetBooks(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, book)
}

// Inserts new book into db.
func (a *App) createBook(w http.ResponseWriter, r *http.Request) {
	var dt model.Book
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateBook(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created book.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates book in db using id from URL.
func (a *App) updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Book
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateBook(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated book.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes book in db using id from URL.
func (a *App) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Book{ID: id}
	if err := dt.DeleteBook(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}