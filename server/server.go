package server

import (
	"fmt"
	"net/http"
)

type Config map[string]interface{}

func Run(config Config) {
	fmt.Println(config["port"])

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
