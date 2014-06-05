package model

import (
	"github.com/repp/armadillo/test"
	"testing"
)

func TestMongoGatewayImplementsDbGateway(t *testing.T) {
	gateway := DbGateway(&MongoGateway{})
	test.AssertNotEqual(t, gateway, nil)
}
