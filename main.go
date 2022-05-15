//Paquete principal del la app LECO(Learn to Code, aplicacion que senseñara de forma visual a programar en pPython a jovenes)
package main

import (

	"fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	c "main/controllers"

)

//Controla las rutas de las peticiones a la API utilizando chi, asi como inicia el servidor en el puerto 3000
func main() {
	fmt.Println("Servidor levantado en el puerto :3000")

    r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, 
	})
	r.Use(cors.Handler)

    r.Use(middleware.Logger)
    r.Get("/diagrams/all/{offset}", c.GetAllDraws)
	r.Get("/diagrams/{id}", c.GetDrawById)
	r.Post("/diagrams/", c.PostDraw)
	r.Put("/diagrams/{id}", c.UpdateDrawById)
	r.Post("/diagrams/execute", c.ExecuteDrawCode)
	r.Get("/diagrams/count", c.GetNumberOfDraws)
    
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}

}
