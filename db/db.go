//Este paquete contiene la configuracion de la conexión a la base de datos
package db

import (

	"log"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"

)

//Esta funcion genera la conexion a la base de datos(Solo importar una vez para evitar creacion de multiples instancias)
func NewClient() *dgo.Dgraph {

	conn, err := dgo.DialSlashEndpoint("https://blue-surf-590507.us-east-1.aws.cloud.dgraph.io/graphql", "YjU5YmE5NDBmMDIzMzAzYmY1NGQwOTAzZGY0NzI1MGU=")
	if err != nil {
	  log.Fatal(err)
	}
	return dgo.NewDgraphClient(api.NewDgraphClient(conn))
}