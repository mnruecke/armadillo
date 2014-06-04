package server

import (
	"fmt"
	"net/http"
	"regexp"
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
	"Create":     Route{"POST", "/{{.api_prefix}}/{{.model_name}}", tmpFunc},
	"Find":       Route{"GET", "/{{.api_prefix}}/{{.model_name}}/:id", tmpFunc},
	"FindAll":    Route{"GET", "/{{.api_prefix}}/{{.model_name}}", tmpFunc},
	"Update":     Route{"PATCH", "/{{.api_prefix}}/{{.model_name}}/:id", tmpFunc},
	"UpdateAll":  Route{"PATCH", "/{{.api_prefix}}/{{.model_name}}", tmpFunc},
	"Replace":    Route{"PUT", "/{{.api_prefix}}/{{.model_name}}/:id", tmpFunc},
	"Destroy":    Route{"DELETE", "/{{.api_prefix}}/{{.model_name}}/:id", tmpFunc},
	"DestroyAll": Route{"DELETE", "//{{.api_prefix}}/{{.model_name}}", tmpFunc},
	"Info":       Route{"OPTIONS", "/{{.api_prefix}}/{{.model_name}}", tmpFunc},
}

func (r *Router) Model(publicName string, modelInstance interface{}, methodRules MethodRules) {
	allowedMethods := allowedModelMethods(methodRules)
	for _, route := range allowedMethods {
		r.appendModelRoute(route, publicName)
	}
}

func (r *Router) Create(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Create"], publicName)
}

func (r *Router) Find(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Find"], publicName)
}

func (r *Router) FindAll(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["FindAll"], publicName)
}

func (r *Router) Update(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Update"], publicName)
}

func (r *Router) UpdateAll(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["UpdateAll"], publicName)
}

func (r *Router) Replace(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Replace"], publicName)
}

func (r *Router) Destroy(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Destroy"], publicName)
}

func (r *Router) DestroyAll(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["DestroyAll"], publicName)
}

func (r *Router) Info(publicName string, modelInstance interface{}) {
	r.appendModelRoute(allModelMethods["Info"], publicName)
}

func (r *Router) Action(path string, handler httpHandler) {
	r.appendNewRoute("POST", fmt.Sprintf("/{{.api_prefix}}/{{ .action_prefix }}/%v", path), handler)
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
	path = safeFormatPath(path)
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

func (r *Router) appendModelRoute(route Route, publicName string) {
	route.Path = strings.Replace(route.Path, "{{.model_name}}", publicName, 1)
	r.appendRoute(route)
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

func safeFormatPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	r := regexp.MustCompile("/{2,}")
	return r.ReplaceAllString(path, "/")
}
