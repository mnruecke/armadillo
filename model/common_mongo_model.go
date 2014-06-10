package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type CommonMongoModel struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Created modelTime     `json:"created"`
	Updated modelTime     `json:"updated"`
}

func (cm *CommonMongoModel) Initialize() {
	currentTime := time.Now()
	cm.SetId(bson.NewObjectId())
	cm.SetCreated(currentTime)
	cm.SetUpdated(currentTime)
}

// Getters and Setters
func (cm *CommonMongoModel) GetId() interface{}     { return cm.Id }
func (cm *CommonMongoModel) GetCreated() time.Time  { return cm.Created.format().Time }
func (cm *CommonMongoModel) GetUpdated() time.Time  { return cm.Updated.format().Time }
func (cm *CommonMongoModel) SetId(id interface{})   { cm.Id = id.(bson.ObjectId) }
func (cm *CommonMongoModel) SetCreated(t time.Time) { cm.Created = modelTime{t}.format() }
func (cm *CommonMongoModel) SetUpdated(t time.Time) { cm.Updated = modelTime{t}.format() }
