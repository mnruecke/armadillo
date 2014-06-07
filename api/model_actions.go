package api

import (
	"fmt"
	"github.com/repp/armadillo/model"
	"net/http"
	"reflect"
)

func GenerateCreate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Creatd a new: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateFind(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateFindAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Found all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateUpdate(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateUpdateAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Updated all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateReplace(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Replace one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateDestroy(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed one: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateDestroyAll(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Destroyed all: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}

func GenerateInfo(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("Info for: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}
