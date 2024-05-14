package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) routes() http.Handler{
    mux := chi.NewRouter()

    mux.Use(middleware.Recoverer)
    mux.Use(app.enableCORS)

    mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("../outputs"))))

    // mux.Get("/", app.uploadPageHandler)
    mux.Post("/upload", app.handleUpload)

    return mux
}
