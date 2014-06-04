package server

import (
	"testing"
)

func TestExtractPathFromTemplate(t *testing.T) {
	pathTemplate := "/{{.api_prefix}}/{{.model_name}}"
	config := Config{
		"api_prefix": "api",
		"model_name": "user",
	}
	path := extractPathFromTemplate(pathTemplate, config)
	assertEqual(t, path, "/api/user/")
}

