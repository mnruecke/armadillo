package model

import (
	"github.com/repp/armadillo/test"
	"labix.org/v2/mgo"
	"testing"
	"time"
)

var TestDialInfo *mgo.DialInfo = &mgo.DialInfo{
	Addrs:    []string{"localhost"},
	Direct:   true,
	Timeout:  10 * time.Second,
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
		Order:      []string{"name"},
		Offset:     &offset,
	}
	err = TestGateway.FindBy(trevor, qry)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, trevor.GetId(), mock1.GetId())
	test.AssertEqual(t, trevor.Name, "Trevor")
}

func TestFindById(t *testing.T) {
	resetDb()

	mock := &MockMongoModel{Name: "Fuu"}
	mock.Initialize()
	insertTestFixture(mock)

	newMock := &MockMongoModel{}
	newMock.SetId(mock.GetId())
	err := TestGateway.FindById(newMock)

	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, newMock.GetCreated(), mock.GetCreated())
	test.AssertEqual(t, newMock.Name, "Fuu")
}

func TestFindAll(t *testing.T) {
	resetDb()

	thing1 := &MockMongoModel{Name: "Thing 1"}
	thing2 := &MockMongoModel{Name: "Thing 2"}
	thing1.Initialize()
	thing2.Initialize()
	insertTestFixture(thing1)
	insertTestFixture(thing2)

	modelsInterface, err := TestGateway.FindAll(&MockMongoModel{})
	models := *modelsInterface.(*[]*MockMongoModel)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 2)
	test.AssertTypeMatch(t, models, []*MockMongoModel{thing1, thing2})
}

func TestFindAllBy(t *testing.T) {
	resetDb()

	thing1 := &MockMongoModel{Name: "Thing"}
	thing2 := &MockMongoModel{Name: "Thing"}
	thing3 := &MockMongoModel{Name: "Bob"}
	thing1.Initialize()
	thing2.Initialize()
	thing3.Initialize()
	insertTestFixture(thing1)
	insertTestFixture(thing2)
	insertTestFixture(thing3)

	modelsInterface, err := TestGateway.FindAllBy(&MockMongoModel{}, Query{Conditions: map[string]interface{}{"name": "Thing"}})
	models := *modelsInterface.(*[]*MockMongoModel)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 2)
	test.AssertTypeMatch(t, models, []*MockMongoModel{thing1, thing2})
	// Todo: extend this test to hit order, offset and limit
}

func TestUpdate(t *testing.T) {
	resetDb()

	mock := &MockMongoModel{Name: "Robert Zimmerman"}
	mock.Initialize()
	insertTestFixture(mock)

	err := TestGateway.Update(mock, map[string]interface{}{"name": "Bob Dylan"})
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, mock.Name, "Bob Dylan")

	foundModel := findMockModel(mock.GetId())
	test.AssertEqual(t, foundModel.Name, "Bob Dylan")
}

func TestUpdateAll(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Neo"}
	mock2 := &MockMongoModel{Name: "Morpheus"}
	mock3 := &MockMongoModel{Name: "Trinity"}
	mock1.Initialize()
	mock2.Initialize()
	mock3.Initialize()
	insertTestFixture(mock1)
	insertTestFixture(mock2)
	insertTestFixture(mock3)

	updateCount, err := TestGateway.UpdateAll(&MockMongoModel{}, map[string]interface{}{"$set": map[string]string{"name": "Agent Smith"}})
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, updateCount, 3)

	foundModel := findMockModel(mock1.GetId())
	test.AssertEqual(t, foundModel.Name, "Agent Smith")
}

func TestUpdateAllWhere(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Neo"}
	mock2 := &MockMongoModel{Name: "Morpheus"}
	mock3 := &MockMongoModel{Name: "Trinity"}
	mock1.Initialize()
	mock2.Initialize()
	mock3.Initialize()
	insertTestFixture(mock1)
	insertTestFixture(mock2)
	insertTestFixture(mock3)

	updates := map[string]interface{}{"$set": map[string]string{"name": "Agent Smith"}}
	conditions := map[string]interface{}{"name": "Morpheus"}
	updateCount, err := TestGateway.UpdateAllWhere(&MockMongoModel{}, updates, conditions)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, updateCount, 1)

	foundModel1 := findMockModel(mock1.GetId())
	foundModel2 := findMockModel(mock2.GetId())
	test.AssertEqual(t, foundModel1.Name, "Neo")
	test.AssertEqual(t, foundModel2.Name, "Agent Smith")
}

