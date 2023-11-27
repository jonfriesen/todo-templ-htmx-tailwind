package service

import (
	"log"
	"net/http"
)

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
	message := "The server encountered a problem and could not process your request"
	http.Error(w, message, http.StatusInternalServerError)
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	http.Error(w, message, http.StatusNotFound)
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}
