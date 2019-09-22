package structutil

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func GetValue(target interface{}, path string) reflect.Value {
	keys := strings.Split(path, ".")
	v := reflect.ValueOf(target)

	for _, k := range keys {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		v = v.FieldByName(k)
	}

	return v
}

func GetValueByJson(target interface{}, path string) reflect.Value {
	keys := strings.Split(path, ".")
	v := reflect.ValueOf(target)

	for _, k := range keys {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			tag := ParseJsonTag(v.Type().Field(i).Tag.Get("json"))

			if tag == k {
				v = field
				break
			}
		}
	}

	return v
}

func Set(target interface{}, path string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	v := GetValue(target, path)
	v.Set(reflect.ValueOf(value))

	return nil
}

func SetByJson(target interface{}, path string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	v := GetValueByJson(target, path)
	v.Set(reflect.ValueOf(value))

	return nil
}

func ParseJsonTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}

	return tag
}

func Get(target interface{}, path string) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	return GetValue(target, path).Interface(), nil
}

func GetByJsonTag(target interface{}, path string) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	return GetValueByJson(target, path).Interface(), nil
}
