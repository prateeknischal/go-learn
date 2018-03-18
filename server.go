package main

import (
	"go-learn/load"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	router := util.Router()
	http.ListenAndServe(":8081", handlers.CombinedLoggingHandler(os.Stdout, router))
}
