package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Faydiamond/microservice_user/internal/user"
	"github.com/Faydiamond/microservice_user/pkg/bootstrap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == " " {
		l.Fatal("paginator limit default is required")
	}

	userRepo := user.NewRepo(l, db)
	userSrv := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userSrv, user.Config{LimPageDef: pagLimDef})

	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	server := &http.Server{
		Addr:         "127.0.0.1:8081", //exit for the port 8081, aattention please
		Handler:      router,
		ReadTimeout:  6 * time.Second,
		WriteTimeout: 6 * time.Second,
	}
	fmt.Println(" serve in 127.0.0.1:8081 ")

	log.Fatal(server.ListenAndServe())

}
