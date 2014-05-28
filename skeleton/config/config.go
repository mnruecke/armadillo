package config

import "github.com/repp/armadillo/server"

var development = server.Config {
	"port": 3000,
	"serve_static_files": false,
}

var production = server.Config {
	"port": 3000,
	"serve_static_files": false,
}

func Load() server.Config {
	if false { // ENV var based
		return production
	}
	return development
}
