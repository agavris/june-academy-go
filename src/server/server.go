package server

import (
	"errors"
	"fmt"
	"github.com/agavris/june-academy-go/src/server/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve(port string) error {
	// if we don't have a port we can't serve
	if port == "" {
		return errors.New("must specify a port")
	}

	// set up the router to handle the different
	// http requests / paths we want to enable
	router := mux.NewRouter()

	router.Handle("/schedule", handler.NewScheduleHandler())

	// set up the address to listen on
	addr := fmt.Sprintf(":%s", port)

	fmt.Println("Listening on", addr)
	return http.ListenAndServe(addr, router)
}
