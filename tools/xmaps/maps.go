package xmaps

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// FieldsToMap converts any struct to map[string]any.
func FieldsToMap(s any) map[string]any {
	return fieldsToMapNested("", s)
}

func fieldsToMapNested(prefix string, src any) map[string]any {
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
	if len(prefix) > 0 {
		prefix += "."
	}

	f := map[string]any{}
	for i := 0; i < fn; i++ {
		fv := ur.Field(i)

		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		if fv.Kind() == reflect.Struct {
			ff := fieldsToMapNested(prefix+ur.Type().Field(i).Name, fv.Interface())
			maps.Copy(f, ff)
		} else if fv.Kind() == reflect.Slice {
			if !fv.IsZero() {
				for j := 0; j < fv.Len(); j++ {
					pr := fmt.Sprintf("%s%s[%d]", prefix, ur.Type().Field(i).Name, j)
					ff := fieldsToMapNested(pr, fv.Index(j).Interface())
					maps.Copy(f, ff)
				}
			}
		} else {
			if !fv.IsZero() {
				f[prefix+ur.Type().Field(i).Name] = fv.Interface()
			}
		}
	}
	return f
}

// FilterMap returns new map only containing keys from the filter slice.
func FilterMap[T comparable, K any](m map[T]K, filter []T) map[T]K {
	result := map[T]K{}
	for k, v := range m {
		if slices.Contains(filter, k) {
			result[k] = v
		}
	}
	return result
}
