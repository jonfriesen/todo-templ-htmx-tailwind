package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/todo"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 100 * time.Second
	defaultShutdownPeriod = 30 * time.Second
	defaultPort           = 3000
)

type Application struct {
	TodoService todo.TodoList
}

func (app *Application) ServeHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", defaultPort),
		Handler:      app.routes(),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	fmt.Printf("starting server on %s\n", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	fmt.Printf("stopped server on %s", srv.Addr)

	return nil
}
