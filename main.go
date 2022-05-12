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
	"os/exec"

	"strings"

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
	r.Put("/{id}", UpdateDrawById)
	r.Post("/execute", ExecuteDrawCode)
	r.Get("/count", GetNumberOfDraws)
    
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("Created"))

}

func UpdateDrawById(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	mu := &api.Mutation{
		CommitNow: true,
	}
	req := &api.Request{CommitNow: true}
	req.Query = `{
					drawflow(func: uid(` + chi.URLParam(r, "id") + `)) {
					  v as uid
					  exportedNodes
					  name
					}
				}`

	var data interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
			
	body := data.(map[string]interface{})

	m := `
		uid(v) <name> "`+ body["name"].(string)+`" .
		uid(v) <exportedNodes> "`+ body["exportedNodes"].(string)+`" .
	`
	mu.SetNquads = []byte(m)
	req.Mutations = []*api.Mutation{mu}

	if _, err := txn.Do(context.Background(), req); err != nil {
		fmt.Println("Ocurrio un error")
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Updated"))
	

}

func ExecuteDrawCode(w http.ResponseWriter, r *http.Request) {
	
	var data interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
			
	body := data.(map[string]interface{})

	m := body["code"].(string);
	cmd := exec.Command("python")
	reader := strings.NewReader(m)
	cmd.Stdin = reader
	output, err := cmd.Output()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(output)))
	

}

func GetNumberOfDraws(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	const q = `
	{
		drawflow(func: has(name)) {
		  uid
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

func newClient() *dgo.Dgraph {

	conn, err := dgo.DialSlashEndpoint("https://blue-surf-590507.us-east-1.aws.cloud.dgraph.io/graphql", "YjU5YmE5NDBmMDIzMzAzYmY1NGQwOTAzZGY0NzI1MGU=")
	if err != nil {
	  log.Fatal(err)
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn))
}