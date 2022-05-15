//Este paquete contiene los controladores de las rutas(CRUD) de la app LECO.
//En este paquete se crea el objeto de conexion a la base de datos.
package controllers

import (

	"fmt"
	"encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"context"
	"os/exec"
	"strings"
	db "main/db"

)

var conn = db.NewClient()

type Drawflow struct{
	name string
	exportedNodes string
}

//Retorna el conjunto de diagramsa creados por el usuario, recibe como parametro el offset, y envia los diagramas en grupos de 9
func GetAllDraws(w http.ResponseWriter, r *http.Request) {
	
	txn := conn.NewTxn()
	
	q := `
	{
		drawflow(func: has(name), first: 9 , offset:`+ chi.URLParam(r, "offset") + `) {
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

//Retorna un solo diagrama respecto a su identificador unico, recibe como parametro este identificador
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

//Crea y almacena el diagrama creado por el usuario, recibe como parametro dentro del cuerpo de la peticion un objeto con la estructura de Drawflow
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

//ACtualiza un diagrama respecto a su identificador unico, recibe como parametro en el cuerpo de la peticion dicho objeto
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

//Ejecuta el codigo generado en python y retorna una respuesta en formato cadena con la salida de la terminal
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

//Retornar el numero total de elementos que hay almacenados
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