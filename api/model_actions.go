package api

import (
	"fmt"
	"github.com/repp/armadillo/model"
	"net/http"
	"reflect"
)

func CreateMethodFor(model model.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(fmt.Sprintf("MODEL: %v", reflect.TypeOf(model))))
		rw.WriteHeader(200)
	}
}
