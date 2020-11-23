package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/base"
	"github.com/typical-go/typical-rest-server/app/book/controller"
	"github.com/typical-go/typical-rest-server/config"
)

// Server server application
type Server struct {
	*echo.Echo
	config.AppConfig
	bookController controller.BookController
}

// NewServer return instance of server
func NewServer(
	config config.AppConfig,
	bookController controller.BookController,
) *Server {

	s := &Server{
		Echo:           echo.New(),
		AppConfig:      config,
		bookController: bookController,
	}
	initMiddlewares(s)
	initRoutes(s)

	return s
}

// BaseCRUDController func
func (s *Server) BaseCRUDController(entity string, crud base.BaseCRUDController) {
	s.GET(fmt.Sprintf("/%s", entity), crud.List)
	s.POST(fmt.Sprintf("/%s", entity), crud.Create)
	s.GET(fmt.Sprintf("/%s/:id", entity), crud.Get)
	s.PUT(fmt.Sprintf("/%s", entity), crud.Update)
	s.DELETE(fmt.Sprintf("/%s/:id", entity), crud.Delete)
}

// Serve start serve http request
func (s *Server) Serve() error {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// gracefull shutdown
	go func() {
		<-gracefulStop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()

	return s.Start(s.Address)
}
