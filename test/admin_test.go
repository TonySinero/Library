package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/library/app"
)

// Test functions

// Tests response if admins table is empty.
// Deletes all records from admins table and sends GET request to /admins endpoint.
func TestEmptyAdminTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/admins", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested admin is non-existent.
// Tests if status code = 404 & response message = "Admin not found".
func TestGetNonExistentAdmin(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/admin/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Admin not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'admin not found'. Got '%s'", m["error"])
	}
}

// Test response on login route.
// Tests if status code = 200.
func TestLoginAdmin(t *testing.T) {
	clearTable()
	addAdmin(1)

	var jsonStr = []byte(`{"email":"testemail1@gmail.com", "password":"password1"}`)
	req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test response when fetching a specific admin.
// Tests if status code = 200.
func TestGetAdmin(t *testing.T) {
	clearTable()
	addAdmin(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/admin/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new admin by manually adding a test admin to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateAdmin(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"email":"testemail1@gmail.com", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/admin", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["email"] != "testemail1@gmail.com" {
		t.Errorf("Expected admin email to be 'testemail1@gmail.com'. Got '%v'", m["email"])
	}

	if m["password"] != "password1" {
		t.Errorf("Expected admin password to be 'password1'. Got '%v'", m["password"])
	}
}

// Test process of updating admin.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateAdmin(t *testing.T) {
	clearTable()
	addAdmin(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/admin/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalAdmin map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalAdmin)

	var jsonStr = []byte(`{"email":"testemail1@gmail.com - updated email", "password": "password1 - updated password"}`)
	req, _ = http.NewRequest("PUT", "/admin/"+testID, bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalAdmin["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalAdmin["id"], m["id"])
	}

	if m["email"] == originalAdmin["email"] {
		t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalAdmin["email"], m["email"], m["email"])
	}

	if m["password"] == originalAdmin["password"] {
		t.Errorf("Expected the password to change from '%v' to '%v'. Got '%v'", originalAdmin["password"], m["password"], m["password"])
	}
}

// Test process of deleting admins.
// Tests if status code = 200.
func TestDeleteAdmin(t *testing.T) {
	clearTable()
	addAdmin(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that admin exists.
	req, _ := http.NewRequest("GET", "/admin/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete user.
	req, _ = http.NewRequest("DELETE", "/admin/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if user still exists.
	req, _ = http.NewRequest("GET", "/admin/"+uuid.NewString(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addAdmin(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		d.Database.Exec("INSERT INTO admins(id, email, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5)", testID, "testemail"+strconv.Itoa(i)+"@gmail.com", "password"+strconv.Itoa(i), timestamp, timestamp)
	}
}
