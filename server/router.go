package server

import (
	"fmt"
	"reflect"
	"strings"
)

type Router struct {
	Routes []Route
}

type Route struct {
	HttpMethod string
	Path       string
	Handler    interface{}
}

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

var tmpFunc func() = func() {}

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
		r.Routes = append(r.Routes, route)
	}
}

func (r *Router) Action(path string, action interface{}) {
	r.appendRoute("POST", fmt.Sprintf("{{ .ActionPath }}/%v", path), action)
}

func (r *Router) Get(path string, action interface{}) {
	r.appendRoute("GET", path, action)
}

func (r *Router) Head(path string, action interface{}) {
	r.appendRoute("HEAD", path, action)
}

func (r *Router) Post(path string, action interface{}) {
	r.appendRoute("POST", path, action)
}

func (r *Router) Put(path string, action interface{}) {
	r.appendRoute("PUT", path, action)
}

func (r *Router) Patch(path string, action interface{}) {
	r.appendRoute("PATCH", path, action)
}

func (r *Router) Delete(path string, action interface{}) {
	r.appendRoute("DELETE", path, action)
}

func (r *Router) Options(path string, action interface{}) {
	r.appendRoute("OPTIONS", path, action)
}

func (r *Router) Trace(path string, action interface{}) {
	r.appendRoute("TRACE", path, action)
}

func (r *Router) appendRoute(method string, path string, action interface{}) {
	if err := validateHandler(action); err != nil {
		panic(err.Error())
	}
	r.Routes = append(r.Routes, Route{method, path, action})
}

func validateHandler(handler interface{}) (e error) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		e = &InvalidHandlerError{handler}
	}
	return
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

type InvalidHandlerError struct {
	handler interface{}
}

func (e *InvalidHandlerError) Error() string {
	return fmt.Sprintf(`"A handler must be a callable function. Got: '%s'"`, e.handler)
}
