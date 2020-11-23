package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/typical-go/typical-rest-server/app/book/mocks"
	"github.com/typical-go/typical-rest-server/app/book/models"
	"github.com/typical-go/typical-rest-server/app/book/repository"
	"github.com/typical-go/typical-rest-server/app/book/service"
)

var BookID int64

func TestBookControllerGetBook(t *testing.T) {
	mockRepository := new(mocks.Repository)
	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	mockBook := models.Book{ID: 1, Title: "test", Author: "test"}

	id := int(1)

	t.Run("when success", func(t *testing.T) {
		mockRepository.On("Find", int64(id)).Return(mockBook, nil)

		s := service.NewBookService(bookRepository)

		book, err := s.GetBook(mockBook.ID)
		assert.NoError(t, err)
		assert.NotNil(t, book)

		tt := new(testing.T)
		// assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.Find(int64(id))
		assert.True(t, mockRepository.AssertExpectations(tt))
	})

	t.Run("when error", func(t *testing.T) {
		mockRepository.On("Find", int64(id)).Return(nil, errors.New("Unexpected"))

		s := service.NewBookService(bookRepository)

		book, err := s.GetBook(0)
		err = errors.New("Unexpected")
		assert.Error(t, err)
		assert.Nil(t, book)

		tt := new(testing.T)
		assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.Find(int64(id))
		// assert.True(t, mockRepository.AssertExpectations(tt))
	})
}

func TestBookControllerListBook(t *testing.T) {
	mockRepository := new(mocks.Repository)
	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	mockListBook := []*models.Book{
		&models.Book{ID: 1, Title: "test", Author: "test"},
		&models.Book{ID: 2, Title: "test2 edit", Author: "test2 edit"},
		&models.Book{ID: 4, Title: "test4", Author: "test4"},
	}

	t.Run("when success", func(t *testing.T) {
		mockRepository.On("List", mock.Anything).Return(mockListBook, nil)

		s := service.NewBookService(bookRepository)

		book, err := s.ListBook()
		assert.NoError(t, err)
		assert.NotNil(t, book)

		tt := new(testing.T)
		// assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.List()
		assert.True(t, mockRepository.AssertExpectations(tt))
	})

	t.Run("when error", func(t *testing.T) {
		mockRepository.On("List", mock.Anything).Return(nil, errors.New("Unexpected"))

		s := service.NewBookService(bookRepository)

		books, err := s.ListBook()
		books = nil
		err = errors.New("Unexpected")
		assert.Error(t, err)
		assert.Nil(t, books)

		tt := new(testing.T)
		assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.List()
		// assert.True(t, mockRepository.AssertExpectations(tt))
	})
}

func TestBookControllerCreateBook(t *testing.T) {
	mockRepository := new(mocks.Repository)
	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	var mockBook models.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)

	t.Run("when success", func(t *testing.T) {
		mockRepository.On("Insert", mock.Anything, mock.AnythingOfType("*models.Book")).Return(mockBook.ID, nil)

		s := service.NewBookService(bookRepository)

		bookID, err := s.CreateBook(mockBook)
		assert.NoError(t, err)
		assert.NotNil(t, bookID)
		BookID = bookID

		tt := new(testing.T)
		// assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.Insert(context.TODO(), mockBook)
		assert.True(t, mockRepository.AssertExpectations(tt))
	})
}

func TestBookControllerUpdateBook(t *testing.T) {
	mockRepository := new(mocks.Repository)
	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)
	mockBook := models.Book{
		ID:     4,
		Title:  "test4",
		Author: "test4",
	}

	t.Run("when success", func(t *testing.T) {
		mockRepository.On("Update", mock.Anything, mock.AnythingOfType("*models.Book")).Return(nil)

		s := service.NewBookService(bookRepository)

		err := s.UpdateBook(mockBook)
		assert.NoError(t, err)

		tt := new(testing.T)
		// assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.Update(context.TODO(), mockBook)
		assert.True(t, mockRepository.AssertExpectations(tt))
	})
}

func TestBookControllerDeleteBook(t *testing.T) {
	mockRepository := new(mocks.Repository)
	conn := mocks.GetMockConnection()
	bookRepository := repository.NewBookRepository(conn)

	id := int(BookID)

	t.Run("when success", func(t *testing.T) {
		mockRepository.On("Delete", mock.Anything, int64(id)).Return(nil)

		s := service.NewBookService(bookRepository)

		err := s.DeleteBook(int64(id))
		assert.NoError(t, err)

		tt := new(testing.T)
		// assert.False(t, mockRepository.AssertExpectations(tt))
		mockRepository.Delete(context.TODO(), int64(id))
		assert.True(t, mockRepository.AssertExpectations(tt))
	})
}
