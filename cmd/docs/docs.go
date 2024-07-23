package main

import (
	"bytes"
	"encoding/json"
	"github.com/swaggo/swag"
	"homework/config"
	"os"
	"strings"
)

func init() {
	docsPATH := os.Getenv("ORDER_JSON_DOCS_PATH")
	file, err := os.ReadFile(docsPATH)
	if err != nil {
		panic(err)
	}

	api, err := config.NewApiConfig()
	if err != nil {
		panic(err)
	}

	template := string(file)

	m := make(map[string]any)
	err = json.NewDecoder(strings.NewReader(template)).Decode(&m)
	if err != nil {
		panic(err)
	}
	m["basePath"] = "/"
	m["host"] = api.HttpENDPOINT

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(m)
	if err != nil {
		panic(err)
	}

	swagger := &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  buf.String(),
	}

	swag.Register(swagger.InstanceName(), swagger)
}
