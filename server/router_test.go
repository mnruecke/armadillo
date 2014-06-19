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
	router.Model("cat", mock, MethodRules{Allow: NewSet("Create")})

	test.AssertEqual(t, len(router.ModelRoutes), 1)
	createRoute := router.ModelRoutes[0]
	test.AssertEqual(t, createRoute.Method, "POST")
	test.AssertEqual(t, createRoute.ModelName, "cat")
}

func TestCreate(t *testing.T) {
	var router Router
	//Create(publicName string, modelInstance model.Model) {
	router.Create("mocks", &test.MockModel{})

	createRoute := router.ModelRoutes[0]
	test.AssertEqual(t, createRoute.Method, "POST")
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
	test.AssertEqual(t, findAllRoute.Method, "GET")
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

func TestAllowedModelMethods(t *testing.T) {
	var r Router
	ruleSet1 := MethodRules{}
	allowedMethods1 := r.allowedModelMethods(ruleSet1)
	// Order isn't ensured so a deep equivalency test isn't possible
	test.AssertEqual(t, len(allowedMethods1), len(r.modelMethods("").(map[string]ModelRoute)))

	ruleSet2 := MethodRules{Allow: NewSet("Create")}
	allowedMethods2 := r.allowedModelMethods(ruleSet2)
	_, createPresent := allowedMethods2["Create"]
	_, deletePresent := allowedMethods2["Delete"]
	test.AssertTrue(t, createPresent)
	test.AssertFalse(t, deletePresent)

	ruleSet3 := MethodRules{Forbid: NewSet("Delete")}
	allowedMethods3 := r.allowedModelMethods(ruleSet3)
	_, createPresent2 := allowedMethods3["Create"]
	_, deletePresent2 := allowedMethods3["Delete"]
	test.AssertTrue(t, createPresent2)
	test.AssertFalse(t, deletePresent2)
}

func TestConvertToRoute(t *testing.T) {
	var r Router
	blankModelRoute := ModelRoute{}
	blankRoute := blankModelRoute.Route()
	test.AssertDeepEqual(t, blankRoute, Route{})

	noHandlerGeneratorRoute := ModelRoute{}
	noHandlerGeneratorRoute.ModelInstance = &test.MockModel{}
	noHandlerRoute := blankModelRoute.Route()
	test.AssertDeepEqual(t, noHandlerRoute, Route{})

	completeModelRoute := ModelRoute{"GET", "/{{.api_prefix}}/{{.model_name}}", r.C.GenerateFindAll, "tacos", &test.MockModel{}}
	completeRoute := completeModelRoute.Route()
	test.AssertEqual(t, completeRoute.Method, "GET")
	test.AssertEqual(t, completeRoute.Path, "/{{.api_prefix}}/tacos")
	ok := httpHandler(completeRoute.Handler)
	test.AssertNotEqual(t, ok, nil)
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
