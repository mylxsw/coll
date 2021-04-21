package coll

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

var (
	// ErrorInvalidDataType invalid data type error
	ErrorInvalidDataType = errors.New("invalid data type")
)

// DataType is collection data type
type DataType int8

const (
	// DataTypeMap represent the collection is a map collection
	DataTypeMap DataType = iota
	// DataTypeArrayOrSlice represent the collection is a array or slice collection
	DataTypeArrayOrSlice
)

// InvalidTypeError describes an invalid argument passed to To.
// (The argument to To must be a non-nil pointer.)
type InvalidTypeError struct {
	Type reflect.Type
}

func (e *InvalidTypeError) Error() string {
	if e.Type == nil {
		return "collection: To(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "collection: To(non-pointer " + e.Type.String() + ")"
	}
	return "collection: To(nil " + e.Type.String() + ")"
}

// Collection is a data collection
type Collection struct {
	dataArray []interface{}
	dataMap   map[interface{}]interface{}
	dataType  reflect.Type
}

// New create a new collection from data
func New(data interface{}) (*Collection, error) {
	collection := Collection{
		dataType: reflect.TypeOf(data),
	}

	dataKind := collection.dataType.Kind()
	if dataKind != reflect.Array && dataKind != reflect.Slice && dataKind != reflect.Map {
		return nil, ErrorInvalidDataType
	}

	dataValue := reflect.ValueOf(data)
	if dataKind == reflect.Map {
		dataMap := make(map[interface{}]interface{}, dataValue.Len())
		for _, key := range dataValue.MapKeys() {
			dataMap[key.Interface()] = dataValue.MapIndex(key).Interface()
		}

		collection.dataMap = dataMap
	} else {
		dataArray := make([]interface{}, dataValue.Len())
		for i := 0; i < dataValue.Len(); i++ {
			dataArray[i] = dataValue.Index(i).Interface()
		}

		collection.dataArray = dataArray
	}

	return &collection, nil
}

// MustNew create a new collection from data with error suppress
func MustNew(data interface{}) *Collection {
	res, err := New(data)
	if err != nil {
		panic(err.Error())
	}

	return res
}

// Sort sort collection with orderBy function, only support array collection
//     compareFunc(val1 interface{}, val2 interface{}) bool
func (collection *Collection) Sort(compareFunc interface{}) *Collection {
	if collection.isMapType() {
		panic("map not support sort")
	}

	if !IsFunction(compareFunc, []int{2, 1}) {
		panic("invalid callback function")
	}

	orderByFuncValue := reflect.ValueOf(compareFunc)
	if orderByFuncValue.Type().Out(0).Kind() != reflect.Bool {
		panic("the return type for compareFunc must be a bool value")
	}

	sortItems := make(sortStructs, len(collection.dataArray))
	for i, v := range collection.dataArray {
		sortItems[i] = sortStruct{
			Compare: func(v1, v2 interface{}) bool {
				arguments := []reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2)}
				return orderByFuncValue.Call(arguments)[0].Bool()
			},
			Value: v,
		}
	}

	sort.Sort(sortItems)

	results := make([]interface{}, len(sortItems))
	for i, v := range sortItems {
		results[i] = v.Value
	}

	return MustNew(results)
}

type sortStruct struct {
	Compare func(v1, v2 interface{}) bool
	Value   interface{}
}

type sortStructs []sortStruct

func (s sortStructs) Len() int {
	return len(s)
}

func (s sortStructs) Less(i, j int) bool {
	return s[i].Compare(s[i].Value, s[j].Value)
}

