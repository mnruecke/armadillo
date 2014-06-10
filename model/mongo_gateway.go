package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"strings"
	"unicode"
)

type MongoGateway struct {
	DialInfo *mgo.DialInfo
	BaseSession *mgo.Session
}

func (g *MongoGateway) Initialize() error {
	var connectionErr error
	g.BaseSession, connectionErr = mgo.DialWithInfo(g.DialInfo)
	return connectionErr
}

func (g *MongoGateway) NewSession() *mgo.Session {
	return  g.BaseSession.Clone()
}

func (g *MongoGateway) Create(m Model) error {
	s := g.NewSession()
	defer s.Close()
	ensureValidId(m)
	return s.DB("").C(collectionName(m)).Insert(m)
}

func (g *MongoGateway) Save(m Model) error {
	s := g.NewSession()
	defer s.Close()
	ensureValidId(m)
	_, err := s.DB("").C(collectionName(m)).UpsertId(m.GetId(), m)
	return err
}

func (g *MongoGateway) FindBy(m Model, q Query) error {
	return nil
}

func (g *MongoGateway) FindById(m Model) error {
	return nil
}

func (g *MongoGateway) FindAll(m Model) ([]Model, error) {
	return []Model{}, nil
}

func (g *MongoGateway) FindAllBy(m Model, q Query) ([]Model, error) {
	return []Model{}, nil
}

func (g *MongoGateway) Update(m Model) error {
	return nil
}

func (g *MongoGateway) UpdateAll(m Model) (int, error) {
	return 0, nil
}

func (g *MongoGateway) UpdateAllWith(m Model, q Query) (int, error) {
	return 0, nil
}

func (g *MongoGateway) DeleteBy(m Model, q Query) (int, error) {
	return 0, nil
}

func (g *MongoGateway) DeleteById(m Model) error {
	return nil
}

func (g *MongoGateway) DeleteAll(m Model) (int, error) {
	return 0, nil
}

func ensureValidId(m Model) {
	objectIdStr := m.GetId().(bson.ObjectId).Hex()
	if !bson.IsObjectIdHex(objectIdStr) {
		m.SetId(bson.NewObjectId())
	}
}

// Is this a good idea? Convention vs Configuration? TODO: decide.
func collectionName(m Model) string {
	modelTypeStr := reflect.TypeOf(m).String()
	starsRemoved := strings.Replace(modelTypeStr, "*", "", -1)
	structName := starsRemoved[(strings.LastIndex(starsRemoved, ".")+1):]

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
