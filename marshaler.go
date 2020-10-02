package rac

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Marshaler interface {
	Marshaler() ([]byte, error)
}

func Marshal(value interface{}) ([]byte, error) {

	rv := reflect.ValueOf(value)

	rType := reflect.TypeOf(value)
	if _, ok := rType.(reflect.Type); !ok {
		rType = rType.Elem()
	}

	var v reflect.Value
	//var elemType reflect.Type

	v = rv.Elem()

	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint

	print(isSlice)

	return []byte{}, nil
}

func MarshalObject(value reflect.Value) ([]byte, error) {

	var arrStr []string

	rt := value.Type()
	ri := reflect.Indirect(value)

	for i := 0; i < ri.NumField(); i++ {

		fieldName := NameMapping(rt.Field(i).Name)

		tag := rt.Field(i).Tag.Get(TagNamespace)

		tags := strings.Split(tag, ",")

		if len(tags) > 0 && len(tags[0]) > 0 {
			fieldName = tags[0]
		}

		if tags[0] == "-" {
			continue
		}

		str := fieldName + " : "
		value := ri.Field(i).Interface()

		switch v := value.(type) {
		case Marshaler:
			data, _ := v.Marshaler()
			str += string(data)
		case int:
			str += fmt.Sprintf("%d", v)
		case time.Time:
			if v.IsZero() {
				str += ""
			} else {
				str += v.Format("2006-01-02T15:04:05")
			}
		case string:
			str += v
		}

		arrStr = append(arrStr, str)

	}

	b := strings.Join(arrStr, "\r\n")

	return []byte(b), nil
}
