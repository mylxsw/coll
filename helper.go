package coll

import (
	"errors"
	"reflect"
	"unsafe"
)

func Unique(origin interface{}, dest interface{}, uniquer interface{}) error {
	return MustNew(origin).Unique(uniquer).All(dest)
}

func Map(origin interface{}, dest interface{}, mapper interface{}) error {
	return MustNew(origin).Map(mapper).All(dest)
}

func Filter(origin interface{}, dest interface{}, filter interface{}) error {
	return MustNew(origin).Filter(filter).All(dest)
}

var ErrTargetIsNil = errors.New("target is nil")
var ErrTargetInvalid = errors.New("target must be a pointer to struct")

// CopyProperties copy exported properties(with same name and type) from source to target
// target must be a pointer to struct
func CopyProperties(source interface{}, target interface{}) error {
	sourceRefVal := reflect.Indirect(reflect.ValueOf(source))
	targetRefVal := reflect.ValueOf(target)

	// 如果 source 为 null，则不需要拷贝任何属性
	if !sourceRefVal.IsValid() {
		return nil
	}

	if !targetRefVal.IsValid() {
		return ErrTargetIsNil
	}

	if targetRefVal.Kind() != reflect.Ptr {
		return ErrTargetInvalid
	}

	targetVal := targetRefVal.Elem()
	targetType := targetVal.Type()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldName := field.Name
		if fieldName[0] < 'A' || fieldName[0] > 'Z' {
			continue
		}

		dst := sourceRefVal.FieldByName(fieldName)
		if !dst.IsValid() || field.Type != dst.Type() {
			continue
		}

		reflect.NewAt(field.Type, unsafe.Pointer(targetVal.Field(i).UnsafeAddr())).Elem().Set(dst)
	}

	return nil
}
