package server

import (
	"github.com/repp/armadillo/model"
	"github.com/repp/armadillo/test"
	"testing"
)

// TODO: write unit/integration tests for Run(), serveStaticFiles() and buildRoutes()

var mockConfig = Config{
	"port":               3000,
	"serve_static_files": true,
	"api_prefix":         "api",
	"action_prefix":      "actions",
	"db": &model.MongoGateway{
		Address:  "localhost",
		Database: "armadillo_test",
	},
}

func TestConvertToRoutes(t *testing.T) {
	model := &test.MockModel{1, "Test Model"}
	mockModelRoutes := []ModelRoute{
		ModelRoute{"Create", "mocks", model},
		ModelRoute{"Find", "mocks", model},
		ModelRoute{"FindAll", "mocks", model},
	}
	routes := convertToRoutes(mockModelRoutes, mockConfig)
	test.AssertEqual(t, len(routes), len(mockModelRoutes))
	// Todo: this better
}

func TestExtractPathFromTemplate(t *testing.T) {
	pathTemplate := "/{{.api_prefix}}/{{.model_name}}"
	config := Config{
		"api_prefix": "api",
		"model_name": "user",
	}
	path := extractPathFromTemplate(pathTemplate, config)
	test.AssertEqual(t, path, "/api/user/")
}
