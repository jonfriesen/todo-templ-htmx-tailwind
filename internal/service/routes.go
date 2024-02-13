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

	mux.Route("/signin", func(mux chi.Router) {
		mux.Get("/", app.signin)
		mux.Post("/", app.signinUser)
	})

	mux.Get("/signout", app.signoutUser)

	mux.Route("/register", func(mux chi.Router) {
		mux.Get("/", app.register)
		mux.Post("/", app.registerUser)
	})

	mux.Route("/todo", func(mux chi.Router) {
		// auth required
		mux.Use(app.checkAuthMiddleware)

		mux.Get("/", app.todo)
		mux.Post("/", app.addItem)
		mux.Put("/{item_id}", app.completeItem)
	})

	return mux
}
