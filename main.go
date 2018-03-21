package main

import (
	"log"
	"net/http"
)

func main(){
	DatabaseConnect()
	router := NewRouter(AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
	defer db.Close()
}


