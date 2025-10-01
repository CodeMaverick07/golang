package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"github.com/codemaverick07/api/internals/app"
	"github.com/codemaverick07/api/internals/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	r := routes.SetupRoutes(app)
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()
	app.Logger.Println("first log")
	http.HandleFunc("/health", app.HealthCheck)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Logger.Printf("running server on port %d \n", port)
	error := server.ListenAndServe()
	if error != nil {
		app.Logger.Fatal("issue with server", err)
	}

}
