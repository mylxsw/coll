package coll_test

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/mylxsw/coll"
)

var testMapData = map[string]string{
	"k1": "v1",
	"k2": "v2",
	"k3": "v3",
	"":   "v4",
	"k5": "",
	"xx": "yy",
}

type Element struct {
	ID     int
	Name   string
	Gender string
}

type Element2 struct {
	ID   int
	Name string
	Age  int
}

func TestInvalidTypeForMap(t *testing.T) {
	cc := coll.MustNew(testMapData)

	collectionWithoutEmpty := cc.Filter(func(value string) bool {
		return value != ""
	}).Filter(func(value string, key string) bool {
		return key != ""
	})

	collectionWithoutEmpty.Each(func(value, key string) {
		if value == "" || key == "" {
			t.Errorf("test failed: %s=>%s", key, value)
		}
	})

	if collectionWithoutEmpty.Size() != 4 || cc.Size() != 6 {
		t.Error("test failed")
	}

	if cc.IsEmpty() {
		t.Error("test failed")
	}

	if !cc.Filter(func(value string) bool {
		return false
	}).IsEmpty() {
		t.Error("test failed")
	}
}

func TestInvalidTypeForArray(t *testing.T) {
	_, err := coll.New("hello")
	if err == nil || err != coll.ErrorInvalidDataType {
		t.Errorf("test failed")
	}

	// coll := MustNew([]string{"hello", "world", "", "you", "are"})
	// fmt.Println(coll.Filter2(func(item string) bool {
	// 	return item != ""
	// }).ToString())

	coll := coll.MustNew([]string{"hello", "world", "", "you", "are"}).Filter(func(item string) bool {
		return item != ""
	})

	if coll.Size() != 4 {
		t.Error("test failed")
	}

	coll.Each(func(item string, index int) {
		if item == "" || index < 0 {
			t.Errorf("test failed: %d:%s", index, item)
		}
	})

	coll.Each(func(item string) {
		if item == "" {
			t.Errorf("test failed: %s", item)
		}
	})

	if coll.IsEmpty() {
		t.Error("test failed")
	}

	if !coll.Filter(func(item string) bool {
		return false
	}).IsEmpty() {
		t.Error("test failed")
	}
}

func TestStringMapCollection(t *testing.T) {
	coll := coll.MustNew(testMapData)
	coll = coll.Filter(func(_, key string) bool {
		return key != ""
	}).Filter(func(value string) bool {
		return value != ""
	}).Map(func(value, key string) string {
		return fmt.Sprintf("<%s(%s)>", value, key)
	})

	coll.Each(func(value, key string) {
		if !regexp.MustCompile(fmt.Sprintf(`^<\w+\(%s\)>$`, key)).MatchString(value) {
			t.Error("test failed")
		}
	})

	joinedValue := coll.Reduce(func(carry string, value, key string) string {
		if coll.MapIndex(key).(string) != value {
			t.Error("test failed")
		}
		return carry + " " + coll.MapIndex(key).(string)
	}, "value: ").(string)

	if len(joinedValue) <= len("value: ") {
		t.Error("test failed")
	}

	coll.Map(func(value string, key string) (string, string) {
		return value, key + "(modified)"
	}).Each(func(_, key string) {
		if !strings.HasSuffix(key, "(modified)") {
			t.Error("test failed")
		}
	})
}

func TestStringCollection(t *testing.T) {
	coll := coll.MustNew([]interface{}{"hello", "world", "", "you", "are"})
	coll = coll.Filter(func(item string) bool {
		return item != ""
	}).Map(func(item string) string {
		return "<" + item + ">"
	})

	if coll.ToString() != "[<hello> <world> <you> <are>]" {
		t.Errorf("test failed: ^%s$", coll.ToString())
	}

	res := coll.Reduce(func(carry string, item string) string {
		return fmt.Sprintf("%s->%s", carry, item)
	}, "")

	if res != "-><hello>-><world>-><you>-><are>" {
		t.Errorf("test failed: ^%s$", res)
	}
}

func TestComplexMapCollection(t *testing.T) {
	elements := map[string]Element{
		"one":   {ID: 1, Name: "hello"},
		"two":   {ID: 2, Name: "world"},
		"three": {ID: 3, Name: ""},
		"four":  {ID: 4, Name: "Tom"},
	}

	coll := coll.MustNew(elements)
	coll = coll.Filter(func(value Element, key string) bool {
		return value.Name != ""
	}).Map(func(value Element) Element2 {
		return Element2{
			ID:   value.ID,
			Name: value.Name,
			Age:  value.ID * 2,
		}
	})

	if coll.Size() != 3 {
		t.Errorf("test failed")
	}

	if !coll.MapHasIndex("one") {
		t.Error("test failed")
	}

	if coll.MapHasIndex("three") {
		t.Error("test failed")
	}

	coll.Each(func(value Element2, key string) {
		if key == "" {
			t.Error("test failed")
		}
	})
}

