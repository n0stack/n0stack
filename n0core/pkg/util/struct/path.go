package structutil

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func GetValue(target reflect.Value, path string) reflect.Value {
	keys := strings.Split(path, ".")

	for _, k := range keys {
		if target.Kind() == reflect.Ptr {
			target = target.Elem()
		}

		target = target.FieldByName(k)
	}

	return target
}

func GetValueByJson(target reflect.Value, path string) reflect.Value {
	keys := strings.Split(path, ".")

	for _, k := range keys {
		if target.Kind() == reflect.Ptr {
			target = target.Elem()
		}

		for i := 0; i < target.NumField(); i++ {
			field := target.Field(i)
			tag := ParseJsonTag(target.Type().Field(i).Tag.Get("json"))

			if tag == k {
				target = field
				break
			}
		}
	}

	return target
}

func Set(target interface{}, path string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return errors.Errorf("target must be ptr")
	}

	v := GetValue(targetValue, path)
	v.Set(reflect.ValueOf(value))

	return nil
}

func SetByJson(target interface{}, path string, value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return errors.Errorf("target must be ptr")
	}

	v := GetValueByJson(targetValue, path)
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

	targetValue := reflect.ValueOf(target)
	return GetValue(targetValue, path).Interface(), nil
}

func GetByJsonTag(target interface{}, path string) (res interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	targetValue := reflect.ValueOf(target)
	return GetValueByJson(targetValue, path).Interface(), nil
}

func UpdateWithMaskUsingJson(target interface{}, source interface{}, paths []string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	targetValue := reflect.ValueOf(target)
	sourceValue := reflect.ValueOf(source)

	if targetValue.Kind() != reflect.Ptr {
		return errors.Errorf("target must be ptr")
	}

	for {
		if sourceValue.Kind() == reflect.Ptr {
			sourceValue = sourceValue.Elem()
		} else {
			break
		}
	}

	if targetValue.Elem().Type() != sourceValue.Type() {
		return errors.Errorf("target and source is not same type: target=%s, source=%s", targetValue.Elem().Type().String(), sourceValue.Type().String())
	}

	for _, p := range paths {
		v, err := GetByJsonTag(source, p)
		if err != nil {
			return err
		}

		if err := SetByJson(target, p, v); err != nil {
			return err
		}
	}

	return nil
}
