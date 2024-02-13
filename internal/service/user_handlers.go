package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/components"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/db"
)

func (app *Application) register(w http.ResponseWriter, r *http.Request) {
	comp := components.Page(components.Center(components.RegistrationForm()))
	comp.Render(r.Context(), w)
}

type newUserForm struct {
	Name            string `form:"name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirmPassword"`
}

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	var form newUserForm
	err := DecodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if form.Password != form.ConfirmPassword {
		app.badRequest(w, r, fmt.Errorf("passwords do not match"))
		return
	}

	if form.Name == "" || form.Email == "" || form.Password == "" {
		app.badRequest(w, r, fmt.Errorf("all fields are required"))
		return
	}

	user, err := app.UserService.CreateUser(r.Context(), &db.CreateUserParams{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Respond with the created user or redirect as needed
	log.Printf("User created successfully: %s", user.ID)
	r.Method = http.MethodGet
	http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
}

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	comp := components.Page(components.Center(components.SigninForm()))
	comp.Render(r.Context(), w)
}

type signinUserForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (app *Application) signinUser(w http.ResponseWriter, r *http.Request) {
	var form signinUserForm
	err := DecodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if form.Email == "" || form.Password == "" {
		app.badRequest(w, r, fmt.Errorf("email and password are required"))
		return
	}

	isValid, user, err := app.UserService.CheckUserPassword(r.Context(), form.Email, form.Password)
	if err != nil || !isValid {
		app.badRequest(w, r, fmt.Errorf("invalid credentials"))
		return
	}

	session, err := app.cookieStore.Get(r, app.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	session.Values["is_authed"] = true
	session.Values["authed_at"] = time.Now().Unix()
	session.Values["user_id"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Implement session or token creation for logged in user
	log.Printf("Signin successful for user: %s", form.Email)

	w.Header().Add("HX-Redirect", "/todo")
}

func (app *Application) signoutUser(w http.ResponseWriter, r *http.Request) {
	session, err := app.cookieStore.Get(r, app.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Resetting session values
	session.Values["is_authed"] = false
	session.Values["authed_at"] = time.Now().Unix() // You might also want to remove this value
	delete(session.Values, "user_id")               // Remove user ID from session

	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	log.Printf("Signout successful")

	w.Header().Add("HX-Redirect", "/signin") // Redirect to login or home page after signout
}
