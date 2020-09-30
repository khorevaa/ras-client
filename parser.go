package rac

import (
	"log"
	"reflect"
	"strings"
)

const TagNamespace = "rac"

type object map[string]string

func (o object) Columns() []string {

	var col []string

	for s, _ := range o {
		col = append(col, s)
	}
	return col
}

func (o object) Scan(ptr map[string]interface{}) error {

	for key, value := range o {

		dest, ok := ptr[key]

		if !ok {
			continue
		}

		err := convertAssign(dest, value)

		if err != nil {
			log.Println("assing err: ", err)
		}
	}

	return nil
}

type decodeState struct {
	data    []byte
	objects []object
}

func (s *decodeState) init(data []byte) {

	s.data = data
	s.objects = []object{}

}

func (s *decodeState) parse() {

	parts := strings.Split(string(s.data), "\n")
	//parts = strings.Split(parts, "\n")

	part := ""
	for _, str := range parts {

		part += str + "\n"

		if len(str) == 0 {

			field := s.parseObject(part)

			if len(part) == 0 {
				continue
			}

			s.objects = append(s.objects, field)

			part = ""
		}

	}

}

func (s *decodeState) parseObject(in string) object {

	result := make(object)

	for _, line := range strings.Split(in, "\n") {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			result[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return result
}

func (s *decodeState) unmarshal(value interface{}) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(value)}
	}

	s.parse()

	rType := reflect.TypeOf(value)
	if _, ok := rType.(reflect.Type); !ok {
		rType = rType.Elem()
	}

	var v reflect.Value
	//var elemType reflect.Type

	v = rv.Elem()

	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint

	tagStore := newTagStore()

	for _, object := range s.objects {

		column := object.Columns()

		ptr := make(map[string]interface{}, len(column))

		var elem reflect.Value

		if isSlice {
			elem = reflectAlloc(v.Type().Elem())
		} else {
			elem = v
		}

		err := tagStore.findPtr(elem, column, ptr)
		if err != nil {
			return err
		}

		err = object.Scan(ptr)
		if err != nil {
			return err
		}
		for i := range ptr {
			ptr[i] = nil
		}

		if isSlice {
			v.Set(reflect.Append(v, elem))
		}

	}

	return nil
}

func reflectAlloc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem())
	}
	return reflect.New(typ).Elem()
}

func Unmarshal(data []byte, v interface{}) error {

	var d decodeState

	d.init(data)
	return d.unmarshal(v)
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "rac: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "rac: Unmarshal(nil " + e.Type.String() + ")"
}
