package server

import (
	"testing"
)

func TestAppendRoute(t *testing.T) {
	var router Router
	router.appendNewRoute("POST", "/dogs", tmpFunc)

	_, present := router.Routes["/dogs"]["POST"]
	assertTrue(t, present)
}

func TestGet(t *testing.T) {
	var router Router
	router.Get("/test", tmpFunc)

	_, present1 := router.Routes["/test"]
	_, present2 := router.Routes["/test"]["GET"]
	assertTrue(t, present1)
	assertTrue(t, present2)
}

// No tests for other methods (ie Post()) as they're just wrappers for appendRoute()

func TestAllowedModelMethods(t *testing.T) {
	ruleSet1 := MethodRules{}
	allowedMethods1 := allowedModelMethods(ruleSet1)
	assertDeepEqual(t, allowedMethods1, allModelMethods)

	ruleSet2 := MethodRules{Allow: NewSet("Create")}
	allowedMethods2 := allowedModelMethods(ruleSet2)
	_, createPresent := allowedMethods2["Create"]
	_, deletePresent := allowedMethods2["Delete"]
	assertTrue(t, createPresent)
	assertFalse(t, deletePresent)

	ruleSet3 := MethodRules{Forbid: NewSet("Delete")}
	allowedMethods3 := allowedModelMethods(ruleSet3)
	_, createPresent2 := allowedMethods3["Create"]
	_, deletePresent2 := allowedMethods3["Delete"]
	assertTrue(t, createPresent2)
	assertFalse(t, deletePresent2)
}

func TestModel(t *testing.T) {
	var router Router
	type Cat struct{ legs int }
	router.Model("cat", Cat{1}, MethodRules{Allow: NewSet("Create")})

	assertEqual(t, len(router.Routes), 1)
	_, present := router.Routes["/cat"]["POST"]
	assertTrue(t, present)
}
