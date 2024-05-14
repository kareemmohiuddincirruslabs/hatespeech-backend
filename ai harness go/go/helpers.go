package main

import (
    "log"
    "net/http"
)

func (app *Application) handleError(w http.ResponseWriter, errMsg string, err error, statusCode int) {
    log.Printf("%s: %v\n", errMsg, err)
    http.Error(w, errMsg, statusCode)
}