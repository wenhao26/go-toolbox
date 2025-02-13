package service

import (
	"fmt"

	"toolbox/2025/book-manager-example/model"
	"toolbox/2025/book-manager-example/repo"
)

// BookSrv 业务服务层
type BookSrv struct {
	repo *repo.BookRepo
}

// NewBookSrv 创建BookSrv实例
func NewBookSrv(repo *repo.BookRepo) *BookSrv {
	return &BookSrv{repo: repo}
}

// CreateBook 创建书籍
func (s *BookSrv) CreateBook(book *model.Book) error {
	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		return fmt.Errorf("invalid book data")
	}
	return s.repo.CreateBook(book)
}

// GetBook 获取书籍
func (s *BookSrv) GetBook(id int) (*model.Book, error) {
	return s.repo.GetBookByID(id)
}
