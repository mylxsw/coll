package coll_test

import (
	"testing"

	"github.com/mylxsw/coll"
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
