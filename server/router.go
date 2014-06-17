package server

import (
	"fmt"
	"github.com/repp/armadillo/api"
	"github.com/repp/armadillo/model"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	C api.ModelMethodConstructor
	Routes map[string]map[string]httpHandler // Path, Method, Handler
}

type Route struct {
	Method  string
	Path    string
	Handler httpHandler
}

type ModelRoute struct {
	Method           string
	Path             string
	HandlerGenerator func(model.Model) func(http.ResponseWriter, *http.Request)
	ModelName        string
	ModelInstance    model.Model
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

func (r *Router) Model(publicName string, modelInstance model.Model, methodRules MethodRules) {
	allowedMethods := r.allowedModelMethods(methodRules)
	for key, _ := range allowedMethods {
		r.ModelRoute(key, publicName, modelInstance)
	}
}

func (r *Router) Create(publicName string, modelInstance model.Model) {
	r.ModelRoute("Create", publicName, modelInstance)
}

func (r *Router) Find(publicName string, modelInstance model.Model) {
	r.ModelRoute("Find", publicName, modelInstance)
}

func (r *Router) FindAll(publicName string, modelInstance model.Model) {
	r.ModelRoute("FindAll", publicName, modelInstance)
}

func (r *Router) Update(publicName string, modelInstance model.Model) {
	r.ModelRoute("Update", publicName, modelInstance)
}

func (r *Router) UpdateAll(publicName string, modelInstance model.Model) {
	r.ModelRoute("UpdateAll", publicName, modelInstance)
}

func (r *Router) Replace(publicName string, modelInstance model.Model) {
	r.ModelRoute("Replace", publicName, modelInstance)
}

func (r *Router) Destroy(publicName string, modelInstance model.Model) {
	r.ModelRoute("Destroy", publicName, modelInstance)
}

func (r *Router) DestroyAll(publicName string, modelInstance model.Model) {
	r.ModelRoute("DestroyAll", publicName, modelInstance)
}

func (r *Router) Info(publicName string, modelInstance model.Model) {
	r.ModelRoute("Info", publicName, modelInstance)
}

func (r *Router) Action(path string, handler httpHandler) {
	r.appendRoute("POST", fmt.Sprintf("/{{.api_prefix}}/{{ .action_prefix }}/%v", path), handler)
}

func (r *Router) Get(path string, handler httpHandler) {
	r.appendRoute("GET", path, handler)
}

func (r *Router) Head(path string, handler httpHandler) {
	r.appendRoute("HEAD", path, handler)
}

func (r *Router) Post(path string, handler httpHandler) {
	r.appendRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler httpHandler) {
	r.appendRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler httpHandler) {
	r.appendRoute("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler httpHandler) {
	r.appendRoute("DELETE", path, handler)
}

func (r *Router) Options(path string, handler httpHandler) {
	r.appendRoute("OPTIONS", path, handler)
}

func (r *Router) Trace(path string, handler httpHandler) {
	r.appendRoute("TRACE", path, handler)
}

func (r *Router) ModelRoute(modelRouteName string, modelName string, modelInstance model.Model) {
	modelRoute := r.modelMethods(modelRouteName).(ModelRoute)
	modelRoute.ModelInstance = modelInstance
	modelRoute.ModelName = modelName
	route := convertToRoute(modelRoute)
	r.appendRoute(route.Method, route.Path, route.Handler)
}

func (r *Router) modelMethods(name string) interface{} {
	allModelMethods := map[string]ModelRoute{
		"Create":     ModelRoute{Method: "POST", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: r.C.GenerateCreate},
		"Find":       ModelRoute{Method: "GET", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: r.C.GenerateFind},
		"FindAll":    ModelRoute{Method: "GET", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: r.C.GenerateFindAll},
		"Update":     ModelRoute{Method: "PATCH", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: r.C.GenerateUpdate},
		"UpdateAll":  ModelRoute{Method: "PATCH", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: r.C.GenerateUpdateAll},
		"Replace":    ModelRoute{Method: "PUT", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: r.C.GenerateReplace},
		"Destroy":    ModelRoute{Method: "DELETE", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: r.C.GenerateDestroy},
		"DestroyAll": ModelRoute{Method: "DELETE", Path: "//{{.api_prefix}}/{{.model_name}}", HandlerGenerator: r.C.GenerateDestroyAll},
		"Info":       ModelRoute{Method: "OPTIONS", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: r.C.GenerateInfo},
	}
	if name == "" {
		return allModelMethods
	}
	if route, present := allModelMethods[name]; present {
		return route
	}
	panic(fmt.Sprintf(`Bad Route Name - "%v"`, name))
}

func (r *Router) appendRoute(method string, path string, handler httpHandler) {
	path = safeFormatPath(path)
	if r.Routes == nil {
		r.Routes = make(map[string]map[string]httpHandler)
	}
	if r.Routes[path] == nil {
		r.Routes[path] = make(map[string]httpHandler)
	}
	r.Routes[path][method] = handler
}

func (r *Router) allowedModelMethods(rules MethodRules) (allowedMethods map[string]ModelRoute) {
	// TODO: this is sketchy, leftover from an older architecture, clean it up!
	allModelMethods := r.modelMethods("").(map[string]ModelRoute)
	if len(rules.Allow) > 0 && len(rules.Forbid) > 0 {
		panic(fmt.Sprintf("Abort: Invalid ModelRules: %v \n Allow or Forbid cannot both contain rules.", rules))
	}

	allowedMethods = make(map[string]ModelRoute)

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

func convertToRoute(mr ModelRoute) Route {
	var handler httpHandler
	if mr.ModelInstance != nil {
		handler = mr.HandlerGenerator(mr.ModelInstance)
	}
	return Route{
		mr.Method,
		strings.Replace(mr.Path, "{{.model_name}}", mr.ModelName, 1),
		handler,
	}
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
