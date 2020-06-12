package coll_test

import (
	"testing"
	"time"

	"github.com/mylxsw/coll"
	"github.com/stretchr/testify/assert"
)

type Struct1 struct {
	Name string
}

type Struct2 struct {
	Name  string
	Count int
}

var sources = []Struct1{
	{Name: "name1"},
	{Name: "name2"},
}

func TestMap(t *testing.T) {
	var dest []Struct2
	err := coll.Map(sources, &dest, func(s1 Struct1) Struct2 {
		return Struct2{
			Name:  s1.Name,
			Count: 100,
		}
	})
	if err != nil {
		t.Error(err)
	}

	if len(sources) != len(dest) {
		t.Error("test failed")
	}

	if sources[0].Name != dest[0].Name {
		t.Error("test failed")
	}
}

func TestFilter(t *testing.T) {
	var dest []Struct1
	err := coll.Filter(sources, &dest, func(s1 Struct1) bool {
		return s1.Name == "name2"
	})
	if err != nil {
		t.Error(err)
	}

	if len(dest) != 1 {
		t.Error("test failed")
	}

	if dest[0].Name != "name2" {
		t.Error("test failed")
	}
}


type UserDomain struct {
	ID           int64
	Name         string
	Email        string
	Password     string
	Gender       string
	Roles        []Role
	CreatedAt    time.Time
	privateField string
}

type Role struct {
	ID   int64
	Name string
}

type UserView struct {
	ID            int64
	Name          string
	Email         string
	Gender        bool
	Roles         []Role
	HasPermission bool
	privateField  string
}

func TestCopyProperties(t *testing.T) {

	{
		userDomain := UserDomain{
			ID:       1,
			Name:     "Tom",
			Email:    "tom@aicode.cc",
			Password: "2356565656",
			Gender:   "MALE",
			Roles: []Role{
				{ID: 2, Name: "admin"},
				{ID: 3, Name: "viewer"},
			},
			CreatedAt:    time.Now(),
			privateField: "145",
		}

		var targetUserView UserView
		assert.NoError(t, coll.CopyProperties(userDomain, &targetUserView))
		assert.Equal(t, userDomain.ID, targetUserView.ID)
		assert.Equal(t, userDomain.Name, targetUserView.Name)
		assert.Equal(t, userDomain.Email, targetUserView.Email)
		assert.Equal(t, false, targetUserView.Gender)
		assert.Equal(t, false, targetUserView.HasPermission)
		assert.Equal(t, userDomain.Roles, targetUserView.Roles)
		assert.Equal(t, "", targetUserView.privateField)
	}

	{
		userDomain := &UserDomain{
			ID:        1,
			Name:      "Tom",
			Email:     "tom@aicode.cc",
			Password:  "2356565656",
			Gender:    "MALE",
			CreatedAt: time.Now(),
		}

		var targetUserView UserView
		assert.NoError(t, coll.CopyProperties(userDomain, &targetUserView))
		assert.Equal(t, userDomain.ID, targetUserView.ID)
		assert.Equal(t, userDomain.Name, targetUserView.Name)
		assert.Equal(t, userDomain.Email, targetUserView.Email)
		assert.Equal(t, false, targetUserView.Gender)
		assert.Equal(t, false, targetUserView.HasPermission)
	}

	{
		userDomain := UserDomain{
			ID:        1,
			Name:      "Tom",
			Email:     "tom@aicode.cc",
			Password:  "2356565656",
			Gender:    "MALE",
			CreatedAt: time.Now(),
		}

		var targetUserView UserView
		assert.Equal(t, coll.ErrTargetInvalid, coll.CopyProperties(userDomain, targetUserView))
		assert.Equal(t, coll.ErrTargetIsNil, coll.CopyProperties(userDomain, nil))
	}

}