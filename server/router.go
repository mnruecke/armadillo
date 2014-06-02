package server

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	Routes map[string]map[string]httpHandler
}

type Route struct {
	Method  string
	Path    string
	Handler httpHandler
}

type httpHandler func(http.ResponseWriter, *http.Request)
type Set map[string]struct{}

func NewSet(elements ...string) Set {
	s := make(Set)
	for _, e := range elements {
		s[e] = struct{}{}
	}
	return s
}

type MethodRules struct {
	Allow  Set
	Forbid Set
}

var tmpFunc httpHandler = func(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello World!"))
	rw.WriteHeader(200)
}

var allModelMethods = map[string]Route{
	"Create":    Route{"POST", "/{{.ModelName}}", tmpFunc},
	"Find":      Route{"GET", "/{{.ModelName}}/:id", tmpFunc},
	"FindAll":   Route{"GET", "/{{.ModelName}}", tmpFunc},
	"Update":    Route{"PATCH", "/{{.ModelName}}/:id", tmpFunc},
	"UpdateAll": Route{"PATCH", "/{{.ModelName}}", tmpFunc},
	"Replace":   Route{"PUT", "/{{.ModelName}}/:id", tmpFunc},
	"Delete":    Route{"DELETE", "/{{.ModelName}}/:id", tmpFunc},
	"DeleteAll": Route{"DELETE", "/{{.ModelName}}", tmpFunc},
	"Options":   Route{"OPTIONS", "/{{.ModelName}}", tmpFunc},
}

func (r *Router) Model(publicName string, modelInstance interface{}, methodRules MethodRules) {
	allowedMethods := allowedModelMethods(methodRules)
	for _, route := range allowedMethods {
		route.Path = strings.Replace(route.Path, "{{.ModelName}}", publicName, 1)
		r.appendRoute(route)
	}
}

func (r *Router) Action(path string, handler httpHandler) {
	r.appendNewRoute("POST", fmt.Sprintf("{{ .ActionPath }}/%v", path), handler)
}

func (r *Router) Get(path string, handler httpHandler) {
	r.appendNewRoute("GET", path, handler)
}

func (r *Router) Head(path string, handler httpHandler) {
	r.appendNewRoute("HEAD", path, handler)
}

func (r *Router) Post(path string, handler httpHandler) {
	r.appendNewRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler httpHandler) {
	r.appendNewRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler httpHandler) {
	r.appendNewRoute("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler httpHandler) {
	r.appendNewRoute("DELETE", path, handler)
}

func (r *Router) Options(path string, handler httpHandler) {
	r.appendNewRoute("OPTIONS", path, handler)
}

func (r *Router) Trace(path string, handler httpHandler) {
	r.appendNewRoute("TRACE", path, handler)
}

func (r *Router) appendNewRoute(method string, path string, handler httpHandler) {
	if r.Routes == nil {
		r.Routes = make(map[string]map[string]httpHandler)
	}
	if r.Routes[path] == nil {
		r.Routes[path] = make(map[string]httpHandler)
	}
	r.Routes[path][method] = handler
}

func (r *Router) appendRoute(route Route) {
	r.appendNewRoute(route.Method, route.Path, route.Handler)
}

func allowedModelMethods(rules MethodRules) (allowedMethods map[string]Route) {
	if len(rules.Allow) > 0 && len(rules.Forbid) > 0 {
		panic(fmt.Sprintf("Abort: Invalid ModelRules: %v \n Allow or Forbid cannot both contain rules.", rules))
	}

	allowedMethods = make(map[string]Route)

	// If Forbidden, get diff between all methods and forbidden, send results to add model routes method
	if len(rules.Forbid) > 0 {
		allowedMethods = allModelMethods
		for key, _ := range rules.Forbid {
			if _, keyExists := allowedMethods[key]; keyExists {
				delete(allowedMethods, key)
			}
		}
	}

	// If Allow, loop through allow, adding routes
	if len(rules.Allow) > 0 {
		for key, _ := range rules.Allow {
			route, ok := allModelMethods[key]
			if ok {
				allowedMethods[key] = route
			}
		}
	}

	// both Allow and Forbid are empty, all methods are allowed.
	if len(rules.Allow) == len(rules.Forbid) {
		allowedMethods = allModelMethods
	}

	return
}