func (s sortStructs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// AsArray convert the collection to an array collection
// Remember: the underlying data type is changed to []interface{}
func (collection *Collection) AsArray() *Collection {
	if collection.isMapType() {
		results := make([]interface{}, len(collection.dataMap))
		i := 0
		for _, value := range collection.dataMap {
			results[i] = value
			i++
		}

		return MustNew(results)
	}

	results := make([]interface{}, len(collection.dataArray))
	for i, value := range collection.dataArray {
		results[i] = value
	}

	return MustNew(results)
}

// AsMap convert collection to a map collection
//     keyFunc(value interface{}) interface{}
//     keyFunc(value interface{}, key interface{}) interface{}
// Remember: the underlying data type is changed to map[interface{}]interface{}
func (collection *Collection) AsMap(keyFunc interface{}) *Collection {
	if !IsFunction(keyFunc, []int{1, 1}, []int{2, 1}) {
		panic("invalid callback function")
	}

	keyFuncValue := reflect.ValueOf(keyFunc)
	keyFuncType := keyFuncValue.Type()
	argumentCount := keyFuncType.NumIn()

	if collection.isMapType() {
		results := make(map[interface{}]interface{})
		for key, value := range collection.dataMap {
			arguments := []reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}
			uniqID := keyFuncValue.Call(arguments[0:argumentCount])[0].Interface()

			if _, ok := results[uniqID]; ok {
				continue
			}

			results[uniqID] = value
		}

		return MustNew(results)
	}

	results := make(map[interface{}]interface{})
	for index, item := range collection.dataArray {
		uniqID := keyFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:argumentCount])[0].Interface()

		if _, ok := results[uniqID]; ok {
			continue
		}

		results[uniqID] = item
	}

	return MustNew(results)
}

// Unique remove duplicated elements from collection
//     uniqFunc(value interface{}) interface{}
//     uniqFunc(value interface{}, key interface{}) interface{}
// Remember: the return collection type is map[interface{}]interface{} for map, []interface{} for array or slices
func (collection *Collection) Unique(uniqFunc interface{}) *Collection {
	if !IsFunction(uniqFunc, []int{1, 1}, []int{2, 1}) {
		panic("invalid callback function")
	}

	uniqFuncValue := reflect.ValueOf(uniqFunc)
	uniqFuncType := uniqFuncValue.Type()
	argumentCount := uniqFuncType.NumIn()

	if collection.isMapType() {
		results := make(map[interface{}]interface{})
		for key, value := range collection.dataMap {
			arguments := []reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}
			uniqID := uniqFuncValue.Call(arguments[0:argumentCount])[0].Interface()

			if _, ok := results[uniqID]; ok {
				continue
			}

			results[uniqID] = value
		}

		return MustNew(results)
	}

	results := make(map[interface{}]interface{})
	for index, item := range collection.dataArray {
		uniqID := uniqFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:argumentCount])[0].Interface()

		if _, ok := results[uniqID]; ok {
			continue
		}

		results[uniqID] = item
	}

	resultsArr := make([]interface{}, len(results))
	i := 0
	for _, r := range results {
		resultsArr[i] = r
		i++
	}

	return MustNew(resultsArr)
}

// GroupBy iterates over elements of collection, grouping all elements by specified conditions
// groupFunc(interface{}) interface{}
// groupFunc(interface{}, int) interface{}
// Remember: the result of All must be type: map[interface{}][]interface{}
func (collection *Collection) GroupBy(groupFunc interface{}) *Collection {
	if !IsFunction(groupFunc, []int{1, 1}, []int{2, 1}) {
		panic("invalid callback function")
	}

	groupFuncValue := reflect.ValueOf(groupFunc)
	groupFuncType := groupFuncValue.Type()
	argumentCount := groupFuncType.NumIn()

	if collection.isMapType() {
		results := make(map[interface{}][]interface{})
		for key, value := range collection.dataMap {
			arguments := []reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}
			groupID := groupFuncValue.Call(arguments[0:argumentCount])[0].Interface()

			if _, ok := results[groupID]; !ok {
				results[groupID] = make([]interface{}, 0)
			}

			results[groupID] = append(results[groupID], value)
		}

		return MustNew(results)
	}

	results := make(map[interface{}][]interface{})
	for index, item := range collection.dataArray {
		groupID := groupFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:argumentCount])[0].Interface()

		if _, ok := results[groupID]; !ok {
			results[groupID] = make([]interface{}, 0)
		}

		results[groupID] = append(results[groupID], item)
	}

	return MustNew(results)
}

