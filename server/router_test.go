package server

import (
	"github.com/repp/armadillo/model"
	"github.com/repp/armadillo/test"
	"net/http"
	"testing"
)

var mockHandler httpHandler = func(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello World!"))
	rw.WriteHeader(200)
}

func TestModel(t *testing.T) {
	var router Router
	mock := model.Model(&test.MockModel{7, "Test"})
	router.Model("cat", mock, MethodRules{Allow: []string{"Create"}})

	test.AssertEqual(t, len(router.ModelRoutes), 1)
	createRoute := router.ModelRoutes[0]
	test.AssertEqual(t, createRoute.Action, "Create")
	test.AssertEqual(t, createRoute.ModelName, "cat")
}

func TestCreate(t *testing.T) {
	var router Router
	//Create(publicName string, modelInstance model.Model) {
	router.Create("mocks", &test.MockModel{})

	createRoute := router.ModelRoutes[0]
	test.AssertEqual(t, createRoute.Action, "Create")
	test.AssertEqual(t, createRoute.ModelName, "mocks")
}

// No tests for other model methods (ie Find()) as they're just wrappers for appendModelRoute()

func TestGet(t *testing.T) {
	var router Router
	router.Get("/test/", mockHandler)

	getRoute := router.Routes[0]
	test.AssertEqual(t, getRoute.Method, "GET")
	test.AssertEqual(t, getRoute.Path, "/test/")
}

func TestAppendModelRoute(t *testing.T) {
	var router Router
	router.appendModelRoute("FindAll", "mocks", &test.MockModel{})

	findAllRoute := router.ModelRoutes[0]
	test.AssertEqual(t, findAllRoute.Action, "FindAll")
	test.AssertEqual(t, findAllRoute.ModelName, "mocks")
}

func TestAppendRoute(t *testing.T) {
	var router Router
	router.appendRoute("POST", "/dogs", mockHandler)

	route := router.Routes[0]
	test.AssertEqual(t, route.Method, "POST")
	test.AssertEqual(t, route.Path, "/dogs/")
}

// No tests for other methods (ie Post()) as they're just wrappers for appendRoute()

func TestSubtractSlice(t *testing.T) {
	a := []string{"fish", "turkey", "walnut"}
	b := []string{"turkey", "fish", "donkey"}
	c := subtractSlice(a, b)
	test.AssertDeepEqual(t, c, []string{"walnut"})
}

func TestAllowedModelMethods(t *testing.T) {
	var r Router
	ruleSet1 := MethodRules{}
	allowedMethods1 := r.allowedModelMethods(ruleSet1)
	// Order isn't ensured so a deep equivalency test isn't possible
	test.AssertEqual(t, len(allowedMethods1), 9) // TODO: remove hard coded number, get list from api

	ruleSet2 := MethodRules{Allow: []string{"Create"}}
	allowedMethods2 := r.allowedModelMethods(ruleSet2)
	test.AssertEqual(t, allowedMethods2[0], "Create")

	ruleSet3 := MethodRules{Forbid: []string{"Destroy"}}
	allowedMethods3 := r.allowedModelMethods(ruleSet3)
	test.AssertEqual(t, len(allowedMethods3), 8)
}

func TestSafeFormatPath(t *testing.T) {
	test.AssertEqual(t, safeFormatPath("dogs"), "/dogs/")
	test.AssertEqual(t, safeFormatPath("//dogs"), "/dogs/")
	test.AssertEqual(t, safeFormatPath("//dogs//"), "/dogs/")
	test.AssertEqual(t, safeFormatPath("dog//more/cat"), "/dog/more/cat/")
	test.AssertEqual(t, safeFormatPath("//////dog///more/cat"), "/dog/more/cat/")
	test.AssertEqual(t, safeFormatPath("//////dog///more//cat////"), "/dog/more/cat/")
	test.AssertEqual(t, safeFormatPath(""), "/")
}
