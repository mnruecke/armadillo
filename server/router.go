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
	Routes      []Route
	ModelRoutes []ModelRoute
}

type Route struct {
	Method  string
	Path    string
	Handler httpHandler
}

type ModelRoute struct {
	Action        string
	ModelName     string
	ModelInstance model.Model
}

type httpHandler func(http.ResponseWriter, *http.Request)

type MethodRules struct {
	Allow  []string
	Forbid []string
}

func (r *Router) Model(publicName string, modelInstance model.Model, methodRules MethodRules) {
	allowedMethods := r.allowedModelMethods(methodRules)
	for _, key := range allowedMethods {
		r.appendModelRoute(key, publicName, modelInstance)
	}
}

func (r *Router) allowedModelMethods(rules MethodRules) []string {
	allowLen := len(rules.Allow)
	forbidLen := len(rules.Forbid)
	allMethods := api.ModelMethodNames()

	// Catch bad MethodRules
	if allowLen > 0 && forbidLen > 0 {
		panic(fmt.Sprintf("Abort: Invalid ModelRules: %v \n Allow or Forbid cannot both contain rules.", rules))
	}

	// If Forbidden, get diff between all methods and forbidden
	if forbidLen > 0 {
		return subtractSlice(allMethods, rules.Forbid)
	}

	// If Allow, just return that.
	if allowLen > 0 {
		return rules.Allow
	}

	// both Allow and Forbid are empty, all methods are allowed.
	return allMethods
}

func subtractSlice(minuend, subtrahend []string) (diff []string) {
	var found bool
	for _, str1 := range minuend {
		found = false
		for _, str2 := range subtrahend {
			if str1 == str2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, str1)
		}
	}
	return
}

func (r *Router) appendRoute(method string, path string, handler httpHandler) {
	route := Route{method, safeFormatPath(path), handler}
	r.Routes = append(r.Routes, route)
}

func (r *Router) appendModelRoute(modelRouteName string, modelName string, modelInstance model.Model) {
	r.ModelRoutes = append(r.ModelRoutes, ModelRoute{modelRouteName, modelName, modelInstance})
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

// Route Creation Methods

func (r *Router) Create(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Create", publicName, modelInstance)
}

func (r *Router) Find(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Find", publicName, modelInstance)
}

func (r *Router) FindAll(publicName string, modelInstance model.Model) {
	r.appendModelRoute("FindAll", publicName, modelInstance)
}

func (r *Router) Update(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Update", publicName, modelInstance)
}

func (r *Router) UpdateAll(publicName string, modelInstance model.Model) {
	r.appendModelRoute("UpdateAll", publicName, modelInstance)
}

func (r *Router) Replace(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Replace", publicName, modelInstance)
}

func (r *Router) Destroy(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Destroy", publicName, modelInstance)
}

func (r *Router) DestroyAll(publicName string, modelInstance model.Model) {
	r.appendModelRoute("DestroyAll", publicName, modelInstance)
}

func (r *Router) Info(publicName string, modelInstance model.Model) {
	r.appendModelRoute("Info", publicName, modelInstance)
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
