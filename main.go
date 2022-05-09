package main

import (
	"fmt"
	"encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
	"log"
	"context"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/go-chi/cors"

)

var conn = newClient()

type Drawflow struct{
	name string
	exportedNodes string
}

func main() {
	fmt.Println("Servidor levantado en el puerto :3000")

    r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

    r.Use(middleware.Logger)
    r.Get("/", GetAllDraws)
	r.Get("/{id}", GetDrawById)
	r.Post("/", PostDraw)
    
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}

func GetAllDraws(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	const q = `
	{
		drawflow(func: has(name)) {
		  uid
		  exportedNodes
		  name
		}
	}
	`

	resp, err := txn.Query(context.Background(), q)

	if err != nil {
		log.Fatal(err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
	

}

func GetDrawById(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	q := `
	{
		drawflow(func: uid(` + chi.URLParam(r, "id") + `)) {
		  uid
		  exportedNodes
		  name
		}
	}
	`

	resp, err := txn.Query(context.Background(), q)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.Json)
	

}

func PostDraw(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	
	var drawflow interface {}
    err := json.NewDecoder(r.Body).Decode(&drawflow)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	txn := conn.NewTxn()

	pb, err := json.Marshal(drawflow)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		SetJson: pb,
		CommitNow: true,
	}

	response, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	w.WriteHeader(200)
	w.Write([]byte("Created"))

}

func newClient() *dgo.Dgraph {

	conn, err := dgo.DialSlashEndpoint("https://blue-surf-590507.us-east-1.aws.cloud.dgraph.io/graphql", "YjU5YmE5NDBmMDIzMzAzYmY1NGQwOTAzZGY0NzI1MGU=")
	if err != nil {
	  log.Fatal(err)
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn))
}