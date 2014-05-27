package main

import (
	"github.com/repp/armadillo/server"
	"{{}}/config"
)

func main() {
	server.Run(config.Load())
}