// Filter iterates over elements of collection, return all element meet the needs
// filter(interface{}) bool
// filter(interface{}, int) bool
func (collection *Collection) Filter(filter interface{}) *Collection {
	if !IsFunction(filter, []int{1, 1}, []int{2, 1}) {
		panic("invalid callback function")
	}

	filterValue := reflect.ValueOf(filter)
	filterType := filterValue.Type()
	argumentCount := filterType.NumIn()

	// 返回值类型必须为bool
	if filterType.Out(0).Kind() != reflect.Bool {
		panic("return argument should be a boolean")
	}

	if collection.isMapType() {
		results := make(map[interface{}]interface{})
		for key, value := range collection.dataMap {
			arguments := []reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}
			if filterValue.Call(arguments[0:argumentCount])[0].Interface().(bool) {
				results[key] = value
			}
		}

		return MustNew(results)
	}
	results := make([]interface{}, 0)
	for index, item := range collection.dataArray {
		if filterValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:argumentCount])[0].Interface().(bool) {
			results = append(results, item)
		}
	}

	return MustNew(results)
}

// Map manipulates an iterate and transforms it to another type.
// mapFunc(value interface{}) interface{}
// mapFunc(value interface{}) (value interface{}, key interface{})
// mapFunc(value interface{}, key interface{}) interface{}
// mapFunc(value interface{}, key interface{}) (value interface{}, key interface{})
func (collection *Collection) Map(mapFunc interface{}) *Collection {
	if !IsFunction(mapFunc, []int{1, 1}, []int{2, 2}, []int{1, 2}, []int{2, 1}) {
		panic("invalid callback function")
	}

	mapFuncValue := reflect.ValueOf(mapFunc)
	mapFuncArgumentCount := mapFuncValue.Type().NumIn()

	if collection.isMapType() {
		results := make(map[interface{}]interface{}, collection.Size())
		for key, value := range collection.dataMap {
			var values []reflect.Value
			arguments := []reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}

			values = mapFuncValue.Call(arguments[0:mapFuncArgumentCount])

			if len(values) == 1 {
				results[key] = values[0].Interface()
			} else {
				results[values[1].Interface()] = values[0].Interface()
			}
		}

		return MustNew(results)
	}

	results := make([]interface{}, len(collection.dataArray))
	for index, item := range collection.dataArray {
		results[index] = mapFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:mapFuncArgumentCount])[0].Interface()
	}

	return MustNew(results)
}

// Reduce Iteratively reduce the array to a single value using a callback function
// reduceFunc(carry interface{}, item interface{}) interface{}
func (collection *Collection) Reduce(reduceFunc interface{}, initial interface{}) interface{} {
	if !IsFunction(reduceFunc, []int{2, 1}, []int{3, 1}) {
		panic("invalid callback function")
	}

	reduceFuncValue := reflect.ValueOf(reduceFunc)
	argumentsCount := reduceFuncValue.Type().NumIn()

	previous := initial
	if collection.isMapType() {
		for key, value := range collection.dataMap {
			arguments := []reflect.Value{reflect.ValueOf(previous), reflect.ValueOf(value), reflect.ValueOf(key)}
			previous = reduceFuncValue.Call(arguments[0:argumentsCount])[0].Interface()
		}
	} else {
		for index, item := range collection.dataArray {
			arguments := []reflect.Value{reflect.ValueOf(previous), reflect.ValueOf(item), reflect.ValueOf(index)}
			previous = reduceFuncValue.Call(arguments[0:argumentsCount])[0].Interface()
		}
	}

	return previous
}

// Items return all items in the collection
func (collection *Collection) Items() interface{} {
	if collection.isMapType() {
		return collection.dataMap
	}

	return collection.dataArray
}

// All get all of the items in the collection.
// the argument result must be a pointer to map or slice
func (collection *Collection) All(result interface{}) (err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	if collection.isMapType() {
		return collection.toMap(result)
	}

	return collection.toArray(result)
}

func (collection *Collection) toMap(result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		return fmt.Errorf("result argument must be a slice address")
	}

	mapv := reflect.MakeMap(resultv.Elem().Type())

	for k, v := range collection.dataMap {
		mapv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}

	resultv.Elem().Set(mapv)

	return nil
}

