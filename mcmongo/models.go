package mcmongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"errors"
)

type ModelBase struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	UpdatedAt time.Time     `json:",omitempty" bson:"updatedAt,omitempty"`
	DeletedAt *time.Time    `json:",omitempty" bson:"deletedAt,omitempty"`
}

type ModelBaseBaseInterface interface {
	GetBaseQuery() *bson.M
}

type ModelBaseInterface interface {
	ModelBaseBaseInterface

	EnsureIndex(collection *mgo.Collection) error // REQUIRED

	Validate(collection *mgo.Collection) []error // MANUALLY

	BeforeInsert(collection *mgo.Collection) error
	AfterInsert(collection *mgo.Collection) error

	BeforeUpdate(collection *mgo.Collection) error
	AfterUpdate(collection *mgo.Collection) error

	BeforeDelete(collection *mgo.Collection) error
	AfterDelete(collection *mgo.Collection) error

	GetId() bson.ObjectId
	SetId(id bson.ObjectId)

	SetUpdatedAtNow()
}

func (this *ModelBase) EnsureIndex(collection *mgo.Collection) error {
	panic(errors.New("indexes must be specified"))
	return nil
}

func (this *ModelBase) Validate(collection *mgo.Collection) []error {
	return nil
}

func (this *ModelBase) BeforeInsert(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) AfterInsert(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) BeforeUpdate(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) AfterUpdate(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) BeforeDelete(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) AfterDelete(collection *mgo.Collection) error {
	return nil
}

func (this *ModelBase) GetId() bson.ObjectId {
	return this.ID
}

func (this *ModelBase) SetId(id bson.ObjectId) {
	this.ID = id
}

func (this *ModelBase) GetBaseQuery() *bson.M {
	return nil
}

func (this *ModelBase) SetUpdatedAtNow() {
	this.UpdatedAt = time.Now()
}

type ModelsListBase []ModelBase

func (this *ModelsListBase) GetBaseQuery() *bson.M {
	return nil
}

type ModelsListBaseInterface interface {
	ModelBaseBaseInterface
}