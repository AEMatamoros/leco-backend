package main

import (
	"fmt"

	"encoding/json"

    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

)

func main() {

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/", ApiGetExample)
	r.Post("/", ApiPostExample)
    http.ListenAndServe(":3000", r)

}

func ApiGetExample(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello World From API!"))

}

func ApiPostExample(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var data interface{}

    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	m := data.(map[string]interface{})

	fmt.Println(m["ok"])
	fmt.Println(m["ok2"])
	fmt.Println(m["ok3"])

	jsonResp, err := json.Marshal(m)
	w.Write(jsonResp)

}