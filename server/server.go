package server

import (
	"fmt"
	"io"
	"net/http"
)

type Config map[string]interface{}

func Run(config Config) {
	if ssf, ssfPresent := config["serve_static_files"]; ssfPresent {
		serveStaticFiles(ssf)
	}

	if r, routerPresent := config["router"]; routerPresent {
		if router, isRouterType := r.(Router); isRouterType {
			buildRoutes(router, config)
		} else {
			// TODO: Log bad router variable
		}
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

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func buildRoutes(router Router, config Config) {

	for path, methodToHandler := range router.Routes {
		http.HandleFunc(path, func(rw http.ResponseWriter, request *http.Request) {
			handler, methodDefinedOnPath := methodToHandler[request.Method]
			if methodDefinedOnPath {
				handler(rw, request)
			} else {
				http.NotFound(rw, request)
			}
		})
	}

}
