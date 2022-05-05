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

)

var conn = newClient()

func main() {

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/", ApiGetExample)
	r.Post("/", ApiPostExample)
    http.ListenAndServe(":3000", r)

}

func ApiGetExample(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	const q = `
	{
		querynode(func: has(id)) {
		  uid
		  id
		  name
		  data{
			value
		  }
		  class
		  html
		  typenode
		  inputs{
			input_1{
			  connections{
				node
				input
				output
			  }
			}
			input_2{
			  connections{
				node
				input
				output
			  }
			}
			input_3{
			  connections{
				node
				input
				output
			  }
			}
			input_4{
			  connections{
				node
				input
				output
			  }
			}
		  }
		  outputs{
			output_1{
			  connections{
				node
				input
				output
			  }
			}
			output_2{
			  connections{
				node
				input
				output
			  }
			}
			output_3{
			  connections{
				node
				input
				output
			  }
			}
			output_4{
			  connections{
				node
				input
				output
			  }
			}
		  }
		  pos_x
		  pos_y
		  
		}
	  }
	`

	resp, err := txn.Query(context.Background(), q)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp.Json))
	w.Header().Set("Content-Type", "text/plain")
	w.Write(resp.Json)
	


	

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

func newClient() *dgo.Dgraph {

	conn, err := dgo.DialSlashEndpoint("https://blue-surf-590475.us-east-1.aws.cloud.dgraph.io/graphql", "NjlkNWU3ODYxNzY5YTVhYjdhNGZkZWNjOTQ5YmJhNzI=")
	if err != nil {
	  log.Fatal(err)
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn))
}