package config

import "github.com/repp/armadillo/server"

var development = server.Config{
	"port":               3000,
	"serve_static_files": true,
}

var production = server.Config{
	"port":               3000,
	"serve_static_files": false,
}

func Load() (config server.Config) {
	if false { // ENV var based
		config = production
	} else {
		config = development
	}

	// If a router isn't specified in the chosen config, use the default created by Routes()
	if _, routerPresent := config["router"]; !routerPresent {
		config["router"] = Routes()
	}
	return
}