func TestDeleteWhere(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Lemming"}
	mock2 := &MockMongoModel{Name: "Dodo"}
	mock3 := &MockMongoModel{Name: "House Cat"}
	mock1.Initialize()
	mock2.Initialize()
	mock3.Initialize()
	insertTestFixture(mock1)
	insertTestFixture(mock2)
	insertTestFixture(mock3)

	err := TestGateway.DeleteWhere(&MockMongoModel{}, map[string]interface{}{"name": "Dodo"})

	// Look for something that shouldn't be there
	var foundDodo MockMongoModel
	notFound := TestGateway.NewSession().DB("").C(collectionName(mock2)).FindId(mock2.GetId()).One(&foundDodo)
	lemming := findMockModel(mock1.GetId())
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, notFound.Error(), "not found")
	test.AssertNotEqual(t, lemming, nil)
}

func TestDeleteId(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Hitler"}
	mock1.Initialize()
	insertTestFixture(mock1)

	err := TestGateway.DeleteById(mock1)

	// Look for something that shouldn't be there
	var foundHitler MockMongoModel
	notFound := TestGateway.NewSession().DB("").C(collectionName(mock1)).FindId(mock1.GetId()).One(&foundHitler)
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, notFound, nil)
	test.AssertEqual(t, notFound.Error(), "not found")
}

func TestDeleteAll(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Lemming"}
	mock2 := &MockMongoModel{Name: "Dodo"}
	mock3 := &MockMongoModel{Name: "White Rhino"}
	mock1.Initialize()
	mock2.Initialize()
	mock3.Initialize()
	insertTestFixture(mock1)
	insertTestFixture(mock2)
	insertTestFixture(mock3)

	deleted, err := TestGateway.DeleteAll(&MockMongoModel{})
	var foundAnimal MockMongoModel
	notFound := TestGateway.NewSession().DB("").C(collectionName(mock1)).FindId(mock1.GetId()).One(&foundAnimal)

	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, deleted, 3)
	test.AssertNotEqual(t, notFound, nil)
	test.AssertEqual(t, notFound.Error(), "not found")
}

func TestDeleteAllWhere(t *testing.T) {
	resetDb()

	mock1 := &MockMongoModel{Name: "Lemming"}
	mock2 := &MockMongoModel{Name: "Dodo"}
	mock3 := &MockMongoModel{Name: "Dodo"}
	mock1.Initialize()
	mock2.Initialize()
	mock3.Initialize()
	insertTestFixture(mock1)
	insertTestFixture(mock2)
	insertTestFixture(mock3)

	deleted, err := TestGateway.DeleteAllWhere(&MockMongoModel{}, map[string]interface{}{"name": "Dodo"})
	var notFoundAnimal MockMongoModel
	notFound := TestGateway.NewSession().DB("").C(collectionName(mock2)).FindId(mock2.GetId()).One(&notFoundAnimal)

	var foundAnimal MockMongoModel
	foundErr := TestGateway.NewSession().DB("").C(collectionName(mock1)).FindId(mock1.GetId()).One(&foundAnimal)

	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, deleted, 2)
	test.AssertNotEqual(t, notFound, nil)
	test.AssertEqual(t, notFound.Error(), "not found")

	test.AssertEqual(t, foundErr, nil)
	test.AssertEqual(t, foundAnimal.Name, "Lemming")
}

// Helpers

func findMockModel(id interface{}) MockMongoModel {
	var foundMock MockMongoModel
	err := TestGateway.NewSession().DB("").C(collectionName(&MockMongoModel{})).FindId(id).One(&foundMock)
	if err != nil {
		panic(err)
	}
	return foundMock
}

func insertTestFixture(m Model) {
	TestGateway.NewSession().DB("").C(collectionName(m)).Insert(m)
}

// Resets the database so each test can be run in isolation
func resetDb() {
	err := TestGateway.NewSession().DB("").DropDatabase()
	if err != nil {
		panic(err)
	}
}
