package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/categolj/blog-feed/handler"
)


func main() {
	http.HandleFunc("/", handler.Feed)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "4000"
	}
	log.Printf(fmt.Sprintf("Listening at %s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

