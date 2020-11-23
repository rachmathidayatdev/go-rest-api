package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/app/book/models"
	"github.com/typical-go/typical-rest-server/pkg/dbtrxn"
)

// BookRepository to get book data from databasesa
type BookRepository interface {
	Find(id int64) (*models.Book, error)
	List() ([]*models.Book, error)
	Insert(ctx context.Context, book models.Book) (lastInsertID int64, err error)
	Update(ctx context.Context, book models.Book) error
	Delete(ctx context.Context, id int64) error
}

//InitBookRepository struct
type InitBookRepository struct {
	conn *sql.DB
}

// NewBookRepository return new instance of BookRepository
func NewBookRepository(conn *sql.DB) BookRepository {
	return &InitBookRepository{
		conn: conn,
	}
}

//Find func
func (r *InitBookRepository) Find(id int64) (book *models.Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).
		From(bookTable).
		Where(sq.Eq{idColumn: id})

	rows, err := builder.RunWith(r.conn).Query()
	if err != nil {
		return book, err
	}

	if rows.Next() {
		book, err = models.ScanBook(rows)
	}

	return book, err
}

//List func
func (r *InitBookRepository) List() (list []*models.Book, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(BookColumns...).From(bookTable).OrderBy("id ASC")

	rows, err := builder.RunWith(r.conn).Query()

	list = make([]*models.Book, 0)

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var book *models.Book
		book, err = models.ScanBook(rows)
		if err != nil {
			return
		}
		list = append(list, book)
	}

	return list, err
}

//Insert func
func (r *InitBookRepository) Insert(ctx context.Context, book models.Book) (lastInsertID int64, err error) {
	trxn, err := dbtrxn.Use(ctx, r.conn)

	query := sq.Insert(bookTable).
		Columns(bookTitleColumn, bookAuthorColumn).
		Values(book.Title, book.Author).
		Suffix("RETURNING \"id\"").
		RunWith(trxn.DB).
		PlaceholderFormat(sq.Dollar)

	err = query.QueryRow().Scan(&book.ID)
	if err != nil {
		trxn.SetError(err)
		return lastInsertID, err
	}

	lastInsertID = book.ID
	return lastInsertID, err
}

//Update func
func (r *InitBookRepository) Update(ctx context.Context, book models.Book) (err error) {
	trxn, err := dbtrxn.Use(ctx, r.conn)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update(bookTable).
		Set(bookTitleColumn, book.Title).
		Set(bookAuthorColumn, book.Author).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: book.ID})

	_, err = builder.RunWith(trxn.DB).Exec()

	if err != nil {
		trxn.SetError(err)
		return err
	}

	return err
}

//Delete func
func (r *InitBookRepository) Delete(ctx context.Context, id int64) (err error) {
	trxn, err := dbtrxn.Use(ctx, r.conn)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Delete(bookTable).
		Where(sq.Eq{idColumn: id})

	_, err = builder.RunWith(trxn.DB).Exec()
	if err != nil {
		trxn.SetError(err)
		return err
	}

	return err
}
