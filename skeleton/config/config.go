package config

import "github.com/repp/armadillo/server"

var development = server.Config{
	"port":               3000,
	"serve_static_files": true,
	"api_prefix":         "/api",
	"action_prefix":      "/actions",
}

var production = server.Config{
	"port":               3000,
	"serve_static_files": false,
	"api_prefix":         "/api",
	"action_prefix":      "/actions",
}

func Load() (config server.Config) {
	if false { // ENV var based
		config = production
	} else {
		config = development
	}
	return
}
