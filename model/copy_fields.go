package model

import (
	"errors"
	"reflect"
)

// CopyFields creates the new instance of src type and  copies specific fields to the new value.
func CopyFields[T any](src T, fields []string) T {
	rr := reflect.New(reflect.TypeOf(src))
	resultVal := rr.Elem()

	ur := reflect.ValueOf(&src).Elem()
	for _, f := range fields {
		fv := resultVal.FieldByName(f)
		if fv.IsValid() && fv.CanSet() {
			fv.Set(ur.FieldByName(f))
		}
	}
	return resultVal.Interface().(T)
}

func CopyDstFields[T any, K any](src T, dst K) error {
	res := reflect.ValueOf(dst)
	if res.Kind() != reflect.Ptr {
		return errors.New("destination is not a pointer")
	}
	resElement := res.Elem()
	if resElement.Kind() != reflect.Struct {
		return errors.New("destination is not a struct")
	}

	srcElement := reflect.ValueOf(src)
	if srcElement.Kind() == reflect.Ptr {
		srcElement = srcElement.Elem()
	}

	for i := 0; i < resElement.NumField(); i++ {
		dstField := resElement.Field(i)
		sfv := srcElement.FieldByName(resElement.Type().Field(i).Name)
		if dstField.IsValid() && dstField.CanSet() && dstField.Kind() == sfv.Kind() {
			dstField.Set(sfv)
		}
	}

	return nil
}