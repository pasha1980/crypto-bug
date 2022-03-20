package utils

import (
	"reflect"
)

func InArray(search interface{}, array interface{}) (exists bool, index interface{}) {
	val := reflect.ValueOf(array)
	val = val.Convert(val.Type())
	typ := reflect.TypeOf(array).Kind()

	switch typ {
	case reflect.Map:
		s := val.MapRange()

		for s.Next() {
			if reflect.DeepEqual(search, s.Value().Interface()) {
				return true, s.Key()
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(search, val.Index(i).Interface()) {
				return true, i
			}
		}
	}

	return false, nil
}
