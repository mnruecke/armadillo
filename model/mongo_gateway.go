package model

import "labix.org/v2/mgo"

type MongoGateway struct {
	DialInfo    *mgo.DialInfo
	BaseSession *mgo.Session
}

func (g *MongoGateway) Initialize() {
	var err error
	g.BaseSession, err = mgo.DialWithInfo(g.DialInfo)
	if err != nil {
		panic(err)
	}
}

func (g *MongoGateway) Create(m Model) error {
	return nil
}

func (g *MongoGateway) Save(m Model) error {
	return nil
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
