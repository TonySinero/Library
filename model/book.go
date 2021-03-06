package model

import (
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
)

// Defines book model.

type BookToAuthors struct {
	BookID         uuid.UUID  `json:"bookId"   validate:"required"   sql:"book_id"`
	AuthorID       uuid.UUID  `json:"authorId" validate:"required" sql:"author_id"`
}
type BookToCategories struct {
	BookID             uuid.UUID  `json:"bookId"   validate:"required"   sql:"book_id"`
	CategoriesID       uuid.UUID  `json:"categoriesId" validate:"required" sql:"categories_id"`
}

type Book struct {
	ID               uuid.UUID `json:"id"       sql:"uuid"`
	Name             string    `json:"name" validate:"required" sql:"name"`
	Category         []Categories `json:"category"`
	Authors          []Author  `json:"authors"`
	Cost             float32   `json:"cost" validate:"required" sql:"cost"`
	PricePerDay      float32   `json:"pricePerDay" validate:"required" sql:"price_per_day"`
	Photo            string    `json:"photo" validate:"required" sql:"photo"`
	YearOfPublishing uint      `json:"yearOfPublishing" validate:"required" sql:"year_of_publishing"`
	NumberOfPages    uint      `json:"numberOfPages" validate:"required" sql:"number_of_pages"`
	Views            uint      `json:"views" validate:"required" sql:"views"`
	CreatedAt        time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" sql:"updated_at"`
}

// Query operations

