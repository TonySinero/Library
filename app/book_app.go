package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
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
const MAX_UPLOAD_SIZE = 1024 * 1024

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
	a.Router.HandleFunc("/post/image", a.PostBookImage).Methods("POST")
	a.Router.HandleFunc("/load/image", a.LoadBookImage).Methods("GET")


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

	book, err := model.GetBooks(d.Database, field, sort, limit, page)
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

func (a *App) PostBookImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := viper.GetString("IMAGE_POST_PATH")
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(fmt.Sprintf("%s/%d%s", path, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprintf(w, "Upload successful")
}

func (a *App) LoadBookImage(w http.ResponseWriter, r *http.Request) {
	path := viper.GetString("IMAGE_LOAD_PATH")
	image := r.URL.Query().Get("image")
	filename := fmt.Sprintf(  "%s/%s", path, image)
	file, err := ioutil.ReadFile(filename)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}else {
		w.Write(file)
	}
}