package main

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

func main() {
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)

	if err := http.ListenAndServe(":8888", swaggerHandler); err != nil {
		log.Fatalln(err)
	}
}