// Gets a specific book by name.
func (dt *Book) GetBook(db *sql.DB) error {
	return db.QueryRow("SELECT id, name, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at FROM book WHERE name=$1",
		dt.Name).Scan(&dt.ID, &dt.Name, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets books. Limit count and start position in db.
func GetBooks(db *sql.DB, field, sort string, limit, page int) ([]Book, error) {

	rows, err := db.Query("SELECT * FROM book ORDER BY $1 ,$2 LIMIT $3 OFFSET $4",
		field ,sort ,limit, limit*(page-1))

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()
	book := []Book{}
	// Store query results into book variable if no errors.
	for rows.Next() {
		var dt Book
		if err := rows.Scan(&dt.ID, &dt.Name, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt);
		err != nil {
			return nil, err
		}
		dt.Category = SelectCategories(db, dt.ID)
		dt.Authors = SelectAuthors(db, dt.ID)
		book = append(book, dt)
	}

	return book, nil
}

// CRUD operations

// Create new book and insert to database.
func (dt *Book) CreateBook(db *sql.DB, categoryId, authorId string, boosNumber int) error {
	if dt.Name == "" {
		return errors.New("name is required")
	}
	if dt.Cost == 0 {
		return errors.New("cost cannot be zero")
	}
	if dt.PricePerDay == 0 {
		return errors.New("pricePerDay is required")
	}
	if dt.Photo == "" {
		return errors.New("photo is required")
	}
	if dt.YearOfPublishing == 0 {
		return errors.New("yearOfPublishing cannot be zero")
	}
	if dt.NumberOfPages == 0 {
		return errors.New("numberOfPages cannot be zero")
	}
	// Scan db after creation if book exists using new book id.

	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO book(name, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, name, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at", dt.Name, dt.Cost, dt.PricePerDay, dt.Photo, dt.YearOfPublishing, dt.NumberOfPages, dt.Views, timestamp, timestamp).Scan(&dt.ID, &dt.Name, &dt.Cost, &dt.PricePerDay, &dt.Photo, &dt.YearOfPublishing, &dt.NumberOfPages, &dt.Views, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	if categoryId != ""{
		categoryID := uuid.MustParse(categoryId)
		createCategory(db, dt.ID, categoryID)
	}
	if authorId != ""{
		authorID := uuid.MustParse(authorId)
		createAuthors(db, dt.ID, authorID)
	}
	if boosNumber != 0{
		createBooks(db, dt.ID, boosNumber)
	}

	return nil
}

func createAuthors(db *sql.DB ,id, authorId uuid.UUID) error{
	err2 := db.QueryRow(
		"INSERT INTO book_authors(book_id, author_id) VALUES($1, $2)", id, authorId).Scan(&id, &authorId)
	if err2 != nil {
		return nil
	}
	return nil
}
func createCategory(db *sql.DB ,id, categoryId uuid.UUID) error{
	err2 := db.QueryRow(
		"INSERT INTO book_categories(book_id, categories_id) VALUES($1, $2)", id, categoryId).Scan(&id, &categoryId)
	if err2 != nil {
		return nil
	}
	return nil
}
func createBooks(db *sql.DB ,id uuid.UUID, booksNumber int) error{
	timestamp := time.Now()
	err2 := db.QueryRow(
		"INSERT INTO books(book_id, number_of_book, created_at, deleted_at ) VALUES($1, $2, $3, $4)", id, booksNumber, timestamp, timestamp).Scan(&id, &booksNumber, &timestamp, &timestamp)
	if err2 != nil {
		return nil
	}
	return nil
}

// Updates a specific book details by id.
func (dt *Book) UpdateBook(db *sql.DB) error {
	if dt.Name == "" {
		return errors.New("name is required")
	}
	if dt.Cost == 0 {
		return errors.New("cost cannot be zero")
	}
	if dt.PricePerDay == 0 {
		return errors.New("pricePerDay is required")
	}
	if dt.Photo == "" {
		return errors.New("photo is required")
	}
	if dt.YearOfPublishing == 0 {
		return errors.New("yearOfPublishing cannot be zero")
	}
	if dt.NumberOfPages == 0 {
		return errors.New("numberOfPages cannot be zero")
	}
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE book SET name=$1, cost=$2, price_per_day=$3, photo=$4, year_of_publishing=$5, number_of_pages=$6, views=$7, updated_at=$8 WHERE id=$9 RETURNING id, name, cost, price_per_day, photo, year_of_publishing, number_of_pages, views, created_at, updated_at", dt.Name, dt.Cost, dt.PricePerDay, dt.Photo, dt.YearOfPublishing, dt.NumberOfPages, dt.Views, timestamp, dt.ID)

	return err
}

// Deletes a specific book by id.
func (dt *Book) DeleteBook(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM book WHERE id=$1", dt.ID)

	return err
}




// Create new category and insert to database.
func (dt *BookToAuthors) CreateBookToAuthors(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO book_authors(book_id, author_id) VALUES($1, $2) RETURNING book_id, author_id", dt.BookID, dt.AuthorID).Scan(&dt.BookID, &dt.AuthorID)
	if err != nil {
		return err
	}

	return nil
}



func (dt *BookToCategories) CreateBookToCategories(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO book_categories(book_id, categories_id) VALUES($1, $2) RETURNING book_id, categories_id", dt.BookID, dt.CategoriesID).Scan(&dt.BookID, &dt.CategoriesID)
	if err != nil {
		return err
	}

	return nil
}

func SelectCategories(db *sql.DB ,id uuid.UUID) []Categories {
	get, err := db.Query("SELECT id, name, created_at FROM categories JOIN book_categories ON categories.id = book_categories.categories_id AND book_categories.book_id = $1", id)
	if err != nil{
		return nil
	}
	defer get.Close()
	category := []Categories{}
	for get.Next() {
		var cat Categories
		err = get.Scan(&cat.ID, &cat.Name, &cat.CreatedAt)
		category = append(category, cat)
	}
	return category
}


func SelectAuthors(db *sql.DB ,id uuid.UUID) []Author {
	get, err := db.Query("SELECT id, firstname, surname, date_of_birth, photo, created_at, updated_at FROM authors JOIN book_authors ON authors.id = book_authors.author_id AND book_authors.book_id = $1", id)
	if err != nil{
		return nil
	}
	defer get.Close()
	authors := []Author{}
	for get.Next() {
		var author Author
		err = get.Scan(&author.ID, &author.Firstname, &author.Surname, &author.DateOfBirth, &author.Photo, &author.CreatedAt, &author.UpdatedAt)
		authors = append(authors, author)
	}
	return authors
}

////n + 1 fix
//func GetBooksID(db *sql.DB) []uuid.UUID {
//	get, err := db.Query("SELECT book.id FROM book")
//	if err != nil{
//		return nil
//	}
//
//	defer get.Close()
//
//	var booksID []uuid.UUID
//	for get.Next() {
//		var bookId uuid.UUID
//		err = get.Scan(&bookId)
//		booksID = append(booksID, bookId)
//	}
//	return booksID
//}
//func SecondSelectCategories(db *sql.DB ,id []uuid.UUID) []Categories {
//	get, err := db.Query("SELECT book_categories.book_id, categories.name, book_categories.categories_id FROM categories JOIN book_categories ON categories.id=book_categories.categories_id AND book_categories.book_id IN ($1)", id)
//	if err != nil{
//		log.Fatalln(err)
//		return nil
//	}
//	defer get.Close()
//
//	category := []Categories{}
//	for get.Next() {
//		var cat Categories
//		err = get.Scan(&cat.ID, &cat.Name, &cat.CreatedAt)
//		category = append(category, cat)
//	}
//	return category
//}