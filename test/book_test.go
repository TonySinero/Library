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

// Tests response if book table is empty.
// Deletes all records from book table and sends GET request to /book endpoint.
func TestEmptyBookTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/books", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested book is non-existent.
// Tests if status code = 404 & response message = "book not found".
func TestGetNonExistentBook(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/book/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'book not found'. Got '%s'", m["error"])
	}
}

// Test response when fetching a specific book.
// Tests if status code = 200.
func TestGetBook(t *testing.T) {
	clearTable()
	addBook(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/book/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new book by manually adding a test book to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateBook(t *testing.T) {
	clearTable()

	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	newData := model.Book{
		Name: "string1",
		Cost: 1.1,
		PricePerDay: 1.5,
		Photo: "string4",
		YearOfPublishing: 1111,
		NumberOfPages: 1222,
	}
	payload, err := json.Marshal(newData)
	if err != nil {
		t.Error("Failed to parse JSON")
	}
	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(payload))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "string1" {
		t.Errorf("Expected book bookname to be 'string1'. Got '%v'", m["name"])
	}
	if m["cost"] != float64(1.1) {
		t.Errorf("Expected book cost to be 1.1 Got '%v'", m["cost"])
	}
	if m["pricePerDay"] != float64(2.2) {
		t.Errorf("Expected book pricePerDay to be 2.2 Got '%v'", m["pricePerDay"])
	}
	if m["photo"] != "string4" {
		t.Errorf("Expected book photo to be 'string4'. Got '%v'", m["photo"])
	}
	if m["yearOfPublishing"] != uint(1111) {
		t.Errorf("Expected book yearOfPublishing to be 1111. Got '%v'", m["yearOfPublishing"])
	}
	if m["numberOfPages"] != uint(1222) {
		t.Errorf("Expected book numberOfPages to be 1222. Got '%v'", m["numberOfPages"])
	}
}

// Test process of updating a book.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateBook(t *testing.T) {
	clearTable()
	addBook(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/book/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	var jsonStr = []byte(`{"name":"string1 - updated name", "categoryID":"uuid - updated category", "authorID":"uuid - updated author", "cost": 1.1 , "pricePerDay": 2.2, "photo":"string4 - updated photo" , "yearOfPublishing": 3333, "numberOfPages": 4444}`)
	req, _ = http.NewRequest("PUT", "/book/"+testID, bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalBook["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalBook["id"], m["id"])
	}

	if m["name"] == originalBook["name"] {
		t.Errorf("Expected the bookname to change from '%v' to '%v'. Got '%v'", originalBook["name"], m["name"], m["name"])
	}
	if m["cost"] == originalBook["cost"] {
		t.Errorf("Expected the cost to change from '%v' to '%v'. Got '%v'", originalBook["cost"], m["cost"], m["cost"])
	}
	if m["pricePerDay"] == originalBook["pricePerDay"] {
		t.Errorf("Expected the pricePerDay to change from '%v' to '%v'. Got '%v'", originalBook["pricePerDay"], m["pricePerDay"], m["pricePerDay"])
	}
	if m["photo"] == originalBook["photo"] {
		t.Errorf("Expected the photo to change from '%v' to '%v'. Got '%v'", originalBook["photo"], m["photo"], m["photo"])
	}
	if m["yearOfPublishing"] == originalBook["yearOfPublishing"] {
		t.Errorf("Expected the yearOfPublishing to change from '%v' to '%v'. Got '%v'", originalBook["yearOfPublishing"], m["yearOfPublishing"], m["yearOfPublishing"])
	}
	if m["numberOfPages"] == originalBook["numberOfPages"] {
		t.Errorf("Expected the numberOfPages to change from '%v' to '%v'. Got '%v'", originalBook["numberOfPages"], m["numberOfPages"], m["numberOfPages"])
	}
}

// Test process of deleting book.
// Tests if status code = 200.
func TestDeleteBook(t *testing.T) {
	clearTable()
	addBook(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that data exists.
	req, _ := http.NewRequest("GET", "/book/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete book.
	req, _ = http.NewRequest("DELETE", "/book/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if book still exists.
	req, _ = http.NewRequest("GET", "/book/"+uuid.NewString(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addBook(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		d.Database.Exec("INSERT INTO book(id, name, cost, price_per_day, photo, year_of_publishing, number_of_pages, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", testID, "string"+strconv.Itoa(i), i, i, "string"+strconv.Itoa(i), i, i, timestamp, timestamp)
	}
}