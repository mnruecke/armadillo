package config

import (
	"github.com/repp/armadillo/model"
	"github.com/repp/armadillo/server"
)

var development = server.Config{
	"port":               3000,
	"serve_static_files": true,
	"api_prefix":         "api",
	"action_prefix":      "actions",
	"db": &model.MongoGateway{
		Address:  "localhost",
		Database: "armadillo_test",
	},
}

var production = server.Config{
	"port":               3000,
	"serve_static_files": false,
	"api_prefix":         "api",
	"action_prefix":      "actions",
	"db": &model.MongoGateway{
		Address:  "localhost",
		Database: "armadillo_production",
		Username: "user",
		Password: "pass",
	},
}

func Load() server.Config {
	if false { // Todo: Make ENV var based
		return production
	}
	return development
}
