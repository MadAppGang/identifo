package model

import "reflect"

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
