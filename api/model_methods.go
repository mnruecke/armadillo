package api

import (
	"fmt"
	"github.com/repp/armadillo/model"
	"net/http"
	"reflect"
)

type ModelMethodConstructor struct {
	Gateway model.DbGateway
}

type ModelMethod struct {
	Method           string
	Path             string
	HandlerGenerator func(model.Model) func(http.ResponseWriter, *http.Request)
}

var MMC = ModelMethodConstructor{}
var AllModelMethods = map[string]ModelMethod{
	"Create":     ModelMethod{Method: "POST", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: MMC.GenerateCreate},
	"Find":       ModelMethod{Method: "GET", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: MMC.GenerateFind},
	"FindAll":    ModelMethod{Method: "GET", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: MMC.GenerateFindAll},
	"Update":     ModelMethod{Method: "PATCH", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: MMC.GenerateUpdate},
	"UpdateAll":  ModelMethod{Method: "PATCH", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: MMC.GenerateUpdateAll},
	"Replace":    ModelMethod{Method: "PUT", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: MMC.GenerateReplace},
	"Destroy":    ModelMethod{Method: "DELETE", Path: "/{{.api_prefix}}/{{.model_name}}/:id", HandlerGenerator: MMC.GenerateDestroy},
	"DestroyAll": ModelMethod{Method: "DELETE", Path: "//{{.api_prefix}}/{{.model_name}}", HandlerGenerator: MMC.GenerateDestroyAll},
	"Info":       ModelMethod{Method: "OPTIONS", Path: "/{{.api_prefix}}/{{.model_name}}", HandlerGenerator: MMC.GenerateInfo},
}

func ModelMethodNames() (names []string) {
	for key, _ := range AllModelMethods {
		names = append(names, key)
	}
	return
}

func (m *ModelMethodConstructor) GenerateCreate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		//a.Gateway.Create(model)
		rw.Write([]byte(fmt.Sprintf("Creatd a new: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateFind(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateFindAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateUpdate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateUpdateAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateReplace(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Replace one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateDestroy(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateDestroyAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (m *ModelMethodConstructor) GenerateInfo(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Info for: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}
