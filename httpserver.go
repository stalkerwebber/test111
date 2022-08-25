package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var router *mux.Router

func Init() {
	router = mux.NewRouter()
	routes()
}

func Serve() {
	addr := "127.0.0.1:4100"
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"force-refresh", "force-redirect"},
		MaxAge:           60 * 60 * 24,
		Debug:            false,
	})
	handler := c.Handler(router)
	srv := &http.Server{
		Handler:      handler,
		Addr:         addr,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
	}

	log.Println("Starting server on http://" + addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
