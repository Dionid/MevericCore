package mcmongo

import (
	"gopkg.in/mgo.v2"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CollectionManagerBaseSt struct {
	DBName         string
	CollectionName string
	OriginSession  *mgo.Session
	Models         []ModelBaseInterface
	ErrNotFound    error
}

type CollectionManagerMgoBaseInterface interface {
	//Find(query interface{}) *mgo.Query
}

type CollectionManagerBaseInterface interface {
	CollectionManagerMgoBaseInterface
	// Collections and Sessions
	GetFullCustomSesAndCol(DBName string, collectionName string) (*mgo.Session, *mgo.Collection)
	GetCustomSesAndCol(collectionName string) (*mgo.Session, *mgo.Collection)
	GetSesAndCol() (*mgo.Session, *mgo.Collection)
	GetColBySes(session *mgo.Session) *mgo.Collection
	// API
	Destroy(query *bson.M) error
	DestroyByStringId(id string) error
	Delete(query *bson.M) error
	DeleteByStringId(id string) error
	// MODEL
	CreateModelBaseQuery(targetQuery *bson.M, model ModelBaseBaseInterface) *bson.M
	FindModel(targetQuery *bson.M, model ModelBaseInterface) error
	FindModelById(id bson.ObjectId, model ModelBaseInterface) error
	FindModelByStringId(id string, model ModelBaseInterface) error
	FindAllModels(targetQuery *bson.M, modelsList ModelsListBaseInterface) error
	InsertModelFullCustom(dbName string, colName string, model ModelBaseInterface) error
	InsertModelCustomCol(colName string, model ModelBaseInterface) error
	UpdateModelFullCustom(dbName string, colName string, model ModelBaseInterface, colQuerier bson.M, data bson.M) error
	UpdateModelCustomCol(colName string, model ModelBaseInterface, colQuerier bson.M, data bson.M) error
	InsertModel(model ModelBaseInterface) error
	UpdateModel(model ModelBaseInterface) error
	SaveModel(model ModelBaseInterface) error
	DestroyModel(model ModelBaseInterface) error
	DeleteModel(model ModelBaseInterface) error
	// INIT
	AddModel(model ModelBaseInterface)
	EnsureModelsIndexes() error
	InitManager(session *mgo.Session, dbName string, colName string)
}

//func (this *CollectionManagerBaseSt) Find(query interface{}) *mgo.Query {
//	session, col := this.GetSesAndCol()
//	defer session.Close()
//
//	return col.Find(query)
//}

func (this *CollectionManagerBaseSt) GetColBySes(session *mgo.Session) *mgo.Collection {
	col := session.DB(this.DBName).C(this.CollectionName)
	if col == nil {
		panic(errors.New("collection could not be created, maybe need to create it manually"))
	}
	return col
}

func (this *CollectionManagerBaseSt) GetFullCustomSesAndCol(DBName string, collectionName string) (*mgo.Session, *mgo.Collection) {
	ses := this.OriginSession.Copy()
	col := ses.DB(DBName).C(collectionName)
	if col == nil {
		panic(errors.New("collection could not be created, maybe need to create it manually"))
	}
	return ses, col
}

func (this *CollectionManagerBaseSt) GetCustomSesAndCol(collectionName string) (*mgo.Session, *mgo.Collection) {
	ses := this.OriginSession.Copy()
	col := ses.DB(this.DBName).C(collectionName)
	if col == nil {
		panic(errors.New("collection could not be created, maybe need to create it manually"))
	}
	return ses, col
}

func (this *CollectionManagerBaseSt) GetSesAndCol() (*mgo.Session, *mgo.Collection) {
	return this.GetCustomSesAndCol(this.CollectionName)
}

// SIMPLE API

func (this *CollectionManagerBaseSt) Destroy(query *bson.M) error {
	ses, col := this.GetSesAndCol()
	defer ses.Close()

	return col.Remove(query)
}

func (this *CollectionManagerBaseSt) DestroyByStringId(id string) error {
	return this.Destroy(&bson.M{"_id": bson.ObjectIdHex(id)})
}

func (this *CollectionManagerBaseSt) Delete(query *bson.M) error {
	ses, col := this.GetSesAndCol()
	defer ses.Close()

	if err := col.Update(query, &bson.M{"$set": &bson.M{"deleted_at": time.Now()}}); err != nil {
		return err
	}

	return nil
}

func (this *CollectionManagerBaseSt) DeleteByStringId(id string) error {
	return this.Delete(&bson.M{"_id": bson.ObjectIdHex(id)})
}

// MODEL

func (this *CollectionManagerBaseSt) CreateModelBaseQuery(targetQuery *bson.M, model ModelBaseBaseInterface) *bson.M {
	query := &bson.M{}
	baseQuery := model.GetBaseQuery()
	if baseQuery != nil {
		for k, v := range *targetQuery {
			(*baseQuery)[k] = v
		}
		query = baseQuery
	} else {
		query = targetQuery
	}
	return query
}

func (this *CollectionManagerBaseSt) FindModel(targetQuery *bson.M, model ModelBaseInterface) error {
	session, col := this.GetSesAndCol()
	defer session.Close()

	query := this.CreateModelBaseQuery(targetQuery, model)

	if err := col.Find(query).One(model); err != nil {
		return err
	}

	return nil
}

func (this *CollectionManagerBaseSt) FindModelById(id bson.ObjectId, model ModelBaseInterface) error {
	return this.FindModel(&bson.M{"_id": id}, model)
}

func (this *CollectionManagerBaseSt) FindModelByStringId(id string, model ModelBaseInterface) error {
	return this.FindModelById(bson.ObjectIdHex(id), model)
}

func (this *CollectionManagerBaseSt) FindAllModels(targetQuery *bson.M, modelsList ModelsListBaseInterface) error {
	session, col := this.GetSesAndCol()
	defer session.Close()

	query := this.CreateModelBaseQuery(targetQuery, modelsList)

	if err := col.Find(query).All(modelsList); err != nil {
		return err
	}

	return nil
}

func (this *CollectionManagerBaseSt) InsertModelFullCustom(dbName string, colName string, model ModelBaseInterface) error {
	ses, col := this.GetFullCustomSesAndCol(dbName, colName)
	defer ses.Close()

	if err := model.BeforeInsert(col); err != nil {
		return err
	}

	model.SetUpdatedAtNow()
	model.SetId(bson.NewObjectId())

	if err := col.Insert(model); err != nil {
		return err
	}

	if err := model.AfterInsert(col); err != nil {
		return err
	}

	return nil
}

func (this *CollectionManagerBaseSt) InsertModelCustomCol(colName string, model ModelBaseInterface) error {
	return this.InsertModelFullCustom(this.DBName, colName, model)
}

func (this *CollectionManagerBaseSt) UpdateModelFullCustom(dbName string, colName string, model ModelBaseInterface, colQuerier bson.M, data bson.M) error {
	ses, col := this.GetFullCustomSesAndCol(dbName, colName)
	defer ses.Close()

	if err := model.BeforeUpdate(col); err != nil {
		return err
	}

	model.SetUpdatedAtNow()

	if err := col.Update(colQuerier, &data); err != nil {
		return err
	}

	if err := model.AfterUpdate(col); err != nil {
		return err
	}

	return nil
}

func (this *CollectionManagerBaseSt) UpdateModelCustomCol(colName string, model ModelBaseInterface, colQuerier bson.M, data bson.M) error {
	return this.UpdateModelFullCustom(this.DBName, colName, model, colQuerier, data)
}

func (this *CollectionManagerBaseSt) InsertModel(model ModelBaseInterface) error {
	return this.InsertModelFullCustom(this.DBName, this.CollectionName, model)
}

func (this *CollectionManagerBaseSt) UpdateModel(model ModelBaseInterface) error {
	colQuerier := bson.M{"_id": model.GetId()}
	return this.UpdateModelFullCustom(this.DBName, this.CollectionName, model, colQuerier, bson.M{"$set": model})
}

func (this *CollectionManagerBaseSt) SaveModel(model ModelBaseInterface) error {
	if model.GetId() == "" {
		return this.InsertModel(model)
	} else {
		return this.UpdateModel(model)
	}
}

func (this *CollectionManagerBaseSt) DestroyModel(model ModelBaseInterface) error {
	return this.Destroy(&bson.M{"_id": model.GetId()})
}

func (this *CollectionManagerBaseSt) DeleteModel(model ModelBaseInterface) error {
	return this.Delete(&bson.M{"_id": model.GetId()})
}


// INIT
// For call EnsureModelsIndexes on every Model after manager initialization (InitManager)
func (this *CollectionManagerBaseSt) AddModel(model ModelBaseInterface) {
	this.Models = append(this.Models, model)
}

func (this *CollectionManagerBaseSt) EnsureModelsIndexes() error {
	ses, col := this.GetSesAndCol()
	defer ses.Close()
	for _, model := range this.Models {
		if err := model.EnsureIndex(col); err != nil {
			return err
		}
	}
	return nil
}

func (this *CollectionManagerBaseSt) InitManager(session *mgo.Session, dbName string, colName string) {
	this.ErrNotFound = mgo.ErrNotFound
	this.OriginSession = session
	this.DBName = dbName
	this.CollectionName = colName

	this.EnsureModelsIndexes()
}