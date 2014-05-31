package server

import (
	"../"
	"testing"
)

func TestValidateHandler(t *testing.T) {
	goodHandler := func() int { return 21 }
	badHandler := map[string]string{
		"bad": "worse",
	}
	type f struct{ name string }
	badHandler2 := f{"bob"}

	err1 := validateHandler(goodHandler)
	err2 := validateHandler(badHandler)
	err3 := validateHandler(badHandler2)

	main.AssertEqual(t, err1, nil)
	main.AssertNotEqual(t, err2, nil)
	main.AssertNotEqual(t, err3, nil)
}

func TestAppendRoute(t *testing.T) {
	var router Router
	router.appendRoute("POST", "/dogs", func() {})

	main.AssertEqual(t, router.Routes[0].HttpMethod, "POST")
	main.AssertEqual(t, router.Routes[0].Path, "/dogs")
}

func TestGet(t *testing.T) {
	var router Router
	routeHandler := func() int { return 1 }
	router.Get("/test", routeHandler)

	main.AssertEqual(t, router.Routes[0].HttpMethod, "GET")
	main.AssertEqual(t, router.Routes[0].Path, "/test")
}

// No tests for other methods (ie Post()) as they're just wrappers for appendRoute()

func TestAllowedModelMethods(t *testing.T) {
	ruleSet1 := MethodRules{}
	allowedMethods1 := allowedModelMethods(ruleSet1)
	main.AssertDeepEqual(t, allowedMethods1, allModelMethods)

	ruleSet2 := MethodRules{Allow: NewSet("Create")}
	allowedMethods2 := allowedModelMethods(ruleSet2)
	_, createPresent := allowedMethods2["Create"]
	_, deletePresent := allowedMethods2["Delete"]
	main.AssertTrue(t, createPresent)
	main.AssertFalse(t, deletePresent)

	ruleSet3 := MethodRules{Forbid: NewSet("Delete")}
	allowedMethods3 := allowedModelMethods(ruleSet3)
	_, createPresent2 := allowedMethods3["Create"]
	_, deletePresent2 := allowedMethods3["Delete"]
	main.AssertTrue(t, createPresent2)
	main.AssertFalse(t, deletePresent2)
}

func TestModel(t *testing.T) {
	var router Router
	type Cat struct{ legs int }
	router.Model("cat", Cat{1}, MethodRules{Allow: NewSet("Create")})
	createRoute := router.Routes[0]

	main.AssertEqual(t, len(router.Routes), 1)
	main.AssertEqual(t, createRoute.HttpMethod, "POST")
	main.AssertEqual(t, createRoute.Path, "/cat")
}
