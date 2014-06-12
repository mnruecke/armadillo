package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"strings"
	"unicode"
)

type MongoGateway struct {
	DialInfo    *mgo.DialInfo
	BaseSession *mgo.Session
}

func (g *MongoGateway) Initialize() error {
	var connectionErr error
	g.BaseSession, connectionErr = mgo.DialWithInfo(g.DialInfo)
	return connectionErr
}

func (g *MongoGateway) NewSession() *mgo.Session {
	return g.BaseSession.Clone()
}

func (g *MongoGateway) Create(m Model) error {
	s := g.NewSession()
	defer s.Close()
	ensureValidId(m)
	return collection(s, m).Insert(m)
}

func (g *MongoGateway) Save(m Model) error {
	s := g.NewSession()
	defer s.Close()
	ensureValidId(m)
	_, err := collection(s, m).UpsertId(m.GetId(), m)
	return err
}

func (g *MongoGateway) FindBy(m Model, q Query) error {
	s := g.NewSession()
	defer s.Close()
	query := queryToMgoQuery(collection(s, m), q)
	query.One(m)
	return nil
}

func (g *MongoGateway) FindById(m Model) error {
	s := g.NewSession()
	defer s.Close()
	return collection(s, m).FindId(m.GetId()).One(m)
}

func (g *MongoGateway) FindAll(m Model) (interface{}, error) {
	s := g.NewSession()
	defer s.Close()
	// Create a pointer to an empty slice of whatever type "m" is
	var results = reflect.New(reflect.SliceOf(reflect.TypeOf(m))).Interface()
	err := collection(s, m).Find(nil).All(results)
	return results, err
}

func (g *MongoGateway) FindAllBy(m Model, q Query) (interface{}, error) {
	s := g.NewSession()
	defer s.Close()
	// Create a pointer to an empty slice of whatever type "m" is
	var results = reflect.New(reflect.SliceOf(reflect.TypeOf(m))).Interface()
	query := queryToMgoQuery(collection(s, m), q)
	query.All(results)
	return results, nil
}

func (g *MongoGateway) Update(m Model, updates map[string]interface{}) error {
	s := g.NewSession()
	defer s.Close()
	err := collection(s, m).UpdateId(m.GetId(), updates)
	if err != nil { return err }
	return g.FindById(m) // Todo: another trip to the db!? uggh. find a better way to map the fields in 'updates' to model
}

func (g *MongoGateway) UpdateAll(m Model, updates map[string]interface{}) (int, error) {
	return g.UpdateAllWhere(m, updates, nil)
}

func (g *MongoGateway) UpdateAllWhere(m Model, updates map[string]interface{}, conditions map[string]interface{}) (int, error) {
	s := g.NewSession()
	defer s.Close()
	changeInfo, err := collection(s, m).UpdateAll(conditions, updates)
	if err != nil {
		return 0, err
	}
	return changeInfo.Updated, nil
}

func (g *MongoGateway) DeleteWhere(m Model, conditions map[string]interface{}) error {
	s := g.NewSession()
	defer s.Close()
	return collection(s,m).Remove(conditions)
}

func (g *MongoGateway) DeleteById(m Model) error {
	s := g.NewSession()
	defer s.Close()
	return collection(s,m).RemoveId(m.GetId())
}

func (g *MongoGateway) DeleteAll(m Model) (int, error) {
	return g.DeleteAllWhere(m, nil)
}

func (g *MongoGateway) DeleteAllWhere(m Model, conditions map[string]interface{}) (int, error) {
	s := g.NewSession()
	defer s.Close()
	changeInfo, err := collection(s,m).RemoveAll(conditions)
	if err != nil {
		return 0, err
	}
	return changeInfo.Removed, nil
}

// Private methods

func queryToMgoQuery(c *mgo.Collection, q Query) *mgo.Query {
	query := c.Find(q.Conditions)
	if q.Order != nil {
		query = query.Sort(q.Order...)
	}
	if q.Offset != nil {
		query = query.Skip(*q.Offset)
	}
	if q.Limit != nil {
		query = query.Limit(*q.Limit)
	}
	return query
}

func ensureValidId(m Model) {
	objectIdStr := m.GetId().(bson.ObjectId).Hex()
	if !bson.IsObjectIdHex(objectIdStr) {
		m.SetId(bson.NewObjectId())
	}
}

func collection(s *mgo.Session, m Model) *mgo.Collection {
	return s.DB("").C(collectionName(m))
}

// Is this a good idea? Convention vs Configuration? TODO: decide.
func collectionName(m Model) string {
	modelTypeStr := reflect.TypeOf(m).String()
	starsRemoved := strings.Replace(modelTypeStr, "*", "", -1)
	structName := starsRemoved[(strings.LastIndex(starsRemoved, ".") + 1):]

	singularSnakeCaseName := ""
	for i := 0; i < len(structName); i++ {
		if unicode.IsUpper(rune(structName[i])) && i != 0 {
			singularSnakeCaseName += "_"
		}
		singularSnakeCaseName += string(structName[i])
	}

	singularSnakeCaseName += "s" // TODO: use (or write) and inflector to handle things like "Box" or "Deer" or "Class"
	return strings.ToLower(singularSnakeCaseName)
}
