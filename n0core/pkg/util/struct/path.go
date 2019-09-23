package structutil

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func GetValue(target reflect.Value, path string) (v reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	keys := strings.Split(path, ".")
	for _, k := range keys {
		if target.Kind() == reflect.Ptr {
			target = target.Elem()
		}

		target = target.FieldByName(k)
	}

	return target, nil
}

func GetValueByJson(target reflect.Value, path string) (v reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

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

	return target, nil
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

	v, err := GetValue(targetValue, path)
	if err != nil {
		return err
	}
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

	v, err := GetValueByJson(targetValue, path)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(value))

	return nil
}

func ParseJsonTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}

	return tag
}

func Get(target interface{}, path string) (interface{}, error) {
	targetValue := reflect.ValueOf(target)
	v, err := GetValue(targetValue, path)
	return v.Interface(), err
}

func GetByJsonTag(target interface{}, path string) (interface{}, error) {

	targetValue := reflect.ValueOf(target)
	v, err := GetValueByJson(targetValue, path)
	return v.Interface(), err
}

func UpdateWithMaskUsingJson(target interface{}, source interface{}, paths []string) error {
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
