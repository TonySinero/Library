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
	"github.com/library/model"
)

// Test functions

// Tests response if data table is empty.
// Deletes all records from data table and sends GET request to /data endpoint.
func TestEmptyUserTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/users", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested user is non-existent.
// Tests if status code = 404 & response message = "user not found".
func TestGetNonExistentUser(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "user not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'user not found'. Got '%s'", m["error"])
	}
}

// Test response when fetching a specific user.
// Tests if status code = 200.
func TestGetUser(t *testing.T) {
	clearTable()
	addUser(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new user by manually adding a test user to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateUser(t *testing.T) {
	clearTable()

	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}


	newData := model.User{
		Firstname: "string1",
		Surname: "string2",
		SecondName: "string3",
		Passport: "string4",
		DateOfBirth: "string5",
		Email: "string6",
		Address: "string7",
		Indebtedness: "string8",
	}
	payload, err := json.Marshal(newData)
	if err != nil {
		t.Error("Failed to parse JSON")
	}
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["firstname"] != "string1" {
		t.Errorf("Expected user firstname to be 'string1'. Got '%v'", m["firstname"])
	}
	if m["surname"] != "string2" {
		t.Errorf("Expected user surname to be 'string2'. Got '%v'", m["surname"])
	}
	if m["secondname"] != "string3" {
		t.Errorf("Expected user secondname to be 'string3'. Got '%v'", m["secondname"])
	}
	if m["passport"] != "string4" {
		t.Errorf("Expected user passport to be 'string4'. Got '%v'", m["passport"])
	}
	if m["dateofbirth"] != "string5" {
		t.Errorf("Expected user dateofbirth to be 'string5'. Got '%v'", m["dateofbirth"])
	}
	if m["email"] != "string6" {
		t.Errorf("Expected user email to be 'string6'. Got '%v'", m["email"])
	}
	if m["address"] != "string7" {
		t.Errorf("Expected user address to be 'string7'. Got '%v'", m["address"])
	}
	if m["indebtedness"] != "string8" {
		t.Errorf("Expected user indebtedness to be 'string8'. Got '%v'", m["indebtedness"])
	}

}

// Test process of updating a data.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateUser(t *testing.T) {
	clearTable()
	addUser(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	var jsonStr = []byte(`{"firstname":"string1 - updated firstname", "surname":"string2 - updated surname", "secondname":"string3 - updated secondname", "passport":"string4 - updated passport", "dateofbirth":"string5 - updated dateofbirth", "email":"string6 - updated email", "address":"string7 - updated address", "indebtedness":"string8 - updated indebtedness",}`)
	req, _ = http.NewRequest("PUT", "/user/"+testID, bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}

	if m["firstname"] == originalUser["firstname"] {
		t.Errorf("Expected the firstname to change from '%v' to '%v'. Got '%v'", originalUser["firstname"], m["firstname"], m["firstname"])
	}

	if m["surname"] == originalUser["surname"] {
		t.Errorf("Expected the surname to change from '%v' to '%v'. Got '%v'", originalUser["surname"], m["surname"], m["surname"])
	}
	if m["secondname"] == originalUser["secondname"] {
		t.Errorf("Expected the secondname to change from '%v' to '%v'. Got '%v'", originalUser["secondname"], m["secondname"], m["secondname"])
	}

	if m["passport"] == originalUser["passport"] {
		t.Errorf("Expected the passport to change from '%v' to '%v'. Got '%v'", originalUser["passport"], m["passport"], m["passport"])
	}
	if m["dateofbirth"] == originalUser["dateofbirth"] {
		t.Errorf("Expected the dateofbirth to change from '%v' to '%v'. Got '%v'", originalUser["dateofbirth"], m["dateofbirth"], m["dateofbirth"])
	}

	if m["email"] == originalUser["email"] {
		t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalUser["email"], m["email"], m["email"])
	}
	if m["address"] == originalUser["address"] {
		t.Errorf("Expected the address to change from '%v' to '%v'. Got '%v'", originalUser["address"], m["address"], m["address"])
	}

	if m["indebtedness"] == originalUser["indebtedness"] {
		t.Errorf("Expected the indebtedness to change from '%v' to '%v'. Got '%v'", originalUser["indebtedness"], m["indebtedness"], m["indebtedness"])
	}
}

// Test process of deleting user.
// Tests if status code = 200.
func TestDeleteUser(t *testing.T) {
	clearTable()
	addUser(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that data exists.
	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete user.
	req, _ = http.NewRequest("DELETE", "/user/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if user still exists.
	req, _ = http.NewRequest("GET", "/user/"+uuid.NewString(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addUser(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		d.Database.Exec("INSERT INTO users(id, firstname, surname, secondname, passport, dateofbirth, email, address, indebtedness, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", testID, "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), timestamp, timestamp)
	}
}
