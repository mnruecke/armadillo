package config

import "github.com/repp/armadillo/server"

var development = server.Config {
	"port": 3000,
}

var production = server.Config {
	"port": 3000,
}

func Load() server.Config {
	if false { // ENV var based
		return production
	}
	return development
}

func Routes() {

}
