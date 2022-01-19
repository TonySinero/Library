package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
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
	a.Router.HandleFunc("/accept/image", a.PostAcceptBookImage).Methods("POST")
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

	acceptance, err := model.GetAcceptances(d.Database, limit, page)
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

func (a *App) PostAcceptBookImage(w http.ResponseWriter, r *http.Request) {
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

		path := viper.GetString("IMAGE_ACCEPT_PATH")
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
