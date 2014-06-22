package server

import (
	"bytes"
	"fmt"
	"github.com/repp/armadillo/api"
	"github.com/repp/armadillo/model"
	"net/http"
	"strings"
	"text/template"
)

type Config map[string]interface{}

func Run(config Config, router Router) {
	buildRoutes(router, config)
	http.ListenAndServe(fmt.Sprintf(":%v", config["port"]), nil)
}

func buildRoutes(router Router, config Config) {
	allRoutes := append(router.Routes, convertToRoutes(router.ModelRoutes, config)...)
	panic(allRoutes)

	//	for pathTemplate, methodToHandler := range router.Routes {
	//		path := extractPathFromTemplate(pathTemplate, config)
	//		http.HandleFunc(path, func(rw http.ResponseWriter, request *http.Request) {
	//			// Check if the current method(GET, POST, etc) has been defined for this path("/api/v1/users")
	//			handler, methodDefinedOnPath := methodToHandler[request.Method]
	//			if methodDefinedOnPath {
	//				handler(rw, request)
	//			} else {
	//				http.NotFound(rw, request)
	//			}
	//		})
	//	}

	if ssf, ssfPresent := config["serve_static_files"]; ssfPresent {
		serveStaticFiles(ssf)
	}

}

func convertToRoutes(modelRoutes []ModelRoute, config Config) (routes []Route) {
	api.MMC.Gateway = config["db"].(model.DbGateway)
	for _, route := range modelRoutes {
		modelMethod := api.AllModelMethods[route.Action]
		path := strings.Replace(modelMethod.Path, "{{.model_name}}", route.ModelName, 1)
		handler := modelMethod.HandlerGenerator(route.ModelInstance)
		routes = append(routes, Route{Method: modelMethod.Method, Path: path, Handler: handler})
	}
	return
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

func extractPathFromTemplate(path string, config Config) string {
	var b bytes.Buffer
	t := template.Must(template.New("path").Parse(path))
	err := t.Execute(&b, config)
	if err != nil {
		panic(err)
	}
	return safeFormatPath(b.String())
}
