package model

import (
	"github.com/repp/armadillo/test"
	"testing"
	"labix.org/v2/mgo"
	"time"
)

var TestDialInfo *mgo.DialInfo = &mgo.DialInfo{
	Addrs: []string{"localhost"},
	Direct: true,
	Timeout: 10*time.Second,
	Database: "armadillo_test",
}
var TestGateway *MongoGateway = &MongoGateway{DialInfo: TestDialInfo}

func TestMongoGatewayImplementsDbGateway(t *testing.T) {
	gateway := DbGateway(&MongoGateway{})
	test.AssertNotEqual(t, gateway, nil)
}

func TestGatewayInitialize(t *testing.T) {
	err := TestGateway.Initialize()
	test.AssertEqual(t, err, nil)

	buildInfo, err := TestGateway.BaseSession.BuildInfo()
	test.AssertEqual(t, err, nil)
	test.AssertTypeMatch(t, buildInfo.Version, "2.6.whatever")
}

func TestNewSession(t *testing.T) {
	newSession := TestGateway.NewSession()
	test.AssertTypeMatch(t, newSession, &mgo.Session{})

	err := newSession.Ping()
	test.AssertEqual(t, err, nil)
}

func TestCollectionName(t *testing.T) {
	mock := &MockMongoModel{Name: "Ryan"}
	collectionName := collectionName(mock)

	test.AssertEqual(t, collectionName, "mock_mongo_models")
}

func TestCreate(t *testing.T) {
	resetDb()
	mock := &MockMongoModel{Name: "Ryan"}
	mock.Initialize()
	err := TestGateway.Create(mock)

	test.AssertEqual(t, err, nil)

	var foundMock MockMongoModel
	err = TestGateway.NewSession().DB("").C(collectionName(mock)).FindId(mock.GetId()).One(&foundMock)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, mock.GetId(), foundMock.GetId())
	test.AssertEqual(t, foundMock.Name, "Ryan")

	err = TestGateway.Create(mock)
	test.AssertNotEqual(t, err, nil) // Duplicate Index Key Error expected
}

func TestSave(t *testing.T) {
	resetDb()

	// New Record using save:
	mock := &MockMongoModel{Name: "Casscius Clay"}
	err := TestGateway.Save(mock)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, mock.Name, "Casscius Clay")
	test.AssertNotEqual(t, mock.GetId(), nil)

	var foundMock MockMongoModel
	err = TestGateway.NewSession().DB("").C(collectionName(mock)).FindId(mock.GetId()).One(&foundMock)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, mock.GetId(), foundMock.GetId())
	test.AssertEqual(t, foundMock.Name, "Casscius Clay")

	//Existing record using save
	mock.Name = "Muhammud Ali"
	err = TestGateway.Save(mock)
	test.AssertEqual(t, err, nil)

	err = TestGateway.NewSession().DB("").C(collectionName(mock)).FindId(mock.GetId()).One(&foundMock)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, foundMock.Name, "Muhammud Ali")
}

func TestFindBy(t *testing.T) {
	resetDb()

	// Setup
	mock1 := &MockMongoModel{Name: "Trevor"}
	mock2 := &MockMongoModel{Name: "Emma"}
	mock1.Initialize()
	mock2.Initialize()

	insertTestFixture(mock1)
	insertTestFixture(mock2)

	// Find Emma
	emma := &MockMongoModel{}
	err := TestGateway.FindBy(emma, Query{Conditions: map[string]interface{}{"name": "Emma"}})
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, emma.GetId(), mock2.GetId())
	test.AssertEqual(t, emma.Name, "Emma")

	// Find the second mock not named "Joyce" order by name, in this case "Trevor"
	trevor := &MockMongoModel{}
	offset := 1
	qry := Query{
		Conditions: map[string]interface{}{"name": map[string]string{"$ne": "Joyce"}},
		Order: []string{"name"},
		Offset: &offset,
	}
	err = TestGateway.FindBy(trevor, qry)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, trevor.GetId(), mock1.GetId())
	test.AssertEqual(t, trevor.Name, "Trevor")

}

func insertTestFixture(m Model) {
	TestGateway.NewSession().DB("").C(collectionName(m)).Insert(m)
}

// Resets the database so each test can be run in isolation
func resetDb() {
	err := TestGateway.NewSession().DB("").DropDatabase()
	if err != nil { panic(err) }
}
