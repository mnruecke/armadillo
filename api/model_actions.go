package api

import (
	"fmt"
	"github.com/repp/armadillo/model"
	"net/http"
	"reflect"
)

type Api struct {
	// NOTES FROM 6/13:
	// Create an Api with a gateway, all route handler functions become functions of 'api' so as to ensure access
	// to Gateway and any other config variables as needed.
	Gateway DbGateway
}

func (a *Api) GenerateCreate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		//a.Gateway.Create(model)
		rw.Write([]byte(fmt.Sprintf("Creatd a new: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateFind(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateFindAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateUpdate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateUpdateAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateReplace(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Replace one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateDestroy(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateDestroyAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func (a *Api) GenerateInfo(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Info for: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}
