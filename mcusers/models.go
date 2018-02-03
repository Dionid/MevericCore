package tztusers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"tztatom/tztcoremgo"

	"time"
	"mevericcore/mcmongo"
)

//easyjson:json
type UserModel struct {
	mcmongo.ModelBase `bson:",inline"`
	Login  string  `bson:",omitempty"`
	Email                *string  `bson:",omitempty"`
	Password             string  `json:",omitempty"`
	IsAdmin              bool    `json:",omitempty" bson:"isAdmin"`
	Phone                *string `json:",omitempty"`
}

func (this *UserModel) EnsureIndex(collection *mgo.Collection) error {
	index := mgo.Index{
		Key:      []string{"email", "login"},
		Unique:   true,
		DropDups: true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return nil
}

func (this *UserModel) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (this *UserModel) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(password))
	return err == nil
}

func (this *UserModel) Validate(collection *mgo.Collection) []error {
	if this.Email == nil || this.Password == "" {
		return []error{errors.New("email and password must be specified")}
	}
	return nil
}

func (this *UserModel) BeforeInsert(collection *mgo.Collection) error {
	hashedPass, err := this.HashPassword(this.Password)
	if err != nil {
		return err
	}
	this.Password = hashedPass
	this.IsAdmin = false
	return nil
}

//easyjson:json
type UsersListModel []UserModel

//easyjson:json
type CompanyModel struct {
	tztcoremgo.ModelBase `bson:",inline"`
	Name                 string          `bson:"name"`
	EmployeesIds         []bson.ObjectId `json:"employeesIds" bson:"employeesIds,omitempty"`
	Employees            *UsersListModel `json:"-" bson:"-"`
}

func (this *CompanyModel) EnsureIndex(collection *mgo.Collection) error {
	index := mgo.Index{
		Key:      []string{"name", "employeesIds"},
		Unique:   true,
		DropDups: true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return nil
}

func (this *CompanyModel) Validate(collection *mgo.Collection) []error {
	if len(this.EmployeesIds) == 0 {
		return []error{errors.New("no emploees")}
	}
	return nil
}

