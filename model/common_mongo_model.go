package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// Aliases let us use modelTime's time attribute as an inline field in MongoDb
type created modelTime
type updated modelTime

type CommonMongoModel struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Created created     `json:"created" bson:",inline"`
	Updated updated     `json:"updated" bson:",inline"`
}

func (cm *CommonMongoModel) Initialize(currentTime time.Time) {
	cm.Id = bson.NewObjectId()
	cm.SetCreated(currentTime)
	cm.SetUpdated(currentTime)
}

// Getters and Setters
func (cm *CommonMongoModel) GetId() interface{}     { return cm.Id }
func (cm *CommonMongoModel) GetCreated() time.Time  { return modelTime(cm.Created).format().Time }
func (cm *CommonMongoModel) GetUpdated() time.Time  { return modelTime(cm.Updated).format().Time }
func (cm *CommonMongoModel) SetId(id interface{})   { cm.Id = id.(bson.ObjectId) }
func (cm *CommonMongoModel) SetCreated(t time.Time) { cm.Created = created(modelTime{t}.format()) }
func (cm *CommonMongoModel) SetUpdated(t time.Time) { cm.Updated = updated(modelTime{t}.format()) }
