package main
import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	// "time"
	"os"

	// mock of sql-driver
	"github.com/DATA-DOG/go-sqlmock"
)

var mock sqlmock.Sqlmock

func setupMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mock_db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	return mock_db, mock
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
 
func TestCreateBook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO books (title, author, year) VALUES (?, ?, ?)")).WithArgs("Moby-Dick", "Herman Melville", 1851).WillReturnResult(sqlmock.NewResult(1, 1))
		id, err := CreateBook("Moby-Dick", "Herman Melville", 1851)
		if err != nil {
			t.Errorf("Expected no error, received: %v", err)
		}
		if id != 1 {
			t.Errorf("Expected id == 1, received: %d", id)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
	
	t.Run("Error", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO books (title, author, year) VALUES (?, ?, ?")).WithArgs("", "Herman Melville", 1851).WillReturnError(errors.New("NOT NULL constraint failed"))

		_, err := CreateBook("", "Herman Melville", 1851)
		if err == nil {
			t.Error("Expected error, received nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
}

func TestGetAllBooks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "is_available", "checked_out_to", "checkout_date", "due_date"}).AddRow(1, "Moby-Dick", "Herman Melville", 1851, true, nil, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books")).WillReturnRows(rows)

		books, err := GetAllBooks()
		if err != nil {
			t.Errorf("Expected no error, received: %v", err)
		}
		if len(books) != 1 {
			t.Errorf("Expected 1 book, received: %d", len(books))
		}
		if len(books) > 0 && books[0].Title != "Moby-Dick" {
			t.Errorf("Expected Moby-Dick, received: %v", books[0].Title)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books")).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "year", "is_available", "checked_out_to", "checkout_date", "due_date"}))
		books, err := GetAllBooks()
		if err != nil {
			t.Errorf("Expected no error, received: %v", err)
		}
		if len(books) != 0 {
			t.Errorf("Expected 0 books, received: %d", len(books))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books")).WillReturnError(errors.New("database error"))

		_, err := GetAllBooks()
		if err == nil {
			t.Error("Expected error, received nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
}


func TestGetBook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "is_available", "checked_out_to", "checkout_date", "due_date"}).AddRow(1, "Moby-Dick", "Herman Melville", 1851, true, nil, nil, nil)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books WHERE id = ?")).WithArgs(1).WillReturnRows(rows)

		book, err := GetBook(1)
		if err != nil {
			t.Errorf("Expected no error, received: %v", err)
		}
		if book == nil || book.ID != 1 {
			t.Errorf("Expected 1 book, received %d", book.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books WHERE id = ?")).WithArgs(999).WillReturnError(sql.ErrNoRows)
		book, err := GetBook(999)
		if err == nil || err.Error() != "no book with given id found" {
			t.Errorf("Expected 'no book with given id' error, given %v", err)
		}

		if book != nil {
			t.Errorf("Expected nil book, received %v", book)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
}

func TestUpdateBook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?")).WithArgs("Moby-Dick Updated Title", "Herman Melville", 1852, 1).WillReturnResult(sqlmock.NewResult(0, 1))
		err := UpdateBook(1, "Moby-Dick Updated Title", "Herman Melville", 1852)
		if err != nil {
			t.Errorf("Expected no error, received %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?")).WithArgs("Moby-Dick Updated Title", "Herman Melville", 1851, 123).WillReturnResult(sqlmock.NewResult(0, 0))

		err := UpdateBook(123, "Moby-Dick Updated Title", "Herman Melville", 1851)
		if err != nil {
			t.Errorf("Expected no error, received %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
}

func TestDeleteBook(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = ?")).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
		err := DeleteBook(1)
		if err != nil {
			t.Errorf("Expected no error, received %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock_db, mock := setupMock(t)
		defer mock_db.Close()
		original_db := db
		db = mock_db
		defer func() { db = original_db }()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM books WHERE id = ?")).WithArgs(123).WillReturnResult(sqlmock.NewResult(0, 0))
		err := DeleteBook(123)
		if err != nil {
			t.Errorf("Expected no error, received %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Expectations were not met: %v", err)
		}
	})
}

func TestCheckoutBook(t *testing.T) {
}

func TestReturnBook(t *testing.T) {
}


