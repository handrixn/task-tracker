package router

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(true)

	return r
}
