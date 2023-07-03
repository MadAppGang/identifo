package model

import (
	"errors"
	"reflect"
)

// CopyFields creates the new instance of src type
// and copies specific fields in 'fields' argument to the new value.
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

// CopyDstFields copy fields from src struct or pointer to struct
// to dst struct. dst Should be a pointer to a struct.
// Only the fields with the same name and type will be copied.
// If the type is mismatch - the field will be ignored.
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

// Filled returns the array of fields which are set in a struct.
// The empty fields are the fields with pointier types and and nil value in it. Zero fields are not empty.
// Fields which are structs are checked as well.
func Filled[T any](src T) []string {
	return filledNested("", src)
}

func filledNested[T any](prefix string, src T) []string {
	ur := reflect.ValueOf(src)

	// if the src is interface, get underlying value behind that.
	if ur.Kind() == reflect.Interface && !ur.IsNil() {
		elm := ur.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			ur = elm
		}
	}

	// if type is pointer - get a value referenced by a pointer
	// maybe do while? for pinter to pinter to pointer case?
	if ur.Kind() == reflect.Ptr {
		ur = ur.Elem()
	}

	fn := ur.NumField()
	result := []string{}
	if len(prefix) > 0 {
		prefix += "."
	}

	for i := 0; i < fn; i++ {
		fv := ur.Field(i)
		if fv.Kind() == reflect.Pointer && !fv.IsNil() {
			if fv.Elem().Kind() == reflect.Struct {
				result = append(result, filledNested(prefix+ur.Type().Field(i).Name, fv.Elem().Interface())...)
			} else {
				result = append(result, prefix+ur.Type().Field(i).Name)
			}
		}
		// else if fv.Kind() == reflect.Struct {
		//
		// }
	}

	return result
}