func (collection *Collection) ToArray() ([]interface{}, error) {
	if collection.isMapType() {
		return collection.AsArray().ToArray()
	}

	var res []interface{}
	if err := collection.toArray(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (collection *Collection) toArray(result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		return fmt.Errorf("result argument must be a slice address")
	}

	elementLen := len(collection.dataArray)
	slicev := reflect.MakeSlice(resultv.Elem().Type(), elementLen, elementLen)

	for i, v := range collection.dataArray {
		slicev.Index(i).Set(reflect.ValueOf(v))
	}

	resultv.Elem().Set(slicev)

	return nil
}

// Each Execute a callback over each item.
func (collection *Collection) Each(eachFunc interface{}) {
	if !IsFunction(eachFunc) {
		panic("invalid callback function")
	}

	eachFuncValue := reflect.ValueOf(eachFunc)
	eachFuncType := eachFuncValue.Type()
	argumentNums := eachFuncType.NumIn()
	if argumentNums == 0 {
		panic("invalid callback function")
	}

	if collection.isMapType() {
		for key, value := range collection.dataMap {
			eachFuncValue.Call([]reflect.Value{reflect.ValueOf(value), reflect.ValueOf(key)}[0:argumentNums])
		}
	} else {
		for index, item := range collection.dataArray {
			eachFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)}[0:argumentNums])
		}
	}
}

// DataType return the data type
func (collection *Collection) DataType() DataType {
	if collection.isMapType() {
		return DataTypeMap
	}

	return DataTypeArrayOrSlice
}

// IsEmpty Determine if the collection is empty or not.
func (collection *Collection) IsEmpty() bool {
	if collection.isMapType() {
		return len(collection.dataMap) == 0
	}
	return len(collection.dataArray) == 0
}

// ToString print the data element
func (collection *Collection) ToString() string {
	if collection.isMapType() {
		return fmt.Sprint(collection.dataMap)
	}
	return fmt.Sprint(collection.dataArray)
}

// Size count the number of items in the collection.
func (collection *Collection) Size() int {
	if collection.isMapType() {
		return len(collection.dataMap)
	}
	return len(collection.dataArray)
}

// Index Get an item from the collection by index.
func (collection *Collection) Index(index int) interface{} {
	if collection.isMapType() {
		return nil
	}

	if !collection.HasIndex(index) {
		return nil
	}

	return reflect.ValueOf(collection.dataArray).Index(index).Interface()
}

// MapIndex get an item from the collection by key
func (collection *Collection) MapIndex(key interface{}) interface{} {
	if !collection.isMapType() {
		return nil
	}

	value := reflect.ValueOf(collection.dataMap).MapIndex(reflect.ValueOf(key))
	if value.IsValid() {
		return value.Interface()
	}

	return nil
}

// MapHasIndex return whether the collection has a key
func (collection *Collection) MapHasIndex(key interface{}) bool {
	if !collection.isMapType() {
		return false
	}

	return reflect.ValueOf(collection.dataMap).MapIndex(reflect.ValueOf(key)).IsValid()
}

// HasIndex return whether the collection has a index
func (collection *Collection) HasIndex(index int) bool {
	if collection.isMapType() {
		return false
	}

	return index >= 0 && index < collection.Size()
}

func (collection *Collection) isMapType() bool {
	return collection.dataType.Kind() == reflect.Map
}

// IsFunction returns if the argument is a function.
func IsFunction(in interface{}, argumentCheck ...[]int) bool {
	funcType := reflect.TypeOf(in)
	if funcType.Kind() != reflect.Func {
		return false
	}

	if len(argumentCheck) == 0 {
		return true
	}

	for _, check := range argumentCheck {
		isValid := false
		if len(check) >= 1 && check[0] >= 0 {
			isValid = funcType.NumIn() == check[0]
		}

		if len(check) == 2 {
			isValid = funcType.NumOut() == check[1]
		}

		if isValid {
			return true
		}
	}

	return false
}
