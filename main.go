package main

import (
	"github.com/Kudzeri/go-mongo-api/config"
	"github.com/Kudzeri/go-mongo-api/controllers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(config.GetClient())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user/:id", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	log.Fatal(http.ListenAndServe(":9000", r))
}
