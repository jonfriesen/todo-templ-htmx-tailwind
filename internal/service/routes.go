package service

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/assets"
)

func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("/static/*", fileServer)

	mux.Get("/", app.home)
	mux.Post("/todo", app.addItem)
	mux.Put("/todo/{item_id}", app.completeItem)

	return mux
}
