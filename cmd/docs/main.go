package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"homework/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.MustNewApiConfig()

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.SwaggerPort), swaggerHandler); err != nil {
		log.Fatalln(err)
	}
}
