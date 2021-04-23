package coll

import (
	"errors"
	"reflect"
	"unsafe"
)

func Unique(origin interface{}, dest interface{}, uniquers ...interface{}) error {
	c := MustNew(origin)
	for _, m := range uniquers {
		c = c.Unique(m)
	}

	return c.All(dest)
}

func Map(origin interface{}, dest interface{}, mappers ...interface{}) error {
	c := MustNew(origin)
	for _, m := range mappers {
		c = c.Map(m)
	}

	return c.All(dest)
}

func Filter(origin interface{}, dest interface{}, filters ...interface{}) error {
	c := MustNew(origin)
	for _, f := range filters {
		c = c.Filter(f)
	}
	return c.All(dest)
}

var ErrTargetIsNil = errors.New("target is nil")
var ErrTargetInvalid = errors.New("target must be a pointer to struct")

// CopyProperties copy exported properties(with same name and type) from source to target
// target must be a pointer to struct
func CopyProperties(source interface{}, targets ...interface{}) error {
	sourceRefVal := reflect.Indirect(reflect.ValueOf(source))
	// 如果 source 为 null，则不需要拷贝任何属性
	if !sourceRefVal.IsValid() {
		return nil
	}

	for _, target := range targets {
		targetRefVal := reflect.ValueOf(target)

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

	}

	return nil
}
