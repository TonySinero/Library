package app

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	app "github.com/library/app/utils"
	"github.com/library/model"
	"net/http"
	"strconv"
)

// Used for validating header tokens.
var mySigningKey = []byte("secret")

// Initialize DB and routes.
func (a *App) AdminInitialize() {
	a.initializeAdminRoutes()
}

// Defines routes.
func (a *App) initializeAdminRoutes() {
	a.Router.HandleFunc("/admin", a.createAdmin).Methods("POST")
	a.Router.HandleFunc("/admin/login", a.loginAdmin).Methods("POST")
	// Authorized routes.
	//a.Router.Handle("/admin/{id}", a.isAuthorized(a.getAdmin)).Methods("GET")
	//a.Router.Handle("/admins", a.isAuthorized(a.getAdmins)).Methods("GET")
	//a.Router.Handle("/admin/{id}", a.isAuthorized(a.updateAdmin)).Methods("PUT")
	//a.Router.Handle("/admin/{id}", a.isAuthorized(a.deleteAdmin)).Methods("DELETE")

	a.Router.HandleFunc("/admins", a.getAdmins).Methods("GET")
	a.Router.HandleFunc("/admin/{id}", a.getAdmin).Methods("GET")
	a.Router.HandleFunc("/admin/{id}", a.updateAdmin).Methods("PUT")
	a.Router.HandleFunc("/admin/{id}", a.deleteAdmin).Methods("DELETE")
}

// Route handlers

// Retrieves admin from db using id from URL.
func (a *App) loginAdmin(w http.ResponseWriter, r *http.Request) {
	var u model.Admin
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	// Find admin in db with email and password from request body.
	if err := u.GetAdminByEmailAndPassword(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if admin not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Admin not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// Generate and send token to client with response header.
	validToken, err := GenerateJWT()
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Add("Token", validToken)
	// Respond with user in db.
	app.RespondWithJSON(w, http.StatusOK, u)
}

// Retrieves admin from db using id from URL.
func (a *App) getAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	u := model.Admin{ID: id}
	if err := u.GetAdmin(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if admin not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Admin not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If admin found respond with admin object.
	app.RespondWithJSON(w, http.StatusOK, u)
}

// Gets list of admin with count and start variables from URL.
func (a *App) getAdmins(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	start, _ := strconv.Atoi(r.URL.Query().Get("start"))

	if count > 1{
		count = count
	}
	if count < 1 {
		count = 20
	}
	// Min start is 0;
	if start < 0 {
		start = 0
	}

	admins, err := model.GetAdmins(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, admins)
}

// Inserts new admin into db.
func (a *App) createAdmin(w http.ResponseWriter, r *http.Request) {
	var u model.Admin
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := u.CreateAdmin(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created admin.
	app.RespondWithJSON(w, http.StatusCreated, u)
}

// Updates admin in db using id from URL.
func (a *App) updateAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var u model.Admin
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateAdmin(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated admin.
	app.RespondWithJSON(w, http.StatusOK, u)
}

// Deletes admin in db using id from URL.
func (a *App) deleteAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	u := model.Admin{ID: id}
	if err := u.DeleteAdmin(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "Admin deleted"})
}

// Helper functions

// Generate JWT
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		// fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
