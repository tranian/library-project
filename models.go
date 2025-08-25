package main

import (
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID		int		`json:"id"`
	Title		string		`json:"title"`
	Author		string		`json:"author"`
	Year		int		`json:"year"`
	IsAvailable	bool		`json:"is_available"`
	CheckedOutTo	*string		`json:"checked_out_to"`
	CheckoutDate	*time.Time	`json:"checkout_date,omitempty"`
	DueDate		*time.Time	`json:"due_date,omitempty"`
}

func CreateBook(title string, author string, year int) (int64, error) {
	result, err := db.Exec("INSERT INTO books (title, author, year) VALUES (?, ?, ?)", title, author, year)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetAllBooks() ([]Book, error) {
	rows, err := db.Query("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.IsAvailable, &b.CheckedOutTo, &b.CheckoutDate, &b.DueDate)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func GetBook(id int) (*Book, error) {
	var b Book
	err := db.QueryRow("SELECT id, title, author, year, is_available, checked_out_to, checkout_date, due_date FROM books WHERE id = ?", id).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.IsAvailable, &b.CheckedOutTo, &b.CheckoutDate, &b.DueDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no book with given id found")
		}
		return nil, err
	}
	return &b, nil
}


func UpdateBook(id int, title string, author string, year int) error {
	_, err := db.Exec("UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?", title, author, year, id)
	return err
}

func DeleteBook(id int) error {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

func CheckoutBook(id int, checked_out_user string) error {
	var is_available bool
	err := db.QueryRow("SELECT is_available FROM books WHERE id = ?", id).Scan(&is_available)
	if err != nil {
		// check if no books
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no book with given id found")
		}
		return err
	}

	// check if book is available
	if !is_available {
		return errors.New("book not available")
	}


	checkout_date := time.Now()
	due_date := checkout_date.AddDate(0, 0, 21)
	result, err := db.Exec("UPDATE books SET is_available = FALSE, checked_out_to = ?, checkout_date = ?, due_date = ? WHERE id = ?", checked_out_user, checkout_date, due_date, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no book updated")
	}
	return nil
}

func ReturnBook(id int, checked_out_user string) error {
	var is_available bool
	err := db.QueryRow("SELECT is_available FROM books WHERE id = ?", id).Scan(&is_available)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("no book with given id found")
		}
		return err
	}
	if is_available {
		return errors.New("book is already available")
	}

	result, err := db.Exec("UPDATE books SET is_available = TRUE, checked_out_to = NULL, checkout_date = NULL, due_date = NULL WHERE id = ? AND checked_out_to= ?", id, checked_out_user)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no book updated")
	}
	return nil
}
