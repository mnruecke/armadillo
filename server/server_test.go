package server

import (
	"github.com/repp/armadillo/test"
	"testing"
)

// TODO: write unit/integration tests for Run(), serveStaticFiles() and buildRoutes()

func TestExtractPathFromTemplate(t *testing.T) {
	pathTemplate := "/{{.api_prefix}}/{{.model_name}}"
	config := Config{
		"api_prefix": "api",
		"model_name": "user",
	}
	path := extractPathFromTemplate(pathTemplate, config)
	test.AssertEqual(t, path, "/api/user/")
}