func TestComplexCollection(t *testing.T) {

	elements := []Element{
		{ID: 1, Name: "hello"},
		{ID: 2, Name: "world"},
		{ID: 3, Name: ""},
		{ID: 4, Name: "Tom"},
	}

	coll := coll.MustNew(elements)
	coll = coll.Filter(func(item Element) bool {
		return item.Name != ""
	}).Map(func(item Element) Element2 {
		return Element2{
			ID:   item.ID,
			Name: item.Name,
			Age:  item.ID * 2,
		}
	})

	if coll.ToString() != "[{1 hello 2} {2 world 4} {4 Tom 8}]" {
		t.Errorf("test failed: ^%s$", coll.ToString())
	}

	res := coll.Reduce(func(carry string, item Element2) string {
		return fmt.Sprintf("%v\n%v", carry, item)
	}, "{0 Start}")

	expectValue := `{0 Start}
{1 hello 2}
{2 world 4}
{4 Tom 8}`

	if res != expectValue {
		t.Errorf("test failed: ^%s$", res)
	}
}

func TestToArray(t *testing.T) {
	col := coll.MustNew([]Element{
		{ID: 11, Name: "guan"},
		{ID: 12, Name: "yi"},
		{ID: 13, Name: "yao"},
	})

	var elements []Element
	if err := col.All(&elements); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for _, v := range elements {
		fmt.Printf("type=%s, id=%d, name=%s\n", reflect.TypeOf(v).Name(), v.ID, v.Name)
	}

	var element2s []Element2
	if err := col.Map(func(ele Element) Element2 {
		return Element2{
			ID:   ele.ID,
			Name: ele.Name,
			Age:  ele.ID * 10,
		}
	}).All(&element2s); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for _, v := range element2s {
		fmt.Printf("type=%s, id=%d, name=%s, age=%d\n", reflect.TypeOf(v).Name(), v.ID, v.Name, v.Age)
	}
}

func TestToMap(t *testing.T) {
	col := coll.MustNew(map[string]Element{
		"guan": {ID: 11, Name: "guan"},
		"yi":   {ID: 12, Name: "yi"},
		"yao":  {ID: 13, Name: "yao"},
	})

	var elements map[string]Element
	if err := col.All(&elements); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for k, v := range elements {
		fmt.Printf("type=%s, k=%s, id=%d, name=%s\n", reflect.TypeOf(v).Name(), k, v.ID, v.Name)
	}

	var element2s map[string]Element2
	if err := col.Map(func(ele Element) Element2 {
		return Element2{
			ID:   ele.ID,
			Name: ele.Name,
			Age:  ele.ID * 10,
		}
	}).All(&element2s); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for k, v := range element2s {
		fmt.Printf("type=%s, k=%s, id=%d, name=%s, age=%d\n", reflect.TypeOf(v).Name(), k, v.ID, v.Name, v.Age)
	}
}

func TestPointerToArray(t *testing.T) {
	col := coll.MustNew([]*Element{
		{ID: 11, Name: "guan"},
		{ID: 12, Name: "yi"},
		{ID: 13, Name: "yao"},
	})

	var elements []*Element
	if err := col.All(&elements); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for _, v := range elements {
		fmt.Printf("type=%s, id=%d, name=%s\n", reflect.TypeOf(v).Kind(), v.ID, v.Name)
	}

	var element2s []*Element2
	if err := col.Map(func(ele *Element) *Element2 {
		return &Element2{
			ID:   ele.ID,
			Name: ele.Name,
			Age:  ele.ID * 10,
		}
	}).All(&element2s); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for _, v := range element2s {
		fmt.Printf("type=%s, id=%d, name=%s, age=%d\n", reflect.TypeOf(v).Kind(), v.ID, v.Name, v.Age)
	}
}

func TestPointerToMap(t *testing.T) {
	col := coll.MustNew(map[string]*Element{
		"guan": {ID: 11, Name: "guan"},
		"yi":   {ID: 12, Name: "yi"},
		"yao":  {ID: 13, Name: "yao"},
	})

	var elements map[string]*Element
	if err := col.All(&elements); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for k, v := range elements {
		fmt.Printf("type=%s, k=%s, id=%d, name=%s\n", reflect.TypeOf(v).Kind(), k, v.ID, v.Name)
	}

	var element2s map[string]*Element2
	if err := col.Map(func(ele *Element) *Element2 {
		return &Element2{
			ID:   ele.ID,
			Name: ele.Name,
			Age:  ele.ID * 10,
		}
	}).All(&element2s); err != nil {
		t.Errorf("test failed: %s", err)
	}

	for k, v := range element2s {
		fmt.Printf("type=%s, k=%s, id=%d, name=%s, age=%d\n", reflect.TypeOf(v).Kind(), k, v.ID, v.Name, v.Age)
	}
}

func TestGroupByForArray(t *testing.T) {
	col := coll.MustNew([]*Element{
		{ID: 11, Name: "guan", Gender: "boy"},
		{ID: 12, Name: "yi", Gender: "boy"},
		{ID: 13, Name: "yao", Gender: "girl"},
		{ID: 15, Name: "a", Gender: "girl"},
		{ID: 18, Name: "b", Gender: "boy"},
		{ID: 23, Name: "c", Gender: "girl"},
		{ID: 24, Name: "dd", Gender: "girl"},
		{ID: 26, Name: "e", Gender: "girl"},
	})

	var element2s map[interface{}][]interface{}
	if err := col.GroupBy(func(ele *Element) string {
		return ele.Gender
	}).All(&element2s); err != nil {
		t.Errorf("test failed: %v", err)
	}

	for k, v := range element2s {
		fmt.Printf("key=%v, value=%v\n", k, v)
	}
}
