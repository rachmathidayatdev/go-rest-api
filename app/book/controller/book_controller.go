package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/base"
	"github.com/typical-go/typical-rest-server/app/book/models"
	"github.com/typical-go/typical-rest-server/app/book/service"
	"github.com/typical-go/typical-rest-server/app/helper/strkit"
)

// BookController handle input related to Book
type BookController interface {
	base.BaseCRUDController
	//Put new route belows
}

//InitBookController struct
type InitBookController struct {
	Service *InitBookServiceInterface
}

//InitBookServiceInterface struct
type InitBookServiceInterface struct {
	Book service.BookService
}

// NewBookController return new instance of book controller
func NewBookController(bookService service.BookService) BookController {
	return &InitBookController{
		Service: &InitBookServiceInterface{
			Book: bookService,
		},
	}
}

//Get func
func (c *InitBookController) Get(ctx echo.Context) error {
	id, err := strkit.ToInt64(ctx.Param("id"))
	if err != nil {
		return invalidID(ctx, err)
	}

	book, err := c.Service.Book.GetBook(id)
	if err != nil {
		return err
	}

	if book == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("book #%d not found", id)})
	}

	return ctx.JSON(http.StatusOK, book)
}

//List func
func (c *InitBookController) List(ctx echo.Context) error {
	books, err := c.Service.Book.ListBook()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, books)
}

//Create func
func (c *InitBookController) Create(ctx echo.Context) (err error) {
	var book models.Book
	err = ctx.Bind(&book)
	if err != nil {
		return err
	}

	err = book.Validate()
	if err != nil {
		return invalidMessage(ctx, err)
	}

	result, err := c.Service.Book.CreateBook(book)
	if err != nil {
		return err
	}

	return insertSuccess(ctx, result)

}

//Update func
func (c *InitBookController) Update(ctx echo.Context) (err error) {
	var book models.Book

	err = ctx.Bind(&book)
	if err != nil {
		return err
	}

	if book.ID <= 0 {
		return invalidID(ctx, err)
	}

	err = book.Validate()
	if err != nil {
		return invalidMessage(ctx, err)
	}

	err = c.Service.Book.UpdateBook(book)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Update success"})
}

//Delete func
func (c *InitBookController) Delete(ctx echo.Context) error {
	id, err := strkit.ToInt64(ctx.Param("id"))
	if err != nil {
		return invalidID(ctx, err)
	}

	err = c.Service.Book.DeleteBook(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Delete #%d done", id)})
}
