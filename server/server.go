package server

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

type Config map[string]interface{}

func Run(config Config, router Router) {
	if ssf, ssfPresent := config["serve_static_files"]; ssfPresent {
		serveStaticFiles(ssf)
	}
	buildRoutes(router, config)
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

func buildRoutes(router Router, config Config) {

	for pathTemplate, methodToHandler := range router.Routes {
		path := extractPathFromTemplate(pathTemplate, config)
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

func extractPathFromTemplate(path string, config Config) string {
	var b bytes.Buffer
	t := template.Must(template.New("path").Parse(path))
	err := t.Execute(&b, config)
	if err != nil {
		panic(err)
	}
	return safeFormatPath(b.String())
}
