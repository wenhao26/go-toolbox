package repo

import (
	"database/sql"
	"fmt"

	"toolbox/2025/book-manager-example/config"
	"toolbox/2025/book-manager-example/model"
)

// BookRepo 数据访问层
type BookRepo struct{}

// NewBookRepo 创建一个BookRepo实例
func NewBookRepo() *BookRepo {
	return &BookRepo{}
}

// CreateBook 创建书籍
func (r *BookRepo) CreateBook(book *model.Book) error {
	query := "INSERT INTO books (title, author, price) VALUES (?, ?, ?)"
	_, err := config.DB.Exec(query, book.Title, book.Author, book.Price)
	return err
}

// GetBookByID 根据ID获取书籍
func (r *BookRepo) GetBookByID(id int) (*model.Book, error) {
	query := "SELECT id, title, author, price FROM books WHERE id = ?"
	row := config.DB.QueryRow(query, id)
	book := &model.Book{}
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, err
	}
	return book, nil
}
