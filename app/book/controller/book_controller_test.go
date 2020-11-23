package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/app/book/controller"
	"github.com/typical-go/typical-rest-server/app/book/mocks"
	"github.com/typical-go/typical-rest-server/app/book/models"
	"github.com/typical-go/typical-rest-server/app/book/repository"
	"github.com/typical-go/typical-rest-server/app/book/service"
)

var BookID int64

func TestBookControllerGet(t *testing.T) {
	mockService := new(mocks.Service)
	mockBook := models.Book{ID: 1, Title: "test", Author: "test"}

	id := int(1)

	mockService.On("GetBook", int64(id)).Return(mockBook, nil)
	tt := new(testing.T)
	assert.False(t, mockService.AssertExpectations(tt))
	mockService.GetBook(int64(id))
	assert.True(t, mockService.AssertExpectations(tt))

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/book/"+strconv.Itoa(id), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("book/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	bookService := service.NewBookService(bookRepository)
	initBookServiceInterface := &controller.InitBookServiceInterface{
		Book: bookService,
	}
	initBookController := &controller.InitBookController{
		Service: initBookServiceInterface,
	}

	err = initBookController.Get(c)
	require.NoError(t, err)

	expect, _ := json.Marshal(&mockBook)
	assert.Equal(t, string(expect)+"\n", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBookControllerList(t *testing.T) {
	mockService := new(mocks.Service)
	mockListBook := []*models.Book{
		&models.Book{ID: 1, Title: "test", Author: "test"},
		&models.Book{ID: 2, Title: "test2 edit", Author: "test2 edit"},
		&models.Book{ID: 4, Title: "test4", Author: "test4"},
	}

	mockService.On("ListBook", mock.Anything).Return(mockListBook, nil)
	tt := new(testing.T)
	assert.False(t, mockService.AssertExpectations(tt))
	mockService.ListBook()
	assert.True(t, mockService.AssertExpectations(tt))

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/book", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	bookService := service.NewBookService(bookRepository)
	initBookServiceInterface := &controller.InitBookServiceInterface{
		Book: bookService,
	}
	initBookController := &controller.InitBookController{
		Service: initBookServiceInterface,
	}

	err = initBookController.List(c)
	require.NoError(t, err)

	expect, _ := json.Marshal(&mockListBook)
	assert.Equal(t, string(expect)+"\n", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBookControllerCreate(t *testing.T) {
	mockService := new(mocks.Service)
	var mockBook models.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)

	mockService.On("CreateBook", mock.Anything, mock.AnythingOfType("*models.Book")).Return(mockBook.ID, nil)
	tt := new(testing.T)
	assert.False(t, mockService.AssertExpectations(tt))
	mockService.CreateBook(context.TODO(), &mockBook)
	assert.True(t, mockService.AssertExpectations(tt))

	jsonParam, err := json.Marshal(mockBook)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/book", strings.NewReader(string(jsonParam)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/book")

	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	bookService := service.NewBookService(bookRepository)
	initBookServiceInterface := &controller.InitBookServiceInterface{
		Book: bookService,
	}
	initBookController := &controller.InitBookController{
		Service: initBookServiceInterface,
	}

	err = initBookController.Create(c)
	require.NoError(t, err)

	var respBody map[string]interface{}
	err = json.NewDecoder(rec.Body).Decode(&respBody)
	assert.NoError(t, err)
	BookID = int64(respBody["data"].(float64))

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestBookControllerUpdate(t *testing.T) {
	mockService := new(mocks.Service)
	mockBook := &models.Book{
		ID:     4,
		Title:  "test4",
		Author: "test4",
	}

	mockService.On("UpdateBook", mock.Anything, mock.AnythingOfType("*models.Book")).Return(nil)
	tt := new(testing.T)
	assert.False(t, mockService.AssertExpectations(tt))
	mockService.UpdateBook(context.TODO(), &mockBook)
	assert.True(t, mockService.AssertExpectations(tt))

	jsonParam, err := json.Marshal(mockBook)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/book", strings.NewReader(string(jsonParam)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/book")

	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	bookService := service.NewBookService(bookRepository)
	initBookServiceInterface := &controller.InitBookServiceInterface{
		Book: bookService,
	}
	initBookController := &controller.InitBookController{
		Service: initBookServiceInterface,
	}

	err = initBookController.Update(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestBookControllerDelete(t *testing.T) {
	mockService := new(mocks.Service)

	id := int(BookID)

	mockService.On("DeleteBook", mock.Anything, int64(id)).Return(nil)
	tt := new(testing.T)
	assert.False(t, mockService.AssertExpectations(tt))
	mockService.DeleteBook(context.TODO(), int64(id))
	assert.True(t, mockService.AssertExpectations(tt))

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/book/"+strconv.Itoa(id), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("book/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	bookService := service.NewBookService(bookRepository)
	initBookServiceInterface := &controller.InitBookServiceInterface{
		Book: bookService,
	}
	initBookController := &controller.InitBookController{
		Service: initBookServiceInterface,
	}

	err = initBookController.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}
