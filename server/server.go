package server

import (
	"fmt"
	"net/http"
)

type Config map[string]interface{}

func Run(config Config) {
	ssf, ok := config["serve_static_files"]
	if ok {
		serveStaticFiles(ssf)
	}

	http.ListenAndServe(fmt.Sprintf(":%v", config["port"]), nil)
}

func serveStaticFiles(staticFiles interface{}) {
	switch v := staticFiles.(type) {
	default:
		panic(fmt.Sprintf(`Abort: unexpected type for "serve_static_files" bool, []string, or map[string]string expected; found %v`, v))
	case bool:
		if staticFiles.(bool) {
			http.Handle("/", http.FileServer(http.Dir("./public")))
		}
	case []string:
		for _, v := range staticFiles.([]string) {
			http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("./%v", v))))
		}
	case map[string]string:
		for k, v := range staticFiles.(map[string]string) {
			http.Handle(k, http.StripPrefix(k, http.FileServer(http.Dir(fmt.Sprintf("./%v", v)))))
		}
	}
}
