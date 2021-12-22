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
		BookName: "string1",
		CategoryID: "string2",
		AuthorID: "string3",
		Cost: 1.1,
		NumberOfBook: 5,
		Photo: "string4",
		PricePerDay: 2.2,
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

	if m["bookname"] != "string1" {
		t.Errorf("Expected book bookname to be 'string1'. Got '%v'", m["bookname"])
	}

	if m["categoryID"] != "string2" {
		t.Errorf("Expected book categoryID to be 'string2'. Got '%v'", m["categoryID"])
	}
	if m["authorID"] != "string3" {
		t.Errorf("Expected book authorID to be 'string3'. Got '%v'", m["authorID"])
	}

	if m["cost"] != float64(1.1) {
		t.Errorf("Expected book cost to be 1.1 Got '%v'", m["cost"])
	}
	if m["numberofbook"] != uint(1) {
		t.Errorf("Expected book numberofbook to be 1. Got '%v'", m["numberofbook"])
	}

	if m["photo"] != "string4" {
		t.Errorf("Expected book photo to be 'string4'. Got '%v'", m["photo"])
	}
	if m["priceperday"] != float64(2.2) {
		t.Errorf("Expected book priceperday to be 2.2 Got '%v'", m["priceperday"])
	}

	if m["yearofpublishing"] != uint(1111) {
		t.Errorf("Expected book yearofpublishing to be 1111. Got '%v'", m["yearofpublishing"])
	}
	if m["numberofpages"] != uint(1222) {
		t.Errorf("Expected book numberofpages to be 1222. Got '%v'", m["numberofpages"])
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

	var jsonStr = []byte(`{"bookname":"string1 - updated bookname", "categoryID":"string2 - updated categoryID", "authorID":"string3 - updated authorID", "cost": 1.1 , "numberofbook": 2, "photo":"string4 - updated photo", "priceperday": 1.1 , "yearofpublishing": 3, "numberofpages": 4}`)
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

	if m["bookname"] == originalBook["bookname"] {
		t.Errorf("Expected the bookname to change from '%v' to '%v'. Got '%v'", originalBook["bookname"], m["bookname"], m["bookname"])
	}

	if m["categoryID"] == originalBook["categoryID"] {
		t.Errorf("Expected the categoryID to change from '%v' to '%v'. Got '%v'", originalBook["categoryID"], m["categoryID"], m["categoryID"])
	}
	if m["authorID"] == originalBook["authorID"] {
		t.Errorf("Expected the authorID to change from '%v' to '%v'. Got '%v'", originalBook["authorID"], m["authorID"], m["authorID"])
	}

	if m["cost"] == originalBook["cost"] {
		t.Errorf("Expected the cost to change from '%v' to '%v'. Got '%v'", originalBook["cost"], m["cost"], m["cost"])
	}
	if m["numberofbook"] == originalBook["numberofbook"] {
		t.Errorf("Expected the numberofbook to change from '%v' to '%v'. Got '%v'", originalBook["numberofbook"], m["numberofbook"], m["numberofbook"])
	}

	if m["photo"] == originalBook["photo"] {
		t.Errorf("Expected the photo to change from '%v' to '%v'. Got '%v'", originalBook["photo"], m["photo"], m["photo"])
	}
	if m["priceperday"] == originalBook["priceperday"] {
		t.Errorf("Expected the priceperday to change from '%v' to '%v'. Got '%v'", originalBook["priceperday"], m["priceperday"], m["priceperday"])
	}

	if m["yearofpublishing"] == originalBook["yearofpublishing"] {
		t.Errorf("Expected the yearofpublishing to change from '%v' to '%v'. Got '%v'", originalBook["yearofpublishing"], m["yearofpublishing"], m["yearofpublishing"])
	}
	if m["numberofpages"] == originalBook["numberofpages"] {
		t.Errorf("Expected the numberofpages to change from '%v' to '%v'. Got '%v'", originalBook["numberofpages"], m["numberofpages"], m["numberofpages"])
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
		d.Database.Exec("INSERT INTO books(id, bookname, categoryID, authorID, cost, numberofbook, photo, priceperday, yearofpublishing, numberofpages, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,  $11, $12)", testID, "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), "string"+strconv.Itoa(i), i, i, "string"+strconv.Itoa(i), i, i, i, timestamp, timestamp)
	}
}