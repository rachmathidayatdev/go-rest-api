package service

import (
	"context"

	"github.com/typical-go/typical-rest-server/app/book/models"
	"github.com/typical-go/typical-rest-server/app/book/repository"
	"github.com/typical-go/typical-rest-server/pkg/dbtrxn"
)

//BookService interface
type BookService interface {
	CreateBook(book models.Book) (int64, error)
	GetBook(id int64) (*models.Book, error)
	ListBook() ([]*models.Book, error)
	UpdateBook(book models.Book) error
	DeleteBook(id int64) error
}

//InitBookService struct
type InitBookService struct {
	Repository *InitBookRepositoryInterface
}

//InitBookRepositoryInterface struct
type InitBookRepositoryInterface struct {
	Book repository.BookRepository
}

// NewBookService return new instance of BookRepository
func NewBookService(bookRepository repository.BookRepository) BookService {
	return &InitBookService{
		Repository: &InitBookRepositoryInterface{
			Book: bookRepository,
		},
	}
}

//GetBook func
func (r *InitBookService) GetBook(id int64) (*models.Book, error) {
	book, err := r.Repository.Book.Find(id)

	return book, err
}

//ListBook func
func (r *InitBookService) ListBook() ([]*models.Book, error) {
	books, err := r.Repository.Book.List()

	if err != nil {
		return books, err
	}

	return books, err
}

//CreateBook func
func (r *InitBookService) CreateBook(book models.Book) (int64, error) {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	result, err := r.Repository.Book.Insert(ctx, book)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return result, err
}

//UpdateBook func
func (r *InitBookService) UpdateBook(book models.Book) error {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	err := r.Repository.Book.Update(ctx, book)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return err
}

//DeleteBook func
func (r *InitBookService) DeleteBook(id int64) error {
	//start transaction
	ctx := context.Background()
	defer dbtrxn.Begin(&ctx)()

	err := r.Repository.Book.Delete(ctx, id)

	//transaction commit or rollback if error
	dbtrxn.Error(ctx)

	return err
}
